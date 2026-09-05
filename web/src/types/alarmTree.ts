export type AlarmSensorType = 'IVS' | 'PIR' | 'MAG' | 'SMK'
export type AlarmSensorStatus = 'online' | 'unarmed' | 'alert' | 'offline'

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
}

export interface AlarmFolderNode {
  id: string
  name: string
  isExpanded: boolean
  alarms: AlarmItem[]
}
