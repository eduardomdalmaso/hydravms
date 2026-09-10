export interface RtspPreset {
  brand: string
  label: string
  path: string
  queryAuth?: boolean
  pattern?: RegExp
}

export const RTSP_PRESETS: RtspPreset[] = [
  { brand: 'Intelbras', label: '[INTELBRAS] /cam/realmonitor?channel=1&subtype=0', path: '/cam/realmonitor?channel=1&subtype=0', pattern: /realmonitor/i },
  { brand: 'Intelbras', label: '[INTELBRAS] Query Auth (.sdp)', path: '/user={USER}&password={PASS}&channel=1&stream=0.sdp', queryAuth: true, pattern: /channel=\d+&stream=0\.sdp/i },
  { brand: 'Hikvision', label: '[HIKVISION] /Streaming/Channels/101', path: '/Streaming/Channels/101', pattern: /Streaming\/Channels/i },
  { brand: 'Hikvision', label: '[HIKVISION / JFL] /h264/ch1/main/av_stream', path: '/h264/ch1/main/av_stream', pattern: /h264\/ch\d+/i },
  { brand: 'Tecvoz', label: '[TECVOZ] Linha TW IP (/profile1)', path: '/profile1', pattern: /profile1/i },
  { brand: 'Tecvoz', label: '[TECVOZ] DVR Linha TW', path: '/chID=1&streamType=main&linkType=tcpa', pattern: /streamType=main/i },
  { brand: 'Axis', label: '[AXIS] /axis-media/media.amp', path: '/axis-media/media.amp?videocodec=h264', pattern: /axis-media/i },
  { brand: 'Vivotek', label: '[VIVOTEK] /live.sdp', path: '/live.sdp', pattern: /live\.sdp/i },
  { brand: 'Reolink', label: '[REOLINK] /Preview_01_main', path: '/Preview_01_main', pattern: /Preview_\d+/i },
  { brand: 'Tapo', label: '[TAPO / TP-LINK] /stream1', path: '/stream1', pattern: /stream[12]/i },
  { brand: 'Multilaser', label: '[MULTILASER] /H264?ch=1&subtype=0', path: '/H264?ch=1&subtype=0', pattern: /H264\?ch=/i },
  { brand: 'Foscam', label: '[FOSCAM] /videoMain', path: '/videoMain', pattern: /videoMain/i },
  { brand: 'Generic', label: '[GENERIC] /live (Padrão)', path: '/live', pattern: /^\/live$/i }
]

export function detectPreset(input: string): RtspPreset | undefined {
  if (!input) return undefined
  return RTSP_PRESETS.find(p => p.pattern && p.pattern.test(input))
}

export function generateChannelPaths(basePath: string, ch: number): { main: string; sub: string } {
  const num = ch || 1
  if (/Streaming\/Channels/i.test(basePath)) return { main: `/Streaming/Channels/${num}01`, sub: `/Streaming/Channels/${num}02` }
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
  const auth = user ? (pass ? `${encodeURIComponent(user)}:${encodeURIComponent(pass)}@` : `${encodeURIComponent(user)}@`) : ''
  return `rtsp://${auth}${host}:${p}${cleanPath}`
}
