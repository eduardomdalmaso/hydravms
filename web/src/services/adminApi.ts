// HydraVMS - Admin API Service for Maps, Layouts, Cluster Nodes, Logs & Users (< 100 lines)
import type { EnterpriseLayoutItem } from '../types/layoutTree'
import type { EnterpriseMapItem } from '../types/mapTree'
import type { EnterpriseRondaItem } from '../types/rondaTree'
import type { UserItem } from '../types/userTree'
import type { ServerNodeItem } from '../types/performanceCluster'
import type { LogEntry } from '../types/systemLogs'
import { fetchFolders } from './api'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'

export async function fetchLiveClusterNodes(fallback: ServerNodeItem[] = []): Promise<ServerNodeItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cluster/nodes`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    if (!Array.isArray(data.nodes) || data.nodes.length === 0) return fallback
    return data.nodes.map((n: any) => ({
      id: n.id, hostname: n.node_name, ip: n.ip_address,
      role: n.node_role === 'gpu_worker' ? 'HYDRASTREAM_GPU_WORKER' : (n.node_role === 'control_plane' ? 'MASTER_VMS' : 'HYDRASTREAM_EDGE'),
      status: n.status === 'online' ? 'ONLINE' : 'OFFLINE', uptime: 'Ativo',
      cpuPercent: n.cpu_usage_pct || 18, cpuModel: 'AMD / Intel Server Core',
      ramUsedGb: 8.4, ramTotalGb: 64,
      gpus: n.gpu_device_info ? [{ id: 'gpu_0', name: n.gpu_device_info, index: 0, tempC: 48, powerWatts: 220, vramUsedGb: 6.2, vramTotalGb: 32.0, computePercent: 25 }] : []
    }))
  } catch { return fallback }
}

export async function probeClusterNode(ip: string, port: number): Promise<any> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cluster/probe`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ip_address: ip, http_port: port }), signal: AbortSignal.timeout(3000)
    })
    return res.ok ? await res.json() : { online: false, error: 'HTTP error' }
  } catch (err: any) { return { online: false, error: err?.message || 'Timeout' } }
}

export async function createRemoteClusterNode(node: { node_name: string; node_role: string; ip_address: string; http_port: number; grpc_port: number; webrtc_port: number; gpu_device_info?: string }): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cluster/nodes`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(node), signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch { return false }
}

export async function deleteRemoteClusterNode(id: string): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cluster/nodes/${id}`, { method: 'DELETE', signal: AbortSignal.timeout(3000) })
    return res.ok
  } catch { return false }
}

export async function fetchLiveLayouts(f: EnterpriseLayoutItem[] = []): Promise<EnterpriseLayoutItem[]> {
  try { const r = await fetch(`${API_BASE}/api/v1/layouts`, { signal: AbortSignal.timeout(3000) }); return r.ok ? (await r.json()).layouts || f : f } catch { return f }
}
export async function fetchLiveMaps(f: EnterpriseMapItem[] = []): Promise<EnterpriseMapItem[]> {
  try { const r = await fetch(`${API_BASE}/api/v1/maps`, { signal: AbortSignal.timeout(3000) }); return r.ok ? (await r.json()).maps || f : f } catch { return f }
}
export async function fetchLiveTours(f: EnterpriseRondaItem[] = []): Promise<EnterpriseRondaItem[]> {
  try { const r = await fetch(`${API_BASE}/api/v1/tours`, { signal: AbortSignal.timeout(3000) }); return r.ok ? (await r.json()).tours || f : f } catch { return f }
}
export async function fetchLiveUsers(f: UserItem[] = []): Promise<UserItem[]> {
  try { const r = await fetch(`${API_BASE}/api/v1/users`, { signal: AbortSignal.timeout(3000) }); return r.ok ? (await r.json()).users || f : f } catch { return f }
}
export async function fetchLiveSystemLogs(f: LogEntry[] = []): Promise<LogEntry[]> {
  try { const r = await fetch(`${API_BASE}/api/v1/system/logs`, { signal: AbortSignal.timeout(3000) }); return r.ok ? (await r.json()).logs || f : f } catch { return f }
}

export { fetchFolders }
