import { ref } from 'vue'
import type { PlaybackSpeed, TimelineSegment } from '../types/mosaic'

const isPlaying = ref(false)
const playbackSpeed = ref<PlaybackSpeed>(1)
const currentTime = ref(Date.now())
const isLive = ref(true)
const seekTrigger = ref(Date.now())
const activePlaybackCameraId = ref<string | null>(null)
const segments = ref<TimelineSegment[]>([])

export function useTimelinePlayback() {
  const togglePlay = () => {
    if (isLive.value) {
      isLive.value = false; isPlaying.value = false
    } else {
      isPlaying.value = !isPlaying.value
    }
  }

  const setSpeed = (speed: PlaybackSpeed) => { playbackSpeed.value = speed }

  const jumpSeconds = (seconds: number) => {
    isLive.value = false
    currentTime.value = Math.min(Date.now(), currentTime.value + seconds * 1000)
    seekTrigger.value = Date.now()
  }

  const stepFrame = (frames: number) => {
    isLive.value = false; isPlaying.value = false
    currentTime.value = Math.min(Date.now(), currentTime.value + (frames * 1000) / 30)
    seekTrigger.value = Date.now()
  }

  const seek = (time: number) => {
    currentTime.value = Math.min(Date.now(), time)
    isLive.value = (Date.now() - currentTime.value) < 1500
    seekTrigger.value = Date.now()
    if (isLive.value) isPlaying.value = true
  }

  const goToLive = () => {
    if (isLive.value) {
      isLive.value = false; isPlaying.value = false
    } else {
      currentTime.value = Date.now(); isLive.value = true
      isPlaying.value = true; playbackSpeed.value = 1
    }
  }

  return {
    isPlaying, playbackSpeed, currentTime, isLive, activePlaybackCameraId, seekTrigger, segments,
    togglePlay, setSpeed, jumpSeconds, stepFrame, goToLive, seek
  }
}

if (typeof window !== 'undefined') {
  setInterval(() => {
    if (isLive.value) {
      currentTime.value = Date.now()
    } else if (isPlaying.value) {
      currentTime.value = Math.min(Date.now(), currentTime.value + 1000 * playbackSpeed.value)
      if (currentTime.value >= Date.now()) isLive.value = true
    }
  }, 1000)
}
