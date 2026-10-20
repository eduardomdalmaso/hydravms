export interface CameraBrandPreset {
  id: string
  label: string
  defaultPort: number
  defaultUser: string
  mainTemplate: string
  subTemplate: string
}

export const CAMERA_BRAND_PRESETS: CameraBrandPreset[] = [
  {
    id: 'intelbras',
    label: 'Intelbras',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/cam/realmonitor?channel={channel}&subtype=0',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/cam/realmonitor?channel={channel}&subtype=1',
  },
  {
    id: 'hikvision',
    label: 'Hikvision / HiLook',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/Streaming/Channels/{channel}01',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/Streaming/Channels/{channel}02',
  },
  {
    id: 'dahua',
    label: 'Dahua / Imou',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/cam/realmonitor?channel={channel}&subtype=0',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/cam/realmonitor?channel={channel}&subtype=1',
  },
  {
    id: 'tapo',
    label: 'TP-Link Tapo / VIGI',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/stream1',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/stream2',
  },
  {
    id: 'reolink',
    label: 'Reolink',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/h264Preview_{channel}_main',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/h264Preview_{channel}_sub',
  },
  {
    id: 'axis',
    label: 'Axis Communications',
    defaultPort: 554,
    defaultUser: 'root',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/axis-media/media.amp',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/axis-media/media.amp?resolution=640x360',
  },
  {
    id: 'uniview',
    label: 'Uniview (UNV)',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/unicast/c{channel}/s0/live',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/unicast/c{channel}/s1/live',
  },
  {
    id: 'yoosee',
    label: 'Yoosee / ONVIF Generico',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/onvif1',
    subTemplate: 'rtsp://{user}:{pass}@{ip}:{port}/onvif2',
  },
  {
    id: 'custom',
    label: 'Manual / Personalizado',
    defaultPort: 554,
    defaultUser: 'admin',
    mainTemplate: 'rtsp://{ip}:{port}/live',
    subTemplate: '',
  },
]

export function buildPresetRTSP(
  presetId: string,
  ip: string,
  port: number,
  user: string,
  pass: string,
  channel = 1
): { main: string; sub: string } {
  const preset = CAMERA_BRAND_PRESETS.find((p) => p.id === presetId) || CAMERA_BRAND_PRESETS[0]
  if (preset.id === 'custom') {
    return { main: '', sub: '' }
  }

  const cleanIP = ip.trim() || '192.168.1.100'
  const cleanPort = port > 0 ? port : 554
  const auth = user ? (pass ? `${user}:${pass}@` : `${user}@`) : ''

  const formatTemplate = (template: string) => {
    if (!template) return ''
    return template
      .replace('{user}:{pass}@', auth)
      .replace('{ip}', cleanIP)
      .replace('{port}', String(cleanPort))
      .replace('{channel}', String(channel))
  }

  return {
    main: formatTemplate(preset.mainTemplate),
    sub: formatTemplate(preset.subTemplate),
  }
}
