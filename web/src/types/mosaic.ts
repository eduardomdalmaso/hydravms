export type GridLayout = 'auto' | '1x1' | '1x2' | '2x2' | '3x3' | '4x4' | '8x8' | '10x10' | '1+5' | '1+7' | '1+12'
export type PlaybackSpeed = 0.5 | 1 | 2 | 3
export type SlotContentType = 'camera' | 'map' | 'carousel'
export type StreamProtocol = 'RTSP' | 'RTMP' | 'ONVIF'

export interface CameraStreamInfo {
  id: string
  name: string
  location: string
  status: 'online' | 'offline' | 'recording'
  protocol?: StreamProtocol
  codec?: string
  has_ptz: boolean
  fps: number
  resolution: string
  main_stream_url?: string
  sub_stream_url?: string
}

export interface MapResource {
  id: string
  name: string
  image_url: string
  cameras_count: number
}

export interface CarouselConfig {
  id: string
  name: string
  camera_ids: string[]
  interval_seconds: number
}

export interface WorkspaceSlot {
  slot_index: number
  type: SlotContentType
  data?: CameraStreamInfo | MapResource | CarouselConfig
}

export interface CustomLayout {
  id: string
  name: string
  grid: GridLayout
  is_system: boolean
  created_by: string
  slots?: WorkspaceSlot[]
}

export interface WorkspaceTab {
  id: string
  name: string
  grid: GridLayout
  is_system: boolean
  is_temporary?: boolean
  is_saved?: boolean
  created_by: string
  slots: WorkspaceSlot[]
}

export interface TimelineSegment {
  id: string
  start_time: number
  end_time: number
  type: 'continuous' | 'motion' | 'ai_alert'
  label?: string
}

export type AlarmStatus = 'online' | 'alert' | 'offline'

export interface AlarmItemInfo {
  id: string
  name: string
  zone: string
  status: AlarmStatus
  type?: string
}
