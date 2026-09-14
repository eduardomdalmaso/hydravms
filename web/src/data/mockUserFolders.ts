import type { UserFolderNode, UserItem } from '../types/userTree'
import { createDefaultUserModules } from './defaultUserModules'

export const initialUserFolders: UserFolderNode[] = []

export const initialRootUsers: UserItem[] = [
  {
    id: '00000000-0000-0000-0000-000000000002',
    username: 'admin',
    fullName: 'Administrador',
    email: 'admin@hydravms.io',
    role: 'admin_master',
    companyScope: 'HYDRA MASTER OPERATIONS',
    groupName: 'Raiz (Sem Grupo)',
    isActive: true,
    createdAt: '2026-09-10',
    lastLogin: 'Agora',
    twoFactorEnabled: true,
    timezone: 'America/Sao_Paulo (UTC-03:00)',
    phonePrimary: '',
    phoneSecondary: '',
    modules: createDefaultUserModules('admin_master', 'HYDRA MASTER OPERATIONS')
  }
]

