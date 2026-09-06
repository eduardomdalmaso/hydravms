<script setup lang="ts">
import type { PaletteBlockDefinition } from '../../../types/workflowTree'
import { paletteBlocks } from '../../../data/mockWorkflowFolders'
import FlowNodeIcon from './FlowNodeIcon.vue'

const props = defineProps<{ isLocked?: boolean }>()
const filtered = (type: 'TRIGGER' | 'CONDITION' | 'ACTION') => paletteBlocks.filter(b => b.type === type)

const handleDragStart = (e: DragEvent, b: PaletteBlockDefinition) => {
  if (props.isLocked || !e.dataTransfer) return
  e.dataTransfer.setData('application/json', JSON.stringify(b))
  e.dataTransfer.setData('text/plain', b.kind)
  e.dataTransfer.effectAllowed = 'copy'
}
</script>

<template>
  <div class="vms-flow-assembly-pane">
    <div class="vms-assembly-header">
      <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">OPCOES DE MONTAGEM (FLOW)</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">ARRASTE OS BLOCOS PARA O CANVAS</span>
    </div>

    <div class="vms-assembly-body">
      <div class="vms-palette-section">
        <span class="vms-palette-section-title">[1. GATILHOS // PLAY]</span>
        <div class="vms-palette-cards">
          <div
            v-for="b in filtered('TRIGGER')" :key="b.kind" class="vms-palette-card"
            :class="{ disabled: isLocked }" :draggable="!isLocked"
            title="Arraste o bloco para o canvas"
            @dragstart="handleDragStart($event, b)"
          >
            <div class="vms-palette-icon-box"><FlowNodeIcon :kind="b.kind" :type="b.type" /></div>
            <div class="vms-palette-info">
              <span class="vms-palette-title">{{ b.title }}</span>
              <span class="vms-palette-summary">{{ b.summary }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="vms-palette-section">
        <span class="vms-palette-section-title">[2. CONDICOES // FILTROS]</span>
        <div class="vms-palette-cards">
          <div
            v-for="b in filtered('CONDITION')" :key="b.kind" class="vms-palette-card"
            :class="{ disabled: isLocked }" :draggable="!isLocked"
            title="Arraste o bloco para o canvas"
            @dragstart="handleDragStart($event, b)"
          >
            <div class="vms-palette-icon-box"><FlowNodeIcon :kind="b.kind" :type="b.type" /></div>
            <div class="vms-palette-info">
              <span class="vms-palette-title">{{ b.title }}</span>
              <span class="vms-palette-summary">{{ b.summary }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="vms-palette-section">
        <span class="vms-palette-section-title">[3. ACOES & TERMINAL]</span>
        <div class="vms-palette-cards">
          <div
            v-for="b in filtered('ACTION')" :key="b.kind" class="vms-palette-card"
            :class="{ disabled: isLocked }" :draggable="!isLocked"
            title="Arraste o bloco para o canvas"
            @dragstart="handleDragStart($event, b)"
          >
            <div class="vms-palette-icon-box"><FlowNodeIcon :kind="b.kind" :type="b.type" /></div>
            <div class="vms-palette-info">
              <span class="vms-palette-title">{{ b.title }}</span>
              <span class="vms-palette-summary">{{ b.summary }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-flow-assembly-pane { width: 270px; background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg); display: flex; flex-direction: column; overflow: hidden; }
.vms-assembly-header { padding: 10px 12px; background: #07080c; border-bottom: 1px solid var(--vms-border); display: flex; flex-direction: column; gap: 2px; }
.vms-assembly-body { padding: 8px 10px; overflow-y: auto; display: flex; flex-direction: column; gap: 12px; flex: 1; }
.vms-palette-section { display: flex; flex-direction: column; gap: 6px; }
.vms-palette-section-title { font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-orange); font-weight: 700; }
.vms-palette-cards { display: flex; flex-direction: column; gap: 5px; }
.vms-palette-card { background: #14171d; border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 8px; cursor: grab; display: flex; align-items: center; gap: 8px; transition: all 0.15s ease; user-select: none; -webkit-user-drag: element; }
.vms-palette-card * { pointer-events: none; }
.vms-palette-card:hover { border-color: var(--vms-neu-accent-orange); background: #1c212a; transform: translateX(2px); }
.vms-palette-card:active { cursor: grabbing; }
.vms-palette-card.disabled { opacity: 0.5; cursor: not-allowed; }
.vms-palette-icon-box { width: 34px; height: 34px; border-radius: 8px; background: #07080c; border: 1px solid var(--vms-border); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.vms-palette-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.vms-palette-title { color: #ffffff; font-size: 11px; font-weight: 700; }
.vms-palette-summary { color: var(--vms-text-dim); font-size: 8.5px; font-family: var(--vms-font-mono); }
</style>
