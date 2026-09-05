export interface StreamItem {
  id: string
  name: string
  protocol: 'RTSP' | 'RTMP' | 'ONVIF'
  url: string
  ip: string
  port: number
  codec: 'H.265' | 'H.264'
  resolution: string
  fps: number
  bitrate: string
  recordMode: 'continuous' | 'motion' | 'ai_event' | 'disabled'
  status: 'online' | 'offline' | 'recording'
  has_ptz: boolean
}

export interface FolderNode {
  id: string
  name: string
  isExpanded: boolean
  streams: StreamItem[]
}
