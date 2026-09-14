<script setup lang="ts">
import { ref } from 'vue'
import type { UserItem } from '../../../types/userTree'
import UserScopeHeader from './UserScopeHeader.vue'
import UserModuleAccordion from './UserModuleAccordion.vue'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const activeTab = ref<'profile' | 'permissions'>('profile')

const handleRoleChange = (role: string) => {
  emit('saved', `[RBAC] Nível do usuário alterado para ${role.toUpperCase()}. Permissões sincronizadas.`)
}
</script>

<template>
  <div class="vms-split-pane" style="gap: 0.85rem;">
    <!-- Top Tabs Header -->
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.6rem;">
      <div class="vms-flex-row" style="gap: 0.45rem;">
        <button
          class="vms-btn vms-btn-sm"
          :class="activeTab === 'profile' ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="font-size: 11px; padding: 4px 12px;"
          @click="activeTab = 'profile'"
        >
          DADOS & CONTATO
        </button>
        <button
          class="vms-btn vms-btn-sm"
          :class="activeTab === 'permissions' ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="font-size: 11px; padding: 4px 12px;"
          @click="activeTab = 'permissions'"
        >
          PERMISSÕES & MÓDULOS
        </button>
      </div>

      <span class="vms-text-mono vms-text-2xs vms-text-dim">
        {{ activeTab === 'profile' ? '// DADOS CADASTRAIS' : '// MATRIZ DE ACESSO' }}
      </span>
    </div>

    <!-- Tab 1: Profile & Contact Data -->
    <div v-if="activeTab === 'profile'" class="vms-flex-col" style="gap: 0.75rem;">
      <UserScopeHeader :user="user" @role-change="handleRoleChange" @saved="(msg) => emit('saved', msg)" />
    </div>

    <!-- Tab 2: Application Scope & Modules Tree -->
    <div v-else-if="activeTab === 'permissions'" class="vms-flex-col" style="gap: 0.6rem; flex: 1; overflow-y: auto;">
      <div class="vms-flex-between" style="padding-bottom: 0.25rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 12px;">ESCOPO DA APLICAÇÃO // MÓDULOS & REGRAS</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// ÁRVORE DE ACESSO A RECURSOS NA PLATAFORMA</span>
        </div>
      </div>

      <div style="display: flex; flex-direction: column; gap: 0.6rem; padding-right: 2px;">
        <UserModuleAccordion
          v-for="mod in (user.modules || [])"
          :key="mod.id"
          :module="mod"
          @toggle="(msg) => emit('saved', msg)"
        />
      </div>
    </div>
  </div>
</template>


