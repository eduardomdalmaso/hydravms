<script setup lang="ts">
import { computed } from 'vue'
import type { CustomLayout, GridLayout } from '../../types/mosaic'

const props = defineProps<{
  x: number
  y: number
  layout?: CustomLayout | null
  isHeaderMenu?: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'rename', layout: CustomLayout): void
  (e: 'duplicate', layout: CustomLayout): void
  (e: 'delete', layout: CustomLayout): void
  (e: 'selectGrid', grid: GridLayout): void
}>()

const availableGrids: GridLayout[] = [
  '1x1', '2x2', '3x3', '4x4', '8x8', '10x10', '1+5', '1+7', '1+12'
]

const menuPos = computed(() => {
  const maxH = typeof window !== 'undefined' ? window.innerHeight : 800
  const maxW = typeof window !== 'undefined' ? window.innerWidth : 1200
  return {
    top: `${Math.min(props.y, maxH - 320)}px`,
    left: `${Math.min(props.x, maxW - 220)}px`
  }
})
</script>

<template>
  <div class="vms-context-menu-backdrop" @click.self="emit('close')" @contextmenu.prevent="emit('close')">
    <div class="vms-context-menu" :style="menuPos">
      <!-- ITEM CONTEXT MENU -->
      <template v-if="layout && !isHeaderMenu">
        <div class="vms-context-menu-header">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">{{ layout.name }}</span>
          <span v-if="layout.is_system" class="vms-badge vms-badge-recording" style="font-size: 8px;">[ADMIN // READ-ONLY]</span>
          <span v-else class="vms-badge vms-badge-online" style="font-size: 8px;">[OPERADOR]</span>
        </div>

        <button v-if="!layout.is_system" class="vms-context-item" @click="emit('rename', layout)">
          <span>[RENOMEAR LAYOUT]</span>
        </button>
        <button class="vms-context-item" @click="emit('duplicate', layout)">
          <span>[DUPLICAR LAYOUT]</span>
        </button>
        <button v-if="!layout.is_system" class="vms-context-item danger" @click="emit('delete', layout)">
          <span>[EXCLUIR LAYOUT]</span>
        </button>
        <div class="vms-context-divider"></div>
      </template>

      <!-- GRID SELECTION QUICK PICKER -->
      <div class="vms-context-menu-section-title">
        <span>SELECIONAR FORMATO DE GRADE:</span>
      </div>
      <div class="vms-context-grid-picker">
        <button
          v-for="g in availableGrids"
          :key="g"
          class="vms-btn vms-btn-ghost vms-btn-sm"
          style="font-size: 10px; padding: 3px 6px;"
          @click="emit('selectGrid', g)"
        >
          [{{ g.toUpperCase() }}]
        </button>
      </div>
    </div>
  </div>
</template>
