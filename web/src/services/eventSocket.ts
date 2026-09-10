// HydraVMS - Real-Time CloudEvents WebSocket Bus (< 100 lines)
import { reactive, readonly } from 'vue'

export interface CloudEvent<T = any> {
  specversion: string
  id: string
  source: string
  type: string
  subject: string
  time: string
  tenant_id: string
  category: string
  severity: 'low' | 'medium' | 'warning' | 'critical'
  data: T
}

type EventListener = (event: CloudEvent) => void

const state = reactive({
  isConnected: false,
  cameraStatuses: {} as Record<string, 'online' | 'offline' | 'reconnecting'>,
  recentEvents: [] as CloudEvent[]
})

const listeners = new Set<EventListener>()
let ws: WebSocket | null = null, isExplicitlyClosed = false, retryAttempts = 0

function getWsUrl(): string {
  const isHttps = typeof window !== 'undefined' && window.location.protocol === 'https:'
  const host = (typeof window !== 'undefined' && window.location.hostname) || 'localhost'
  const token = (typeof localStorage !== 'undefined' && localStorage.getItem('vms_token')) || 'default_token'
  return `${isHttps ? 'wss:' : 'ws:'}//${host}:8083/ws/v1/live?token=${encodeURIComponent(token)}`
}

export function initEventSocket() {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return
  isExplicitlyClosed = false

  try {
    ws = new WebSocket(getWsUrl())
    ws.onopen = () => {
      state.isConnected = true
      retryAttempts = 0
      ws?.send(JSON.stringify({ action: 'subscribe', topics: ['*'] }))
    }

    ws.onmessage = (e) => {
      try {
        const parsed = JSON.parse(e.data) as CloudEvent
        if (parsed && parsed.type) {
          handleIncomingEvent(parsed)
        }
      } catch {}
    }

    ws.onclose = () => {
      state.isConnected = false
      if (!isExplicitlyClosed) {
        const delay = Math.min(15000, Math.pow(1.5, retryAttempts++) * 1000 + 500)
        setTimeout(initEventSocket, delay)
      }
    }

    ws.onerror = () => ws?.close()
  } catch {
    setTimeout(initEventSocket, 3000)
  }
}

function handleIncomingEvent(evt: CloudEvent) {
  state.recentEvents.unshift(evt)
  if (state.recentEvents.length > 50) state.recentEvents.pop()

  if (evt.type === 'system.camera.offline' && evt.subject) {
    state.cameraStatuses[evt.subject] = 'offline'
  } else if (evt.type === 'system.camera.online' && evt.subject) {
    state.cameraStatuses[evt.subject] = 'online'
  }

  listeners.forEach(fn => { try { fn(evt) } catch {} })
}

export function useEventBus() {
  const subscribe = (fn: EventListener) => {
    listeners.add(fn)
    return () => listeners.delete(fn)
  }

  return {
    state: readonly(state),
    subscribe,
    getCameraStatus: (camId: string) => state.cameraStatuses[camId]
  }
}
