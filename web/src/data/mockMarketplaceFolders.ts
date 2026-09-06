import type { AnalyticFolderNode } from '../types/marketplace'

export const mockMarketplaceFolders: AnalyticFolderNode[] = [
  {
    id: 'fld-lpr-01',
    name: 'Portaria & Cancelas',
    plugin_id: 'hydra-analytics-lpr',
    instances: [{
      id: 'inst-lpr-01', plugin_id: 'hydra-analytics-lpr', plugin_name: 'LPR / ALPR Mercosul & Antigas',
      name: 'LPR // Portaria Principal - Pista 1', camera_id: 'cam_portaria_01', camera_name: 'CAM_01 // Portaria Principal',
      stream_type: 'main_1080p', hardware_target: 'rtx_5090_cuda', confidence_threshold: 0.70, roi_mode: 'custom_polygon',
      specific_params: { region: 'mercosul', sahi_enabled: true }, is_active: true, fps_rate: 29.8, detections_count: 86, created_at: '2026-09-05T08:30:00Z'
    }]
  },
  {
    id: 'fld-lpr-02',
    name: 'Garagem & Subsolo',
    plugin_id: 'hydra-analytics-lpr',
    instances: [{
      id: 'inst-lpr-02', plugin_id: 'hydra-analytics-lpr', plugin_name: 'LPR / ALPR Mercosul & Antigas',
      name: 'LPR // Entrada Garagem Subsolo', camera_id: 'cam_garagem_02', camera_name: 'CAM_02 // Garagem Subsolo',
      stream_type: 'main_1080p', hardware_target: 'rtx_5090_cuda', confidence_threshold: 0.65, roi_mode: 'full_frame',
      specific_params: { region: 'mercosul', sahi_enabled: false }, is_active: true, fps_rate: 30.0, detections_count: 62, created_at: '2026-09-05T09:15:00Z'
    }]
  },
  {
    id: 'fld-face-01',
    name: 'Controle de Catracas',
    plugin_id: 'hydra-analytics-face',
    instances: [{
      id: 'inst-face-01', plugin_id: 'hydra-analytics-face', plugin_name: 'Reconhecimento Facial & Acesso VIP',
      name: 'FACIAL // Catracas Hall de Entrada', camera_id: 'cam_hall_03', camera_name: 'CAM_03 // Catracas Recepção',
      stream_type: 'main_1080p', hardware_target: 'rtx_5090_cuda', confidence_threshold: 0.80, roi_mode: 'custom_polygon',
      specific_params: { min_face_size: 64, anti_spoof: true }, is_active: true, fps_rate: 28.5, detections_count: 94, created_at: '2026-09-05T10:00:00Z'
    }]
  },
  {
    id: 'fld-ppe-01',
    name: 'Galpões Logísticos',
    plugin_id: 'hydra-analytics-ppe',
    instances: [{
      id: 'inst-ppe-01', plugin_id: 'hydra-analytics-ppe', plugin_name: 'EPI / Segurança do Trabalho',
      name: 'EPI // Área de Carga e Galpão B', camera_id: 'cam_galpao_04', camera_name: 'CAM_04 // Galpão Logística B',
      stream_type: 'main_1080p', hardware_target: 'rtx_5090_cuda', confidence_threshold: 0.65, roi_mode: 'full_frame',
      specific_params: { require_helmet: true, require_vest: true }, is_active: true, fps_rate: 30.0, detections_count: 42, created_at: '2026-09-05T11:20:00Z'
    }]
  }
]
