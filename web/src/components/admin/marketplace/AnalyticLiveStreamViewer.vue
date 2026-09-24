<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import type { ZoneConfig } from '../../../types/marketplace'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useEventBus, initEventSocket } from '../../../services/eventSocket'

const props = defineProps<{
  cameraId: string
  cameraName: string
  isActive: boolean
  zones?: ZoneConfig[]
  fps?: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const imgRef = ref<HTMLImageElement | null>(null)
const snapshotUrl = ref('')

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

const renderCanvas = () => {
  if (!canvasRef.value || !props.isActive) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const w = canvasRef.value.width, h = canvasRef.value.height
  ctx.clearRect(0, 0, w, h)

  // Draw configured Zones
  if (props.zones) {
    props.zones.forEach((z, idx) => {
      if (!z.polygon || z.polygon.length < 3) return
      ctx.beginPath()
      ctx.strokeStyle = idx === 0 ? 'rgba(0, 240, 255, 0.85)' : 'rgba(252, 238, 10, 0.85)'
      ctx.fillStyle = idx === 0 ? 'rgba(0, 240, 255, 0.12)' : 'rgba(252, 238, 10, 0.12)'
      ctx.lineWidth = 2
      z.polygon.forEach((pt, pIdx) => {
        const px = pt.x * w, py = pt.y * h
        pIdx === 0 ? ctx.moveTo(px, py) : ctx.lineTo(px, py)
      })
      ctx.closePath()
      ctx.stroke()
      ctx.fill()
    })
  }

  // Draw Real Detections from Backend Engine
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

const loadSnapshot = () => {
  if (props.cameraId) snapshotUrl.value = getCameraSnapshotUrl(props.cameraId, false)
}

watch(() => props.cameraId, loadSnapshot, { immediate: true })
watch(() => props.isActive, (active) => {
  if (active) animFrameId = requestAnimationFrame(renderCanvas)
  else if (animFrameId) { cancelAnimationFrame(animFrameId); animFrameId = null }
}, { immediate: true })

onMounted(() => {
  loadSnapshot()
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
  animFrameId = requestAnimationFrame(renderCanvas)
})

onUnmounted(() => {
  if (animFrameId) cancelAnimationFrame(animFrameId)
  if (unsubscribe) unsubscribe()
})
</script>

<template>
  <div class="vms-card" style="position: relative; overflow: hidden; background: #000; border: 1px solid var(--vms-border); border-radius: 8px; aspect-ratio: 16/9; display: flex; align-items: center; justify-content: center;">
    <img ref="imgRef" :src="snapshotUrl" style="width: 100%; height: 100%; object-fit: contain;" alt="Fluxo Câmera" />
    <canvas ref="canvasRef" width="800" height="450" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none;" />

    <div style="position: absolute; top: 8px; left: 10px; display: flex; gap: 8px; align-items: center; z-index: 10;">
      <span class="vms-badge" :class="isActive ? 'vms-badge-green' : 'vms-badge-orange'" style="font-size: 9px; font-weight: bold;">
        {{ isActive ? '● AO VIVO' : '■ PAUSADO' }}
      </span>
      <span class="vms-text-mono vms-text-2xs" style="color: #fff; background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 4px;">
        {{ cameraName }} // {{ fps || 15 }} FPS
      </span>
    </div>

    <div v-if="!isActive" style="position: absolute; inset: 0; background: rgba(7, 8, 12, 0.75); backdrop-filter: blur(2px); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; z-index: 20;">
      <span class="vms-text-mono vms-text-sm" style="color: #fcee0a; font-weight: bold;">[PAUSADO // FLUXO INTERROMPIDO]</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">CLIQUE EM "RETOMAR" PARA REINICIAR A DETECÇÃO EM TEMPO REAL</span>
    </div>
  </div>
</template>
