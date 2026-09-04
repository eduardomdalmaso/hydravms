<script setup lang="ts">
import { ref } from 'vue'
import { useLiveWorkspace } from '../../composables/useLiveWorkspace'
import WorkspaceSidebar from '../../components/mosaic/WorkspaceSidebar.vue'
import WorkspaceTopTabs from '../../components/mosaic/WorkspaceTopTabs.vue'
import DynamicMosaicGrid from '../../components/mosaic/DynamicMosaicGrid.vue'
import WorkspaceModals from '../../components/mosaic/WorkspaceModals.vue'
import type { CameraStreamInfo, CustomLayout, GridLayout } from '../../types/mosaic'

const {
  layoutMgr, workspaceTabs, currentLayout, slots, maxSlots, selectedCameraForPlayback,
  addCameraToNextFreeSlot, addMapToNextFreeSlot, addCarouselToNextFreeSlot, clearSlot, swapSlots, selectLayout
} = useLiveWorkspace()

const isCarouselOpen = ref(false)
const isPlaybackOpen = ref(false)
const contextMenu = ref<{ x: number; y: number; layout?: CustomLayout | null; isHeader?: boolean } | null>(null)
const layoutToRename = ref<CustomLayout | null>(null)

const handleOpenPlayback = (cam: CameraStreamInfo) => {
  selectedCameraForPlayback.value = cam; isPlaybackOpen.value = true
}
const handleSidebarSelectCamera = (cam: CameraStreamInfo) => {
  addCameraToNextFreeSlot(cam); handleOpenPlayback(cam)
}
const handleContextMenu = (p: { event: MouseEvent; layout?: CustomLayout; isHeader?: boolean }) => {
  contextMenu.value = { x: p.event.clientX, y: p.event.clientY, layout: p.layout || null, isHeader: p.isHeader }
}
const handleSelectGrid = (grid: GridLayout) => {
  if (contextMenu.value?.layout && !contextMenu.value.layout.is_system) {
    contextMenu.value.layout.grid = grid; layoutMgr.saveOperatorLayouts()
  } else if (workspaceTabs.activeTab.value) {
    workspaceTabs.setTabGrid(workspaceTabs.activeTabId.value, grid)
  } else { workspaceTabs.createBlankTab(grid) }
  contextMenu.value = null
}
</script>

<template>
  <div class="vms-mosaic-container" style="display: flex; flex-direction: row; height: 100%; width: 100%; overflow: hidden; position: relative;">
    <WorkspaceSidebar
      :layouts="layoutMgr.layouts.value" :activeLayoutId="workspaceTabs.activeTabId.value"
      @selectCamera="handleSidebarSelectCamera" @selectMap="addMapToNextFreeSlot"
      @selectCarousel="addCarouselToNextFreeSlot" @openCarouselModal="isCarouselOpen = true"
      @selectLayout="selectLayout" @layoutContextMenu="handleContextMenu"
    />
    <div style="flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; min-height: 0; overflow: hidden; position: relative;">
      <WorkspaceTopTabs
        :tabs="workspaceTabs.openTabs.value" :activeTabId="workspaceTabs.activeTabId.value"
        @selectTab="workspaceTabs.activeTabId.value = $event" @closeTab="workspaceTabs.closeTab($event)"
        @saveTab="workspaceTabs.saveTab($event)" @newTab="workspaceTabs.createBlankTab('2x2')"
      />
      <div style="flex: 1; min-height: 0; width: 100%; display: flex; overflow: hidden;">
        <DynamicMosaicGrid
          :layout="currentLayout" :slots="slots" :maxSlots="maxSlots" :selectedCameraId="selectedCameraForPlayback?.id"
          @selectCamera="handleOpenPlayback" @clearSlot="clearSlot" @swapSlots="swapSlots"
        />
      </div>
      <WorkspaceModals
        :contextMenu="contextMenu" :layoutToRename="layoutToRename" :isCarouselOpen="isCarouselOpen"
        :isPlaybackOpen="isPlaybackOpen" :selectedCamera="selectedCameraForPlayback"
        @closeContextMenu="contextMenu = null" @renameLayout="layoutToRename = $event; contextMenu = null"
        @saveRename="(name) => { layoutMgr.renameLayout(layoutToRename!.id, name); layoutToRename = null }"
        @cancelRename="layoutToRename = null" @duplicateLayout="layoutMgr.duplicateLayout($event); contextMenu = null"
        @deleteLayout="layoutMgr.deleteLayout($event); contextMenu = null" @selectGrid="handleSelectGrid"
        @closeCarousel="isCarouselOpen = false" @carouselCreated="addCarouselToNextFreeSlot" @closePlayback="isPlaybackOpen = false"
      />
    </div>
  </div>
</template>
