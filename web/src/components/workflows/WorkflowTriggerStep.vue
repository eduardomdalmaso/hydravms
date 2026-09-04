<script setup lang="ts">
const props = defineProps<{
  triggerType: string
  workflowName: string
}>()

const emit = defineEmits<{
  (e: 'update:triggerType', val: string): void
  (e: 'update:workflowName', val: string): void
}>()

const triggerOptions = [
  { id: 'ai_event', title: 'Detecção de IA em Tempo Real', desc: 'Invasão, cruzamento de linha, LPR, EPI, fumaça' },
  { id: 'camera_offline', title: 'Câmera Offline / Perda de Sinal', desc: 'Queda de conexão RTSP ou timeout de stream' },
  { id: 'storage_warning', title: 'Alerta de Armazenamento (StorageGuard)', desc: 'Disco atingiu 80%, 90% ou modo expurgo' }
]
</script>

<template>
  <div class="vms-flex-col">
    <div class="vms-form-group">
      <label class="vms-label">Nome do Workflow</label>
      <input
        :value="workflowName"
        class="vms-input"
        placeholder="Ex: Alerta de Invasão Noturna - Portaria"
        @input="emit('update:workflowName', ($event.target as HTMLInputElement).value)"
      />
    </div>
    <div class="vms-form-group">
      <label class="vms-label">Gatilho Inicial (Trigger)</label>
      <div class="vms-grid-container" style="grid-template-columns: 1fr; gap: 0.5rem;">
        <div
          v-for="opt in triggerOptions"
          :key="opt.id"
          class="vms-card vms-card-elevated"
          :style="{
            padding: '0.75rem 1rem',
            cursor: 'pointer',
            borderColor: triggerType === opt.id ? 'var(--vms-primary)' : 'var(--vms-border-light)'
          }"
          @click="emit('update:triggerType', opt.id)"
        >
          <div class="vms-flex-between">
            <span class="vms-text-sm vms-font-semibold" :style="{ color: triggerType === opt.id ? 'var(--vms-primary)' : 'var(--vms-text-main)' }">
              {{ opt.title }}
            </span>
            <input type="radio" :checked="triggerType === opt.id" name="trigger" class="vms-checkbox" />
          </div>
          <span class="vms-text-xs vms-text-muted" style="margin-top: 0.25rem;">{{ opt.desc }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
