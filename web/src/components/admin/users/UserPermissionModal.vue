<script setup lang="ts">
import { ref, watch } from 'vue'
import type { UserPermissionRule } from '../../../types/userTree'

const props = defineProps<{ isOpen: boolean; rule?: UserPermissionRule | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', rule: UserPermissionRule): void }>()

const form = ref<UserPermissionRule>({
  id: '', name: '', category: 'cameras', isActive: true, scope: 'CAM_01, CAM_02'
})

watch(() => props.isOpen, (open) => {
  if (!open) return
  if (props.rule) form.value = { ...props.rule }
  else {
    form.value = {
      id: `PERM_0${Math.floor(Math.random() * 90) + 10}`,
      name: 'Diretriz de Acesso de Cameras',
      category: 'cameras',
      isActive: true,
      scope: 'CAM_01, CAM_02'
    }
  }
})

const handleSave = () => {
  if (!form.value.name) form.value.name = 'Diretriz de Permissao'
  emit('save', { ...form.value })
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="max-width: 540px; width: 95%; max-height: 90vh; overflow-y: auto; display: flex; flex-direction: column; gap: 1rem;">
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">CADASTRO DE DIRETRIZ DE PERMISSAO // RBAC</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// ID: {{ form.id }} (DEFINA ESCOPO DE ACESSO)</span>
        </div>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Nome da Diretriz de Acesso</label>
        <input v-model="form.name" class="vms-auth-input" placeholder="Ex: Liberar Visualizacao Cameras Portaria" />
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Categoria de Permissao</label>
        <select v-model="form.category" class="vms-auth-input">
          <option value="cameras">[CAMERAS] Visualizacao ao Vivo de Streams</option>
          <option value="ptz">[PTZ] Controle de Joystick e Presets</option>
          <option value="export">[EXPORTACAO] Download de Gravacoes e Evidencias MP4</option>
          <option value="alarms">[ALARMES] Reconhecimento e Desarme de Sensores</option>
          <option value="system">[SISTEMA] Acesso ao Painel Administrativo</option>
        </select>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Escopo / Alvos Autorizados</label>
        <input v-model="form.scope" class="vms-auth-input" placeholder="Ex: CAM_01, CAM_02 ou ZONA_NORTE ou ALL" />
      </div>

      <div class="vms-modal-footer" style="padding: 0.75rem 0 0 0; margin-top: 0.5rem; border-top: 1px solid var(--vms-border); display: flex; justify-content: flex-end; gap: 0.75rem; background: transparent;">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
