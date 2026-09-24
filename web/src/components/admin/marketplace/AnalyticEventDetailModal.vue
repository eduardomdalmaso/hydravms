<script setup lang="ts">
import { ref, computed } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{
  isOpen: boolean
  event: AnalyticEventRecord | null
}>()

const emit = defineEmits<{ (e: 'close'): void }>()
const copied = ref(false)

const formattedJson = computed(() => {
  if (!props.event) return ''
  const standardPayload = {
    id: props.event.id,
    timestamp: props.event.timestamp,
    camera_id: props.event.camera_id,
    camera_name: props.event.camera_name,
    plugin_id: props.event.plugin_id,
    plugin_name: props.event.plugin_name,
    event_type: props.event.event_type,
    object_label: props.event.object_label,
    details: props.event.details,
    snapshot_url: props.event.snapshot_url || null,
    bbox: props.event.bbox || null,
    metadata: props.event.raw_payload || {}
  }
  return JSON.stringify(standardPayload, null, 2)
})

const copyJson = () => {
  if (!formattedJson.value) return
  navigator.clipboard.writeText(formattedJson.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
</script>

<template>
  <div v-if="isOpen && event" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content vms-event-detail-modal">
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">DETALHES FORENSES DO EVENTO</h4>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">ID: {{ event.id }} // {{ event.plugin_name }}</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 1rem; padding: 1rem 0; flex: 1; overflow-y: auto;">
        <!-- Header com Snapshot e Dados Principais -->
        <div class="vms-flex-row" style="gap: 1rem; align-items: center;">
          <div style="width: 140px; height: 95px; background: #060910; border-radius: 6px; border: 1px dashed rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden;">
            <img v-if="event.snapshot_url" :src="event.snapshot_url" style="width: 100%; height: 100%; object-fit: cover;" alt="Event Snapshot" />
            <div v-else class="vms-flex-col" style="align-items: center; gap: 4px;">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
              <span class="vms-text-mono vms-text-2xs" style="color: #64748b;">SNAP CROP</span>
            </div>
          </div>

          <div class="vms-flex-col" style="gap: 0.35rem; flex: 1;">
            <div class="vms-flex-between" style="align-items: center;">
              <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">{{ event.object_label.toUpperCase() }}</span>
              <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ (event.event_type || 'DETECCAO').toUpperCase() }}</span>
            </div>
            <span class="vms-text-xs vms-text-dim">{{ event.details }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">CÂMERA: {{ event.camera_name }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">TIMESTAMP UTC: {{ event.timestamp }}</span>
          </div>
        </div>

        <!-- JSON Payload Viewer -->
        <div class="vms-flex-col" style="gap: 0.35rem; flex: 1;">
          <div class="vms-flex-between" style="align-items: center;">
            <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">ESQUEMA JSON PADRÃO (GOLD LAYER)</span>
            <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 2px 8px;" @click="copyJson">
              {{ copied ? 'COPIADO!' : 'COPIAR JSON' }}
            </button>
          </div>
          <pre class="vms-json-payload-pre">{{ formattedJson }}</pre>
        </div>
      </div>

      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.75rem;">
        <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">[INTEGRIDADE VALIDADA] // REGISTRO IMUTÁVEL</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">FECHAR</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-event-detail-modal { width: 620px; max-width: 95vw; height: 500px; display: flex; flex-direction: column; background: #0b0e14; border: 1px solid rgba(255, 94, 58, 0.35); box-shadow: 0 20px 50px rgba(0,0,0,0.85); border-radius: 8px; }
.vms-json-payload-pre { flex: 1; background: #05070a; border: 1px solid rgba(0, 240, 255, 0.25); border-radius: 6px; padding: 0.75rem 1rem; color: #a5f3fc; font-family: var(--vms-font-mono); font-size: 11px; line-height: 1.45; overflow: auto; margin: 0; }
</style>
