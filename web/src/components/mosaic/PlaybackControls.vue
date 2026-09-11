<script setup lang="ts">
import { useI18n } from '../../composables/useI18n'
import type { PlaybackSpeed } from '../../types/mosaic'

defineProps<{
  isPlaying: boolean
  speed: PlaybackSpeed
  isLive: boolean
  isLocked?: boolean
}>()

const emit = defineEmits<{
  (e: 'togglePlay'): void
  (e: 'setSpeed', speed: PlaybackSpeed): void
  (e: 'jumpSeconds', seconds: number): void
  (e: 'toggleLock'): void
  (e: 'exportClip'): void
}>()
const { t } = useI18n()

const speedOptions: { speed: PlaybackSpeed; label: string }[] = [
  { speed: 0.5, label: '0.5X' },
  { speed: 1, label: '1X' },
  { speed: 2, label: '2X' },
  { speed: 3, label: '3X' }
]
</script>

<template>
  <div class="vms-timeline-controls">
    <div class="vms-flex-row" style="gap: 0.25rem;">
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Voltar 60s" @click="emit('jumpSeconds', -60)">&lt;&lt;&lt;</button>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Voltar 10s" @click="emit('jumpSeconds', -10)">&lt;&lt;</button>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Voltar 1s" @click="emit('jumpSeconds', -1)">&lt;</button>
    </div>

    <button
      class="vms-btn vms-btn-sm"
      :class="isPlaying ? 'vms-btn-primary' : 'vms-btn-secondary'"
      style="padding: 0.4rem 1.25rem; font-weight: 700;"
      @click="emit('togglePlay')"
    >
      {{ isPlaying ? t('pause') : t('play') }}
    </button>

    <div class="vms-flex-row" style="gap: 0.25rem;">
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Avançar 1s" @click="emit('jumpSeconds', 1)">&gt;</button>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Avançar 10s" @click="emit('jumpSeconds', 10)">&gt;&gt;</button>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" title="Avançar 60s" @click="emit('jumpSeconds', 60)">&gt;&gt;&gt;</button>
    </div>

    <div class="vms-flex-row" style="gap: 0.25rem; margin-left: 0.5rem;">
      <button
        v-for="opt in speedOptions"
        :key="opt.speed"
        class="vms-btn vms-btn-sm"
        :class="speed === opt.speed ? 'vms-btn-primary' : 'vms-btn-ghost'"
        style="padding: 0.2rem 0.4rem; font-size: 10px;"
        @click="emit('setSpeed', opt.speed)"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>
