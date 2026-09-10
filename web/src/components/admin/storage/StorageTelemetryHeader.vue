<script setup lang="ts">
import DatabaseStorageIcon from './DatabaseStorageIcon.vue'

defineProps<{
  totalCapacityGb: number
  totalUsedGb: number
  overallPercentage: number
  hotUsagePercent: number
}>()

const emit = defineEmits<{
  (e: 'openAddModal'): void
  (e: 'triggerDrain'): void
}>()

const formatGb = (gb: number) => gb >= 1024 ? `${(gb / 1024).toFixed(1)} TB` : `${gb} GB`
</script>

<template>
  <div class="vms-storage-telemetry-card">
    <div class="vms-flex-between" style="align-items: flex-start; gap: 1rem; flex-wrap: wrap;">
      <div class="vms-flex-col" style="gap: 4px;">
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-weight: 700;">
          [STORAGEGUARD // CAPACIDADE DE ARMAZENAMENTO]
        </span>
        <div class="vms-flex-row" style="gap: 12px; align-items: baseline;">
          <span class="vms-h2" style="color: #ffffff; font-family: var(--vms-font-jetbrains);">
            {{ formatGb(totalUsedGb) }}
          </span>
          <span class="vms-text-mono vms-text-xs vms-text-dim">
            DE {{ formatGb(totalCapacityGb) }} TOTAL ({{ overallPercentage }}% OCUPADO)
          </span>
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 8px;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-weight: 700;" title="Esvaziar buffer quente para os HDs" @click="emit('triggerDrain')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><path d="M12 2v20M17 17l-5 5-5-5"/></svg>
          <span>DRENAR BUFFER NVMe</span>
        </button>
        <button
          class="vms-btn vms-btn-ghost vms-btn-sm"
          style="padding: 6px 10px; display: inline-flex; align-items: center; gap: 5px; border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 4px; background: rgba(255, 94, 58, 0.08);"
          title="Adicionar Storage"
          @click="emit('openAddModal')"
        >
          <DatabaseStorageIcon :size="18" color="#ff5e3a" />
          <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="3.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </div>
    </div>

    <!-- Storage Watermark Gauge -->
    <div class="vms-flex-col" style="gap: 6px; margin-top: 12px;">
      <div class="vms-progress-track">
        <div class="vms-progress-fill" :style="{ width: `${Math.min(100, overallPercentage)}%` }"></div>
        <!-- Marcadores de Seguranca StorageGuard -->
        <div class="vms-marker mark-80" title="80% // Inicio de purga automatica"></div>
        <div class="vms-marker mark-95" title="95% // Circuit Breaker"></div>
      </div>
      <div class="vms-flex-between vms-text-mono vms-text-2xs vms-text-dim">
        <span>0 GB</span>
        <span style="color: rgba(255, 94, 58, 0.7);">[80% PURGA RECOMENDADA]</span>
        <span style="color: #ff5e3a;">[95% CIRCUIT BREAKER]</span>
        <span>{{ formatGb(totalCapacityGb) }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-storage-telemetry-card {
  background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg);
  padding: 1rem 1.25rem; display: flex; flex-direction: column;
}
.vms-progress-track {
  width: 100%; height: 10px; background: #07080c; border: 1px solid var(--vms-border);
  border-radius: 999px; position: relative; overflow: hidden;
}
.vms-progress-fill {
  height: 100%; background: linear-gradient(90deg, #ff5e3a 0%, #ffffff 100%);
  border-radius: 999px; transition: width 0.3s ease;
}
.vms-marker { position: absolute; top: 0; bottom: 0; width: 2px; z-index: 2; }
.mark-80 { left: 80%; background: rgba(255, 255, 255, 0.4); }
.mark-95 { left: 95%; background: #ff5e3a; }
</style>
