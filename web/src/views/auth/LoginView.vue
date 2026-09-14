<script setup lang="ts">
import type { LoginCredentials } from '../../types/auth'
import LoginForm from '../../components/auth/LoginForm.vue'
import { useBranding } from '../../composables/useBranding'

const props = defineProps<{
  isLoading: boolean
  errorMessage?: string | null
}>()

const emit = defineEmits<{ (e: 'login', creds: LoginCredentials): void }>()
const { branding } = useBranding()
</script>

<template>
  <div class="vms-auth-container">
    <div class="vms-auth-card">
      <div class="vms-auth-logo-circle" :style="{ background: 'rgba(255, 94, 58, 0.12)', border: `1px solid ${branding.brandColor || 'rgba(255, 94, 58, 0.35)'}`, display: 'flex', alignItems: 'center', justifyContent: 'center', width: '64px', height: '64px', borderRadius: '50%' }">
        <img :src="branding.loginLogo || '/hydra.svg'" alt="Logo" style="width: 38px; height: 38px; object-fit: contain;" />
      </div>
      <div class="vms-flex-col" style="text-align: center; gap: 0.25rem;">
        <h1 class="vms-h1" style="font-size: 20px;">{{ branding.systemName || 'HYDRA VMS' }}</h1>
        <span class="vms-text-xs vms-text-muted">{{ branding.slogan || 'Video Management & AI Studio' }}</span>
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
