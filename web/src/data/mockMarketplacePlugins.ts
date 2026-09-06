import type { PluginManifest } from '../types/marketplace'

export const mockMarketplacePlugins: PluginManifest[] = [
  {
    id: 'hydra-analytics-lpr',
    name: 'LPR / ALPR Mercosul & Antigas',
    version: '1.4.2',
    author: 'Hydra Vision Lab',
    category: 'traffic',
    description: 'Leitura de placas veiculares Mercosul e padrão cinza com SAHI e YOLO11.',
    runtime: 'python3',
    hardware_req: 'CUDA 13.3 // RTX 5090',
    permissions: ['stream:read', 'events:publish', 'shm:read', 'gpu:cuda'],
    is_official: true, is_installed: true, status: 'running', instances_count: 2, events_count: 148,
    default_config: { region: 'mercosul', sahi_enabled: true, confidence: 0.65 }
  },
  {
    id: 'hydra-analytics-face',
    name: 'Reconhecimento Facial & Acesso VIP',
    version: '2.1.0',
    author: 'Hydra Vision Lab',
    category: 'access_control',
    description: 'Identificação facial em tempo real com anti-spoofing e verificação de watchlist.',
    runtime: 'python3',
    hardware_req: 'CUDA 13.3 // RTX 5090',
    permissions: ['stream:read', 'events:publish', 'shm:read', 'gpu:cuda'],
    is_official: true, is_installed: true, status: 'running', instances_count: 1, events_count: 94,
    default_config: { min_face_size: 48, anti_spoof: true, confidence: 0.75 }
  },
  {
    id: 'hydra-analytics-ppe',
    name: 'EPI / Segurança do Trabalho',
    version: '1.2.5',
    author: 'Hydra Safety AI',
    category: 'safety',
    description: 'Detecção contínua de capacete, colete reflexivo e máscara de proteção.',
    runtime: 'python3',
    hardware_req: 'CUDA 13.3 // RTX 5090',
    permissions: ['stream:read', 'events:publish', 'shm:read', 'gpu:cuda'],
    is_official: true, is_installed: true, status: 'running', instances_count: 1, events_count: 42,
    default_config: { require_helmet: true, require_vest: true, confidence: 0.60 }
  },
  {
    id: 'hydra-analytics-perimeter',
    name: 'Invasão Perimetral & Cerca Virtual',
    version: '3.0.1',
    author: 'Hydra Security Corp',
    category: 'analytics',
    description: 'Cruzamento de linha, intrusão em zona proibida e loitering com tracking persistente.',
    runtime: 'binary_elf',
    hardware_req: 'CPU // ZERO-COPY SHM',
    permissions: ['stream:read', 'events:publish', 'shm:read'],
    is_official: true, is_installed: false, status: 'available', instances_count: 0, events_count: 0,
    default_config: { dwell_seconds: 5, direction: 'both', confidence: 0.70 }
  },
  {
    id: 'hydra-analytics-fire',
    name: 'Detecção de Fogo & Fumaça',
    version: '1.0.8',
    author: 'Pyros Threat Shield',
    category: 'safety',
    description: 'Identificação precoce de chamas ativas e colunas de fumaça em galpões e pátios.',
    runtime: 'python3',
    hardware_req: 'CUDA 13.3 // RTX 5090',
    permissions: ['stream:read', 'events:publish', 'shm:read', 'gpu:cuda'],
    is_official: true, is_installed: false, status: 'available', instances_count: 0, events_count: 0,
    default_config: { smoke_threshold: 0.55, fire_threshold: 0.65 }
  },
  {
    id: 'hydra-analytics-flow',
    name: 'Contador de Fluxo de Pessoas & Veículos',
    version: '2.4.0',
    author: 'Hydra Traffic Mesh',
    category: 'traffic',
    description: 'Contagem bidirecional em linhas de passagem para métricas de ocupação e tráfego.',
    runtime: 'binary_elf',
    hardware_req: 'CPU // ZERO-COPY SHM',
    permissions: ['stream:read', 'events:publish', 'shm:read'],
    is_official: true, is_installed: false, status: 'available', instances_count: 0, events_count: 0,
    default_config: { direction_filter: 'both', aggregation_minutes: 15 }
  }
]
