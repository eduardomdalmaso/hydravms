import { ref } from 'vue'
import type { PlaybackSpeed, TimelineSegment } from '../types/mosaic'

export function useTimelinePlayback() {
  const isPlaying = ref(false)
  const playbackSpeed = ref<PlaybackSpeed>(1)
  const currentTime = ref(Date.now())
  const isLive = ref(true)

  const segments = ref<TimelineSegment[]>([])

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
