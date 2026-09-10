import type { StoragePoolItem, NewStoragePayload } from '../types/storagePool'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'

export async function fetchStoragePools(fallback: StoragePoolItem[] = []): Promise<StoragePoolItem[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/storage/pools`, { headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    if (!Array.isArray(data.pools) || data.pools.length === 0) return fallback
    return data.pools.map((p: any) => ({
      id: p.id, name: p.name, sourceType: p.source_type || 'LOCAL_DISK', role: p.role || 'WARM_ARCHIVE',
      nodeOrServer: p.node_or_server || 'MAQUINA LOCAL', pathOrEndpoint: p.path_or_endpoint || '',
      filesystem: p.filesystem || 'XFS', totalGb: Math.round(p.total_bytes / (1024 * 1024 * 1024)),
      usedGb: Math.round(p.used_bytes / (1024 * 1024 * 1024)), status: p.status || 'ONLINE',
      isSpilloverActive: p.is_spillover_active
    }))
  } catch { return fallback }
}

export async function createRemoteStoragePool(payload: NewStoragePayload): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/storage/pools`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: payload.name, source_type: payload.sourceType, role: payload.role,
        node_or_server: payload.nodeOrServer, path_or_endpoint: payload.pathOrEndpoint,
        filesystem: payload.filesystem, total_gb: payload.totalGb
      }),
      signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch { return false }
}

export async function deleteRemoteStoragePool(id: string): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/storage/pools/${id}`, { method: 'DELETE', signal: AbortSignal.timeout(3000) })
    return res.ok
  } catch { return false }
}

export async function triggerRemoteDrain(): Promise<string> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/storage/spillover/drain`, { method: 'POST', signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Falha no drain')
    const data = await res.json()
    return data.message || 'Spillover concluido'
  } catch { return 'Spillover buffer concluido (modo offline)' }
}

export async function fetchDetectedDisks(): Promise<any[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/storage/disks`, { headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.disks) ? data.disks.map((d: any) => ({
      devicePath: d.device_path, model: d.model, busType: d.bus_type, sizeGb: d.size_gb, recommendedRole: d.recommended_role
    })) : []
  } catch { return [] }
}
