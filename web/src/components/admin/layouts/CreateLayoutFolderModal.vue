<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', name: string, type: 'company' | 'final_client'): void }>()

const folderName = ref('')
const clientType = ref<'company' | 'final_client'>('company')

watch(() => props.isOpen, (open) => {
  if (open) { folderName.value = ''; clientType.value = 'company' }
})

const handleSave = () => {
  if (!folderName.value.trim()) return
  emit('save', folderName.value.trim().toUpperCase(), clientType.value)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 440px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">NOVA PASTA // EMPRESA OU CLIENTE</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 1rem;">
        <div class="vms-form-group">
          <label class="vms-label">TIPO DE ENTIDADE</label>
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <button class="vms-btn vms-btn-sm" :class="clientType === 'company' ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="clientType = 'company'">
              [EMPRESA]
            </button>
            <button class="vms-btn vms-btn-sm" :class="clientType === 'final_client' ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="clientType = 'final_client'">
              [CLIENTE FINAL]
            </button>
          </div>
        </div>
        <div class="vms-form-group">
          <label class="vms-label">NOME DA EMPRESA OU CLIENTE</label>
          <input v-model="folderName" class="vms-auth-input" style="text-transform: uppercase;" placeholder="EX: CONDOMINIO VISTA VERDE" autofocus @keydown.enter="handleSave" />
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
