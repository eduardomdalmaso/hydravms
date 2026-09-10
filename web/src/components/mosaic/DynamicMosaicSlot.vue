<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount } from "vue"
import type { WorkspaceSlot, CameraStreamInfo } from "../../types/mosaic"
import { useWebRTCPlayer } from "../../composables/useWebRTCPlayer"
import { useTimelinePlayback } from "../../composables/useTimelinePlayback"
import SlotLinkedAlarm from "./SlotLinkedAlarm.vue"

const props = defineProps<{ slot: WorkspaceSlot; isActive?: boolean; isHero?: boolean }>()
const emit = defineEmits<{ (e: "selectCamera", cam: CameraStreamInfo): void; (e: "clear", idx: number): void }>()

const videoRef = ref<HTMLVideoElement | null>(null), isGearOpen = ref(false)
const decoderMode = ref<"MSE" | "H264">("MSE"), retryKey = ref(Date.now()), isImgLoading = ref(true)
const { isPlaying: isLivePlaying, start: startLive, stop: stopLive } = useWebRTCPlayer(videoRef)
const { isLive, currentTime, isPlaying, playbackSpeed, activePlaybackCameraId } = useTimelinePlayback()

const isPlayback = computed(() => {
  const cam = props.slot.type === 'camera' ? props.slot.data as CameraStreamInfo : null
  return !!cam && !isLive.value && activePlaybackCameraId.value === cam.id
})
const playbackSrc = computed(() => !isPlayback.value ? '' : `http://localhost:8083/api/v1/cameras/${(props.slot.data as CameraStreamInfo).id}/recordings/stream?t=${currentTime.value}`)

const syncStream = () => {
  isImgLoading.value = true
  if (isPlayback.value) {
    stopLive()
    if (videoRef.value && videoRef.value.src !== playbackSrc.value) {
      videoRef.value.src = playbackSrc.value; videoRef.value.playbackRate = playbackSpeed.value
      if (isPlaying.value) videoRef.value.play().catch(() => {})
      else videoRef.value.pause()
    }
  } else if (props.slot.type === 'camera' && props.slot.data && decoderMode.value === 'MSE') {
    if (videoRef.value && videoRef.value.src.includes('recordings')) videoRef.value.src = ''
    startLive((props.slot.data as CameraStreamInfo).id, props.isHero)
  } else { stopLive() }
}

const onImgError = () => { setTimeout(() => { retryKey.value = Date.now() }, 1500) }
const handleVis = () => { if (!document.hidden) { retryKey.value = Date.now(); syncStream() } }
document.addEventListener('visibilitychange', handleVis)
window.addEventListener('focus', handleVis)
onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', handleVis)
  window.removeEventListener('focus', handleVis)
})

watch(() => [props.slot.data, props.isHero, decoderMode.value, isLive.value, isPlayback.value, playbackSrc.value, isPlaying.value, playbackSpeed.value], syncStream, { deep: true })
onMounted(syncStream)
</script>

<template>
  <div class="vms-slot" :class="{ active: isActive, 'vms-slot-hero': isHero }" draggable="true" @click="slot.type === 'camera' && slot.data && emit('selectCamera', slot.data as CameraStreamInfo)">
    <button v-if="slot.data" class="vms-slot-close-btn" title="Fechar Slot" @click.stop="emit('clear', slot.slot_index)">
      <svg width="9" height="9" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
    </button>

    <!-- CAMERA TYPE -->
    <template v-if="slot.type === 'camera' && slot.data">
      <div class="vms-slot-video" style="background: #000000; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; position: relative; overflow: hidden;">
        <video v-show="decoderMode === 'MSE' && isPlaying" ref="videoRef" autoplay playsinline muted style="width: 100%; height: 100%; object-fit: contain; display: block;"></video>
        <img v-if="decoderMode === 'H264'" v-show="!isImgLoading" :src="`http://localhost:8080/api/v1/streams/${(slot.data as CameraStreamInfo).id}/mjpeg?k=${retryKey}`" alt="" style="width: 100%; height: 100%; object-fit: contain; display: block;" @load="isImgLoading = false" @error="onImgError" />
        <div v-if="(decoderMode === 'MSE' && !isPlaying) || (decoderMode === 'H264' && isImgLoading)" class="vms-offline-sphere-container"><div class="vms-ubuntu-spinner"></div></div>
      </div>

      <div class="vms-slot-hud">
        <div class="vms-flex-row" style="gap: 0.35rem; margin-top: 1.1rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.9); color: #fff; font-size: 11px;">{{ (slot.data as CameraStreamInfo).name }}</span>
          <div class="vms-gear-wrapper" style="position: relative;" @mouseenter="isGearOpen = true" @mouseleave="isGearOpen = false">
            <button class="vms-gear-btn" title="Informacoes tecnicas" @click.stop="isGearOpen = !isGearOpen">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            </button>
            <div v-if="isGearOpen" class="vms-gear-popover">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">CODEC: {{ (slot.data as CameraStreamInfo).codec || 'H.264' }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">RESOLUCAO: {{ (slot.data as CameraStreamInfo).resolution }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">MEDIA FPS: {{ (slot.data as CameraStreamInfo).fps }} FPS</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">DECODER: {{ decoderMode }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-muted);">PROTOCOLO: {{ ((slot.data as CameraStreamInfo).protocol || 'RTSP').toUpperCase() }}</span>
            </div>
          </div>
          <SlotLinkedAlarm :cameraId="(slot.data as CameraStreamInfo).id" />
          <button class="vms-decoder-pill-btn" title="Alternar decodificador" @click.stop="decoderMode = decoderMode === 'MSE' ? 'H264' : 'MSE'">[{{ decoderMode }}]</button>
        </div>
      </div>
    </template>

    <!-- MAP / CAROUSEL / EMPTY -->
    <template v-else-if="slot.type === 'map' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #111419;"><span class="vms-badge vms-channel-pill-telegram">[MAPA INTERATIVO]</span><span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span></div>
    </template>
    <template v-else-if="slot.type === 'carousel' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #0c0f14;"><span class="vms-badge" style="background: rgba(255, 94, 58, 0.15); color: #ff5e3a !important; border: 1px solid rgba(255, 94, 58, 0.35);">[RONDA ATIVA // {{ (slot.data as any).interval_seconds }}S]</span><span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span></div>
    </template>
    <template v-else>
      <div style="display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; user-select: none;"><span class="vms-text-mono vms-text-xs vms-text-dim">[ + ]</span></div>
    </template>
  </div>
</template>
