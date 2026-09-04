<script setup lang="ts">
import type { UserTenant } from '../../types/auth'

const props = defineProps<{ tenants: UserTenant[] }>()
const emit = defineEmits<{ (e: 'select', tenant: UserTenant): void }>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-col" style="text-align: center; gap: 0.25rem;">
      <h3 class="vms-h3" style="color: #fff;">Selecione o Tenant de Acesso</h3>
      <p class="vms-text-xs vms-text-muted">Sua conta possui acesso a múltiplos ambientes corporativos.</p>
    </div>
    <div class="vms-flex-col" style="gap: 0.75rem;">
      <div
        v-for="t in tenants"
        :key="t.id"
        class="vms-tenant-card"
        @click="emit('select', t)"
      >
        <div class="vms-flex-col" style="gap: 0.125rem;">
          <span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ t.name }}</span>
          <span class="vms-text-xs vms-text-muted">{{ t.slug }}.hydravms.corp</span>
        </div>
        <span class="vms-badge" :class="t.role === 'admin' ? 'vms-badge-online' : 'vms-badge-info'">
          [{{ t.role.toUpperCase() }}]
        </span>
      </div>
    </div>
  </div>
</template>
