<script setup lang="ts">
import { ref } from 'vue'

interface CameraPaletteItem {
  id: string; name: string; codec: string; res: string; fps: number; ip: string
}

const props = defineProps<{ usedCameraIds?: string[] }>()
const hoveredCam = ref<CameraPaletteItem | null>(null), hoverX = ref(0), hoverY = ref(0)

const availableCameras: CameraPaletteItem[] = [
  { id: 'cam_01', name: 'Portaria Principal', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.101' },
  { id: 'cam_02', name: 'Estacionamento Visitantes', codec: 'H.264', res: '1080P', fps: 25, ip: '192.168.1.102' },
  { id: 'cam_03', name: 'Corredor de Cargas & Docas', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.103' },
  { id: 'cam_04', name: 'Perímetro dos Fundos', codec: 'H.265', res: '4K', fps: 30, ip: '192.168.1.104' },
  { id: 'cam_05', name: 'Racks Servidores NOC', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.2.101' },
  { id: 'cam_root_01', name: 'Almoxarifado Geral', codec: 'H.265', res: '1080P', fps: 30, ip: '192.168.1.109' },
  { id: 'cam_root_02', name: 'Refeitório Central', codec: 'H.264', res: '1080P', fps: 25, ip: '192.168.1.110' }
]

const isUsed = (id: string) => props.usedCameraIds?.includes(id) ?? false
const onMouseEnter = (c: CameraPaletteItem, ev: MouseEvent) => { if (!isUsed(c.id)) { hoveredCam.value = c; hoverX.value = ev.clientX; hoverY.value = ev.clientY } }
const onMouseMove = (ev: MouseEvent) => { hoverX.value = ev.clientX; hoverY.value = ev.clientY }
const onDragStart = (c: CameraPaletteItem, ev: DragEvent) => {
  if (isUsed(c.id)) { ev.preventDefault(); return }
  hoveredCam.value = null; ev.dataTransfer?.setData('application/json', JSON.stringify(c))
  const ghost = document.getElementById(`ghost_${c.id}`); if (ghost && ev.dataTransfer?.setDragImage) ev.dataTransfer.setDragImage(ghost, 80, 45)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.5rem; flex: 1; min-height: 0;">
    <div class="vms-flex-between" style="align-items: center;">
      <span class="vms-text-mono vms-text-xs vms-font-bold" style="color: var(--vms-neu-accent-orange);">
        // CAMERAS DISPONIVEIS (ARRASTE PARA O GRID ACIMA)
      </span>
      <span class="vms-badge vms-badge-neutral" style="font-size: 8px;">{{ availableCameras.length - (usedCameraIds?.length || 0) }} DISPONIVEIS</span>
    </div>

    <!-- App Grid Mode with Grayed-Out Used Cameras -->
    <div class="vms-desktop-grid" style="grid-template-columns: repeat(auto-fill, minmax(110px, 1fr)); gap: 0.6rem; overflow-y: auto; padding: 4px;">
      <div
        v-for="c in availableCameras"
        :key="c.id"
        class="vms-desktop-app-card"
        :draggable="!isUsed(c.id)"
        :style="{
          padding: '0.6rem 0.4rem', position: 'relative',
          opacity: isUsed(c.id) ? '0.35' : '1',
          filter: isUsed(c.id) ? 'grayscale(1)' : 'none',
          cursor: isUsed(c.id) ? 'not-allowed' : 'grab',
          borderColor: isUsed(c.id) ? 'rgba(255, 255, 255, 0.05)' : undefined,
          background: isUsed(c.id) ? 'rgba(255, 255, 255, 0.015)' : undefined
        }"
        :title="isUsed(c.id) ? 'Esta câmera já está vinculada a um slot no grid' : 'Arraste para um slot no grid'"
        @dragstart="onDragStart(c, $event)"
        @mouseenter="onMouseEnter(c, $event)"
        @mousemove="onMouseMove"
        @mouseleave="hoveredCam = null"
      >
        <div style="width: 38px; height: 38px; border-radius: 8px; background: rgba(255, 94, 58, 0.14); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center;">
          <svg width="20" height="20" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        </div>
        <span class="vms-font-medium" style="color: #fff; font-size: 10px; line-height: 1.2; text-align: center; max-width: 100%; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ c.name }}</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: isUsed(c.id) ? 'var(--vms-text-dim)' : 'var(--vms-neu-accent-orange)' }" style="font-size: 8.5px;">
          {{ isUsed(c.id) ? '[NO GRID]' : `[${c.id.toUpperCase()}]` }}
        </span>

        <div :id="`ghost_${c.id}`" style="position: fixed; top: -9999px; left: -9999px; width: 160px; height: 90px; background: #07080c; border: 2px solid #ff5e3a; border-radius: 6px; display: flex; flex-direction: column; align-items: center; justify-content: center; z-index: 9999; padding: 4px;">
          <svg width="24" height="24" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
          <span style="color: #fff; font-size: 10px; font-weight: 700; margin-top: 4px; white-space: nowrap;">{{ c.name }}</span>
          <span style="color: #00ff9d; font-size: 8px; font-family: monospace;">[SNAPSHOT LIVE // {{ c.res }}]</span>
        </div>
      </div>
    </div>

    <!-- Fixed Floating Snapshot Popover on Hover -->
    <Teleport to="body">
      <div
        v-if="hoveredCam"
        style="position: fixed; z-index: 99999; pointer-events: none; transform: translate(-50%, -105%); width: 220px; background: #0c0f16; border: 1px solid var(--vms-neu-accent-orange); border-radius: 6px; padding: 6px; box-shadow: 0 10px 30px rgba(0,0,0,0.9); display: flex; flex-direction: column; gap: 4px;"
        :style="{ left: `${hoverX}px`, top: `${hoverY - 12}px` }"
      >
        <div style="width: 100%; aspect-ratio: 16 / 9; background: #07080c; border-radius: 4px; display: flex; align-items: center; justify-content: center; position: relative; border: 1px solid rgba(255, 94, 58, 0.35);">
          <svg width="28" height="28" viewBox="0 0 576 512" fill="#ff5e3a" style="opacity: 0.9;"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
          <span class="vms-badge vms-badge-success" style="position: absolute; top: 3px; left: 3px; font-size: 7px; padding: 0 3px;">[LIVE {{ hoveredCam.fps }}FPS]</span>
          <span class="vms-badge" style="position: absolute; bottom: 3px; right: 3px; font-size: 7px; padding: 0 3px; background: rgba(255, 94, 58, 0.15); color: #ff5e3a; border: 1px solid rgba(255, 94, 58, 0.35);">{{ hoveredCam.res }}</span>
        </div>
        <span class="vms-font-bold" style="color: #fff; font-size: 10px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ hoveredCam.name }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8px;">IP: {{ hoveredCam.ip }} // CODEC: {{ hoveredCam.codec }}</span>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 7.5px;">// ARRASTE PARA UM SLOT ACIMA</span>
      </div>
    </Teleport>
  </div>
</template>
