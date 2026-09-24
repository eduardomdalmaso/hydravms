<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Point2D, Line2D, AnalyticMode } from '../../../types/marketplace'

const props = defineProps<{
  snapshotUrl?: string
  mode: AnalyticMode
  points: Point2D[]
  line: Line2D
}>()

const emit = defineEmits<{
  (e: 'update:points', val: Point2D[]): void
  (e: 'update:line', val: Line2D): void
}>()

const svgRef = ref<SVGSVGElement | null>(null)
const draggingIndex = ref<number | null>(null)
const draggingLinePoint = ref<'p1' | 'p2' | null>(null)

const polygonPointsStr = computed(() =>
  props.points.map(p => `${p.x * 100},${p.y * 100}`).join(' ')
)

const handleMouseDownPoint = (idx: number) => { draggingIndex.value = idx }
const handleMouseDownLine = (pt: 'p1' | 'p2') => { draggingLinePoint.value = pt }
const handleMouseUp = () => { draggingIndex.value = null; draggingLinePoint.value = null }

const handleMouseMove = (e: MouseEvent) => {
  if (!svgRef.value) return
  if (draggingIndex.value === null && draggingLinePoint.value === null) return
  const rect = svgRef.value.getBoundingClientRect()
  const rawX = (e.clientX - rect.left) / rect.width
  const rawY = (e.clientY - rect.top) / rect.height
  const x = Math.max(0.02, Math.min(0.98, rawX))
  const y = Math.max(0.02, Math.min(0.98, rawY))

  if (draggingIndex.value !== null) {
    const next = props.points.map((p, i) => i === draggingIndex.value ? { x, y } : p)
    emit('update:points', next)
  } else if (draggingLinePoint.value !== null) {
    const nextLine = { ...props.line, [draggingLinePoint.value]: { x, y } }
    emit('update:line', nextLine)
  }
}
</script>

<template>
  <div class="vms-roi-canvas-box" @mouseup="handleMouseUp" @mouseleave="handleMouseUp" @mousemove="handleMouseMove">
    <img v-if="snapshotUrl" :src="snapshotUrl" class="vms-roi-snapshot" alt="Camera Snapshot" />
    <div v-else class="vms-roi-placeholder">
      <span class="vms-text-mono vms-text-2xs vms-text-dim">// GRID VIRTUAL // SNAPSHOT DA CÂMERA</span>
    </div>

    <svg ref="svgRef" viewBox="0 0 100 100" preserveAspectRatio="none" class="vms-roi-svg">
      <!-- 1. Linha Bidirecional para Contagem -->
      <template v-if="mode === 'counting'">
        <line
          :x1="line.p1.x * 100" :y1="line.p1.y * 100"
          :x2="line.p2.x * 100" :y2="line.p2.y * 100"
          stroke="#00f0ff" stroke-width="0.8" stroke-dasharray="2,1"
        />
        <circle
          :cx="line.p1.x * 100" :cy="line.p1.y * 100" r="2.2" fill="#00f0ff" stroke="#fff" stroke-width="0.5"
          class="vms-roi-handle" @mousedown.prevent="handleMouseDownLine('p1')"
        />
        <text :x="line.p1.x * 100" :y="line.p1.y * 100 - 3" fill="#00f0ff" font-size="3" font-weight="bold" text-anchor="middle">A (IN)</text>
        <circle
          :cx="line.p2.x * 100" :cy="line.p2.y * 100" r="2.2" fill="#ff5e3a" stroke="#fff" stroke-width="0.5"
          class="vms-roi-handle" @mousedown.prevent="handleMouseDownLine('p2')"
        />
        <text :x="line.p2.x * 100" :y="line.p2.y * 100 - 3" fill="#ff5e3a" font-size="3" font-weight="bold" text-anchor="middle">B (OUT)</text>
      </template>

      <!-- 2. Polígono de 6 Pontos (Intrusão, Multidão, Tempo Excedido) -->
      <template v-else>
        <polygon :points="polygonPointsStr" fill="rgba(255, 94, 58, 0.22)" stroke="#ff5e3a" stroke-width="0.6" />
        <g v-for="(p, idx) in points" :key="idx">
          <circle
            :cx="p.x * 100" :cy="p.y * 100" r="2.0" fill="#ff5e3a" stroke="#ffffff" stroke-width="0.5"
            class="vms-roi-handle" @mousedown.prevent="handleMouseDownPoint(idx)"
          />
          <text :x="p.x * 100" :y="p.y * 100 - 2.5" fill="#ffffff" font-size="2.6" font-weight="bold" text-anchor="middle">P{{ idx + 1 }}</text>
        </g>
      </template>
    </svg>
  </div>
</template>

<style scoped>
.vms-roi-canvas-box { position: relative; width: 100%; height: 230px; background: #07090e; border: 1px solid var(--vms-border); border-radius: 6px; overflow: hidden; user-select: none; }
.vms-roi-snapshot { width: 100%; height: 100%; object-fit: cover; }
.vms-roi-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: radial-gradient(circle, rgba(255,94,58,0.06) 0%, rgba(7,9,14,0.95) 100%); }
.vms-roi-svg { position: absolute; top: 0; left: 0; width: 100%; height: 100%; z-index: 10; cursor: crosshair; }
.vms-roi-handle { cursor: grab; transition: transform 0.1s; }
.vms-roi-handle:active { cursor: grabbing; fill: #00ff9d; }
</style>
