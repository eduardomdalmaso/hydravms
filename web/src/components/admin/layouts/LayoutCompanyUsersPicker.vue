<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ companyName: string; allowedUserIds: string[] }>()
const emit = defineEmits<{ (e: 'toggleUser', userId: string): void }>()

const companyUsers = computed(() => [
  { id: 'usr_02', name: 'gestor_alpha_seguranca', role: 'GESTOR EMPRESA' },
  { id: 'usr_03', name: 'operador_portaria', role: 'OPERADOR CENTRAL' },
  { id: 'usr_05', name: 'auditor_externo_ti', role: 'AUDITOR DE TI' }
])

const allowedInCompanyCount = computed(() =>
  companyUsers.value.filter(u => props.allowedUserIds?.includes(u.id)).length
)
</script>

<template>
  <div class="vms-form-group" style="background: #07080c !important; padding: 0.65rem; border-radius: 6px; border: 1px solid var(--vms-border);">
    <div class="vms-flex-between" style="align-items: center; margin-bottom: 0.45rem;">
      <span class="vms-text-mono vms-text-2xs vms-font-bold" style="color: var(--vms-neu-accent-orange);">
        // USUARIOS COM PERMISSAO NA EMPRESA
      </span>
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border); font-size: 8px;">
        {{ allowedInCompanyCount }} / {{ companyUsers.length }} LIBERADOS
      </span>
    </div>

    <!-- Div system with user info on left and checkbox on right -->
    <div class="vms-flex-col" style="gap: 0.35rem; max-height: 140px; overflow-y: auto;">
      <div
        v-for="u in companyUsers"
        :key="u.id"
        class="vms-flex-between"
        :style="{
          padding: '6px 10px',
          borderRadius: '4px',
          background: allowedUserIds.includes(u.id) ? '#14171d' : 'rgba(255, 255, 255, 0.02)',
          border: allowedUserIds.includes(u.id) ? '1px solid rgba(255, 94, 58, 0.35)' : '1px solid var(--vms-border)',
          cursor: 'pointer',
          alignItems: 'center'
        }"
        @click="emit('toggleUser', u.id)"
      >
        <div class="vms-flex-col" style="gap: 1px;">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">{{ u.name }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">[{{ u.role }}]</span>
        </div>

        <!-- Checkbox placed on the right, replacing the text description -->
        <input
          type="checkbox"
          class="vms-checkbox"
          :checked="allowedUserIds.includes(u.id)"
          @click.stop
          @change="emit('toggleUser', u.id)"
        />
      </div>
    </div>
  </div>
</template>
