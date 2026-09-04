import { ref } from 'vue'
import type { NotificationWorkflow, NotificationChannel } from '../types/workflow'

export function useWorkflows() {
  const workflows = ref<NotificationWorkflow[]>([])
  const channels = ref<NotificationChannel[]>([])
  const isLoading = ref(false)
  const isSaving = ref(false)
  const testStatus = ref<string | null>(null)

  const fetchWorkflows = async () => {
    isLoading.value = true
    try {
      // Mock / API call
      workflows.value = [
        {
          id: 'wf-1',
          name: 'Notificação de Invasão de Perímetro (Noite)',
          description: 'Dispara Telegram e WebSocket quando pessoa detectada',
          is_enabled: true,
          trigger_type: 'ai_event',
          filter_conditions: {
            event_types: ['person_intrusion'],
            min_confidence: 0.85,
            severities: ['critical']
          },
          cooldown_seconds: 30,
          actions_pipeline: [
            { channel_id: 'ch-tg-1', channel_type: 'telegram', channel_name: 'Telegram - Central 24h' },
            { channel_id: 'ch-ws-1', channel_type: 'websocket', channel_name: 'WebSocket Live Feed' }
          ],
          last_triggered_at: '2026-09-03T19:30:00Z',
          created_at: '2026-09-01T10:00:00Z'
        }
      ]
    } finally {
      isLoading.value = false
    }
  }

  const fetchChannels = async () => {
    channels.value = [
      {
        id: 'ch-tg-1',
        name: 'Telegram - Central 24h',
        channel_type: 'telegram',
        is_active: true,
        config_json: { bot_token: '••••••••', chat_id: '-1001928374', include_snapshot: true },
        created_at: '2026-09-01T08:00:00Z'
      },
      {
        id: 'ch-ws-1',
        name: 'WebSocket Live Feed',
        channel_type: 'websocket',
        is_active: true,
        config_json: {},
        created_at: '2026-09-01T08:00:00Z'
      }
    ]
  }

  const toggleWorkflow = async (workflow: NotificationWorkflow) => {
    workflow.is_enabled = !workflow.is_enabled
  }

  const testTelegramChannel = async (_token: string, _chatId: string) => {
    testStatus.value = 'sending'
    setTimeout(() => { testStatus.value = 'success' }, 1000)
  }

  return {
    workflows,
    channels,
    isLoading,
    isSaving,
    testStatus,
    fetchWorkflows,
    fetchChannels,
    toggleWorkflow,
    testTelegramChannel
  }
}
