<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMarketplace } from '../../composables/useMarketplace'
import AnalyticFilterBar from '../../components/admin/marketplace/AnalyticFilterBar.vue'
import AnalyticEventsTable from '../../components/admin/marketplace/AnalyticEventsTable.vue'
import AnalyticEventsGrid from '../../components/admin/marketplace/AnalyticEventsGrid.vue'
import AnalyticEventDetailModal from '../../components/admin/marketplace/AnalyticEventDetailModal.vue'
import type { AnalyticEventRecord } from '../../types/marketplace'

const props = defineProps<{ pluginId: string }>()
const { plugins, events, toast, showToast, getPluginById } = useMarketplace()

const plugin = computed(() => getPluginById(props.pluginId))
const expCameraId = ref('ALL'), expSeverity = ref('ALL'), expSearch = ref(''), expMinConfidence = ref(0.5), expViewMode = ref<'table' | 'grid'>('table')
const activeEventModal = ref<AnalyticEventRecord | null>(null)

const filteredEvents = computed(() => events.value.filter(e => {
  if (e.plugin_id !== props.pluginId) return false
  const mC = expCameraId.value === 'ALL' || e.camera_id === expCameraId.value
  const mS = expSeverity.value === 'ALL' || e.severity === expSeverity.value
  const mConf = e.confidence >= expMinConfidence.value
  const mQ = !expSearch.value || e.object_label.toLowerCase().includes(expSearch.value.toLowerCase()) || e.details.toLowerCase().includes(expSearch.value.toLowerCase())
  return mC && mS && mConf && mQ
}))

const handleExport = (format: 'csv' | 'json') => {
  const name = plugin.value ? plugin.value.name.toLowerCase().replace(/[^a-z0-9]/g, '_') : 'eventos'
  if (format === 'json') {
    const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(filteredEvents.value, null, 2))
    const a = document.createElement('a'); a.href = dataStr; a.download = `${name}_${Date.now()}.json`; a.click()
    showToast('[EXPORT] Arquivo JSON gerado com sucesso')
  } else {
    const headers = ['id', 'timestamp', 'camera', 'object', 'confidence', 'severity', 'details']
    const rows = filteredEvents.value.map(e => [e.id, e.timestamp, e.camera_name, e.object_label, e.confidence, e.severity, `"${e.details}"`])
    const csv = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    const a = document.createElement('a'); a.href = encodeURI(csv); a.download = `${name}_${Date.now()}.csv`; a.click()
    showToast('[EXPORT] Arquivo CSV gerado com sucesso')
  }
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <Transition name="vms-toast">
      <div v-if="toast" class="vms-toast-notification" title="Clique para fechar" @click="toast = null"><span>{{ toast }}</span></div>
    </Transition>

    <div v-if="plugin" class="vms-flex-between" style="align-items: center; flex-wrap: wrap; gap: 1rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <div class="vms-flex-row" style="align-items: center; gap: 0.65rem;">
          <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange); margin: 0;">{{ plugin.name.toUpperCase() }} // EVENTOS</h3>
          <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ filteredEvents.length }} EVENTOS ENCONTRADOS</span>
        </div>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">EXPLORAÇÃO DE DADOS, TELEMETRIA E CROPS FORENSES VIA FILTROS</span>
      </div>
    </div>

    <AnalyticFilterBar
      :plugins="plugins.filter(p => p.id === props.pluginId)"
      :selected-plugin-id="props.pluginId"
      :selected-camera-id="expCameraId"
      :selected-severity="expSeverity"
      :search-query="expSearch"
      :min-confidence="expMinConfidence"
      :view-mode="expViewMode"
      @update:selected-plugin-id="() => {}"
      @update:selected-camera-id="expCameraId = $event"
      @update:selected-severity="expSeverity = $event"
      @update:search-query="expSearch = $event"
      @update:min-confidence="expMinConfidence = $event"
      @update:view-mode="expViewMode = $event"
      @export="handleExport"
    />

    <AnalyticEventsTable v-if="expViewMode === 'table'" :events="filteredEvents" @inspect="activeEventModal = $event" />
    <AnalyticEventsGrid v-else :events="filteredEvents" @inspect="activeEventModal = $event" />
    <AnalyticEventDetailModal :is-open="!!activeEventModal" :event="activeEventModal" @close="activeEventModal = null" />
  </div>
</template>
