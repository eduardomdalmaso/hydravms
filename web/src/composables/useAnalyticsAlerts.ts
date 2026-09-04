import { ref, computed } from 'vue'
import type { AnalyticsAlert } from '../types/admin'

export function useAnalyticsAlerts() {
  const isDrawerOpen = ref(false)
  const alerts = ref<AnalyticsAlert[]>([
    {
      id: 'alt_001',
      camera_id: 'cam_01',
      camera_name: 'Portaria Principal (Entrada)',
      event_type: 'INTRUSAO PERIMETRO (HUMANO)',
      severity: 'critical',
      confidence: 0.94,
      timestamp: '20:39:12',
      is_acknowledged: false
    },
    {
      id: 'alt_002',
      camera_id: 'cam_04',
      camera_name: 'Perimetro dos Fundos',
      event_type: 'CRUZAMENTO LINHA VIRTUAL',
      severity: 'warning',
      confidence: 0.88,
      timestamp: '20:37:45',
      is_acknowledged: false
    },
    {
      id: 'alt_003',
      camera_id: 'cam_03',
      camera_name: 'Corredor de Cargas & Docas',
      event_type: 'OBSTACULO EM ZONA DE EMBARQUE',
      severity: 'info',
      confidence: 0.91,
      timestamp: '20:32:00',
      is_acknowledged: true
    }
  ])

  const unreadCount = computed(() => alerts.value.filter(a => !a.is_acknowledged).length)

  const toggleDrawer = () => {
    isDrawerOpen.value = !isDrawerOpen.value
  }

  const acknowledgeAlert = (id: string) => {
    const target = alerts.value.find(a => a.id === id)
    if (target) target.is_acknowledged = true
  }

  const clearAllAlerts = () => {
    alerts.value.forEach(a => { a.is_acknowledged = true })
  }

  return {
    isDrawerOpen,
    alerts,
    unreadCount,
    toggleDrawer,
    acknowledgeAlert,
    clearAllAlerts
  }
}
