import { ref, computed } from 'vue'
import type { WorkspaceTab, CustomLayout, CameraStreamInfo, GridLayout } from '../types/mosaic'
import { useLayoutManager } from './useLayoutManager'

export function useWorkspaceTabs(layoutMgr: ReturnType<typeof useLayoutManager>, user: string = 'operador') {
  const openTabs = ref<WorkspaceTab[]>([])
  const activeTabId = ref<string>('')
  const activeTab = computed(() => openTabs.value.find(t => t.id === activeTabId.value) || null)

  const openLayoutTab = (l: CustomLayout) => {
    const existing = openTabs.value.find(t => t.id === l.id)
    if (existing) { activeTabId.value = existing.id; return }
    const newTab: WorkspaceTab = { id: l.id, name: l.name, grid: l.grid, is_system: l.is_system, is_saved: true, created_by: l.created_by, slots: l.slots ? [...l.slots] : [] }
    openTabs.value.push(newTab); activeTabId.value = newTab.id
  }

  const addCameraWithAutoTab = (cam: CameraStreamInfo) => {
    // 1. Se a aba ativa for 1x1 e já tiver essa câmera, mantém sem duplicar
    if (activeTab.value && activeTab.value.grid === '1x1' && activeTab.value.slots.length > 0 && (activeTab.value.slots[0].data as any)?.id === cam.id) {
      return
    }

    // 2. Se nenhuma aba estiver aberta, ou se a aba ativa for do sistema (Admin protegida)
    if (!activeTab.value || activeTab.value.is_system) {
      const existingTemp = openTabs.value.find(t => t.is_temporary && !t.is_saved)
      if (existingTemp) {
        activeTabId.value = existingTemp.id
        appendCameraToTab(existingTemp, cam)
        return
      }
      const count = openTabs.value.filter(t => t.is_temporary || !t.is_system).length + 1
      const newTab: WorkspaceTab = { id: `tab_temp_${Date.now()}`, name: `Layout ${count}`, grid: '1x1', is_system: false, is_temporary: true, is_saved: false, created_by: user, slots: [{ slot_index: 0, type: 'camera', data: cam }] }
      openTabs.value.push(newTab); activeTabId.value = newTab.id
      return
    }

    // 3. Adiciona na aba editável ativa
    appendCameraToTab(activeTab.value, cam)
  }

  const appendCameraToTab = (tab: WorkspaceTab, cam: CameraStreamInfo) => {
    const emptyIdx = tab.slots.findIndex(s => !s.data)
    if (emptyIdx >= 0) {
      tab.slots[emptyIdx] = { slot_index: emptyIdx, type: 'camera', data: cam }
    } else {
      tab.slots.push({ slot_index: tab.slots.length, type: 'camera', data: cam })
    }
    const activeCount = tab.slots.filter(s => s.data).length
    const isSpecialLayout = tab.grid.startsWith('1+') || tab.grid === '1x2'
    if (!isSpecialLayout) {
      if (activeCount <= 1) tab.grid = '1x1'
      else if (activeCount <= 4 && tab.grid === '1x1') tab.grid = '2x2'
      else if (activeCount > 4 && activeCount <= 9 && (tab.grid === '1x1' || tab.grid === '2x2')) tab.grid = '3x3'
      else if (activeCount > 9 && activeCount <= 16 && (tab.grid === '1x1' || tab.grid === '2x2' || tab.grid === '3x3')) tab.grid = '4x4'
      else if (activeCount > 16 && activeCount <= 25) tab.grid = '5x5'
    }
  }

  const setTabGrid = (tabId: string, grid: GridLayout) => {
    const target = openTabs.value.find(t => t.id === tabId)
    if (target) { target.grid = grid; target.slots.forEach((s, idx) => { s.slot_index = idx }) }
  }

  const saveTab = (tabId: string) => {
    const target = openTabs.value.find(t => t.id === tabId)
    if (!target) return
    target.is_temporary = false; target.is_saved = true
    layoutMgr.createOperatorLayout(target.grid, target.slots)
    layoutMgr.renameLayout(layoutMgr.activeLayoutId.value, target.name)
  }

  const closeTab = (tabId: string) => {
    const idx = openTabs.value.findIndex(t => t.id === tabId)
    if (idx < 0) return
    openTabs.value.splice(idx, 1)
    if (openTabs.value.length === 0) { activeTabId.value = ''; return }
    if (activeTabId.value === tabId) activeTabId.value = openTabs.value[Math.max(0, idx - 1)].id
  }

  const createBlankTab = (grid: GridLayout = '2x2') => {
    const count = openTabs.value.length + 1
    const newTab: WorkspaceTab = { id: `tab_temp_${Date.now()}`, name: `Layout ${count}`, grid, is_system: false, is_temporary: true, is_saved: false, created_by: user, slots: [] }
    openTabs.value.push(newTab); activeTabId.value = newTab.id
  }

  const purgeCameraFromTabs = (cameraId: string) => {
    openTabs.value.forEach(tab => {
      tab.slots.forEach(s => {
        if (s.type === 'camera' && (s.data as any)?.id === cameraId) {
          s.data = undefined
        }
      })
    })
  }

  return { openTabs, activeTabId, activeTab, openLayoutTab, addCameraWithAutoTab, setTabGrid, saveTab, closeTab, createBlankTab, purgeCameraFromTabs }
}
