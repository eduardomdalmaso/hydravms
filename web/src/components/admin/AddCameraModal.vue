<script setup lang="ts">
import { ref, watch } from 'vue'
import { CAMERA_BRAND_PRESETS, buildPresetRTSP, type CameraBrandPreset } from '../../data/cameraBrandPresets'

defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'add', cam: { name: string; rtsp_url: string; protocol: string; serverId: string; switchId: string }): void }>()

const name = ref('')
const selectedBrand = ref('intelbras')
const hostIP = ref('')
const hostPort = ref(554)
const authUser = ref('admin')
const authPass = ref('')
const rtspUrl = ref('rtsp://')
const protocol = ref('RTSP')
const serverId = ref('edge_01')
const switchId = ref('sw_01')

const updateGeneratedURL = () => {
  if (selectedBrand.value === 'custom') return
  const { main } = buildPresetRTSP(selectedBrand.value, hostIP.value, hostPort.value, authUser.value, authPass.value)
  if (main) rtspUrl.value = main
}

watch([selectedBrand, hostIP, hostPort, authUser, authPass], () => {
  updateGeneratedURL()
})

const onBrandChange = () => {
  const p = CAMERA_BRAND_PRESETS.find((x: CameraBrandPreset) => x.id === selectedBrand.value)
  if (p) {
    hostPort.value = p.defaultPort
    authUser.value = p.defaultUser
    updateGeneratedURL()
  }
}

const handleSave = () => {
  if (!name.value) return
  emit('add', { name: name.value, rtsp_url: rtspUrl.value, protocol: protocol.value, serverId: serverId.value, switchId: switchId.value })
  name.value = ''
  hostIP.value = ''
  authPass.value = ''
  rtspUrl.value = 'rtsp://'
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 580px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CADASTRAR FLUXO NA TOPOLOGIA</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body" style="display: flex; flex-direction: column; gap: 0.75rem;">
        <div class="vms-form-group">
          <label class="vms-label">Fabricante / Preset</label>
          <select v-model="selectedBrand" class="vms-auth-input" @change="onBrandChange">
            <option v-for="b in CAMERA_BRAND_PRESETS" :key="b.id" :value="b.id">{{ b.label }}</option>
          </select>
        </div>
        <div class="vms-flex-row" style="gap: 0.75rem;">
          <div class="vms-form-group" style="flex: 2;">
            <label class="vms-label">Nome da Camera / Local</label>
            <input v-model="name" class="vms-auth-input" placeholder="Ex: Portaria Leste Bloco C" autofocus />
          </div>
          <div class="vms-form-group" style="flex: 1.5;">
            <label class="vms-label">Endereco IP</label>
            <input v-model="hostIP" class="vms-auth-input" placeholder="192.168.1.100" />
          </div>
        </div>
        <div class="vms-flex-row" style="gap: 0.75rem;" v-if="selectedBrand !== 'custom'">
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Usuario</label>
            <input v-model="authUser" class="vms-auth-input" />
          </div>
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Senha</label>
            <input v-model="authPass" type="password" class="vms-auth-input" placeholder="Opcional" />
          </div>
          <div class="vms-form-group" style="flex: 0.8;">
            <label class="vms-label">Porta RTSP</label>
            <input v-model.number="hostPort" type="number" class="vms-auth-input" />
          </div>
        </div>
        <div class="vms-form-group">
          <label class="vms-label">URL do Fluxo RTSP (Montada)</label>
          <input v-model="rtspUrl" class="vms-auth-input" style="font-family: var(--font-mono, monospace); font-size: 0.82rem;" />
        </div>
        <div class="vms-flex-row" style="gap: 0.75rem;">
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
