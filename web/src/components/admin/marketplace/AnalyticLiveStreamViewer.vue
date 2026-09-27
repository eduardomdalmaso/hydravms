<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import type { ZoneConfig } from '../../../types/marketplace'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useWebRTCPlayer } from '../../../composables/useWebRTCPlayer'
import { useLiveDetections } from '../../../composables/useLiveDetections'

const props = defineProps<{
  cameraId: string
  cameraName: string
  isActive: boolean
  zones?: ZoneConfig[]
  fps?: number
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const currentFrameUrl = ref('')

const { isPlaying, start, stop } = useWebRTCPlayer(videoRef)
const { getDetections } = useLiveDetections()
const cameraDetections = computed(() => getDetections(props.cameraId).value)

let animFrameId: number | null = null
let fpsTimerId: number | null = null

const renderCanvas = () => {
  if (!canvasRef.value || !props.isActive) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const w = canvasRef.value.width, h = canvasRef.value.height
  ctx.clearRect(0, 0, w, h)

  // Draw configured Zones: Orange Dashed Line, NO background
  if (props.zones) {
    props.zones.forEach((z) => {
      if (!z.polygon || z.polygon.length < 3) return
      ctx.beginPath()
      ctx.strokeStyle = '#ff5e3a'
      ctx.setLineDash([6, 4])
      ctx.lineWidth = 2
      z.polygon.forEach((pt, pIdx) => {
        const px = pt.x * w, py = pt.y * h
        pIdx === 0 ? ctx.moveTo(px, py) : ctx.lineTo(px, py)
      })
      ctx.closePath()
      ctx.stroke()
      ctx.setLineDash([])
    })
  }

  // Draw ONLY clean bounding boxes (no text, no labels)
  const boxes = cameraDetections.value
  boxes.forEach(b => {
    // box is [x_pct, y_pct, w_pct, h_pct] where each is 0..100
    const bx = (b.box[0] / 100) * w
    const by = (b.box[1] / 100) * h
    const bw = (b.box[2] / 100) * w
    const bh = (b.box[3] / 100) * h

    ctx.strokeStyle = b.color || '#ff5e3a'
    ctx.lineWidth = 1.2
    ctx.strokeRect(bx, by, bw, bh)
  })

  animFrameId = requestAnimationFrame(renderCanvas)
}

const startFpsStream = () => {
  if (fpsTimerId) clearInterval(fpsTimerId)
  if (!props.cameraId || !props.isActive) return

  const targetFps = Math.max(15, Math.min(50, props.fps || 15))
  const intervalMs = Math.floor(1000 / targetFps)

  const fetchNextFrame = () => {
    const base = getCameraSnapshotUrl(props.cameraId, false)
    const nextUrl = `${base}&_t=${Date.now()}`
    const img = new Image()
    img.onload = () => { currentFrameUrl.value = nextUrl }
    img.src = nextUrl
  }

  fetchNextFrame()
  fpsTimerId = window.setInterval(fetchNextFrame, intervalMs)
}

const initStream = () => {
  if (props.cameraId && props.isActive) {
    start(props.cameraId, true)
    startFpsStream()
    if (!animFrameId) animFrameId = requestAnimationFrame(renderCanvas)
  }
}

watch([() => props.cameraId, () => props.fps], initStream, { immediate: true })
watch(() => props.isActive, (active) => {
  if (active) {
    initStream()
  } else {
    stop()
    if (fpsTimerId) { clearInterval(fpsTimerId); fpsTimerId = null }
    if (animFrameId) { cancelAnimationFrame(animFrameId); animFrameId = null }
  }
}, { immediate: true })

onMounted(() => { initStream() })
onUnmounted(() => {
  stop()
  if (fpsTimerId) clearInterval(fpsTimerId)
  if (animFrameId) cancelAnimationFrame(animFrameId)
})
</script>

<template>
  <div class="vms-card" style="position: relative; overflow: hidden; background: #000; border: 1px solid var(--vms-border); border-radius: 8px; aspect-ratio: 16/9; display: flex; align-items: center; justify-content: center;">
    <!-- Live Video Element (WebRTC WHEP) -->
    <video
      ref="videoRef"
      autoplay
      muted
      playsinline
      style="width: 100%; height: 100%; object-fit: contain; background: #000;"
      :style="{ opacity: isPlaying ? 1 : 0 }"
    />

    <!-- Realtime FPS Frame Stream Fallback -->
    <img
      v-if="!isPlaying"
      :src="currentFrameUrl"
      style="position: absolute; inset: 0; width: 100%; height: 100%; object-fit: contain;"
      alt="Fluxo de Vídeo"
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
