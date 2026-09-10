import { ref } from 'vue'
import type { PlaybackSpeed, TimelineSegment } from '../types/mosaic'

const isPlaying = ref(false)
const playbackSpeed = ref<PlaybackSpeed>(1)
const currentTime = ref(Date.now())
const isLive = ref(true)
const activePlaybackCameraId = ref<string | null>(null)
const segments = ref<TimelineSegment[]>([])

export function useTimelinePlayback() {

  const togglePlay = () => {
    if (isLive.value) {
      isLive.value = false
      isPlaying.value = false
    } else {
      isPlaying.value = !isPlaying.value
    }
  }

  const setSpeed = (speed: PlaybackSpeed) => {
    playbackSpeed.value = speed
  }

  const jumpSeconds = (seconds: number) => {
    isLive.value = false
    currentTime.value = Math.min(Date.now(), currentTime.value + seconds * 1000)
    isPlaying.value = true
  }

  const stepFrame = (frames: number) => {
    isLive.value = false
    isPlaying.value = false
    currentTime.value = Math.min(Date.now(), currentTime.value + (frames * 1000) / 30)
  }

  const seek = (time: number) => {
    currentTime.value = Math.min(Date.now(), time)
    isLive.value = Math.abs(currentTime.value - Date.now()) < 3000
    if (!isLive.value) isPlaying.value = true
  }

  const goToLive = () => {
    if (isLive.value) {
      isLive.value = false
      isPlaying.value = false
    } else {
      currentTime.value = Date.now()
      isLive.value = true
      isPlaying.value = true
      playbackSpeed.value = 1
    }
  }

  let ticker: any = null
  if (typeof window !== 'undefined') {
    ticker = setInterval(() => {
      if (isLive.value) {
        currentTime.value = Date.now()
      } else if (isPlaying.value) {
        currentTime.value = Math.min(Date.now(), currentTime.value + 1000 * playbackSpeed.value)
        if (currentTime.value >= Date.now()) isLive.value = true
      }
    }, 1000)
  }

  return {
    isPlaying, playbackSpeed, currentTime, isLive, activePlaybackCameraId, segments,
    togglePlay, setSpeed, jumpSeconds, stepFrame, goToLive, seek
  }
}
