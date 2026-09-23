<script setup lang="ts">
import type { FolderNode } from '../../../types/streamTree'
import StreamWizardSnapshotPane from './StreamWizardSnapshotPane.vue'
import RtspProtocolForm from './wizard/RtspProtocolForm.vue'
import OnvifProtocolForm from './wizard/OnvifProtocolForm.vue'
import RtmpProtocolForm from './wizard/RtmpProtocolForm.vue'
import LoopProtocolForm from './wizard/LoopProtocolForm.vue'

defineProps<{
  form: {
    name: string; protocol: 'RTSP' | 'RTMP' | 'ONVIF' | 'LOOP'; url: string; path: string;
    ip: string; port: number; user: string; pass: string; streamKey: string; folderId: string;
  }
  folders: FolderNode[]
  isTesting: boolean
  hasSnapshot: boolean
  snapshotUrl?: string
  authRequired?: boolean
  testError?: string | null
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
  <div style="display: flex; flex-direction: row; gap: 1.25rem; width: 100%; min-height: 380px; align-items: stretch; box-sizing: border-box; flex-wrap: wrap;">
    <div class="vms-flex-col" style="flex: 1 1 360px; min-width: 280px; justify-content: space-between; gap: 0.65rem; box-sizing: border-box;">
      <div class="vms-flex-col" style="gap: 0.65rem;">
        <div class="vms-flex-row" style="gap: 0.5rem;">
          <div class="vms-form-group" style="width: 130px; flex-shrink: 0;">
            <label class="vms-label">Protocolo</label>
            <select v-model="form.protocol" class="vms-select" style="height: 38px; padding: 0 0.75rem; font-size: 13px;" @change="emit('resetSnapshot')">
              <option value="RTSP">RTSP</option>
              <option value="ONVIF">ONVIF</option>
              <option value="RTMP">RTMP</option>
              <option value="LOOP">LOOP (Arquivo)</option>
            </select>
          </div>
          <div class="vms-form-group" style="flex: 1; min-width: 140px;">
            <label class="vms-label">Nome da Camera</label>
            <input v-model="form.name" class="vms-auth-input" style="height: 38px; box-sizing: border-box;" placeholder="Ex: Portaria Leste" autofocus />
          </div>
        </div>

        <RtspProtocolForm v-if="form.protocol === 'RTSP'" v-model:url="form.url" v-model:ip="form.ip" v-model:port="form.port" v-model:path="form.path" v-model:user="form.user" v-model:pass="form.pass" />
        <OnvifProtocolForm v-else-if="form.protocol === 'ONVIF'" v-model:ip="form.ip" v-model:port="form.port" v-model:user="form.user" v-model:pass="form.pass" />
        <RtmpProtocolForm v-else-if="form.protocol === 'RTMP'" v-model:stream-key="form.streamKey" />
        <LoopProtocolForm v-else-if="form.protocol === 'LOOP'" v-model:url="form.url" />

        <div class="vms-flex-row" style="gap: 0.5rem; align-items: flex-end;">
          <div class="vms-form-group" style="flex: 1; min-width: 140px; margin-bottom: 0;">
            <label class="vms-label">Pasta Destino</label>
            <select v-model="form.folderId" class="vms-select" style="height: 38px; padding: 0 0.75rem; font-size: 13px; width: 100%;">
              <option value="">[RAIZ] Sem Pasta (Área Principal)</option>
              <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
            </select>
          </div>
          <button class="vms-btn vms-btn-primary" style="height: 38px; padding: 0 1.25rem; white-space: nowrap; margin-bottom: 0; box-sizing: border-box; flex-shrink: 0;" :disabled="isTesting" @click="emit('test')">
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
      :test-error="testError"
      :detected-codec="detectedCodec" 
      :detected-resolution="detectedResolution" 
      :detected-fps="detectedFps" 
      :latency-ms="latencyMs" 
    />
  </div>
</template>
