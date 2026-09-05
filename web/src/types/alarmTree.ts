export type AlarmSensorType = 'IVS' | 'PIR' | 'MAG' | 'SMK'
export type AlarmSensorStatus = 'online' | 'unarmed' | 'alert' | 'offline'

export interface AlarmZoneRule {
  id: string
  name: string
  actionType: string
  isActive: boolean
  delaySeconds: number
  targetOutput?: string
}

export interface AlarmItem {
  id: string
  name: string
  zone: string
  type: AlarmSensorType
  status: AlarmSensorStatus
  sensitivity: number
  linkedCameraId?: string
  linkedCameraName?: string
  lastTrigger?: string
  rules?: AlarmZoneRule[]
}

export interface AlarmFolderNode {
  id: string
  name: string
  isExpanded: boolean
  alarms: AlarmItem[]
}
