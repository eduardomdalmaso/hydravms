<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import L from 'leaflet'
import type { MapMarkerItem } from '../../../types/mapTree'
import { getCameraTooltipHtml, getAlarmTooltipHtml } from './mapTooltipHelper'

const props = defineProps<{
  initialCenter: [number, number]
  initialZoom: number
  markers: MapMarkerItem[]
  flyToTarget?: { lat: number; lng: number } | null
  isLocked?: boolean
}>()

const emit = defineEmits<{
  (e: 'dropCoords', coords: { lat: number; lng: number }): void
  (e: 'removeMarker', markerId: string): void
}>()

const mapContainer = ref<HTMLElement | null>(null)
let map: L.Map | null = null
const leafletMarkers = new Map<string, L.Marker>()

const createIcon = (type: 'CAMERA' | 'ALARME') => {
  const isCam = type === 'CAMERA'
  const bg = isCam ? '#ff5e3a' : '#07080c'
  const border = isCam ? '#ffffff' : '#ff5e3a'
  const svg = isCam
    ? '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><path d="M23 7l-7 5 7 5V7z"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>'
    : '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>'
  return L.divIcon({
    className: 'vms-custom-pin',
    html: `<div style="background: ${bg}; width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; box-shadow: 0 0 12px rgba(255, 94, 58, 0.5); border: 2px solid ${border};">${svg}</div>`,
    iconSize: [28, 28], iconAnchor: [14, 14]
  })
}

const syncMarkers = () => {
  if (!map) return
  leafletMarkers.forEach(m => m.remove()); leafletMarkers.clear()
  props.markers.forEach(mk => {
    const marker = L.marker([mk.lat, mk.lng], { icon: createIcon(mk.type), draggable: !props.isLocked })
    const tipHtml = mk.type === 'CAMERA' ? getCameraTooltipHtml(mk.name, mk.lat, mk.lng) : getAlarmTooltipHtml(mk.name, mk.lat, mk.lng)
    marker.bindTooltip(tipHtml, { direction: 'top', className: 'vms-map-tooltip', opacity: 1, offset: L.point(0, -16) })
    marker.addTo(map!)
    leafletMarkers.set(mk.id, marker)
  })
}

onMounted(() => {
  if (!mapContainer.value) return
  map = L.map(mapContainer.value, { zoomControl: false }).setView(props.initialCenter, props.initialZoom)
  L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
    maxZoom: 19, attribution: '&copy; OpenStreetMap &copy; CARTO'
  }).addTo(map)
  L.control.zoom({ position: 'bottomright' }).addTo(map)
  syncMarkers()
})

onUnmounted(() => { if (map) { map.remove(); map = null } })
watch(() => props.markers, syncMarkers, { deep: true })
watch(() => props.flyToTarget, (t) => { if (map && t) map.flyTo([t.lat, t.lng], 16) })

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  if (!map || !mapContainer.value || props.isLocked) return
  const rect = mapContainer.value.getBoundingClientRect()
  const pt = map.containerPointToLatLng([e.clientX - rect.left, e.clientY - rect.top])
  emit('dropCoords', { lat: pt.lat, lng: pt.lng })
}
</script>

<template>
  <div ref="mapContainer" class="vms-map-canvas-container" @dragover.prevent @drop="handleDrop" />
</template>

<style scoped>
.vms-map-canvas-container { width: 100%; height: 100%; min-height: 520px; background: #07080c; border-radius: 8px; overflow: hidden; position: relative; }
</style>
