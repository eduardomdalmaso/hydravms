import { ref, computed } from 'vue'
import type { GridLayout, CameraStreamInfo, MapResource, CarouselConfig } from '../types/mosaic'
import { useLayoutManager } from './useLayoutManager'
import { useWorkspaceTabs } from './useWorkspaceTabs'

export function useLiveWorkspace(username: string = 'operador') {
  const layoutMgr = useLayoutManager(username)
  const workspaceTabs = useWorkspaceTabs(layoutMgr, username)
  const selectedCameraForPlayback = ref<CameraStreamInfo | null>(null)

  const currentLayout = computed<GridLayout>(() => workspaceTabs.activeTab.value?.grid || '2x2')
  const slots = computed(() => workspaceTabs.activeTab.value?.slots || [])

  const maxSlots = computed(() => {
    if (!workspaceTabs.activeTab.value) return 0
    switch (currentLayout.value) {
      case '1x1': return 1
      case '1x2': return 2
      case '2x2': return 4
      case '3x3': return 9
      case '4x4': return 16
      case '8x8': return 64
      case '10x10': return 100
      case '1+5': return 6
      case '1+7': return 8
      case '1+12': return 13
      default: return 4
    }
  })

  const addCameraToNextFreeSlot = (camera: CameraStreamInfo) => {
    workspaceTabs.addCameraWithAutoTab(camera)
  }

  const addMapToNextFreeSlot = (map: MapResource) => {
    if (!workspaceTabs.activeTab.value) workspaceTabs.createBlankTab('2x2')
    const s = workspaceTabs.activeTab.value!.slots
    const emptyIdx = s.findIndex(slot => !slot.data)
    if (emptyIdx >= 0) s[emptyIdx] = { slot_index: emptyIdx, type: 'map', data: map }
    else s.push({ slot_index: s.length, type: 'map', data: map })
  }

  const addCarouselToNextFreeSlot = (carousel: CarouselConfig) => {
    if (!workspaceTabs.activeTab.value) workspaceTabs.createBlankTab('2x2')
    const s = workspaceTabs.activeTab.value!.slots
    const emptyIdx = s.findIndex(slot => !slot.data)
    if (emptyIdx >= 0) s[emptyIdx] = { slot_index: emptyIdx, type: 'carousel', data: carousel }
    else s.push({ slot_index: s.length, type: 'carousel', data: carousel })
  }

  const clearSlot = (index: number) => {
    if (!workspaceTabs.activeTab.value) return
    const s = workspaceTabs.activeTab.value.slots
    if (s[index]) s[index] = { slot_index: index, type: 'camera', data: undefined }
  }

  const swapSlots = (fromIndex: number, toIndex: number) => {
    if (!workspaceTabs.activeTab.value) return
    const s = workspaceTabs.activeTab.value.slots
    if (!s[fromIndex] || !s[toIndex]) return
    const temp = s[fromIndex]
    s[fromIndex] = { ...s[toIndex], slot_index: fromIndex }
    s[toIndex] = { ...temp, slot_index: toIndex }
  }

  const selectLayout = (layoutId: string) => {
    const target = layoutMgr.layouts.value.find(l => l.id === layoutId)
    if (target) {
      workspaceTabs.openLayoutTab(target)
    }
  }

  return {
    layoutMgr,
    workspaceTabs,
    currentLayout,
    slots,
    maxSlots,
    selectedCameraForPlayback,
    addCameraToNextFreeSlot,
    addMapToNextFreeSlot,
    addCarouselToNextFreeSlot,
    clearSlot,
    swapSlots,
    selectLayout
  }
}
