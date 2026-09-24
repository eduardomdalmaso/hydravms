<script setup lang="ts">
import type { AnalyticEventRecord } from '../../../types/marketplace'

defineProps<{ events: AnalyticEventRecord[] }>()
const emit = defineEmits<{ (e: 'inspect', event: AnalyticEventRecord): void }>()

const formatTimestamp = (iso: string) => {
  if (!iso) return '--:--:--'
  try {
    const d = new Date(iso)
    return d.toLocaleString('pt-BR', { timeZone: 'UTC', dateStyle: 'short', timeStyle: 'medium' })
  } catch {
    return iso
  }
}
</script>

<template>
  <div v-if="events.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 4rem 2rem; border: 1px dashed rgba(255,255,255,0.08); border-radius: 6px; background: rgba(0,0,0,0.2);">
    <div class="vms-flex-col" style="gap: 0.5rem; align-items: center;">
      <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">// NENHUM EVENTO ENCONTRADO</span>
      <span class="vms-text-2xs vms-text-dim">Nenhum registro corresponde aos filtros de câmera, objeto ou período selecionados.</span>
    </div>
  </div>

  <div v-else class="vms-explorer-events-grid">
    <div
      v-for="evt in events"
      :key="evt.id"
      class="vms-event-card"
      @click="emit('inspect', evt)"
    >
      <!-- Thumbnail / Crop Box Preview -->
      <div class="vms-event-crop-preview">
        <img v-if="evt.snapshot_url" :src="evt.snapshot_url" class="vms-event-thumb-img" alt="Event Crop" />
        <div v-else class="vms-event-thumb-placeholder">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="rgba(255, 94, 58, 0.4)" stroke-width="1.5">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
          </svg>
          <span class="vms-text-mono vms-text-2xs" style="color: #64748b;">SNAP CROP</span>
        </div>

        <!-- Tag do Objeto Detectado -->
        <span class="vms-badge vms-badge-orange" style="position: absolute; top: 6px; left: 6px; font-size: 9px; font-weight: bold;">
          [{{ evt.object_label.toUpperCase() }}]
        </span>

        <!-- Tag do Modo / Tipo de Evento -->
        <span class="vms-badge vms-badge-secondary" style="position: absolute; top: 6px; right: 6px; font-size: 8.5px;">
          {{ (evt.event_type || 'DETECCAO').toUpperCase() }}
        </span>
      </div>

      <!-- Informações do Card -->
      <div class="vms-flex-col" style="gap: 4px; padding: 0.25rem 0;">
        <div class="vms-flex-between" style="align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: #ffffff; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ evt.camera_name }}
          </span>
        </div>

        <span class="vms-text-2xs vms-text-dim" style="line-height: 1.35; height: 26px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical;">
          {{ evt.details }}
        </span>
      </div>

      <!-- Rodapé com Timestamp e Ação -->
      <div class="vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.45rem; margin-top: auto; align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          {{ formatTimestamp(evt.timestamp) }}
        </span>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 2px 6px;">
          INSPECIONAR
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-explorer-events-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.85rem; }
.vms-event-card { background: var(--vms-surface); border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.65rem; display: flex; flex-direction: column; gap: 0.4rem; cursor: pointer; transition: transform 0.15s, border-color 0.15s; }
.vms-event-card:hover { transform: translateY(-2px); border-color: rgba(255, 94, 58, 0.4); }
.vms-event-crop-preview { width: 100%; height: 110px; background: #060910; border-radius: 4px; overflow: hidden; position: relative; display: flex; align-items: center; justify-content: center; }
.vms-event-thumb-img { width: 100%; height: 100%; object-fit: cover; }
.vms-event-thumb-placeholder { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; width: 100%; height: 100%; }
</style>
