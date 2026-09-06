<script setup lang="ts">
import type { PluginManifest } from '../../../types/marketplace'
import PluginCard from './PluginCard.vue'

defineProps<{ plugins: PluginManifest[]; searchQuery: string; selectedCategory: string; statusFilter: string }>()

const emit = defineEmits<{
  (e: 'update:searchQuery', val: string): void; (e: 'update:selectedCategory', val: string): void
  (e: 'update:statusFilter', val: any): void; (e: 'install', id: string): void
  (e: 'uninstall', id: string): void; (e: 'toggle', id: string): void
  (e: 'details', plugin: PluginManifest): void
}>()

const categories = [
  { id: 'ALL', label: 'TODOS' },
  { id: 'traffic', label: 'TRÁFEGO & VEÍCULOS' },
  { id: 'access_control', label: 'ACESSO & FACIAL' },
  { id: 'safety', label: 'SEGURANÇA & EPI' },
  { id: 'analytics', label: 'PERÍMETRO & VMD' }
]
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Filter bar -->
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.75rem; align-items: center;">
      <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
        <button
          v-for="cat in categories"
          :key="cat.id"
          class="vms-btn vms-btn-sm"
          :class="selectedCategory === cat.id ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="font-size: 11px; padding: 4px 10px;"
          @click="emit('update:selectedCategory', cat.id)"
        >
          {{ cat.label }}
        </button>
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
          style="padding: 4px 10px; font-size: 11px; width: 180px;"
          placeholder="Buscar analítico..."
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
