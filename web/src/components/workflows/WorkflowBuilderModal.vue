<script setup lang="ts">
import { ref } from 'vue'
import type { NotificationChannel, WorkflowAction, NotificationWorkflow } from '../../types/workflow'
import WorkflowTriggerStep from './WorkflowTriggerStep.vue'
import WorkflowConditionStep from './WorkflowConditionStep.vue'
import WorkflowActionStep from './WorkflowActionStep.vue'

const props = defineProps<{ channels: NotificationChannel[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'created', wf: Partial<NotificationWorkflow>): void; (e: 'openTelegram'): void }>()

const step = ref(1)
const workflowName = ref('')
const triggerType = ref('ai_event')
const selectedEventType = ref('person_intrusion')
const selectedSeverity = ref('critical')
const minConfidence = ref(0.85)
const cooldownSeconds = ref(30)
const selectedActions = ref<WorkflowAction[]>([])

const handleToggleAction = (ch: NotificationChannel) => {
  const idx = selectedActions.value.findIndex(a => a.channel_id === ch.id)
  if (idx >= 0) selectedActions.value.splice(idx, 1)
  else selectedActions.value.push({ channel_id: ch.id, channel_type: ch.channel_type, channel_name: ch.name })
}

const handleSave = () => {
  if (!workflowName.value) return
  emit('created', {
    name: workflowName.value,
    trigger_type: triggerType.value,
    filter_conditions: { event_types: [selectedEventType.value], min_confidence: minConfidence.value, severities: [selectedSeverity.value as any] },
    cooldown_seconds: cooldownSeconds.value,
    actions_pipeline: selectedActions.value,
    is_enabled: true
  })
  emit('close')
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog vms-modal-lg">
      <div class="vms-modal-header">
        <h3 class="vms-h3">NOVO FLUXO DE NOTIFICACAO DE EVENTOS</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="emit('close')">[FECHAR]</button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-stepper">
          <div class="vms-step-item" :class="{ active: step === 1 }"><span class="vms-step-number">1</span> GATILHO (TRIGGER)</div>
          <div class="vms-step-item" :class="{ active: step === 2 }"><span class="vms-step-number">2</span> FILTROS & CONDICOES</div>
          <div class="vms-step-item" :class="{ active: step === 3 }"><span class="vms-step-number">3</span> ACOES & DESTINOS</div>
        </div>
        <WorkflowTriggerStep v-if="step === 1" v-model:triggerType="triggerType" v-model:workflowName="workflowName" />
        <WorkflowConditionStep v-else-if="step === 2" v-model:minConfidence="minConfidence" v-model:cooldownSeconds="cooldownSeconds" v-model:selectedEventType="selectedEventType" v-model:selectedSeverity="selectedSeverity" />
        <WorkflowActionStep v-else-if="step === 3" :channels="channels" :selectedActions="selectedActions" @toggleAction="handleToggleAction" @openTelegramModal="emit('openTelegram')" />
      </div>
      <div class="vms-modal-footer">
        <button v-if="step > 1" class="vms-btn vms-btn-secondary" @click="step--">VOLTAR</button>
        <button v-if="step < 3" class="vms-btn vms-btn-primary" @click="step++">AVANCAR</button>
        <button v-else class="vms-btn vms-btn-primary" @click="handleSave">SALVAR E ATIVAR WORKFLOW</button>
      </div>
    </div>
  </div>
</template>
