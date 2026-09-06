import { ref, computed } from 'vue'
import type { EnterpriseMapItem, DeployableDevice, MapMarkerItem } from '../types/mapTree'
import { initialDeployableDevices } from '../data/mockMapFolders'

export function useMapEditorState(mapItem: EnterpriseMapItem, onNotify: (msg: string) => void) {
  const devices = ref<DeployableDevice[]>([...initialDeployableDevices])
  const activeTab = ref<'CAMERAS' | 'ALARMES' | null>(null)
  const draggedDevice = ref<DeployableDevice | null>(null)
  const mapCenter = ref<[number, number]>([...mapItem.initialCenter])
  const mapZoom = ref<number>(mapItem.initialZoom)

  const camerasList = computed(() => devices.value.filter(d => d.type === 'CAMERA'))
  const alarmsList = computed(() => devices.value.filter(d => d.type === 'ALARME'))
  const currentTabDevices = computed(() => activeTab.value === 'CAMERAS' ? camerasList.value : (activeTab.value === 'ALARMES' ? alarmsList.value : []))

  const toggleTab = (tab: 'CAMERAS' | 'ALARMES') => {
    activeTab.value = activeTab.value === tab ? null : tab
  }

  const handleDeviceDragStart = (dev: DeployableDevice) => {
    if (mapItem.is_locked) return
    draggedDevice.value = dev
  }

  const addMarkerAt = (dev: DeployableDevice, lat: number, lng: number) => {
    if (mapItem.is_locked) { onNotify('[TRAVADO] Desbloqueie o mapa para posicionar dispositivos'); return }
    const existing = mapItem.markers.find(m => m.deviceId === dev.id)
    if (existing) {
      existing.lat = lat; existing.lng = lng
      onNotify(`[PIN ATUALIZADO] ${dev.name} reposicionado para [${lat.toFixed(5)}, ${lng.toFixed(5)}]`)
    } else {
      const newMarker: MapMarkerItem = {
        id: `m_${Date.now()}`,
        deviceId: dev.id,
        name: dev.name,
        type: dev.type,
        lat,
        lng,
        angle: dev.type === 'CAMERA' ? 45 : undefined,
        status: 'ONLINE'
      }
      mapItem.markers.push(newMarker)
      onNotify(`[PIN FIXADO] ${dev.name} adicionado em [${lat.toFixed(5)}, ${lng.toFixed(5)}]`)
    }
  }

  const applySavedCoords = (dev: DeployableDevice) => {
    if (mapItem.is_locked) { onNotify('[TRAVADO] Desbloqueie o mapa para posicionar'); return }
    if (!dev.hasCoords || dev.lat === undefined || dev.lng === undefined) {
      onNotify(`[SEM COORDENADAS] ${dev.name} nao possui geolocalizacao salva`)
      return
    }
    addMarkerAt(dev, dev.lat, dev.lng)
    mapCenter.value = [dev.lat, dev.lng]
  }

  const removeMarker = (markerId: string) => {
    if (mapItem.is_locked) { onNotify('[TRAVADO] Desbloqueie o mapa para remover dispositivos'); return }
    const idx = mapItem.markers.findIndex(m => m.id === markerId)
    if (idx >= 0) {
      const name = mapItem.markers[idx].name
      mapItem.markers.splice(idx, 1)
      onNotify(`[PIN REMOVIDO] ${name} desvinculado do mapa`)
    }
  }

  return {
    devices, activeTab, draggedDevice, mapCenter, mapZoom, camerasList, alarmsList,
    currentTabDevices, toggleTab, handleDeviceDragStart, addMarkerAt, applySavedCoords, removeMarker
  }
}
