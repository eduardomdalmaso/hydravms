import type { FolderNode, StreamItem } from '../types/streamTree'

export const initialRootStreams: StreamItem[] = [
  { id: 'cam_root_01', name: 'Almoxarifado Geral', protocol: 'RTSP', url: 'rtsp://192.168.1.109:554/live', ip: '192.168.1.109', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, bitrate: '4.0 Mbps', recordMode: 'continuous', status: 'online', has_ptz: false },
  { id: 'cam_root_02', name: 'Refeitorio Central', protocol: 'RTSP', url: 'rtsp://192.168.1.110:554/live', ip: '192.168.1.110', port: 554, codec: 'H.264', resolution: '1080P', fps: 25, bitrate: '3.5 Mbps', recordMode: 'motion', status: 'online', has_ptz: false }
]

export const initialFolders: FolderNode[] = [
  {
    id: 'f_01',
    name: 'PORTARIA & ACESSOS PRINCIPAIS',
    isExpanded: true,
    streams: [
      { id: 'cam_01', name: 'Portaria Principal (Entrada)', protocol: 'RTSP', url: 'rtsp://192.168.1.101:554/live', ip: '192.168.1.101', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, bitrate: '4.2 Mbps', recordMode: 'continuous', status: 'online', has_ptz: true },
      { id: 'cam_02', name: 'Estacionamento Visitantes', protocol: 'RTSP', url: 'rtsp://192.168.1.102:554/live', ip: '192.168.1.102', port: 554, codec: 'H.264', resolution: '1080P', fps: 25, bitrate: '3.8 Mbps', recordMode: 'motion', status: 'online', has_ptz: false }
    ]
  },
  {
    id: 'f_02',
    name: 'GALPAO DE CARGAS & DOCAS',
    isExpanded: true,
    streams: [
      { id: 'cam_03', name: 'Corredor de Cargas & Docas', protocol: 'RTSP', url: 'rtsp://192.168.1.103:554/live', ip: '192.168.1.103', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, bitrate: '4.5 Mbps', recordMode: 'continuous', status: 'recording', has_ptz: false },
      { id: 'cam_04', name: 'Perimetro dos Fundos', protocol: 'RTSP', url: 'rtsp://192.168.1.104:554/live', ip: '192.168.1.104', port: 554, codec: 'H.265', resolution: '4K', fps: 30, bitrate: '8.1 Mbps', recordMode: 'ai_event', status: 'recording', has_ptz: true }
    ]
  },
  {
    id: 'f_03',
    name: 'DATACENTER & NOC',
    isExpanded: false,
    streams: [
      { id: 'cam_05', name: 'Racks Servidores Sala 01', protocol: 'RTSP', url: 'rtsp://192.168.2.101:554/live', ip: '192.168.2.101', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, bitrate: '3.9 Mbps', recordMode: 'continuous', status: 'online', has_ptz: false }
    ]
  }
]
