// HydraVMS - Resilient API & WebSocket Client with Graceful Fallback (< 100 lines)
import type { RegisteredCamera } from '../types/admin'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'
const WS_BASE = import.meta.env.VITE_WS_URL || 'ws://localhost:8083/ws/v1/live'

export interface ApiFolder {
  id: string
  tenant_id: string
  module: string
  name: string
  parent_id?: string
  color_hex?: string
}

export async function fetchFolders(module: string, fallback: ApiFolder[] = []): Promise<ApiFolder[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/folders?module=${encodeURIComponent(module)}`, {
      headers: { 'Accept': 'application/json' },
      signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.folders) && data.folders.length > 0 ? data.folders : fallback
  } catch {
    return fallback
  }
}

export async function createRemoteFolder(module: string, name: string, parentId?: string): Promise<ApiFolder | null> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/folders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ module, name, parent_id: parentId }),
      signal: AbortSignal.timeout(3000)
    })
    return res.ok ? await res.json() : null
  } catch {
    return null
  }
}

export async function fetchCameras(fallback: RegisteredCamera[] = []): Promise<RegisteredCamera[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cameras`, {
      headers: { 'Accept': 'application/json' },
      signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.cameras) && data.cameras.length > 0 ? data.cameras : fallback
  } catch {
    return fallback
  }
}

export function connectLiveWebSocket(onMessage: (data: any) => void, onStatusChange?: (online: boolean) => void) {
  let ws: WebSocket | null = null, retryCount = 0, isStopped = false

  const connect = () => {
    if (isStopped) return
    try {
      ws = new WebSocket(WS_BASE)
      ws.onopen = () => { retryCount = 0; onStatusChange?.(true) }
      ws.onmessage = (e) => {
        try { onMessage(JSON.parse(e.data)) } catch { onMessage(e.data) }
      }
      ws.onclose = () => {
        onStatusChange?.(false)
        if (!isStopped) {
          const delay = Math.min(30000, Math.pow(1.5, retryCount++) * 1000 + Math.random() * 500)
          setTimeout(connect, delay)
        }
      }
      ws.onerror = () => { ws?.close() }
    } catch {
      setTimeout(connect, 3000)
    }
  }

  connect()
  return {
    disconnect: () => { isStopped = true; ws?.close() }
  }
}
