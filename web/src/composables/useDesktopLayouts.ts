import { ref, computed } from 'vue'
import type { LayoutFolderNode, EnterpriseLayoutItem } from '../types/layoutTree'
import { initialLayoutFolders } from '../data/mockLayoutFolders'
import { initialRootLayouts } from '../data/mockLayoutRoot'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'

export function useDesktopLayouts() {
  const searchQuery = ref(''), currentFolderId = ref<string | null>(null)
  const rootLayouts = ref<EnterpriseLayoutItem[]>([...initialRootLayouts]), folders = ref<LayoutFolderNode[]>([...initialLayoutFolders])
  const selectedLayout = ref<EnterpriseLayoutItem | null>(null), notification = ref<string | null>(null)
  const isFolderModalOpen = ref(false), isWizardOpen = ref(false), draggedLayout = ref<EnterpriseLayoutItem | null>(null)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalLayouts = computed(() => rootLayouts.value.length + folders.value.reduce((acc, f) => acc + f.layouts.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || f.layouts.some(l => l.name.toLowerCase().includes(searchQuery.value.toLowerCase()))))
  const displayedLayouts = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.layouts : rootLayouts.value
    return src.filter(l => !q || l.name.toLowerCase().includes(q) || l.companyScope.toLowerCase().includes(q))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { if (notification.value === msg) notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }
  const handleDragStart = (l: EnterpriseLayoutItem) => { draggedLayout.value = l }
  const handleDropOnFolder = (targetFolderId: string) => { if (draggedLayout.value) { moveLayoutToFolder(draggedLayout.value.id, targetFolderId); draggedLayout.value = null } }

  const moveLayoutToFolder = (layoutId: string, targetFolderId?: string) => {
    let item: EnterpriseLayoutItem | null = null
    const rIdx = rootLayouts.value.findIndex(l => l.id === layoutId)
    if (rIdx >= 0) item = rootLayouts.value.splice(rIdx, 1)[0]
    else { for (const f of folders.value) { const idx = f.layouts.findIndex(l => l.id === layoutId); if (idx >= 0) { item = f.layouts.splice(idx, 1)[0]; break } } }
    if (!item) return
    if (!targetFolderId) { item.folderId = undefined; item.companyScope = 'GLOBAL'; rootLayouts.value.push(item) }
    else { const t = folders.value.find(f => f.id === targetFolderId); if (t) { item.folderId = targetFolderId; item.companyScope = t.name; t.layouts.push(item) } }
    showNotification(`Layout "${item.name}" movido com sucesso.`)
  }

  const toggleLockLayout = (layoutId: string) => {
    const apply = (l: EnterpriseLayoutItem) => { l.is_locked = !l.is_locked; showNotification(`[CADEADO] "${l.name}" está ${l.is_locked ? 'TRAVADO // SOMENTE LEITURA' : 'DESTRAVADO'}`) }
    const r = rootLayouts.value.find(l => l.id === layoutId); if (r) return apply(r)
    for (const f of folders.value) { const l = f.layouts.find(i => i.id === layoutId); if (l) return apply(l) }
  }

  const deleteLayoutById = (id: string) => {
    rootLayouts.value = rootLayouts.value.filter(l => l.id !== id); folders.value.forEach(f => { f.layouts = f.layouts.filter(l => l.id !== id) })
    if (selectedLayout.value?.id === id) selectedLayout.value = null
    showNotification('Layout removido.')
  }

  const requestDeleteFolder = (id: string) => { const f = folders.value.find(fold => fold.id === id); if (f) openDeletePrompt(f.id, f.name, f.layouts.length) }
  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const f = folders.value.find(fold => fold.id === folderToDelete.value!.id)
    if (f) { f.layouts.forEach(l => { l.folderId = undefined; l.companyScope = 'GLOBAL' }); rootLayouts.value.push(...f.layouts); folders.value = folders.value.filter(item => item.id !== f.id) }
    if (currentFolderId.value === folderToDelete.value.id) currentFolderId.value = null
    showNotification(`Pasta removida. ${folderToDelete.value.itemCount} layouts movidos para a raiz.`)
    closeDeletePrompt()
  }

  const handleSaveFolder = (name: string, clientType: 'company' | 'final_client' = 'company') => {
    folders.value.push({ id: `f_lay_${Date.now()}`, name: name.toUpperCase(), clientType, isExpanded: true, layouts: [] })
    showNotification(`Pasta "${name}" criada.`)
    isFolderModalOpen.value = false
  }

  const handleSaveLayout = (layout: EnterpriseLayoutItem) => {
    if (layout.folderId) {
      const f = folders.value.find(fold => fold.id === layout.folderId)
      if (f) { const idx = f.layouts.findIndex(l => l.id === layout.id); if (idx >= 0) f.layouts[idx] = layout; else f.layouts.push(layout) }
    } else {
      const idx = rootLayouts.value.findIndex(l => l.id === layout.id); if (idx >= 0) rootLayouts.value[idx] = layout; else rootLayouts.value.push(layout)
    }
    selectedLayout.value = null
    showNotification(`Layout "${layout.name}" salvo com sucesso.`)
    isWizardOpen.value = false
  }

  return {
    searchQuery, currentFolderId, rootLayouts, folders, selectedLayout, notification, isFolderModalOpen,
    isWizardOpen, isConfirmDeleteOpen, folderToDelete, contextMenu, currentFolder, totalLayouts,
    displayedFolders, displayedLayouts, showNotification, openContextMenu, handleDragStart, handleDropOnFolder,
    moveLayoutToFolder, toggleLockLayout, deleteLayoutById, requestDeleteFolder, confirmDeleteFolder,
    handleSaveFolder, handleSaveLayout
  }
}
