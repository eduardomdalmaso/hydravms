<script setup lang="ts">
import { computed } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{
  events: AnalyticEventRecord[]
  selectedCameraId: string
  selectedObject: string
  startDate: string
  startTime: string
  endDate: string
  endTime: string
}>()

const emit = defineEmits<{
  (e: 'update:selectedCameraId', val: string): void
  (e: 'update:selectedObject', val: string): void
  (e: 'update:startDate', val: string): void
  (e: 'update:startTime', val: string): void
  (e: 'update:endDate', val: string): void
  (e: 'update:endTime', val: string): void
  (e: 'reset'): void
  (e: 'export', format: 'csv' | 'json'): void
}>()

// Extrai dinamicamente as câmeras únicas presentes nos eventos salvos
const availableCameras = computed(() => {
  const map = new Map<string, string>()
  props.events.forEach(e => { if (e.camera_id) map.set(e.camera_id, e.camera_name || e.camera_id) })
  return Array.from(map.entries()).map(([id, name]) => ({ id, name }))
})

// Extrai dinamicamente os objetos únicos presentes nos eventos salvos
const availableObjects = computed(() => {
  const set = new Set<string>()
  props.events.forEach(e => { if (e.object_label) set.add(e.object_label) })
  return Array.from(set).sort()
})

const applyQuickRange = (type: 'today' | '24h' | 'all') => {
  const now = new Date()
  if (type === 'all') {
    emit('update:startDate', '')
    emit('update:startTime', '')
    emit('update:endDate', '')
    emit('update:endTime', '')
    return
  }
  const todayStr = now.toISOString().split('T')[0]
  emit('update:endDate', todayStr)
  emit('update:endTime', now.toTimeString().slice(0, 5))
  if (type === 'today') {
    emit('update:startDate', todayStr)
    emit('update:startTime', '00:00')
  } else if (type === '24h') {
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
    emit('update:startDate', yesterday.toISOString().split('T')[0])
    emit('update:startTime', yesterday.toTimeString().slice(0, 5))
  }
}
</script>

<template>
  <div class="vms-card" style="padding: 0.85rem 1.15rem; display: flex; flex-direction: column; gap: 0.75rem;">
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.75rem; align-items: center;">
      <!-- Filtros Dinâmicos (Câmera & Objeto baseados em registros) -->
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap; align-items: center;">
        <!-- Dropdown Dinâmico de Câmera -->
        <select
          :value="selectedCameraId"
          class="vms-auth-input"
          style="padding: 4px 8px; font-size: 11px; width: 170px;"
          @change="emit('update:selectedCameraId', ($event.target as HTMLSelectElement).value)"
        >
          <option value="ALL">[CÂMERAS: TODAS]</option>
          <option v-for="c in availableCameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>

        <!-- Dropdown Dinâmico de Objeto -->
        <select
          :value="selectedObject"
          class="vms-auth-input"
          style="padding: 4px 8px; font-size: 11px; width: 150px;"
          @change="emit('update:selectedObject', ($event.target as HTMLSelectElement).value)"
        >
          <option value="ALL">[OBJETO: TODOS]</option>
          <option v-for="obj in availableObjects" :key="obj" :value="obj">{{ obj.toUpperCase() }}</option>
        </select>
      </div>

      <!-- Filtro Temporal Dinâmico por Data e Hora -->
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center; flex-wrap: wrap;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">DE:</span>
        <input
          :value="startDate"
          type="date"
          class="vms-auth-input"
          style="padding: 3px 6px; font-size: 10.5px; width: 110px;"
          @input="emit('update:startDate', ($event.target as HTMLInputElement).value)"
        />
        <input
          :value="startTime"
          type="time"
          class="vms-auth-input"
          style="padding: 3px 6px; font-size: 10.5px; width: 75px;"
          @input="emit('update:startTime', ($event.target as HTMLInputElement).value)"
        />

        <span class="vms-text-mono vms-text-2xs vms-text-dim">ATÉ:</span>
        <input
          :value="endDate"
          type="date"
          class="vms-auth-input"
          style="padding: 3px 6px; font-size: 10.5px; width: 110px;"
          @input="emit('update:endDate', ($event.target as HTMLInputElement).value)"
        />
        <input
          :value="endTime"
          type="time"
          class="vms-auth-input"
          style="padding: 3px 6px; font-size: 10.5px; width: 75px;"
          @input="emit('update:endTime', ($event.target as HTMLInputElement).value)"
        />

        <!-- Botões de Atalho Temporal -->
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 3px 6px;" @click="applyQuickRange('today')">HOJE</button>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 3px 6px;" @click="applyQuickRange('24h')">24H</button>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9.5px; padding: 3px 6px;" @click="applyQuickRange('all')">LIMPAR</button>
      </div>

      <!-- Exportação JSON / CSV -->
      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('export', 'json')">JSON</button>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('export', 'csv')">CSV</button>
      </div>
    </div>
  </div>
</template>
