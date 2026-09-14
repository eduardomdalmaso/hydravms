<script setup lang="ts">
import { ref, computed } from 'vue'
import type { GridLayout, WorkspaceSlot, CameraStreamInfo } from '../../types/mosaic'
import DynamicMosaicSlot from './DynamicMosaicSlot.vue'
import { useBranding } from '../../composables/useBranding'

const props = defineProps<{
  layout: GridLayout
  slots: WorkspaceSlot[]
  selectedCameraId?: string
  maxSlots: number
}>()

const emit = defineEmits<{
  (e: 'selectCamera', cam: CameraStreamInfo): void
  (e: 'clearSlot', index: number): void
  (e: 'swapSlots', fromIndex: number, toIndex: number): void
}>()

const { branding } = useBranding()
const draggedSlotIndex = ref<number | null>(null)

const effectiveGridClass = computed(() => {
  switch (props.layout) {
    case '1x1': return 'vms-grid-1x1'
    case '1x2': return 'vms-grid-1x2'
    case '2x2': return 'vms-grid-2x2'
    case '3x3': return 'vms-grid-3x3'
    case '4x4': return 'vms-grid-4x4'
    case '8x8': return 'vms-grid-8x8'
    case '10x10': return 'vms-grid-10x10'
    case '1+5': return 'vms-grid-1-5'
    case '1+7': return 'vms-grid-1-7'
    case '1+12': return 'vms-grid-1-12'
    default: return 'vms-grid-2x2'
  }
})

const onDrop = (targetIndex: number) => {
  if (draggedSlotIndex.value !== null && draggedSlotIndex.value !== targetIndex) {
    emit('swapSlots', draggedSlotIndex.value, targetIndex)
  }
  draggedSlotIndex.value = null
}
</script>

<template>
  <div v-if="maxSlots === 0" class="vms-empty-mosaic-canvas" style="background: #181b22; display: flex; align-items: center; justify-content: center; flex: 1; height: 100%; border: 1px dashed rgba(255, 255, 255, 0.08);">
    <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.5rem; user-select: none; opacity: 0.6;">
      <img :src="branding.headerLogo || '/hydra.svg'" alt="Logo" style="width: 38px; height: 38px; object-fit: contain; filter: grayscale(1) brightness(0.65);" />
      <span class="vms-text-mono vms-text-sm vms-font-bold" style="color: #94a3b8; letter-spacing: 1px; text-transform: uppercase;">
        {{ branding.companyName || branding.systemName || 'HYDRA VMS' }}
      </span>
      <span class="vms-text-mono vms-text-xs" style="color: #64748b;">
        [ SELECIONE UM LAYOUT OU ARRASTE UMA CÂMERA ]
      </span>
    </div>
  </div>
  <div
    v-else
    class="vms-mosaic-grid"
    :class="effectiveGridClass"
  >
    <DynamicMosaicSlot
      v-for="index in maxSlots"
      :key="index"
      :slot="slots[index - 1] || { slot_index: index - 1, type: 'camera' }"
      :isActive="slots[index - 1]?.type === 'camera' && (slots[index - 1]?.data as any)?.id === selectedCameraId"
      :isHero="(effectiveGridClass === 'vms-grid-1-5' || effectiveGridClass === 'vms-grid-1-7' || effectiveGridClass === 'vms-grid-1-12') && index === 1"
      @selectCamera="emit('selectCamera', $event)"
      @clear="emit('clearSlot', $event)"
      @dragstart="draggedSlotIndex = index - 1"
      @dragover.prevent
      @drop="onDrop(index - 1)"
    />
  </div>
</template>
