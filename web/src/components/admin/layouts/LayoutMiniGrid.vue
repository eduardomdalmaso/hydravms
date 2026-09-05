<script setup lang="ts">
import { ref, computed } from 'vue'
import type { GridLayout } from '../../../types/mosaic'
import type { LayoutSlotItem } from '../../../types/layoutTree'

const props = defineProps<{ grid: GridLayout; slots: LayoutSlotItem[] }>()
const emit = defineEmits<{ (e: 'dropCamera', slotIdx: number, cam: any): void; (e: 'clearSlot', slotIdx: number): void }>()

const dragOverSlot = ref<number | null>(null)

const slotCount = computed(() => {
  if (props.grid === '1x1') return 1
  if (props.grid === '1x2') return 2
  if (props.grid === '2x2') return 4
  if (props.grid === '1+5') return 6
  if (props.grid === '3x3') return 9
  return 16
})

const isHero = (slotIdx: number) => props.grid === '1+5' && slotIdx === 0

const gridTemplateStyle = computed(() => {
  if (props.grid === '1x1') return { gridTemplateColumns: '1fr', maxWidth: '480px', aspectRatio: '16 / 9' }
  if (props.grid === '1x2') return { gridTemplateColumns: 'repeat(2, 1fr)', maxWidth: '540px', aspectRatio: '16 / 9' }
  if (props.grid === '2x2') return { gridTemplateColumns: 'repeat(2, 1fr)', gridTemplateRows: 'repeat(2, 1fr)', maxWidth: '540px', aspectRatio: '16 / 10' }
  if (props.grid === '1+5') return { gridTemplateColumns: 'repeat(3, 1fr)', gridTemplateRows: 'repeat(3, 1fr)', maxWidth: '580px', aspectRatio: '16 / 10' }
  if (props.grid === '3x3') return { gridTemplateColumns: 'repeat(3, 1fr)', gridTemplateRows: 'repeat(3, 1fr)', maxWidth: '580px', aspectRatio: '16 / 10' }
  return { gridTemplateColumns: 'repeat(4, 1fr)', gridTemplateRows: 'repeat(4, 1fr)', maxWidth: '600px', aspectRatio: '16 / 10' }
})

const getCameraInSlot = (slotIdx: number) => props.slots.find(s => s.slotIndex === slotIdx)

const onDrop = (slotIdx: number, ev: DragEvent) => {
  ev.preventDefault(); dragOverSlot.value = null
  const raw = ev.dataTransfer?.getData('application/json')
  if (raw) { try { const cam = JSON.parse(raw); emit('dropCamera', slotIdx, cam) } catch {} }
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.5rem; width: 100%;">
    <div class="vms-flex-between" style="align-items: center;">
      <span class="vms-text-mono vms-text-xs vms-font-bold" style="color: var(--vms-neu-accent-orange);">
        // PREVIA DO MOSAICO // GRADE {{ grid }} (SOLTE AS CAMERAS NOS SLOTS)
      </span>
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border); font-size: 8.5px;">GRADE AMPLIADA</span>
    </div>

    <!-- Generous prominent non-stretched centered mini-grid -->
    <div
      class="vms-grid-container"
      :style="{ ...gridTemplateStyle, width: '100%', margin: '0 auto', gap: '6px', background: '#07080c', padding: '8px', borderRadius: '8px', border: '1px solid var(--vms-border)' }"
    >
      <div
        v-for="idx in slotCount"
        :key="idx"
        class="vms-flex-col"
        :style="{
          gridColumn: isHero(idx - 1) ? '1 / span 2' : undefined,
          gridRow: isHero(idx - 1) ? '1 / span 2' : undefined,
          border: dragOverSlot === idx - 1 ? '1px dashed #ff5e3a' : '1px dashed rgba(255, 255, 255, 0.15)',
          background: dragOverSlot === idx - 1 ? 'rgba(255, 94, 58, 0.12)' : getCameraInSlot(idx - 1)?.cameraId ? '#0e121a' : 'rgba(255, 255, 255, 0.02)',
          borderRadius: '6px', position: 'relative', overflow: 'hidden', justifyContent: 'center', alignItems: 'center', minWidth: 0, minHeight: 0
        }"
        @dragover.prevent="dragOverSlot = idx - 1"
        @dragleave="dragOverSlot = null"
        @drop="onDrop(idx - 1, $event)"
      >
        <template v-if="getCameraInSlot(idx - 1)?.cameraId">
          <div style="position: absolute; top: 4px; left: 5px; display: flex; gap: 4px; align-items: center;">
            <span class="vms-badge vms-badge-success" style="font-size: 7px; padding: 1px 3px;">[LIVE]</span>
            <span class="vms-text-mono vms-text-2xs" style="color: #ff5e3a; font-size: 8px;">{{ isHero(idx - 1) ? 'HERO' : `S${idx}` }}</span>
          </div>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="position: absolute; top: 2px; right: 2px; padding: 1px 5px; font-size: 10px; color: #ff5e3a;" title="Remover" @click.stop="emit('clearSlot', idx - 1)">
            ×
          </button>
          <svg :width="isHero(idx - 1) ? '34' : (grid === '4x4' ? '18' : '26')" :height="isHero(idx - 1) ? '34' : (grid === '4x4' ? '18' : '26')" viewBox="0 0 576 512" fill="#ff5e3a" style="opacity: 0.85; margin-top: 4px;">
            <path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/>
          </svg>
          <span class="vms-font-bold" :style="{ fontSize: isHero(idx - 1) ? '12px' : (grid === '4x4' ? '8.5px' : '10px') }" style="color: #fff; margin-top: 2px; max-width: 90%; text-align: center; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            {{ getCameraInSlot(idx - 1)?.cameraName }}
          </span>
        </template>
        <template v-else>
          <span class="vms-text-mono vms-text-2xs vms-text-dim" :style="{ fontSize: isHero(idx - 1) ? '11px' : (grid === '4x4' ? '8.5px' : '9.5px') }" style="text-align: center;">
            {{ isHero(idx - 1) ? 'DESTAQUE // HERO' : `SLOT ${idx}` }}
          </span>
        </template>
      </div>
    </div>
  </div>
</template>
