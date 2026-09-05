<script setup lang="ts">
import type { RecordingMode } from '../../../types/recordingSchedule'

const props = defineProps<{ schedule: boolean[][]; mode: RecordingMode }>()
const emit = defineEmits<{ (e: 'update:schedule', val: boolean[][]): void }>()

const days = ['DOM', 'SEG', 'TER', 'QUA', 'QUI', 'SEX', 'SAB']
const hours = Array.from({ length: 24 }, (_, i) => i)

const toggleCell = (d: number, h: number) => {
  const next = props.schedule.map((row, rIdx) => rIdx === d ? row.map((c, cIdx) => cIdx === h ? !c : c) : [...row])
  emit('update:schedule', next)
}

const toggleRow = (d: number) => {
  const allSet = props.schedule[d].every(Boolean)
  const next = props.schedule.map((row, rIdx) => rIdx === d ? Array(24).fill(!allSet) : [...row])
  emit('update:schedule', next)
}

const applyPreset = (preset: 'all' | 'work' | 'clear') => {
  const next = Array.from({ length: 7 }, (_, d) =>
    Array.from({ length: 24 }, (_, h) => {
      if (preset === 'all') return true
      if (preset === 'clear') return false
      return d >= 1 && d <= 5 && h >= 8 && h < 18
    })
  )
  emit('update:schedule', next)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.5rem; background: #07080c; border: 1px solid var(--vms-border); border-radius: 8px; padding: 0.75rem;">
    <div class="vms-flex-between">
      <span class="vms-text-mono vms-text-2xs vms-font-bold" style="color: #ffffff;">// GRADE DE HORARIOS (DIAS X 0-24H)</span>
      <div class="vms-flex-row" style="gap: 0.4rem;">
        <button type="button" class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 2px 8px;" @click="applyPreset('all')">24/7</button>
        <button type="button" class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 2px 8px;" @click="applyPreset('work')">08-18H</button>
        <button type="button" class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 2px 8px; color: #8b94a0 !important;" @click="applyPreset('clear')">LIMPAR</button>
      </div>
    </div>

    <!-- 24h Axis Header with All Numbers 0 to 23 in Bold White -->
    <div style="display: grid; grid-template-columns: 36px repeat(24, 1fr); gap: 2px; text-align: center;">
      <span class="vms-text-mono vms-text-2xs vms-font-bold" style="color: #ffffff;">DIA</span>
      <span v-for="h in hours" :key="h" class="vms-text-mono vms-font-bold" style="font-size: 9px; color: #ffffff;">{{ h }}</span>
    </div>

    <!-- Rows (DOM .. SAB) -->
    <div v-for="(day, dIdx) in days" :key="day" style="display: grid; grid-template-columns: 36px repeat(24, 1fr); gap: 2px; align-items: center;">
      <button type="button" class="vms-text-mono vms-text-2xs vms-font-bold" style="background: transparent; border: none; color: #ffffff; cursor: pointer; text-align: left; padding: 0;" @click="toggleRow(dIdx)">{{ day }}</button>
      <div
        v-for="h in hours"
        :key="h"
        :style="{
          background: schedule[dIdx]?.[h] ? '#ff5e3a' : 'rgba(255, 255, 255, 0.04)',
          border: schedule[dIdx]?.[h] ? '1px solid #ff7250' : '1px solid rgba(255, 255, 255, 0.06)',
          boxShadow: schedule[dIdx]?.[h] ? '0 0 6px rgba(255, 94, 58, 0.35)' : 'none'
        }"
        style="height: 14px; border-radius: 2px; cursor: pointer; transition: all 0.1s ease;"
        :title="`${day} ${h}:00 - ${h + 1}:00`"
        @click="toggleCell(dIdx, h)"
      />
    </div>
  </div>
</template>
