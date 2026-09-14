import { ref, computed, onMounted, watch } from 'vue'
import type { UserFolderNode, UserItem } from '../types/userTree'
import { initialUserFolders, initialRootUsers } from '../data/mockUserFolders'
import { createDefaultUserModules } from '../data/defaultUserModules'
import type { ContextMenuTarget } from '../components/admin/TreeContextMenu.vue'
import { useFolderModalState } from './useDesktopFolderOps'
import { fetchFolders, fetchUsers } from '../services/api'

const STORAGE_KEY = 'hydravms_admin_users_v2'

export function useDesktopUsers() {
  const searchQuery = ref(''), folders = ref<UserFolderNode[]>(initialUserFolders), rootUsers = ref<UserItem[]>(initialRootUsers)
  const currentFolderId = ref<string | null>(null), selectedUser = ref<UserItem | null>(null), draggedUser = ref<UserItem | null>(null)
  const notification = ref<string | null>(null), isFolderModalOpen = ref(false), isWizardOpen = ref(false)
  const contextMenu = ref<{ isOpen: boolean; x: number; y: number; target: ContextMenuTarget }>({ isOpen: false, x: 0, y: 0, target: { type: 'canvas' } })
  const { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt } = useFolderModalState()

  const persistState = () => {
    if (typeof window === 'undefined') return
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ folders: folders.value, rootUsers: rootUsers.value }))
    } catch {}
  }

  onMounted(async () => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        try {
          const parsed = JSON.parse(saved)
          if (parsed && Array.isArray(parsed.folders) && Array.isArray(parsed.rootUsers) && (parsed.folders.length > 0 || parsed.rootUsers.length > 0)) {
            folders.value = parsed.folders
            rootUsers.value = parsed.rootUsers
            return
          }
        } catch {}
      }
    }

    const [dbF, dbU] = await Promise.all([fetchFolders('users'), fetchUsers(initialRootUsers)])
    const folderList: UserFolderNode[] = (dbF && dbF.length > 0) ? dbF.map(f => ({ id: f.id, name: f.name, isExpanded: true, users: [] })) : [...initialUserFolders]
    const unassigned: UserItem[] = [], rawList = (dbU && dbU.length > 0) ? dbU : initialRootUsers

    rawList.forEach((raw: any) => {
      const role = (raw.role === 'admin' || raw.role === 'super_admin' || raw.role === 'admin_master') ? 'admin_master' : (raw.role || 'operator')
      const scope = raw.companyScope || 'HYDRA MASTER OPERATIONS'
      const uItem: UserItem = {
        id: raw.id || `usr_${Date.now()}`, username: raw.username || raw.email || 'usuario', fullName: raw.fullName || raw.name || raw.username || 'Usuário',
        email: raw.email || `${raw.username || 'user'}@hydravms.io`, role: role as any, companyScope: scope, groupName: raw.groupName || 'Raiz (Sem Grupo)',
        isActive: raw.isActive !== false, createdAt: raw.createdAt || '2026-09-10', lastLogin: raw.lastLogin || 'Recentemente', twoFactorEnabled: raw.twoFactorEnabled !== false,
        timezone: raw.timezone || 'America/Sao_Paulo (UTC-03:00)', phonePrimary: raw.phonePrimary || '+55 (11) 98765-4321', phoneSecondary: raw.phoneSecondary || '',
        modules: raw.modules || createDefaultUserModules(role as any, scope)
      }
      const target = folderList.find(f => f.name.toLowerCase() === (uItem.groupName || '').toLowerCase())
      if (target) target.users.push(uItem); else unassigned.push(uItem)
    })
    folders.value = folderList; rootUsers.value = unassigned
    persistState()
  })

  watch([folders, rootUsers], () => { persistState() }, { deep: true })

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
    rootUsers.value = rootUsers.value.filter(u => u.id !== userId); folders.value.forEach(f => { f.users = f.users.filter(u => u.id !== userId) })
    if (targetFolderId) {
      const target = folders.value.find(f => f.id === targetFolderId)
      if (target) { user.groupName = target.name; target.users.push(user); showNotification(`Usuário "${user.username}" movido para "${target.name}".`) }
    } else { user.groupName = 'Raiz (Sem Grupo)'; rootUsers.value.push(user); showNotification(`Usuário "${user.username}" movido para a Raiz.`) }
  }

  const requestDeleteFolder = (id: string) => { const f = folders.value.find(item => item.id === id); if (f && f.users.length > 0) openDeletePrompt(f.id, f.name, f.users.length); else deleteFolderById(id) }
  const confirmDeleteFolder = () => {
    if (!folderToDelete.value) return
    const f = folders.value.find(item => item.id === folderToDelete.value!.id)
    if (f) { f.users.forEach(u => { u.groupName = 'Raiz (Sem Grupo)' }); rootUsers.value.push(...f.users); folders.value = folders.value.filter(item => item.id !== f.id); if (currentFolderId.value === f.id) currentFolderId.value = null; showNotification(`Grupo removido. ${f.users.length} usuários movidos para a raiz.`) }
    closeDeletePrompt()
  }

  const deleteFolderById = (id: string) => { folders.value = folders.value.filter(f => f.id !== id); if (currentFolderId.value === id) currentFolderId.value = null; showNotification('Grupo removido.') }
  const deleteUserById = (id: string) => {
    rootUsers.value = rootUsers.value.filter(u => u.id !== id); folders.value.forEach(f => { f.users = f.users.filter(u => u.id !== id) })
    if (selectedUser.value?.id === id) selectedUser.value = null; showNotification('Usuário removido.')
  }

  const handleSaveFolder = (name: string) => { folders.value.push({ id: `f_usr_0${folders.value.length + 1}`, name, isExpanded: true, users: [] }); showNotification(`Grupo "${name}" criado.`) }
  const handleSaveUser = (user: Partial<UserItem>, folderId?: string) => {
    const role = user.role || 'operator', scope = user.companyScope || 'EMPRESA PADRÃO'
    const newUser: UserItem = {
      id: `usr_0${totalUsers.value + 1}`, username: user.username || 'novo_usuario', fullName: user.fullName || 'Usuário Operador',
      email: user.email || `${user.username || 'user'}@hydravms.internal`, role, companyScope: scope, groupName: 'Raiz (Sem Grupo)',
      isActive: true, createdAt: '2026-09-05', lastLogin: 'Nunca', twoFactorEnabled: false, modules: createDefaultUserModules(role, scope)
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
