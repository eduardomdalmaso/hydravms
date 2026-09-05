import { ref, computed } from 'vue'
import type { AlarmFolderNode, AlarmItem } from '../types/alarmTree'
import { initialAlarmFolders, initialRootAlarms } from '../data/mockAlarmFolders'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'

export function useDesktopAlarms() {
  const searchQuery = ref(''), folders = ref<AlarmFolderNode[]>(initialAlarmFolders), rootAlarms = ref<AlarmItem[]>(initialRootAlarms)
  const currentFolderId = ref<string | null>(null), selectedAlarm = ref<AlarmItem | null>(null), draggedAlarm = ref<AlarmItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalAlarms = computed(() => rootAlarms.value.length + folders.value.reduce((acc, f) => acc + f.alarms.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase())).sort((a, b) => a.name.localeCompare(b.name)))
  const displayedAlarms = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.alarms : rootAlarms.value
    return src.filter(a => !q || a.name.toLowerCase().includes(q) || a.zone.toLowerCase().includes(q)).sort((a, b) => a.name.localeCompare(b.name))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }
  const handleDragStart = (a: AlarmItem) => { draggedAlarm.value = a }
  const handleDropOnFolder = (targetFolderId: string) => { if (draggedAlarm.value) { moveAlarmToFolder(draggedAlarm.value.id, targetFolderId); draggedAlarm.value = null } }

  const moveAlarmToFolder = (alarmId: string, targetFolderId: string | null) => {
    let alarm: AlarmItem | undefined = rootAlarms.value.find(a => a.id === alarmId)
    if (!alarm) { for (const f of folders.value) { alarm = f.alarms.find(a => a.id === alarmId); if (alarm) break } }
    if (!alarm) return
    rootAlarms.value = rootAlarms.value.filter(a => a.id !== alarmId)
    folders.value.forEach(f => { f.alarms = f.alarms.filter(a => a.id !== alarmId) })
    if (targetFolderId) {
      const target = folders.value.find(f => f.id === targetFolderId)
      if (target) { target.alarms.push(alarm); showNotification(`Sensor "${alarm.name}" movido para "${target.name}".`) }
    } else {
      rootAlarms.value.push(alarm)
      showNotification(`Sensor "${alarm.name}" movido para a Raiz.`)
    }
  }

  const requestDeleteFolder = (id: string) => {
    const folder = folders.value.find(f => f.id === id)
    if (!folder) return
    if (folder.alarms.length > 0) openDeletePrompt(folder.id, folder.name, folder.alarms.length)
    else deleteFolderById(id)
  }

  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const folder = folders.value.find(f => f.id === folderToDelete.value!.id)
    if (folder) {
      rootAlarms.value.push(...folder.alarms)
      folders.value = folders.value.filter(f => f.id !== folder.id)
      if (currentFolderId.value === folder.id) currentFolderId.value = null
      showNotification(`Zona removida. ${folder.alarms.length} sensores movidos para a raiz.`)
    }
    closeDeletePrompt()
  }

  const deleteFolderById = (id: string) => { folders.value = folders.value.filter(f => f.id !== id); if (currentFolderId.value === id) currentFolderId.value = null; showNotification('Zona removida.') }
  const deleteAlarmById = (id: string) => {
    rootAlarms.value = rootAlarms.value.filter(a => a.id !== id); folders.value.forEach(f => { f.alarms = f.alarms.filter(a => a.id !== id) })
    if (selectedAlarm.value?.id === id) selectedAlarm.value = null
    showNotification('Sensor removido.')
  }

  const handleSaveFolder = (name: string) => { folders.value.push({ id: `f_alm_0${folders.value.length + 1}`, name, isExpanded: true, alarms: [] }); showNotification(`Zona "${name}" criada.`) }
  const handleSaveAlarm = (alarm: Partial<AlarmItem>, folderId?: string) => {
    const newAlarm: AlarmItem = {
      id: `alm_0${totalAlarms.value + 1}`, name: alarm.name || 'Novo Sensor', zone: alarm.zone || 'Zona Geral', type: alarm.type || 'IVS',
      status: 'online', sensitivity: alarm.sensitivity || 80, linkedCameraId: alarm.linkedCameraId || 'cam_01', linkedCameraName: alarm.linkedCameraName || 'CAM_01 Portaria', lastTrigger: 'Agora'
    }
    const target = folders.value.find(f => f.id === folderId)
    if (target) target.alarms.push(newAlarm); else rootAlarms.value.push(newAlarm)
    selectedAlarm.value = newAlarm; showNotification(`Sensor "${newAlarm.name}" cadastrado.`)
  }

  return {
    searchQuery, folders, rootAlarms, currentFolderId, currentFolder, selectedAlarm, draggedAlarm, notification, contextMenu,
    isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalAlarms, displayedFolders, displayedAlarms,
    openContextMenu, handleDragStart, handleDropOnFolder, moveAlarmToFolder, requestDeleteFolder, confirmDeleteFolder,
    deleteFolderById, deleteAlarmById, handleSaveFolder, handleSaveAlarm, showNotification
  }
}
