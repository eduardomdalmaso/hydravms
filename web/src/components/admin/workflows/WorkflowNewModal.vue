<script setup lang="ts">
import { ref } from 'vue'
import type { EnterpriseWorkflowItem, WorkflowFolderNode } from '../../../types/workflowTree'

const props = defineProps<{
  folders: WorkflowFolderNode[]
  currentFolderId?: string | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', payload: { wf: EnterpriseWorkflowItem; targetFolderId: string | null }): void
}>()

const wfName = ref('')
const targetFolder = ref<string>(props.currentFolderId || '')

const handleSubmit = () => {
  if (!wfName.value.trim()) return
  const folderObj = props.folders.find(f => f.id === targetFolder.value)
  const scopeName = folderObj ? folderObj.name : 'RAIZ // GLOBAL'

  const newWf: EnterpriseWorkflowItem = {
    id: `wf_${Date.now()}`,
    name: wfName.value.trim(),
    companyScope: scopeName,
    folderId: targetFolder.value || undefined,
    is_locked: false,
    status: 'ATIVO',
    cooldownSeconds: 30,
    nodes: [
      { id: `node_${Date.now()}_1`, type: 'TRIGGER', kind: 'AI_DETECT', title: 'DETECCAO IA', summary: 'Fluxo de Video Camera', position: { x: 100, y: 140 }, config: { targetClass: 'TODAS' } },
      { id: `node_${Date.now()}_2`, type: 'ACTION', kind: 'TERMINAL_STOP', title: 'PARE / FIM', summary: 'Finalizar execucao', position: { x: 420, y: 140 }, config: {} }
    ],
    edges: [],
    createdAt: new Date().toISOString().split('T')[0]
  }

  emit('create', { wf: newWf, targetFolderId: targetFolder.value || null })
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">NOVO APLICATIVO // WORKFLOW</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">NOME DO WORKFLOW</label>
          <input v-model="wfName" class="vms-auth-input" placeholder="Ex: Alerta Invasao Docas" autofocus />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">PASTA DE DESTINO</label>
          <select v-model="targetFolder" class="vms-auth-input">
            <option value="">[RAIZ // SEM PASTA]</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">
              [PASTA] {{ f.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" :disabled="!wfName.trim()" @click="handleSubmit">CRIAR APLICATIVO</button>
      </div>
    </div>
  </div>
</template>
