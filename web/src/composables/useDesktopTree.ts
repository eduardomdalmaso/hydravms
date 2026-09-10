import { ref, computed, onMounted } from 'vue'
import type { FolderNode, StreamItem } from '../types/streamTree'
import { initialFolders, initialRootStreams } from '../data/mockStreamFolders'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'
import { fetchFolders, fetchCameras, createRemoteFolder, createRemoteCamera, deleteRemoteCamera, deleteRemoteFolder } from '../services/api'
import { buildStreamTree, createNewStreamItem } from '../utils/streamTreeUtils'

export function useDesktopTree() {
  const searchQuery = ref(''), folders = ref<FolderNode[]>(initialFolders), rootStreams = ref<StreamItem[]>(initialRootStreams)
  const currentFolderId = ref<string | null>(null), selectedStream = ref<StreamItem | null>(null), draggedStream = ref<StreamItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const loadData = async () => {
    const [dbF, dbC] = await Promise.all([fetchFolders('cameras'), fetchCameras()])
    const tree = buildStreamTree(dbF || [], dbC || [])
    folders.value = tree.folders; rootStreams.value = tree.rootStreams
  }
  onMounted(loadData)

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalStreams = computed(() => rootStreams.value.length + folders.value.reduce((acc, f) => acc + f.streams.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase())).sort((a, b) => a.name.localeCompare(b.name)))
  const displayedStreams = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.streams : rootStreams.value
    return src.filter(s => !q || s.name.toLowerCase().includes(q) || s.ip.includes(q)).sort((a, b) => a.name.localeCompare(b.name))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }
  const handleDragStart = (s: StreamItem) => { draggedStream.value = s }
  const handleDropOnFolder = (targetFolderId: string) => { if (draggedStream.value) { moveStreamToFolder(draggedStream.value.id, targetFolderId); draggedStream.value = null } }

  const moveStreamToFolder = (streamId: string, targetFolderId: string | null) => {
    let stream: StreamItem | undefined = rootStreams.value.find(s => s.id === streamId)
    if (!stream) { for (const f of folders.value) { stream = f.streams.find(s => s.id === streamId); if (stream) break } }
    if (!stream) return
    rootStreams.value = rootStreams.value.filter(s => s.id !== streamId)
    folders.value.forEach(f => { f.streams = f.streams.filter(s => s.id !== streamId) })
    if (targetFolderId) {
      const target = folders.value.find(f => f.id === targetFolderId)
      if (target) { target.streams.push(stream); showNotification(`Fluxo "${stream.name}" movido para "${target.name}".`) }
    } else { rootStreams.value.push(stream); showNotification(`Fluxo "${stream.name}" movido para a Raiz.`) }
  }

  const requestDeleteFolder = (id: string) => {
    const f = folders.value.find(fold => fold.id === id)
    if (f && f.streams.length > 0) openDeletePrompt(f.id, f.name, f.streams.length); else deleteFolderById(id)
  }

  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const f = folders.value.find(fold => fold.id === folderToDelete.value!.id)
    if (f) {
      rootStreams.value.push(...f.streams); folders.value = folders.value.filter(item => item.id !== f.id)
      if (currentFolderId.value === f.id) currentFolderId.value = null
      showNotification(`Pasta removida. ${f.streams.length} fluxos movidos para a raiz.`)
    }
    closeDeletePrompt()
  }

  const deleteFolderById = async (id: string) => {
    if (id.includes('-')) await deleteRemoteFolder(id)
    folders.value = folders.value.filter(f => f.id !== id); if (currentFolderId.value === id) currentFolderId.value = null; showNotification('Pasta removida.')
  }
  const deleteStreamById = async (id: string) => {
    await deleteRemoteCamera(id)
    rootStreams.value = rootStreams.value.filter(s => s.id !== id); folders.value.forEach(f => { f.streams = f.streams.filter(s => s.id !== id) })
    if (selectedStream.value?.id === id) selectedStream.value = null; showNotification('Fluxo removido.')
  }

  const handleSaveFolder = async (name: string) => {
    const res = await createRemoteFolder('cameras', name)
    folders.value.push({ id: res?.id || `f_${Date.now()}`, name, isExpanded: true, streams: [] }); showNotification(`Pasta ${name} criada.`)
  }

  const handleSaveStream = async (stream: Partial<StreamItem>, folderId?: string) => {
    const newStream = createNewStreamItem(stream, totalStreams.value + 1)
    const saved = await createRemoteCamera({
      id: newStream.id, name: newStream.name, protocol: newStream.protocol?.toLowerCase() || 'rtsp',
      rtsp_url: newStream.url, onvif_ip: newStream.ip, onvif_port: newStream.port, location: newStream.locationName || '',
      status: 'online', resolution: newStream.resolution || '1920x1080', fps: newStream.fps || 30.0,
      bitrate_kbps: 4096, codec: newStream.codec || 'H.265', has_ptz: !!newStream.has_ptz,
      folder_id: (folderId && folderId.includes('-')) ? folderId : undefined
    })
    if (saved) newStream.id = saved.id
    const target = folders.value.find(f => f.id === folderId)
    if (target) target.streams.push(newStream); else rootStreams.value.push(newStream)
    selectedStream.value = newStream; showNotification('Fluxo salvo com sucesso.')
  }

  return {
    searchQuery, folders, rootStreams, currentFolderId, currentFolder, selectedStream, notification, contextMenu,
    isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalStreams, displayedFolders, displayedStreams,
    openContextMenu, handleDragStart, handleDropOnFolder, moveStreamToFolder, requestDeleteFolder, confirmDeleteFolder,
    deleteFolderById, deleteStreamById, handleSaveFolder, handleSaveStream, showNotification
  }
}
