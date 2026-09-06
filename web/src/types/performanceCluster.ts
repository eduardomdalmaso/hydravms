export type ServerNodeType = 'MASTER_VMS' | 'HYDRASTREAM_EDGE' | 'HYDRASTREAM_GPU_WORKER'
export type ServerNodeStatus = 'ONLINE' | 'DEGRADED' | 'OFFLINE'

export interface GpuTelemetry {
  id: string
  name: string
  index: number
  tempC: number
  powerWatts: number
  vramUsedGb: number
  vramTotalGb: number
  computePercent: number
}

export interface ServerNodeItem {
  id: string
  hostname: string
  ip: string
  role: ServerNodeType
  status: ServerNodeStatus
  uptime: string
  cpuPercent: number
  cpuModel: string
  ramUsedGb: number
  ramTotalGb: number
  gpus: GpuTelemetry[]
}
