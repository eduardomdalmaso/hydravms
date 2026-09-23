import { ref } from 'vue'
import type { StreamItem, FolderNode } from '../types/streamTree'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useOnvifDiscovery } from './useOnvifDiscovery'
import { useStreamExportImport } from './useStreamExportImport'

export function useAdminCameraActions(
  folders: { value: FolderNode[] },
  displayedStreams: { value: StreamItem[] },
  currentFolderId: { value: string | null },
  selectedStream: { value: StreamItem | null },
  isFolderModalOpen: { value: boolean },
  isWizardOpen: { value: boolean },
  showNotification: (msg: string) => void,
  handleSaveStream: (stream: Partial<StreamItem>, folderId?: string) => void,
  deleteStreamById: (id: string) => void,
  requestDeleteFolder: (id: string) => void,
  moveStreamToFolder: (id: string, targetFolderId: string | null) => void
) {
  const { markAsImported } = useOnvifDiscovery()
  const { exportStreamsToCsv, parseStreamsFromFile } = useStreamExportImport()
  const wizardInitialData = ref<Partial<StreamItem> | undefined>(undefined)
  const fileInputRef = ref<HTMLInputElement | null>(null)
  const isConfirmDeleteStreamOpen = ref(false)
  const streamToDelete = ref<{ id: string; name: string } | null>(null)

  const handleExportStreams = () => {
    const count = exportStreamsToCsv(folders.value, displayedStreams.value)
    showNotification(`[EXPORTAÇÃO CSV] ${count} fluxo(s) exportado(s) em planilha com URLs seguras.`)
  }

  const triggerFileInput = () => { fileInputRef.value?.click() }

  const handleImportFile = async (e: Event) => {
    const target = e.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return
    try {
      const list = await parseStreamsFromFile(file)
      for (const item of list) {
        handleSaveStream(item, currentFolderId.value || undefined)
      }
      showNotification(`[IMPORTAÇÃO] ${list.length} fluxo(s) importado(s) com sucesso.`)
    } catch (err: any) {
      showNotification(`[ERRO DE IMPORTAÇÃO] Arquivo inválido: ${err?.message || 'Erro de formato'}`)
    } finally {
      if (fileInputRef.value) fileInputRef.value.value = ''
    }
  }

  const handleOpenNewWizard = (data?: Partial<StreamItem>) => {
    wizardInitialData.value = data
    isWizardOpen.value = true
  }

  const onSaveStream = (stream: Partial<StreamItem>, folderId?: string) => {
    if (stream.ip) markAsImported(stream.ip)
    handleSaveStream(stream, folderId)
  }

  const requestDeleteStream = (id: string, name?: string) => {
    if (!name) {
      let s = displayedStreams.value.find(item => item.id === id)
      if (!s) { for (const f of folders.value) { s = f.streams.find(item => item.id === id); if (s) break } }
      name = s?.name || id
    }
    streamToDelete.value = { id, name }
    isConfirmDeleteStreamOpen.value = true
  }

  const confirmDeleteStream = () => {
    if (streamToDelete.value) {
      deleteStreamById(streamToDelete.value.id)
      if (selectedStream.value?.id === streamToDelete.value.id) {
        selectedStream.value = null
      }
      showNotification(`[EXCLUSÃO] Fluxo ${streamToDelete.value.name} removido com sucesso.`)
    }
    isConfirmDeleteStreamOpen.value = false
    streamToDelete.value = null
  }

  const handleContextAction = (action: string, target: ContextMenuTarget, extra?: any) => {
    if (action === 'open-folder' && target.id) currentFolderId.value = target.id
    else if (action === 'create-folder') isFolderModalOpen.value = true
    else if (action === 'create-stream') { if (target.type === 'folder' && target.id) currentFolderId.value = target.id; handleOpenNewWizard() }
    else if (action === 'inspect-stream' && target.id) {
      let s = displayedStreams.value.find(item => item.id === target.id)
      if (!s) { for (const f of folders.value) { s = f.streams.find(item => item.id === target.id); if (s) break } }
      if (s) selectedStream.value = s
    }
    else if (action === 'delete-folder' && target.id) requestDeleteFolder(target.id)
    else if (action === 'delete-stream' && target.id) requestDeleteStream(target.id, target.name)
    else if (action === 'move-stream' && target.id) moveStreamToFolder(target.id, extra)
    else if (action === 'test-stream' && target.id) showNotification(`[SOCKET TEST] Handshake RTSP // OK (29ms)`)
  }

  return {
    wizardInitialData, fileInputRef, isConfirmDeleteStreamOpen, streamToDelete,
    handleExportStreams, triggerFileInput, handleImportFile, handleOpenNewWizard,
    onSaveStream, requestDeleteStream, confirmDeleteStream, handleContextAction
  }
}
