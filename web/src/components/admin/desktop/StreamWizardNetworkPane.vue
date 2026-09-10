<script setup lang="ts">
import type { FolderNode } from '../../../types/streamTree'
import StreamWizardSnapshotPane from './StreamWizardSnapshotPane.vue'
import RtspProtocolForm from './wizard/RtspProtocolForm.vue'
import OnvifProtocolForm from './wizard/OnvifProtocolForm.vue'
import RtmpProtocolForm from './wizard/RtmpProtocolForm.vue'

defineProps<{
  form: {
    name: string; protocol: 'RTSP' | 'RTMP' | 'ONVIF'; url: string; path: string;
    ip: string; port: number; user: string; pass: string; streamKey: string; folderId: string;
  }
  folders: FolderNode[]
  isTesting: boolean
  hasSnapshot: boolean
  snapshotUrl?: string
  authRequired?: boolean
  detectedCodec: string
  detectedResolution: string
  detectedFps: number
  latencyMs: number
}>()

const emit = defineEmits<{
  (e: 'test'): void
  (e: 'resetSnapshot'): void
}>()
</script>

<template>
  <div style="display: flex; flex-direction: row; gap: 1.5rem; width: 100%; height: 380px; align-items: stretch;">
    <div class="vms-flex-col" style="flex: 1; min-width: 0; justify-content: space-between; height: 380px; gap: 0.65rem;">
      <div class="vms-flex-col" style="gap: 0.65rem;">
        <div class="vms-flex-row" style="gap: 0.5rem;">
          <div class="vms-form-group" style="width: 120px;">
            <label class="vms-label">Protocolo</label>
            <select v-model="form.protocol" class="vms-auth-input" @change="emit('resetSnapshot')">
              <option value="RTSP">RTSP</option>
              <option value="ONVIF">ONVIF</option>
              <option value="RTMP">RTMP</option>
            </select>
          </div>
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Nome da Camera</label>
            <input v-model="form.name" class="vms-auth-input" placeholder="Ex: Portaria Leste" autofocus />
          </div>
        </div>

        <RtspProtocolForm v-if="form.protocol === 'RTSP'" v-model:url="form.url" v-model:ip="form.ip" v-model:port="form.port" v-model:path="form.path" v-model:user="form.user" v-model:pass="form.pass" />
        <OnvifProtocolForm v-else-if="form.protocol === 'ONVIF'" v-model:ip="form.ip" v-model:port="form.port" v-model:user="form.user" v-model:pass="form.pass" />
        <RtmpProtocolForm v-else-if="form.protocol === 'RTMP'" v-model:stream-key="form.streamKey" />

        <div class="vms-flex-row" style="gap: 0.5rem; align-items: flex-end;">
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Pasta Destino</label>
            <select v-model="form.folderId" class="vms-auth-input">
              <option value="">[RAIZ] Sem Pasta (Área Principal)</option>
              <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
            </select>
          </div>
          <button class="vms-btn vms-btn-primary" style="height: 38px; padding: 0 1.5rem; white-space: nowrap;" :disabled="isTesting" @click="emit('test')">
            <span v-if="isTesting">TESTANDO...</span>
            <span v-else>TESTAR</span>
          </button>
        </div>
      </div>
    </div>

    <StreamWizardSnapshotPane 
      :is-testing="isTesting" 
      :has-snapshot="hasSnapshot" 
      :snapshot-url="snapshotUrl"
      :auth-required="authRequired"
      :detected-codec="detectedCodec" 
      :detected-resolution="detectedResolution" 
      :detected-fps="detectedFps" 
      :latency-ms="latencyMs" 
    />
  </div>
</template>
