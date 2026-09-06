export type MapKind = 'MAP_OPENSOURCE' | 'PLANTA_BAIXA'
export type DeviceType = 'CAMERA' | 'ALARME'

export interface MapMarkerItem {
  id: string
  deviceId: string
  name: string
  type: DeviceType
  lat: number
  lng: number
  angle?: number
  status: 'ONLINE' | 'OFFLINE' | 'DISPARADO'
}

export interface EnterpriseMapItem {
  id: string
  name: string
  type: MapKind
  companyScope: string
  folderId?: string
  is_locked: boolean
  initialCenter: [number, number]
  initialZoom: number
  floorplanFile?: string
  markers: MapMarkerItem[]
  createdAt: string
  description?: string
}

export interface MapFolderNode {
  id: string
  name: string
  clientType: 'company' | 'final_client'
  isExpanded: boolean
  maps: EnterpriseMapItem[]
}

export interface DeployableDevice {
  id: string
  name: string
  type: DeviceType
  hasCoords: boolean
  lat?: number
  lng?: number
  status: string
  details?: string
}
