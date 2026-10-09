<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import type { ZoneConfig } from '../../../types/marketplace'
import { getCameraMjpegUrl } from '../../../utils/streamUrls'
import { useLiveDetections } from '../../../composables/useLiveDetections'

const props = defineProps<{
  cameraId: string
  cameraName: string
  isActive: boolean
  zones?: ZoneConfig[]
  fps?: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const liveStreamUrl = ref('')

const { getDetections, startLiveAnalyticsPipeline } = useLiveDetections()
const cameraDetections = computed(() => getDetections(props.cameraId).value)

let animFrameId: number | null = null
let reconnectTimer: any = null

const renderCanvas = () => {
  if (!props.isActive) { animFrameId = null; return }
  if (!canvasRef.value) { animFrameId = requestAnimationFrame(renderCanvas); return }
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) { animFrameId = requestAnimationFrame(renderCanvas); return }

  const w = canvasRef.value.width, h = canvasRef.value.height
  ctx.clearRect(0, 0, w, h)

  if (props.zones) {
    props.zones.forEach((z) => {
      if (!z.polygon || z.polygon.length < 3) return
      ctx.beginPath()
      ctx.strokeStyle = '#ff5e3a'
      ctx.setLineDash([6, 4])
      ctx.lineWidth = 2
      z.polygon.forEach((pt, pIdx) => {
        pIdx === 0 ? ctx.moveTo(pt.x * w, pt.y * h) : ctx.lineTo(pt.x * w, pt.y * h)
      })
      ctx.closePath()
      ctx.stroke()
      ctx.setLineDash([])
    })
  }

  // Draw ONLY clean bounding boxes (solid orange, no text clutter)
  const boxes = cameraDetections.value
  if (boxes && boxes.length > 0) {
    boxes.forEach(b => {
      const bx = (b.box[0] / 100) * w
      const by = (b.box[1] / 100) * h
      const bw = (b.box[2] / 100) * w
      const bh = (b.box[3] / 100) * h

      ctx.strokeStyle = b.color || '#ff5e3a'
      ctx.lineWidth = 2
      ctx.strokeRect(bx, by, bw, bh)
    })
  }

  animFrameId = requestAnimationFrame(renderCanvas)
}

const updateStreamSource = () => {
  if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null }
  if (props.cameraId && props.isActive) {
    liveStreamUrl.value = `${getCameraMjpegUrl(props.cameraId)}?t=${Date.now()}`
    startLiveAnalyticsPipeline()
    if (!animFrameId) animFrameId = requestAnimationFrame(renderCanvas)
  } else {
    liveStreamUrl.value = ''
    if (animFrameId) {
      cancelAnimationFrame(animFrameId)
      animFrameId = null
    }
  }
}

const handleVisibilityOrFocus = () => {
  if (document.visibilityState === 'visible' && props.isActive) {
    updateStreamSource()
  }
}

const handleStreamError = () => {
  if (props.isActive && !reconnectTimer) {
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      updateStreamSource()
    }, 300)
  }
}

watch([() => props.cameraId, () => props.isActive], updateStreamSource, { immediate: true })

onMounted(() => {
  startLiveAnalyticsPipeline()
  updateStreamSource()
  if (!animFrameId) animFrameId = requestAnimationFrame(renderCanvas)
  document.addEventListener('visibilitychange', handleVisibilityOrFocus)
  window.addEventListener('focus', handleVisibilityOrFocus)
})

onUnmounted(() => {
  liveStreamUrl.value = ''
  if (animFrameId) {
    cancelAnimationFrame(animFrameId)
    animFrameId = null
  }
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  document.removeEventListener('visibilitychange', handleVisibilityOrFocus)
  window.removeEventListener('focus', handleVisibilityOrFocus)
})
</script>

<template>
  <div class="vms-card" style="position: relative; overflow: hidden; background: #000; border: 1px solid var(--vms-border); border-radius: 8px; aspect-ratio: 16/9; display: flex; align-items: center; justify-content: center;">
    <!-- Direct Real-time OpenCV MJPEG Stream -->
    <img
      v-if="liveStreamUrl"
      :src="liveStreamUrl"
      @error="handleStreamError"
      style="width: 100%; height: 100%; object-fit: contain; display: block;"
      alt="Fluxo em Tempo Real OpenCV"
    />

    <!-- HUD Vector & Bounding Boxes Canvas (Clean BBox, Zero Text Clutter) -->
    <canvas ref="canvasRef" width="800" height="450" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none; z-index: 10;" />

    <!-- Paused Mask -->
    <div v-if="!isActive" style="position: absolute; inset: 0; background: rgba(7, 8, 12, 0.75); backdrop-filter: blur(2px); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; z-index: 30;">
      <span class="vms-text-mono vms-text-sm" style="color: #fcee0a; font-weight: bold;">[PAUSADO // FLUXO INTERROMPIDO]</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">CLIQUE EM "RETOMAR" PARA REINICIAR A DETECÇÃO EM TEMPO REAL</span>
    </div>
  </div>
</template>
