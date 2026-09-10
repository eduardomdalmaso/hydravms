<script setup lang="ts">
import { computed } from 'vue'
import type { StoragePoolItem } from '../../../types/storagePool'

const props = defineProps<{ pool: StoragePoolItem }>()
const emit = defineEmits<{ (e: 'remove', id: string): void }>()

const formatGb = (gb: number) => gb >= 1024 ? `${(gb / 1024).toFixed(1)} TB` : `${gb} GB`
const percent = (used: number, total: number) => Math.round((used / total) * 100)
</script>

<template>
  <div class="vms-storage-pool-card">
    <div class="vms-flex-between" style="align-items: flex-start; gap: 8px;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-weight: 700;">
          [{{ pool.role }}] // {{ pool.sourceType }}
        </span>
        <h4 class="vms-font-bold vms-text-sm" style="color: #ffffff; margin: 0;">{{ pool.name }}</h4>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ pool.nodeOrServer }}</span>
      </div>

      <div class="vms-flex-row" style="gap: 4px;">
        <button
          class="vms-btn vms-btn-ghost vms-btn-sm"
          style="padding: 3px 6px; color: var(--vms-neu-accent-orange);"
          title="Desanexar storage do HydraVMS"
          @click="emit('remove', pool.id)"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    </div>

    <!-- Path & Filesystem -->
    <div class="vms-storage-meta-box">
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="word-break: break-all;">
        {{ pool.pathOrEndpoint }}
      </span>
      <div class="vms-flex-between vms-text-mono vms-text-2xs" style="margin-top: 4px;">
        <span style="color: #ffffff;">FS: {{ pool.filesystem }}</span>
        <span style="color: var(--vms-neu-accent-orange);">STATUS: {{ pool.status }}</span>
      </div>
    </div>

    <!-- Usage bar -->
    <div class="vms-flex-col" style="gap: 4px;">
      <div class="vms-flex-between vms-text-mono vms-text-2xs">
        <span style="color: #ffffff;">USO: {{ formatGb(pool.usedGb) }}</span>
        <span class="vms-text-dim">TOTAL: {{ formatGb(pool.totalGb) }} ({{ percent(pool.usedGb, pool.totalGb) }}%)</span>
      </div>
      <div class="vms-usage-track">
        <div class="vms-usage-fill" :style="{ width: `${percent(pool.usedGb, pool.totalGb)}%` }"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-storage-pool-card {
  background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg);
  padding: 1rem; display: flex; flex-direction: column; gap: 10px; transition: border-color 0.2s ease;
}
.vms-storage-pool-card:hover { border-color: rgba(255, 94, 58, 0.5); }
.vms-storage-meta-box {
  background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 8px;
}
.vms-usage-track {
  width: 100%; height: 6px; background: #07080c; border: 1px solid var(--vms-border);
  border-radius: 999px; overflow: hidden;
}
.vms-usage-fill {
  height: 100%; background: #ff5e3a; border-radius: 999px; transition: width 0.3s ease;
}
</style>
