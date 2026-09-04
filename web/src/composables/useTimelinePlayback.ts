import { ref } from 'vue'
import type { PlaybackSpeed, TimelineSegment } from '../types/mosaic'

export function useTimelinePlayback() {
  const isPlaying = ref(true)
  const playbackSpeed = ref<PlaybackSpeed>(1)
  const currentTime = ref(Date.now() - 3600000) // 1 hour ago
  const isLive = ref(false)

  const segments = ref<TimelineSegment[]>([
    { id: 's1', start_time: Date.now() - 86400000, end_time: Date.now(), type: 'continuous' },
    { id: 's2', start_time: Date.now() - 7200000, end_time: Date.now() - 6600000, type: 'motion', label: 'Movimento' },
    { id: 's3', start_time: Date.now() - 3600000, end_time: Date.now() - 3400000, type: 'ai_alert', label: 'Invasão IA' }
  ])

  const togglePlay = () => {
    isPlaying.value = !isPlaying.value
  }

  const setSpeed = (speed: PlaybackSpeed) => {
    playbackSpeed.value = speed
  }

  const jumpSeconds = (seconds: number) => {
    currentTime.value += seconds * 1000
    isLive.value = false
  }

  const stepFrame = (frames: number) => {
    isPlaying.value = false
    currentTime.value += (frames * 1000) / 30 // ~33ms per frame @ 30fps
  }

  const goToLive = () => {
    currentTime.value = Date.now()
    isLive.value = true
    isPlaying.value = true
    playbackSpeed.value = 1
  }

  return {
    isPlaying,
    playbackSpeed,
    currentTime,
    isLive,
    segments,
    togglePlay,
    setSpeed,
    jumpSeconds,
    stepFrame,
    goToLive
  }
}
