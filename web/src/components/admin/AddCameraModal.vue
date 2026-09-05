<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'add', cam: { name: string; rtsp_url: string; protocol: string; serverId: string; switchId: string }): void }>()

const name = ref('')
const rtspUrl = ref('rtsp://')
const protocol = ref('RTSP')
const serverId = ref('edge_01')
const switchId = ref('sw_01')

const handleSave = () => {
  if (!name.value) return
  emit('add', { name: name.value, rtsp_url: rtspUrl.value, protocol: protocol.value, serverId: serverId.value, switchId: switchId.value })
  name.value = ''
  rtspUrl.value = 'rtsp://'
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CADASTRAR FLUXO NA TOPOLOGIA</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome da Camera / Local</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: Portaria Leste Bloco C" autofocus />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">URL do Fluxo RTSP / RTMP</label>
          <input v-model="rtspUrl" class="vms-auth-input" />
        </div>
        <div class="vms-flex-row" style="gap: 1rem;">
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Servidor Edge</label>
            <select v-model="serverId" class="vms-auth-input">
              <option value="edge_01">[EDGE-SERVER-01] Galpao Principal</option>
              <option value="edge_02">[EDGE-SERVER-02] CPD Central</option>
            </select>
          </div>
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Switch PoE / Segmento</label>
            <select v-model="switchId" class="vms-auth-input">
              <option value="sw_01">[SWITCH-POE-A] Acessos</option>
              <option value="sw_02">[SWITCH-POE-B] Galpao</option>
              <option value="sw_03">[SWITCH-POE-C] CPD Central</option>
            </select>
          </div>
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
