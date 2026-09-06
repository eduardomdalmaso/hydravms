export type StorageSourceType = 'LOCAL_DISK' | 'NETWORK_NAS' | 'OBJECT_S3'
export type StorageRole = 'HOT_BUFFER' | 'WARM_ARCHIVE' | 'SNAPSHOTS' | 'DATABASE'
export type StorageStatus = 'ONLINE' | 'STANDBY' | 'MAINTENANCE'

export interface StoragePoolItem {
  id: string
  name: string
  sourceType: StorageSourceType
  role: StorageRole
  nodeOrServer: string
  pathOrEndpoint: string
  filesystem: string
  totalGb: number
  usedGb: number
  status: StorageStatus
  isSpilloverActive?: boolean
  retentionDays?: number
}

export interface UnallocatedDiskDevice {
  devicePath: string
  model: string
  sizeGb: number
  busType: 'NVMe' | 'SATA' | 'SAS' | 'USB'
  status: 'NAO_FORMATADO' | 'PARTICAO_LIVRE'
}

export interface NewStoragePayload {
  name: string
  sourceType: StorageSourceType
  role: StorageRole
  nodeOrServer: string
  pathOrEndpoint: string
  filesystem: string
  totalGb: number
  retentionDays?: number
}
