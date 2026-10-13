<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { CameraStreamInfo } from '../../types/mosaic'
import { useWSPCanvasPlayer } from '../../composables/useWSPCanvasPlayer'
import { getCameraMjpegUrl } from '../../utils/streamUrls'

const props = defineProps<{
  camera?: CameraStreamInfo
  isActive?: boolean
  isHero?: boolean
}>()

const emit = defineEmits<{ (e: 'select', cam: CameraStreamInfo): void }>()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const useMjpeg = ref(false)
const isImgLoading = ref(true)
const { start: startWsp, stop: stopWsp, error: wspError, isPlaying, isConnecting } = useWSPCanvasPlayer(canvasRef)

watch(wspError, (err) => {
  if (err) useMjpeg.value = true
})

const initStream = () => {
  if (props.camera) {
    useMjpeg.value = false
    isImgLoading.value = true
    const hasSub = props.camera.sub_stream_url !== undefined ? Boolean(props.camera.sub_stream_url) : true
    startWsp(props.camera.id, props.isHero, hasSub)
  } else {
    stopWsp()
  }
}

watch(() => [props.camera?.id, props.isHero], initStream)
onMounted(initStream)
</script>

<template>
  <div
    class="vms-slot"
    :class="{ active: isActive, 'vms-slot-hero': isHero }"
    @click="camera && emit('select', camera)"
  >
    <template v-if="camera">
      <div class="vms-slot-video" style="background: #000; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; overflow: hidden; position: relative;">
        <canvas
          v-show="!useMjpeg"
          ref="canvasRef"
          style="width: 100%; height: 100%; object-fit: contain; display: block;"
        ></canvas>
        <img
          v-if="useMjpeg"
          :src="getCameraMjpegUrl(camera.id, 4, 20)"
          alt="Camera Stream"
          style="width: 100%; height: 100%; object-fit: contain; display: block;"
          @load="isImgLoading = false"
        />
        <div v-if="(!useMjpeg && (!isPlaying || isConnecting)) || (useMjpeg && isImgLoading)" class="vms-offline-sphere-container">
          <div class="vms-ubuntu-spinner"></div>
        </div>
      </div>
      <div class="vms-slot-hud">
        <div class="vms-flex-row" style="gap: 0.5rem;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.85); color: #fff; font-size: 11px;">
            {{ camera.name }}
          </span>
          <span v-if="camera.has_ptz" class="vms-badge vms-badge-info" style="font-size: 10px;">[PTZ]</span>
        </div>
        <div class="vms-flex-between">
          <span class="vms-text-mono vms-text-xs vms-text-dim" style="background: rgba(0,0,0,0.6); padding: 2px 4px; border-radius: 2px;">
            [{{ useMjpeg ? 'MJPEG' : 'H264' }}]
          </span>
          <span class="vms-badge vms-badge-recording" style="font-size: 10px;">[REC]</span>
        </div>
      </div>
    </template>
    <template v-else>
      <span class="vms-text-xs vms-text-dim">[SLOT VAZIO]</span>
    </template>
  </div>
</template>
