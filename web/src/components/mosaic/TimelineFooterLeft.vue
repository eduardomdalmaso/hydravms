<script setup lang="ts">
defineProps<{
  isExportMode?: boolean
  exportDurationText?: string
}>()
</script>

<template>
  <div class="vms-vezha-controls-left">
    <span
      v-if="isExportMode"
      class="vms-badge vms-badge-orange"
      style="font-size: 10px; font-weight: 700;"
    >
      [RECORTE: {{ exportDurationText || '10 MIN' }}]
    </span>
  </div>
</template>
