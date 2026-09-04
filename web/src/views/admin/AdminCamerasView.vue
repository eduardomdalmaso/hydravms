<script setup lang="ts">
import { ref } from 'vue'
import type { RegisteredCamera } from '../../types/admin'

const cameras = ref<RegisteredCamera[]>([
  { id: 'cam_01', name: 'Portaria Principal (Entrada)', rtsp_url: 'rtsp://192.168.1.101:554/live', ip: '192.168.1.101', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, has_ptz: true, status: 'online', group: 'Acessos' },
  { id: 'cam_02', name: 'Estacionamento Visitantes', rtsp_url: 'rtsp://192.168.1.102:554/live', ip: '192.168.1.102', port: 554, codec: 'H.264', resolution: '1080P', fps: 25, has_ptz: false, status: 'online', group: 'Estacionamento' },
  { id: 'cam_03', name: 'Corredor de Cargas & Docas', rtsp_url: 'rtsp://192.168.1.103:554/live', ip: '192.168.1.103', port: 554, codec: 'H.265', resolution: '1080P', fps: 30, has_ptz: false, status: 'recording', group: 'Galpao' },
  { id: 'cam_04', name: 'Perimetro dos Fundos', rtsp_url: 'rtsp://192.168.1.104:554/live', ip: '192.168.1.104', port: 554, codec: 'H.265', resolution: '4K', fps: 30, has_ptz: true, status: 'recording', group: 'Perimetro' }
])

const isModalOpen = ref(false)
const newCamName = ref('')
const newCamRtsp = ref('rtsp://')

const handleAddCamera = () => {
  if (!newCamName.value) return
  cameras.value.push({
    id: `cam_0${cameras.value.length + 1}`,
    name: newCamName.value,
    rtsp_url: newCamRtsp.value,
    ip: '192.168.1.110',
    port: 554,
    codec: 'H.265',
    resolution: '1080P',
    fps: 30,
    has_ptz: false,
    status: 'online',
    group: 'Geral'
  })
  isModalOpen.value = false
  newCamName.value = ''
  newCamRtsp.value = 'rtsp://'
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE CAMERAS // RTSP & ONVIF</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PAGINA DE GESTAO DE FLUXOS E HARDWARE (TOTAL: {{ cameras.length }})</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVA CAMERA]</button>
    </div>

    <!-- Cameras Table -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead>
          <tr><th>ID</th><th>NOME</th><th>URL RTSP</th><th>CODEC</th><th>RESOLUCAO</th><th>STATUS</th><th>ACOES</th></tr>
        </thead>
        <tbody>
          <tr v-for="cam in cameras" :key="cam.id">
            <td class="vms-text-mono">{{ cam.id }}</td>
            <td class="vms-font-semibold" style="color: #fff;">{{ cam.name }}</td>
            <td class="vms-text-mono vms-text-2xs vms-text-dim">{{ cam.rtsp_url }}</td>
            <td><span class="vms-badge vms-badge-info">{{ cam.codec }}</span></td>
            <td class="vms-text-mono vms-text-xs">{{ cam.resolution }} @ {{ cam.fps }}fps</td>
            <td><span class="vms-badge vms-badge-online">[{{ cam.status.toUpperCase() }}]</span></td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Nova Camera -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CADASTRAR NOVA CAMERA RTSP</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="isModalOpen = false">[X]</button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group">
            <label class="vms-label">Nome da Camera / Local</label>
            <input v-model="newCamName" class="vms-auth-input" placeholder="Ex: Recepcao Bloco B" />
          </div>
          <div class="vms-form-group">
            <label class="vms-label">URL do Fluxo RTSP / ONVIF</label>
            <input v-model="newCamRtsp" class="vms-auth-input" />
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddCamera">SALVAR CAMERA</button>
        </div>
      </div>
    </div>
  </div>
</template>
