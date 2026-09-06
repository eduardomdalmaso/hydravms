<script setup lang="ts">
import { ref } from 'vue'
import type { EnterpriseMapItem, DeployableDevice } from '../../../types/mapTree'
import { useMapEditorState } from '../../../composables/useMapEditorState'
import MapOsmCanvas from './MapOsmCanvas.vue'
import MapTopRightTabs from './MapTopRightTabs.vue'
import MapDeviceDrawer from './MapDeviceDrawer.vue'
import MapLocationSearchBar from './MapLocationSearchBar.vue'

const props = defineProps<{ mapItem: EnterpriseMapItem }>()
const emit = defineEmits<{ (e: 'notify', msg: string): void }>()

const flyTarget = ref<{ lat: number; lng: number } | null>(null)
const {
  activeTab, draggedDevice, camerasList, alarmsList, currentTabDevices,
  toggleTab, handleDeviceDragStart, addMarkerAt, applySavedCoords, removeMarker
} = useMapEditorState(props.mapItem, (m) => emit('notify', m))

const handleCanvasDrop = (coords: { lat: number; lng: number }) => {
  if (draggedDevice.value) {
    addMarkerAt(draggedDevice.value, coords.lat, coords.lng)
    draggedDevice.value = null
  }
}

const handleSelectLocation = (loc: { lat: number; lng: number; label: string }) => {
  flyTarget.value = { lat: loc.lat, lng: loc.lng }
  emit('notify', `[GEOCODER] Centralizado em: ${loc.label.slice(0, 45)}...`)
}

const handleQuickAdd = (type: 'CAMERA' | 'ALARME', loc: { lat: number; lng: number; label: string }) => {
  const pool = type === 'CAMERA' ? camerasList.value : alarmsList.value
  const candidate = pool.find(d => !props.mapItem.markers.some(m => m.deviceId === d.id)) || pool[0]
  if (candidate) addMarkerAt(candidate, loc.lat, loc.lng)
}
</script>

<template>
  <div class="vms-map-editor-wrapper">
    <!-- Mapa OpenSource Interativo com Leaflet -->
    <MapOsmCanvas
      :initial-center="mapItem.initialCenter"
      :initial-zoom="mapItem.initialZoom"
      :markers="mapItem.markers"
      :fly-to-target="flyTarget"
      :is-locked="mapItem.is_locked"
      @drop-coords="handleCanvasDrop"
      @remove-marker="removeMarker"
    />

    <!-- Busca de Localização Estilo Google Maps (Lado Esquerdo) -->
    <MapLocationSearchBar
      @select-location="handleSelectLocation"
      @quick-add="handleQuickAdd"
    />

    <!-- Abas Superiores Direitas [CAMERAS] e [ALARMES] -->
    <MapTopRightTabs
      :active-tab="activeTab"
      :cameras-count="camerasList.length"
      :alarms-count="alarmsList.length"
      @toggle="toggleTab"
    />

    <!-- Gaveta Retrátil com Transição Suave -->
    <Transition name="vms-toast">
      <MapDeviceDrawer
        v-if="activeTab"
        :active-tab="activeTab"
        :devices="currentTabDevices"
        :is-locked="mapItem.is_locked"
        @close="activeTab = null"
        @drag-start="handleDeviceDragStart"
        @apply-coords="applySavedCoords"
      />
    </Transition>
  </div>
</template>

<style scoped>
.vms-map-editor-wrapper {
  position: relative; width: 100%; height: calc(100vh - 165px);
  min-height: 540px; border-radius: 8px; border: 1px solid var(--vms-border);
  overflow: hidden; background: #07080c;
}
</style>
