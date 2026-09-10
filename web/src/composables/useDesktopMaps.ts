import { ref, computed, onMounted } from 'vue'
import type { MapFolderNode, EnterpriseMapItem } from '../types/mapTree'
import { initialMapFolders, initialRootMaps } from '../data/mockMapFolders'
import { fetchFolders } from '../services/api'

export function useDesktopMaps() {
  const searchQuery = ref(''), folders = ref<MapFolderNode[]>(initialMapFolders), rootMaps = ref<EnterpriseMapItem[]>(initialRootMaps)
  const currentFolderId = ref<string | null>(null), selectedMap = ref<EnterpriseMapItem | null>(null), draggedMap = ref<EnterpriseMapItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isNewMapModalOpen = ref(false)

  onMounted(async () => {
    const dbF = await fetchFolders('maps')
    if (dbF.length > 0) folders.value = dbF.map(f => ({ id: f.id, name: f.name, clientType: 'company', isExpanded: true, maps: [] }))
  })

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalMaps = computed(() => rootMaps.value.length + folders.value.reduce((a, f) => a + f.maps.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase())).sort((a, b) => a.name.localeCompare(b.name)))
  const displayedMaps = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.maps : rootMaps.value
    return src.filter(m => !q || m.name.toLowerCase().includes(q) || m.companyScope.toLowerCase().includes(q)).sort((a, b) => a.name.localeCompare(b.name))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { if (notification.value === msg) notification.value = null }, 3500) }

  const handleDragStart = (m: EnterpriseMapItem) => { draggedMap.value = m }
  const handleDropOnFolder = (folderId: string) => {
    if (!draggedMap.value) return
    const mapId = draggedMap.value.id
    let item: EnterpriseMapItem | undefined = rootMaps.value.find(m => m.id === mapId)
    if (!item) { for (const f of folders.value) { item = f.maps.find(m => m.id === mapId); if (item) break } }
    if (!item) return
    rootMaps.value = rootMaps.value.filter(m => m.id !== mapId)
    folders.value.forEach(f => f.maps = f.maps.filter(m => m.id !== mapId))
    const dest = folders.value.find(f => f.id === folderId)
    if (dest) { dest.maps.push(item); showNotification(`[MAPA] "${item.name}" movido para "${dest.name}"`) }
    draggedMap.value = null
  }

  const handleCreateFolder = (name: string, clientType: 'company' | 'final_client') => {
    folders.value.push({ id: `mf_${Date.now()}`, name, clientType, isExpanded: true, maps: [] })
    isFolderModalOpen.value = false; showNotification(`[PASTA] "${name}" criada com sucesso`)
  }

  const handleCreateMap = (newMap: EnterpriseMapItem) => {
    if (currentFolder.value) currentFolder.value.maps.push(newMap); else rootMaps.value.push(newMap)
    isNewMapModalOpen.value = false; selectedMap.value = newMap
    showNotification(`[NOVO MAPA] "${newMap.name}" criado com sucesso`)
  }

  const handleSaveMap = (m: EnterpriseMapItem) => { showNotification(`[SALVO] Mapa "${m.name}" salvo com sucesso`) }

  const handleToggleLock = () => {
    if (!selectedMap.value) return
    selectedMap.value.is_locked = !selectedMap.value.is_locked
    showNotification(selectedMap.value.is_locked ? `[TRAVADO COM CADEADO] Mapa "${selectedMap.value.name}" bloqueado` : `[DESTRAVADO] Mapa "${selectedMap.value.name}" liberado`)
  }

  return {
    searchQuery, folders, rootMaps, currentFolderId, currentFolder, selectedMap, draggedMap,
    notification, isFolderModalOpen, isNewMapModalOpen, totalMaps, displayedFolders, displayedMaps,
    showNotification, handleDragStart, handleDropOnFolder, handleCreateFolder, handleCreateMap,
    handleSaveMap, handleToggleLock
  }
}
