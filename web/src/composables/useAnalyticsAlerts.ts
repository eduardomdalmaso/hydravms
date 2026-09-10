import { ref, computed } from 'vue'
import type { AnalyticsAlert } from '../types/admin'
import { useEventBus, type CloudEvent } from '../services/eventSocket'

const isDrawerOpen = ref(false)
const alerts = ref<AnalyticsAlert[]>([])

// Singleton event listener across the entire app
const eventBus = useEventBus()
eventBus.subscribe((evt: CloudEvent) => {
  const isOffline = evt.type === 'system.camera.offline'
  const isOnline = evt.type === 'system.camera.online'
  const isAI = evt.type.startsWith('ai.')

  let title = evt.type.toUpperCase()
  if (isOffline) title = 'CAMERA OFFLINE // SINAL RTSP INTERROMPIDO'
  else if (isOnline) title = 'CAMERA ONLINE // SINAL RESTABELECIDO'
  else if (isAI) title = `DETECCAO IA // ${evt.data?.class_name?.toUpperCase() || 'ALERTA'}`

  alerts.value.unshift({
    id: evt.id || `alert_${Date.now()}`,
    camera_id: evt.subject || 'cam_01',
    camera_name: evt.data?.camera_name || evt.subject || 'CAMERA',
    event_type: title,
    severity: (evt.severity as any) || (isOffline ? 'critical' : 'medium'),
    confidence: evt.data?.confidence || 1.0,
    timestamp: new Date().toLocaleTimeString(),
    snapshot_url: evt.data?.snapshot_url,
    is_acknowledged: false
  })

  // Keep last 50 alerts
  if (alerts.value.length > 50) alerts.value.pop()
})

export function useAnalyticsAlerts() {
  const unreadCount = computed(() => alerts.value.filter(a => !a.is_acknowledged).length)
  const toggleDrawer = () => { isDrawerOpen.value = !isDrawerOpen.value }
  const acknowledgeAlert = (id: string) => {
    const target = alerts.value.find(a => a.id === id)
    if (target) target.is_acknowledged = true
  }
  const clearAllAlerts = () => { alerts.value.forEach(a => { a.is_acknowledged = true }) }

  return { isDrawerOpen, alerts, unreadCount, toggleDrawer, acknowledgeAlert, clearAllAlerts }
}
