export type AlertSeverity = 'info' | 'warning' | 'critical'

export interface AnalyticsAlert {
  id: string
  camera_id: string
  camera_name: string
  event_type: string
  severity: AlertSeverity
  confidence: number
  timestamp: string
  snapshot_url?: string
  is_acknowledged: boolean
}

export interface RegisteredCamera {
  id: string
  name: string
  rtsp_url: string
  ip: string
  port: number
  codec: 'H.264' | 'H.265'
  resolution: string
  fps: number
  has_ptz: boolean
  status: 'online' | 'offline' | 'recording'
  group: string
}

export interface AdminUser {
  id: string
  username: string
  role: 'admin' | 'supervisor' | 'operator'
  email: string
  allowed_cameras: string[]
  is_active: boolean
  created_at: string
}

export interface StorageVolume {
  id: string
  mount_point: string
  total_space_gb: number
  used_space_gb: number
  retention_days: number
  status: 'healthy' | 'warning' | 'error'
}
