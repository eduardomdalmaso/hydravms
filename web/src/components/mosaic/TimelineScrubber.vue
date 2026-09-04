<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TimelineSegment } from '../../types/mosaic'
import TimelineRuler from './TimelineRuler.vue'

const props = defineProps<{ currentTime: number; segments: TimelineSegment[] }>()
const emit = defineEmits<{ (e: 'seek', time: number): void }>()

const zoomMinutes = ref(60)
const isHovering = ref(false)
const hoverX = ref(0)
const hoverInfo = ref('')

const windowStart = computed(() => props.currentTime - (zoomMinutes.value * 60000) / 2)
const windowEnd = computed(() => props.currentTime + (zoomMinutes.value * 60000) / 2)

// Sample discrete 1-min chunks with realistic gaps
const visibleChunks = computed(() => {
  const base = props.currentTime - 1800000
  const chunks = []
  for (let i = 0; i < 40; i++) {
    if (i === 12 || i === 13 || i === 25) continue // Intended GAP
    const start = base + i * 60000
    const left = Math.max(0, ((start - windowStart.value) / (windowEnd.value - windowStart.value)) * 100)
    const width = ((60000) / (windowEnd.value - windowStart.value)) * 100
    if (left < 100 && left + width > 0) chunks.push({ id: `c-${i}`, left, width: Math.max(0.2, width) })
  }
  return chunks
})

const handleWheel = (e: WheelEvent) => {
  if (e.deltaY < 0) zoomMinutes.value = Math.max(1, Math.round(zoomMinutes.value / 2))
  else zoomMinutes.value = Math.min(1440, Math.round(zoomMinutes.value * 2))
}

const handleMouseMove = (e: MouseEvent) => {
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  hoverX.value = e.clientX - rect.left
  const ratio = hoverX.value / rect.width
  const targetDate = new Date(windowStart.value + ratio * (windowEnd.value - windowStart.value))
  const isRecorded = visibleChunks.value.some(c => ratio * 100 >= c.left && ratio * 100 <= c.left + c.width)
  hoverInfo.value = `${targetDate.toLocaleTimeString('pt-BR')} ${isRecorded ? '[GRAVADO // 60S]' : '[GAP // SEM GRAVACAO]'}`
  isHovering.value = true
}

const handleTimelineClick = (e: MouseEvent) => {
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  const ratio = (e.clientX - rect.left) / rect.width
  emit('seek', windowStart.value + ratio * (windowEnd.value - windowStart.value))
}
</script>

<template>
  <div class="vms-timeline-wrapper" @wheel.prevent="handleWheel" @mousemove="handleMouseMove" @mouseleave="isHovering = false" @click="handleTimelineClick">
    <div v-if="isHovering" class="vms-timeline-tooltip" :style="{ left: `${hoverX}px` }">
      {{ hoverInfo }} [ZOOM: {{ zoomMinutes >= 60 ? `${Math.round(zoomMinutes / 60)}H` : `${zoomMinutes} MIN` }}]
    </div>
    <div class="vms-timeline-track">
      <div v-for="c in visibleChunks" :key="c.id" class="vms-timeline-recorded-bar" :style="{ left: `${c.left}%`, width: `${c.width}%` }"></div>
      <div class="vms-timeline-needle" style="left: 50%;"></div>
    </div>
    <TimelineRuler :zoomMinutes="zoomMinutes" :centerTime="currentTime" />
  </div>
</template>
