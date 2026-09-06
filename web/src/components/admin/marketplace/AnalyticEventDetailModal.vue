<script setup lang="ts">
import type { AnalyticEventRecord } from '../../../types/marketplace'

defineProps<{
  isOpen: boolean
  event: AnalyticEventRecord | null
}>()

const emit = defineEmits<{ (e: 'close'): void }>()
</script>

<template>
  <div v-if="isOpen && event" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content" style="max-width: 600px; width: 100%;">
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">DETALHES FORENSES DO EVENTO</h4>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">ID: {{ event.id }} // {{ event.plugin_name }}</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 1rem; padding: 1.25rem 0;">
        <div class="vms-flex-row" style="gap: 1rem; align-items: center;">
          <div style="width: 140px; height: 100px; background: #060910; border-radius: 4px; border: 1px dashed rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center; position: relative;">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
            <span class="vms-text-mono vms-text-2xs" style="position: absolute; bottom: 4px; color: #64748b;">SNAP CROP</span>
          </div>

          <div class="vms-flex-col" style="gap: 0.35rem; flex: 1;">
            <div class="vms-flex-between">
              <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">{{ event.object_label }}</span>
              <span class="vms-badge" :class="event.severity === 'critical' ? 'vms-badge-danger' : (event.severity === 'warning' ? 'vms-badge-orange' : 'vms-badge-blue')" style="font-size: 10px;">
                [{{ event.severity.toUpperCase() }}]
              </span>
            </div>
            <span class="vms-text-xs vms-text-dim">{{ event.details }}</span>
            <span class="vms-text-mono vms-text-2xs" style="color: #00f0ff;">CONFIANÇA: {{ Math.round(event.confidence * 100) }}% // {{ event.camera_name }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">TIMESTAMP: {{ event.timestamp }}</span>
          </div>
        </div>

        <div class="vms-flex-col" style="gap: 0.35rem;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">PAYLOAD METADATA (JSON RAW)</span>
          <pre class="vms-text-mono vms-text-2xs" style="background: #080a0f; padding: 0.85rem; border-radius: 4px; border: 1px solid var(--vms-border); overflow-x: auto; color: #a5f3fc; margin: 0;">{{ JSON.stringify(event.raw_payload, null, 2) }}</pre>
        </div>
      </div>

      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.85rem;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">HASH FORENSE: SHA-256 VALIDADO</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">FECHAR</button>
      </div>
    </div>
  </div>
</template>
