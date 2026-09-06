import type { StoragePoolItem, UnallocatedDiskDevice } from '../types/storagePool'

export const initialStoragePools: StoragePoolItem[] = [
  {
    id: 'sp_01',
    name: 'Samsung 990 PRO NVMe 2TB',
    sourceType: 'LOCAL_DISK',
    role: 'HOT_BUFFER',
    nodeOrServer: 'MAQUINA LOCAL // NÓ 01',
    pathOrEndpoint: '/var/lib/hydravms/hot-01',
    filesystem: 'XFS (allocsize=64M)',
    totalGb: 1920,
    usedGb: 480,
    status: 'ONLINE',
    isSpilloverActive: true,
    retentionDays: 2
  },
  {
    id: 'sp_02',
    name: 'WD Purple Pro 18TB (/dev/sdb)',
    sourceType: 'LOCAL_DISK',
    role: 'WARM_ARCHIVE',
    nodeOrServer: 'MAQUINA LOCAL // NÓ 01',
    pathOrEndpoint: '/var/lib/hydravms/warm-01',
    filesystem: 'XFS (noatime)',
    totalGb: 18000,
    usedGb: 8200,
    status: 'ONLINE',
    retentionDays: 30
  },
  {
    id: 'sp_03',
    name: 'PostgreSQL Events DB (/dev/sda2)',
    sourceType: 'LOCAL_DISK',
    role: 'DATABASE',
    nodeOrServer: 'MAQUINA LOCAL // NÓ 01',
    pathOrEndpoint: '/var/lib/hydravms/db-events',
    filesystem: 'EXT4 (writeback)',
    totalGb: 500,
    usedGb: 54,
    status: 'ONLINE',
    retentionDays: 90
  },
  {
    id: 'sp_04',
    name: 'NAS Synology RS4021xs+ (Rack)',
    sourceType: 'NETWORK_NAS',
    role: 'WARM_ARCHIVE',
    nodeOrServer: 'SERVIDOR REMOTO // 192.168.1.150',
    pathOrEndpoint: 'nfs://192.168.1.150/volume1/cftv_longterm',
    filesystem: 'NFSv4',
    totalGb: 64000,
    usedGb: 28400,
    status: 'ONLINE',
    retentionDays: 60
  },
  {
    id: 'sp_05',
    name: 'Storage S3 Snapshots & Evidencias IA',
    sourceType: 'OBJECT_S3',
    role: 'SNAPSHOTS',
    nodeOrServer: 'CLUSTER S3 // minio.corp.lan',
    pathOrEndpoint: 's3://hydravms-snapshots',
    filesystem: 'S3_API',
    totalGb: 20000,
    usedGb: 2400,
    status: 'ONLINE',
    retentionDays: 180
  }
]

export const detectedLocalDisks: UnallocatedDiskDevice[] = [
  { devicePath: '/dev/sdc', model: 'Seagate SkyHawk AI 16TB', sizeGb: 16000, busType: 'SATA', status: 'NAO_FORMATADO' },
  { devicePath: '/dev/nvme1n1', model: 'Kingston KC3000 NVMe 2TB', sizeGb: 2048, busType: 'NVMe', status: 'PARTICAO_LIVRE' }
]
