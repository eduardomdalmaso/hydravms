<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  zoomMinutes: number
  centerTime: number
}>()

const ticks = computed(() => {
  const result: string[] = []
  const count = 7
  const startMs = props.centerTime - (props.zoomMinutes * 60000) / 2
  const stepMs = (props.zoomMinutes * 60000) / (count - 1)

  for (let i = 0; i < count; i++) {
    const d = new Date(startMs + i * stepMs)
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    const ss = String(d.getSeconds()).padStart(2, '0')

    if (props.zoomMinutes <= 5) {
      result.push(`${hh}:${mm}:${ss}`)
    } else {
      result.push(`${hh}:${mm}`)
    }
  }
  return result
})
</script>

<template>
  <div class="vms-timeline-ruler">
    <span v-for="(tick, idx) in ticks" :key="idx" class="vms-timeline-tick">
      {{ tick }}
    </span>
  </div>
</template>
