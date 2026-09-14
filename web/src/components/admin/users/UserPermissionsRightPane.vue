<script setup lang="ts">
import type { UserItem } from '../../../types/userTree'
import UserScopeHeader from './UserScopeHeader.vue'
import UserModuleAccordion from './UserModuleAccordion.vue'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const handleRoleChange = (role: string) => {
  emit('saved', `[RBAC] Nível do usuário alterado para ${role.toUpperCase()}. Permissões sincronizadas.`)
}
</script>

<template>
  <div class="vms-split-pane" style="gap: 0.85rem;">
    <!-- Root Role & Company Scope Header -->
    <UserScopeHeader :user="user" @role-change="handleRoleChange" />

    <!-- Application Scope & Modules Section -->
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">ESCOPO DA APLICAÇÃO // MÓDULOS & REGRAS</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// ÁRVORE DE ACESSO A RECURSOS NA PLATAFORMA</span>
      </div>
    </div>

    <!-- Tree of Application Modules & Sub-permissions -->
    <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 0.6rem; padding-right: 2px;">
      <UserModuleAccordion
        v-for="mod in (user.modules || [])"
        :key="mod.id"
        :module="mod"
        @toggle="(msg) => emit('saved', msg)"
      />
    </div>
  </div>
</template>

