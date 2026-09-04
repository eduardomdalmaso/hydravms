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
      class="vms-badge"
      style="background: rgba(0, 240, 255, 0.15); color: #00f0ff; border: 1px solid rgba(0, 240, 255, 0.35); font-size: 10px; font-weight: 700;"
    >
      [RECORTE: {{ exportDurationText || '10 MIN' }}]
    </span>
  </div>
</template>
