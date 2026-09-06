<script setup lang="ts">
import type { AnalyticsAlert } from '../../types/admin'

defineProps<{ alerts: AnalyticsAlert[] }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'acknowledge', id: string): void
  (e: 'clearAll'): void
}>()
</script>

<template>
  <aside class="vms-alerts-panel">
    <!-- Header -->
    <div class="vms-flex-between" style="padding: 0.75rem 1rem; border-bottom: 1px solid var(--vms-border); background: #0c0e14;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-sm vms-font-bold" style="color: var(--vms-neu-accent-orange);">ALERTAS DE IA</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">DETECCOES AO VIVO</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px; padding: 2px 6px;" @click="emit('clearAll')">[LIMPAR]</button>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 6px;" title="Fechar" @click="emit('close')">✕</button>
      </div>
    </div>

    <!-- Square Alerts Grid -->
    <TransitionGroup name="vms-alert" tag="div" class="vms-alerts-square-grid">
      <div
        v-for="alt in alerts"
        :key="alt.id"
        class="vms-alert-square-card"
        :class="[alt.severity, { acknowledged: alt.is_acknowledged }]"
        :title="alt.is_acknowledged ? 'Alerta reconhecido' : 'Clique para reconhecer'"
        @click="emit('acknowledge', alt.id)"
      >
        <!-- Top Row: Severidade & Hora -->
        <div class="vms-flex-between vms-text-mono vms-text-2xs" style="width: 100%;">
          <span class="vms-alert-square-badge" :class="alt.severity">
            {{ alt.severity.toUpperCase() }}
          </span>
          <span class="vms-text-dim" style="font-size: 10px;">{{ alt.timestamp.split(' ').pop() }}</span>
        </div>

        <!-- Center: Reticulo de IA Minimalista -->
        <div style="flex: 1; display: flex; align-items: center; justify-content: center; width: 100%;">
          <svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="1.6">
            <path d="M4 8V4m0 0h4M4 4l5 5m11-5h-4m4 0v4m0-4l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5h-4m4 0v-4m0 4l-5-5"/>
            <circle cx="12" cy="12" r="2" fill="#ff5e3a"/>
          </svg>
        </div>

        <!-- Bottom: Evento & Camera -->
        <div class="vms-flex-col" style="width: 100%; gap: 1px;">
          <span class="vms-font-bold vms-text-xs vms-truncate" style="color: #ffffff;" :title="alt.event_type">
            {{ alt.event_type }}
          </span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim vms-truncate" :title="alt.camera_name">
            {{ alt.camera_name }}
          </span>
        </div>
      </div>
    </TransitionGroup>
  </aside>
</template>
