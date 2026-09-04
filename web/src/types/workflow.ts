export type ChannelType = 'telegram' | 'webhook' | 'email' | 'mqtt' | 'websocket'

export interface TelegramConfig {
  bot_token: string
  chat_id: string
  include_snapshot: boolean
}

export interface NotificationChannel {
  id: string
  name: string
  channel_type: ChannelType
  is_active: boolean
  config_json: TelegramConfig | Record<string, unknown>
  created_at: string
}

export interface WorkflowFilterConditions {
  camera_ids?: string[]
  event_types?: string[]
  min_confidence?: number
  severities?: ('info' | 'warning' | 'critical')[]
}

export interface WorkflowAction {
  channel_id: string
  channel_type: ChannelType
  channel_name: string
}

export interface NotificationWorkflow {
  id: string
  name: string
  description?: string
  is_enabled: boolean
  trigger_type: string
  filter_conditions: WorkflowFilterConditions
  cooldown_seconds: number
  actions_pipeline: WorkflowAction[]
  last_triggered_at?: string
  created_at: string
}
