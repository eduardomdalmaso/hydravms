import { ref, computed } from 'vue'
import type { DiscoveredOnvifCamera } from '../types/onvifDiscovery'
import type { StreamItem } from '../types/streamTree'

const initialDiscovered: DiscoveredOnvifCamera[] = [
  {
    id: 'onvif_discov_01', name: 'Intelbras VIP 3230 B', manufacturer: 'Intelbras',
    model: 'VIP 3230 B VF', ip: '192.168.1.145', port: 80, macAddress: 'A0:B1:C2:33:44:55',
    profiles: [
      { name: 'MainStream_1080P', token: 'profile_1', resolution: '1920x1080', codec: 'H.265', rtspUri: 'rtsp://admin:admin@192.168.1.145:554/cam/realmonitor?channel=1&subtype=0' },
      { name: 'SubStream_480P', token: 'profile_2', resolution: '640x480', codec: 'H.264', rtspUri: 'rtsp://admin:admin@192.168.1.145:554/cam/realmonitor?channel=1&subtype=1' }
    ],
    hasPtz: false, isImported: false
  },
  {
    id: 'onvif_discov_02', name: 'Hikvision DS-2CD2043', manufacturer: 'Hikvision',
    model: 'DS-2CD2043G2-I', ip: '192.168.1.188', port: 80, macAddress: '54:B8:0A:11:22:33',
    profiles: [
      { name: 'HD_4MP_Main', token: 'Profile_1', resolution: '2560x1440', codec: 'H.265', rtspUri: 'rtsp://admin:pass@192.168.1.188:554/Streaming/Channels/101' },
      { name: 'SD_SubStream', token: 'Profile_2', resolution: '640x360', codec: 'H.264', rtspUri: 'rtsp://admin:pass@192.168.1.188:554/Streaming/Channels/102' }
    ],
    hasPtz: false, isImported: false
  },
  {
    id: 'onvif_discov_03', name: 'Dahua Speed Dome PTZ', manufacturer: 'Dahua',
    model: 'SD49225XA-HNR', ip: '192.168.1.210', port: 80, macAddress: '3C:EF:8C:77:88:99',
    profiles: [
      { name: 'MainStream_4K', token: 'Main_01', resolution: '3840x2160', codec: 'H.265', rtspUri: 'rtsp://admin:admin@192.168.1.210:554/cam/realmonitor?channel=1&subtype=0' },
      { name: 'SubStream_VGA', token: 'Sub_01', resolution: '720x480', codec: 'H.264', rtspUri: 'rtsp://admin:admin@192.168.1.210:554/cam/realmonitor?channel=1&subtype=1' }
    ],
    hasPtz: true, isImported: false
  },
  {
    id: 'onvif_discov_04', name: 'Axis Dome Seguranca', manufacturer: 'Axis',
    model: 'M3045-V Network Camera', ip: '192.168.1.225', port: 80, macAddress: '00:40:8C:55:66:77',
    profiles: [
      { name: 'Profile_1080P', token: 'prof_main', resolution: '1920x1080', codec: 'H.264', rtspUri: 'rtsp://root:pass@192.168.1.225:554/axis-media/media.amp' }
    ],
    hasPtz: false, isImported: false
  }
]

const discoveredDevices = ref<DiscoveredOnvifCamera[]>(initialDiscovered)
const lastScanTime = ref<string>('Hoje, 18:45')

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
    }, 1800)
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
