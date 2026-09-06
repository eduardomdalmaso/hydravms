import type { RondaFolderNode } from '../types/rondaTree'

export const initialRondaFolders: RondaFolderNode[] = [
  {
    id: 'rf_01',
    name: 'EMPRESA // ALPHA SEGURANCA',
    clientType: 'company',
    isExpanded: true,
    rondas: [
      {
        id: 'ronda_01',
        name: 'Ronda Portarias & Recepcao',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        folderId: 'rf_01',
        is_locked: false,
        status: 'ATIVO',
        transition: 'CORTE SECO',
        allowedUserIds: ['usr_02', 'usr_03'],
        createdAt: '2026-09-01 10:30',
        streams: [
          { id: 'rs_1', cameraId: 'cam_01', cameraName: 'Portaria Principal (Entrada)', resolution: '1080P', fps: 30, intervalSeconds: 10, orderIndex: 0 },
          { id: 'rs_2', cameraId: 'cam_02', cameraName: 'Estacionamento Visitantes', resolution: '1080P', fps: 25, intervalSeconds: 15, orderIndex: 1 },
          { id: 'rs_3', cameraId: 'cam_05', cameraName: 'Racks Servidores Sala 01', resolution: '1080P', fps: 30, intervalSeconds: 8, orderIndex: 2 }
        ]
      }
    ]
  },
  {
    id: 'rf_02',
    name: 'CLIENTE // CONDOMINIO VISTA VERDE',
    clientType: 'final_client',
    isExpanded: true,
    rondas: [
      {
        id: 'ronda_02',
        name: 'Ronda Noturna Perimetro',
        companyScope: 'CLIENTE // CONDOMINIO VISTA VERDE',
        folderId: 'rf_02',
        is_locked: true,
        status: 'ATIVO',
        transition: 'CROSSFADE',
        allowedUserIds: ['usr_03'],
        createdAt: '2026-09-02 14:15',
        streams: [
          { id: 'rs_4', cameraId: 'cam_03', cameraName: 'Corredor de Cargas & Docas', resolution: '1080P', fps: 30, intervalSeconds: 12, orderIndex: 0 },
          { id: 'rs_5', cameraId: 'cam_04', cameraName: 'Perimetro dos Fundos', resolution: '4K', fps: 30, intervalSeconds: 20, orderIndex: 1 }
        ]
      }
    ]
  }
]
