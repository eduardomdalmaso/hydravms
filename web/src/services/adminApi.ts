// HydraVMS - Admin API Service for Maps, Layouts, Cluster Nodes, Logs & Users (< 100 lines)
import type { EnterpriseLayoutItem } from '../types/layoutTree'
import type { EnterpriseMapItem } from '../types/mapTree'
import type { EnterpriseRondaItem } from '../types/rondaTree'
import type { UserItem } from '../types/userTree'
import type { ServerNodeItem } from '../types/performanceCluster'
import type { LogEntry } from '../types/systemLogs'
import { fetchFolders } from './api'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'

export async function fetchLiveLayouts(fallback: EnterpriseLayoutItem[] = []): Promise<EnterpriseLayoutItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/layouts`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.layouts) && data.layouts.length > 0 ? data.layouts : fallback
  } catch {
    return fallback
  }
}

export async function fetchLiveMaps(fallback: EnterpriseMapItem[] = []): Promise<EnterpriseMapItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/maps`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.maps) && data.maps.length > 0 ? data.maps : fallback
  } catch {
    return fallback
  }
}

export async function fetchLiveTours(fallback: EnterpriseRondaItem[] = []): Promise<EnterpriseRondaItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/tours`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.tours) && data.tours.length > 0 ? data.tours : fallback
  } catch {
    return fallback
  }
}

export async function fetchLiveUsers(fallback: UserItem[] = []): Promise<UserItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/users`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.users) && data.users.length > 0 ? data.users : fallback
  } catch {
    return fallback
  }
}

export async function fetchLiveClusterNodes(fallback: ServerNodeItem[] = []): Promise<ServerNodeItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cluster/nodes`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.nodes) && data.nodes.length > 0 ? data.nodes : fallback
  } catch {
    return fallback
  }
}

export async function fetchLiveSystemLogs(fallback: LogEntry[] = []): Promise<LogEntry[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/system/logs`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) return fallback
    const data = await res.json()
    return Array.isArray(data.logs) && data.logs.length > 0 ? data.logs : fallback
  } catch {
    return fallback
  }
}

export { fetchFolders }
