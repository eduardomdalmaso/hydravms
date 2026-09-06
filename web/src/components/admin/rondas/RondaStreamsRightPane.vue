<script setup lang="ts">
import { ref, computed } from 'vue'
import type { EnterpriseRondaItem } from '../../../types/rondaTree'
import RondaStreamItemCard from './RondaStreamItemCard.vue'
import RondaAvailablePalette from './RondaAvailablePalette.vue'

const props = defineProps<{ ronda: EnterpriseRondaItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
const isDragOver = ref(false)

const totalDuration = computed(() => props.ronda.streams.reduce((acc, s) => acc + (s.intervalSeconds || 0), 0))
const usedCameraIds = computed(() => props.ronda.streams.map(s => s.cameraId))

const handleUpdateInterval = (idx: number, val: number) => { props.ronda.streams[idx].intervalSeconds = Math.max(1, Math.min(300, val)) }
const handleMoveUp = (i: number) => { if (i > 0) props.ronda.streams.splice(i - 1, 0, props.ronda.streams.splice(i, 1)[0]) }
const handleMoveDown = (i: number) => { if (i < props.ronda.streams.length - 1) props.ronda.streams.splice(i + 1, 0, props.ronda.streams.splice(i, 1)[0]) }
const handleRemove = (idx: number) => { props.ronda.streams.splice(idx, 1) }

const handleAddCamera = (cam: { id: string; name: string; resolution: string; fps: number }) => {
  props.ronda.streams.push({
    id: `rs_${Date.now()}_${props.ronda.streams.length}`,
    cameraId: cam.id,
    cameraName: cam.name,
    resolution: cam.resolution,
    fps: cam.fps,
    intervalSeconds: 10,
    orderIndex: props.ronda.streams.length
  })
}

const onDrop = (ev: DragEvent) => {
  ev.preventDefault(); isDragOver.value = false
  const raw = ev.dataTransfer?.getData('application/json')
  if (raw) { try { handleAddCamera(JSON.parse(raw)) } catch {} }
}
</script>

<template>
  <div class="vms-split-pane">
    <div class="vms-split-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">
          SEQUENCIA DE FLUXOS & TEMPOS INDIVIDUAIS
        </span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          // DEFINA O TEMPO DE CADA CAMERA E ARRASTE NOVOS FLUXOS
        </span>
      </div>
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border);">
        CICLO: {{ totalDuration }}s // {{ ronda.streams.length }} CAMS
      </span>
    </div>

    <!-- Patrol Sequence Flow Path -->
    <div v-if="ronda.streams.length > 0" class="vms-flex-row" style="gap: 0.35rem; align-items: center; overflow-x: auto; padding: 0.45rem 0.65rem; background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; flex-shrink: 0;">
      <template v-for="(s, idx) in ronda.streams" :key="s.id">
        <span class="vms-text-mono vms-text-2xs" style="color: #fff; white-space: nowrap;">
          [{{ idx + 1 }}] {{ s.cameraName.split(' ')[0] }}
          <strong style="color: #ff5e3a;">({{ s.intervalSeconds }}s)</strong>
        </span>
        <span v-if="idx < ronda.streams.length - 1" style="color: var(--vms-neu-accent-orange); font-size: 10px;">➔</span>
      </template>
    </div>

    <!-- Vertical List / Drop Zone of Cameras in Ronda -->
    <div
      class="vms-flex-col"
      style="gap: 0.4rem; overflow-y: auto; max-height: 240px; flex: 1; padding: 4px; border-radius: 6px; transition: all 0.2s ease;"
      :style="{
        background: isDragOver ? 'rgba(255, 94, 58, 0.05)' : 'transparent',
        border: isDragOver ? '1px dashed #ff5e3a' : '1px dashed transparent'
      }"
      @dragover.prevent="isDragOver = true"
      @dragleave="isDragOver = false"
      @drop="onDrop"
    >
      <div v-if="ronda.streams.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="text-align: center; padding: 2rem; border: 1px dashed var(--vms-border); border-radius: 6px;">
        // ARRASTE AS CAMERAS AQUI OU CLIQUE EM [+] NA PALETA ABAIXO
      </div>
      <RondaStreamItemCard
        v-for="(s, idx) in ronda.streams"
        :key="s.id"
        :stream="s"
        :index="idx"
        :is-first="idx === 0"
        :is-last="idx === ronda.streams.length - 1"
        @update-interval="(val) => handleUpdateInterval(idx, val)"
        @move-up="handleMoveUp(idx)"
        @move-down="handleMoveDown(idx)"
        @remove="handleRemove(idx)"
      />
    </div>

    <!-- Available Cameras Palette with Drag & Click -->
    <RondaAvailablePalette :used-camera-ids="usedCameraIds" @add-camera="handleAddCamera" />
  </div>
</template>
