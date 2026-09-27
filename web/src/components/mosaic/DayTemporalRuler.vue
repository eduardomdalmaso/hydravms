<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentTime: number
  recordedRanges?: { start: number; end: number }[]
}>()
const emit = defineEmits<{ (e: 'seek', timestamp: number): void }>()

const months = ['JAN', 'FEV', 'MAR', 'ABR', 'MAI', 'JUN', 'JUL', 'AGO', 'SET', 'OUT', 'NOV', 'DEZ']

const recordedDays = computed(() => {
  if (!props.recordedRanges || props.recordedRanges.length === 0) return []

  const currentDayStr = new Date(props.currentTime).toDateString()
  const map = new Map<string, { key: string; timestamp: number; dayNumber: string; monthName: string; isActive: boolean }>()

  props.recordedRanges.forEach(r => {
    const d = new Date(r.start)
    const key = d.toDateString()
    if (!map.has(key)) {
      map.set(key, {
        key,
        timestamp: r.start,
        dayNumber: String(d.getDate()).padStart(2, '0'),
        monthName: months[d.getMonth()] || '',
        isActive: key === currentDayStr
      })
    }
  })

  return Array.from(map.values()).sort((a, b) => a.timestamp - b.timestamp)
})

const selectDay = (targetTimestamp: number) => {
  emit('seek', Math.min(Date.now(), targetTimestamp))
}
</script>

<template>
  <div class="vms-days-strip-container" title="Datas com Gravações Salvas no Sistema">
    <div v-if="recordedDays.length > 0" class="vms-days-strip">
      <div
        v-for="d in recordedDays"
        :key="d.key"
        class="vms-day-item"
        :class="{ active: d.isActive }"
        @click="selectDay(d.timestamp)"
      >
        <span class="vms-day-item-number">{{ d.dayNumber }}</span>
        <span class="vms-text-mono vms-text-2xs" style="font-size: 8px; color: var(--vms-text-dim); line-height: 1;">{{ d.monthName }}</span>
        <div v-if="d.isActive" class="vms-day-active-bar"></div>
      </div>
    </div>
    <div v-else class="vms-flex-center" style="color: var(--vms-text-dim); font-size: 11px; font-family: var(--vms-font-mono);">
      // NENHUMA DATA COM GRAVAÇÃO SALVA
    </div>
  </div>
</template>
