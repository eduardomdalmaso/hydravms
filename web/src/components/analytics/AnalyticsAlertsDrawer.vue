<script setup lang="ts">
import type { AnalyticsAlert } from '../../types/admin'

const props = defineProps<{
  alerts: AnalyticsAlert[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'acknowledge', id: string): void
  (e: 'clearAll'): void
}>()
</script>

<template>
  <aside class="vms-alerts-panel">
    <!-- Header -->
    <div class="vms-flex-between" style="padding: 0.85rem 1rem; border-bottom: 1px solid var(--vms-border); background: #16191f;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-sm vms-font-semibold" style="color: var(--vms-neu-accent-orange);">ALERTAS DE ANALITICOS // IA</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">DETECCOES AO VIVO</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;" @click="emit('clearAll')">[LIMPAR]</button>
        <button class="vms-ubuntu-close-btn" style="position: relative; top: 0; right: 0;" title="Fechar" @click="emit('close')">
          <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <path d="M2 2L10 10M10 2L2 10" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Alerts List -->
    <div style="flex: 1; padding: 0.85rem; overflow-y: auto; display: flex; flex-direction: column; gap: 0.75rem;">
      <div
        v-for="alt in alerts"
        :key="alt.id"
        class="vms-alert-item"
        :class="alt.severity"
        :style="{ opacity: alt.is_acknowledged ? 0.55 : 1 }"
      >
        <div class="vms-flex-between">
          <span
            class="vms-badge"
            :class="{
              'vms-alert-badge-critical': alt.severity === 'critical',
              'vms-alert-badge-warning': alt.severity === 'warning',
              'vms-alert-badge-info': alt.severity === 'info'
            }"
          >
            [{{ alt.severity.toUpperCase() }}]
          </span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ alt.timestamp }}</span>
        </div>

        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">{{ alt.event_type }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ alt.camera_name }} // {{ Math.round(alt.confidence * 100) }}% CONF</span>
        </div>

        <!-- Snapshot Placeholder -->
        <div class="vms-alert-snapshot-placeholder">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// SNAPSHOT AI BOUNDING BOX</span>
        </div>

        <div class="vms-flex-between" style="margin-top: 0.25rem;">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">ID: {{ alt.id }}</span>
          <button
            v-if="!alt.is_acknowledged"
            class="vms-btn vms-btn-secondary vms-btn-sm"
            style="font-size: 11px; padding: 0.2rem 0.5rem;"
            @click="emit('acknowledge', alt.id)"
          >
            [RECONHECER]
          </button>
          <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">[RECONHECIDO]</span>
        </div>
      </div>
    </div>
  </aside>
</template>
