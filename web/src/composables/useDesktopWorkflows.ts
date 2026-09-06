import { ref, computed } from 'vue'
import type { WorkflowFolderNode, EnterpriseWorkflowItem } from '../types/workflowTree'
import { initialWorkflowFolders, initialRootWorkflows } from '../data/mockWorkflowFolders'

export function useDesktopWorkflows() {
  const searchQuery = ref(''), folders = ref<WorkflowFolderNode[]>(initialWorkflowFolders)
  const rootWorkflows = ref<EnterpriseWorkflowItem[]>(initialRootWorkflows)
  const currentFolderId = ref<string | null>(null), selectedWorkflow = ref<EnterpriseWorkflowItem | null>(null)
  const draggedWorkflow = ref<EnterpriseWorkflowItem | null>(null), notification = ref<string | null>(null)
  const isFolderModalOpen = ref(false), isNewWorkflowModalOpen = ref(false)

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalWorkflows = computed(() => rootWorkflows.value.length + folders.value.reduce((a, f) => a + f.workflows.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase())).sort((a, b) => a.name.localeCompare(b.name)))
  const displayedWorkflows = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.workflows : rootWorkflows.value
    return src.filter(w => !q || w.name.toLowerCase().includes(q) || w.companyScope.toLowerCase().includes(q)).sort((a, b) => a.name.localeCompare(b.name))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { if (notification.value === msg) notification.value = null }, 3500) }
  const handleDragStart = (w: EnterpriseWorkflowItem) => { draggedWorkflow.value = w }

  const handleDropOnFolder = (folderId: string) => {
    if (!draggedWorkflow.value) return
    const id = draggedWorkflow.value.id
    let item: EnterpriseWorkflowItem | undefined = rootWorkflows.value.find(w => w.id === id)
    if (!item) { for (const f of folders.value) { item = f.workflows.find(w => w.id === id); if (item) break } }
    if (!item) return
    rootWorkflows.value = rootWorkflows.value.filter(w => w.id !== id)
    folders.value.forEach(f => f.workflows = f.workflows.filter(w => w.id !== id))
    const dest = folders.value.find(f => f.id === folderId)
    if (dest) { dest.workflows.push(item); showNotification(`[WORKFLOW] "${item.name}" movido para "${dest.name}"`) }
    draggedWorkflow.value = null
  }

  const handleCreateFolder = (name: string, clientType: 'company' | 'final_client') => {
    folders.value.push({ id: `wff_${Date.now()}`, name, clientType, isExpanded: true, workflows: [] })
    isFolderModalOpen.value = false; showNotification(`[PASTA] "${name}" criada com sucesso`)
  }

  const handleCreateWorkflow = (payload: { wf: EnterpriseWorkflowItem; targetFolderId: string | null } | EnterpriseWorkflowItem) => {
    const newWf = 'wf' in payload ? payload.wf : payload
    const targetId = 'targetFolderId' in payload ? payload.targetFolderId : currentFolderId.value
    if (targetId) {
      const dest = folders.value.find(f => f.id === targetId)
      if (dest) dest.workflows.push(newWf); else rootWorkflows.value.push(newWf)
    } else {
      rootWorkflows.value.push(newWf)
    }
    isNewWorkflowModalOpen.value = false; selectedWorkflow.value = newWf
    showNotification(`[NOVO WORKFLOW] "${newWf.name}" criado com sucesso`)
  }

  const handleSaveWorkflow = (w: EnterpriseWorkflowItem) => { showNotification(`[SALVO] Workflow "${w.name}" atualizado com sucesso`) }

  const handleToggleLock = () => {
    if (!selectedWorkflow.value) return
    selectedWorkflow.value.is_locked = !selectedWorkflow.value.is_locked
    showNotification(selectedWorkflow.value.is_locked ? `[TRAVADO COM CADEADO] Workflow "${selectedWorkflow.value.name}" bloqueado` : `[DESTRAVADO] Workflow "${selectedWorkflow.value.name}" liberado`)
  }

  return {
    searchQuery, folders, rootWorkflows, currentFolderId, currentFolder, selectedWorkflow,
    draggedWorkflow, notification, isFolderModalOpen, isNewWorkflowModalOpen, totalWorkflows,
    displayedFolders, displayedWorkflows, showNotification, handleDragStart, handleDropOnFolder,
    handleCreateFolder, handleCreateWorkflow, handleSaveWorkflow, handleToggleLock
  }
}
