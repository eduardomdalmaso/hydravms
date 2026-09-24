<script setup lang="ts">
import { ref, computed } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{
  events: AnalyticEventRecord[]; selectedCameraId: string; selectedObject: string;
  startDate: string; startTime: string; endDate: string; endTime: string;
}>()

const emit = defineEmits<{
  (e: 'update:selectedCameraId' | 'update:selectedObject' | 'update:startDate' | 'update:startTime' | 'update:endDate' | 'update:endTime', val: string): void
  (e: 'reset'): void; (e: 'export', format: 'csv' | 'json'): void
}>()

const isDatePickerOpen = ref(false)

const availableCameras = computed(() => {
  const map = new Map<string, string>()
  props.events.forEach(e => { if (e.camera_id) map.set(e.camera_id, e.camera_name || e.camera_id) })
  return Array.from(map.entries()).map(([id, name]) => ({ id, name }))
})

const availableObjects = computed(() => {
  const set = new Set<string>()
  props.events.forEach(e => { if (e.object_label) set.add(e.object_label) })
  return Array.from(set).sort()
})

const dateRangeSummary = computed(() => {
  if (!props.startDate && !props.endDate) return '[PERÍODO: TODOS]'
  const start = props.startDate ? `${props.startDate.slice(5)} ${props.startTime || '00:00'}` : 'INÍCIO'
  const end = props.endDate ? `${props.endDate.slice(5)} ${props.endTime || '23:59'}` : 'FIM'
  return `[${start} -> ${end}]`
})

const applyQuickRange = (type: 'today' | '24h' | 'all') => {
  const now = new Date()
  if (type === 'all') {
    emit('update:startDate', ''); emit('update:startTime', ''); emit('update:endDate', ''); emit('update:endTime', '')
    isDatePickerOpen.value = false
    return
  }
  const todayStr = now.toISOString().split('T')[0]
  emit('update:endDate', todayStr)
  emit('update:endTime', now.toTimeString().slice(0, 5))
  if (type === 'today') {
    emit('update:startDate', todayStr); emit('update:startTime', '00:00')
  } else if (type === '24h') {
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
    emit('update:startDate', yesterday.toISOString().split('T')[0])
    emit('update:startTime', yesterday.toTimeString().slice(0, 5))
  }
  isDatePickerOpen.value = false
}
</script>

<template>
  <div class="vms-card" style="padding: 0.75rem 1rem; position: relative;">
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.75rem; align-items: center;">
      <!-- Filtros Dinâmicos (Câmera & Objeto) -->
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap; align-items: center;">
        <select
          :value="selectedCameraId"
          class="vms-auth-input"
          style="padding: 5px 8px; font-size: 11px; width: 175px;"
          @change="emit('update:selectedCameraId', ($event.target as HTMLSelectElement).value)"
        >
          <option value="ALL">[CÂMERAS: TODAS]</option>
          <option v-for="c in availableCameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>

        <select
          :value="selectedObject"
          class="vms-auth-input"
          style="padding: 5px 8px; font-size: 11px; width: 155px;"
          @change="emit('update:selectedObject', ($event.target as HTMLSelectElement).value)"
        >
          <option value="ALL">[OBJETO: TODOS]</option>
          <option v-for="obj in availableObjects" :key="obj" :value="obj">{{ obj.toUpperCase() }}</option>
        </select>

        <!-- Botão Seletor de Período Temporal que abre Dropdown com 2 Calendários -->
        <div style="position: relative;">
          <button
            class="vms-btn vms-btn-secondary"
            style="font-size: 11px; padding: 5px 10px; font-family: var(--vms-font-mono); color: var(--vms-neu-accent-orange); display: flex; align-items: center; gap: 6px;"
            @click="isDatePickerOpen = !isDatePickerOpen"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            <span>{{ dateRangeSummary }}</span>
          </button>

          <!-- Popover / Dropdown com 2 Calendários e Horários -->
          <div
            v-if="isDatePickerOpen"
            class="vms-card"
            style="position: absolute; top: calc(100% + 6px); left: 0; z-index: 50; padding: 0.85rem; width: 440px; background: #0c0f17; border: 1px solid var(--vms-border); box-shadow: 0 10px 25px rgba(0,0,0,0.6); display: flex; flex-direction: column; gap: 0.75rem;"
          >
            <div class="vms-flex-between" style="border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 0.4rem;">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">FILTRO TEMPORAL // DATA & HORA</span>
              <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 1px 5px;" @click="isDatePickerOpen = false">FECHAR</button>
            </div>

            <div class="vms-flex-row" style="gap: 0.85rem;">
              <!-- Início -->
              <div class="vms-flex-col" style="gap: 4px; flex: 1;">
                <span class="vms-text-mono vms-text-2xs vms-text-dim">INÍCIO:</span>
                <input :value="startDate" type="date" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:startDate', ($event.target as HTMLInputElement).value)" />
                <input :value="startTime" type="time" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:startTime', ($event.target as HTMLInputElement).value)" />
              </div>

              <!-- Fim -->
              <div class="vms-flex-col" style="gap: 4px; flex: 1;">
                <span class="vms-text-mono vms-text-2xs vms-text-dim">FIM:</span>
                <input :value="endDate" type="date" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:endDate', ($event.target as HTMLInputElement).value)" />
                <input :value="endTime" type="time" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:endTime', ($event.target as HTMLInputElement).value)" />
              </div>
            </div>

            <!-- Atalhos Rápidos e Ações -->
            <div class="vms-flex-between" style="border-top: 1px solid rgba(255,255,255,0.08); padding-top: 0.5rem; align-items: center;">
              <div class="vms-flex-row" style="gap: 0.35rem;">
                <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9px; padding: 3px 6px;" @click="applyQuickRange('today')">HOJE</button>
                <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9px; padding: 3px 6px;" @click="applyQuickRange('24h')">ÚLTIMAS 24H</button>
                <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 3px 6px;" @click="applyQuickRange('all')">LIMPAR</button>
              </div>
              <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-size: 9.5px; padding: 3px 10px;" @click="isDatePickerOpen = false">APLICAR</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Exportação JSON / CSV -->
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 5px 9px;" @click="emit('export', 'json')">JSON</button>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 5px 9px;" @click="emit('export', 'csv')">CSV</button>
      </div>
    </div>
  </div>
</template>
