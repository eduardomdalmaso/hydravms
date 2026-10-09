<script setup lang="ts">
import { ref, computed } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{
  events: AnalyticEventRecord[]; selectedCameraId: string; selectedAgeGroup: string;
  searchKnownFace: string; selectedFaceType: 'ALL' | 'KNOWN' | 'UNKNOWN';
  startDate: string; startTime: string; endDate: string; endTime: string;
}>()

const emit = defineEmits<{
  (e: 'update:selectedCameraId' | 'update:selectedAgeGroup' | 'update:searchKnownFace' | 'update:startDate' | 'update:startTime' | 'update:endDate' | 'update:endTime', val: string): void
  (e: 'update:selectedFaceType', val: 'ALL' | 'KNOWN' | 'UNKNOWN'): void
  (e: 'reset'): void; (e: 'export', format: 'csv' | 'json'): void
}>()

const isDatePickerOpen = ref(false)

const availableCameras = computed(() => {
  const map = new Map<string, string>()
  props.events.forEach(e => { if (e.camera_id) map.set(e.camera_id, e.camera_name || e.camera_id) })
  return Array.from(map.entries()).map(([id, name]) => ({ id, name }))
})

const ageGroups = [
  { id: 'ALL', label: '[IDADE: TODAS]' },
  { id: 'child', label: 'CRIANÇA (0-12)' },
  { id: 'youth', label: 'JOVEM (13-24)' },
  { id: 'adult', label: 'ADULTO (25-59)' },
  { id: 'senior', label: 'IDOSO (60+)' }
]

const dateRangeSummary = computed(() => {
  if (!props.startDate && !props.endDate) return '[PERÍODO: TODOS]'
  const s = props.startDate ? `${props.startDate.slice(5)} ${props.startTime || '00:00'}` : 'INÍCIO'
  const end = props.endDate ? `${props.endDate.slice(5)} ${props.endTime || '23:59'}` : 'FIM'
  return `[${s} -> ${end}]`
})

const applyQuickRange = (type: 'today' | '24h' | 'all') => {
  const now = new Date()
  if (type === 'all') {
    emit('update:startDate', ''); emit('update:startTime', ''); emit('update:endDate', ''); emit('update:endTime', '')
    isDatePickerOpen.value = false; return
  }
  const todayStr = now.toISOString().split('T')[0]
  emit('update:endDate', todayStr); emit('update:endTime', now.toTimeString().slice(0, 5))
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
    <div class="vms-flex-between" style="flex-wrap: wrap; gap: 0.65rem; align-items: center;">
      <!-- Filtros Biométricos (Câmera, Face Conhecida, Idade) -->
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap; align-items: center; flex: 1;">
        <!-- Seletor de Câmeras -->
        <select :value="selectedCameraId" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px; width: 150px;" @change="emit('update:selectedCameraId', ($event.target as HTMLSelectElement).value)">
          <option value="ALL">[CÂMERAS: TODAS]</option>
          <option v-for="c in availableCameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>

        <!-- Campo de Busca por Face Conhecida (Nome) -->
        <div style="position: relative; width: 190px;">
          <input :value="searchKnownFace" type="text" class="vms-auth-input" style="padding: 5px 8px 5px 24px; font-size: 11px; width: 100%; box-sizing: border-box;" placeholder="Buscar face conhecida..." @input="emit('update:searchKnownFace', ($event.target as HTMLInputElement).value)" />
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="position: absolute; left: 7px; top: 8px; color: var(--vms-neu-accent-orange); pointer-events: none;"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </div>

        <!-- Status de Reconhecimento: Todos, Reconhecido, Desconhecido -->
        <select :value="selectedFaceType" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px; width: 150px;" @change="emit('update:selectedFaceType', ($event.target as HTMLSelectElement).value as any)">
          <option value="ALL">[STATUS: TODOS]</option>
          <option value="KNOWN">[RECONHECIDO]</option>
          <option value="UNKNOWN">[DESCONHECIDO]</option>
        </select>

        <!-- Faixa Etária / Idade -->
        <select :value="selectedAgeGroup" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px; width: 140px;" @change="emit('update:selectedAgeGroup', ($event.target as HTMLSelectElement).value)">
          <option v-for="ag in ageGroups" :key="ag.id" :value="ag.id">{{ ag.label }}</option>
        </select>

        <!-- Período Temporal -->
        <div style="position: relative;">
          <button class="vms-btn vms-btn-secondary" style="font-size: 11px; padding: 5px 10px; font-family: var(--vms-font-mono); color: var(--vms-neu-accent-orange); display: flex; align-items: center; gap: 6px;" @click="isDatePickerOpen = !isDatePickerOpen">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            <span>{{ dateRangeSummary }}</span>
          </button>

          <div v-if="isDatePickerOpen" class="vms-card" style="position: absolute; top: calc(100% + 6px); left: 0; z-index: 50; padding: 0.85rem; width: 420px; background: #0c0f17; border: 1px solid var(--vms-border); box-shadow: 0 10px 25px rgba(0,0,0,0.6); display: flex; flex-direction: column; gap: 0.75rem;">
            <div class="vms-flex-between" style="border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 0.4rem;">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">FILTRO TEMPORAL // DATA & HORA</span>
              <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 1px 5px;" @click="isDatePickerOpen = false">FECHAR</button>
            </div>
            <div class="vms-flex-row" style="gap: 0.85rem;">
              <div class="vms-flex-col" style="gap: 4px; flex: 1;">
                <span class="vms-text-mono vms-text-2xs vms-text-dim">INÍCIO:</span>
                <input :value="startDate" type="date" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:startDate', ($event.target as HTMLInputElement).value)" />
                <input :value="startTime" type="time" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:startTime', ($event.target as HTMLInputElement).value)" />
              </div>
              <div class="vms-flex-col" style="gap: 4px; flex: 1;">
                <span class="vms-text-mono vms-text-2xs vms-text-dim">FIM:</span>
                <input :value="endDate" type="date" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:endDate', ($event.target as HTMLInputElement).value)" />
                <input :value="endTime" type="time" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:endTime', ($event.target as HTMLInputElement).value)" />
              </div>
            </div>
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
