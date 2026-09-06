import type { AnalyticEventRecord } from '../types/marketplace'

export const mockMarketplaceEvents: AnalyticEventRecord[] = [
  {
    id: 'evt-101',
    timestamp: '2026-09-06T11:14:22Z',
    plugin_id: 'hydra-analytics-lpr',
    plugin_name: 'LPR / ALPR Mercosul & Antigas',
    camera_id: 'cam_portaria_01',
    camera_name: 'CAM_01 // Portaria Principal',
    event_type: 'vehicle_plate_read',
    severity: 'info',
    confidence: 0.96,
    object_label: 'BRA2E19',
    details: 'Veículo Toyota Corolla // Padrão Mercosul // Liberado',
    raw_payload: { plate: 'BRA2E19', type: 'car', speed_kmh: 18, country: 'BRA' }
  },
  {
    id: 'evt-102',
    timestamp: '2026-09-06T11:12:05Z',
    plugin_id: 'hydra-analytics-lpr',
    plugin_name: 'LPR / ALPR Mercosul & Antigas',
    camera_id: 'cam_garagem_02',
    camera_name: 'CAM_02 // Garagem Subsolo',
    event_type: 'vehicle_plate_read',
    severity: 'warning',
    confidence: 0.88,
    object_label: 'XYZ9876',
    details: 'Veículo Ford Cargo // Placa Cinza // Visitante Não Autorizado',
    raw_payload: { plate: 'XYZ9876', type: 'truck', speed_kmh: 12, country: 'BRA' }
  },
  {
    id: 'evt-103',
    timestamp: '2026-09-06T11:08:44Z',
    plugin_id: 'hydra-analytics-face',
    plugin_name: 'Reconhecimento Facial & Acesso VIP',
    camera_id: 'cam_hall_03',
    camera_name: 'CAM_03 // Catracas Recepção',
    event_type: 'face_identified',
    severity: 'info',
    confidence: 0.94,
    object_label: 'Carlos Silveira',
    details: 'Engenheiro Chefe // Acesso Catraca 02 Liberado',
    raw_payload: { person_id: 'USR-891', role: 'engineer', liveness_score: 0.99 }
  },
  {
    id: 'evt-104',
    timestamp: '2026-09-06T11:01:17Z',
    plugin_id: 'hydra-analytics-face',
    plugin_name: 'Reconhecimento Facial & Acesso VIP',
    camera_id: 'cam_hall_03',
    camera_name: 'CAM_03 // Catracas Recepção',
    event_type: 'face_unknown',
    severity: 'warning',
    confidence: 0.72,
    object_label: 'Desconhecido',
    details: 'Rosto não cadastrado na base corporativa // Direcionado à Recepção',
    raw_payload: { person_id: null, liveness_score: 0.92 }
  },
  {
    id: 'evt-105',
    timestamp: '2026-09-06T10:45:30Z',
    plugin_id: 'hydra-analytics-ppe',
    plugin_name: 'EPI / Segurança do Trabalho',
    camera_id: 'cam_galpao_04',
    camera_name: 'CAM_04 // Galpão Logística B',
    event_type: 'ppe_violation',
    severity: 'critical',
    confidence: 0.89,
    object_label: 'Sem Capacete',
    details: 'Operador em zona de empilhadeira sem capacete regulamentar',
    raw_payload: { missing: ['helmet'], has_vest: true, zone: 'DOCK_04' }
  },
  {
    id: 'evt-106',
    timestamp: '2026-09-06T10:30:12Z',
    plugin_id: 'hydra-analytics-ppe',
    plugin_name: 'EPI / Segurança do Trabalho',
    camera_id: 'cam_galpao_04',
    camera_name: 'CAM_04 // Galpão Logística B',
    event_type: 'ppe_compliant',
    severity: 'info',
    confidence: 0.92,
    object_label: 'EPI Conforme',
    details: 'Operador com capacete e colete reflexivo detectados',
    raw_payload: { missing: [], has_helmet: true, has_vest: true }
  }
]
