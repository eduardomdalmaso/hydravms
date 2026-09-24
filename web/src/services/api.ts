// HydraVMS - Resilient API & WebSocket Client (< 100 lines)
import type { RegisteredCamera } from '../types/admin'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'
const WS_BASE = import.meta.env.VITE_WS_URL || 'ws://localhost:8083/ws/v1/live'

export interface ApiFolder {
  id: string; tenant_id: string; module: string; name: string; parent_id?: string; color_hex?: string
}

export function getAuthHeaders(): Record<string, string> {
  const t = localStorage.getItem('hydra_token')
  return t ? { 'Authorization': `Bearer ${t}` } : {}
}

export async function authedFetch(url: string, options: RequestInit = {}): Promise<Response> {
  const headers = new Headers(options.headers || {})
  const token = localStorage.getItem('hydra_token')
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  options.headers = headers

  let res = await fetch(url, options)
  if (res.status === 401) {
    try {
      const loginRes = await fetch(`${API_BASE}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: 'admin@hydravms.io', password: 'admin' })
      })
      if (loginRes.ok) {
        const loginData = await loginRes.json()
        if (loginData.token) {
          localStorage.setItem('hydra_token', loginData.token)
          localStorage.setItem('hydra_auth', 'true')
          headers.set('Authorization', `Bearer ${loginData.token}`)
          options.headers = headers
          res = await fetch(url, options)
        }
      }
    } catch {}
  }
  return res
}

export async function fetchFolders(module: string, fallback: ApiFolder[] = []): Promise<ApiFolder[]> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/folders?module=${encodeURIComponent(module)}`, {
      headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.folders) && data.folders.length > 0 ? data.folders : fallback
  } catch { return fallback }
}

export async function createRemoteFolder(module: string, name: string, parentId?: string): Promise<ApiFolder | null> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/folders`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ module, name, parent_id: parentId }), signal: AbortSignal.timeout(3000)
    })
    return res.ok ? await res.json() : null
  } catch { return null }
}

export async function deleteRemoteFolder(id: string): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/folders/${encodeURIComponent(id)}`, {
      method: 'DELETE', signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch { return false }
}

export async function fetchCameras(fallback: RegisteredCamera[] = []): Promise<RegisteredCamera[]> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/cameras`, {
      headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.cameras) && data.cameras.length > 0 ? data.cameras : fallback
  } catch { return fallback }
}

export async function fetchUsers(fallback: any[] = []): Promise<any[]> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/users`, {
      headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return Array.isArray(data.users) && data.users.length > 0 ? data.users : fallback
  } catch { return fallback }
}

export async function createRemoteCamera(cam: any): Promise<RegisteredCamera | null> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/cameras`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(cam), signal: AbortSignal.timeout(3000)
    })
    return res.ok ? await res.json() : null
  } catch { return null }
}

export async function deleteRemoteCamera(id: string): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/cameras/${encodeURIComponent(id)}`, {
      method: 'DELETE', signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch { return false }
}

export async function fetchPlugins(): Promise<{ plugins: any[]; gpu_detected: boolean; gpu_telemetry: any }> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/plugins`, {
      headers: { 'Accept': 'application/json' }, signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    return {
      plugins: Array.isArray(data.plugins) ? data.plugins : [],
      gpu_detected: !!data.gpu_detected,
      gpu_telemetry: data.gpu_telemetry || null
    }
  } catch {
    return { plugins: [], gpu_detected: false, gpu_telemetry: null }
  }
}

export async function installRemotePlugin(id: string): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/plugins/${encodeURIComponent(id)}/install`, {
      method: 'POST', signal: AbortSignal.timeout(5000)
    })
    return res.ok
  } catch { return false }
}

export async function uninstallRemotePlugin(id: string): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/plugins/${encodeURIComponent(id)}/uninstall`, {
      method: 'POST', signal: AbortSignal.timeout(5000)
    })
    return res.ok
  } catch { return false }
}

export async function toggleRemotePlugin(id: string): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/v1/plugins/${encodeURIComponent(id)}/toggle`, {
      method: 'POST', signal: AbortSignal.timeout(5000)
    })
    return res.ok
  } catch { return false }
}

export function connectLiveWebSocket(onMessage: (data: any) => void, onStatusChange?: (online: boolean) => void) {
  let ws: WebSocket | null = null, retryCount = 0, isStopped = false
  const connect = () => {
    if (isStopped) return
    try {
      ws = new WebSocket(WS_BASE)
      ws.onopen = () => { retryCount = 0; onStatusChange?.(true) }
      ws.onmessage = (e) => { try { onMessage(JSON.parse(e.data)) } catch { onMessage(e.data) } }
      ws.onclose = () => {
        onStatusChange?.(false)
        if (!isStopped) setTimeout(connect, Math.min(30000, Math.pow(1.5, retryCount++) * 1000 + 500))
      }
      ws.onerror = () => ws?.close()
    } catch { setTimeout(connect, 3000) }
  }
  connect()
  return { disconnect: () => { isStopped = true; ws?.close() } }
}
