<script setup lang="ts">
import type { NotificationChannel, WorkflowAction } from '../../types/workflow'

const props = defineProps<{
  channels: NotificationChannel[]
  selectedActions: WorkflowAction[]
}>()

const emit = defineEmits<{
  (e: 'toggleAction', channel: NotificationChannel): void
  (e: 'openTelegramModal'): void
}>()

const isChannelSelected = (channelId: string) => {
  return props.selectedActions.some(a => a.channel_id === channelId)
}
</script>

<template>
  <div class="vms-flex-col">
    <div class="vms-flex-between">
      <label class="vms-label">Destinos de Notificacao (Pipeline de Acoes)</label>
      <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('openTelegramModal')">
        [+ ADICIONAR BOT TELEGRAM]
      </button>
    </div>
    <div class="vms-grid-container" style="grid-template-columns: 1fr; gap: 0.5rem;">
      <div
        v-for="ch in channels"
        :key="ch.id"
        class="vms-card vms-card-elevated"
        :style="{
          padding: '0.75rem 1rem',
          cursor: 'pointer',
          borderColor: isChannelSelected(ch.id) ? 'var(--vms-primary)' : 'var(--vms-border-light)'
        }"
        @click="emit('toggleAction', ch)"
      >
        <div class="vms-flex-between">
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <span
              class="vms-badge"
              :class="ch.channel_type === 'telegram' ? 'vms-channel-pill-telegram' : 'vms-badge-info'"
            >
              [{{ ch.channel_type.toUpperCase() }}]
            </span>
            <span class="vms-text-sm vms-font-semibold">{{ ch.name }}</span>
          </div>
          <input type="checkbox" :checked="isChannelSelected(ch.id)" class="vms-checkbox" />
        </div>
      </div>
    </div>
  </div>
</template>
