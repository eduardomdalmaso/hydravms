<script setup lang="ts">
import { ref, computed } from 'vue'

interface CamItem { id: string; name: string; codec: string; res: string; fps: number; ip: string }

const props = defineProps<{ usedCameraIds: string[] }>()
const emit = defineEmits<{ (e: 'addCamera', cam: { id: string; name: string; resolution: string; fps: number }): void }>()

const searchQuery = ref('')
const hoveredCam = ref<CamItem | null>(null), hoverX = ref(0), hoverY = ref(0)

const allCameras: CamItem[] = [
  { id: 'cam_01', name: 'Portaria Principal', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.101' },
  { id: 'cam_02', name: 'Estacionamento Visitantes', codec: 'H.264', res: '1080P', fps: 25, ip: '192.168.1.102' },
  { id: 'cam_03', name: 'Corredor Cargas & Docas', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.103' },
  { id: 'cam_04', name: 'Perímetro dos Fundos', codec: 'H.265', res: '4K', fps: 30, ip: '192.168.1.104' },
  { id: 'cam_05', name: 'Racks Servidores NOC', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.2.101' },
  { id: 'cam_root_01', name: 'Almoxarifado Geral', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.109' },
  { id: 'cam_root_02', name: 'Refeitório Central', codec: 'H.264', res: '1080P', fps: 25, ip: '192.168.1.110' }
]

const availableCameras = computed(() =>
  allCameras.filter(c => !searchQuery.value || c.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
)

const onMouseEnter = (c: CamItem, ev: MouseEvent) => { hoveredCam.value = c; hoverX.value = ev.clientX; hoverY.value = ev.clientY }
const onMouseMove = (ev: MouseEvent) => { hoverX.value = ev.clientX; hoverY.value = ev.clientY }
const onDragStart = (c: CamItem, ev: DragEvent) => {
  hoveredCam.value = null
  ev.dataTransfer?.setData('application/json', JSON.stringify({ id: c.id, name: c.name, resolution: c.res, fps: c.fps }))
}
const handleAdd = (c: CamItem) => emit('addCamera', { id: c.id, name: c.name, resolution: c.res, fps: c.fps })
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.5rem; background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.65rem;">
    <div class="vms-flex-between" style="align-items: center;">
      <span class="vms-text-mono vms-text-2xs vms-font-bold" style="color: var(--vms-neu-accent-orange);">
        // ADICIONAR FLUXOS A RONDA (CLIQUE OU ARRASTE)
      </span>
      <input v-model="searchQuery" class="vms-auth-input" style="width: 140px; font-size: 10px; padding: 2px 6px; height: 22px;" placeholder="Filtrar câmeras..." />
    </div>

    <div class="vms-desktop-grid" style="grid-template-columns: repeat(auto-fill, minmax(110px, 1fr)); gap: 0.5rem; max-height: 120px; overflow-y: auto; padding: 2px;">
      <div
        v-for="c in availableCameras"
        :key="c.id"
        class="vms-desktop-app-card"
        draggable="true"
        style="padding: 0.5rem 0.35rem; position: relative; cursor: grab;"
        @dragstart="onDragStart(c, $event)"
        @click="handleAdd(c)"
        @mouseenter="onMouseEnter(c, $event)"
        @mousemove="onMouseMove"
        @mouseleave="hoveredCam = null"
      >
        <span class="vms-badge vms-badge-orange" style="position: absolute; top: 4px; right: 4px; font-size: 8px; padding: 1px 4px;">+</span>

        <div style="width: 34px; height: 34px; border-radius: 8px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center;">
          <svg width="18" height="18" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        </div>

        <span class="vms-font-medium" style="color: #fff; font-size: 9.5px; line-height: 1.2; text-align: center; max-width: 100%; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ c.name }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8px;">{{ c.res }} @ {{ c.fps }}FPS</span>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="hoveredCam"
        style="position: fixed; z-index: 99999; pointer-events: none; transform: translate(-50%, -105%); width: 220px; background: #0c0f16; border: 1px solid var(--vms-neu-accent-orange); border-radius: 6px; padding: 6px; box-shadow: 0 10px 30px rgba(0,0,0,0.9); display: flex; flex-direction: column; gap: 4px;"
        :style="{ left: `${hoverX}px`, top: `${hoverY - 12}px` }"
      >
        <div style="width: 100%; aspect-ratio: 16 / 9; background: #07080c; border-radius: 4px; display: flex; align-items: center; justify-content: center; position: relative; border: 1px solid rgba(255, 94, 58, 0.35);">
          <svg width="28" height="28" viewBox="0 0 576 512" fill="#ff5e3a" style="opacity: 0.9;"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
          <span class="vms-badge vms-badge-online" style="position: absolute; top: 3px; left: 3px; font-size: 7px; padding: 0 3px;">[LIVE {{ hoveredCam.fps }}FPS]</span>
          <span class="vms-badge" style="position: absolute; bottom: 3px; right: 3px; font-size: 7px; padding: 0 3px; background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border);">{{ hoveredCam.res }}</span>
        </div>
        <span class="vms-font-bold" style="color: #fff; font-size: 10px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ hoveredCam.name }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8px;">IP: {{ hoveredCam.ip }} // CODEC: {{ hoveredCam.codec }}</span>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 7.5px;">// CLIQUE PARA ADICIONAR A RONDA</span>
      </div>
    </Teleport>
  </div>
</template>
