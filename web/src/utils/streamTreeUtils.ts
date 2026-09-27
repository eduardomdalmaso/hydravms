// HydraVMS - Stream Tree Mapping Utilities (< 100 lines)
import type { FolderNode, StreamItem } from '../types/streamTree'
import type { RegisteredCamera } from '../types/admin'
import type { ApiFolder } from '../services/api'

export function buildStreamTree(dbFolders: ApiFolder[], dbCameras: RegisteredCamera[]): { folders: FolderNode[]; rootStreams: StreamItem[] } {
  const folderMap = new Map<string, FolderNode>()
  dbFolders.forEach(f => {
    folderMap.set(f.id, { id: f.id, name: f.name, isExpanded: true, streams: [] })
  })

  const rootStreams: StreamItem[] = []
  dbCameras.forEach(c => {
    const subUrl = c.sub_stream_url || (c.rtsp_url && c.rtsp_url.includes('/stream1') ? c.rtsp_url.replace('/stream1', '/stream2') : '')
    const item: StreamItem = {
      id: c.id,
      name: c.name,
      protocol: (c.protocol?.toUpperCase() as any) || 'RTSP',
      url: c.rtsp_url || 'rtsp://',
      subUrl: subUrl || undefined,
      ip: '127.0.0.1',
      port: 554,
      codec: c.codec || 'H.265',
      subCodec: c.sub_codec || 'H.264',
      resolution: c.resolution || '1080P',
      subResolution: c.sub_resolution || '640x360 (SD)',
      fps: c.fps || 30,
      subFps: c.sub_fps || 15,
      bitrate: `${((c.bitrate_kbps || 2048) / 1024).toFixed(1)} Mbps`,
      subBitrate: `${((c.sub_bitrate_kbps || 512) / 1024).toFixed(1)} Mbps`,
      recordMode: 'disabled',
      status: (c.status as any) || 'online',
      has_ptz: c.has_ptz,
      snapshotUrl: (c as any).snapshot_url || `http://localhost:8083/api/v1/cameras/${c.id}/snapshot`,
      analyticsCount: 0,
      eventsCount: 0
    }
    if (c.folder_id && folderMap.has(c.folder_id)) {
      folderMap.get(c.folder_id)!.streams.push(item)
    } else {
      rootStreams.push(item)
    }
  })

  return { folders: Array.from(folderMap.values()), rootStreams }
}

export function createNewStreamItem(stream: Partial<StreamItem>, index: number): StreamItem {
  const id = stream.id || `cam_${index.toString().padStart(2, '0')}`
  const url = stream.url || 'rtsp://'
  const subUrl = stream.subUrl || (url.includes('/stream1') ? url.replace('/stream1', '/stream2') : '')
  return {
    id,
    name: stream.name || 'Nova Camera',
    protocol: stream.protocol || 'RTSP',
    url,
    subUrl: subUrl || undefined,
    ip: stream.ip || '192.168.1.100',
    port: stream.port || 554,
    codec: stream.codec || 'H.265',
    subCodec: stream.subCodec || 'H.264',
    resolution: stream.resolution || '1080P',
    subResolution: stream.subResolution || '640x360 (SD)',
    fps: stream.fps || 30,
    subFps: stream.subFps || 15,
    bitrate: stream.bitrate || '4.0 Mbps',
    subBitrate: stream.subBitrate || '512 Kbps',
    recordMode: stream.recordMode || 'disabled',
    status: stream.status || 'online',
    has_ptz: !!stream.has_ptz,
    snapshotUrl: stream.snapshotUrl || `http://localhost:8083/api/v1/cameras/${id}/snapshot`,
    analyticsCount: stream.analyticsCount || 0,
    eventsCount: stream.eventsCount || 0
  }
}

