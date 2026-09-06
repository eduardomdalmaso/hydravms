<script setup lang="ts">
import type { PluginManifest } from '../../../types/marketplace'

defineProps<{
  plugins: PluginManifest[]; selectedPluginId: string; selectedCameraId: string
  selectedSeverity: string; searchQuery: string; minConfidence: number; viewMode: 'table' | 'grid'
}>()

const emit = defineEmits<{
  (e: 'update:selectedPluginId', val: string): void; (e: 'update:selectedCameraId', val: string): void
  (e: 'update:selectedSeverity', val: string): void; (e: 'update:searchQuery', val: string): void
  (e: 'update:minConfidence', val: number): void; (e: 'update:viewMode', val: 'table' | 'grid'): void
  (e: 'export', format: 'csv' | 'json'): void
}>()

const cameras = [
  { id: 'ALL', name: '[CÂMERAS: TODAS]' },
  { id: 'cam_portaria_01', name: 'CAM_01 // Portaria Principal' },
  { id: 'cam_garagem_02', name: 'CAM_02 // Garagem Subsolo' },
  { id: 'cam_hall_03', name: 'CAM_03 // Catracas Recepção' },
  { id: 'cam_galpao_04', name: 'CAM_04 // Galpão Logística B' }
]
</script>

<template>
  <div class="vms-card" style="padding: 0.85rem 1.15rem; display: flex; flex-direction: column; gap: 0.75rem;">
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.75rem; align-items: center;">
      <!-- Filters Group -->
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap; align-items: center;">
        <select :value="selectedPluginId" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; width: 170px;" @change="emit('update:selectedPluginId', ($event.target as HTMLSelectElement).value)">
          <option value="ALL">[ANALÍTICO: TODOS]</option>
          <option v-for="p in plugins" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>

        <select :value="selectedCameraId" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; width: 170px;" @change="emit('update:selectedCameraId', ($event.target as HTMLSelectElement).value)">
          <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>

        <select :value="selectedSeverity" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; width: 140px;" @change="emit('update:selectedSeverity', ($event.target as HTMLSelectElement).value)">
          <option value="ALL">[SEVERIDADE: TODAS]</option>
          <option value="info">[NORMAL / INFO]</option>
          <option value="warning">[ALERTA / AVISO]</option>
          <option value="critical">[CRÍTICO / PERIGO]</option>
        </select>

        <input :value="searchQuery" class="vms-auth-input" style="padding: 4px 10px; font-size: 11px; width: 170px;" placeholder="Buscar placa, pessoa..." @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" />
      </div>

      <!-- Right controls & export -->
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-2xs vms-text-dim">CONF: {{ Math.round(minConfidence * 100) }}%</span>
          <input type="range" min="0.1" max="0.95" step="0.05" :value="minConfidence" style="width: 70px; accent-color: var(--vms-neu-accent-orange);" @input="emit('update:minConfidence', parseFloat(($event.target as HTMLInputElement).value))" />
        </div>

        <div class="vms-flex-row" style="gap: 2px;">
          <button class="vms-btn vms-btn-sm" :class="viewMode === 'table' ? 'vms-btn-primary' : 'vms-btn-secondary'" style="padding: 4px 8px; font-size: 10px;" @click="emit('update:viewMode', 'table')">TABELA</button>
          <button class="vms-btn vms-btn-sm" :class="viewMode === 'grid' ? 'vms-btn-primary' : 'vms-btn-secondary'" style="padding: 4px 8px; font-size: 10px;" @click="emit('update:viewMode', 'grid')">CARDS</button>
        </div>

        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('export', 'csv')">CSV</button>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('export', 'json')">JSON</button>
      </div>
    </div>
  </div>
</template>
