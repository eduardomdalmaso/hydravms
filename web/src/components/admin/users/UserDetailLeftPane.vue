<script setup lang="ts">
import type { UserItem } from '../../../types/userTree'

defineProps<{ user: UserItem }>()
</script>

<template>
  <div class="vms-split-pane">
    <!-- Header -->
    <div class="vms-split-header">
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div style="width: 36px; height: 36px; border-radius: 8px; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center;">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: #fff; font-size: 13px;">{{ user.id.toUpperCase() }} // {{ user.username }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">DEPARTAMENTO // {{ user.groupName.toUpperCase() }}</span>
        </div>
      </div>
    </div>

    <!-- Telemetry & User Specs Grid -->
    <div class="vms-telemetry-grid" style="flex: 1; align-content: start;">
      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">PERFIL DE ACESSO (RBAC)</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">
          [{{ user.role === 'admin_master' ? 'ADMIN MASTER' : (user.role === 'company_admin' ? 'GESTOR EMPRESA' : (user.role === 'operator' ? 'OPERADOR' : 'CLIENTE FINAL')) }}]
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ACESSO AO ADMIN CENTER</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" :style="{ color: user.role === 'admin_master' ? '#00ff9d' : (user.role === 'company_admin' ? '#fcee0a' : '#ff003c') }">
          {{ user.role === 'admin_master' ? 'LIBERADO (TOTAL)' : (user.role === 'company_admin' ? 'PARCIAL (TENANT)' : 'BLOQUEADO') }}
        </span>
      </div>

      <div class="vms-telemetry-card" style="grid-column: 1 / -1;">
        <span class="vms-text-dim vms-text-2xs">ESCOPO DE EMPRESA & CLIENTES VINCULADOS</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #38bdf8;">{{ user.companyScope }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">NOME COMPLETO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ user.fullName }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">EMAIL CORPORATIVO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #38bdf8;">{{ user.email }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">STATUS DA CONTA</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">
          {{ user.isActive ? 'ATIVO (AUTORIZADO)' : 'BLOQUEADO' }}
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">DATA DE CADASTRO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ user.createdAt }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ULTIMO ACESSO AO VMS</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ user.lastLogin }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">AUTENTICACAO 2FA</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" :style="{ color: user.twoFactorEnabled ? '#00ff9d' : 'var(--vms-text-dim)' }">
          {{ user.twoFactorEnabled ? 'HABILITADO (TOTP)' : 'DESATIVADO' }}
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">NIVEL DE AUDITORIA</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #00ff9d;">TOTAL (LOG IMUTAVEL)</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">SESSAO ATIVA</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">TOKEN JWT (TLS 1.3)</span>
      </div>
    </div>
  </div>
</template>
