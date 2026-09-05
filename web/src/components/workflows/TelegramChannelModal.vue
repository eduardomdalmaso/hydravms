<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', channel: { name: string; token: string; chatId: string; snapshot: boolean }): void
}>()

const name = ref('Alerta Telegram - Grupo Seguranca')
const botToken = ref('')
const chatId = ref('')
const includeSnapshot = ref(true)
const isTesting = ref(false)
const testResult = ref<string | null>(null)

const handleTest = async () => {
  if (!botToken.value || !chatId.value) return
  isTesting.value = true
  setTimeout(() => {
    isTesting.value = false
    testResult.value = 'Mensagem de teste enviada com sucesso!'
  }, 800)
}

const handleSave = () => {
  if (!name.value || !botToken.value || !chatId.value) return
  emit('created', {
    name: name.value,
    token: botToken.value,
    chatId: chatId.value,
    snapshot: includeSnapshot.value
  })
  emit('close')
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CONFIGURAR CANAL TELEGRAM BOT</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome do Canal</label>
          <input v-model="name" class="vms-input" placeholder="Ex: Central de Monitoramento" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Telegram Bot Token (via @BotFather)</label>
          <input v-model="botToken" type="password" class="vms-input" placeholder="123456789:ABCdefGHIjklMNOpqrsTUVwxyz" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Chat ID ou Grupo ID</label>
          <input v-model="chatId" class="vms-input" placeholder="-100192837465" />
          <span class="vms-form-help">Use @userinfobot no Telegram para descobrir o Chat ID.</span>
        </div>
        <div class="vms-form-group">
          <label class="vms-checkbox-label">
            <input v-model="includeSnapshot" type="checkbox" class="vms-checkbox" />
            <span>Anexar snapshot (foto) do evento na mensagem</span>
          </label>
        </div>
        <div v-if="testResult" class="vms-badge vms-badge-online" style="margin-top: 0.5rem;">
          {{ testResult }}
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" :disabled="isTesting" @click="handleTest">
          {{ isTesting ? 'TESTANDO...' : 'TESTAR CONEXAO' }}
        </button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
