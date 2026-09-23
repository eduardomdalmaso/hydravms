import { ref, computed } from 'vue'
import type { DiscoveredOnvifCamera } from '../types/onvifDiscovery'
import type { StreamItem } from '../types/streamTree'
import { discoverOnvifDevices } from '../services/adminApi'

const discoveredDevices = ref<DiscoveredOnvifCamera[]>([])
const lastScanTime = ref<string>('')

export function useOnvifDiscovery() {
  const isScanning = ref(false)
  const searchFilter = ref('')

  const filteredDevices = computed(() => {
    const q = searchFilter.value.toLowerCase()
    return discoveredDevices.value.filter(d => !d.isImported && (!q || d.name.toLowerCase().includes(q) || d.ip.includes(q) || d.manufacturer.toLowerCase().includes(q)))
  })

  const availableCount = computed(() => discoveredDevices.value.filter(d => !d.isImported).length)

  const scanNetwork = async () => {
    isScanning.value = true
    try {
      const rawDevices = await discoverOnvifDevices()
      if (Array.isArray(rawDevices)) {
        discoveredDevices.value = rawDevices.map((d: any) => ({
          id: d.id || d.device_id || `onvif_${d.ip_address || d.ip}_${d.port || 80}`,
          name: d.name || `Câmera ONVIF ${d.ip_address || d.ip}`,
          manufacturer: d.manufacturer || 'ONVIF',
          model: d.model || 'IP Camera',
          ip: d.ip || d.ip_address || '127.0.0.1',
          port: d.port || 80,
          macAddress: d.macAddress || d.hardware_id || d.serial_number || 'AA:BB:CC:DD:EE:FF',
          profiles: (Array.isArray(d.profiles) ? d.profiles : []).map((p: any) => ({
            name: p.name || 'MainProfile',
            token: p.token || 'profile_0',
            resolution: p.resolution || (p.width && p.height ? `${p.width}x${p.height}` : '1080P'),
            codec: (p.codec || p.encoding || 'H.265').toUpperCase().includes('264') ? 'H.264' : 'H.265',
            rtspUri: p.rtspUri || p.rtsp_stream || d.rtsp_url || `rtsp://${d.ip || d.ip_address || '127.0.0.1'}:554/live`
          })),
          hasPtz: !!(d.hasPtz ?? d.has_ptz),
          isImported: !!d.isImported
        }))
      }
      lastScanTime.value = new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    } finally {
      isScanning.value = false
    }
  }

  const convertToStreamItem = (camera: DiscoveredOnvifCamera): Partial<StreamItem> => {
    const profiles = Array.isArray(camera.profiles) ? camera.profiles : []
    const mainProf = profiles[0]
    return {
      name: camera.name || 'Câmera ONVIF', protocol: 'ONVIF', ip: camera.ip || '', port: camera.port || 80,
      url: mainProf?.rtspUri || `rtsp://${camera.ip || '127.0.0.1'}:554/live`,
      codec: (mainProf?.codec as 'H.265' | 'H.264') || 'H.265',
      resolution: mainProf ? mainProf.resolution : '1080P',
      fps: 30, bitrate: '4.0 Mbps', recordMode: 'continuous', status: 'online', has_ptz: !!camera.hasPtz
    }
  }

  const markAsImported = (idOrIp: string) => {
    const dev = discoveredDevices.value.find(d => d.id === idOrIp || d.ip === idOrIp)
    if (dev) dev.isImported = true
  }

  return {
    isScanning, discoveredDevices, filteredDevices, availableCount, lastScanTime,
    searchFilter, scanNetwork, convertToStreamItem, markAsImported
  }
}
