<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Point2D, Line2D, ZoneConfig } from '../../../types/marketplace'

const props = defineProps<{
  snapshotUrl?: string
  zones: ZoneConfig[]
  activeZoneId: string
  zoneColors: string[]
}>()

const emit = defineEmits<{
  (e: 'update:points', payload: { zoneId: string; points: Point2D[] }): void
  (e: 'update:line', payload: { zoneId: string; line: Line2D }): void
  (e: 'select:zone', zoneId: string): void
}>()

const svgRef = ref<SVGSVGElement | null>(null)
const draggingIndex = ref<number | null>(null)
const draggingLinePoint = ref<'p1' | 'p2' | null>(null)
const hasImageError = ref(false)

const activeZone = computed(() => props.zones.find(z => z.id === props.activeZoneId) || props.zones[0])
const activeColor = computed(() => {
  const idx = props.zones.findIndex(z => z.id === props.activeZoneId)
  return props.zoneColors[idx % props.zoneColors.length] || '#ff5e3a'
})

const getPolygonStr = (points: Point2D[]) => points.map(p => `${p.x * 100},${p.y * 100}`).join(' ')

const handleMouseDownPoint = (idx: number) => { draggingIndex.value = idx }
const handleMouseDownLine = (pt: 'p1' | 'p2') => { draggingLinePoint.value = pt }
const handleMouseUp = () => { draggingIndex.value = null; draggingLinePoint.value = null }

const handleMouseMove = (e: MouseEvent) => {
  if (!svgRef.value || !activeZone.value) return
  if (draggingIndex.value === null && draggingLinePoint.value === null) return
  const rect = svgRef.value.getBoundingClientRect()
  const x = Math.max(0.02, Math.min(0.98, (e.clientX - rect.left) / rect.width))
  const y = Math.max(0.02, Math.min(0.98, (e.clientY - rect.top) / rect.height))

  if (draggingIndex.value !== null) {
    const next = activeZone.value.polygon.map((p, i) => i === draggingIndex.value ? { x, y } : p)
    emit('update:points', { zoneId: activeZone.value.id, points: next })
  } else if (draggingLinePoint.value !== null && activeZone.value.line) {
    const nextLine = { ...activeZone.value.line, [draggingLinePoint.value]: { x, y } }
    emit('update:line', { zoneId: activeZone.value.id, line: nextLine })
  }
}
</script>

<template>
  <div class="vms-roi-canvas-box" @mouseup="handleMouseUp" @mouseleave="handleMouseUp" @mousemove="handleMouseMove">
    <img v-if="snapshotUrl && !hasImageError" :src="snapshotUrl" class="vms-roi-snapshot" alt="Camera Snapshot" @error="hasImageError = true" @load="hasImageError = false" />
    <div v-else class="vms-roi-placeholder">
      <div class="vms-flex-col" style="align-items: center; gap: 4px;">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="rgba(255, 94, 58, 0.4)" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// GRID VIRTUAL // SNAPSHOT AO VIVO</span>
      </div>
    </div>

    <svg ref="svgRef" viewBox="0 0 100 100" preserveAspectRatio="none" class="vms-roi-svg">
      <!-- Renderiza todas as Zonas com suas respectivas cores -->
      <g v-for="(z, zIdx) in zones" :key="z.id" @click="emit('select:zone', z.id)">
        <!-- Zona Inativa: Linha / Polígono Suave -->
        <template v-if="z.id !== activeZoneId">
          <line
            v-if="z.mode === 'counting' && z.line"
            :x1="z.line.p1.x * 100" :y1="z.line.p1.y * 100"
            :x2="z.line.p2.x * 100" :y2="z.line.p2.y * 100"
            :stroke="zoneColors[zIdx % zoneColors.length]" stroke-width="0.6" stroke-dasharray="2,1" opacity="0.6"
          />
          <polygon
            v-else
            :points="getPolygonStr(z.polygon)"
            :fill="`${zoneColors[zIdx % zoneColors.length]}18`"
            :stroke="zoneColors[zIdx % zoneColors.length]" stroke-width="0.5" stroke-dasharray="1,1" opacity="0.6"
          />
        </template>

        <!-- Zona Ativa: Totalmente interativa com Handles -->
        <template v-else>
          <template v-if="z.mode === 'counting' && z.line">
            <line
              :x1="z.line.p1.x * 100" :y1="z.line.p1.y * 100"
              :x2="z.line.p2.x * 100" :y2="z.line.p2.y * 100"
              :stroke="activeColor" stroke-width="0.9" stroke-dasharray="2,1"
            />
            <circle :cx="z.line.p1.x * 100" :cy="z.line.p1.y * 100" r="2.2" :fill="activeColor" stroke="#fff" stroke-width="0.5" class="vms-roi-handle" @mousedown.prevent="handleMouseDownLine('p1')" />
            <text :x="z.line.p1.x * 100" :y="z.line.p1.y * 100 - 3" :fill="activeColor" font-size="3" font-weight="bold" text-anchor="middle">A (IN)</text>
            <circle :cx="z.line.p2.x * 100" :cy="z.line.p2.y * 100" r="2.2" :fill="activeColor" stroke="#fff" stroke-width="0.5" class="vms-roi-handle" @mousedown.prevent="handleMouseDownLine('p2')" />
            <text :x="z.line.p2.x * 100" :y="z.line.p2.y * 100 - 3" :fill="activeColor" font-size="3" font-weight="bold" text-anchor="middle">B (OUT)</text>
          </template>

          <template v-else>
            <polygon :points="getPolygonStr(z.polygon)" :fill="`${activeColor}33`" :stroke="activeColor" stroke-width="0.8" />
            <g v-for="(p, pIdx) in z.polygon" :key="pIdx">
              <circle :cx="p.x * 100" :cy="p.y * 100" r="2.0" :fill="activeColor" stroke="#ffffff" stroke-width="0.5" class="vms-roi-handle" @mousedown.prevent="handleMouseDownPoint(pIdx)" />
              <text :x="p.x * 100" :y="p.y * 100 - 2.5" fill="#ffffff" font-size="2.6" font-weight="bold" text-anchor="middle">P{{ pIdx + 1 }}</text>
            </g>
          </template>
        </template>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.vms-roi-canvas-box { position: relative; width: 100%; height: 100%; min-height: 280px; background: #07090e; border: 1px solid var(--vms-border); border-radius: 6px; overflow: hidden; user-select: none; }
.vms-roi-snapshot { width: 100%; height: 100%; object-fit: cover; }
.vms-roi-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: radial-gradient(circle, rgba(255,94,58,0.06) 0%, rgba(7,9,14,0.95) 100%); }
.vms-roi-svg { position: absolute; top: 0; left: 0; width: 100%; height: 100%; z-index: 10; cursor: crosshair; }
.vms-roi-handle { cursor: grab; }
.vms-roi-handle:active { cursor: grabbing; fill: #ffffff; }
</style>
