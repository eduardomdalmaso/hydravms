import type { UserFolderNode, UserItem } from '../types/userTree'
import { createDefaultUserModules } from './defaultUserModules'

export const initialUserFolders: UserFolderNode[] = [
  {
    id: 'f_usr_01',
    name: 'OPERAÇÃO MATRIZ',
    isExpanded: true,
    users: [
      {
        id: 'usr_002',
        username: 'operador.matriz',
        fullName: 'Operador Principal',
        email: 'operador@hydravms.internal',
        role: 'operator',
        companyScope: 'EMPRESA ALFA (MATRIZ)',
        groupName: 'OPERAÇÃO MATRIZ',
        isActive: true,
        createdAt: '2026-09-08',
        lastLogin: 'Há 12 minutos',
        twoFactorEnabled: false,
        modules: createDefaultUserModules('operator', 'EMPRESA ALFA (MATRIZ)')
      }
    ]
  },
  {
    id: 'f_usr_02',
    name: 'SUPERVISÃO REGIONAL',
    isExpanded: true,
    users: [
      {
        id: 'usr_003',
        username: 'supervisor.sp',
        fullName: 'Supervisor Regional SP',
        email: 'supervisor.sp@hydravms.internal',
        role: 'company_admin',
        companyScope: 'REGIONAL SÃO PAULO',
        groupName: 'SUPERVISÃO REGIONAL',
        isActive: true,
        createdAt: '2026-09-06',
        lastLogin: 'Hoje, 08:30',
        twoFactorEnabled: true,
        modules: createDefaultUserModules('company_admin', 'REGIONAL SÃO PAULO')
      }
    ]
  }
]

export const initialRootUsers: UserItem[] = [
  {
    id: '00000000-0000-0000-0000-000000000002',
    username: 'admin',
    fullName: 'Admin',
    email: 'admin@hydravms.io',
    role: 'admin_master',
    companyScope: 'HYDRA MASTER OPERATIONS',
    groupName: 'Raiz (Sem Grupo)',
    isActive: true,
    createdAt: '2026-09-10',
    lastLogin: 'Agora',
    twoFactorEnabled: true,
    modules: createDefaultUserModules('admin_master', 'HYDRA MASTER OPERATIONS')
  }
]

