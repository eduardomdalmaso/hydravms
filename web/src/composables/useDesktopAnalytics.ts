import { ref, computed } from 'vue'
import type { AnalyticFolderNode, AnalyticInstance } from '../types/marketplace'
import { mockMarketplaceFolders } from '../data/mockMarketplaceFolders'
import { useMarketplace } from './useMarketplace'

const allFolders = ref<AnalyticFolderNode[]>([...mockMarketplaceFolders])
const rootInstances = ref<AnalyticInstance[]>([])

export function useDesktopAnalytics(pluginId: string) {
  const { showToast } = useMarketplace()
  const searchQuery = ref('')
  const currentFolderId = ref<string | null>(null)
  const selectedInstance = ref<AnalyticInstance | null>(null)
  const isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const draggedInstance = ref<AnalyticInstance | null>(null)

  const pluginFolders = computed(() => allFolders.value.filter(f => f.plugin_id === pluginId))
  const currentFolder = computed(() => pluginFolders.value.find(f => f.id === currentFolderId.value) || null)
  const currentRootInstances = computed(() => rootInstances.value.filter(i => i.plugin_id === pluginId))

  const totalInstances = computed(() => {
    return currentRootInstances.value.length + pluginFolders.value.reduce((acc, f) => acc + f.instances.length, 0)
  })

  const displayedFolders = computed(() => {
    return pluginFolders.value
      .filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
      .sort((a, b) => a.name.localeCompare(b.name))
  })

  const displayedInstances = computed(() => {
    const q = searchQuery.value.toLowerCase()
    const list = currentFolder.value ? currentFolder.value.instances : currentRootInstances.value
    return list.filter(i => !q || i.name.toLowerCase().includes(q) || i.camera_name.toLowerCase().includes(q))
  })

  const handleDragStart = (inst: AnalyticInstance) => { draggedInstance.value = inst }

  const handleDropOnFolder = (folderId: string) => {
    if (!draggedInstance.value) return
    const id = draggedInstance.value.id
    rootInstances.value = rootInstances.value.filter(i => i.id !== id)
    allFolders.value.forEach(f => { f.instances = f.instances.filter(i => i.id !== id) })
    const dest = allFolders.value.find(f => f.id === folderId)
    if (dest) {
      dest.instances.push(draggedInstance.value)
      showToast(`[ANALÍTICO] "${draggedInstance.value.name}" movido para "${dest.name}"`)
    }
    draggedInstance.value = null
  }

  const handleCreateFolder = (name: string) => {
    allFolders.value.push({ id: `fld-${Date.now()}`, name, plugin_id: pluginId, instances: [] })
    isFolderModalOpen.value = false
    showToast(`[PASTA] "${name}" criada com sucesso`)
  }

  const handleCreateInstance = (inst: AnalyticInstance, targetFolderId?: string | null) => {
    const target = targetFolderId || currentFolderId.value
    if (target) {
      const f = allFolders.value.find(x => x.id === target)
      if (f) f.instances.push(inst); else rootInstances.value.push(inst)
    } else {
      rootInstances.value.push(inst)
    }
    isWizardOpen.value = false
    selectedInstance.value = inst
    showToast(`[NOVO ANALÍTICO] "${inst.name}" configurado`)
  }

  const handleSaveInstance = (inst: AnalyticInstance) => {
    selectedInstance.value = null
    showToast(`[SALVO] Analítico "${inst.name}" atualizado`)
  }

  const handleDeleteInstance = (id: string) => {
    rootInstances.value = rootInstances.value.filter(i => i.id !== id)
    allFolders.value.forEach(f => { f.instances = f.instances.filter(i => i.id !== id) })
    selectedInstance.value = null
    showToast(`[EXCLUÍDO] Analítico removido com sucesso`)
  }

  return {
    searchQuery, currentFolderId, currentFolder, selectedInstance, isFolderModalOpen, isWizardOpen,
    totalInstances, displayedFolders, displayedInstances, pluginFolders, handleDragStart,
    handleDropOnFolder, handleCreateFolder, handleCreateInstance, handleSaveInstance, handleDeleteInstance
  }
}
