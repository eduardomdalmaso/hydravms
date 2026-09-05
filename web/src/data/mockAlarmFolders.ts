import type { AlarmFolderNode, AlarmItem } from '../types/alarmTree'

export const initialAlarmFolders: AlarmFolderNode[] = [
  {
    id: 'f_alm_01',
    name: 'PERIMETRO EXTERNO',
    isExpanded: true,
    alarms: [
      { id: 'alm_01', name: 'Sensor Muro Norte', zone: 'Zona 01 // Muro Norte', type: 'IVS', status: 'online', sensitivity: 85, linkedCameraId: 'cam_01', linkedCameraName: 'CAM_01 Portaria', lastTrigger: '19:42:10' },
      { id: 'alm_02', name: 'Barreira Leste IVS', zone: 'Zona 01 // Cerca Leste', type: 'IVS', status: 'alert', sensitivity: 90, linkedCameraId: 'cam_04', linkedCameraName: 'CAM_04 Perimetro', lastTrigger: '20:15:00' }
    ]
  },
  {
    id: 'f_alm_02',
    name: 'GALPAO & DOCAS',
    isExpanded: true,
    alarms: [
      { id: 'alm_03', name: 'Barreira Patio Cargas', zone: 'Zona 02 // Patio Docas', type: 'PIR', status: 'online', sensitivity: 75, linkedCameraId: 'cam_03', linkedCameraName: 'CAM_03 Docas', lastTrigger: '18:30:22' },
      { id: 'alm_04', name: 'Sensor Portao Docas', zone: 'Zona 02 // Entrada Carga', type: 'MAG', status: 'online', sensitivity: 100, linkedCameraId: 'cam_03', linkedCameraName: 'CAM_03 Docas', lastTrigger: '17:10:05' }
    ]
  },
  {
    id: 'f_alm_03',
    name: 'DATACENTER & CPD',
    isExpanded: true,
    alarms: [
      { id: 'alm_05', name: 'Porta Sala Servidores', zone: 'Zona 03 // CPD Central', type: 'MAG', status: 'alert', sensitivity: 100, linkedCameraId: 'cam_02', linkedCameraName: 'CAM_02 Servidores', lastTrigger: '20:25:12' },
      { id: 'alm_06', name: 'Detector Fumaca Racks', zone: 'Zona 03 // Sala Racks', type: 'SMK', status: 'offline', sensitivity: 95, linkedCameraId: 'cam_02', linkedCameraName: 'CAM_02 Servidores', lastTrigger: 'Ontem 23:10' }
    ]
  }
]

export const initialRootAlarms: AlarmItem[] = [
  { id: 'alm_07', name: 'Sensor Hall Entrada', zone: 'Recepcao Principal', type: 'PIR', status: 'online', sensitivity: 80, linkedCameraId: 'cam_01', linkedCameraName: 'CAM_01 Portaria', lastTrigger: '20:28:40' }
]
