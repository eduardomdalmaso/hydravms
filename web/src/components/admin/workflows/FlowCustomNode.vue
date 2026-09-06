<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import type { FlowNodeItem } from '../../../types/workflowTree'
import FlowNodeIcon from './FlowNodeIcon.vue'

defineProps<{
  id: string
  data: FlowNodeItem
  selected?: boolean
}>()
</script>

<template>
  <div
    class="vms-flow-icon-node"
    :class="[data.type.toLowerCase(), { selected }]"
    :title="data.kind === 'TERMINAL_STOP' ? '[TERMINAL] PARE // Finalizador' : `[${data.type}] ${data.title}`"
  >
    <!-- Entradas na esquerda (se for Y OU ou AND, exibe 2 entradas) -->
    <template v-if="data.kind === 'BRANCH_OR' || data.kind === 'LOGIC_AND'">
      <Handle type="target" id="in_top" :position="Position.Left" style="top: 25%;" class="vms-handle left" />
      <Handle type="target" id="in_bot" :position="Position.Left" style="top: 75%;" class="vms-handle left" />
    </template>
    <Handle
      v-else-if="data.type !== 'TRIGGER'"
      type="target"
      :position="Position.Left"
      class="vms-handle left"
    />

    <FlowNodeIcon :kind="data.kind" :type="data.type" />

    <!-- Saidas na direita (se for Divisao em Y, exibe 2 saidas) -->
    <template v-if="data.kind === 'SPLIT_Y'">
      <Handle type="source" id="out_top" :position="Position.Right" style="top: 25%;" class="vms-handle right" />
      <Handle type="source" id="out_bot" :position="Position.Right" style="top: 75%;" class="vms-handle right" />
    </template>
    <Handle
      v-else-if="data.kind !== 'TERMINAL_STOP'"
      type="source"
      :position="Position.Right"
      class="vms-handle right"
    />
  </div>
</template>

<style scoped>
.vms-flow-icon-node { width: 62px; height: 62px; border-radius: 14px; background: #0f131a; border: 1.5px solid var(--vms-border); display: flex; align-items: center; justify-content: center; position: relative; cursor: grab; box-shadow: 0 6px 18px rgba(0, 0, 0, 0.6); transition: all 0.2s ease; }
.vms-flow-icon-node:hover { border-color: rgba(255, 94, 58, 0.6); transform: scale(1.05); }
.vms-flow-icon-node.selected { border-color: #ff5e3a; background: #1a202c; box-shadow: 0 0 20px rgba(255, 94, 58, 0.45); }
.vms-flow-icon-node.trigger { border-color: rgba(255, 94, 58, 0.45); }
.vms-handle { width: 9px; height: 9px; background: #ff5e3a; border: 2px solid #07080c; border-radius: 50%; }
</style>
