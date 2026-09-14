<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useLiveWorkspace } from '../../composables/useLiveWorkspace'
import { useMultiMonitor } from '../../composables/useMultiMonitor'
import WorkspaceSidebar from '../../components/mosaic/WorkspaceSidebar.vue'
import WorkspaceTopTabs from '../../components/mosaic/WorkspaceTopTabs.vue'
import DynamicMosaicGrid from '../../components/mosaic/DynamicMosaicGrid.vue'
import WorkspaceModals from '../../components/mosaic/WorkspaceModals.vue'
import type { CameraStreamInfo, CustomLayout, GridLayout } from '../../types/mosaic'

const {
  layoutMgr, workspaceTabs, currentLayout, slots, maxSlots, selectedCameraForPlayback,
  addCameraToNextFreeSlot, addMapToNextFreeSlot, addCarouselToNextFreeSlot, clearSlot, swapSlots, selectLayout
} = useLiveWorkspace()

const { isPopout, activeMonitorNumber, openInPopout, dispatchToMonitor, initChannelListener } = useMultiMonitor()

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

const handlePopoutTab = (tabId: string) => {
  const tab = workspaceTabs.openTabs.value.find(t => t.id === tabId)
  const targetMonitor = tab?.target_monitor || 2
  openInPopout(tab?.id || tabId, targetMonitor)
}

const handlePopoutLayout = (layout: CustomLayout) => {
  const targetMonitor = layout.target_monitor || 2
  openInPopout(layout.id, targetMonitor)
  contextMenu.value = null
}

const handleDispatchLayout = (layout: CustomLayout) => {
  const targetMonitor = layout.target_monitor || 2
  dispatchToMonitor(layout.id, targetMonitor)
  contextMenu.value = null
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch(() => {})
  } else {
    document.exitFullscreen().catch(() => {})
  }
}

onMounted(() => {
  initChannelListener((layoutId) => {
    selectLayout(layoutId)
  })
  if (typeof window !== 'undefined') {
    const hash = window.location.hash || ''
    const match = hash.match(/layout=([^&]+)/)
    if (match && match[1]) {
      selectLayout(match[1])
    }
  }
})
</script>

<template>
  <div class="vms-mosaic-container" style="display: flex; flex-direction: row; height: 100%; width: 100%; overflow: hidden; position: relative;">
    <WorkspaceSidebar
      v-if="!isPopout"
      :layouts="layoutMgr.layouts.value" :activeLayoutId="workspaceTabs.activeTabId.value"
      @selectCamera="handleSidebarSelectCamera" @selectMap="addMapToNextFreeSlot"
      @selectCarousel="addCarouselToNextFreeSlot" @openCarouselModal="isCarouselOpen = true"
      @selectLayout="selectLayout" @layoutContextMenu="handleContextMenu"
    />
    <div style="flex: 1; display: flex; flex-direction: column; height: 100%; min-width: 0; min-height: 0; overflow: hidden; position: relative;">
      <!-- Popout HUD Badge -->
      <div v-if="isPopout" style="position: absolute; top: 6px; right: 10px; z-index: 50; display: flex; align-items: center; gap: 8px; background: rgba(15, 23, 42, 0.9); border: 1px solid rgba(255, 94, 58, 0.4); padding: 3px 8px; border-radius: 4px; backdrop-filter: blur(8px);">
        <span class="vms-status-dot online"></span>
        <span class="vms-text-mono vms-text-2xs" style="color: #ff5e3a; font-weight: 700;">MONITOR 0{{ activeMonitorNumber }} // VIDEO WALL</span>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 1px 4px; color: #cbd5e1;" title="Tela Cheia" @click="toggleFullscreen">[ TELA CHEIA ]</button>
      </div>

      <WorkspaceTopTabs
        :tabs="workspaceTabs.openTabs.value" :activeTabId="workspaceTabs.activeTabId.value"
        @selectTab="workspaceTabs.activeTabId.value = $event" @closeTab="workspaceTabs.closeTab($event)"
        @saveTab="workspaceTabs.saveTab($event)" @newTab="workspaceTabs.createBlankTab('2x2')"
        @popoutTab="handlePopoutTab"
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
        @closeContextMenu="contextMenu = null" @popoutLayout="handlePopoutLayout" @dispatchLayout="handleDispatchLayout"
        @renameLayout="layoutToRename = $event; contextMenu = null"
        @saveRename="(name) => { layoutMgr.renameLayout(layoutToRename!.id, name); layoutToRename = null }"
        @cancelRename="layoutToRename = null" @duplicateLayout="layoutMgr.duplicateLayout($event); contextMenu = null"
        @deleteLayout="layoutMgr.deleteLayout($event); contextMenu = null" @selectGrid="handleSelectGrid"
        @closeCarousel="isCarouselOpen = false" @carouselCreated="addCarouselToNextFreeSlot" @closePlayback="isPlaybackOpen = false"
      />
    </div>
  </div>
</template>
