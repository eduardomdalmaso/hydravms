<script setup lang="ts">
import type { AnalyticMode, Point2D, Line2D } from '../../../types/marketplace'
import AnalyticRoiCanvas from './AnalyticRoiCanvas.vue'

defineProps<{
  name: string
  camera: string
  cameras: { id: string; name: string }[]
  mode: AnalyticMode
  selectedClasses: string[]
  points: Point2D[]
  line: Line2D
  snapshotUrl: string
}>()

const emit = defineEmits<{
  (e: 'update:name', val: string): void
  (e: 'update:camera', val: string): void
  (e: 'update:mode', val: AnalyticMode): void
  (e: 'toggle:class', id: string): void
  (e: 'update:points', val: Point2D[]): void
  (e: 'update:line', val: Line2D): void
}>()

const availableClasses = [
  { id: 'person', label: 'Pessoa' }, { id: 'cell_phone', label: 'Celular' },
  { id: 'car', label: 'Carro' }, { id: 'motorcycle', label: 'Moto' },
  { id: 'bicycle', label: 'Bicicleta' }, { id: 'truck', label: 'Caminhão' },
  { id: 'bus', label: 'Ônibus' }, { id: 'bird', label: 'Pássaros' },
  { id: 'cat', label: 'Gato' }, { id: 'dog', label: 'Cachorro' }
]

const modes = [
  { id: 'intrusion', label: 'Intrusão' }, { id: 'crowd', label: 'Multidão' },
  { id: 'counting', label: 'Contagem' }, { id: 'dwell_time', label: 'Tempo Excedido' }
] as const
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.85rem;">
    <div class="vms-flex-row" style="gap: 0.75rem;">
      <div class="vms-flex-col" style="gap: 0.2rem; flex: 1.2;">
        <label class="vms-text-xs vms-font-semibold">NOME DO ANALÍTICO:</label>
        <input :value="name" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;" placeholder="Ex: Portaria Principal // Intrusão..." @input="emit('update:name', ($event.target as HTMLInputElement).value)" />
      </div>
      <div class="vms-flex-col" style="gap: 0.2rem; flex: 1;">
        <label class="vms-text-xs vms-font-semibold">CÂMERA DE VÍDEO:</label>
        <select :value="camera" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;" @change="emit('update:camera', ($event.target as HTMLSelectElement).value)">
          <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
      </div>
    </div>

    <!-- Modos de Análise -->
    <div class="vms-flex-col" style="gap: 0.35rem;">
      <label class="vms-text-xs vms-font-semibold">MODO DE ANÁLISE:</label>
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
        <button v-for="m in modes" :key="m.id" class="vms-btn vms-btn-sm" :class="mode === m.id ? 'vms-btn-primary' : 'vms-btn-secondary'" style="padding: 4px 10px; font-size: 11px;" @click="emit('update:mode', m.id)">{{ m.label }}</button>
      </div>
    </div>

    <!-- Canvas com Polígono / Linha -->
    <AnalyticRoiCanvas :snapshot-url="snapshotUrl" :mode="mode" :points="points" :line="line" @update:points="emit('update:points', $event)" @update:line="emit('update:line', $event)" />

    <!-- Objetos -->
    <div class="vms-flex-col" style="gap: 0.3rem;">
      <label class="vms-text-xs vms-font-semibold">OBJETOS PARA ANALISAR:</label>
      <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
        <span v-for="cls in availableClasses" :key="cls.id" class="vms-chip" :class="{ active: selectedClasses.includes(cls.id) }" @click="emit('toggle:class', cls.id)">{{ cls.label }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-chip { font-size: 10px; font-family: var(--vms-font-mono); padding: 3px 8px; border-radius: 4px; background: rgba(255,255,255,0.06); color: #8b94a0; cursor: pointer; border: 1px solid transparent; }
.vms-chip.active { background: rgba(255,94,58,0.2); color: #ff5e3a; border-color: rgba(255,94,58,0.4); font-weight: bold; }
</style>
