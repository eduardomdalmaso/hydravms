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
    <div v-if="plugins.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 3.5rem 2rem; border: 1px dashed rgba(255,255,255,0.08); border-radius: 6px; background: rgba(0,0,0,0.2);">
      <div class="vms-flex-col" style="gap: 0.5rem; align-items: center;">
        <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">// CATÁLOGO DE ANALÍTICOS VAZIO</span>
        <span class="vms-text-2xs vms-text-dim" style="max-width: 440px; line-height: 1.5;">
          Nenhum pacote de analítico disponível localmente. Novos modelos compilados aparecerão aqui quando exportados pelo HydraForge ou publicados no Hugging Face.
        </span>
      </div>
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
