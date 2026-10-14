<script setup lang="ts">
import type { FolderNode } from '../../../types/streamTree'
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
}>()

const emit = defineEmits<{
  (e: 'resetSnapshot'): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.85rem; width: 100%; max-width: 680px; margin: 0 auto; box-sizing: border-box;">
    <div class="vms-flex-row" style="gap: 0.5rem;">
      <div class="vms-form-group" style="width: 140px; flex-shrink: 0;">
        <label class="vms-label">Protocolo</label>
        <select v-model="form.protocol" class="vms-select" style="height: 38px; padding: 0 0.75rem; font-size: 13px;" @change="emit('resetSnapshot')">
          <option value="RTSP">RTSP</option>
          <option value="ONVIF">ONVIF</option>
          <option value="RTMP">RTMP</option>
          <option value="LOOP">LOOP (Arquivo)</option>
        </select>
      </div>
      <div class="vms-form-group" style="flex: 1; min-width: 160px;">
        <label class="vms-label">Nome da Camera</label>
        <input v-model="form.name" class="vms-auth-input" style="height: 38px; box-sizing: border-box;" placeholder="Ex: Portaria Leste" autofocus />
      </div>
    </div>

    <RtspProtocolForm v-if="form.protocol === 'RTSP'" v-model:url="form.url" v-model:ip="form.ip" v-model:port="form.port" v-model:path="form.path" v-model:user="form.user" v-model:pass="form.pass" />
    <OnvifProtocolForm v-else-if="form.protocol === 'ONVIF'" v-model:ip="form.ip" v-model:port="form.port" v-model:user="form.user" v-model:pass="form.pass" />
    <RtmpProtocolForm v-else-if="form.protocol === 'RTMP'" v-model:stream-key="form.streamKey" />
    <LoopProtocolForm v-else-if="form.protocol === 'LOOP'" v-model:url="form.url" />

    <div class="vms-form-group" style="width: 100%; margin-bottom: 0;">
      <label class="vms-label">Pasta Destino</label>
      <select v-model="form.folderId" class="vms-select" style="height: 38px; padding: 0 0.75rem; font-size: 13px; width: 100%;">
        <option value="">[RAIZ] Sem Pasta (Área Principal)</option>
        <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
      </select>
    </div>
  </div>
</template>
