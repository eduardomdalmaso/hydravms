<script setup lang="ts">
import { ref } from 'vue'
import type { PlaybackSpeed } from '../../types/mosaic'

defineProps<{ isPlaying: boolean; speed: PlaybackSpeed }>()

const emit = defineEmits<{
  (e: 'togglePlay'): void
  (e: 'setSpeed', speed: PlaybackSpeed): void
  (e: 'jumpSeconds', seconds: number): void
}>()

const isSpeedOpen = ref(false)
const selectSpeed = (s: PlaybackSpeed) => {
  emit('setSpeed', s)
  isSpeedOpen.value = false
}
</script>

<template>
  <div class="vms-vezha-controls-center">
    <button class="vms-jump-btn" title="Voltar 60s" @click="emit('jumpSeconds', -60)">&lt;&lt;&lt;</button>
    <button class="vms-jump-btn" title="Voltar 10s" @click="emit('jumpSeconds', -10)">&lt;&lt;</button>
    <button class="vms-jump-btn" title="Voltar 1s / Quadro" @click="emit('jumpSeconds', -1)">&lt;</button>

    <button
      class="vms-neu-circle vms-neu-accent-hero"
      style="width: 36px; height: 36px; font-size: 13px; cursor: pointer; margin: 0 0.25rem;"
      :title="isPlaying ? 'Pausar Reprodução' : 'Iniciar Reprodução'"
      @click="emit('togglePlay')"
    >
      <span v-if="isPlaying" style="letter-spacing: -1px; font-weight: 900; color: #ffffff;">❚❚</span>
      <span v-else style="margin-left: 2px; color: #ffffff;">▶</span>
    </button>

    <button class="vms-jump-btn" title="Avançar 1s / Quadro" @click="emit('jumpSeconds', 1)">&gt;</button>
    <button class="vms-jump-btn" title="Avançar 10s" @click="emit('jumpSeconds', 10)">&gt;&gt;</button>
    <button class="vms-jump-btn" title="Avançar 60s" @click="emit('jumpSeconds', 60)">&gt;&gt;&gt;</button>

    <!-- Speed Popup Trigger >> -->
    <div
      class="vms-speed-trigger-container"
      style="position: relative; margin-left: 0.4rem;"
      @mouseenter="isSpeedOpen = true"
      @mouseleave="isSpeedOpen = false"
    >
      <button
        class="vms-btn vms-btn-sm vms-btn-secondary"
        style="padding: 0.25rem 0.55rem; font-size: 10px; font-family: var(--vms-font-roboto); display: flex; align-items: center; gap: 0.25rem;"
        title="Alterar Velocidade de Reprodução"
        @click.stop="isSpeedOpen = !isSpeedOpen"
      >
        <span>{{ speed }}X</span>
        <span style="font-weight: 700; color: var(--vms-neu-accent-orange);">&gt;&gt;</span>
      </button>

      <!-- Speed Popover Above Trigger -->
      <div v-if="isSpeedOpen" class="vms-speed-popover" @mouseenter="isSpeedOpen = true">
        <button
          v-for="s in ([1, 2, 3] as const)"
          :key="s"
          class="vms-speed-option"
          :class="{ active: speed === s }"
          @click.stop="selectSpeed(s)"
        >
          {{ s }}X
        </button>
      </div>
    </div>
  </div>
</template>
