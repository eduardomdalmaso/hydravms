import type { ServerNodeItem } from '../types/performanceCluster'

export const initialServerNodes: ServerNodeItem[] = [
  {
    id: 'node_master',
    hostname: 'srv-hydravms-core-01',
    ip: '192.168.1.10',
    role: 'MASTER_VMS',
    status: 'ONLINE',
    uptime: '18d 04h 12m',
    cpuPercent: 24,
    cpuModel: 'AMD EPYC 9354 32-Core Processor',
    ramUsedGb: 34.2,
    ramTotalGb: 128,
    gpus: []
  },
  {
    id: 'node_ingest_01',
    hostname: 'hydrastream-edge-ingest-01',
    ip: '192.168.1.20',
    role: 'HYDRASTREAM_EDGE',
    status: 'ONLINE',
    uptime: '14d 19h 45m',
    cpuPercent: 46,
    cpuModel: 'AMD Ryzen 9 9950X 16-Core Processor',
    ramUsedGb: 22.8,
    ramTotalGb: 64,
    gpus: []
  },
  {
    id: 'node_gpu_worker_01',
    hostname: 'hydrastream-ai-worker-01',
    ip: '192.168.1.30',
    role: 'HYDRASTREAM_GPU_WORKER',
    status: 'ONLINE',
    uptime: '09d 08h 30m',
    cpuPercent: 58,
    cpuModel: 'Intel Xeon w9-3495X 56-Core Processor',
    ramUsedGb: 78.4,
    ramTotalGb: 256,
    gpus: [
      {
        id: 'gpu_0_worker1',
        name: 'NVIDIA GeForce RTX 5090',
        index: 0,
        tempC: 58,
        powerWatts: 295,
        vramUsedGb: 14.8,
        vramTotalGb: 32.0,
        computePercent: 74
      },
      {
        id: 'gpu_1_worker1',
        name: 'NVIDIA GeForce RTX 5090',
        index: 1,
        tempC: 62,
        powerWatts: 340,
        vramUsedGb: 18.2,
        vramTotalGb: 32.0,
        computePercent: 82
      }
    ]
  },
  {
    id: 'node_gpu_worker_02',
    hostname: 'hydrastream-ai-worker-02',
    ip: '192.168.1.31',
    role: 'HYDRASTREAM_GPU_WORKER',
    status: 'ONLINE',
    uptime: '06d 12h 10m',
    cpuPercent: 42,
    cpuModel: 'AMD EPYC 7763 64-Core Processor',
    ramUsedGb: 51.0,
    ramTotalGb: 128,
    gpus: [
      {
        id: 'gpu_0_worker2',
        name: 'NVIDIA GeForce RTX 4090',
        index: 0,
        tempC: 54,
        powerWatts: 245,
        vramUsedGb: 12.1,
        vramTotalGb: 24.0,
        computePercent: 61
      }
    ]
  }
]
