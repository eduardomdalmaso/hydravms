import { ref, computed } from 'vue'
import type { FolderNode, StreamItem } from '../types/streamTree'
import { initialFolders, initialRootStreams } from '../data/mockStreamFolders'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'

export function useDesktopTree() {
  const searchQuery = ref(''), folders = ref<FolderNode[]>(initialFolders), rootStreams = ref<StreamItem[]>(initialRootStreams)
  const currentFolderId = ref<string | null>(null), selectedStream = ref<StreamItem | null>(null), draggedStream = ref<StreamItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

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
    } else {
      rootStreams.value.push(stream)
      showNotification(`Fluxo "${stream.name}" movido para a Raiz.`)
    }
  }

  const requestDeleteFolder = (id: string) => {
    const folder = folders.value.find(f => f.id === id)
    if (!folder) return
    if (folder.streams.length > 0) openDeletePrompt(folder.id, folder.name, folder.streams.length)
    else deleteFolderById(id)
  }

  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const folder = folders.value.find(f => f.id === folderToDelete.value!.id)
    if (folder) {
      rootStreams.value.push(...folder.streams)
      folders.value = folders.value.filter(f => f.id !== folder.id)
      if (currentFolderId.value === folder.id) currentFolderId.value = null
      showNotification(`Pasta removida. ${folder.streams.length} fluxos movidos para a raiz.`)
    }
    closeDeletePrompt()
  }

  const deleteFolderById = (id: string) => { folders.value = folders.value.filter(f => f.id !== id); if (currentFolderId.value === id) currentFolderId.value = null; showNotification('Pasta removida.') }
  const deleteStreamById = (id: string) => {
    rootStreams.value = rootStreams.value.filter(s => s.id !== id); folders.value.forEach(f => { f.streams = f.streams.filter(s => s.id !== id) })
    if (selectedStream.value?.id === id) selectedStream.value = null
    showNotification('Fluxo removido.')
  }

  const handleSaveFolder = (name: string) => { folders.value.push({ id: `f_0${folders.value.length + 1}`, name, isExpanded: true, streams: [] }); showNotification(`Pasta ${name} criada.`) }
  const handleSaveStream = (stream: Partial<StreamItem>, folderId?: string) => {
    const newStream: StreamItem = {
      id: `cam_0${totalStreams.value + 1}`, name: stream.name || 'Nova Camera', protocol: stream.protocol || 'RTSP', url: stream.url || 'rtsp://', ip: stream.ip || '192.168.1.100',
      port: stream.port || 554, codec: stream.codec || 'H.265', resolution: stream.resolution || '1080P', fps: stream.fps || 30, bitrate: stream.bitrate || '4.0 Mbps',
      recordMode: stream.recordMode || 'continuous', status: 'online', has_ptz: false
    }
    const target = folders.value.find(f => f.id === folderId)
    if (target) target.streams.push(newStream); else rootStreams.value.push(newStream)
    selectedStream.value = newStream; showNotification('Fluxo salvo.')
  }

  return {
    searchQuery, folders, rootStreams, currentFolderId, currentFolder, selectedStream, notification, contextMenu,
    isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalStreams, displayedFolders, displayedStreams,
    openContextMenu, handleDragStart, handleDropOnFolder, moveStreamToFolder, requestDeleteFolder, confirmDeleteFolder,
    deleteFolderById, deleteStreamById, handleSaveFolder, handleSaveStream, showNotification
  }
}
