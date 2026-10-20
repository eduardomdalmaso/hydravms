export interface RtspPreset {
  brand: string
  brandName: string
  label: string
  path: string
  queryAuth?: boolean
  pattern?: RegExp
}

export const RTSP_PRESETS: RtspPreset[] = [
  { brand: 'Intelbras', brandName: 'Intelbras', label: '[INTELBRAS] /cam/realmonitor?channel=1&subtype=0', path: '/cam/realmonitor?channel=1&subtype=0', pattern: /realmonitor/i },
  { brand: 'Intelbras', brandName: 'Intelbras (Query Auth)', label: '[INTELBRAS] Query Auth (.sdp)', path: '/user={USER}&password={PASS}&channel=1&stream=0.sdp', queryAuth: true, pattern: /channel=\d+&stream=0\.sdp/i },
  { brand: 'Hikvision', brandName: 'Hikvision', label: '[HIKVISION] /Streaming/Channels/101', path: '/Streaming/Channels/101', pattern: /Streaming\/Channels/i },
  { brand: 'Hikvision', brandName: 'Hikvision / JFL', label: '[HIKVISION / JFL] /h264/ch1/main/av_stream', path: '/h264/ch1/main/av_stream', pattern: /h264\/ch\d+/i },
  { brand: 'Tecvoz', brandName: 'Tecvoz (Linha IP)', label: '[TECVOZ] Linha TW IP (/profile1)', path: '/profile1', pattern: /profile1/i },
  { brand: 'Tecvoz', brandName: 'Tecvoz (Linha DVR)', label: '[TECVOZ] DVR Linha TW', path: '/chID=1&streamType=main&linkType=tcpa', pattern: /streamType=main/i },
  { brand: 'Axis', brandName: 'Axis', label: '[AXIS] /axis-media/media.amp', path: '/axis-media/media.amp?videocodec=h264', pattern: /axis-media/i },
  { brand: 'Vivotek', brandName: 'Vivotek', label: '[VIVOTEK] /live.sdp', path: '/live.sdp', pattern: /live\.sdp/i },
  { brand: 'Reolink', brandName: 'Reolink', label: '[REOLINK] /Preview_01_main', path: '/Preview_01_main', pattern: /Preview_\d+/i },
  { brand: 'Tapo', brandName: 'TP-Link / Tapo', label: '[TAPO / TP-LINK] /stream1', path: '/stream1', pattern: /stream[12]/i },
  { brand: 'Multilaser', brandName: 'Multilaser', label: '[MULTILASER] /H264?ch=1&subtype=0', path: '/H264?ch=1&subtype=0', pattern: /H264\?ch=/i },
  { brand: 'Foscam', brandName: 'Foscam', label: '[FOSCAM] /videoMain', path: '/videoMain', pattern: /videoMain/i },
  { brand: 'Generic', brandName: 'Genérico / Padrão', label: '[GENERIC] /live (Padrão)', path: '/live', pattern: /^\/live$/i }
]

export function detectPreset(input: string): RtspPreset | undefined {
  if (!input) return undefined
  return RTSP_PRESETS.find(p => p.pattern && p.pattern.test(input))
}

export function generateChannelPaths(basePath: string, ch: number): { main: string; sub: string } {
  const hikMatch = basePath.match(/Streaming\/[Cc]hannels\/(\d+)01/i)
  const baseNum = hikMatch ? parseInt(hikMatch[1]) : 0
  const num = baseNum > 0 ? (ch === 1 ? baseNum : (baseNum + ch - 1)) : (ch || 1)

  if (/Streaming\/[Cc]hannels/i.test(basePath)) return { main: `/Streaming/channels/${num}01`, sub: `/Streaming/channels/${num}02` }
  if (/realmonitor/i.test(basePath)) return { main: `/cam/realmonitor?channel=${num}&subtype=0`, sub: `/cam/realmonitor?channel=${num}&subtype=1` }
  if (/channel=\d+/i.test(basePath)) return { main: basePath.replace(/channel=\d+/, `channel=${num}`), sub: basePath.replace(/channel=\d+/, `channel=${num}`).replace(/stream=0/, 'stream=1') }
  if (/Preview_\d+/i.test(basePath)) { const pad = num.toString().padStart(2, '0'); return { main: `/Preview_${pad}_main`, sub: `/Preview_${pad}_sub` } }
  if (/stream\d+/i.test(basePath)) return { main: `/stream${num}`, sub: num === 1 ? '/stream2' : `/stream${num}_sub` }

  if (/H264\?ch=/i.test(basePath)) return { main: `/H264?ch=${num}&subtype=0`, sub: `/H264?ch=${num}&subtype=1` }
  return { main: `/ch${num}/main`, sub: `/ch${num}/sub` }
}

export function parseRtspUrl(raw: string): { user?: string; pass?: string; ip?: string; port?: number; path?: string } | null {
  try {
    const clean = raw.trim().replace(/^rtsp:\/\//i, '')
    const slashIdx = clean.indexOf('/')
    const authority = slashIdx !== -1 ? clean.slice(0, slashIdx) : clean
    const path = slashIdx !== -1 ? clean.slice(slashIdx) : '/live'
    let user = '', pass = '', hostPort = authority
    if (authority.includes('@')) {
      const [auth, hp] = authority.split('@')
      hostPort = hp
      if (auth.includes(':')) { [user, pass] = auth.split(':') } else { user = auth }
    }
    let ip = hostPort, port = 554
    if (hostPort.includes(':')) {
      const [h, p] = hostPort.split(':')
      ip = h; port = parseInt(p) || 554
    }
    return { user: decodeURIComponent(user), pass: decodeURIComponent(pass), ip, port, path }
  } catch { return null }
}

export function buildRtspUrl(preset: RtspPreset | undefined, ip: string, port: number, path: string, user: string, pass: string): string {
  const host = ip || '192.168.1.100', p = port || 554
  if (preset?.queryAuth) {
    const filled = preset.path.replace('{USER}', encodeURIComponent(user || 'admin')).replace('{PASS}', encodeURIComponent(pass || 'admin'))
    return `rtsp://${host}:${p}${filled.startsWith('/') ? filled : '/' + filled}`
  }
  const cleanPath = path ? (path.startsWith('/') ? path : '/' + path) : '/live'
  const auth = user ? (pass ? `${user}:${pass}@` : `${user}@`) : ''
  return `rtsp://${auth}${host}:${p}${cleanPath}`
}

export function deriveSubStreamUrl(mainUrl: string): string {
  if (!mainUrl) return ''
  // 1. Hikvision: /Streaming/Channels/X01 -> /Streaming/Channels/X02
  if (/Streaming\/[Cc]hannels\/\d+01/i.test(mainUrl)) {
    return mainUrl.replace(/(Streaming\/[Cc]hannels\/\d+)01/i, '$102')
  }
  // 2. Intelbras / Dahua: subtype=0 -> subtype=1
  if (/subtype=0/i.test(mainUrl)) {
    return mainUrl.replace(/subtype=0/i, 'subtype=1')
  }
  // 3. Tapo: /stream1 -> /stream2
  if (/\/stream1(\b|$)/i.test(mainUrl)) {
    return mainUrl.replace(/\/stream1(\b|$)/i, '/stream2$1')
  }
  // 4. Reolink: _main -> _sub
  if (/_main(\b|$)/i.test(mainUrl)) {
    return mainUrl.replace(/_main(\b|$)/i, '_sub$1')
  }
  // 5. Generic /main -> /sub
  if (/\/main(\b|$)/i.test(mainUrl)) {
    return mainUrl.replace(/\/main(\b|$)/i, '/sub$1')
  }
  // 6. Generic /profile1 -> /profile2
  if (/profile1/i.test(mainUrl)) {
    return mainUrl.replace(/profile1/i, 'profile2')
  }
  return `${mainUrl}_sub`
}

export function detectBrandFromUrl(urlOrPath: string): { brand: string; model?: string } {
  if (!urlOrPath) return { brand: '' }
  if (/Streaming\/[Cc]hannels/i.test(urlOrPath)) return { brand: 'Hikvision', model: 'IP Camera Series' }
  if (/realmonitor/i.test(urlOrPath)) return { brand: 'Intelbras', model: 'VIP Series' }
  if (/stream[12]/i.test(urlOrPath)) return { brand: 'TP-Link Tapo', model: 'Tapo Series' }
  if (/Preview_\d+/i.test(urlOrPath) || /h264Preview_/i.test(urlOrPath)) return { brand: 'Reolink', model: 'RLC Series' }
  if (/axis-media/i.test(urlOrPath)) return { brand: 'Axis Communications', model: 'M/P Series' }
  if (/profile1/i.test(urlOrPath)) return { brand: 'Tecvoz', model: 'TW Series' }
  return { brand: '' }
}

