import { ref, computed } from 'vue'
import type { GridLayout, CameraStreamInfo } from '../types/mosaic'

export function useMosaic() {
  const currentLayout = ref<GridLayout>('2x2')
  const selectedCamera = ref<CameraStreamInfo | null>(null)
  const isPlaybackDrawerOpen = ref(false)

  const cameras = ref<CameraStreamInfo[]>([
    {
      id: 'cam_01',
      name: 'Portaria Principal (Entrada)',
      location: 'Acesso A',
      status: 'recording',
      has_ptz: true,
      main_stream_url: 'webrtc://vms/live/cam_01_main',
      sub_stream_url: 'webrtc://vms/live/cam_01_sub',
      fps: 30,
      resolution: '1080P'
    },
    {
      id: 'cam_02',
      name: 'Estacionamento Visitantes',
      location: 'Pátio Externo',
      status: 'recording',
      has_ptz: false,
      main_stream_url: 'webrtc://vms/live/cam_02_main',
      sub_stream_url: 'webrtc://vms/live/cam_02_sub',
      fps: 25,
      resolution: '1080P'
    },
    {
      id: 'cam_03',
      name: 'Corredor de Cargas & Docas',
      location: 'Galpão 02',
      status: 'recording',
      has_ptz: false,
      main_stream_url: 'webrtc://vms/live/cam_03_main',
      sub_stream_url: 'webrtc://vms/live/cam_03_sub',
      fps: 30,
      resolution: '1080P'
    },
    {
      id: 'cam_04',
      name: 'Perímetro dos Fundos',
      location: 'Cerca Elétrica',
      status: 'recording',
      has_ptz: true,
      main_stream_url: 'webrtc://vms/live/cam_04_main',
      sub_stream_url: 'webrtc://vms/live/cam_04_sub',
      fps: 30,
      resolution: '4K'
    }
  ])

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
