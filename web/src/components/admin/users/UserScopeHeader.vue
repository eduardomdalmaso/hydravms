<script setup lang="ts">
import type { UserItem, UserRole } from '../../../types/userTree'
import { createDefaultUserModules } from '../../../data/defaultUserModules'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'roleChange', newRole: UserRole): void }>()

const handleRoleSelect = (role: UserRole) => {
  props.user.role = role
  props.user.modules = createDefaultUserModules(role, props.user.companyScope)
  emit('roleChange', role)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.65rem; background: rgba(255, 255, 255, 0.02); border: 1px solid var(--vms-border); border-radius: 8px; padding: 0.75rem 0.85rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-dim vms-text-2xs">NIVEL HIERARQUICO GLOBAL // RBAC RAIZ</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">
          {{ user.role === 'admin_master' ? 'ADMIN MASTER // PLATAFORMA TOTAL' : (user.role === 'company_admin' ? 'GESTOR DE EMPRESA // TENANT' : (user.role === 'operator' ? 'OPERADOR DE MONITORAMENTO' : 'CLIENTE FINAL // APENAS VISUALIZACAO')) }}
        </span>
      </div>
      <span class="vms-badge" :class="user.role === 'admin_master' ? 'vms-badge-online' : (user.role === 'company_admin' ? 'vms-badge-warning' : 'vms-badge-neutral')" style="font-size: 9px;">
        {{ user.role === 'admin_master' ? '[ADMIN CENTER: TOTAL]' : (user.role === 'company_admin' ? '[ADMIN CENTER: EMPRESA]' : '[ADMIN CENTER: BLOQUEADO]') }}
      </span>
    </div>

    <!-- Role Switcher Pills -->
    <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
      <button
        v-for="r in ([
          { id: 'admin_master' as UserRole, label: 'ADMIN MASTER' },
          { id: 'company_admin' as UserRole, label: 'EMPRESA / TENANT' },
          { id: 'operator' as UserRole, label: 'OPERADOR' },
          { id: 'client_viewer' as UserRole, label: 'CLIENTE FINAL' }
        ])"
        :key="r.id"
        class="vms-btn vms-btn-sm"
        :class="user.role === r.id ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-size: 10px; padding: 3px 8px;"
        @click="handleRoleSelect(r.id)"
      >
        {{ r.label }}
      </button>
    </div>

    <!-- Company / Tenant Scope Bar -->
    <div class="vms-flex-between" style="border-top: 1px dashed var(--vms-border); padding-top: 0.5rem;">
      <span class="vms-text-dim vms-text-2xs">ESCOPO DE EMPRESA / CLIENTE VINCULADO:</span>
      <span class="vms-text-mono vms-text-2xs vms-font-semibold" style="color: var(--vms-text-regular);">
        {{ user.companyScope }}
      </span>
    </div>
  </div>
</template>
