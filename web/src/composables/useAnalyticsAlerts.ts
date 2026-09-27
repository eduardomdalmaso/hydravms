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
  else if (evt.type.includes('zone_intrusion') || evt.type.includes('intrusion')) title = 'INTRUSÃO // DETECÇÃO DE PESSOA'
  else if (evt.type.includes('crowd')) title = 'MULTIDÃO // AGLOMERAÇÃO'
  else if (isAI) title = `DETECÇÃO IA // ${evt.data?.class_name?.toUpperCase() || 'PESSOA'}`

  const rawBbox = evt.data?.bbox
  let bboxCoords: [number, number, number, number] | undefined = undefined
  if (rawBbox) {
    if (Array.isArray(rawBbox)) {
      bboxCoords = rawBbox as [number, number, number, number]
    } else if (typeof rawBbox === 'object' && rawBbox.width) {
      bboxCoords = [
        (rawBbox.x_center - rawBbox.width / 2) * 100,
        (rawBbox.y_center - rawBbox.height / 2) * 100,
        rawBbox.width * 100,
        rawBbox.height * 100
      ]
    }
  }

  const camId = evt.subject || evt.data?.camera_id || 'cam_01'
  const camName = evt.data?.camera_name || (camId === 'cam_01' || camId === 'c113' ? 'C113' : camId.toUpperCase())

  alerts.value.unshift({
    id: evt.id || `alert_${Date.now()}`,
    camera_id: camId,
    camera_name: camName,
    event_type: title,
    severity: (evt.severity as any) || (isOffline ? 'critical' : 'warning'),
    confidence: evt.data?.confidence || 0.92,
    timestamp: new Date().toLocaleTimeString(),
    snapshot_url: evt.data?.snapshot_url || `http://localhost:8080/api/v1/streams/${camId}/snapshot`,
    bbox: bboxCoords,
    is_acknowledged: false
  })

  // Keep last 50 alerts
  if (alerts.value.length > 50) alerts.value.pop()
})

export function useAnalyticsAlerts() {
  const unreadCount = computed(() => alerts.value.filter(a => !a.is_acknowledged).length)
  const toggleDrawer = () => { isDrawerOpen.value = !isDrawerOpen.value }
  const acknowledgeAlert = (id: string) => {
    alerts.value = alerts.value.filter(a => a.id !== id)
  }
  const clearAllAlerts = () => {
    alerts.value = []
  }

  return { isDrawerOpen, alerts, unreadCount, toggleDrawer, acknowledgeAlert, clearAllAlerts }
}
