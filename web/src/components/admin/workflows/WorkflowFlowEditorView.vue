<script setup lang="ts">
import type { EnterpriseWorkflowItem } from "../../../types/workflowTree"
import FlowPaletteLeftPane from "./FlowPaletteLeftPane.vue"
import FlowCanvasCenter from "./FlowCanvasCenter.vue"
import FlowOptionsRightPane from "./FlowOptionsRightPane.vue"
import { useFlowEditorState } from "../../../composables/useFlowEditorState"

const props = defineProps<{ workflow: EnterpriseWorkflowItem }>()
const emit = defineEmits<{ (e: "notify", msg: string): void }>()

const {
  nodes,
  edges,
  selectedNodeId,
  selectedNode,
  selectNode,
  addNodeAtPosition,
  removeNode,
  removeEdge,
  onConnect,
  onNodeDragStop,
  testFlowExecution,
} = useFlowEditorState(props.workflow, (msg: string) => emit("notify", msg))
</script>

<template>
  <div class="flow-editor-layout">
    <FlowPaletteLeftPane :is-locked="workflow.is_locked" />
    <FlowCanvasCenter
      :nodes="nodes"
      :edges="edges"
      :selected-node-id="selectedNodeId"
      :is-locked="workflow.is_locked"
      @select-node="selectNode"
      @remove-node="removeNode"
      @remove-edge="removeEdge"
      @test-flow="testFlowExecution"
      @connect="onConnect"
      @node-drag-stop="onNodeDragStop"
      @drop-block="({ block, pos }) => addNodeAtPosition(block, pos)"
    />
    <FlowOptionsRightPane
      :workflow="workflow"
      :selected-node="selectedNode"
      :is-locked="workflow.is_locked"
    />
  </div>
</template>

<style scoped>
.flow-editor-layout {
  display: flex;
  gap: 1rem;
  height: calc(100vh - 165px);
  width: 100%;
  overflow: hidden;
  box-sizing: border-box;
}
</style>
