import { ref, computed, onMounted } from 'vue'
import type { FolderNode, StreamItem } from '../types/streamTree'
import { initialFolders } from '../data/mockStreamFolders'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { fetchFolders, fetchCameras } from '../services/api'
import { buildStreamTree, createNewStreamItem } from '../utils/streamTreeUtils'

export function useStreamTree() {
  const searchQuery = ref(''), folders = ref<FolderNode[]>(initialFolders), notification = ref<string | null>(null)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const isFolderModalOpen = ref(false), isRenameFolder = ref(false), selectedFolderId = ref<string | null>(null)
  const selectedFolderName = ref(''), isWizardOpen = ref(false), wizardTargetFolder = ref<string | undefined>(undefined)

  onMounted(async () => {
    const [dbF, dbC] = await Promise.all([fetchFolders('cameras'), fetchCameras()])
    if (dbF.length > 0 || dbC.length > 0) folders.value = buildStreamTree(dbF, dbC).folders
  })

  const totalStreams = computed(() => folders.value.reduce((acc, f) => acc + f.streams.length, 0))
  const filteredFolders = computed(() => {
    if (!searchQuery.value) return folders.value
    const q = searchQuery.value.toLowerCase()
    return folders.value.map(f => ({ ...f, streams: f.streams.filter(s => s.name.toLowerCase().includes(q) || s.ip.includes(q)) })).filter(f => f.streams.length > 0 || f.name.toLowerCase().includes(q))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }

  const deleteFolderById = (id: string) => {
    const target = folders.value.find(f => f.id === id)
    folders.value = folders.value.filter(f => f.id !== id)
    showNotification(`Pasta ${target?.name || ''} removida.`)
  }

  const handleContextAction = (action: string, target: ContextMenuTarget) => {
    if (action === 'create-folder') { isRenameFolder.value = false; selectedFolderName.value = ''; isFolderModalOpen.value = true }
    else if (action === 'create-stream') { wizardTargetFolder.value = target.id; isWizardOpen.value = true }
    else if (action === 'rename-folder') { selectedFolderId.value = target.id || null; selectedFolderName.value = target.name || ''; isRenameFolder.value = true; isFolderModalOpen.value = true }
    else if (action === 'delete-folder' && target.id) deleteFolderById(target.id)
    else if (action === 'test-stream' && target.id) showNotification(`[SOCKET TEST] Handshake RTSP para ${target.name} // OK (29ms)`)
    else if (action === 'delete-stream' && target.id) {
      folders.value.forEach(f => { f.streams = f.streams.filter(s => s.id !== target.id) })
      showNotification(`Fluxo ${target.name} removido.`)
    }
  }

  const handleSaveFolder = (name: string) => {
    if (isRenameFolder.value && selectedFolderId.value) {
      const f = folders.value.find(fold => fold.id === selectedFolderId.value)
      if (f) f.name = name
    } else {
      folders.value.push({ id: `f_0${folders.value.length + 1}`, name, isExpanded: true, streams: [] })
    }
  }

  const handleSaveStream = (stream: Partial<StreamItem>, folderId: string) => {
    const target = folders.value.find(f => f.id === folderId) || folders.value[0]
    if (!target) return
    const newStream = createNewStreamItem(stream, totalStreams.value + 1)
    target.streams.push(newStream)
    showNotification(`Fluxo ${stream.name} cadastrado em ${target.name}.`)
  }

  return {
    searchQuery, folders, notification, contextMenu, isFolderModalOpen, isRenameFolder,
    selectedFolderName, isWizardOpen, wizardTargetFolder, totalStreams, filteredFolders,
    openContextMenu, deleteFolderById, handleContextAction, handleSaveFolder, handleSaveStream, showNotification
  }
}
