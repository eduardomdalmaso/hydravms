<script setup lang="ts">
import { ref } from 'vue'
import { usePerformanceCluster } from '../../composables/usePerformanceCluster'
import { probeClusterNode } from '../../services/adminApi'
import ServerPerformanceRow from '../../components/admin/performance/ServerPerformanceRow.vue'
import PerformanceAlertsModal from '../../components/admin/performance/PerformanceAlertsModal.vue'
import CreateClusterNodeModal from '../../components/admin/performance/CreateClusterNodeModal.vue'

const { servers, countdown, addNode } = usePerformanceCluster()
const showConfigModal = ref(false), showCreateModal = ref(false), toastMessage = ref<string | null>(null), isAutoScanning = ref(false)

const showToast = (msg: string) => { toastMessage.value = msg; setTimeout(() => { toastMessage.value = null }, 4000) }

const handleSaveConfig = (cfg: { cpuLimit: number; vramMinPercent: number; ramLimit: number }) => {
  showConfigModal.value = false; showToast(`[ALERTAS] CPU > ${cfg.cpuLimit}% // VRAM < ${cfg.vramMinPercent}%`)
}

const handleSaveNode = async (data: any) => {
  const ok = await addNode(data)
  showCreateModal.value = false
  showToast(ok ? `[CLUSTER] Servidor ${data.node_name} conectado.` : `[ERRO] Falha ao registrar nó.`)
}

const handleAutoDiscover = async () => {
  isAutoScanning.value = true
  const [streamProbe, forgeProbe] = await Promise.all([probeClusterNode('127.0.0.1', 8080), probeClusterNode('127.0.0.1', 8081)])
  isAutoScanning.value = false
  let added = 0
  const streamIdx = servers.value.filter(s => s.role === 'HYDRASTREAM_EDGE').length
  if (streamProbe.online && !servers.value.some(s => s.ip === '127.0.0.1' && s.role === 'HYDRASTREAM_EDGE')) {
    await addNode({ node_name: `HYDRA-STREAM-NODE-${streamIdx}`, node_role: 'edge_ingest', ip_address: '127.0.0.1', http_port: 8080, grpc_port: 50051, webrtc_port: 8889, gpu_device_info: streamProbe.gpu_model || 'NVIDIA GeForce RTX 5090' })
    added++
  }
  const forgeIdx = servers.value.filter(s => s.role === 'HYDRASTREAM_GPU_WORKER').length
  if (forgeProbe.online && !servers.value.some(s => s.ip === '127.0.0.1' && s.role === 'HYDRASTREAM_GPU_WORKER')) {
    await addNode({ node_name: `HYDRA-FORGE-NODE-${forgeIdx}`, node_role: 'gpu_worker', ip_address: '127.0.0.1', http_port: 8081, grpc_port: 50051, webrtc_port: 8889, gpu_device_info: forgeProbe.gpu_model || 'NVIDIA GeForce RTX 5090' })
    added++
  }
  showToast(added > 0 ? `[AUTO-DISCOVERY] ${added} nó(s) detectados e adicionados com sucesso!` : `[AUTO-DISCOVERY] Nenhum novo nó detectado em localhost.`)
}
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1rem;">
    <div v-if="toastMessage" class="vms-badge vms-badge-orange" style="padding: 8px 14px; font-weight: bold; border-radius: 4px;">{{ toastMessage }}</div>
    <div class="vms-flex-between" style="align-items: center; gap: 1rem; flex-wrap: wrap;">
      <div class="vms-flex-col" style="gap: 2px;">
        <h2 class="vms-h2" style="color: #ffffff; margin: 0; font-family: var(--vms-font-jetbrains);">PERFORMANCE & NÓS DO CLUSTER</h2>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">GESTÃO DE INSTÂNCIAS HYDRASTREAM (INGESTÃO) E HYDRAFORGE (IA GPU)</span>
      </div>
      <div class="vms-flex-row" style="gap: 8px; align-items: center;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" :disabled="isAutoScanning" style="font-size: 11px; padding: 4px 10px;" @click="handleAutoDiscover">
          {{ isAutoScanning ? 'BUSCANDO...' : '🔍 AUTO-DETECTAR INSTÂNCIAS' }}
        </button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-size: 11px; padding: 4px 10px;" @click="showCreateModal = true">
          + CONECTAR NOVO NÓ
        </button>
        <div class="vms-badge vms-badge-orange" style="font-family: var(--vms-font-jetbrains); font-weight: 700; padding: 4px 10px; font-size: 11px;">VERIFICAÇÃO: {{ countdown }}s</div>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 5px 8px; border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 4px; background: rgba(255, 94, 58, 0.08);" title="Configurar Alertas" @click="showConfigModal = true">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
        </button>
      </div>
    </div>
    <div v-if="servers.length === 0" class="vms-flex-col" style="align-items: center; justify-content: center; padding: 3rem 1rem; border: 1px dashed var(--vms-border); border-radius: var(--vms-radius-md); background: #0c0e14; gap: 0.75rem;">
      <span class="vms-text-mono vms-text-xs vms-text-dim">// NENHUM NÓ HYDRASTREAM OU HYDRAFORGE CADASTRADO NO CLUSTER.</span>
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" :disabled="isAutoScanning" @click="handleAutoDiscover">🔍 AUTO-DETECTAR INSTÂNCIAS LOCAIS</button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" @click="showCreateModal = true">+ INSERIR IP MANUALMENTE</button>
      </div>
    </div>
    <div v-else class="vms-flex-col" style="gap: 8px; margin-top: 4px;">
      <ServerPerformanceRow v-for="server in servers" :key="server.id" :node="server" />
    </div>
    <PerformanceAlertsModal v-if="showConfigModal" @close="showConfigModal = false" @save="handleSaveConfig" />
    <CreateClusterNodeModal v-if="showCreateModal" @close="showCreateModal = false" @save="handleSaveNode" />
  </div>
</template>
