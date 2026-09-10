import { ref, computed } from 'vue'
import type { AnalyticsAlert } from '../types/admin'

export function useAnalyticsAlerts() {
  const isDrawerOpen = ref(false)
  const alerts = ref<AnalyticsAlert[]>([])

  const unreadCount = computed(() => alerts.value.filter(a => !a.is_acknowledged).length)
  const toggleDrawer = () => { isDrawerOpen.value = !isDrawerOpen.value }
  const acknowledgeAlert = (id: string) => {
    const target = alerts.value.find(a => a.id === id)
    if (target) target.is_acknowledged = true
  }
  const clearAllAlerts = () => { alerts.value.forEach(a => { a.is_acknowledged = true }) }

  return { isDrawerOpen, alerts, unreadCount, toggleDrawer, acknowledgeAlert, clearAllAlerts }
}
