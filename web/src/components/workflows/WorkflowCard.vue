<script setup lang="ts">
import type { NotificationWorkflow } from '../../types/workflow'

const props = defineProps<{ workflow: NotificationWorkflow }>()
const emit = defineEmits<{ (e: 'toggle', wf: NotificationWorkflow): void; (e: 'test', wf: NotificationWorkflow): void }>()
</script>

<template>
  <div class="vms-card" :style="{ opacity: workflow.is_enabled ? 1 : 0.6 }">
    <div class="vms-card-header">
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <span class="vms-status-led" :class="workflow.is_enabled ? 'online' : 'offline'" :title="workflow.is_enabled ? 'Ativo' : 'Pausado'"></span>
        <h3 class="vms-h3">{{ workflow.name }}</h3>
      </div>
      <label class="vms-checkbox-label" title="Ativar/Desativar">
        <input type="checkbox" :checked="workflow.is_enabled" class="vms-checkbox" @change="emit('toggle', workflow)" />
      </label>
    </div>
    <div class="vms-card-body vms-flex-col">
      <p class="vms-text-xs vms-text-muted">{{ workflow.description || 'Disparo automatizado de analiticos' }}</p>
      <div class="vms-flow-container">
        <div class="vms-flow-node">
          <span class="vms-flow-node-title">GATILHO</span>
          <span class="vms-flow-node-value">{{ workflow.trigger_type }}</span>
        </div>
        <span class="vms-flow-arrow">//</span>
        <div class="vms-flow-node">
          <span class="vms-flow-node-title">FILTRO IA</span>
          <span class="vms-flow-node-value">CONF &gt; {{ Math.round((workflow.filter_conditions.min_confidence || 0.8) * 100) }}%</span>
        </div>
        <span class="vms-flow-arrow">//</span>
        <div class="vms-flow-node">
          <span class="vms-flow-node-title">CANAIS ({{ workflow.actions_pipeline.length }})</span>
          <div class="vms-flex-row" style="gap: 0.25rem; flex-wrap: wrap;">
            <span
              v-for="act in workflow.actions_pipeline"
              :key="act.channel_id"
              class="vms-channel-pill"
              :class="{ 'vms-channel-pill-telegram': act.channel_type === 'telegram' }"
            >
              [{{ act.channel_type.toUpperCase() }}]
            </span>
          </div>
        </div>
      </div>
    </div>
    <div class="vms-card-footer vms-flex-between">
      <span class="vms-text-mono vms-text-xs vms-text-dim">COOLDOWN: {{ workflow.cooldown_seconds }}S</span>
      <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('test', workflow)">[DISPARAR TESTE]</button>
    </div>
  </div>
</template>
