<script setup lang="ts">
import { ref } from 'vue'
import type { LoginCredentials } from '../../types/auth'

const props = defineProps<{ isLoading: boolean; errorMessage?: string | null }>()
const emit = defineEmits<{ (e: 'submit', creds: LoginCredentials): void }>()

const username = ref('admin')
const password = ref('')
const rememberMe = ref(true)

const handleSubmit = () => {
  emit('submit', { username: username.value, password: password.value, remember_me: rememberMe.value })
}
</script>

<template>
  <form class="vms-flex-col" style="gap: 1.25rem;" @submit.prevent="handleSubmit">
    <div class="vms-form-group">
      <label class="vms-label">Usuario</label>
      <input v-model="username" type="text" required class="vms-auth-input" placeholder="admin" autofocus />
    </div>
    <div class="vms-form-group">
      <div class="vms-flex-between">
        <label class="vms-label">Senha</label>
        <a href="#" class="vms-text-2xs" style="color: var(--vms-neu-accent-orange);">Esqueceu a senha?</a>
      </div>
      <input v-model="password" type="password" required class="vms-auth-input" placeholder="••••••••••••" />
    </div>
    <div class="vms-flex-between">
      <label class="vms-checkbox-label">
        <input v-model="rememberMe" type="checkbox" class="vms-checkbox" />
        <span class="vms-text-xs">Lembrar neste navegador</span>
      </label>
    </div>
    <div v-if="errorMessage" class="vms-badge vms-badge-recording" style="justify-content: center; width: 100%;">
      {{ errorMessage }}
    </div>
    <button type="submit" class="vms-auth-submit-btn vms-neu-accent-hero" :disabled="isLoading">
      {{ isLoading ? 'Autenticando...' : 'ENTRAR NO SISTEMA' }}
    </button>
  </form>
</template>
