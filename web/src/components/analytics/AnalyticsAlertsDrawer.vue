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
          <img v-if="alt.snapshot_url" :src="alt.snapshot_url" alt="Snapshot" class="vms-snapshot-img" />
          <div v-else class="vms-snapshot-cctv">
            <div class="vms-cctv-scanline"></div>
            <div class="vms-cctv-bbox" style="top: 24%; left: 34%; width: 32%; height: 56%;">
              <span class="vms-cctv-tag">{{ (alt.confidence * 100).toFixed(0) }}% // {{ alt.severity.toUpperCase() }}</span>
            </div>
            <span class="vms-cctv-hud-cam">{{ alt.camera_name }}</span>
            <span class="vms-cctv-hud-time">{{ alt.timestamp }}</span>
          </div>
        </div>

        <div class="vms-snapshot-label">
          {{ alt.event_type }}
        </div>
      </div>
    </TransitionGroup>
  </aside>
</template>
