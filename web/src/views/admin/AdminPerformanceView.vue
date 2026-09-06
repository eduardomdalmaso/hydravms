<script setup lang="ts">
import { ref } from 'vue'
import { usePerformanceCluster } from '../../composables/usePerformanceCluster'
import ServerPerformanceRow from '../../components/admin/performance/ServerPerformanceRow.vue'
import PerformanceAlertsModal from '../../components/admin/performance/PerformanceAlertsModal.vue'

const { servers, countdown } = usePerformanceCluster()
const showConfigModal = ref(false)
const toastMessage = ref<string | null>(null)

const handleSaveConfig = (cfg: { cpuLimit: number; vramMinPercent: number; ramLimit: number }) => {
  showConfigModal.value = false
  toastMessage.value = `[ALERTAS ATIVADOS] CPU > ${cfg.cpuLimit}% // VRAM livre < ${cfg.vramMinPercent}% // RAM > ${cfg.ramLimit}%`
  setTimeout(() => { toastMessage.value = null }, 4000)
}
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1rem;">
    <!-- Toast Feedback Notification -->
    <div v-if="toastMessage" class="vms-badge vms-badge-orange" style="padding: 8px 14px; font-weight: bold; border-radius: 4px;">
      {{ toastMessage }}
    </div>

    <!-- Header com Titulo, Contador Regressivo de 10s e Engrenagem -->
    <div class="vms-flex-between" style="align-items: center; gap: 1rem; flex-wrap: wrap;">
      <div class="vms-flex-col" style="gap: 2px;">
        <h2 class="vms-h2" style="color: #ffffff; margin: 0; font-family: var(--vms-font-jetbrains);">
          PERFORMANCE DO CLUSTER
        </h2>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          MONITORAMENTO DE CPU, MEMORIA RAM E VRAM DAS GPUS POR SERVIDOR
        </span>
      </div>

      <div class="vms-flex-row" style="gap: 8px; align-items: center;">
        <div class="vms-badge vms-badge-orange" style="font-family: var(--vms-font-jetbrains); font-weight: 700; padding: 4px 10px; font-size: 11px;">
          VERIFICAÇÃO: {{ countdown }}s
        </div>
        <button
          class="vms-btn vms-btn-ghost vms-btn-sm"
          style="padding: 5px 8px; border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 4px; background: rgba(255, 94, 58, 0.08);"
          title="Configurar Alertas de Performance"
          @click="showConfigModal = true"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Lista em Formato de Linhas -->
    <div class="vms-flex-col" style="gap: 8px; margin-top: 4px;">
      <ServerPerformanceRow
        v-for="server in servers"
        :key="server.id"
        :node="server"
      />
    </div>

    <!-- Modal de Configuracao de Alertas -->
    <PerformanceAlertsModal
      v-if="showConfigModal"
      @close="showConfigModal = false"
      @save="handleSaveConfig"
    />
  </div>
</template>
