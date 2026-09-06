export type PluginCategory = 'analytics' | 'access_control' | 'safety' | 'traffic'

export type PluginStatus = 'available' | 'installed' | 'running' | 'stopped' | 'updating'

export interface PluginManifest {
  id: string
  name: string
  version: string
  author: string
  category: PluginCategory
  description: string
  runtime: 'python3' | 'binary_elf' | 'docker'
  hardware_req: 'CUDA 13.3 // RTX 5090' | 'CPU // ZERO-COPY SHM' | 'AVX2'
  permissions: string[]
  is_official: boolean
  is_installed: boolean
  status: PluginStatus
  instances_count: number
  events_count: number
  default_config: Record<string, any>
}

export interface AnalyticInstance {
  id: string
  plugin_id: string
  plugin_name: string
  name: string
  camera_id: string
  camera_name: string
  stream_type: 'main_1080p' | 'sub_stream'
  hardware_target: 'rtx_5090_cuda' | 'cpu_shm'
  confidence_threshold: number
  roi_mode: 'full_frame' | 'custom_polygon'
  specific_params: Record<string, any>
  is_active: boolean
  fps_rate: number
  detections_count: number
  created_at: string
}

export interface AnalyticEventRecord {
  id: string
  timestamp: string
  plugin_id: string
  plugin_name: string
  camera_id: string
  camera_name: string
  event_type: string
  severity: 'info' | 'warning' | 'critical'
  confidence: number
  object_label: string
  details: string
  snapshot_url?: string
  bbox?: [number, number, number, number]
  raw_payload: Record<string, any>
}

export interface AnalyticFolderNode {
  id: string
  name: string
  plugin_id: string
  instances: AnalyticInstance[]
}

export type MarketplaceTab = 'catalog' | 'create' | 'explorer'
