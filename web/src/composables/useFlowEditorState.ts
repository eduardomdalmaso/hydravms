import { ref, shallowRef, computed, watch } from 'vue'
import type { Node, Edge, Connection } from '@vue-flow/core'
import type { EnterpriseWorkflowItem, FlowNodeItem, PaletteBlockDefinition } from '../types/workflowTree'

export function useFlowEditorState(workflow: EnterpriseWorkflowItem, onNotify: (msg: string) => void) {
  const selectedNodeId = ref<string | null>(null)
  const nodes = shallowRef<Node[]>([])
  const edges = shallowRef<Edge[]>([])

  const syncFromWorkflow = () => {
    nodes.value = workflow.nodes.map((n, i) => ({
      id: n.id, type: 'custom', position: n.position || { x: 80 + i * 220, y: 120 }, data: n
    }))
    edges.value = (workflow.edges || []).map(e => ({
      id: e.id, source: e.source, target: e.target, sourceHandle: (e as any).sourceHandle, targetHandle: (e as any).targetHandle, animated: true, style: { stroke: '#ff5e3a', strokeWidth: 2 }
    }))
  }

  syncFromWorkflow()
  watch(() => workflow.id, syncFromWorkflow)

  const selectedNode = computed(() => {
    if (!selectedNodeId.value) return null
    return workflow.nodes.find(item => item.id === selectedNodeId.value) || null
  })

  const selectNode = (id: string) => {
    selectedNodeId.value = id || null
  }

  const addNodeAtPosition = (block: PaletteBlockDefinition, pos: { x: number; y: number }) => {
    if (workflow.is_locked) {
      onNotify('[TRAVADO] Desbloqueie o workflow para adicionar blocos')
      return
    }
    const newNode: FlowNodeItem = {
      id: `node_${Date.now()}`, type: block.type, kind: block.kind, title: block.title,
      summary: block.summary, position: pos, config: JSON.parse(JSON.stringify(block.defaultConfig))
    }
    workflow.nodes.push(newNode)
    const currentPositions = new Map(nodes.value.map(n => [n.id, n.position]))
    nodes.value = workflow.nodes.map(n => ({
      id: n.id, type: 'custom',
      position: currentPositions.get(n.id) || n.position || { x: 100, y: 100 },
      data: n
    }))
    if (newNode.kind !== 'TERMINAL_STOP') selectedNodeId.value = newNode.id
    onNotify(`[FLOW] Bloco "${block.title}" inserido no canvas`)
  }

  const removeNode = (nodeId: string) => {
    if (workflow.is_locked) {
      onNotify('[TRAVADO] Desbloqueie o workflow para remover blocos')
      return
    }
    const idx = workflow.nodes.findIndex(n => n.id === nodeId)
    if (idx >= 0) {
      const removed = workflow.nodes.splice(idx, 1)[0]
      nodes.value = nodes.value.filter(n => n.id !== nodeId)
      edges.value = edges.value.filter(e => e.source !== nodeId && e.target !== nodeId)
      workflow.edges = (workflow.edges || []).filter(e => e.source !== nodeId && e.target !== nodeId)
      if (selectedNodeId.value === nodeId) selectedNodeId.value = null
      onNotify(`[FLOW] Bloco "${removed.title}" removido`)
    }
  }

  const onConnect = (c: Connection) => {
    if (workflow.is_locked) return
    if (!workflow.edges) workflow.edges = []
    const edgeId = `e_${c.source}_${c.sourceHandle || ''}_${c.target}_${c.targetHandle || ''}`
    const edgeObj: any = { id: edgeId, source: c.source, target: c.target, sourceHandle: c.sourceHandle, targetHandle: c.targetHandle, animated: true }
    workflow.edges.push(edgeObj)
    edges.value = [...edges.value, { ...edgeObj, style: { stroke: '#ff5e3a', strokeWidth: 2 } }]
  }

  const onNodeDragStop = (e: any) => {
    const pos = { x: e.node.position.x, y: e.node.position.y }
    const n = workflow.nodes.find(item => item.id === e.node.id)
    if (n) n.position = pos
    const refNode = nodes.value.find(item => item.id === e.node.id)
    if (refNode) refNode.position = pos
  }

  const removeEdge = (edgeId: string) => {
    if (workflow.is_locked) { onNotify('[TRAVADO] Desbloqueie o workflow para alterar conexoes'); return }
    edges.value = edges.value.filter(e => e.id !== edgeId)
    workflow.edges = (workflow.edges || []).filter(e => e.id !== edgeId)
    onNotify('[FLOW] Interligacao removida')
  }

  const testFlowExecution = () => onNotify(`[SIMULACAO FLOW] Pipeline executado // ${nodes.value.length} nos, ${edges.value.length} conexoes OK`)

  return {
    nodes, edges, selectedNodeId, selectedNode, selectNode,
    addNodeAtPosition, removeNode, removeEdge, onConnect, onNodeDragStop, testFlowExecution
  }
}
