export interface RondaStreamItem {
  id: string
  cameraId: string
  cameraName: string
  resolution: string
  fps: number
  intervalSeconds: number
  orderIndex: number
}

export interface EnterpriseRondaItem {
  id: string
  name: string
  companyScope: string
  folderId?: string
  is_locked: boolean
  status: 'ATIVO' | 'PAUSADO'
  transition: 'CORTE SECO' | 'CROSSFADE' | 'FADE PRETO'
  allowedUserIds: string[]
  streams: RondaStreamItem[]
  createdAt: string
  description?: string
}

export interface RondaFolderNode {
  id: string
  name: string
  clientType: 'company' | 'final_client'
  isExpanded: boolean
  rondas: EnterpriseRondaItem[]
}
