import type { UserFolderNode, UserItem } from '../types/userTree'
import { initialUserFolders, initialRootUsers } from '../data/mockUserFolders'
import { createDefaultUserModules } from '../data/defaultUserModules'
import { fetchFolders, fetchUsers } from '../services/api'

export async function loadUsersAndFolders(storageKey: string): Promise<{ folders: UserFolderNode[]; rootUsers: UserItem[] }> {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem(storageKey)
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        if (parsed && Array.isArray(parsed.folders) && Array.isArray(parsed.rootUsers) && (parsed.folders.length > 0 || parsed.rootUsers.length > 0)) {
          return { folders: parsed.folders, rootUsers: parsed.rootUsers }
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

  return { folders: folderList, rootUsers: unassigned }
}
