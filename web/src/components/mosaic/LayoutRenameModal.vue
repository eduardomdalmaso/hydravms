<script setup lang="ts">
import { ref } from 'vue'
import type { CustomLayout } from '../../types/mosaic'

const props = defineProps<{
  layout: CustomLayout
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', newName: string): void
}>()

const nameInput = ref(props.layout.name)

const handleSave = () => {
  if (!nameInput.value.trim()) return
  emit('save', nameInput.value.trim())
  emit('close')
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 420px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">RENOMEAR MATRIZ DE VISUALIZACAO</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Novo Nome da Matriz</label>
          <input v-model="nameInput" class="vms-auth-input" placeholder="Ex: Portaria & Perimetro 2x2" autofocus @keydown.enter="handleSave" />
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
