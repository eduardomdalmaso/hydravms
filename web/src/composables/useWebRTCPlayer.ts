import { ref, onBeforeUnmount, type Ref } from 'vue'

export function useWebRTCPlayer(videoRef: Ref<HTMLVideoElement | null>) {
  const isPlaying = ref(false)
  const isConnecting = ref(false)
  const error = ref<string | null>(null)
  let pc: RTCPeerConnection | null = null
  let activeStreamId = ''
  let activeIsHero = false

  let activeHasSubStream = true

  const stop = () => {
    if (pc) { pc.close(); pc = null }
    if (videoRef.value) { videoRef.value.srcObject = null }
    isPlaying.value = false
    isConnecting.value = false
  }

  const start = async (streamId: string, isHero = false, hasSubStream = true) => {
    stop()
    if (!videoRef.value || !streamId) return
    activeStreamId = streamId
    activeIsHero = isHero
    activeHasSubStream = hasSubStream
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
          if (!document.hidden) setTimeout(() => start(activeStreamId, activeIsHero, activeHasSubStream), 1500)
        }
      }

      const offer = await pc.createOffer()
      await pc.setLocalDescription(offer)

      const shouldUseSub = !isHero && hasSubStream
      const path = shouldUseSub ? `${streamId}_sub` : streamId
      const endpoints = [
        `http://localhost:8889/${path}/whep`,
        `http://localhost:8889/${streamId}/whep`,
        `http://localhost:8080/whep/${path}`,
        `http://localhost:8080/api/v1/streams/${path}/whep`,
        `http://localhost:8080/whep/${streamId}`
      ]

      let answerSdp = ''
      for (const url of endpoints) {
        try {
          const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/sdp' },
            body: offer.sdp
          })
          if (res.ok) {
            answerSdp = await res.text()
            if (answerSdp) break
          }
        } catch (_) {
          // Network error or connection refused on this candidate port, try next
        }
      }

      if (!answerSdp) throw new Error('WHEP negotiation failed on all endpoints')
      await pc.setRemoteDescription({ type: 'answer', sdp: answerSdp })
    } catch (err: any) {
      error.value = err.message || 'WHEP failed'
      isConnecting.value = false
      isPlaying.value = false
      if (!document.hidden) setTimeout(() => start(activeStreamId, activeIsHero), 5000)
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
