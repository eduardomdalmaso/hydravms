<script setup lang="ts">
import { ref, computed } from 'vue'
import type { GridLayout, WorkspaceSlot, CameraStreamInfo } from '../../types/mosaic'
import DynamicMosaicSlot from './DynamicMosaicSlot.vue'

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
  <div v-if="maxSlots === 0" class="vms-empty-mosaic-canvas">
    <span class="vms-text-mono vms-text-sm vms-text-dim" style="user-select: none;">[ + ]</span>
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
