<script setup lang="ts">
import { useSystemLogs } from '../../composables/useSystemLogs'
import LogsFilterBar from '../../components/admin/logs/LogsFilterBar.vue'
import LogsTable from '../../components/admin/logs/LogsTable.vue'
import LogsDetailModal from '../../components/admin/logs/LogsDetailModal.vue'

const {
  viewerRole, activeTenantId, categoryFilter, levelFilter, searchQuery,
  selectedLog, filteredLogs, metrics, exportAsJson
} = useSystemLogs()
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1.25rem;">
    <!-- Telemetry Header -->
    <div class="vms-card vms-flex-between" style="padding: 1rem 1.25rem; border: 1px solid var(--vms-border); flex-wrap: wrap; gap: 1rem;">
      <div class="vms-flex-col" style="gap: 4px;">
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-weight: 700;">
          [REGISTRO FORENSE & AUDITORIA // LOGS DE SISTEMA E CLIENTES]
        </span>
        <div class="vms-flex-row" style="gap: 10px; align-items: baseline;">
          <span class="vms-h2" style="color: #ffffff; font-family: var(--vms-font-jetbrains);">
            {{ metrics.total }} REGISTROS
          </span>
          <span class="vms-text-mono vms-text-xs vms-text-dim">
            ({{ metrics.system }} SISTEMA // {{ metrics.audit }} AUDITORIA // {{ metrics.critical }} CRITICOS)
          </span>
        </div>
      </div>

      <!-- Scope Indicator Badge -->
      <div class="vms-flex-row" style="gap: 8px; align-items: center;">
        <span class="vms-badge" :class="viewerRole === 'SUPERADMIN' ? 'vms-badge-orange' : 'vms-badge-secondary'" style="font-weight: 700;">
          {{ viewerRole === 'SUPERADMIN' ? '[ACESSO: SUPERADMIN // GLOBAL]' : `[ACESSO: ADMIN // ${activeTenantId.toUpperCase()}]` }}
        </span>
      </div>
    </div>

    <!-- Filter Bar -->
    <LogsFilterBar
      :viewerRole="viewerRole"
      :activeTenantId="activeTenantId"
      :categoryFilter="categoryFilter"
      :levelFilter="levelFilter"
      :searchQuery="searchQuery"
      @update:viewerRole="viewerRole = $event"
      @update:activeTenantId="activeTenantId = $event"
      @update:categoryFilter="categoryFilter = $event"
      @update:levelFilter="levelFilter = $event"
      @update:searchQuery="searchQuery = $event"
      @export="exportAsJson"
    />

    <!-- Notice if Admin Scope -->
    <div v-if="viewerRole === 'ADMIN'" class="vms-card" style="padding: 8px 12px; background: rgba(255, 94, 58, 0.08); border: 1px solid var(--vms-neu-accent-orange); font-size: 0.75rem; color: #ffffff;">
      <span class="vms-text-mono" style="color: var(--vms-neu-accent-orange); font-weight: 700;">[ESCOPO RESTRITO] </span>
      Visualizacao filtrada estritamente para auditoria de acoes dos clientes e operadores do seu tenant.
    </div>

    <!-- Table or Empty State -->
    <LogsTable
      v-if="filteredLogs.length > 0"
      :logs="filteredLogs"
      @selectLog="selectedLog = $event"
    />
    <div v-else class="vms-card vms-flex-col vms-flex-center" style="padding: 3rem; text-align: center; gap: 0.5rem;">
      <span class="vms-text-muted vms-text-mono vms-text-sm">[NENHUM LOG LOCALIZADO COM OS FILTROS APLICADOS]</span>
      <span class="vms-text-2xs vms-text-dim">Ajuste os filtros de busca, categoria, nivel ou selecione outro perfil.</span>
    </div>

    <!-- Forensic Detail Modal -->
    <LogsDetailModal
      v-if="selectedLog"
      :log="selectedLog"
      @close="selectedLog = null"
    />
  </div>
</template>
