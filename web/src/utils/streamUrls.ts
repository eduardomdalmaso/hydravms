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
  const wspPort = (typeof import.meta !== 'undefined' && import.meta.env?.VITE_WSP_PORT) || '8086'
  const targetId = (!isHero && hasSub) ? `${cameraId}_sub` : cameraId
  return `ws://${host}:${wspPort}/live/${encodeURIComponent(targetId)}`
}

export const getCameraSnapshotUrl = (cameraId: string, refresh = false): string => {
  const base = getStreamBaseUrl()
  const q = refresh ? `?refresh=true&t=${Date.now()}` : `?t=${Date.now()}`
  return `${base}/api/v1/streams/${encodeURIComponent(cameraId)}/snapshot.jpg${q}`
}

export const getCameraMjpegUrl = (cameraId: string, quality = 4, fps = 20): string => {
  const base = getStreamBaseUrl()
  return `${base}/api/v1/streams/${encodeURIComponent(cameraId)}/mjpeg?q=${quality}&fps=${fps}`
}
