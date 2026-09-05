import type { LayoutFolderNode } from '../types/layoutTree'

export const initialLayoutFolders: LayoutFolderNode[] = [
  {
    id: 'f_lay_01',
    name: 'EMPRESA // ALPHA SEGURANCA',
    clientType: 'company',
    isExpanded: true,
    layouts: [
      {
        id: 'lay_alpha_01',
        name: 'Mural Alpha Operacional 2x2',
        grid: '2x2',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        folderId: 'f_lay_01',
        is_locked: true,
        created_by: 'adminMaster',
        createdAt: '2026-03-05',
        targetScope: 'specific_users',
        allowedUserIds: ['usr_02', 'usr_03'],
        slots: [
          { slotIndex: 0, cameraId: 'cam_01', cameraName: 'Portaria Principal' },
          { slotIndex: 1, cameraId: 'cam_02', cameraName: 'Estacionamento Visitantes' },
          { slotIndex: 2, cameraId: 'cam_03', cameraName: 'Corredor de Cargas' },
          { slotIndex: 3, cameraId: 'cam_05', cameraName: 'Racks Servidores NOC' }
        ]
      },
      {
        id: 'lay_alpha_02',
        name: 'Destaque Docas & Logística 1+5',
        grid: '1+5',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        folderId: 'f_lay_01',
        is_locked: true,
        created_by: 'adminMaster',
        createdAt: '2026-03-06',
        targetScope: 'all_company',
        allowedUserIds: ['usr_02'],
        slots: [
          { slotIndex: 0, cameraId: 'cam_03', cameraName: 'Docas' },
          { slotIndex: 1, cameraId: 'cam_01', cameraName: 'Portaria' },
          { slotIndex: 2, cameraId: 'cam_02', cameraName: 'Estacionamento' },
          { slotIndex: 3, cameraId: 'cam_04', cameraName: 'Perímetro' },
          { slotIndex: 4, cameraId: 'cam_root_01', cameraName: 'Almoxarifado' },
          { slotIndex: 5, cameraId: 'cam_root_02', cameraName: 'Refeitório' }
        ]
      }
    ]
  },
  {
    id: 'f_lay_02',
    name: 'CLIENTE // CONDOMINIO VISTA VERDE',
    clientType: 'final_client',
    isExpanded: true,
    layouts: [
      {
        id: 'lay_vv_01',
        name: 'Portaria & Acessos Vista Verde',
        grid: '2x2',
        companyScope: 'CLIENTE // CONDOMINIO VISTA VERDE',
        folderId: 'f_lay_02',
        is_locked: true,
        created_by: 'gestor_alpha_seguranca',
        createdAt: '2026-03-10',
        targetScope: 'specific_users',
        allowedUserIds: ['usr_04', 'usr_03'],
        slots: [
          { slotIndex: 0, cameraId: 'cam_01', cameraName: 'Portaria Principal' },
          { slotIndex: 1, cameraId: 'cam_02', cameraName: 'Estacionamento' }
        ]
      }
    ]
  },
  {
    id: 'f_lay_03',
    name: 'CLIENTE // EDIFICIO HORIZONTE',
    clientType: 'final_client',
    isExpanded: false,
    layouts: [
      {
        id: 'lay_horiz_01',
        name: 'Visão Geral Portaria Horizonte',
        grid: '1x2',
        companyScope: 'CLIENTE // EDIFICIO HORIZONTE',
        folderId: 'f_lay_03',
        is_locked: true,
        created_by: 'gestor_alpha_seguranca',
        createdAt: '2026-03-12',
        targetScope: 'specific_users',
        allowedUserIds: ['usr_03'],
        slots: [
          { slotIndex: 0, cameraId: 'cam_01', cameraName: 'Portaria Principal' },
          { slotIndex: 1, cameraId: 'cam_04', cameraName: 'Perímetro' }
        ]
      }
    ]
  }
]
