import { ref, onBeforeUnmount, type Ref } from 'vue'
import { getStreamWspUrl, getStreamBaseUrl } from '../utils/streamUrls'

export function useWSPCanvasPlayer(canvasRef: Ref<HTMLCanvasElement | null>) {
  const isPlaying = ref(false), isConnecting = ref(false), error = ref<string | null>(null)
  let ws: WebSocket | null = null, decoder: any = null, isStopped = false
  let activeStreamId = '', activeIsHero = false, activeHasSub = true, activeTargetId = ''
  let hasKeyframe = false, pendingFrame: any = null, rafId: number | null = null, currentCodec = ''
  let isStarting = false

  const renderLoop = () => {
    if (isStopped) return
    if (pendingFrame && canvasRef.value) {
      const f = pendingFrame
      pendingFrame = null
      if (canvasRef.value.width !== f.displayWidth || canvasRef.value.height !== f.displayHeight) {
        canvasRef.value.width = f.displayWidth; canvasRef.value.height = f.displayHeight
      }
      canvasRef.value.getContext('2d')?.drawImage(f, 0, 0, f.displayWidth, f.displayHeight)
      f.close()
    }
    rafId = requestAnimationFrame(renderLoop)
  }

  const stop = () => {
    isStopped = true; isStarting = false; hasKeyframe = false; activeTargetId = ''; currentCodec = ''
    if (rafId) { cancelAnimationFrame(rafId); rafId = null }
    if (pendingFrame) { try { pendingFrame.close() } catch (_) {}; pendingFrame = null }
    if (ws) { ws.onmessage = null; ws.onerror = null; ws.onclose = null; try { ws.close() } catch (_) {}; ws = null }
    if (decoder) { try { decoder.close() } catch (_) {}; decoder = null }
    isPlaying.value = false; isConnecting.value = false
  }

  const initDecoder = () => {
    if (typeof window === 'undefined' || !(window as any).VideoDecoder) return null
    try {
      const codec = activeIsHero ? 'avc1.640033' : 'avc1.42c01f'
      const dec = new (window as any).VideoDecoder({
        output: (frame: any) => {
          if (!canvasRef.value || isStopped) { frame.close(); return }
          if (pendingFrame) pendingFrame.close()
          pendingFrame = frame
          isPlaying.value = true; isConnecting.value = false
        },
        error: (e: any) => {
          console.warn('[WebCodecs]', e?.message || e)
          error.value = e?.message || 'Decode error'
          try { dec.close() } catch (_) {}
          decoder = null
        }
      })
      dec.configure({ codec, optimizeForLatency: true })
      return dec
    } catch (e: any) { error.value = e?.message || 'VideoDecoder failed'; return null }
  }

  const start = async (streamId: string, isHero = false, hasSub = true) => {
    if (!canvasRef.value || !streamId) { stop(); return }
    const targetId = (!isHero && hasSub) ? `${streamId}_sub` : streamId
    if (isStarting && activeTargetId === targetId) return
    if (activeTargetId === targetId && activeIsHero === isHero && ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) && !isStopped) return
    stop()
    isStopped = false; isStarting = true; hasKeyframe = false; activeStreamId = streamId; activeIsHero = isHero; activeHasSub = hasSub
    activeTargetId = targetId; isConnecting.value = true; error.value = null
    decoder = initDecoder()
    if (!decoder) { error.value = 'WebCodecs not supported'; isConnecting.value = false; isStarting = false; return }
    if (!rafId) rafId = requestAnimationFrame(renderLoop)

    let frameSeq = 0
    try {
      await fetch(`${getStreamBaseUrl()}/api/v1/streams/${encodeURIComponent(targetId)}/h264/start`, { method: 'POST' }).catch(() => {})
      if (isStopped || activeTargetId !== targetId) return

      ws = new WebSocket(getStreamWspUrl(streamId, isHero, hasSub))
      ws.binaryType = 'arraybuffer'
      ws.onmessage = (event) => {
        if (isStopped || !(event.data instanceof ArrayBuffer) || event.data.byteLength <= 9) return
        const buf = event.data as ArrayBuffer, view = new DataView(buf), isKey = view.getUint8(0) === 0x01
        if (!isKey && !hasKeyframe) return
        hasKeyframe = true
        if (!decoder || decoder.state !== 'configured') {
          decoder = initDecoder()
          if (!decoder || !isKey) return
        }
        const ts = Math.round(performance.now() * 1000) + (frameSeq++)
        const nalu = new Uint8Array(buf, 9)
        try {
          decoder.decode(new (window as any).EncodedVideoChunk({
            type: isKey ? 'key' : 'delta', timestamp: ts, data: nalu
          }))
        } catch (_) {}
      }
      ws.onerror = () => { if (!isStopped) { error.value = 'WSP error'; isConnecting.value = false } }
      ws.onclose = () => {
        if (!isStopped && !document.hidden && activeStreamId) setTimeout(() => { if (!isStopped) start(activeStreamId, activeIsHero, activeHasSub) }, 1000)
      }
    } catch (err: any) { error.value = err?.message || 'WebSocket init failed'; isConnecting.value = false }
    finally { isStarting = false }
  }

  const handleVis = () => { if (!document.hidden && activeStreamId && (!ws || ws.readyState === WebSocket.CLOSED)) start(activeStreamId, activeIsHero, activeHasSub) }
  document.addEventListener('visibilitychange', handleVis)
  window.addEventListener('focus', handleVis)
  onBeforeUnmount(() => { document.removeEventListener('visibilitychange', handleVis); window.removeEventListener('focus', handleVis); stop() })

  return { isPlaying, isConnecting, error, start, stop }
}
