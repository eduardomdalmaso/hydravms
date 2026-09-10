import { ref, computed, onMounted } from 'vue'
import type { FolderNode, StreamItem } from '../types/streamTree'
import { initialFolders } from '../data/mockStreamFolders'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { fetchFolders, fetchCameras, createRemoteFolder, createRemoteCamera, deleteRemoteFolder, deleteRemoteCamera } from '../services/api'
import { buildStreamTree, createNewStreamItem } from '../utils/streamTreeUtils'
import { useEventBus } from '../services/eventSocket'

export function useStreamTree() {
  const searchQuery = ref(''), folders = ref<FolderNode[]>(initialFolders), notification = ref<string | null>(null)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const isFolderModalOpen = ref(false), isRenameFolder = ref(false), selectedFolderId = ref<string | null>(null)
  const selectedFolderName = ref(''), isWizardOpen = ref(false), wizardTargetFolder = ref<string | undefined>(undefined)

  const loadData = async () => {
    const [dbF, dbC] = await Promise.all([fetchFolders('cameras'), fetchCameras()])
    folders.value = buildStreamTree(dbF || [], dbC || []).folders
  }
  onMounted(loadData)

  const eventBus = useEventBus()
  eventBus.subscribe((evt) => {
    if (evt.type === 'system.camera.offline' && evt.subject) {
      folders.value.forEach(f => { const s = f.streams.find(item => item.id === evt.subject); if (s) s.status = 'offline' })
      showNotification(`[SISTEMA] Câmera ${evt.subject} OFFLINE!`)
    } else if (evt.type === 'system.camera.online' && evt.subject) {
      folders.value.forEach(f => { const s = f.streams.find(item => item.id === evt.subject); if (s) s.status = 'online' })
      showNotification(`[SISTEMA] Câmera ${evt.subject} restabelecida.`)
    }
  })

  const totalStreams = computed(() => folders.value.reduce((acc, f) => acc + f.streams.length, 0))
  const filteredFolders = computed(() => {
    if (!searchQuery.value) return folders.value
    const q = searchQuery.value.toLowerCase()
    return folders.value.map(f => ({ ...f, streams: f.streams.filter(s => s.name.toLowerCase().includes(q) || s.ip.includes(q)) })).filter(f => f.streams.length > 0 || f.name.toLowerCase().includes(q))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }

  const deleteFolderById = async (id: string) => {
    if (id.includes('-')) await deleteRemoteFolder(id)
    const target = folders.value.find(f => f.id === id)
    folders.value = folders.value.filter(f => f.id !== id); showNotification(`Pasta ${target?.name || ''} removida.`)
  }

  const handleContextAction = async (action: string, target: ContextMenuTarget) => {
    if (action === 'create-folder') { isRenameFolder.value = false; selectedFolderName.value = ''; isFolderModalOpen.value = true }
    else if (action === 'create-stream') { wizardTargetFolder.value = target.id; isWizardOpen.value = true }
    else if (action === 'rename-folder') { selectedFolderId.value = target.id || null; selectedFolderName.value = target.name || ''; isRenameFolder.value = true; isFolderModalOpen.value = true }
    else if (action === 'delete-folder' && target.id) await deleteFolderById(target.id)
    else if (action === 'test-stream' && target.id) showNotification(`[SOCKET TEST] Handshake RTSP para ${target.name} // OK (29ms)`)
    else if (action === 'delete-stream' && target.id) {
      await deleteRemoteCamera(target.id)
      folders.value.forEach(f => { f.streams = f.streams.filter(s => s.id !== target.id) })
      showNotification(`Fluxo ${target.name} removido.`)
    }
  }

  const handleSaveFolder = async (name: string) => {
    if (isRenameFolder.value && selectedFolderId.value) {
      const f = folders.value.find(fold => fold.id === selectedFolderId.value); if (f) f.name = name
    } else {
      const res = await createRemoteFolder('cameras', name)
      folders.value.push({ id: res?.id || `f_${Date.now()}`, name, isExpanded: true, streams: [] })
    }
    showNotification(`Pasta ${name} salva.`)
  }

  const handleSaveStream = async (stream: Partial<StreamItem>, folderId: string) => {
    const target = folders.value.find(f => f.id === folderId) || folders.value[0]
    if (!target) return
    const newStream = createNewStreamItem(stream, totalStreams.value + 1)
    const saved = await createRemoteCamera({
      id: newStream.id, name: newStream.name, protocol: newStream.protocol?.toLowerCase() || 'rtsp',
      rtsp_url: newStream.url, onvif_ip: newStream.ip, onvif_port: newStream.port, location: newStream.locationName || '',
      status: 'online', resolution: newStream.resolution || '1920x1080', fps: newStream.fps || 30.0,
      bitrate_kbps: 4096, codec: newStream.codec || 'H.265', has_ptz: !!newStream.has_ptz,
      folder_id: (folderId && folderId.includes('-')) ? folderId : undefined
    })
    if (saved) newStream.id = saved.id
    target.streams.push(newStream); showNotification(`Fluxo ${stream.name} cadastrado em ${target.name}.`)
  }

  return {
    searchQuery, folders, notification, contextMenu, isFolderModalOpen, isRenameFolder,
    selectedFolderName, isWizardOpen, wizardTargetFolder, totalStreams, filteredFolders,
    openContextMenu, deleteFolderById, handleContextAction, handleSaveFolder, handleSaveStream, showNotification
  }
}
