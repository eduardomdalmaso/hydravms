<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWorkflows } from '../../composables/useWorkflows'
import type { NotificationWorkflow } from '../../types/workflow'
import WorkflowCard from '../../components/workflows/WorkflowCard.vue'
import WorkflowBuilderModal from '../../components/workflows/WorkflowBuilderModal.vue'
import TelegramChannelModal from '../../components/workflows/TelegramChannelModal.vue'

const { workflows, channels, fetchWorkflows, fetchChannels, toggleWorkflow } = useWorkflows()
const isBuilderOpen = ref(false)
const isTelegramOpen = ref(false)

onMounted(async () => {
  await Promise.all([fetchWorkflows(), fetchChannels()])
})

const handleCreateWorkflow = (wf: Partial<NotificationWorkflow>) => {
  workflows.value.unshift({
    id: `wf-${Date.now()}`,
    name: wf.name || 'Novo Workflow',
    description: wf.description,
    is_enabled: true,
    trigger_type: wf.trigger_type || 'ai_event',
    filter_conditions: wf.filter_conditions || {},
    cooldown_seconds: wf.cooldown_seconds || 30,
    actions_pipeline: wf.actions_pipeline || [],
    created_at: new Date().toISOString()
  })
}

const handleCreateTelegram = (data: { name: string; token: string; chatId: string; snapshot: boolean }) => {
  channels.value.push({
    id: `ch-tg-${Date.now()}`,
    name: data.name,
    channel_type: 'telegram',
    is_active: true,
    config_json: { bot_token: data.token, chat_id: data.chatId, include_snapshot: data.snapshot },
    created_at: new Date().toISOString()
  })
}

const handleTestWorkflow = (_wf: NotificationWorkflow) => {
  window.alert('Disparo de teste enviado aos canais configurados!')
}
</script>

<template>
  <div class="vms-content-area">
    <div class="vms-page-header">
      <div class="vms-page-title-group">
        <h1 class="vms-h1">FLUXOS DE NOTIFICACAO & AUTOMACAO</h1>
        <span class="vms-text-sm vms-text-muted">Configure regras de disparo de alertas analiticos para Telegram, WebSockets e Webhooks.</span>
      </div>
      <div class="vms-page-actions">
        <button class="vms-btn vms-btn-secondary" @click="isTelegramOpen = true">[+ CONFIGURAR BOT TELEGRAM]</button>
        <button class="vms-btn vms-btn-primary" @click="isBuilderOpen = true">[+ CRIAR NOVO WORKFLOW]</button>
      </div>
    </div>
    <div class="vms-workflow-grid">
      <WorkflowCard
        v-for="wf in workflows"
        :key="wf.id"
        :workflow="wf"
        @toggle="toggleWorkflow"
        @test="handleTestWorkflow"
      />
    </div>
    <WorkflowBuilderModal
      v-if="isBuilderOpen"
      :channels="channels"
      @close="isBuilderOpen = false"
      @created="handleCreateWorkflow"
      @openTelegram="() => { isBuilderOpen = false; isTelegramOpen = true }"
    />
    <TelegramChannelModal
      v-if="isTelegramOpen"
      @close="isTelegramOpen = false"
      @created="handleCreateTelegram"
    />
  </div>
</template>
