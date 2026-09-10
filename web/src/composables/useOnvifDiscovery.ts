import { ref, computed } from 'vue'
import type { DiscoveredOnvifCamera } from '../types/onvifDiscovery'
import type { StreamItem } from '../types/streamTree'

const discoveredDevices = ref<DiscoveredOnvifCamera[]>([])
const lastScanTime = ref<string>('Nunca')

export function useOnvifDiscovery() {
  const isScanning = ref(false)
  const searchFilter = ref('')

  const filteredDevices = computed(() => {
    const q = searchFilter.value.toLowerCase()
    return discoveredDevices.value.filter(d => !d.isImported && (!q || d.name.toLowerCase().includes(q) || d.ip.includes(q) || d.manufacturer.toLowerCase().includes(q)))
  })

  const availableCount = computed(() => discoveredDevices.value.filter(d => !d.isImported).length)

  const scanNetwork = () => {
    isScanning.value = true
    setTimeout(() => {
      isScanning.value = false
      lastScanTime.value = new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    }, 1200)
  }

  const convertToStreamItem = (camera: DiscoveredOnvifCamera): Partial<StreamItem> => {
    const mainProf = camera.profiles[0]
    return {
      name: camera.name, protocol: 'ONVIF', ip: camera.ip, port: camera.port,
      url: mainProf ? mainProf.rtspUri : `rtsp://${camera.ip}:554/live`,
      codec: (mainProf?.codec as 'H.265' | 'H.264') || 'H.265',
      resolution: mainProf ? mainProf.resolution : '1080P',
      fps: 30, bitrate: '4.0 Mbps', recordMode: 'continuous', status: 'online', has_ptz: camera.hasPtz
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
