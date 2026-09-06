<script setup lang="ts">
import type { RondaStreamItem } from '../../../types/rondaTree'

const props = defineProps<{
  stream: RondaStreamItem
  index: number
  isFirst: boolean
  isLast: boolean
}>()

const emit = defineEmits<{
  (e: 'updateInterval', val: number): void
  (e: 'moveUp'): void
  (e: 'moveDown'): void
  (e: 'remove'): void
}>()

const quickIntervals = [5, 10, 15, 30]
</script>

<template>
  <div
    class="vms-flex-between"
    style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.5rem 0.75rem; align-items: center;"
  >
    <!-- Left: Index + Camera Info -->
    <div class="vms-flex-row" style="gap: 0.6rem; align-items: center; min-width: 0;">
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border); font-size: 8.5px; font-weight: 700;">
        [{{ String(index + 1).padStart(2, '0') }}]
      </span>

      <div class="vms-flex-col" style="gap: 1px; min-width: 0;">
        <span class="vms-font-semibold vms-text-xs" style="color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {{ stream.cameraName }}
        </span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          {{ stream.resolution }} @ {{ stream.fps }} FPS // ID: {{ stream.cameraId }}
        </span>
      </div>
    </div>

    <!-- Right: Individual Interval Stepper + Actions -->
    <div class="vms-flex-row" style="gap: 0.5rem; align-items: center; flex-shrink: 0;">
      <!-- Quick Interval Pills -->
      <div class="vms-flex-row" style="gap: 0.2rem;">
        <button
          v-for="sec in quickIntervals"
          :key="sec"
          class="vms-btn vms-btn-sm"
          :class="stream.intervalSeconds === sec ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="font-size: 8.5px; padding: 2px 5px; min-width: 24px;"
          @click="emit('updateInterval', sec)"
        >
          {{ sec }}s
        </button>
      </div>

      <!-- Custom Number Input -->
      <div class="vms-flex-row" style="align-items: center; gap: 2px;">
        <input
          :value="stream.intervalSeconds"
          type="number"
          min="1"
          max="300"
          class="vms-auth-input"
          style="width: 44px; font-size: 10.5px; padding: 2px 4px; text-align: center; height: 24px;"
          @input="emit('updateInterval', Number(($event.target as HTMLInputElement).value) || 5)"
        />
        <span class="vms-text-mono vms-text-2xs vms-text-dim">s</span>
      </div>

      <!-- Move Up / Down Reorder -->
      <div class="vms-flex-row" style="gap: 1px;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" :disabled="isFirst" style="padding: 2px 5px;" title="Mover para cima" @click="emit('moveUp')">
          <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m18 15-6-6-6 6"/></svg>
        </button>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" :disabled="isLast" style="padding: 2px 5px;" title="Mover para baixo" @click="emit('moveDown')">
          <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m6 9 6 6 6-6"/></svg>
        </button>
      </div>

      <!-- Remove button -->
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 6px; color: #ff003c;" title="Remover da ronda" @click="emit('remove')">
        ×
      </button>
    </div>
  </div>
</template>
