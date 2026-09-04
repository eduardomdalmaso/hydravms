<script setup lang="ts">
import type { LoginCredentials } from '../../types/auth'
import LoginForm from '../../components/auth/LoginForm.vue'

const props = defineProps<{
  isLoading: boolean
  errorMessage?: string | null
}>()

const emit = defineEmits<{ (e: 'login', creds: LoginCredentials): void }>()
</script>

<template>
  <div class="vms-auth-container">
    <div class="vms-auth-card">
      <div class="vms-auth-logo-circle">
        <span class="vms-text-mono" style="font-weight: 700; font-size: 24px;">H</span>
      </div>
      <div class="vms-flex-col" style="text-align: center; gap: 0.25rem;">
        <h1 class="vms-h1" style="font-size: 20px;">HYDRA VMS</h1>
        <span class="vms-text-xs vms-text-muted">Video Management & AI Studio</span>
      </div>
      <LoginForm
        :isLoading="isLoading"
        :errorMessage="errorMessage"
        @submit="emit('login', $event)"
      />
      <div class="vms-flex-between" style="border-top: 1px solid rgba(255,255,255,0.04); padding-top: 1rem;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">BUILD-STABLE // v1.0.0</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">TOKEN ENFORCED</span>
      </div>
    </div>
  </div>
</template>
