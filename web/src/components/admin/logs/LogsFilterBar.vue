<script setup lang="ts">
import type { LogCategory, LogLevel, ViewerRole } from '../../../types/systemLogs'

defineProps<{
  viewerRole: ViewerRole
  activeTenantId: string
  categoryFilter: LogCategory | 'ALL'
  levelFilter: LogLevel | 'ALL'
  searchQuery: string
}>()

const emit = defineEmits<{
  (e: 'update:viewerRole', role: ViewerRole): void
  (e: 'update:activeTenantId', tenant: string): void
  (e: 'update:categoryFilter', cat: LogCategory | 'ALL'): void
  (e: 'update:levelFilter', lvl: LogLevel | 'ALL'): void
  (e: 'update:searchQuery', q: string): void
  (e: 'export'): void
}>()
</script>

<template>
  <div class="vms-card vms-flex-col" style="padding: 0.85rem 1rem; gap: 0.75rem; border: 1px solid var(--vms-border);">
    <div class="vms-flex-between" style="gap: 0.75rem; flex-wrap: wrap;">
      <!-- Search input -->
      <div style="flex: 1; min-width: 240px; position: relative;">
        <input
          :value="searchQuery"
          class="vms-auth-input"
          placeholder="Buscar por usuario, acao, IP ou detalhe..."
          @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        />
      </div>

      <!-- Role / Context Simulation Selector -->
      <div class="vms-flex-row" style="gap: 8px; align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PERFIL ATIVO:</span>
        <select
          :value="viewerRole === 'SUPERADMIN' ? 'SUPERADMIN' : activeTenantId"
          class="vms-auth-input"
          style="width: auto; padding: 4px 10px; font-weight: 700; color: var(--vms-neu-accent-orange);"
          @change="(e) => {
            const val = (e.target as HTMLSelectElement).value
            if (val === 'SUPERADMIN') {
              emit('update:viewerRole', 'SUPERADMIN')
            } else {
              emit('update:viewerRole', 'ADMIN')
              emit('update:activeTenantId', val)
            }
          }"
        >
          <option value="SUPERADMIN">[SUPERADMIN] ACESSO GLOBAL E SISTEMA</option>
          <option value="tenant_alpha">[ADMIN] ALPHA SEGURANCA // CLIENTES</option>
          <option value="tenant_beta">[ADMIN] BETA LOGISTICA // CLIENTES</option>
        </select>
      </div>
    </div>

    <!-- Secondary Filters: Categoria, Nivel, Exportar -->
    <div class="vms-flex-between" style="gap: 0.5rem; flex-wrap: wrap;">
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
        <!-- Category Filter -->
        <select
          :value="categoryFilter"
          class="vms-auth-input"
          style="width: auto; padding: 4px 8px;"
          @change="emit('update:categoryFilter', ($event.target as HTMLSelectElement).value as any)"
        >
          <option value="ALL">[CATEGORIA: TODAS]</option>
          <option value="SYSTEM">[SISTEMA] HARDWARE & SERVICOS</option>
          <option value="AUDIT">[AUDITORIA] ACOES DE USUARIOS</option>
        </select>

        <!-- Level Filter -->
        <select
          :value="levelFilter"
          class="vms-auth-input"
          style="width: auto; padding: 4px 8px;"
          @change="emit('update:levelFilter', ($event.target as HTMLSelectElement).value as any)"
        >
          <option value="ALL">[NIVEL: TODOS]</option>
          <option value="INFO">[INFO] REGULAR</option>
          <option value="WARNING">[ALERTA] ATENCAO</option>
          <option value="CRITICAL">[CRITICO] GRAVE</option>
        </select>
      </div>

      <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-weight: 700;" @click="emit('export')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"/></svg>
        <span>EXPORTAR JSON</span>
      </button>
    </div>
  </div>
</template>
