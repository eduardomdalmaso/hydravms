import type { FolderNode, StreamItem } from '../types/streamTree'

export function useStreamExportImport() {
  const exportStreamsToJson = (folders: FolderNode[], rootStreams: StreamItem[]) => {
    const allStreams: Array<Partial<StreamItem> & { folder?: string; maskedUrl: string }> = []
    
    // Process Root Streams
    for (const s of rootStreams) {
      const safeId = s.id || s.name.toLowerCase().replace(/[^a-z0-9_]/g, '_')
      allStreams.push({
        name: s.name,
        protocol: s.protocol,
        maskedUrl: `rtsp://${window.location.hostname || 'localhost'}:8554/live/${safeId}`,
        codec: s.codec,
        resolution: s.resolution,
        fps: s.fps,
        bitrate: s.bitrate,
        recordMode: s.recordMode,
        has_ptz: s.has_ptz,
        latitude: s.latitude,
        longitude: s.longitude,
        locationName: s.locationName
      })
    }

    // Process Folder Streams
    for (const f of folders) {
      for (const s of f.streams) {
        const safeId = s.id || s.name.toLowerCase().replace(/[^a-z0-9_]/g, '_')
        allStreams.push({
          folder: f.name,
          name: s.name,
          protocol: s.protocol,
          maskedUrl: `rtsp://${window.location.hostname || 'localhost'}:8554/live/${safeId}`,
          codec: s.codec,
          resolution: s.resolution,
          fps: s.fps,
          bitrate: s.bitrate,
          recordMode: s.recordMode,
          has_ptz: s.has_ptz,
          latitude: s.latitude,
          longitude: s.longitude,
          locationName: s.locationName
        })
      }
    }

    const jsonStr = JSON.stringify({
      version: '1.0',
      exported_at: new Date().toISOString(),
      security_note: 'URLs físicas e credenciais mascaradas via HydraStream Reverse-Proxy',
      total_streams: allStreams.length,
      streams: allStreams
    }, null, 2)

    const blob = new Blob([jsonStr], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `hydravms_streams_export_${new Date().toISOString().slice(0, 10)}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    return allStreams.length
  }

  const parseStreamsFromJson = async (file: File): Promise<Partial<StreamItem>[]> => {
    const text = await file.text()
    const data = JSON.parse(text)
    const rawList = Array.isArray(data) ? data : (data.streams || [])
    return rawList.map((item: any, idx: number) => ({
      id: item.id || `imported_${Date.now()}_${idx}`,
      name: item.name || `Stream Importado ${idx + 1}`,
      protocol: item.protocol || 'RTSP',
      url: item.url || item.maskedUrl || 'rtsp://192.168.1.100:554/live',
      ip: item.ip || '192.168.1.100',
      port: item.port || 554,
      codec: item.codec || 'H.264',
      resolution: item.resolution || '1080P',
      fps: item.fps || 30,
      bitrate: item.bitrate || '4.0 Mbps',
      recordMode: item.recordMode || 'disabled',
      status: 'online' as const,
      has_ptz: !!item.has_ptz,
      latitude: item.latitude || -23.550520,
      longitude: item.longitude || -46.633308,
      locationName: item.locationName || ''
    }))
  }

  return {
    exportStreamsToJson,
    parseStreamsFromJson
  }
}
