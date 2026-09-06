<script setup lang="ts">
import { ref, computed } from 'vue'
import type { DeployableDevice } from '../../../types/mapTree'
import MapDeviceCard from './MapDeviceCard.vue'

const props = defineProps<{
  activeTab: 'CAMERAS' | 'ALARMES'
  devices: DeployableDevice[]
  isLocked?: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'dragStart', dev: DeployableDevice): void
  (e: 'applyCoords', dev: DeployableDevice): void
}>()

const filterQuery = ref('')
const filteredDevices = computed(() => {
  const q = filterQuery.value.toLowerCase()
  return props.devices.filter(d => !q || d.name.toLowerCase().includes(q) || (d.details && d.details.toLowerCase().includes(q)))
})
</script>

<template>
  <div class="vms-map-device-drawer">
    <div class="vms-drawer-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">
          {{ activeTab === 'CAMERAS' ? 'CAMERAS DISPONIVEIS' : 'SENSORES & ALARMES' }}
        </span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ARRASTE E SOLTE NO MAPA</span>
      </div>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 6px;" @click="emit('close')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <input v-model="filterQuery" class="vms-auth-input" style="font-size: 11px; padding: 4px 8px; margin: 6px 8px;" placeholder="Filtrar dispositivos..." />

    <div class="vms-drawer-body">
      <div v-if="filteredDevices.length === 0" class="vms-text-dim vms-text-xs" style="text-align: center; padding: 1.5rem 0;">
        Nenhum dispositivo encontrado
      </div>
      <MapDeviceCard
        v-for="d in filteredDevices"
        :key="d.id"
        :device="d"
        :is-locked="isLocked"
        @drag-start="emit('dragStart', d)"
        @apply-coords="emit('applyCoords', d)"
      />
    </div>
  </div>
</template>

<style scoped>
.vms-map-device-drawer {
  position: absolute; top: 58px; right: 12px; z-index: 1000; width: 310px;
  max-height: calc(100% - 70px); background: rgba(14, 17, 23, 0.94);
  backdrop-filter: blur(12px); border: 1px solid var(--vms-border); border-radius: 8px;
  display: flex; flex-direction: column; box-shadow: 0 8px 32px rgba(0, 0, 0, 0.7);
  animation: drawerSlide 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes drawerSlide {
  from { opacity: 0; transform: translateY(-8px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.vms-drawer-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 8px 12px; border-bottom: 1px solid var(--vms-border); background: #07080c;
}
.vms-drawer-body {
  padding: 8px; overflow-y: auto; display: flex; flex-direction: column; gap: 6px; flex: 1;
}
</style>
