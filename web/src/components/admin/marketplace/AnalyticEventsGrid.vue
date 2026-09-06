<script setup lang="ts">
import type { AnalyticEventRecord } from '../../../types/marketplace'

defineProps<{ events: AnalyticEventRecord[] }>()
const emit = defineEmits<{ (e: 'inspect', event: AnalyticEventRecord): void }>()
</script>

<template>
  <div v-if="events.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 3rem;">
    // NENHUM EVENTO REGISTRADO COM ESTES CRITÉRIOS DE BUSCA
  </div>

  <div v-else class="vms-explorer-events-grid">
    <div
      v-for="evt in events"
      :key="evt.id"
      class="vms-event-card"
      @click="emit('inspect', evt)"
    >
      <!-- Simulated Forensic Crop Box -->
      <div class="vms-event-crop-preview">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="rgba(255, 94, 58, 0.4)" stroke-width="1.5">
          <rect x="3" y="3" width="18" height="18" rx="2"/>
          <circle cx="8.5" cy="8.5" r="1.5"/>
          <polyline points="21 15 16 10 5 21"/>
        </svg>
        <span class="vms-text-mono vms-text-2xs" style="position: absolute; bottom: 6px; left: 8px; color: #64748b;">
          CROP // {{ evt.camera_id }}
        </span>
        <span class="vms-badge" :class="evt.severity === 'critical' ? 'vms-badge-danger' : (evt.severity === 'warning' ? 'vms-badge-orange' : 'vms-badge-blue')" style="position: absolute; top: 6px; right: 6px; font-size: 9px;">
          [{{ evt.severity.toUpperCase() }}]
        </span>
      </div>

      <div class="vms-flex-between" style="align-items: flex-start;">
        <div class="vms-flex-col" style="gap: 1px;">
          <span class="vms-text-xs vms-font-semibold" style="color: #ffffff;">{{ evt.object_label }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ evt.camera_name }}</span>
        </div>
        <span class="vms-text-mono vms-text-xs" style="color: #00f0ff; font-weight: bold;">
          {{ Math.round(evt.confidence * 100) }}%
        </span>
      </div>

      <p class="vms-text-2xs vms-text-dim" style="margin: 0; line-height: 1.3; height: 28px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;">
        {{ evt.details }}
      </p>

      <div class="vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.45rem; margin-top: auto;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ evt.timestamp.split('T')[1].replace('Z', '') }}</span>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 2px 6px;">INSPECIONAR</button>
      </div>
    </div>
  </div>
</template>
