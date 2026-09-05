import { ref, computed } from 'vue'
import type { UserFolderNode, UserItem } from '../types/userTree'
import { initialUserFolders, initialRootUsers } from '../data/mockUserFolders'
import { createDefaultUserModules } from '../data/defaultUserModules'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'

export function useDesktopUsers() {
  const searchQuery = ref(''), folders = ref<UserFolderNode[]>(initialUserFolders), rootUsers = ref<UserItem[]>(initialRootUsers)
  const currentFolderId = ref<string | null>(null), selectedUser = ref<UserItem | null>(null), draggedUser = ref<UserItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const currentFolder = computed(() => folders.value.find(f => f.id === currentFolderId.value) || null)
  const totalUsers = computed(() => rootUsers.value.length + folders.value.reduce((acc, f) => acc + f.users.length, 0))
  const displayedFolders = computed(() => folders.value.filter(f => !searchQuery.value || f.name.toLowerCase().includes(searchQuery.value.toLowerCase())).sort((a, b) => a.name.localeCompare(b.name)))
  const displayedUsers = computed(() => {
    const q = searchQuery.value.toLowerCase(), src = currentFolder.value ? currentFolder.value.users : rootUsers.value
    return src.filter(u => !q || u.username.toLowerCase().includes(q) || u.fullName.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)).sort((a, b) => a.username.localeCompare(b.username))
  })

  const showNotification = (msg: string) => { notification.value = msg; setTimeout(() => { notification.value = null }, 3500) }
  const openContextMenu = (e: MouseEvent, target: ContextMenuTarget) => { contextMenu.value = { isOpen: true, x: e.clientX, y: e.clientY, target } }
  const handleDragStart = (u: UserItem) => { draggedUser.value = u }
  const handleDropOnFolder = (targetFolderId: string) => { if (draggedUser.value) { moveUserToFolder(draggedUser.value.id, targetFolderId); draggedUser.value = null } }

  const moveUserToFolder = (userId: string, targetFolderId: string | null) => {
    let user: UserItem | undefined = rootUsers.value.find(u => u.id === userId)
    if (!user) { for (const f of folders.value) { user = f.users.find(u => u.id === userId); if (user) break } }
    if (!user) return
    rootUsers.value = rootUsers.value.filter(u => u.id !== userId)
    folders.value.forEach(f => { f.users = f.users.filter(u => u.id !== userId) })
    if (targetFolderId) {
      const target = folders.value.find(f => f.id === targetFolderId)
      if (target) { user.groupName = target.name; target.users.push(user); showNotification(`Usuário "${user.username}" movido para "${target.name}".`) }
    } else {
      user.groupName = 'Raiz (Sem Grupo)'; rootUsers.value.push(user); showNotification(`Usuário "${user.username}" movido para a Raiz.`)
    }
  }

  const requestDeleteFolder = (id: string) => {
    const folder = folders.value.find(f => f.id === id)
    if (!folder) return
    if (folder.users.length > 0) openDeletePrompt(folder.id, folder.name, folder.users.length)
    else deleteFolderById(id)
  }

  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const folder = folders.value.find(f => f.id === folderToDelete.value!.id)
    if (folder) {
      folder.users.forEach(u => { u.groupName = 'Raiz (Sem Grupo)' })
      rootUsers.value.push(...folder.users)
      folders.value = folders.value.filter(f => f.id !== folder.id)
      if (currentFolderId.value === folder.id) currentFolderId.value = null
      showNotification(`Grupo removido. ${folder.users.length} usuários movidos para a raiz.`)
    }
    closeDeletePrompt()
  }

  const deleteFolderById = (id: string) => { folders.value = folders.value.filter(f => f.id !== id); if (currentFolderId.value === id) currentFolderId.value = null; showNotification('Grupo removido.') }
  const deleteUserById = (id: string) => {
    rootUsers.value = rootUsers.value.filter(u => u.id !== id); folders.value.forEach(f => { f.users = f.users.filter(u => u.id !== id) })
    if (selectedUser.value?.id === id) selectedUser.value = null
    showNotification('Usuário removido.')
  }

  const handleSaveFolder = (name: string) => { folders.value.push({ id: `f_usr_0${folders.value.length + 1}`, name, isExpanded: true, users: [] }); showNotification(`Grupo "${name}" criado.`) }
  const handleSaveUser = (user: Partial<UserItem>, folderId?: string) => {
    const role = user.role || 'operator', scope = user.companyScope || 'EMPRESA PADRAO'
    const newUser: UserItem = {
      id: `usr_0${totalUsers.value + 1}`, username: user.username || 'novo_usuario', fullName: user.fullName || 'Usuário Operador',
      email: user.email || `${user.username || 'user'}@hydravms.internal`, role, companyScope: scope,
      groupName: 'Raiz (Sem Grupo)', isActive: true, createdAt: '2026-09-05', lastLogin: 'Nunca', twoFactorEnabled: false,
      modules: createDefaultUserModules(role, scope)
    }
    const target = folders.value.find(f => f.id === folderId)
    if (target) { newUser.groupName = target.name; target.users.push(newUser) } else rootUsers.value.push(newUser)
    selectedUser.value = newUser; showNotification(`Usuário "${newUser.username}" cadastrado.`)
  }

  return {
    searchQuery, folders, rootUsers, currentFolderId, currentFolder, selectedUser, draggedUser, notification, contextMenu,
    isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalUsers, displayedFolders, displayedUsers,
    openContextMenu, handleDragStart, handleDropOnFolder, moveUserToFolder, requestDeleteFolder, confirmDeleteFolder,
    deleteFolderById, deleteUserById, handleSaveFolder, handleSaveUser, showNotification
  }
}
