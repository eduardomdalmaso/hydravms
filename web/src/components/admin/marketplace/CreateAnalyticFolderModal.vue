<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', folderName: string): void
}>()

const folderName = ref('')

const handleConfirm = () => {
  if (!folderName.value.trim()) return
  emit('create', folderName.value.trim())
  folderName.value = ''
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content" style="max-width: 440px; width: 100%;">
      <div class="vms-modal-header vms-flex-between">
        <h4 class="vms-h4" style="margin: 0; color: #ffffff;">CRIAR NOVA PASTA DE ANALÍTICOS</h4>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 1rem; padding: 1.25rem 0;">
        <div class="vms-flex-col" style="gap: 0.35rem;">
          <label class="vms-text-xs vms-font-semibold">NOME DA PASTA (ZONA / SETOR):</label>
          <input
            v-model="folderName"
            class="vms-auth-input"
            style="padding: 8px 12px; font-size: 12px;"
            placeholder="Ex: Portaria Principal, Setor Logística..."
            @keyup.enter="handleConfirm"
          />
        </div>
      </div>

      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.85rem;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" :disabled="!folderName.trim()" @click="handleConfirm">
          SALVAR
        </button>
      </div>
    </div>
  </div>
</template>
