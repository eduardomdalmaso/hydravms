import type { UserFolderNode, UserItem } from '../types/userTree'
import { createDefaultUserModules } from './defaultUserModules'

export const initialRootUsers: UserItem[] = [
  {
    id: 'usr_04',
    username: 'cliente_condominio',
    fullName: 'Cliente Final // Condomínio Vista Verde',
    email: 'sindico@vistaverde.corp',
    role: 'client_viewer',
    companyScope: 'EMPRESA: ALPHA // CLIENTE: CONDOMINIO VISTA VERDE',
    groupName: 'Raiz (Sem Grupo)',
    isActive: true,
    createdAt: '2026-03-10',
    lastLogin: 'Hoje 09:12',
    twoFactorEnabled: true,
    modules: createDefaultUserModules('client_viewer', 'CONDOMINIO VISTA VERDE (4 CAMERAS)')
  },
  {
    id: 'usr_05',
    username: 'auditor_externo_ti',
    fullName: 'Auditor Externo de Segurança',
    email: 'auditoria@ti-secure.com',
    role: 'operator',
    companyScope: 'GLOBAL // AUDITORIA READ-ONLY',
    groupName: 'Raiz (Sem Grupo)',
    isActive: true,
    createdAt: '2026-04-01',
    lastLogin: 'Ontem 17:45',
    twoFactorEnabled: true,
    modules: createDefaultUserModules('operator', 'LOGS & VISUALIZACAO GERAL')
  }
]

export const initialUserFolders: UserFolderNode[] = [
  {
    id: 'f_usr_01',
    name: 'ADMINISTRADORES // CONTROLE TOTAL',
    isExpanded: true,
    users: [
      {
        id: 'usr_01',
        username: 'adminMaster',
        fullName: 'Administrador Master da Plataforma',
        email: 'master@hydravms.internal',
        role: 'admin_master',
        companyScope: 'GLOBAL // TODAS AS EMPRESAS & CLIENTES',
        groupName: 'ADMINISTRADORES',
        isActive: true,
        createdAt: '2026-01-10',
        lastLogin: 'Agora (Online)',
        twoFactorEnabled: true,
        modules: createDefaultUserModules('admin_master', 'ACESSO TOTAL IRRESTRITO')
      }
    ]
  },
  {
    id: 'f_usr_02',
    name: 'GESTAO DE EMPRESAS // CLIENTES',
    isExpanded: true,
    users: [
      {
        id: 'usr_02',
        username: 'gestor_alpha_seguranca',
        fullName: 'Gestor Regional Alpha Segurança',
        email: 'gestor@alphaseg.com',
        role: 'company_admin',
        companyScope: 'EMPRESA: ALPHA SEGURANCA // 12 CLIENTES',
        groupName: 'GESTAO DE EMPRESAS',
        isActive: true,
        createdAt: '2026-02-01',
        lastLogin: 'Hoje 10:15',
        twoFactorEnabled: true,
        modules: createDefaultUserModules('company_admin', 'ALPHA SEGURANCA (12 CLIENTES)')
      }
    ]
  },
  {
    id: 'f_usr_03',
    name: 'OPERADORES DA PORTARIA',
    isExpanded: true,
    users: [
      {
        id: 'usr_03',
        username: 'operador_portaria',
        fullName: 'Operador Central Portaria 01',
        email: 'portaria@hydravms.internal',
        role: 'operator',
        companyScope: 'CLIENTE: EDIFICIO HORIZONTE // CAM_01 A CAM_08',
        groupName: 'OPERADORES DA PORTARIA',
        isActive: true,
        createdAt: '2026-02-15',
        lastLogin: 'Hoje 08:30',
        twoFactorEnabled: false,
        modules: createDefaultUserModules('operator', 'EDIFICIO HORIZONTE')
      }
    ]
  }
]
