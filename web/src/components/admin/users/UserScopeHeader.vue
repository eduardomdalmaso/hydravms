<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { UserItem, UserRole } from '../../../types/userTree'
import { createDefaultUserModules } from '../../../data/defaultUserModules'
import { getAvailableTimezones, detectLocalTimezone } from '../../../utils/timezones'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'roleChange', newRole: UserRole): void }>()

const availableTimezones = ref<string[]>([])

onMounted(() => {
  availableTimezones.value = getAvailableTimezones()
  if (!props.user.timezone) {
    props.user.timezone = detectLocalTimezone()
  }
})

const handleRoleSelect = (role: UserRole) => {
  props.user.role = role
  props.user.modules = createDefaultUserModules(role, props.user.companyScope)
  emit('roleChange', role)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.65rem; background: rgba(255, 255, 255, 0.02); border: 1px solid var(--vms-border); border-radius: 8px; padding: 0.75rem 0.85rem;">
    <!-- Profile & Role Selector -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-dim vms-text-2xs">PERFIL & DADOS DE CONTATO</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">
          {{ user.fullName || user.username }}
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.35rem;">
        <button
          v-for="r in ([
            { id: 'company_admin' as UserRole, label: 'GESTOR' },
            { id: 'operator' as UserRole, label: 'OPERADOR' },
            { id: 'client_viewer' as UserRole, label: 'VISUALIZADOR' }
          ])"
          :key="r.id"
          class="vms-btn vms-btn-sm"
          :class="user.role === r.id ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="font-size: 9px; padding: 2px 7px;"
          @click="handleRoleSelect(r.id)"
        >
          {{ r.label }}
        </button>
      </div>
    </div>

    <!-- Useful Contact & Localization Grid -->
    <div class="vms-flex-col" style="gap: 0.45rem; border-top: 1px dashed var(--vms-border); padding-top: 0.5rem;">
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <div class="vms-form-group" style="flex: 1.2;">
          <label class="vms-label" style="font-size: 9px;">E-mail</label>
          <input v-model="user.email" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;" placeholder="email@empresa.com" />
        </div>
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Fuso Horário</label>
          <select v-model="user.timezone" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;">
            <option v-for="tz in availableTimezones" :key="tz" :value="tz">{{ tz }}</option>
          </select>
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem;">
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Telefone 1 (Principal)</label>
          <input v-model="user.phonePrimary" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;" placeholder="(00) 00000-0000" />
        </div>
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Telefone 2 (Opcional)</label>
          <input v-model="user.phoneSecondary" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;" placeholder="(00) 0000-0000" />
        </div>
      </div>
    </div>
  </div>
</template>

