import type { GridLayout } from './mosaic'

export interface LayoutSlotItem {
  slotIndex: number
  cameraId?: string
  cameraName?: string
  streamUrl?: string
}

export interface EnterpriseLayoutItem {
  id: string
  name: string
  grid: GridLayout
  companyScope: string
  folderId?: string
  is_locked: boolean
  created_by: string
  createdAt: string
  targetScope: 'all_company' | 'specific_users'
  allowedUserIds: string[]
  slots: LayoutSlotItem[]
  description?: string
}

export interface LayoutFolderNode {
  id: string
  name: string
  clientType: 'company' | 'final_client'
  isExpanded: boolean
  layouts: EnterpriseLayoutItem[]
}
