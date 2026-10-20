<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount, nextTick } from "vue"
import type { WorkspaceSlot, CameraStreamInfo } from "../../types/mosaic"
import { useWSPCanvasPlayer } from "../../composables/useWSPCanvasPlayer"
import { useTimelinePlayback } from "../../composables/useTimelinePlayback"
import { useBranding } from "../../composables/useBranding"
import { useLiveDetections } from "../../composables/useLiveDetections"

const props = defineProps<{ slot: WorkspaceSlot; isActive?: boolean; isHero?: boolean; isExpanded?: boolean }>()
const emit = defineEmits<{
  (e: "selectCamera", cam: CameraStreamInfo): void
  (e: "clear", idx: number): void
  (e: "expand", cam: CameraStreamInfo): void
}>()
const { branding } = useBranding()
const { getDetections } = useLiveDetections()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null), isGearOpen = ref(false)
const manualProfile = ref<"SUB" | "MAIN" | null>(null)
const { start: startWsp, stop: stopWsp, isPlaying: isWspPlaying, isConnecting: isWspConnecting } = useWSPCanvasPlayer(canvasRef)
const { isLive, currentTime, isPlaying, playbackSpeed, activePlaybackCameraId, seekTrigger } = useTimelinePlayback()

const effectiveProfile = computed<"SUB" | "MAIN">(() => {
  if (manualProfile.value) return manualProfile.value
  return (props.isHero || props.isExpanded) ? "MAIN" : "SUB"
})

const handleExpand = () => {
  if (props.slot.type === 'camera' && props.slot.data) {
    manualProfile.value = 'MAIN'
    emit('expand', props.slot.data as CameraStreamInfo)
    handleSeekOrSwitch()
  }
}

const cameraDetections = computed(() => {
  if (props.slot.type === 'camera' && props.slot.data) {
    return getDetections((props.slot.data as CameraStreamInfo).id).value
  }
  return []
})

const isPlayback = computed(() => {
  const cam = props.slot.type === 'camera' ? props.slot.data as CameraStreamInfo : null
  return !!cam && !isLive.value && activePlaybackCameraId.value === cam.id
})

const isStreamBuffering = computed(() => {
  if (isPlayback.value) return !isPlaying.value
  return !isWspPlaying.value || isWspConnecting.value
})

const toggleProfile = () => {
  manualProfile.value = effectiveProfile.value === 'SUB' ? 'MAIN' : 'SUB'
  handleSeekOrSwitch()
}

const handleSeekOrSwitch = async () => {
  await nextTick()
  if (isPlayback.value && videoRef.value) {
    stopWsp()
    const cam = props.slot.data as CameraStreamInfo
    videoRef.value.src = `http://localhost:8083/api/v1/cameras/${cam.id}/recordings/stream?t=${currentTime.value}`
    videoRef.value.playbackRate = playbackSpeed.value
    if (isPlaying.value) videoRef.value.play().catch(() => {})
    else videoRef.value.pause()
  } else if (!isPlayback.value && props.slot.type === 'camera' && props.slot.data) {
    const cam = props.slot.data as CameraStreamInfo
    const hasSub = cam.protocol !== 'RTMP'
    const isHero = effectiveProfile.value === 'MAIN'
    startWsp(cam.id, isHero, hasSub)
  } else if (!isPlayback.value && !props.slot.data) {
    stopWsp()
  }
}

const currentCamId = computed(() => (props.slot.type === 'camera' && props.slot.data) ? (props.slot.data as CameraStreamInfo).id : '')

const handleVis = () => { if (!document.hidden) handleSeekOrSwitch() }
document.addEventListener('visibilitychange', handleVis)
window.addEventListener('focus', handleVis)
onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', handleVis)
  window.removeEventListener('focus', handleVis)
  stopWsp()
})

watch(() => [currentCamId.value, props.isHero], () => {
  manualProfile.value = null
  handleSeekOrSwitch()
})
watch(() => [isPlayback.value, seekTrigger.value], handleSeekOrSwitch)
watch(() => [isPlaying.value, playbackSpeed.value], () => {
  if (isPlayback.value && videoRef.value) {
    videoRef.value.playbackRate = playbackSpeed.value
    if (isPlaying.value && videoRef.value.paused) videoRef.value.play().catch(() => {})
    else if (!isPlaying.value && !videoRef.value.paused) videoRef.value.pause()
  }
})
onMounted(handleSeekOrSwitch)
</script>

<template>
  <div
    class="vms-slot"
    :class="{ active: isActive, 'vms-slot-hero': isHero }"
    draggable="true"
    @click="slot.type === 'camera' && slot.data && emit('selectCamera', slot.data as CameraStreamInfo)"
  >
    <div v-if="slot.data" class="vms-slot-top-actions" @click.stop>
      <button
        v-if="slot.type === 'camera'"
        class="vms-slot-action-btn vms-slot-expand-btn"
        :title="isExpanded ? 'Restaurar Grade' : 'Expandir Câmera (Resolução MAIN)'"
        @click.stop="handleExpand"
      >
        <svg v-if="isExpanded" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="4 14 10 14 10 20" /><polyline points="20 10 14 10 14 4" />
          <line x1="14" y1="10" x2="21" y2="3" /><line x1="3" y1="21" x2="10" y2="14" />
        </svg>
        <svg v-else width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="9 3 3 3 3 9" /><line x1="3" y1="3" x2="10" y2="10" />
          <polyline points="15 3 21 3 21 9" /><line x1="21" y1="3" x2="14" y2="10" />
          <polyline points="9 21 3 21 3 15" /><line x1="3" y1="21" x2="10" y2="14" />
          <polyline points="15 21 21 21 21 15" /><line x1="21" y1="21" x2="14" y2="14" />
        </svg>
      </button>
      <button v-if="!isExpanded" class="vms-slot-action-btn vms-slot-close-btn" title="Fechar Slot" @click.stop="emit('clear', slot.slot_index)">
        <svg width="9" height="9" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
      </button>
    </div>
    <template v-if="slot.type === 'camera' && slot.data">
      <div class="vms-slot-video" style="background: #000; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; position: relative; overflow: hidden;">
        <canvas v-show="!isPlayback" ref="canvasRef" style="width: 100%; height: 100%; object-fit: contain; display: block;"></canvas>
        <video v-if="isPlayback" ref="videoRef" playsinline muted autoplay style="width: 100%; height: 100%; object-fit: contain; display: block;" @ended="handleSeekOrSwitch"></video>
        <svg v-if="cameraDetections.length > 0" viewBox="0 0 100 100" preserveAspectRatio="none" style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; z-index: 10;">
          <rect
            v-for="b in cameraDetections"
            :key="b.id"
            :x="b.box[0]"
            :y="b.box[1]"
            :width="b.box[2]"
            :height="b.box[3]"
            fill="rgba(255, 94, 58, 0.06)"
            :stroke="b.color || '#ff5e3a'"
            stroke-width="0.8"
          />
        </svg>
        <div v-if="isStreamBuffering" class="vms-offline-sphere-container"><div class="vms-ubuntu-spinner"></div></div>
      </div>
      <div class="vms-slot-hud" @click.stop>
        <div class="vms-flex-row" style="gap: 0.35rem; margin-top: 1.1rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.9); color: #fff; font-size: 11px;">{{ (slot.data as CameraStreamInfo).name }}</span>
          <div class="vms-gear-wrapper" style="position: relative;" @mouseenter="isGearOpen = true" @mouseleave="isGearOpen = false">
            <button class="vms-gear-btn" title="Informacoes tecnicas" @click.stop="isGearOpen = !isGearOpen">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83-2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            </button>
            <div v-if="isGearOpen" class="vms-gear-popover">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">FLUXO: {{ effectiveProfile }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">RESOLUCAO: {{ effectiveProfile === 'MAIN' ? ((slot.data as CameraStreamInfo).resolution || '2560x1440') : '640x360' }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">MEDIA FPS: {{ (slot.data as CameraStreamInfo).fps || 25 }} FPS</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">DECODER: WebCodecs (HW)</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-muted);">PROTOCOLO: {{ ((slot.data as CameraStreamInfo).protocol || 'RTSP').toUpperCase() }}</span>
            </div>
          </div>
          <button class="vms-decoder-pill-btn" :title="'Alternar qualidade (SUB/MAIN)'" @click.stop="toggleProfile">[{{ effectiveProfile }}]</button>
        </div>
      </div>
    </template>
    <template v-else-if="slot.type === 'map' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #111419;"><span class="vms-badge vms-channel-pill-telegram">[MAPA INTERATIVO]</span><span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span></div>
    </template>
    <template v-else-if="slot.type === 'carousel' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #0c0f14;"><span class="vms-badge" style="background: rgba(255, 94, 58, 0.15); color: #ff5e3a !important; border: 1px solid rgba(255, 94, 58, 0.35);">[RONDA ATIVA // {{ (slot.data as any).interval_seconds }}S]</span><span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span></div>
    </template>
    <template v-else>
      <div class="vms-empty-slot-watermark">
        <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.35rem; user-select: none; opacity: 0.55; transition: opacity 0.2s;">
          <img :src="branding.headerLogo || '/hydra.svg'" alt="Logo" style="width: 26px; height: 26px; object-fit: contain; filter: grayscale(1) brightness(0.65);" />
          <span class="vms-text-mono vms-text-xs vms-font-bold" style="color: #94a3b8; letter-spacing: 0.8px; text-transform: uppercase;">
            {{ branding.companyName || branding.systemName || 'HYDRA VMS' }}
          </span>
          <span class="vms-text-mono vms-text-2xs" style="color: #64748b;">
            [ + ARRASTE UMA CÂMERA ]
          </span>
        </div>
      </div>
    </template>
  </div>
</template>

