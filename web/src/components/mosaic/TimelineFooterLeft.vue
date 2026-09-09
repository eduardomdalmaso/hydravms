<script setup lang="ts">
defineProps<{
  isExportMode?: boolean
  exportDurationText?: string
}>()
const emit = defineEmits<{ (e: 'openModal'): void }>()
</script>

<template>
  <div class="vms-vezha-controls-left">
    <button
      v-if="isExportMode"
      class="vms-btn vms-btn-sm vms-btn-primary"
      style="padding: 0.25rem 0.65rem; font-size: 10px; font-weight: 700; display: flex; align-items: center; gap: 0.35rem; box-shadow: 0 0 10px rgba(255, 94, 58, 0.4);"
      title="Ajustar Horário Exato e Baixar Trecho"
      @click="emit('openModal')"
    >
      <span>BAIXAR [{{ exportDurationText || '10 MIN' }}]</span>
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
    </button>
  </div>
</template>
