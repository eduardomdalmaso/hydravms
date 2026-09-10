<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { drawCanvasTimeline } from '../../composables/useCanvasTimelinePainter'

const props = defineProps<{
  currentTime: number
  isExportMode?: boolean
  exportStart?: number
  exportEnd?: number
  recordedRanges?: { start: number; end: number }[]
}>()
const emit = defineEmits<{
  (e: 'seek', time: number): void
  (e: 'updateExportRange', range: { start: number; end: number }): void
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const zoomMinutes = ref(1440)
const dragType = ref<'none' | 'start' | 'end' | 'range' | 'pan'>('none')
const dragStartX = ref(0), dragStartTime = ref(0)
const initStart = ref(0), initEnd = ref(0)

const render = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  canvas.width = canvas.offsetWidth * window.devicePixelRatio
  canvas.height = canvas.offsetHeight * window.devicePixelRatio
  ctx.scale(window.devicePixelRatio, window.devicePixelRatio)
  drawCanvasTimeline(
    ctx, canvas.offsetWidth, canvas.offsetHeight, props.currentTime,
    zoomMinutes.value, !!props.isExportMode, props.exportStart, props.exportEnd,
    props.recordedRanges || []
  )
}

onMounted(render)
watch([() => props.currentTime, () => props.isExportMode, () => props.exportStart, () => props.exportEnd, () => props.recordedRanges, zoomMinutes], render)

const onMouseDown = (e: MouseEvent) => {
  if (!canvasRef.value) return
  const rect = canvasRef.value.getBoundingClientRect(), cssW = canvasRef.value.offsetWidth
  const windowSpan = zoomMinutes.value * 60000, windowStart = props.currentTime - windowSpan / 2
  const clickX = e.clientX - rect.left
  dragStartX.value = e.clientX; dragStartTime.value = props.currentTime
  initStart.value = props.exportStart || 0; initEnd.value = props.exportEnd || 0

  if (props.isExportMode && props.exportStart && props.exportEnd) {
    const x1 = ((props.exportStart - windowStart) / windowSpan) * cssW
    const x2 = ((props.exportEnd - windowStart) / windowSpan) * cssW
    if (Math.abs(clickX - x1) <= 12) { dragType.value = 'start'; return }
    if (Math.abs(clickX - x2) <= 12) { dragType.value = 'end'; return }
    if (clickX > Math.min(x1, x2) && clickX < Math.max(x1, x2)) { dragType.value = 'range'; return }
  }
  dragType.value = 'pan'
}

const onMouseMove = (e: MouseEvent) => {
  if (dragType.value === 'none' || !canvasRef.value) return
  const cssW = canvasRef.value.offsetWidth, windowSpan = zoomMinutes.value * 60000
  const windowStart = props.currentTime - windowSpan / 2
  const mouseTime = windowStart + ((e.clientX - canvasRef.value.getBoundingClientRect().left) / cssW) * windowSpan
  const deltaShift = ((e.clientX - dragStartX.value) / cssW) * windowSpan

  if (dragType.value === 'range') {
    emit('updateExportRange', { start: initStart.value + deltaShift, end: initEnd.value + deltaShift })
  } else if (dragType.value === 'start' && props.exportEnd) {
    emit('updateExportRange', { start: Math.min(props.exportEnd - 10000, mouseTime), end: props.exportEnd })
  } else if (dragType.value === 'end' && props.exportStart) {
    emit('updateExportRange', { start: props.exportStart, end: Math.max(props.exportStart + 10000, mouseTime) })
  } else if (dragType.value === 'pan') {
    emit('seek', Math.min(Date.now(), dragStartTime.value - deltaShift))
  }
}

const onMouseUp = (e: MouseEvent) => {
  if (dragType.value === 'pan' && Math.abs(e.clientX - dragStartX.value) < 4 && canvasRef.value) {
    const ratio = (e.clientX - canvasRef.value.getBoundingClientRect().left) / canvasRef.value.offsetWidth
    emit('seek', props.currentTime - (zoomMinutes.value * 60000) / 2 + ratio * (zoomMinutes.value * 60000))
  }
  dragType.value = 'none'
}
</script>

<template>
  <div class="vms-vezha-canvas-area">
    <div class="vms-vezha-scale-col">
      <button class="vms-vezha-scale-btn" title="Zoom In" @click="zoomMinutes = Math.max(1, Math.round(zoomMinutes / 2))">+</button>
      <button class="vms-vezha-scale-btn" title="24 Horas" @click="zoomMinutes = Math.min(1440, Math.round(zoomMinutes * 2))">-</button>
    </div>
    <div
      class="vms-vezha-canvas-wrapper"
      :style="{ cursor: dragType === 'range' ? 'move' : isExportMode ? 'ew-resize' : dragType !== 'none' ? 'grabbing' : 'grab' }"
      title="Arraste para mover o período ou ajuste as barras de início/fim do recorte"
      @wheel.prevent="e => zoomMinutes = e.deltaY < 0 ? Math.max(1, Math.round(zoomMinutes / 1.5)) : Math.min(1440, Math.round(zoomMinutes * 1.5))"
      @mousedown="onMouseDown" @mousemove="onMouseMove" @mouseup="onMouseUp" @mouseleave="dragType = 'none'"
    ><canvas ref="canvasRef" class="vms-vezha-canvas"></canvas></div>
  </div>
</template>
