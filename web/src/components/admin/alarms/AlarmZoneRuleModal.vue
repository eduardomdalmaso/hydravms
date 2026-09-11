<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AlarmZoneRule } from '../../../types/alarmTree'

const props = defineProps<{ isOpen: boolean; rule?: AlarmZoneRule | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', rule: AlarmZoneRule): void }>()

const form = ref<AlarmZoneRule>({
  id: '', name: '', actionType: 'notification', isActive: true, delaySeconds: 0, targetOutput: 'NOTIF_PUSH'
})

watch(() => props.isOpen, (open) => {
  if (!open) return
  if (props.rule) form.value = { ...props.rule }
  else {
    form.value = {
      id: `ZON_0${Math.floor(Math.random() * 90) + 10}`,
      name: 'Regra de Disparo de Zona',
      actionType: 'notification',
      isActive: true,
      delaySeconds: 0,
      targetOutput: 'NOTIF_PUSH_APP'
    }
  }
})

const handleIntegerInput = (ev: Event) => {
  const val = (ev.target as HTMLInputElement).value.replace(/\D/g, '')
  form.value.delaySeconds = val === '' ? 0 : parseInt(val, 10)
}

const handleSave = () => {
  if (!form.value.name) form.value.name = 'Regra de Zona'
  emit('save', { ...form.value })
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="max-width: 540px; width: 95%; max-height: 90vh; overflow-y: auto; display: flex; flex-direction: column; gap: 1rem;">
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">CADASTRO DE REGRA DE ZONA & AÇÕES</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// ID: {{ form.id }} (DIRETRIZES DE RESPOSTA)</span>
        </div>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Nome do Cadastro / Regra</label>
        <input v-model="form.name" class="vms-auth-input" placeholder="Ex: Notificação Push + Foto de Câmera" />
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Tipo de Ação ao Disparar</label>
        <select v-model="form.actionType" class="vms-auth-input">
          <option value="notification">[NOTIFICAÇÃO] Push App, Telegram & Webhook</option>
          <option value="relay">[RELAY] Acionamento de Sirene / Portão I/O</option>
          <option value="ptz_preset">[PTZ] Mover Câmera Vinculada para Preset</option>
          <option value="recording">[GRAVAÇÃO] Forçar Gravação de Alta Taxa (4K)</option>
        </select>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Destino / Saída Vinculada</label>
        <input v-model="form.targetOutput" class="vms-auth-input" placeholder="Ex: RELAY_01 // BORNE GUARITA ou APP_MOBILE" />
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Retardo de Acionamento (segundos)</label>
        <input type="text" inputmode="numeric" :value="form.delaySeconds" class="vms-auth-input" placeholder="0 = Imediato" @input="handleIntegerInput($event)" />
      </div>

      <div class="vms-modal-footer" style="padding: 0.75rem 0 0 0; margin-top: 0.5rem; border-top: 1px solid var(--vms-border); display: flex; justify-content: flex-end; gap: 0.75rem; background: transparent;">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
