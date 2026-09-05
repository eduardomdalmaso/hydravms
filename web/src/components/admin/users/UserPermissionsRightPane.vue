<script setup lang="ts">
import { ref } from 'vue'
import type { UserItem, UserPermissionRule } from '../../../types/userTree'
import UserScopeHeader from './UserScopeHeader.vue'
import UserModuleAccordion from './UserModuleAccordion.vue'
import UserPermissionModal from './UserPermissionModal.vue'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const isModalOpen = ref(false)

const handleRoleChange = (role: string) => {
  emit('saved', `[RBAC] Nível do usuário alterado para ${role.toUpperCase()}. Permissões sincronizadas.`)
}

const handleSaveCustomRule = (rule: UserPermissionRule) => {
  if (props.user.modules && props.user.modules.length > 0) {
    const targetMod = props.user.modules.find(m => m.id === 'mosaic_live') || props.user.modules[0]
    targetMod.permissions.push({ id: rule.id, name: `${rule.name} // [${rule.scope}]`, isEnabled: true })
    targetMod.isEnabled = true
  }
  emit('saved', `Diretriz "${rule.name}" adicionada ao escopo do usuário.`)
}
</script>

<template>
  <div class="vms-split-pane" style="gap: 0.85rem;">
    <!-- Root Role & Company Scope Header -->
    <UserScopeHeader :user="user" @role-change="handleRoleChange" />

    <!-- Application Scope & Modules Section -->
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">ESCOPO DA APLICACAO // MODULOS & REGRAS</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// ARVORE DE ACESSO A RECURSOS NA APLICACAO TODA</span>
      </div>
      <button class="vms-btn vms-btn-primary" style="padding: 0.3rem 0.7rem; font-size: 16px; font-weight: 700; line-height: 1; display: flex; align-items: center; justify-content: center;" title="Nova Diretriz de Escopo" @click="isModalOpen = true">
        <span>+</span>
      </button>
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

    <UserPermissionModal :is-open="isModalOpen" @close="isModalOpen = false" @save="handleSaveCustomRule" />
  </div>
</template>
