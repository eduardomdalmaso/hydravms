<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{ events: AnalyticEventRecord[] }>()
const emit = defineEmits<{ (e: 'inspect', event: AnalyticEventRecord): void }>()

const pageSize = ref(12)
const currentPage = ref(1)

const totalPages = computed(() => Math.max(1, Math.ceil(props.events.length / pageSize.value)))

watch([() => props.events.length, pageSize], () => {
  if (currentPage.value > totalPages.value) currentPage.value = 1
})

const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return props.events.slice(start, start + pageSize.value)
})

const formatTimestamp = (iso: string) => {
  if (!iso) return '--:--:--'
  try {
    const d = new Date(iso)
    return d.toLocaleString('pt-BR', { timeZone: 'UTC', dateStyle: 'short', timeStyle: 'medium' })
  } catch { return iso }
}

const getBadgeColor = (type: string) => {
  const t = (type || '').toUpperCase()
  if (t.includes('INTRUSÃO')) return '#ff003c'
  if (t.includes('MULTIDÃO')) return '#fcee0a'
  if (t.includes('VEÍCULO')) return '#00f0ff'
  return '#ff5e3a'
}

const getBadgeStyle = (type: string) => {
  const color = getBadgeColor(type)
  return { background: `${color}22`, color, border: `1px solid ${color}55` }
}
</script>

<template>
  <div v-if="events.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 4rem 2rem; border: 1px dashed rgba(255,255,255,0.08); border-radius: 6px; background: rgba(0,0,0,0.2);">
    <div class="vms-flex-col" style="gap: 0.5rem; align-items: center;">
      <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">// NENHUM EVENTO ENCONTRADO</span>
      <span class="vms-text-2xs vms-text-dim">Nenhum registro corresponde aos filtros de câmera, objeto ou período selecionados.</span>
    </div>
  </div>

  <div v-else class="vms-flex-col" style="gap: 1rem;">
    <!-- Grid de Cards Puros -->
    <div class="vms-explorer-events-grid">
      <div v-for="evt in paginatedEvents" :key="evt.id" class="vms-event-card" @click="emit('inspect', evt)">
        <div class="vms-event-crop-preview">
          <img v-if="evt.snapshot_url" :src="evt.snapshot_url" class="vms-event-thumb-img" alt="Event Crop" />
          <div v-else class="vms-event-thumb-placeholder">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="rgba(255, 94, 58, 0.4)" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
            <span class="vms-text-mono vms-text-2xs" style="color: #64748b;">SNAP CROP</span>
          </div>
          <svg v-if="evt.bbox" viewBox="0 0 100 100" preserveAspectRatio="none" style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; z-index: 5;">
            <rect :x="evt.bbox[0]" :y="evt.bbox[1]" :width="evt.bbox[2]" :height="evt.bbox[3]" fill="none" :stroke="getBadgeColor(evt.event_type)" stroke-width="1.2" stroke-dasharray="3, 1.5" />
          </svg>
          <span class="vms-badge" :style="[{ position: 'absolute', top: '6px', left: '6px', fontSize: '9px', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '4px' }, getBadgeStyle(evt.event_type)]">
            <svg v-if="evt.event_type.includes('MULTIDÃO') || evt.object_label.includes('multidão')" width="11" height="11" viewBox="0 0 24 24" fill="currentColor"><path d="M16.5 12c1.38 0 2.49-1.12 2.49-2.5S17.88 7 16.5 7C15.12 7 14 8.12 14 9.5s1.12 2.5 2.5 2.5zM9 11c1.66 0 2.99-1.34 2.99-3S10.66 5 9 5C7.34 5 6 6.34 6 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm7.5-1c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/></svg>
            <svg v-else width="11" height="11" viewBox="0 0 24 24" fill="currentColor"><circle cx="13.5" cy="4.5" r="2"/><path d="M13.8 8.2c-.4-.4-.9-.7-1.5-.7-.8 0-1.5.4-1.9 1l-2.9 3.9c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l2-2.7v4.5l-3.2 4.3c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l3.4-4.5 2.1 3v4.4c0 .6.4 1 1 1s1-.4 1-1v-5c0-.4-.2-.7-.5-.9l-2.4-3.4.5-4.4 2 1.5c.2.2.5.2.8.2.3 0 .6-.1.8-.3.4-.4.4-1 0-1.4l-3.2-2.3z"/></svg>
            <span>{{ evt.event_type.toUpperCase().includes('INTRUSÃO') ? 'INTRUSÃO' : evt.object_label.toUpperCase() }}</span>
          </span>
          <span class="vms-badge" :style="[{ position: 'absolute', top: '6px', right: '6px', fontSize: '8.5px', fontWeight: 'bold' }, getBadgeStyle(evt.event_type)]">{{ (evt.event_type || 'INTRUSÃO').toUpperCase() }}</span>
        </div>

        <div class="vms-flex-col" style="gap: 4px; padding: 0.25rem 0;">
          <span class="vms-text-xs vms-font-semibold" style="color: #ffffff; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ evt.camera_name }}</span>
          <span class="vms-text-2xs vms-text-dim" style="line-height: 1.35; height: 26px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical;">{{ evt.details }}</span>
        </div>

        <div class="vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.45rem; margin-top: auto; align-items: center;">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ formatTimestamp(evt.timestamp) }}</span>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 2px 6px;">INSPECIONAR</button>
        </div>
      </div>
    </div>

    <!-- Barra de Paginação HUD (Padrão 12 / 24 / 48) -->
    <div class="vms-flex-between" style="align-items: center; border-top: 1px solid rgba(255,255,255,0.08); padding-top: 0.75rem; flex-wrap: wrap; gap: 0.75rem;">
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ITENS / PÁGINA:</span>
        <button v-for="size in [12, 24, 48]" :key="size" class="vms-btn vms-btn-sm" :class="pageSize === size ? 'vms-btn-primary' : 'vms-btn-secondary'" style="font-size: 9.5px; padding: 3px 7px;" @click="pageSize = size">{{ size }}</button>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 4px 8px;" :disabled="currentPage <= 1" @click="currentPage--">ANTERIOR</button>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">PÁGINA {{ currentPage }} DE {{ totalPages }} // TOTAL: {{ events.length }}</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 4px 8px;" :disabled="currentPage >= totalPages" @click="currentPage++">PRÓXIMO</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-explorer-events-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.85rem; }
.vms-event-card { background: var(--vms-surface); border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.65rem; display: flex; flex-direction: column; gap: 0.4rem; cursor: pointer; transition: transform 0.15s, border-color 0.15s; }
.vms-event-card:hover { transform: translateY(-2px); border-color: rgba(255, 94, 58, 0.4); }
.vms-event-crop-preview { width: 100%; aspect-ratio: 16 / 9; background: #060910; border-radius: 4px; overflow: hidden; position: relative; display: flex; align-items: center; justify-content: center; }
.vms-event-thumb-img { width: 100%; height: 100%; object-fit: cover; display: block; }
.vms-event-thumb-placeholder { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; width: 100%; height: 100%; }
</style>
