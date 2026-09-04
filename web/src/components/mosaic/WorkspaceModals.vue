<script setup lang="ts">
import type { CustomLayout, GridLayout, CarouselConfig, CameraStreamInfo } from '../../types/mosaic'
import LayoutContextMenu from './LayoutContextMenu.vue'
import LayoutRenameModal from './LayoutRenameModal.vue'
import CarouselConfigModal from './CarouselConfigModal.vue'
import PlaybackDrawer from './PlaybackDrawer.vue'

defineProps<{
  contextMenu: { x: number; y: number; layout?: CustomLayout | null; isHeader?: boolean } | null
  layoutToRename: CustomLayout | null
  isCarouselOpen: boolean
  isPlaybackOpen: boolean
  selectedCamera: CameraStreamInfo | null
}>()

const emit = defineEmits<{
  (e: 'closeContextMenu'): void
  (e: 'renameLayout', layout: CustomLayout): void
  (e: 'saveRename', newName: string): void
  (e: 'cancelRename'): void
  (e: 'duplicateLayout', layout: CustomLayout): void
  (e: 'deleteLayout', id: string): void
  (e: 'selectGrid', grid: GridLayout): void
  (e: 'closeCarousel'): void
  (e: 'carouselCreated', c: CarouselConfig): void
  (e: 'closePlayback'): void
}>()
</script>

<template>
  <div>
    <!-- Context Menu -->
    <LayoutContextMenu
      v-if="contextMenu"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :layout="contextMenu.layout"
      :isHeaderMenu="contextMenu.isHeader"
      @close="emit('closeContextMenu')"
      @rename="emit('renameLayout', $event)"
      @duplicate="emit('duplicateLayout', $event)"
      @delete="emit('deleteLayout', $event.id)"
      @selectGrid="emit('selectGrid', $event)"
    />

    <!-- Rename Modal -->
    <LayoutRenameModal
      v-if="layoutToRename"
      :layout="layoutToRename"
      @close="emit('cancelRename')"
      @save="emit('saveRename', $event)"
    />

    <!-- Carousel Modal -->
    <CarouselConfigModal
      v-if="isCarouselOpen"
      @close="emit('closeCarousel')"
      @created="emit('carouselCreated', $event)"
    />

    <!-- Playback Drawer -->
    <PlaybackDrawer
      v-if="isPlaybackOpen && selectedCamera"
      :camera="selectedCamera"
      @close="emit('closePlayback')"
      @exportClip="() => {}"
    />
  </div>
</template>
