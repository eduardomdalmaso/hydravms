<script setup lang="ts">
import type { EnterpriseWorkflowItem, FlowNodeItem } from '../../../types/workflowTree'
import FlowNodeConfigFields from './FlowNodeConfigFields.vue'

defineProps<{
  workflow: EnterpriseWorkflowItem
  selectedNode: FlowNodeItem | null
  isLocked?: boolean
}>()
</script>

<template>
  <div class="vms-flow-options-pane">
    <div class="vms-options-header">
      <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">
        {{ selectedNode?.kind === 'TERMINAL_STOP' ? '[PARE // FIM]' : (selectedNode ? `PARAMETROS // ${selectedNode.title}` : 'CONFIGURACOES GLOBAIS') }}
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">
        {{ selectedNode?.kind === 'TERMINAL_STOP' ? 'BLOCO SEM PARAMETROS' : (selectedNode ? `[${selectedNode.type}] // ID: ${selectedNode.id.slice(-6)}` : 'WORKFLOW SELECIONADO') }}
      </span>
    </div>

    <div class="vms-options-body">
      <!-- 1. Bloco Terminal PARE: SEM OPÇÕES / FORMULÁRIO -->
      <template v-if="selectedNode?.kind === 'TERMINAL_STOP'">
        <div style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 14px; display: flex; flex-direction: column; align-items: center; gap: 8px; text-align: center; margin-top: 10px;">
          <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><rect x="6" y="6" width="12" height="12" rx="2"/></svg>
          <span class="vms-font-bold vms-text-xs" style="color: #ffffff;">BLOCO TERMINAL (PARE)</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">Este bloco apenas encerra o fluxo e não possui opções ou parâmetros configuráveis.</span>
        </div>
      </template>

      <!-- 2. Nó Normal com Parâmetros -->
      <template v-else-if="selectedNode">
        <div class="vms-form-group">
          <label class="vms-label">TITULO DO BLOCO</label>
          <input v-model="selectedNode.title" class="vms-auth-input" :disabled="isLocked" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">RESUMO / OBSERVACAO</label>
          <input v-model="selectedNode.summary" class="vms-auth-input" :disabled="isLocked" />
        </div>
        <FlowNodeConfigFields :node="selectedNode" :is-locked="isLocked" />
      </template>

      <!-- 3. Parâmetros Globais do Workflow (Quando nenhum nó estiver selecionado) -->
      <template v-else>
        <div class="vms-form-group">
          <label class="vms-label">NOME DA INTEGRACAO</label>
          <input v-model="workflow.name" class="vms-auth-input" :disabled="isLocked" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">ESCOPO EMPRESA / CLIENTE</label>
          <input v-model="workflow.companyScope" class="vms-auth-input" :disabled="isLocked" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">STATUS OPERACIONAL</label>
          <select v-model="workflow.status" class="vms-auth-input" :disabled="isLocked">
            <option value="ATIVO">[ATIVO]</option>
            <option value="PAUSADO">[PAUSADO]</option>
          </select>
        </div>
        <div style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 8px;">
          <span class="vms-text-dim vms-text-2xs">// DICA: Clique em qualquer bloco no centro do Flow para parametrizar seus atributos.</span>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.vms-flow-options-pane { width: 300px; background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg); display: flex; flex-direction: column; overflow: hidden; }
.vms-options-header { padding: 10px 12px; background: #07080c; border-bottom: 1px solid var(--vms-border); display: flex; flex-direction: column; gap: 2px; }
.vms-options-body { padding: 12px; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; flex: 1; }
</style>
