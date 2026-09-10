import { ref, computed, onMounted } from 'vue'
import type { GridLayout, CameraStreamInfo } from '../types/mosaic'
import { fetchCameras } from '../services/api'

export function useMosaic() {
  const currentLayout = ref<GridLayout>('2x2')
  const selectedCamera = ref<CameraStreamInfo | null>(null)
  const isPlaybackDrawerOpen = ref(false)
  const cameras = ref<CameraStreamInfo[]>([])

  onMounted(async () => {
    const liveCams = await fetchCameras()
    if (liveCams.length > 0) {
      cameras.value = liveCams.map(c => ({
        id: c.id,
        name: c.name,
        location: c.location || c.ip || 'Rede Local',
        status: c.status === 'online' ? 'recording' : 'offline',
        has_ptz: c.has_ptz || false,
        main_stream_url: c.rtsp_url || '',
        sub_stream_url: c.rtsp_url || '',
        fps: c.fps || 30,
        resolution: c.resolution || '1080P'
      }))
    }
  })

  const maxSlots = computed(() => {
    switch (currentLayout.value) {
      case '1x1': return 1
      case '2x2': return 4
      case '3x3': return 9
      case '4x4': return 16
      case '1+5': return 6
    }
  })

  const selectCameraForPlayback = (cam: CameraStreamInfo) => {
    selectedCamera.value = cam
    isPlaybackDrawerOpen.value = true
  }

  const closePlaybackDrawer = () => {
    isPlaybackDrawerOpen.value = false
    selectedCamera.value = null
  }

  return {
    currentLayout,
    cameras,
    selectedCamera,
    isPlaybackDrawerOpen,
    maxSlots,
    selectCameraForPlayback,
    closePlaybackDrawer
  }
}
