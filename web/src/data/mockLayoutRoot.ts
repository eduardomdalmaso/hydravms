import type { EnterpriseLayoutItem } from '../types/layoutTree'

export const initialRootLayouts: EnterpriseLayoutItem[] = [
  {
    id: 'lay_root_01',
    name: 'Mosaico Global Padrão (2x2)',
    grid: '2x2',
    companyScope: 'GLOBAL // TODAS AS EMPRESAS',
    is_locked: true,
    created_by: 'adminMaster',
    createdAt: '2026-03-01',
    targetScope: 'all_company',
    allowedUserIds: ['usr_01', 'usr_02', 'usr_03', 'usr_04', 'usr_05'],
    slots: [
      { slotIndex: 0, cameraId: 'cam_01', cameraName: 'Portaria Principal' },
      { slotIndex: 1, cameraId: 'cam_02', cameraName: 'Estacionamento Visitantes' },
      { slotIndex: 2, cameraId: 'cam_03', cameraName: 'Corredor de Cargas' },
      { slotIndex: 3, cameraId: 'cam_04', cameraName: 'Perímetro dos Fundos' }
    ]
  }
]
