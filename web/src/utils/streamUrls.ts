/**
 * Dynamic resolution of Media/Stream Plane URLs for multi-node clusters and remote networks.
 */
export const getStreamBaseUrl = (): string => {
  if (typeof import.meta !== 'undefined' && import.meta.env?.VITE_STREAM_URL) {
    return (import.meta.env.VITE_STREAM_URL as string).replace(/\/+$/, '')
  }
  if (typeof window !== 'undefined' && window.location) {
    return `${window.location.protocol}//${window.location.hostname}:8080`
  }
  return 'http://localhost:8080'
}

export const getStreamWspUrl = (cameraId: string, isHero = false, hasSub = true): string => {
  const host = typeof window !== 'undefined' && window.location ? window.location.hostname : 'localhost'
  const wspPort = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_WSP_PORT) || '8085'
  const targetId = (!isHero && hasSub) ? `${cameraId}_sub` : cameraId
  return `ws://${host}:${wspPort}/live/${encodeURIComponent(targetId)}`
}

export const getCameraSnapshotUrl = (cameraId: string, refresh = false): string => {
  const base = getStreamBaseUrl()
  const q = refresh ? `?refresh=true&t=${Date.now()}` : `?t=${Date.now()}`
  return `${base}/api/v1/streams/${encodeURIComponent(cameraId)}/snapshot.jpg${q}`
}

export const getCameraMjpegUrl = (cameraId: string, quality = 4, fps = 20, isHero = false, hasSub = true): string => {
  const base = getStreamBaseUrl()
  const targetId = (!isHero && hasSub) ? `${cameraId}_sub` : cameraId
  return `${base}/api/v1/streams/${encodeURIComponent(targetId)}/mjpeg?q=${quality}&fps=${fps}`
}

export const getStreamRelayRtspUrl = (cameraId: string, isSub = false): string => {
  const host = typeof window !== 'undefined' && window.location ? window.location.hostname : 'localhost'
  const rtspPort = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_RTSP_RELAY_PORT) || '8554'
  const targetId = isSub ? `${cameraId}_sub` : cameraId
  return `rtsp://${host}:${rtspPort}/live/${encodeURIComponent(targetId)}`
}

export const maskRtspUrl = (url: string, mask = true): string => {
  if (!url || !mask) return url
  return url.replace(/^(rtsp[s]?:\/\/[^:@/]+):([^@/]+)@/, '$1:••••••••@')
}

export const copyToClipboard = async (text: string): Promise<boolean> => {
  if (!text) return false
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
      return true
    }
  } catch {}
  return false
}
