<script setup lang="ts">
import type { ServerNodeItem } from '../../../types/performanceCluster'

defineProps<{ node: ServerNodeItem }>()
const ramPercent = (used: number, total: number) => Math.round((used / total) * 100)
const vramPercent = (used: number, total: number) => Math.round((used / total) * 100)

const roleShort = (role: string) => {
  if (role === 'MASTER_VMS') return 'VMS CORE'
  if (role === 'HYDRASTREAM_EDGE') return 'EDGE INGEST'
  return 'GPU WORKER'
}
</script>

<template>
  <div class="vms-server-row">
    <!-- Coluna 1: Identificacao do Servidor -->
    <div class="vms-flex-col" style="min-width: 220px; flex: 1.2; gap: 2px;">
      <div class="vms-flex-row" style="gap: 6px; align-items: center;">
        <span class="vms-status-led online" style="margin-right: 2px;"></span>
        <span class="vms-font-bold vms-text-sm" style="color: #ffffff;">{{ node.hostname }}</span>
        <span class="vms-badge vms-badge-orange" style="font-size: 9px; padding: 1px 5px;">
          [{{ roleShort(node.role) }}]
        </span>
      </div>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">IP: {{ node.ip }} // {{ node.uptime }}</span>
    </div>

    <!-- Coluna 2: CPU em Porcentagem -->
    <div class="vms-flex-col" style="min-width: 140px; flex: 1; gap: 3px;">
      <div class="vms-flex-between vms-text-mono vms-text-2xs">
        <span class="vms-text-dim">CPU:</span>
        <span class="vms-font-bold" style="color: #ffffff;">{{ node.cpuPercent }}%</span>
      </div>
      <div class="vms-neu-track" style="height: 5px;">
        <div class="vms-neu-track-fill" :style="{ width: `${node.cpuPercent}%` }"></div>
      </div>
    </div>

    <!-- Coluna 3: RAM Usada -->
    <div class="vms-flex-col" style="min-width: 170px; flex: 1.2; gap: 3px;">
      <div class="vms-flex-between vms-text-mono vms-text-2xs">
        <span class="vms-text-dim">RAM USADA:</span>
        <span class="vms-font-bold" style="color: #ffffff;">
          {{ node.ramUsedGb.toFixed(1) }} / {{ node.ramTotalGb }} GB ({{ ramPercent(node.ramUsedGb, node.ramTotalGb) }}%)
        </span>
      </div>
      <div class="vms-neu-track" style="height: 5px;">
        <div class="vms-neu-track-fill" :style="{ width: `${ramPercent(node.ramUsedGb, node.ramTotalGb)}%` }"></div>
      </div>
    </div>

    <!-- Coluna 4: VRAM das GPUs por GPU -->
    <div class="vms-flex-col" style="min-width: 280px; flex: 2; gap: 4px;">
      <div v-if="node.gpus.length > 0" class="vms-flex-col" style="gap: 4px;">
        <div v-for="gpu in node.gpus" :key="gpu.id" class="vms-gpu-mini-item">
          <div class="vms-flex-between vms-text-mono vms-text-2xs">
            <span style="color: var(--vms-neu-accent-orange); font-weight: 700;">
              GPU {{ gpu.index }}: {{ gpu.name }}
            </span>
            <span style="color: #ffffff;">
              VRAM: {{ gpu.vramUsedGb.toFixed(1) }} / {{ gpu.vramTotalGb.toFixed(1) }} GB ({{ vramPercent(gpu.vramUsedGb, gpu.vramTotalGb) }}%)
            </span>
          </div>
          <div class="vms-neu-track" style="height: 4px; margin-top: 2px;">
            <div class="vms-neu-track-fill" :style="{ width: `${vramPercent(gpu.vramUsedGb, gpu.vramTotalGb)}%` }"></div>
          </div>
        </div>
      </div>
      <div v-else class="vms-text-mono vms-text-2xs vms-text-dim" style="padding: 4px 0;">
        [SEM GPU // PROCESSAMENTO CPU]
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-server-row {
  background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-md);
  padding: 0.85rem 1.1rem; display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap;
  transition: border-color 0.2s ease, background 0.2s ease;
}
.vms-server-row:hover { border-color: rgba(255, 94, 58, 0.4); background: #121620; }
.vms-gpu-mini-item {
  background: rgba(0, 0, 0, 0.35); border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 4px; padding: 4px 8px;
}
</style>
