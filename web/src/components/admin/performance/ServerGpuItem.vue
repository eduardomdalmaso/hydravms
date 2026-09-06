<script setup lang="ts">
import type { GpuTelemetry } from '../../../types/performanceCluster'

const props = defineProps<{ gpu: GpuTelemetry }>()
const vramPercent = (used: number, total: number) => Math.round((used / total) * 100)
</script>

<template>
  <div class="vms-gpu-item">
    <div class="vms-flex-between" style="gap: 8px; flex-wrap: wrap;">
      <div class="vms-flex-row" style="gap: 6px; align-items: center;">
        <span class="vms-badge vms-badge-orange" style="font-size: 9px; padding: 2px 5px; font-weight: 700;">
          GPU {{ gpu.index }}
        </span>
        <span class="vms-font-semibold vms-text-xs" style="color: #ffffff;">{{ gpu.name }}</span>
      </div>
      <div class="vms-flex-row vms-text-mono vms-text-2xs" style="gap: 8px; color: var(--vms-neu-accent-orange);">
        <span>{{ gpu.tempC }}°C</span>
        <span>//</span>
        <span>{{ gpu.powerWatts }}W</span>
        <span>//</span>
        <span style="color: #ffffff;">COMPUTE: {{ gpu.computePercent }}%</span>
      </div>
    </div>

    <!-- VRAM Bar -->
    <div class="vms-flex-col" style="gap: 3px; margin-top: 6px;">
      <div class="vms-flex-between vms-text-mono vms-text-2xs">
        <span style="color: #ffffff;">VRAM: {{ gpu.vramUsedGb.toFixed(1) }} GB</span>
        <span class="vms-text-dim">TOTAL: {{ gpu.vramTotalGb.toFixed(1) }} GB ({{ vramPercent(gpu.vramUsedGb, gpu.vramTotalGb) }}%)</span>
      </div>
      <div class="vms-neu-track" style="height: 5px;">
        <div class="vms-neu-track-fill" :style="{ width: `${vramPercent(gpu.vramUsedGb, gpu.vramTotalGb)}%` }"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-gpu-item {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
  padding: 8px 10px;
  display: flex;
  flex-direction: column;
}
</style>
