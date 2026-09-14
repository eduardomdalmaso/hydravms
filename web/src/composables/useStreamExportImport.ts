import type { FolderNode, StreamItem } from '../types/streamTree'

export function useStreamExportImport() {
  const exportStreamsToCsv = (folders: FolderNode[], rootStreams: StreamItem[]) => {
    const rows: string[] = ['Nome,Protocolo,URL Mascarada (Segura),Codec,Resolucao,FPS,Pasta']
    const appendRow = (s: StreamItem, folder = 'Raiz') => {
      const safeId = s.id || s.name.toLowerCase().replace(/[^a-z0-9_]/g, '_')
      const masked = `rtsp://${window.location.hostname || 'localhost'}:8554/live/${safeId}`
      rows.push(`"${s.name}","${s.protocol}","${masked}","${s.codec}","${s.resolution}","${s.fps}","${folder}"`)
    }
    rootStreams.forEach(s => appendRow(s))
    folders.forEach(f => f.streams.forEach(s => appendRow(s, f.name)))

    const blob = new Blob(['\uFEFF' + rows.join('\r\n')], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob), a = document.createElement('a')
    a.href = url; a.download = `hydravms_cameras_${new Date().toISOString().slice(0, 10)}.csv`
    document.body.appendChild(a); a.click(); document.body.removeChild(a); URL.revokeObjectURL(url)
    return rows.length - 1
  }

  const parseStreamsFromFile = async (file: File): Promise<Partial<StreamItem>[]> => {
    const text = await file.text()
    if (file.name.endsWith('.json')) {
      const data = JSON.parse(text)
      const rawList = Array.isArray(data) ? data : (data.streams || [])
      return rawList.map((item: any, idx: number) => ({
        id: item.id || `imp_${Date.now()}_${idx}`, name: item.name || `Camera ${idx + 1}`,
        protocol: item.protocol || 'RTSP', url: item.url || item.maskedUrl || 'rtsp://192.168.1.100:554/live',
        ip: item.ip || '192.168.1.100', port: item.port || 554, codec: item.codec || 'H.264',
        resolution: item.resolution || '1080P', fps: item.fps || 30, bitrate: '4.0 Mbps',
        recordMode: 'disabled', status: 'online' as const, has_ptz: !!item.has_ptz
      }))
    }
    // Parse CSV or raw TXT line-by-line
    const lines = text.split(/\r?\n/).map(l => l.trim()).filter(l => l.length > 0)
    const result: Partial<StreamItem>[] = []
    const isHeader = (line: string) => /nome|name|url|ip|protocol|camer/i.test(line)
    const startIdx = (lines.length > 0 && isHeader(lines[0])) ? 1 : 0

    for (let i = startIdx; i < lines.length; i++) {
      const line = lines[i]
      const cols = line.includes(';') ? line.split(';').map(c => c.replace(/^"|"$/g, '').trim()) : line.split(',').map(c => c.replace(/^"|"$/g, '').trim())
      if (cols.length === 1 && (cols[0].startsWith('rtsp://') || cols[0].startsWith('rtmp://') || cols[0].endsWith('.mp4'))) {
        result.push({ name: `Camera ${i + 1}`, protocol: cols[0].endsWith('.mp4') ? 'LOOP' : (cols[0].startsWith('rtmp://') ? 'RTMP' : 'RTSP'), url: cols[0], ip: cols[0].split('@').pop()?.split('/')[0]?.split(':')[0] || '192.168.1.100', port: 554, codec: 'H.264', resolution: '1080P', fps: 30, bitrate: '4.0 Mbps', recordMode: 'disabled', status: 'online' })
      } else if (cols.length >= 2) {
        const name = cols[0] || `Camera ${i + 1}`, urlOrIp = cols[1] || 'rtsp://192.168.1.100:554/live'
        const proto = (cols[2] as any) || (urlOrIp.endsWith('.mp4') ? 'LOOP' : (urlOrIp.startsWith('rtmp://') ? 'RTMP' : 'RTSP'))
        result.push({ name, protocol: proto, url: urlOrIp.includes('://') || urlOrIp.endsWith('.mp4') ? urlOrIp : `rtsp://${urlOrIp}:554/live`, ip: urlOrIp.replace(/rtsp:\/\/|rtmp:\/\//g, '').split(':')[0] || '192.168.1.100', port: parseInt(cols[3]) || 554, codec: (cols[4] as any) || 'H.264', resolution: cols[5] || '1080P', fps: 30, bitrate: '4.0 Mbps', recordMode: 'disabled', status: 'online' })
      }
    }
    return result
  }

  return { exportStreamsToCsv, parseStreamsFromFile }
}
