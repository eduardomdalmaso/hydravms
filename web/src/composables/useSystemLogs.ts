import { ref, computed } from 'vue'
import type { LogEntry, LogCategory, LogLevel, ViewerRole } from '../types/systemLogs'
import { mockSystemLogs } from '../data/mockSystemLogs'

export function useSystemLogs() {
  const logs = ref<LogEntry[]>([...mockSystemLogs])
  const viewerRole = ref<ViewerRole>('SUPERADMIN')
  const activeTenantId = ref<string>('tenant_alpha')
  const categoryFilter = ref<LogCategory | 'ALL'>('ALL')
  const levelFilter = ref<LogLevel | 'ALL'>('ALL')
  const searchQuery = ref('')
  const selectedLog = ref<LogEntry | null>(null)

  const filteredLogs = computed(() => {
    return logs.value.filter(log => {
      // 1. Escopo de permissao: Superadmin ve tudo; Admin ve apenas seus clientes/usuarios
      if (viewerRole.value === 'ADMIN' && log.tenantId !== activeTenantId.value) {
        return false
      }
      // 2. Filtro de categoria
      if (categoryFilter.value !== 'ALL' && log.category !== categoryFilter.value) {
        return false
      }
      // 3. Filtro de nivel
      if (levelFilter.value !== 'ALL' && log.level !== levelFilter.value) {
        return false
      }
      // 4. Busca textual
      const q = searchQuery.value.trim().toLowerCase()
      if (!q) return true
      return (
        log.actor.toLowerCase().includes(q) ||
        log.action.toLowerCase().includes(q) ||
        log.target.toLowerCase().includes(q) ||
        log.details.toLowerCase().includes(q) ||
        log.ipAddress.toLowerCase().includes(q) ||
        log.tenantName.toLowerCase().includes(q)
      )
    })
  })

  const metrics = computed(() => ({
    total: filteredLogs.value.length,
    system: filteredLogs.value.filter(l => l.category === 'SYSTEM').length,
    audit: filteredLogs.value.filter(l => l.category === 'AUDIT').length,
    critical: filteredLogs.value.filter(l => l.level === 'CRITICAL').length,
    warning: filteredLogs.value.filter(l => l.level === 'WARNING').length
  }))

  const exportAsJson = () => {
    const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(filteredLogs.value, null, 2))
    const dlAnchor = document.createElement('a')
    dlAnchor.setAttribute('href', dataStr)
    dlAnchor.setAttribute('download', `hydravms_logs_${Date.now()}.json`)
    dlAnchor.click()
  }

  return {
    logs, viewerRole, activeTenantId, categoryFilter, levelFilter, searchQuery,
    selectedLog, filteredLogs, metrics, exportAsJson
  }
}
