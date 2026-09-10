import { ref, onBeforeUnmount, type Ref } from 'vue'

export function useWebRTCPlayer(videoRef: Ref<HTMLVideoElement | null>) {
  const isPlaying = ref(false)
  const isConnecting = ref(false)
  const error = ref<string | null>(null)
  let pc: RTCPeerConnection | null = null
  let activeStreamId = ''
  let activeIsHero = false

  const stop = () => {
    if (pc) { pc.close(); pc = null }
    if (videoRef.value) { videoRef.value.srcObject = null }
    isPlaying.value = false
    isConnecting.value = false
  }

  const start = async (streamId: string, isHero = false) => {
    stop()
    if (!videoRef.value || !streamId) return
    activeStreamId = streamId
    activeIsHero = isHero
    isConnecting.value = true
    error.value = null

    try {
      pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
      pc.addTransceiver('video', { direction: 'recvonly' })
      pc.addTransceiver('audio', { direction: 'recvonly' })

      pc.ontrack = (event) => {
        if (videoRef.value && event.streams[0]) {
          videoRef.value.srcObject = event.streams[0]
          videoRef.value.play().catch(() => {})
          isPlaying.value = true
          isConnecting.value = false
        }
      }

      pc.onconnectionstatechange = () => {
        if (pc?.connectionState === 'failed' || pc?.connectionState === 'disconnected') {
          stop()
          if (!document.hidden) setTimeout(() => start(activeStreamId, activeIsHero), 1500)
        }
      }

      const offer = await pc.createOffer()
      await pc.setLocalDescription(offer)

      let path = isHero ? streamId : `${streamId}_sub`
      let res = await fetch(`http://localhost:8889/${path}/whep`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/sdp' },
        body: offer.sdp
      })

      // Fallback automático para Main se Sub-Stream não existir (404)
      if (!res.ok && !isHero && res.status === 404) {
        path = streamId
        res = await fetch(`http://localhost:8889/${path}/whep`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/sdp' },
          body: offer.sdp
        })
      }

      if (!res.ok) throw new Error(`WHEP status ${res.status}`)
      const answerSdp = await res.text()
      await pc.setRemoteDescription({ type: 'answer', sdp: answerSdp })
    } catch (err: any) {
      error.value = err.message || 'WHEP failed'
      isConnecting.value = false
      isPlaying.value = false
      if (!document.hidden) setTimeout(() => start(activeStreamId, activeIsHero), 3000)
    }
  }

  const handleResume = () => {
    if (document.hidden || !activeStreamId) return
    if (videoRef.value && videoRef.value.paused) videoRef.value.play().catch(() => {})
    if (!pc || pc.connectionState === 'failed' || pc.connectionState === 'disconnected' || !isPlaying.value) {
      start(activeStreamId, activeIsHero)
    }
  }

  document.addEventListener('visibilitychange', handleResume)
  window.addEventListener('focus', handleResume)

  onBeforeUnmount(() => {
    document.removeEventListener('visibilitychange', handleResume)
    window.removeEventListener('focus', handleResume)
    stop()
  })

  return { isPlaying, isConnecting, error, start, stop }
}
