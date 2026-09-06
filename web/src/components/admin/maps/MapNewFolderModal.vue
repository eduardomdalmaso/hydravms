<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', name: string, clientType: 'company' | 'final_client'): void
}>()

const folderName = ref('')
const clientType = ref<'company' | 'final_client'>('company')

const handleSubmit = () => {
  if (!folderName.value.trim()) return
  emit('create', folderName.value.trim(), clientType.value)
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">NOVA PASTA // EMPRESA OU CLIENTE</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">NOME DA PASTA</label>
          <input v-model="folderName" class="vms-auth-input" placeholder="Ex: EMPRESA // GAMMA MONITORAMENTO" autofocus />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">ESCOPO DE VINCULO</label>
          <div class="vms-flex-row" style="gap: 1rem;">
            <label class="vms-flex-row" style="gap: 0.4rem; align-items: center; cursor: pointer;">
              <input v-model="clientType" type="radio" value="company" />
              <span class="vms-text-xs">EMPRESA DE SEGURANCA</span>
            </label>
            <label class="vms-flex-row" style="gap: 0.4rem; align-items: center; cursor: pointer;">
              <input v-model="clientType" type="radio" value="final_client" />
              <span class="vms-text-xs">CLIENTE FINAL</span>
            </label>
          </div>
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" :disabled="!folderName.trim()" @click="handleSubmit">CRIAR PASTA</button>
      </div>
    </div>
  </div>
</template>
