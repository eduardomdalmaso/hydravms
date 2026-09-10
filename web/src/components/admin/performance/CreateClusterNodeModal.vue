<script setup lang="ts">
import { ref } from 'vue'
import { probeClusterNode } from '../../../services/adminApi'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', nodeData: { node_name: string; node_role: string; ip_address: string; http_port: number; grpc_port: number; webrtc_port: number; gpu_device_info?: string }): void
}>()

const nodeName = ref(''), nodeRole = ref<'edge_ingest' | 'gpu_worker' | 'control_plane'>('edge_ingest')
const ipAddress = ref('127.0.0.1'), httpPort = ref(8080), grpcPort = ref(50051), webrtcPort = ref(8889), gpuInfo = ref('')
const isProbing = ref(false), probeResult = ref<{ online: boolean; latency_ms?: number; app_name?: string; gpu_model?: string; error?: string } | null>(null)

const handleProbe = async () => {
  isProbing.value = true; probeResult.value = null
  const res = await probeClusterNode(ipAddress.value, Number(httpPort.value))
  isProbing.value = false; probeResult.value = res
  if (res.online) {
    if (res.suggested_name && !nodeName.value) nodeName.value = res.suggested_name
    if (res.suggested_role) nodeRole.value = res.suggested_role
    if (res.gpu_model) gpuInfo.value = res.gpu_model
    if (res.grpc_port) grpcPort.value = res.grpc_port
    if (res.webrtc_port) webrtcPort.value = res.webrtc_port
  }
}

const handleSubmit = () => {
  if (!nodeName.value || !ipAddress.value) return
  emit('save', {
    node_name: nodeName.value.toUpperCase(), node_role: nodeRole.value, ip_address: ipAddress.value,
    http_port: Number(httpPort.value), grpc_port: Number(grpcPort.value), webrtc_port: Number(webrtcPort.value), gpu_device_info: gpuInfo.value || undefined
  })
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="width: 460px; padding: 1.25rem; background: #0e1117; border: 1px solid var(--vms-border); border-radius: 6px;">
      <div class="vms-flex-between" style="align-items: center; border-bottom: 1px solid var(--vms-border); padding-bottom: 0.6rem; margin-bottom: 0.85rem;">
        <h3 class="vms-h3" style="color: #fff; margin: 0; font-size: 13px; font-family: var(--vms-font-jetbrains);">[CONECTAR NÓ HYDRASTREAM / FORGE]</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="emit('close')">✕</button>
      </div>

      <form class="vms-flex-col" style="gap: 0.65rem;" @submit.prevent="handleSubmit">
        <div style="display: grid; grid-template-columns: 2fr 1fr auto; gap: 0.4rem; align-items: flex-end;">
          <div>
            <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">IP DO SERVIDOR</label>
            <input v-model="ipAddress" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 12px;" placeholder="127.0.0.1" required />
          </div>
          <div>
            <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">PORTA HTTP</label>
            <input v-model.number="httpPort" type="number" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 12px;" required />
          </div>
          <button type="button" class="vms-btn vms-btn-secondary vms-btn-sm" :disabled="isProbing" style="height: 28px; font-size: 10px; padding: 0 8px;" @click="handleProbe">
            {{ isProbing ? 'TESTANDO...' : '⚡ TESTAR CONEXÃO' }}
          </button>
        </div>

        <div v-if="probeResult" class="vms-badge" :style="{ background: probeResult.online ? 'rgba(0, 255, 157, 0.1)' : 'rgba(255, 0, 60, 0.1)', color: probeResult.online ? '#00ff9d' : '#ff003c', border: `1px solid ${probeResult.online ? '#00ff9d44' : '#ff003c44'}`, padding: '6px 10px', fontSize: '10.5px' }">
          <span v-if="probeResult.online">[OK {{ probeResult.latency_ms }}ms] {{ probeResult.app_name }} {{ probeResult.gpu_model ? '// ' + probeResult.gpu_model : '' }}</span>
          <span v-else>[FALHA] {{ probeResult.error }}</span>
        </div>

        <div>
          <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">NOME DA INSTÂNCIA</label>
          <input v-model="nodeName" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 12px;" placeholder="[NODE] HYDRASTREAM-LOCAL" required />
        </div>

        <div>
          <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">PAPEL DO NÓ NO CLUSTER</label>
          <select v-model="nodeRole" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 11px;">
            <option value="edge_ingest">[HYDRASTREAM] NÓ DE INGESTÃO RTSP (DATA PLANE)</option>
            <option value="gpu_worker">[HYDRAFORGE] NÓ DE IA / ANALÍTICOS (GPU WORKER)</option>
            <option value="control_plane">[HYDRAVMS] CONTROL PLANE SECUNDÁRIO</option>
          </select>
        </div>

        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.4rem;">
          <div>
            <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">PORTA gRPC</label>
            <input v-model.number="grpcPort" type="number" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 12px;" required />
          </div>
          <div>
            <label class="vms-text-2xs" style="color: var(--vms-text-dim); display: block; margin-bottom: 3px;">PORTA WebRTC (WHEP)</label>
            <input v-model.number="webrtcPort" type="number" class="vms-auth-input" style="width: 100%; height: 28px; font-size: 12px;" required />
          </div>
        </div>

        <div class="vms-flex-row" style="gap: 0.5rem; justify-content: flex-end; margin-top: 0.4rem;">
          <button type="button" class="vms-btn vms-btn-ghost vms-btn-sm" @click="emit('close')">CANCELAR</button>
          <button type="submit" class="vms-btn vms-btn-primary vms-btn-sm">SALVAR</button>
        </div>
      </form>
    </div>
  </div>
</template>
