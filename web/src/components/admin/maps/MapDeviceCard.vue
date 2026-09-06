<script setup lang="ts">
import type { DeployableDevice } from '../../../types/mapTree'

defineProps<{ device: DeployableDevice; isLocked?: boolean }>()
const emit = defineEmits<{
  (e: 'dragStart', dev: DeployableDevice): void
  (e: 'applyCoords', dev: DeployableDevice): void
}>()
</script>

<template>
  <div
    class="vms-map-device-card"
    :draggable="!isLocked"
    @dragstart="emit('dragStart', device)"
  >
    <div class="vms-flex-row" style="gap: 8px; align-items: center;">
      <!-- Ícone Câmera ou Alarme -->
      <div class="vms-device-icon" style="color: #ff5e3a;">
        <svg v-if="device.type === 'CAMERA'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M23 7l-7 5 7 5V7z"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/>
        </svg>
      </div>

      <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 0;">
        <span class="vms-font-semibold vms-text-xs" style="color: #ffffff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;" :title="device.name">
          {{ device.name }}
        </span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          {{ device.details || device.type }}
        </span>
      </div>
    </div>

    <div class="vms-flex-between" style="align-items: center; margin-top: 4px;">
      <span v-if="device.hasCoords" class="vms-badge" style="font-size: 8px; background: rgba(255, 94, 58, 0.12); color: #ff5e3a; border: 1px solid rgba(255, 94, 58, 0.4);">
        [COM COORDENADAS]
      </span>
      <span v-else class="vms-badge" style="font-size: 8px; background: #07080c; color: var(--vms-text-dim); border: 1px solid var(--vms-border);">
        [ARRASTE NO MAPA]
      </span>

      <button
        v-if="device.hasCoords"
        class="vms-btn vms-btn-ghost vms-btn-sm"
        style="padding: 2px 6px; font-size: 9px; color: #ff5e3a;"
        title="Posicionar nas coordenadas salvas"
        :disabled="isLocked"
        @click="emit('applyCoords', device)"
      >
        FIXAR COORDS
      </button>
    </div>
  </div>
</template>

<style scoped>
.vms-map-device-card {
  background: #0e1117; border: 1px solid var(--vms-border); border-radius: 6px;
  padding: 8px; display: flex; flex-direction: column; gap: 4px; cursor: grab;
  transition: border-color 0.15s ease, background 0.15s ease;
}
.vms-map-device-card:hover { border-color: var(--vms-neu-accent-orange); background: #14171d; }
.vms-map-device-card:active { cursor: grabbing; }
.vms-device-icon {
  width: 28px; height: 28px; border-radius: 4px; background: #07080c;
  display: flex; align-items: center; justify-content: center; border: 1px solid var(--vms-border);
}
</style>
