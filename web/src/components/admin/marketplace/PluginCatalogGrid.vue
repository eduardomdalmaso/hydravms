<script setup lang="ts">
import type { PluginManifest } from '../../../types/marketplace'
import PluginCard from './PluginCard.vue'

defineProps<{ plugins: PluginManifest[]; searchQuery: string; statusFilter: string }>()

const emit = defineEmits<{
  (e: 'update:searchQuery', val: string): void
  (e: 'update:statusFilter', val: any): void
  (e: 'install', id: string): void
  (e: 'uninstall', id: string): void
  (e: 'toggle', id: string): void
  (e: 'details', plugin: PluginManifest): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Filter and Search bar -->
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.75rem; align-items: center;">
      <div class="vms-flex-row" style="gap: 0.65rem; align-items: center;">
        <span class="vms-text-xs vms-font-semibold" style="color: #ffffff; letter-spacing: 0.5px;">CATÁLOGO DE MODELOS SOTA</span>
        <span class="vms-badge vms-badge-secondary" style="font-size: 10px;">{{ plugins.length }} DISPONÍVEIS</span>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <select
          :value="statusFilter"
          class="vms-auth-input"
          style="padding: 4px 8px; font-size: 11px; width: 140px;"
          @change="emit('update:statusFilter', ($event.target as HTMLSelectElement).value)"
        >
          <option value="ALL">[STATUS: TODOS]</option>
          <option value="installed">[INSTALADOS]</option>
          <option value="available">[DISPONÍVEIS]</option>
        </select>

        <input
          :value="searchQuery"
          class="vms-auth-input"
          style="padding: 4px 10px; font-size: 11px; width: 220px;"
          placeholder="Buscar analítico ou modelo..."
          @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        />
      </div>
    </div>

    <!-- Plugins Grid -->
    <div v-if="plugins.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 3rem;">
      // NENHUM ANALÍTICO ENCONTRADO COM OS FILTROS SELECIONADOS
    </div>
    <div v-else class="vms-marketplace-grid">
      <PluginCard
        v-for="p in plugins"
        :key="p.id"
        :plugin="p"
        @install="emit('install', $event)"
        @uninstall="emit('uninstall', $event)"
        @toggle="emit('toggle', $event)"
        @details="emit('details', $event)"
      />
    </div>
  </div>
</template>
