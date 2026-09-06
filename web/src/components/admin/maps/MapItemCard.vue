<script setup lang="ts">
import type { EnterpriseMapItem } from '../../../types/mapTree'

defineProps<{ mapItem: EnterpriseMapItem }>()
const emit = defineEmits<{
  (e: 'open', m: EnterpriseMapItem): void
  (e: 'dragStart', m: EnterpriseMapItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragStart', mapItem)"
    @click="emit('open', mapItem)"
  >
    <!-- Lock Indicator -->
    <div v-if="mapItem.is_locked" style="position: absolute; top: 6px; right: 6px;" title="[TRAVADO COM CADEADO]">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="#ff5e3a">
        <path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/>
      </svg>
    </div>

    <!-- 44x44 Orange Icon Container (Standard across all desktop app cards) -->
    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg v-if="mapItem.type === 'MAP_OPENSOURCE'" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
        <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"/>
        <line x1="8" y1="2" x2="8" y2="18"/><line x1="16" y1="6" x2="16" y2="22"/>
      </svg>
      <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
        <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/><path d="M15 9v12"/>
      </svg>
    </div>

    <!-- White Name with exact font typography -->
    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ mapItem.name }}
    </span>

    <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; justify-content: center; flex-wrap: wrap;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
        [{{ mapItem.type === 'MAP_OPENSOURCE' ? 'MAPA' : 'PLANTA' }}]
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8.5px;">
        // {{ mapItem.markers.length }} PINOS
      </span>
    </div>
  </div>
</template>
