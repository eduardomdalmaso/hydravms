<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{ isLoading: boolean; errorMessage?: string | null }>()
const emit = defineEmits<{ (e: 'verify', code: string): void; (e: 'back'): void }>()

const code = ref('')

const handleVerify = () => {
  emit('verify', code.value)
}
</script>

<template>
  <form class="vms-flex-col" style="gap: 1.25rem;" @submit.prevent="handleVerify">
    <div class="vms-flex-col" style="text-align: center; gap: 0.25rem;">
      <h3 class="vms-h3" style="color: #fff;">Autenticação em Duas Etapas</h3>
      <p class="vms-text-xs vms-text-muted">Digite o código de 6 dígitos gerado pelo app autenticador (TOTP).</p>
    </div>
    <div class="vms-form-group">
      <input
        v-model="code"
        type="text"
        maxlength="6"
        required
        class="vms-auth-input vms-text-mono"
        style="text-align: center; font-size: 24px; letter-spacing: 0.5em; font-weight: 700;"
        placeholder="000000"
        autofocus
      />
    </div>
    <div v-if="errorMessage" class="vms-badge vms-badge-recording" style="justify-content: center; width: 100%;">
      {{ errorMessage }}
    </div>
    <button type="submit" class="vms-auth-submit-btn vms-neu-accent-hero" :disabled="isLoading">
      {{ isLoading ? 'Verificando...' : 'CONFIRMAR CÓDIGO' }}
    </button>
    <button type="button" class="vms-btn vms-btn-ghost vms-btn-sm" style="width: 100%;" @click="emit('back')">
      [VOLTAR]
    </button>
  </form>
</template>
