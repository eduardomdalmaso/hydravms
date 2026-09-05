<script setup lang="ts">
import { computed } from 'vue'
import type { EnterpriseLayoutItem } from '../../../types/layoutTree'
import LayoutMiniGrid from './LayoutMiniGrid.vue'
import LayoutCameraPalette from './LayoutCameraPalette.vue'

const props = defineProps<{ layout: EnterpriseLayoutItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const usedCameraIds = computed(() => {
  return props.layout.slots.map(s => s.cameraId).filter((id): id is string => !!id)
})

const handleDropCamera = (slotIdx: number, cam: any) => {
  if (usedCameraIds.value.includes(cam.id)) {
    emit('saved', `[DUPLICADO] Câmera "${cam.name}" já está em uso na grade.`)
    return
  }
  const existing = props.layout.slots.find(s => s.slotIndex === slotIdx)
  if (existing) {
    existing.cameraId = cam.id
    existing.cameraName = cam.name
  } else {
    props.layout.slots.push({ slotIndex: slotIdx, cameraId: cam.id, cameraName: cam.name })
  }
  emit('saved', `[AUTO-SAVE] Câmera "${cam.name}" vinculada ao Slot ${slotIdx + 1}`)
}

const handleClearSlot = (slotIdx: number) => {
  const idx = props.layout.slots.findIndex(s => s.slotIndex === slotIdx)
  if (idx >= 0) {
    const removed = props.layout.slots.splice(idx, 1)[0]
    emit('saved', `[AUTO-SAVE] Slot ${slotIdx + 1} liberado (${removed.cameraName || ''})`)
  }
}
</script>

<template>
  <div class="vms-split-pane">
    <div class="vms-split-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">
          MAPEAR CAMERAS NOS SLOTS
        </span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          // AUTO-SAVE ATIVO (SOLTE A CAMERA NO SLOT PARA GRAVAR)
        </span>
      </div>
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border);">GRADE: {{ layout.grid }}</span>
    </div>

    <LayoutMiniGrid
      :grid="layout.grid"
      :slots="layout.slots"
      @drop-camera="handleDropCamera"
      @clear-slot="handleClearSlot"
    />

    <LayoutCameraPalette :used-camera-ids="usedCameraIds" />
  </div>
</template>
