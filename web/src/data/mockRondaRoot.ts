import type { EnterpriseRondaItem } from '../types/rondaTree'

export const initialRootRondas: EnterpriseRondaItem[] = [
  {
    id: 'ronda_root_01',
    name: 'Ronda Geral das Docas',
    companyScope: 'GLOBAL // TODAS AS UNIDADES',
    is_locked: false,
    status: 'ATIVO',
    transition: 'CORTE SECO',
    allowedUserIds: ['usr_02', 'usr_03', 'usr_05'],
    createdAt: '2026-09-03 09:00',
    streams: [
      { id: 'rs_root_1', cameraId: 'cam_root_01', cameraName: 'Almoxarifado Geral', resolution: '1080P', fps: 30, intervalSeconds: 10, orderIndex: 0 },
      { id: 'rs_root_2', cameraId: 'cam_root_02', cameraName: 'Refeitorio Central', resolution: '1080P', fps: 25, intervalSeconds: 10, orderIndex: 1 }
    ]
  }
]
