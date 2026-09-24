<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import type { ZoneConfig } from '../../../types/marketplace'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useWebRTCPlayer } from '../../../composables/useWebRTCPlayer'
import { useEventBus, initEventSocket } from '../../../services/eventSocket'

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

interface DetectionBox {
  id: string
  label: string
  conf: number
  x: number
  y: number
  w: number
  h: number
  expiresAt: number
}

const activeDetections = ref<DetectionBox[]>([])
const { subscribe } = useEventBus()
let unsubscribe: (() => void) | null = null
let animFrameId: number | null = null
let fpsTimerId: number | null = null

const renderCanvas = () => {
  if (!canvasRef.value || !props.isActive) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const w = canvasRef.value.width, h = canvasRef.value.height
  ctx.clearRect(0, 0, w, h)

  // Draw configured Zones: Orange Dashed Line, NO background color
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

  // Draw Live Detections from Backend Engine
  const now = Date.now()
  activeDetections.value = activeDetections.value.filter(d => d.expiresAt > now)

  activeDetections.value.forEach(d => {
    const bx = d.x * w, by = d.y * h, bw = d.w * w, bh = d.h * h

    ctx.strokeStyle = '#00ff9d'; ctx.lineWidth = 2
    ctx.strokeRect(bx, by, bw, bh)

    // Corner HUD Reticles
    const len = 8
    ctx.strokeStyle = '#00f0ff'; ctx.lineWidth = 3
    ctx.beginPath()
    ctx.moveTo(bx, by + len); ctx.lineTo(bx, by); ctx.lineTo(bx + len, by)
    ctx.moveTo(bx + bw - len, by); ctx.lineTo(bx + bw, by); ctx.lineTo(bx + bw, by + len)
    ctx.stroke()

    // Label Badge
    ctx.fillStyle = 'rgba(0, 0, 0, 0.85)'
    ctx.fillRect(bx, by - 18, 110, 18)
    ctx.fillStyle = '#00ff9d'
    ctx.font = '10px "JetBrains Mono", monospace'
    ctx.fillText(`[${d.label}] ${(d.conf * 100).toFixed(0)}%`, bx + 4, by - 5)
  })

  animFrameId = requestAnimationFrame(renderCanvas)
}

const startFpsStream = () => {
  if (fpsTimerId) clearInterval(fpsTimerId)
  if (!props.cameraId || !props.isActive) return

  // Fluid UI refresh rate (minimum 15 FPS for smooth monitoring)
  const targetFps = Math.max(15, Math.min(50, props.fps || 15))
  const intervalMs = Math.floor(1000 / targetFps)

  const fetchNextFrame = () => {
    const base = getCameraSnapshotUrl(props.cameraId, false)
    const nextUrl = `${base}&_t=${Date.now()}`
    const img = new Image()
    img.onload = () => {
      currentFrameUrl.value = nextUrl
    }
    img.src = nextUrl
  }

  fetchNextFrame()
  fpsTimerId = window.setInterval(fetchNextFrame, intervalMs)
}

const initStream = () => {
  if (props.cameraId) {
    if (props.isActive) {
      start(props.cameraId, true)
      startFpsStream()
      animFrameId = requestAnimationFrame(renderCanvas)
    }
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

onMounted(() => {
  initStream()
  initEventSocket()
  unsubscribe = subscribe((evt) => {
    if (!props.isActive) return
    const isTargetCam = evt.subject === props.cameraId || evt.data?.camera_id === props.cameraId
    if (isTargetCam && evt.data?.bbox) {
      const [x, y, w, h] = evt.data.bbox
      activeDetections.value.push({
        id: evt.id || String(Date.now()),
        label: (evt.data.object_label || evt.data.class || 'PESSOA').toUpperCase(),
        conf: evt.data.confidence || 0.95,
        x, y, w, h,
        expiresAt: Date.now() + 1500
      })
    }
  })
})

onUnmounted(() => {
  stop()
  if (fpsTimerId) clearInterval(fpsTimerId)
  if (animFrameId) cancelAnimationFrame(animFrameId)
  if (unsubscribe) unsubscribe()
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

    <!-- Realtime FPS Frame Stream (when WebRTC is establishing or in fallback mode) -->
    <img
      v-if="!isPlaying"
      :src="currentFrameUrl"
      style="position: absolute; inset: 0; width: 100%; height: 100%; object-fit: contain;"
      alt="Fluxo de Vídeo"
    />

    <!-- HUD Vector & Bounding Boxes Canvas -->
    <canvas ref="canvasRef" width="800" height="450" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none; z-index: 10;" />

    <!-- Top Stream Status Badge (Clean HUD) -->
    <div style="position: absolute; top: 8px; left: 10px; z-index: 20;">
      <span class="vms-badge" :class="isActive ? 'vms-badge-green' : 'vms-badge-orange'" style="font-size: 9px; font-weight: bold;">
        {{ isActive ? '● AO VIVO' : '■ PAUSADO' }}
      </span>
    </div>

    <!-- Paused Mask -->
    <div v-if="!isActive" style="position: absolute; inset: 0; background: rgba(7, 8, 12, 0.75); backdrop-filter: blur(2px); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; z-index: 30;">
      <span class="vms-text-mono vms-text-sm" style="color: #fcee0a; font-weight: bold;">[PAUSADO // FLUXO INTERROMPIDO]</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">CLIQUE EM "RETOMAR" PARA REINICIAR A DETECÇÃO EM TEMPO REAL</span>
    </div>
  </div>
</template>
