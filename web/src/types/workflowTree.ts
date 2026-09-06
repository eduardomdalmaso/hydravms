export type FlowNodeType = 'TRIGGER' | 'CONDITION' | 'ACTION'

export interface FlowEdgeItem {
  id: string
  source: string
  target: string
  label?: string
  animated?: boolean
}

export interface FlowNodeItem {
  id: string
  type: FlowNodeType
  kind: string
  title: string
  summary: string
  position?: { x: number; y: number }
  config: Record<string, any>
}

export interface EnterpriseWorkflowItem {
  id: string
  name: string
  companyScope: string
  folderId?: string
  is_locked: boolean
  status: 'ATIVO' | 'PAUSADO'
  cooldownSeconds: number
  nodes: FlowNodeItem[]
  edges?: FlowEdgeItem[]
  createdAt: string
  description?: string
}

export interface WorkflowFolderNode {
  id: string
  name: string
  clientType: 'company' | 'final_client'
  isExpanded: boolean
  workflows: EnterpriseWorkflowItem[]
}

export interface PaletteBlockDefinition {
  kind: string
  type: FlowNodeType
  title: string
  summary: string
  defaultConfig: Record<string, any>
}
