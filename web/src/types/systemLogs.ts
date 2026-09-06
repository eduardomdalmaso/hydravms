export type LogCategory = 'SYSTEM' | 'AUDIT'
export type LogLevel = 'INFO' | 'WARNING' | 'CRITICAL'
export type ViewerRole = 'SUPERADMIN' | 'ADMIN'

export interface LogEntry {
  id: string
  timestamp: string
  category: LogCategory
  level: LogLevel
  actor: string
  tenantId: string
  tenantName: string
  action: string
  target: string
  details: string
  ipAddress: string
}

export interface LogFilterState {
  category: LogCategory | 'ALL'
  level: LogLevel | 'ALL'
  searchQuery: string
  viewerRole: ViewerRole
  activeTenantId: string
}
