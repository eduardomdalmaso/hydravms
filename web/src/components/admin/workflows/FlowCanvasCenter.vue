<script setup lang="ts">
import { ref } from 'vue'
import { VueFlow, useVueFlow, type Node, type Edge, type Connection } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import type { PaletteBlockDefinition } from '../../../types/workflowTree'
import FlowCustomNode from './FlowCustomNode.vue'
import FlowNodeContextMenu from './FlowNodeContextMenu.vue'
import '../../../assets/css/vue-flow-dark.css'

const props = defineProps<{ nodes: Node[]; edges: Edge[]; selectedNodeId: string | null; isLocked?: boolean }>()
const emit = defineEmits<{
  (e: 'selectNode', id: string): void; (e: 'removeNode', id: string): void; (e: 'removeEdge', id: string): void; (e: 'testFlow'): void
  (e: 'connect', conn: Connection): void; (e: 'nodeDragStop', evt: any): void
  (e: 'dropBlock', payload: { block: PaletteBlockDefinition; pos: { x: number; y: number } }): void
}>()

const { project } = useVueFlow()
const stageRef = ref<HTMLElement | null>(null)
const contextMenu = ref<{ x: number; y: number; nodeId: string; title: string } | null>(null)

const onNodeContextMenu = (e: any) => {
  e.event?.preventDefault?.()
  const x = e.event?.clientX ?? (e.event?.touches?.[0]?.clientX || 100)
  const y = e.event?.clientY ?? (e.event?.touches?.[0]?.clientY || 100)
  contextMenu.value = { x, y, nodeId: e.node.id, title: e.node.data?.title || 'BLOCO' }
}

const handleDeleteFromContext = () => {
  if (contextMenu.value) { emit('removeNode', contextMenu.value.nodeId); contextMenu.value = null }
}

const onDragOver = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer) event.dataTransfer.dropEffect = 'copy'
}

const onDrop = (event: DragEvent) => {
  event.preventDefault()
  event.stopPropagation()
  if (props.isLocked) return
  const data = event.dataTransfer?.getData('application/json')
  if (!data) return
  try {
    const block = JSON.parse(data) as PaletteBlockDefinition
    const rect = stageRef.value?.getBoundingClientRect() || (event.currentTarget as HTMLElement)?.getBoundingClientRect()
    const pos = project ? project({ x: event.clientX - (rect?.left || 0), y: event.clientY - (rect?.top || 0) }) : { x: 200, y: 150 }
    emit('dropBlock', { block, pos })
  } catch (err) { console.error(err) }
}
</script>

<template>
  <div class="vms-flow-canvas-wrapper">
    <div class="vms-canvas-top-toolbar">
      <div class="vms-flex-row" style="gap: 8px; align-items: center;">
        <span class="vms-font-bold vms-text-xs" style="color: #ffffff;">CANVAS INTERATIVO (VUE FLOW)</span>
        <span class="vms-badge" style="font-size: 8px; background: #07080c; border: 1px solid var(--vms-border); color: #ff5e3a;">
          {{ nodes.length }} BLOCOS // {{ edges.length }} CONEXOES
        </span>
      </div>
      <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 3px 8px; font-weight: bold;" @click="emit('testFlow')">
        TESTAR DISPARO
      </button>
    </div>

    <div ref="stageRef" class="vms-canvas-stage" @dragover="onDragOver" @dragenter.prevent @drop="onDrop">
      <VueFlow
        :nodes="nodes" :edges="edges" :default-viewport="{ zoom: 1 }" :min-zoom="0.2" :max-zoom="3"
        :nodes-draggable="!isLocked" :nodes-connectable="!isLocked"
        @connect="emit('connect', $event)" @node-drag-stop="emit('nodeDragStop', $event)"
        @node-click="emit('selectNode', $event.node.id)" @pane-click="emit('selectNode', ''); contextMenu = null"
        @edge-click="emit('removeEdge', $event.edge.id)"
        @node-context-menu="onNodeContextMenu"
      >
        <Background pattern-color="#ff5e3a" :gap="20" :size="1" />
        <Controls position="bottom-right" />
        <template #node-custom="nodeProps">
          <FlowCustomNode :id="nodeProps.id" :data="nodeProps.data" :selected="selectedNodeId === nodeProps.id" />
        </template>
      </VueFlow>
    </div>

    <FlowNodeContextMenu
      v-if="contextMenu" :x="contextMenu.x" :y="contextMenu.y" :node-title="contextMenu.title"
      :is-locked="isLocked" @delete="handleDeleteFromContext" @close="contextMenu = null"
    />
  </div>
</template>

<style scoped>
.vms-flow-canvas-wrapper { flex: 1; min-width: 0; background: #07080c; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg); display: flex; flex-direction: column; overflow: hidden; height: 100%; }
.vms-canvas-top-toolbar { padding: 8px 12px; background: #0e1117; border-bottom: 1px solid var(--vms-border); display: flex; justify-content: space-between; align-items: center; }
.vms-canvas-stage { flex: 1; position: relative; width: 100%; height: 100%; overflow: hidden; }
</style>
