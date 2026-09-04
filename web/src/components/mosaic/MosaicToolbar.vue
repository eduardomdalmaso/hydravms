<script setup lang="ts">
import type { GridLayout } from '../../types/mosaic'

const props = defineProps<{ layout: GridLayout; totalCameras?: number }>()
const emit = defineEmits<{ (e: 'update:layout', val: GridLayout): void; (e: 'fullscreen'): void }>()

const layouts: { id: GridLayout; label: string }[] = [
  { id: 'auto', label: 'AUTO' },
  { id: '1x1', label: '1x1' },
  { id: '1x2', label: '1x2' },
  { id: '2x2', label: '2x2' },
  { id: '1+5', label: '1+5' },
  { id: '3x3', label: '3x3' },
  { id: '4x4', label: '4x4' }
]
</script>

<template>
  <div class="vms-flex-between" style="padding: 0.35rem 0.75rem; background-color: var(--vms-bg-surface); border-bottom: 1px solid var(--vms-border);">
    <div class="vms-speed-selector">
      <button
        v-for="l in layouts"
        :key="l.id"
        class="vms-btn vms-btn-ghost vms-btn-sm"
        :style="{
          backgroundColor: layout === l.id ? 'var(--vms-primary)' : 'transparent',
          color: layout === l.id ? '#ffffff' : 'var(--vms-text-muted)',
          fontWeight: layout === l.id ? '600' : '400'
        }"
        @click="emit('update:layout', l.id)"
      >
        {{ l.label }}
      </button>
    </div>
    <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('fullscreen')">[TELA CHEIA]</button>
  </div>
</template>
