<script setup lang="ts">
import { ref } from 'vue'
import type { PluginManifest, AnalyticEventRecord } from '../../../types/marketplace'
import AnalyticFilterBar from './AnalyticFilterBar.vue'
import AnalyticEventsTable from './AnalyticEventsTable.vue'
import AnalyticEventsGrid from './AnalyticEventsGrid.vue'
import AnalyticEventDetailModal from './AnalyticEventDetailModal.vue'

const props = defineProps<{
  plugins: PluginManifest[]
  events: AnalyticEventRecord[]
  selectedPluginId: string
  selectedCameraId: string
  selectedSeverity: string
  searchQuery: string
  minConfidence: number
  viewMode: 'table' | 'grid'
}>()

const emit = defineEmits<{
  (e: 'update:selectedPluginId', val: string): void
  (e: 'update:selectedCameraId', val: string): void
  (e: 'update:selectedSeverity', val: string): void
  (e: 'update:searchQuery', val: string): void
  (e: 'update:minConfidence', val: number): void
  (e: 'update:viewMode', val: 'table' | 'grid'): void
  (e: 'notify', msg: string): void
}>()

const activeEventModal = ref<AnalyticEventRecord | null>(null)

const handleExport = (format: 'csv' | 'json') => {
  if (format === 'json') {
    const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(props.events, null, 2))
    const dlAnchor = document.createElement('a')
    dlAnchor.setAttribute('href', dataStr)
    dlAnchor.setAttribute('download', `hydra_analytics_${Date.now()}.json`)
    dlAnchor.click()
    emit('notify', '[EXPORT] Arquivo JSON gerado com sucesso')
  } else {
    const headers = ['id', 'timestamp', 'camera', 'plugin', 'object', 'confidence', 'severity', 'details']
    const rows = props.events.map(e => [e.id, e.timestamp, e.camera_name, e.plugin_name, e.object_label, e.confidence, e.severity, `"${e.details}"`])
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    const dlAnchor = document.createElement('a')
    dlAnchor.setAttribute('href', encodeURI(csvContent))
    dlAnchor.setAttribute('download', `hydra_analytics_${Date.now()}.csv`)
    dlAnchor.click()
    emit('notify', '[EXPORT] Arquivo CSV gerado com sucesso')
  }
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <AnalyticFilterBar
      :plugins="plugins"
      :selected-plugin-id="selectedPluginId"
      :selected-camera-id="selectedCameraId"
      :selected-severity="selectedSeverity"
      :search-query="searchQuery"
      :min-confidence="minConfidence"
      :view-mode="viewMode"
      @update:selected-plugin-id="emit('update:selectedPluginId', $event)"
      @update:selected-camera-id="emit('update:selectedCameraId', $event)"
      @update:selected-severity="emit('update:selectedSeverity', $event)"
      @update:search-query="emit('update:searchQuery', $event)"
      @update:min-confidence="emit('update:minConfidence', $event)"
      @update:view-mode="emit('update:viewMode', $event)"
      @export="handleExport"
    />

    <AnalyticEventsTable
      v-if="viewMode === 'table'"
      :events="events"
      @inspect="activeEventModal = $event"
    />
    <AnalyticEventsGrid
      v-else
      :events="events"
      @inspect="activeEventModal = $event"
    />

    <AnalyticEventDetailModal
      :is-open="!!activeEventModal"
      :event="activeEventModal"
      @close="activeEventModal = null"
    />
  </div>
</template>
