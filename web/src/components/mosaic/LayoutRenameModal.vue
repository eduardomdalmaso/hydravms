<script setup lang="ts">
import { ref } from 'vue'
import type { CustomLayout } from '../../types/mosaic'

const props = defineProps<{ layout: CustomLayout }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', newName: string): void
}>()

const currentName = ref(props.layout.name)

const handleSave = () => {
  if (currentName.value.trim()) {
    emit('save', currentName.value.trim())
  }
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 380px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">RENOMEAR LAYOUT DO OPERADOR</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="emit('close')">[X]</button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome do Layout</label>
          <input
            v-model="currentName"
            class="vms-auth-input"
            autofocus
            placeholder="Ex: Minha Visao Portaria"
            @keyup.enter="handleSave"
          />
        </div>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          GRADE: {{ layout.grid.toUpperCase() }} // PERSISTIDO NO PERFIL
        </span>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR NOME</button>
      </div>
    </div>
  </div>
</template>
