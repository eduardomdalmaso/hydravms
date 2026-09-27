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
            <svg v-if="event.bbox" viewBox="0 0 100 100" preserveAspectRatio="none" style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; z-index: 5;">
              <rect :x="event.bbox[0]" :y="event.bbox[1]" :width="event.bbox[2]" :height="event.bbox[3]" fill="none" stroke="#ff003c" stroke-width="1.2" stroke-dasharray="3, 1.5" />
            </svg>
          </div>

          <div class="vms-flex-col" style="gap: 0.35rem; flex: 1;">
            <div class="vms-flex-between" style="align-items: center;">
              <div class="vms-flex-row" style="align-items: center; gap: 5px;">
                <svg v-if="event.event_type.includes('MULTIDÃO') || event.object_label.includes('multidão')" width="14" height="14" viewBox="0 0 24 24" fill="#fcee0a"><path d="M16.5 12c1.38 0 2.49-1.12 2.49-2.5S17.88 7 16.5 7C15.12 7 14 8.12 14 9.5s1.12 2.5 2.5 2.5zM9 11c1.66 0 2.99-1.34 2.99-3S10.66 5 9 5C7.34 5 6 6.34 6 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm7.5-1c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/></svg>
                <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="#ff003c"><circle cx="13.5" cy="4.5" r="2"/><path d="M13.8 8.2c-.4-.4-.9-.7-1.5-.7-.8 0-1.5.4-1.9 1l-2.9 3.9c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l2-2.7v4.5l-3.2 4.3c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l3.4-4.5 2.1 3v4.4c0 .6.4 1 1 1s1-.4 1-1v-5c0-.4-.2-.7-.5-.9l-2.4-3.4.5-4.4 2 1.5c.2.2.5.2.8.2.3 0 .6-.1.8-.3.4-.4.4-1 0-1.4l-3.2-2.3z"/></svg>
                <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">{{ event.event_type.toUpperCase().includes('INTRUSÃO') ? 'INTRUSÃO' : event.object_label.toUpperCase() }}</span>
              </div>
              <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ (event.event_type || 'INTRUSÃO').toUpperCase() }}</span>
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
