import { ref, onBeforeUnmount, type Ref } from 'vue'
import { getStreamWspUrl, getStreamBaseUrl } from '../utils/streamUrls'

export function useWSPCanvasPlayer(canvasRef: Ref<HTMLCanvasElement | null>) {
  const isPlaying = ref(false), isConnecting = ref(false), error = ref<string | null>(null)
  let ws: WebSocket | null = null, decoder: any = null, isStopped = false
  let activeStreamId = '', activeIsHero = false, activeHasSub = true, activeTargetId = ''

  const stop = () => {
    isStopped = true
    if (activeTargetId) {
      const idToStop = activeTargetId
      activeTargetId = ''
      fetch(`${getStreamBaseUrl()}/api/v1/streams/${encodeURIComponent(idToStop)}/h264/stop`, { method: 'POST' }).catch(() => {})
    }
    if (ws) {
      ws.onmessage = null; ws.onerror = null; ws.onclose = null
      try { ws.close() } catch (_) {}
      ws = null
    }
    if (decoder) { try { decoder.close() } catch (_) {}; decoder = null }
    isPlaying.value = false; isConnecting.value = false
  }

  const initDecoder = () => {
    if (typeof window === 'undefined' || !(window as any).VideoDecoder) return null
    try {
      const dec = new (window as any).VideoDecoder({
        output: (frame: any) => {
          if (!canvasRef.value || isStopped) { frame.close(); return }
          if (canvasRef.value.width !== frame.displayWidth || canvasRef.value.height !== frame.displayHeight) {
            canvasRef.value.width = frame.displayWidth
            canvasRef.value.height = frame.displayHeight
          }
          const ctx = canvasRef.value.getContext('2d')
          ctx?.drawImage(frame, 0, 0, frame.displayWidth, frame.displayHeight)
          frame.close()
          isPlaying.value = true; isConnecting.value = false
        },
        error: (e: any) => { error.value = e?.message || 'Decode error' }
      })
      dec.configure({ codec: 'avc1.42E01E', optimizeForLatency: true })
      return dec
    } catch (e: any) {
      error.value = e?.message || 'VideoDecoder failed'
      return null
    }
  }

  const start = async (streamId: string, isHero = false, hasSub = true) => {
    stop()
    if (!canvasRef.value || !streamId) return
    isStopped = false; activeStreamId = streamId; activeIsHero = isHero; activeHasSub = hasSub
    const targetId = (!isHero && hasSub) ? `${streamId}_sub` : streamId
    activeTargetId = targetId
    isConnecting.value = true; error.value = null
    decoder = initDecoder()
    if (!decoder) { error.value = 'WebCodecs not supported'; isConnecting.value = false; return }

    try {
      await fetch(`${getStreamBaseUrl()}/api/v1/streams/${encodeURIComponent(targetId)}/h264/start`, { method: 'POST' }).catch(() => {})
      if (isStopped || activeTargetId !== targetId) return

      ws = new WebSocket(getStreamWspUrl(streamId, isHero, hasSub))
      ws.binaryType = 'arraybuffer'
      ws.onmessage = (event) => {
        if (isStopped) return
        if (typeof event.data === 'string') {
          try {
            const msg = JSON.parse(event.data)
            if (msg.error) { error.value = `WSP: ${msg.error}`; isConnecting.value = false }
          } catch (_) {}
          return
        }
        if (!(event.data instanceof ArrayBuffer) || !decoder || decoder.state !== 'configured') return
        const buf = event.data as ArrayBuffer
        if (buf.byteLength <= 9) return
        const view = new DataView(buf)
        const isKey = view.getUint8(0) === 0x01
        const ts = Number(view.getBigUint64(1))
        const nalu = new Uint8Array(buf, 9)
        try {
          decoder.decode(new (window as any).EncodedVideoChunk({
            type: isKey ? 'key' : 'delta', timestamp: ts, data: nalu
          }))
        } catch (_) {}
      }
      ws.onerror = () => { if (!isStopped) { error.value = 'WSP error'; isConnecting.value = false } }
      ws.onclose = () => {
        if (!isStopped && !document.hidden && activeStreamId) {
          setTimeout(() => { if (!isStopped) start(activeStreamId, activeIsHero, activeHasSub) }, 3000)
        }
      }
    } catch (err: any) {
      error.value = err?.message || 'WebSocket init failed'
      isConnecting.value = false
    }
  }

  const handleVis = () => {
    if (!document.hidden && activeStreamId && (!ws || ws.readyState === WebSocket.CLOSED)) {
      start(activeStreamId, activeIsHero, activeHasSub)
    }
  }
  document.addEventListener('visibilitychange', handleVis)
  window.addEventListener('focus', handleVis)
  onBeforeUnmount(() => {
    document.removeEventListener('visibilitychange', handleVis)
    window.removeEventListener('focus', handleVis)
    stop()
  })

  return { isPlaying, isConnecting, error, start, stop }
}
