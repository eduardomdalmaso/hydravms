<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount } from "vue"
import type { WorkspaceSlot, CameraStreamInfo } from "../../types/mosaic"
import { useWSPCanvasPlayer } from "../../composables/useWSPCanvasPlayer"
import { useTimelinePlayback } from "../../composables/useTimelinePlayback"
import { useBranding } from "../../composables/useBranding"
import { getCameraMjpegUrl } from "../../utils/streamUrls"
import { useLiveDetections } from "../../composables/useLiveDetections"

const props = defineProps<{ slot: WorkspaceSlot; isActive?: boolean; isHero?: boolean }>()
const emit = defineEmits<{ (e: "selectCamera", cam: CameraStreamInfo): void; (e: "clear", idx: number): void }>()
const { branding } = useBranding()
const { getDetections } = useLiveDetections()
const canvasRef = ref<HTMLCanvasElement | null>(null), imgRef = ref<HTMLImageElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null), isGearOpen = ref(false)
const savedMode = (typeof localStorage !== 'undefined' && localStorage.getItem('vms_decoder_mode')) as "H264" | "MJPEG" | null
const decoderMode = ref<"H264" | "MJPEG">(savedMode || "MJPEG"), retryKey = ref(Date.now()), isImgLoading = ref(true)
const { start: startWsp, stop: stopWsp, isPlaying: isWspPlaying, isConnecting: isWspConnecting } = useWSPCanvasPlayer(canvasRef)
const { isLive, currentTime, isPlaying, playbackSpeed, activePlaybackCameraId, seekTrigger } = useTimelinePlayback()

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
  if (decoderMode.value === 'MJPEG') return isImgLoading.value
  return !isWspPlaying.value || isWspConnecting.value
})

const checkMjpegFrame = () => {
  let attempts = 0
  const poll = () => {
    attempts++
    if (imgRef.value && imgRef.value.naturalWidth > 0) {
      isImgLoading.value = false
      return
    }
    if (attempts < 25) setTimeout(poll, 150)
    else isImgLoading.value = false
  }
  setTimeout(poll, 200)
}

const toggleDecoder = () => {
  if (decoderMode.value === 'MJPEG' && imgRef.value) {
    imgRef.value.src = ''
  }
  decoderMode.value = decoderMode.value === 'H264' ? 'MJPEG' : 'H264'
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('vms_decoder_mode', decoderMode.value)
  }
}

const handleSeekOrSwitch = () => {
  if (isPlayback.value && videoRef.value) {
    stopWsp()
    const cam = props.slot.data as CameraStreamInfo
    videoRef.value.src = `http://localhost:8083/api/v1/cameras/${cam.id}/recordings/stream?t=${currentTime.value}`
    videoRef.value.playbackRate = playbackSpeed.value
    if (isPlaying.value) videoRef.value.play().catch(() => {})
    else videoRef.value.pause()
  } else if (!isPlayback.value && props.slot.type === 'camera' && props.slot.data) {
    const cam = props.slot.data as CameraStreamInfo
    const hasSub = cam.protocol !== 'RTMP' && Boolean(cam.sub_stream_url)
    if (decoderMode.value === 'H264') {
      if (imgRef.value) imgRef.value.src = ''
      startWsp(cam.id, props.isHero, hasSub)
    } else {
      stopWsp()
      isImgLoading.value = true
      retryKey.value = Date.now()
      checkMjpegFrame()
    }
  }
}

const syncStream = () => {
  isImgLoading.value = true
  if (isPlayback.value && videoRef.value) {
    videoRef.value.playbackRate = playbackSpeed.value
    if (isPlaying.value && videoRef.value.paused) videoRef.value.play().catch(() => {})
    else if (!isPlaying.value && !videoRef.value.paused) videoRef.value.pause()
  } else { handleSeekOrSwitch() }
}

const retryImg = () => setTimeout(() => { retryKey.value = Date.now(); checkMjpegFrame() }, 1500)
const handleVis = () => { if (!document.hidden) { retryKey.value = Date.now(); handleSeekOrSwitch() } }
document.addEventListener('visibilitychange', handleVis)
window.addEventListener('focus', handleVis)
onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', handleVis)
  window.removeEventListener('focus', handleVis)
  if (imgRef.value) imgRef.value.src = ''
  stopWsp()
})

watch(() => [isPlayback.value, seekTrigger.value], handleSeekOrSwitch)
watch(() => [props.slot.data, props.isHero, decoderMode.value, isPlaying.value, playbackSpeed.value], syncStream, { deep: true })
onMounted(() => { handleSeekOrSwitch(); if (decoderMode.value === 'MJPEG') checkMjpegFrame() })
</script>

<template>
  <div class="vms-slot" :class="{ active: isActive, 'vms-slot-hero': isHero }" draggable="true" @click="slot.type === 'camera' && slot.data && emit('selectCamera', slot.data as CameraStreamInfo)">
    <button v-if="slot.data" class="vms-slot-close-btn" title="Fechar Slot" @click.stop="emit('clear', slot.slot_index)">
      <svg width="9" height="9" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
    </button>
    <template v-if="slot.type === 'camera' && slot.data">
      <div class="vms-slot-video" style="background: #000; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; position: relative; overflow: hidden;">
        <canvas v-show="!isPlayback && decoderMode === 'H264'" ref="canvasRef" style="width: 100%; height: 100%; object-fit: contain; display: block;"></canvas>
        <img v-if="!isPlayback && decoderMode === 'MJPEG'" ref="imgRef" :src="`${getCameraMjpegUrl((slot.data as CameraStreamInfo).id, 4, 20)}&k=${retryKey}`" alt="" style="width: 100%; height: 100%; object-fit: contain; display: block;" @load="isImgLoading = false" @error="retryImg" />
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
      <div class="vms-slot-hud">
        <div class="vms-flex-row" style="gap: 0.35rem; margin-top: 1.1rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.9); color: #fff; font-size: 11px;">{{ (slot.data as CameraStreamInfo).name }}</span>
          <div class="vms-gear-wrapper" style="position: relative;" @mouseenter="isGearOpen = true" @mouseleave="isGearOpen = false">
            <button class="vms-gear-btn" title="Informacoes tecnicas" @click.stop="isGearOpen = !isGearOpen">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83-2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            </button>
            <div v-if="isGearOpen" class="vms-gear-popover">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">CODEC: {{ (slot.data as CameraStreamInfo).codec || 'H.264' }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">RESOLUCAO: {{ (slot.data as CameraStreamInfo).resolution }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">MEDIA FPS: {{ (slot.data as CameraStreamInfo).fps }} FPS</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">DECODER: {{ decoderMode }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-muted);">PROTOCOLO: {{ ((slot.data as CameraStreamInfo).protocol || 'RTSP').toUpperCase() }}</span>
            </div>
          </div>
          <button class="vms-decoder-pill-btn" :title="'Alternar decodificador (' + decoderMode + ')'" @click.stop="toggleDecoder">[{{ decoderMode }}]</button>
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
