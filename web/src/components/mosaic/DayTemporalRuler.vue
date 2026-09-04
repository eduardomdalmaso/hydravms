<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ currentTime: number }>()
const emit = defineEmits<{ (e: 'seek', timestamp: number): void }>()

const visibleDays = computed(() => {
  const days = []
  for (let offset = -4; offset <= 4; offset++) {
    const timestamp = props.currentTime + offset * 86400000
    const d = new Date(timestamp)
    days.push({
      offset,
      timestamp,
      dayNumber: String(d.getDate()).padStart(2, '0'),
      isActive: offset === 0,
      isFuture: timestamp > Date.now() + 86400000
    })
  }
  return days
})

const selectDay = (targetTimestamp: number) => {
  emit('seek', Math.min(Date.now(), targetTimestamp))
}

const handleDayWheel = (e: WheelEvent) => {
  const deltaDays = e.deltaY > 0 ? 1 : -1
  emit('seek', Math.min(Date.now(), props.currentTime + deltaDays * 86400000))
}
</script>

<template>
  <div
    class="vms-days-strip-container"
    title="Régua de Dias: Role o mouse ou clique para alternar o dia"
    @wheel.prevent="handleDayWheel"
  >
    <div class="vms-days-strip">
      <div
        v-for="d in visibleDays"
        :key="d.offset"
        class="vms-day-item"
        :class="{ active: d.isActive, disabled: d.isFuture }"
        @click="!d.isFuture && selectDay(d.timestamp)"
      >
        <span class="vms-day-item-number">{{ d.dayNumber }}</span>
        <div v-if="d.isActive" class="vms-day-active-bar"></div>
      </div>
    </div>
  </div>
</template>
