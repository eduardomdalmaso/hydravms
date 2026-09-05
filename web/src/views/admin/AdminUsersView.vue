<script setup lang="ts">
import { ref } from 'vue'
import type { AdminUser } from '../../types/admin'

const users = ref<AdminUser[]>([
  { id: 'usr_01', username: 'admin', role: 'admin', email: 'admin@hydravms.internal', allowed_cameras: ['ALL'], is_active: true, created_at: '2026-01-10' },
  { id: 'usr_02', username: 'operador_portaria', role: 'operator', email: 'portaria@hydravms.internal', allowed_cameras: ['cam_01', 'cam_02'], is_active: true, created_at: '2026-02-15' },
  { id: 'usr_03', username: 'supervisor_seguranca', role: 'supervisor', email: 'supervisor@hydravms.internal', allowed_cameras: ['ALL'], is_active: true, created_at: '2026-03-01' }
])

const isModalOpen = ref(false)
const newUsername = ref('')
const newRole = ref<'operator' | 'supervisor'>('operator')

const handleAddUser = () => {
  if (!newUsername.value) return
  users.value.push({
    id: `usr_0${users.value.length + 1}`,
    username: newUsername.value,
    role: newRole.value,
    email: `${newUsername.value}@hydravms.internal`,
    allowed_cameras: ['cam_01'],
    is_active: true,
    created_at: '2026-09-03'
  })
  isModalOpen.value = false
  newUsername.value = ''
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">GESTAO DE USUARIOS & PERMISSOES (RBAC)</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PAGINA DE CONTROLE DE ACESSO AO VMS OPERACIONAL</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVO USUARIO]</button>
    </div>

    <!-- Users Table -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead>
          <tr><th>USUARIO</th><th>PERFIL</th><th>EMAIL</th><th>CAMERAS LIBERADAS</th><th>STATUS</th><th>ACOES</th></tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td class="vms-font-semibold" style="color: #fff;">{{ u.username }}</td>
            <td><span class="vms-badge" :class="u.role === 'admin' ? 'vms-badge-recording' : 'vms-badge-info'">[{{ u.role.toUpperCase() }}]</span></td>
            <td class="vms-text-mono vms-text-2xs vms-text-dim">{{ u.email }}</td>
            <td class="vms-text-mono vms-text-xs">{{ u.allowed_cameras.join(', ') }}</td>
            <td style="text-align: center;"><span class="vms-status-led" :class="u.is_active ? 'online' : 'offline'" :title="u.is_active ? 'Ativo' : 'Inativo'"></span></td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Novo Usuario -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CADASTRAR NOVO USUARIO DO SISTEMA</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="isModalOpen = false">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group">
            <label class="vms-label">Nome de Usuario</label>
            <input v-model="newUsername" class="vms-auth-input" placeholder="Ex: operador_galpao" />
          </div>
          <div class="vms-form-group">
            <label class="vms-label">Nivel de Acesso</label>
            <div class="vms-flex-row" style="gap: 1rem;">
              <button class="vms-btn vms-btn-sm" :class="newRole === 'operator' ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="newRole = 'operator'">OPERADOR</button>
              <button class="vms-btn vms-btn-sm" :class="newRole === 'supervisor' ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="newRole = 'supervisor'">SUPERVISOR</button>
            </div>
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddUser">SALVAR</button>
        </div>
      </div>
    </div>
  </div>
</template>
