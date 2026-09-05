<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  isOpen: boolean
  initialName?: string
  isRename?: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', name: string): void
}>()

const folderName = ref('')

watch(() => props.isOpen, (open) => {
  if (open) {
    folderName.value = (props.initialName || '').toUpperCase()
  }
})

const handleSave = () => {
  if (!folderName.value.trim()) return
  emit('save', folderName.value.trim().toUpperCase())
  folderName.value = ''
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 420px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">{{ isRename ? 'RENOMEAR PASTA' : 'NOVA PASTA DA TOPOLOGIA' }}</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome da Pasta (Caixa Alta)</label>
          <input v-model="folderName" class="vms-auth-input" style="text-transform: uppercase;" placeholder="EX: PORTARIA PRINCIPAL, GALPAO 01" autofocus @keydown.enter="handleSave" />
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
