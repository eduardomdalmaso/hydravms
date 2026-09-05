<script setup lang="ts">
import { ref, watch } from 'vue'
import type { UserItem, UserFolderNode, UserRole } from '../../../types/userTree'

const props = defineProps<{
  isOpen: boolean
  targetFolderId?: string
  folders: UserFolderNode[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', user: Partial<UserItem>, folderId?: string): void
}>()

const username = ref('')
const fullName = ref('')
const email = ref('')
const role = ref<UserRole>('operator')
const companyScope = ref('GLOBAL // TODAS AS EMPRESAS')
const selectedFolderId = ref(props.targetFolderId || '')

watch(() => props.isOpen, (open) => {
  if (open) {
    username.value = ''
    fullName.value = ''
    email.value = ''
    role.value = 'operator'
    companyScope.value = 'EMPRESA MATRIZ // CLIENTE PADRAO'
    selectedFolderId.value = props.targetFolderId || ''
  }
})

const handleSave = () => {
  if (!username.value.trim()) return
  emit('save', {
    username: username.value.trim().toLowerCase(),
    fullName: fullName.value.trim() || 'Usuário Operador',
    email: email.value.trim() || `${username.value.trim().toLowerCase()}@hydravms.internal`,
    role: role.value,
    companyScope: companyScope.value.trim() || 'GLOBAL'
  }, selectedFolderId.value || undefined)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="max-width: 520px; width: 95%;">
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRAR NOVO USUARIO</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-flex-col" style="gap: 0.85rem; padding: 0.5rem 0;">
        <div class="vms-form-group">
          <label class="vms-label">Nome de Usuario (Login)</label>
          <input v-model="username" class="vms-auth-input" placeholder="Ex: gestor_filial_01" autofocus />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Nome Completo</label>
          <input v-model="fullName" class="vms-auth-input" placeholder="Ex: Carlos Silva de Andrade" />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Nível Hierárquico RBAC</label>
          <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
            <button v-for="r in ([{ id: 'admin_master', label: 'MASTER' }, { id: 'company_admin', label: 'EMPRESA' }, { id: 'operator', label: 'OPERADOR' }, { id: 'client_viewer', label: 'CLIENTE' }] as const)" :key="r.id" class="vms-btn vms-btn-sm" :class="role === r.id ? 'vms-btn-primary' : 'vms-btn-secondary'" style="font-size: 10px;" @click="role = r.id">{{ r.label }}</button>
          </div>
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Empresa / Clientes Vinculados</label>
          <input v-model="companyScope" class="vms-auth-input" placeholder="Ex: EMPRESA ALPHA // CLIENTE CONDOMINIO" />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Grupo / Departamento</label>
          <select v-model="selectedFolderId" class="vms-auth-input">
            <option value="">[RAIZ] Sem Grupo Definido</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">[GRUPO] {{ f.name }}</option>
          </select>
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
