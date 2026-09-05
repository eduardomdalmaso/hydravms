import type { ServerNode } from '../components/admin/TreeServerCard.vue'

export const initialTopology: ServerNode[] = [
  {
    id: 'edge_01',
    name: 'EDGE-SERVER-01',
    ip: '10.0.4.10',
    location: 'Galpao Principal',
    switches: [
      {
        id: 'sw_01',
        name: 'SWITCH-POE-A [ACESSOS]',
        ip: '192.168.1.20',
        cameras: [
          { id: 'cam_01', name: 'Portaria Principal (Entrada)', ip: '192.168.1.101', codec: 'H.265', res: '1080P', fps: 30, bitrate: '4.2 Mbps', status: 'online' },
          { id: 'cam_02', name: 'Estacionamento Visitantes', ip: '192.168.1.102', codec: 'H.264', res: '1080P', fps: 25, bitrate: '3.8 Mbps', status: 'online' }
        ]
      },
      {
        id: 'sw_02',
        name: 'SWITCH-POE-B [GALPAO]',
        ip: '192.168.1.30',
        cameras: [
          { id: 'cam_03', name: 'Corredor de Cargas & Docas', ip: '192.168.1.103', codec: 'H.265', res: '1080P', fps: 30, bitrate: '4.5 Mbps', status: 'recording' },
          { id: 'cam_04', name: 'Perimetro dos Fundos', ip: '192.168.1.104', codec: 'H.265', res: '4K', fps: 30, bitrate: '8.1 Mbps', status: 'recording' }
        ]
      }
    ]
  },
  {
    id: 'edge_02',
    name: 'EDGE-SERVER-02',
    ip: '10.0.4.11',
    location: 'CPD Central',
    switches: [
      {
        id: 'sw_03',
        name: 'SWITCH-POE-C [CPD / DATACENTER]',
        ip: '192.168.2.10',
        cameras: [
          { id: 'cam_05', name: 'Racks Servidores Sala 01', ip: '192.168.2.101', codec: 'H.265', res: '1080P', fps: 30, bitrate: '3.9 Mbps', status: 'online' },
          { id: 'cam_06', name: 'NOC / Sala de Controle', ip: '192.168.2.102', codec: 'H.265', res: '1080P', fps: 30, bitrate: '3.5 Mbps', status: 'online' }
        ]
      }
    ]
  }
]
