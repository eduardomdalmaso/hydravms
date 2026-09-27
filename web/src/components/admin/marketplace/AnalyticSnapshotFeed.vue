<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useLiveDetections } from '../../../composables/useLiveDetections'
import type { ZoneConfig } from '../../../types/marketplace'

const props = defineProps<{
  cameraId: string
  cameraName: string
  analyticName: string
  isActive?: boolean
  zones?: ZoneConfig[]
}>()

interface LocalSnapshot {
  id: string
  time: string
  label: string
  conf: number
  zone: string
  url: string
  bbox?: [number, number, number, number]
}

const snapshots = ref<LocalSnapshot[]>([])
const { getDetections } = useLiveDetections()
const liveBoxes = getDetections(props.cameraId)
let lastSnapTime = 0

// Point-in-polygon helper to detect if person is inside ROI zone
const isPointInPolygon = (x: number, y: number, poly: { x: number; y: number }[]) => {
  let inside = false
  for (let i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    const xi = poly[i].x * 100, yi = poly[i].y * 100
    const xj = poly[j].x * 100, yj = poly[j].y * 100
    const intersect = ((yi > y) !== (yj > y)) && (x < (xj - xi) * (y - yi) / (yj - yi) + xi)
    if (intersect) inside = !inside
  }
  return inside
}

watch(liveBoxes, (boxes) => {
  if (props.isActive === false || !boxes || boxes.length === 0) return
  const now = Date.now()
  if (now - lastSnapTime < 3000) return // Cooldown 3s between snapshot captures

  for (const b of boxes) {
    if (b.category === 'person') {
      const centerX = b.box[0] + b.box[2] / 2
      const centerY = b.box[1] + b.box[3] / 2

      let zoneName = 'ÁREA MONITORADA'
      let isIntrusion = false

      if (props.zones && props.zones.length > 0) {
        for (const z of props.zones) {
          if (z.polygon && z.polygon.length >= 3 && isPointInPolygon(centerX, centerY, z.polygon)) {
            zoneName = `ROI // ${z.name.toUpperCase()}`
            isIntrusion = true
            break
          }
        }
      } else {
        zoneName = 'ROI // INTRUSÃO'
        isIntrusion = true
      }

      if (isIntrusion) {
        lastSnapTime = now
        const snapUrl = `${getCameraSnapshotUrl(props.cameraId, true)}&_t=${now}`
        const time = new Date().toLocaleTimeString()
        const conf = Math.round((b.confidence || 0.85) * 100)

        snapshots.value.unshift({
          id: `snap-${now}-${Math.random().toString(36).substr(2, 4)}`,
          time,
          label: 'INTRUSÃO // PESSOA',
          conf,
          zone: zoneName,
          url: snapUrl,
          bbox: b.box
        })

        if (snapshots.value.length > 30) snapshots.value.pop()
        break
      }
    }
  }
}, { deep: true })
</script>

<template>
  <div class="vms-card vms-flex-col" style="flex: 1; height: 100%; max-height: calc(100vh - 210px); display: flex; flex-direction: column; background: #0b0e14; border: 1px solid var(--vms-border); border-radius: 8px; overflow: hidden;">
    <div class="vms-flex-between" style="padding: 0.75rem 1rem; border-bottom: 1px solid var(--vms-border); background: rgba(255,255,255,0.02);">
      <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">SNAPSHOTS // DETECÇÕES</span>
      <span class="vms-badge" :class="isActive !== false ? 'vms-badge-green' : 'vms-badge-orange'" style="font-size: 9px;">
        {{ isActive !== false ? `${snapshots.length} EVENTOS` : 'PAUSADO' }}
      </span>
    </div>

    <div class="vms-flex-col" style="padding: 0.75rem; gap: 0.65rem; overflow-y: auto; flex: 1;">
      <div v-if="snapshots.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="text-align: center; padding: 2rem 0;">
        // AGUARDANDO DETECÇÃO DE PESSOA NO ROI...
      </div>

      <div
        v-for="s in snapshots"
        :key="s.id"
        class="vms-card"
        style="padding: 0.5rem; background: #121824; border: 1px solid rgba(255, 0, 60, 0.25); border-radius: 6px; display: flex; gap: 0.65rem; align-items: center;"
      >
        <div style="width: 72px; height: 50px; border-radius: 4px; overflow: hidden; background: #000; flex-shrink: 0; border: 1px solid rgba(255, 0, 60, 0.4); position: relative;">
          <img :src="s.url" style="width: 100%; height: 100%; object-fit: cover;" alt="Snapshot" />
          <svg v-if="s.bbox" viewBox="0 0 100 100" preserveAspectRatio="none" style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none;">
            <rect :x="s.bbox[0]" :y="s.bbox[1]" :width="s.bbox[2]" :height="s.bbox[3]" fill="none" stroke="#ff003c" stroke-width="1.2" stroke-dasharray="2, 1" />
          </svg>
        </div>

        <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 0;">
          <div class="vms-flex-between" style="align-items: center;">
            <span class="vms-badge vms-badge-red" style="font-size: 8.5px; font-weight: bold; display: flex; align-items: center; gap: 3px;">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><circle cx="13.5" cy="4.5" r="2"/><path d="M13.8 8.2c-.4-.4-.9-.7-1.5-.7-.8 0-1.5.4-1.9 1l-2.9 3.9c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l2-2.7v4.5l-3.2 4.3c-.3.4-.2 1 .2 1.3.4.3 1 .2 1.3-.2l3.4-4.5 2.1 3v4.4c0 .6.4 1 1 1s1-.4 1-1v-5c0-.4-.2-.7-.5-.9l-2.4-3.4.5-4.4 2 1.5c.2.2.5.2.8.2.3 0 .6-.1.8-.3.4-.4.4-1 0-1.4l-3.2-2.3z"/></svg>
              <span>[{{ s.label }}] {{ s.conf }}%</span>
            </span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ s.time }}</span>
          </div>
          <span class="vms-text-mono vms-text-2xs" style="color: #ff5e3a; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ s.zone }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
