export type UserRole = 'admin_master' | 'company_admin' | 'operator' | 'client_viewer'

export interface UserSubPermission {
  id: string
  name: string
  isEnabled: boolean
}

export interface UserPermissionRule {
  id: string
  name: string
  category: string
  isActive: boolean
  scope: string
}

export interface UserModulePermission {
  id: string
  name: string
  isEnabled: boolean
  description: string
  scopeTarget?: string
  permissions: UserSubPermission[]
}

export interface UserItem {
  id: string
  username: string
  fullName: string
  email: string
  role: UserRole
  companyScope: string
  groupName: string
  isActive: boolean
  createdAt: string
  lastLogin: string
  twoFactorEnabled: boolean
  modules?: UserModulePermission[]
}

export interface UserFolderNode {
  id: string
  name: string
  isExpanded: boolean
  users: UserItem[]
}
