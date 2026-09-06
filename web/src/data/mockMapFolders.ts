import type { MapFolderNode, EnterpriseMapItem, DeployableDevice } from '../types/mapTree'

export const initialDeployableDevices: DeployableDevice[] = [
  { id: 'cam_01', name: 'Portaria Principal (Entrada)', type: 'CAMERA', hasCoords: true, lat: -23.55052, lng: -46.633308, status: 'ONLINE', details: '1080P @ 30 FPS // RTSP' },
  { id: 'cam_02', name: 'Estacionamento Visitantes', type: 'CAMERA', hasCoords: true, lat: -23.55110, lng: -46.634000, status: 'ONLINE', details: '1080P @ 25 FPS // PTZ' },
  { id: 'cam_03', name: 'Corredor de Cargas & Docas', type: 'CAMERA', hasCoords: false, status: 'ONLINE', details: '1080P @ 30 FPS // H.265' },
  { id: 'cam_04', name: 'Perimetro dos Fundos', type: 'CAMERA', hasCoords: true, lat: -23.54980, lng: -46.632200, status: 'ONLINE', details: '4K @ 30 FPS // SAHI' },
  { id: 'cam_05', name: 'Racks Datacenter Sala 01', type: 'CAMERA', hasCoords: false, status: 'ONLINE', details: '1080P @ 30 FPS // RTSP' },
  { id: 'alm_01', name: 'Sensor Barreira Infravermelho 01', type: 'ALARME', hasCoords: true, lat: -23.55010, lng: -46.632900, status: 'ONLINE', details: 'ZONA_01 // PERIMETRO' },
  { id: 'alm_02', name: 'Detector de Fumaca Galpao A', type: 'ALARME', hasCoords: false, status: 'ONLINE', details: 'ZONA_02 // INCENDIO' },
  { id: 'alm_03', name: 'Sensor Abertura Portao Docas', type: 'ALARME', hasCoords: true, lat: -23.55150, lng: -46.634500, status: 'ONLINE', details: 'MAGNETICO // DOCAS' },
  { id: 'alm_04', name: 'Sirene Alarme Perimetro Norte', type: 'ALARME', hasCoords: false, status: 'ONLINE', details: 'SIRENE // 120DB' }
]

export const initialRootMaps: EnterpriseMapItem[] = [
  {
    id: 'map_root_01',
    name: 'Perimetro Fabril & Acessos Externos',
    type: 'MAP_OPENSOURCE',
    companyScope: 'EMPRESA // ALPHA SEGURANCA',
    is_locked: false,
    initialCenter: [-23.55052, -46.633308],
    initialZoom: 16,
    markers: [
      { id: 'm_01', deviceId: 'cam_01', name: 'Portaria Principal (Entrada)', type: 'CAMERA', lat: -23.55052, lng: -46.633308, angle: 45, status: 'ONLINE' },
      { id: 'm_02', deviceId: 'alm_01', name: 'Sensor Barreira Infravermelho 01', type: 'ALARME', lat: -23.55010, lng: -46.632900, status: 'ONLINE' }
    ],
    createdAt: '2026-03-01'
  }
]

export const initialMapFolders: MapFolderNode[] = [
  {
    id: 'mf_01',
    name: 'EMPRESA // ALPHA SEGURANCA',
    clientType: 'company',
    isExpanded: true,
    maps: [
      {
        id: 'map_01',
        name: 'Matriz Sao Paulo // Centro Operacional',
        type: 'MAP_OPENSOURCE',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        is_locked: false,
        initialCenter: [-23.55052, -46.633308],
        initialZoom: 16,
        markers: [
          { id: 'm_03', deviceId: 'cam_01', name: 'Portaria Principal', type: 'CAMERA', lat: -23.55052, lng: -46.633308, angle: 90, status: 'ONLINE' },
          { id: 'm_04', deviceId: 'cam_02', name: 'Estacionamento Visitantes', type: 'CAMERA', lat: -23.55110, lng: -46.634000, angle: 180, status: 'ONLINE' }
        ],
        createdAt: '2026-03-02'
      },
      {
        id: 'map_02',
        name: 'Galpao Logistica - Setor A',
        type: 'PLANTA_BAIXA',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        is_locked: true,
        initialCenter: [-23.55000, -46.633000],
        initialZoom: 17,
        floorplanFile: 'galpao_setor_a.png',
        markers: [],
        createdAt: '2026-03-03'
      }
    ]
  },
  {
    id: 'mf_02',
    name: 'CLIENTE // CONDOMINIO VISTA VERDE',
    clientType: 'final_client',
    isExpanded: false,
    maps: [
      {
        id: 'map_03',
        name: 'Perimetro do Condominio & Guaritas',
        type: 'MAP_OPENSOURCE',
        companyScope: 'CLIENTE // CONDOMINIO VISTA VERDE',
        is_locked: false,
        initialCenter: [-23.58000, -46.680000],
        initialZoom: 15,
        markers: [
          { id: 'm_05', deviceId: 'cam_04', name: 'Perimetro dos Fundos', type: 'CAMERA', lat: -23.58050, lng: -46.680800, angle: 270, status: 'ONLINE' }
        ],
        createdAt: '2026-03-04'
      }
    ]
  }
]
