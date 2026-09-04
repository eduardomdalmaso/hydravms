<script setup lang="ts">
import type { GridLayout, CameraStreamInfo } from '../../types/mosaic'
import MosaicSlot from './MosaicSlot.vue'

const props = defineProps<{
  layout: GridLayout
  cameras: CameraStreamInfo[]
  selectedCameraId?: string
  maxSlots: number
}>()

const emit = defineEmits<{ (e: 'selectCamera', cam: CameraStreamInfo): void }>()
</script>

<template>
  <div
    class="vms-mosaic-grid"
    :class="{
      'vms-grid-1x1': layout === '1x1',
      'vms-grid-2x2': layout === '2x2',
      'vms-grid-3x3': layout === '3x3',
      'vms-grid-4x4': layout === '4x4',
      'vms-grid-1-5': layout === '1+5'
    }"
  >
    <MosaicSlot
      v-for="index in maxSlots"
      :key="index"
      :camera="cameras[index - 1]"
      :isActive="cameras[index - 1]?.id === selectedCameraId"
      :isHero="layout === '1+5' && index === 1"
      @select="emit('selectCamera', $event)"
    />
  </div>
</template>
