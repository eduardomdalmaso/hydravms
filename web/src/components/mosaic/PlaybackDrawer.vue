<script setup lang="ts">
import { ref, computed } from 'vue'
import type { CameraStreamInfo } from '../../types/mosaic'
import { useTimelinePlayback } from '../../composables/useTimelinePlayback'
import VezhaCanvasTimeline from './VezhaCanvasTimeline.vue'
import DayTemporalRuler from './DayTemporalRuler.vue'
import TimelineFooterLeft from './TimelineFooterLeft.vue'
import TimelineFooterCenter from './TimelineFooterCenter.vue'
import TimelineFooterRight from './TimelineFooterRight.vue'

defineProps<{ camera: CameraStreamInfo }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'exportClip'): void }>()

const {
  isPlaying, playbackSpeed, currentTime, isLive,
  togglePlay, setSpeed, jumpSeconds, goToLive
} = useTimelinePlayback()

const isExportMode = ref(false)
const isSyncActive = ref(false)
const exportRange = ref<{ start: number; end: number }>({
  start: currentTime.value - 300000,
  end: currentTime.value + 300000
})

const exportDurationText = computed(() => {
  const diffSec = Math.max(0, Math.round((exportRange.value.end - exportRange.value.start) / 1000))
  const m = Math.floor(diffSec / 60), s = diffSec % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')} MIN`
})

const toggleExportMode = () => {
  isExportMode.value = !isExportMode.value
  if (isExportMode.value) {
    exportRange.value = { start: currentTime.value - 300000, end: currentTime.value + 300000 }
  }
}

const notify = (msg: string) => {
  if (typeof window !== 'undefined') window.alert(msg)
}
</script>

<template>
  <div class="vms-vezha-timeline">
    <!-- Top Bar: Pure Day Temporal Ruler + Ubuntu Close Button -->
    <div class="vms-drawer-header" style="padding: 0.35rem 0.85rem; display: flex; align-items: center; gap: 0.75rem;">
      <DayTemporalRuler :currentTime="currentTime" @seek="currentTime = $event" />
      <button class="vms-ubuntu-close-btn" style="position: relative; top: 0; right: 0;" title="Fechar Reprodução" @click="emit('close')">
        <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
      </button>
    </div>

    <!-- HTML5 Canvas 24-Hour Timeline with Dual-Bar Export Range Selection -->
    <VezhaCanvasTimeline
      :currentTime="currentTime"
      :isExportMode="isExportMode"
      :exportStart="exportRange.start"
      :exportEnd="exportRange.end"
      @seek="currentTime = $event"
      @updateExportRange="exportRange = $event"
    />

    <!-- Vezha 3-Zone Footer Controls -->
    <div class="vms-vezha-footer">
      <TimelineFooterLeft
        :isExportMode="isExportMode"
        :exportDurationText="exportDurationText"
      />
      <TimelineFooterCenter
        :isPlaying="isPlaying" :speed="playbackSpeed"
        @togglePlay="togglePlay" @setSpeed="setSpeed" @jumpSeconds="jumpSeconds"
      />
      <TimelineFooterRight
        :isLive="isLive"
        :isExportMode="isExportMode"
        :isSyncActive="isSyncActive"
        @showLayers="notify('Filtro de Camadas: Contínuo, Movimento e IA Ativos')"
        @toggleExport="toggleExportMode"
        @toggleSync="isSyncActive = !isSyncActive"
        @goLive="goToLive"
      />
    </div>
  </div>
</template>
