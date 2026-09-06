import { ref, computed } from 'vue'
import type { RondaFolderNode, EnterpriseRondaItem } from '../types/rondaTree'
import { initialRondaFolders } from '../data/mockRondaFolders'
import { initialRootRondas } from '../data/mockRondaRoot'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'

export function useDesktopRondas() {
  const searchQuery = ref(''), currentFolderId = ref<string | null>(null)
  const rootRondas = ref<EnterpriseRondaItem[]>([...initialRootRondas]), folders = ref<RondaFolderNode[]>([...initialRondaFolders])
  const selectedRonda = ref<EnterpriseRondaItem | null>(null), notification = ref<string | null>(null)
  const isFolderModalOpen = ref(false), isWizardOpen = ref(false), draggedRonda = ref<EnterpriseRondaItem | null>(null)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalRondas = computed(() => rootRondas.value.length + folders.value.reduce((acc, f) => acc + f.rondas.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || f.rondas.some(r => r.name.toLowerCase().includes(searchQuery.value.toLowerCase()))))
  const displayedRondas = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.rondas : rootRondas.value
    return src.filter(r => !q || r.name.toLowerCase().includes(q) || r.companyScope.toLowerCase().includes(q))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { if (notification.value === msg) notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }
  const handleDragStart = (r: EnterpriseRondaItem) => { draggedRonda.value = r }
  const handleDropOnFolder = (targetFolderId: string) => { if (draggedRonda.value) { moveRondaToFolder(draggedRonda.value.id, targetFolderId); draggedRonda.value = null } }

  const moveRondaToFolder = (rondaId: string, targetFolderId?: string) => {
    let item: EnterpriseRondaItem | null = null
    const rIdx = rootRondas.value.findIndex(r => r.id === rondaId)
    if (rIdx >= 0) item = rootRondas.value.splice(rIdx, 1)[0]
    else { for (const f of folders.value) { const idx = f.rondas.findIndex(r => r.id === rondaId); if (idx >= 0) { item = f.rondas.splice(idx, 1)[0]; break } } }
    if (!item) return
    if (!targetFolderId) { item.folderId = undefined; item.companyScope = 'GLOBAL'; rootRondas.value.push(item) }
    else { const t = folders.value.find(f => f.id === targetFolderId); if (t) { item.folderId = targetFolderId; item.companyScope = t.name; t.rondas.push(item) } }
    showNotification(`Ronda "${item.name}" movida com sucesso.`)
  }

  const deleteRondaById = (id: string) => {
    rootRondas.value = rootRondas.value.filter(r => r.id !== id); folders.value.forEach(f => { f.rondas = f.rondas.filter(r => r.id !== id) })
    if (selectedRonda.value?.id === id) selectedRonda.value = null
    showNotification('Ronda removida.')
  }

  const requestDeleteFolder = (id: string) => { const f = folders.value.find(fold => fold.id === id); if (f) openDeletePrompt(f.id, f.name, f.rondas.length) }
  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const f = folders.value.find(fold => fold.id === folderToDelete.value!.id)
    if (f) { f.rondas.forEach(r => { r.folderId = undefined; r.companyScope = 'GLOBAL' }); rootRondas.value.push(...f.rondas); folders.value = folders.value.filter(item => item.id !== f.id) }
    if (currentFolderId.value === folderToDelete.value.id) currentFolderId.value = null
    showNotification(`Pasta removida. ${folderToDelete.value.itemCount} rondas movidas para a raiz.`)
    closeDeletePrompt()
  }

  const handleSaveFolder = (name: string, clientType: 'company' | 'final_client') => {
    folders.value.push({ id: `rf_${Date.now()}`, name, clientType, isExpanded: true, rondas: [] })
    isFolderModalOpen.value = false
    showNotification(`Pasta "${name}" criada com sucesso.`)
  }

  const handleSaveRonda = (ronda: EnterpriseRondaItem) => {
    let exists = false
    const rIdx = rootRondas.value.findIndex(r => r.id === ronda.id)
    if (rIdx >= 0) { rootRondas.value[rIdx] = { ...ronda }; exists = true }
    else { for (const f of folders.value) { const idx = f.rondas.findIndex(r => r.id === ronda.id); if (idx >= 0) { f.rondas[idx] = { ...ronda }; exists = true; break } } }
    if (!exists) { if (ronda.folderId) { const f = folders.value.find(fold => fold.id === ronda.folderId); if (f) f.rondas.push(ronda); else rootRondas.value.push(ronda) } else rootRondas.value.push(ronda) }
    selectedRonda.value = null
    isWizardOpen.value = false
    showNotification(`[RONDA] "${ronda.name}" salva com sucesso.`)
  }

  return {
    searchQuery, folders, currentFolderId, currentFolder, selectedRonda, notification, contextMenu,
    isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalRondas, displayedFolders,
    displayedRondas, openContextMenu, handleDragStart, handleDropOnFolder, moveRondaToFolder,
    requestDeleteFolder, confirmDeleteFolder, deleteRondaById, handleSaveFolder, handleSaveRonda, showNotification
  }
}
