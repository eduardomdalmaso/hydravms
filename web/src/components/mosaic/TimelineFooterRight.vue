<script setup lang="ts">
defineProps<{
  isLive: boolean
  isExportMode?: boolean
  isSyncActive?: boolean
}>()

const emit = defineEmits<{
  (e: 'showLayers'): void
  (e: 'toggleExport'): void
  (e: 'toggleSync'): void
  (e: 'goLive'): void
}>()
</script>

<template>
  <div class="vms-vezha-controls-right" style="display: flex; align-items: center; gap: 0.35rem;">
    <!-- Camadas Button -->
    <button
      class="vms-btn vms-btn-sm vms-btn-secondary"
      style="padding: 0.25rem 0.55rem; font-size: 10px;"
      title="Filtrar Camadas de Gravação (Contínuo, Movimento, IA)"
      @click="emit('showLayers')"
    >
      [CAMADAS]
    </button>

    <!-- Scissors Icon Only Button (Toggle Export Period Selection) -->
    <button
      class="vms-btn vms-btn-sm"
      :class="isExportMode ? 'vms-btn-primary' : 'vms-btn-secondary'"
      style="display: flex; align-items: center; justify-content: center; padding: 0.25rem 0.55rem;"
      :style="isExportMode ? 'box-shadow: 0 0 10px rgba(0, 240, 255, 0.45); border-color: #00f0ff; color: #00f0ff; background: rgba(0, 240, 255, 0.2);' : ''"
      title="Recortar e Exportar Trecho de Vídeo (Selecionar Período)"
      @click="emit('toggleExport')"
    >
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="6" cy="6" r="3"></circle>
        <circle cx="6" cy="18" r="3"></circle>
        <line x1="20" y1="4" x2="8.12" y2="15.88"></line>
        <line x1="14.47" y1="14.48" x2="20" y2="20"></line>
        <line x1="8.12" y1="8.12" x2="12" y2="12"></line>
      </svg>
    </button>

    <!-- Sync Button (Visual Toggle without alert popup) -->
    <button
      class="vms-btn vms-btn-sm"
      :class="isSyncActive ? 'vms-btn-primary' : 'vms-btn-secondary'"
      style="padding: 0.25rem 0.55rem; font-size: 10px;"
      :style="isSyncActive ? 'box-shadow: 0 0 10px rgba(255, 94, 58, 0.45); border-color: var(--vms-neu-accent-orange);' : ''"
      title="Sincronizar Todas as Câmeras no Mesmo Horário"
      @click="emit('toggleSync')"
    >
      [SYNC]
    </button>

    <!-- Live Button -->
    <button
      class="vms-btn vms-btn-sm"
      :class="isLive ? 'vms-btn-danger' : 'vms-btn-secondary'"
      style="padding: 0.25rem 0.55rem; font-size: 10px;"
      title="Ir para o fluxo Ao Vivo"
      @click="emit('goLive')"
    >
      [LIVE]
    </button>
  </div>
</template>
