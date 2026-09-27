<script setup lang="ts">
import type { AnalyticsAlert } from '../../types/admin'

defineProps<{ alerts: AnalyticsAlert[] }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'acknowledge', id: string): void
  (e: 'clearAll'): void
}>()

const getEventColor = (type: string) => {
  if (type.includes('INTRUSÃO') || type.includes('OFFLINE')) return '#ff003c'
  if (type.includes('MULTIDÃO')) return '#fcee0a'
  if (type.includes('VEÍCULO')) return '#00f0ff'
  return '#ff5e3a'
}
</script>

<template>
  <aside class="vms-alerts-panel">
    <!-- Header: Apenas ALERTA -->
    <div class="vms-flex-between" style="padding: 0.75rem 1rem; border-bottom: 1px solid var(--vms-border); background: #0c0e14;">
      <span class="vms-text-sm vms-font-bold" style="color: var(--vms-neu-accent-orange); letter-spacing: 0.5px;">ALERTA</span>
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button v-if="alerts.length > 0" class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px; padding: 2px 6px;" @click="emit('clearAll')">[LIMPAR]</button>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    </div>

    <!-- Sequential Snapshots Feed -->
    <div v-if="alerts.length === 0" class="vms-flex-col vms-flex-center" style="padding: 3rem 1rem; color: var(--vms-text-dim); text-align: center; gap: 0.5rem;">
      <span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-regular);">[SEM ALERTAS ATIVOS]</span>
      <span class="vms-text-2xs" style="opacity: 0.6;">Aguardando eventos da IA</span>
    </div>

    <TransitionGroup v-else name="vms-alert" tag="div" class="vms-alerts-snapshots-feed">
      <div
        v-for="alt in alerts"
        :key="alt.id"
        class="vms-alert-snapshot-card"
        :class="{ acknowledged: alt.is_acknowledged }"
        :title="alt.is_acknowledged ? 'Alerta reconhecido' : 'Clique para reconhecer'"
        @click="emit('acknowledge', alt.id)"
      >
        <!-- Snapshot Frame -->
        <div class="vms-snapshot-frame">
          <img
            :src="alt.snapshot_url || `http://localhost:8080/api/v1/streams/${alt.camera_id || 'cam_01'}/snapshot`"
            alt="Snapshot"
            class="vms-snapshot-img"
            @error="($event.target as HTMLElement).style.display = 'none'"
          />
          <svg
            v-if="alt.bbox"
            viewBox="0 0 100 100"
            preserveAspectRatio="none"
            style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; z-index: 5;"
          >
            <rect
              :x="alt.bbox[0]"
              :y="alt.bbox[1]"
              :width="alt.bbox[2]"
              :height="alt.bbox[3]"
              fill="none"
              :stroke="getEventColor(alt.event_type)"
              stroke-width="1.2"
              stroke-dasharray="3, 1.5"
            />
          </svg>
          <div class="vms-snapshot-hud-badge" :style="{ background: getEventColor(alt.event_type), color: '#000' }">
            {{ (alt.confidence * 100).toFixed(0) }}% // {{ alt.severity.toUpperCase() }}
          </div>
          <span class="vms-cctv-hud-cam">{{ alt.camera_name }}</span>
          <span class="vms-cctv-hud-time">{{ alt.timestamp }}</span>
        </div>

        <div class="vms-snapshot-label" :style="{ color: getEventColor(alt.event_type), display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '6px' }">
          <svg v-if="alt.event_type.includes('MULTIDÃO')" width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><path d="M16.5 12c1.38 0 2.49-1.12 2.49-2.5S17.88 7 16.5 7C15.12 7 14 8.12 14 9.5s1.12 2.5 2.5 2.5zM9 11c1.66 0 2.99-1.34 2.99-3S10.66 5 9 5C7.34 5 6 6.34 6 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm7.5-1c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/></svg>
          <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><circle cx="13.5" cy="4.5" r="2"/><path d="M13.8 8.2c-.4-.4-.9-.7-1.5-.7-.8 0-1.5.4-1.9 1l-2.9 3.9c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l2-2.7v4.5l-3.2 4.3c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l3.4-4.5 2.1 3v4.4c0 .6.4 1 1 1s1-.4 1-1v-5c0-.4-.2-.7-.5-.9l-2.4-3.4.5-4.4 2 1.5c.2.2.5.2.8.2.3 0 .6-.1.8-.3.4-.4.4-1 0-1.4l-3.2-2.3z"/></svg>
          <span>{{ alt.event_type }}</span>
        </div>
      </div>
    </TransitionGroup>
  </aside>
</template>
