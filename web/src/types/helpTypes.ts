export interface HelpSection {
  id: string
  title: string
  tag: string
  description: string
  steps: string[]
  animationType: 'streams_scan' | 'grid_slots' | 'rbac_tree' | 'brand_palette' | 'flow_pipes' | 'cluster_nodes' | 'storage_pool'
  example: { label: string; scenario: string; result: string }
}

export interface PageHelpGuide {
  pageId: string
  title: string
  subtitle: string
  badge: string
  sections: HelpSection[]
}
