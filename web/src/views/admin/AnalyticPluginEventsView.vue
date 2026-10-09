<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMarketplace } from '../../composables/useMarketplace'
import AnalyticFilterBar from '../../components/admin/marketplace/AnalyticFilterBar.vue'
import AnalyticFacialFilterBar from '../../components/admin/marketplace/AnalyticFacialFilterBar.vue'
import AnalyticEventsGrid from '../../components/admin/marketplace/AnalyticEventsGrid.vue'
import AnalyticFacialEventsGrid from '../../components/admin/marketplace/AnalyticFacialEventsGrid.vue'
import AnalyticEventDetailModal from '../../components/admin/marketplace/AnalyticEventDetailModal.vue'
import type { AnalyticEventRecord } from '../../types/marketplace'

const props = defineProps<{ pluginId: string }>()
const { events, toast, showToast, getPluginById } = useMarketplace()

const plugin = computed(() => getPluginById(props.pluginId))
const isFacial = computed(() => props.pluginId === 'facial_recognition_sota' || plugin.value?.name?.toLowerCase().includes('facial'))

const expCameraId = ref('ALL'), expObject = ref('ALL')
const expAgeGroup = ref('ALL'), expSearchFace = ref(''), expFaceType = ref<'ALL' | 'KNOWN' | 'UNKNOWN'>('ALL')
const expStartDate = ref(''), expStartTime = ref(''), expEndDate = ref(''), expEndTime = ref('')
const activeEventModal = ref<AnalyticEventRecord | null>(null)

const pluginEvents = computed(() => events.value.filter(e => e.plugin_id === props.pluginId))

const filteredEvents = computed(() => {
  const list = pluginEvents.value.filter(e => {
    const matchCamera = expCameraId.value === 'ALL' || e.camera_id === expCameraId.value
    let matchFacial = true
    if (isFacial.value) {
      if (expSearchFace.value.trim()) {
        const q = expSearchFace.value.toLowerCase().trim()
        const pName = (e.person_name || '').toLowerCase()
        if (!pName.includes(q)) matchFacial = false
      }
      if (expFaceType.value === 'KNOWN' && e.face_status !== 'Reconhecido') matchFacial = false
      if (expFaceType.value === 'UNKNOWN' && e.face_status === 'Reconhecido') matchFacial = false
      if (expAgeGroup.value !== 'ALL' && e.age !== undefined) {
        if (expAgeGroup.value === 'child' && (e.age < 0 || e.age > 12)) matchFacial = false
        if (expAgeGroup.value === 'youth' && (e.age < 13 || e.age > 24)) matchFacial = false
        if (expAgeGroup.value === 'adult' && (e.age < 25 || e.age > 59)) matchFacial = false
        if (expAgeGroup.value === 'senior' && e.age < 60) matchFacial = false
      }
    } else {
      if (expObject.value !== 'ALL' && e.object_label.toLowerCase() !== expObject.value.toLowerCase()) return false
    }

    if (!matchCamera || !matchFacial) return false

    if (e.timestamp) {
      const evtTime = new Date(e.timestamp).getTime()
      if (expStartDate.value) {
        const sTime = new Date(`${expStartDate.value}T${expStartTime.value || '00:00'}:00`).getTime()
        if (!isNaN(sTime) && evtTime < sTime) return false
      }
      if (expEndDate.value) {
        const eTime = new Date(`${expEndDate.value}T${expEndTime.value || '23:59'}:59`).getTime()
        if (!isNaN(eTime) && evtTime > eTime) return false
      }
    }
    return true
  })
  return list.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
})

const handleExport = (format: 'csv' | 'json') => {
  const name = plugin.value ? plugin.value.name.toLowerCase().replace(/[^a-z0-9]/g, '_') : 'eventos'
  if (format === 'json') {
    const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(filteredEvents.value, null, 2))
    const a = document.createElement('a'); a.href = dataStr; a.download = `${name}_${Date.now()}.json`; a.click()
    showToast('[EXPORT] Arquivo JSON gerado com sucesso')
  } else {
    const headers = isFacial.value ? ['id', 'timestamp', 'camera_id', 'camera_name', 'status_reconhecimento', 'person_name', 'age', 'gender', 'match_rate', 'details'] : ['id', 'timestamp', 'camera_id', 'camera_name', 'event_type', 'object', 'details']
    const rows = filteredEvents.value.map(e => isFacial.value ? [e.id, e.timestamp, e.camera_id, `"${e.camera_name}"`, `"${e.face_status || 'Desconhecido'}"`, `"${e.person_name || 'Desconhecido'}"`, e.age || 'N/D', e.gender || 'N/D', e.confidence, `"${e.details}"`] : [e.id, e.timestamp, e.camera_id, `"${e.camera_name}"`, e.event_type, e.object_label, `"${e.details}"`])
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

    <!-- Header -->
    <div v-if="plugin" class="vms-flex-between" style="align-items: center; flex-wrap: wrap; gap: 1rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <div class="vms-flex-row" style="align-items: center; gap: 0.65rem;">
          <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange); margin: 0;">{{ plugin.name.toUpperCase() }} // EVENTOS</h3>
          <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ filteredEvents.length }} EVENTOS ENCONTRADOS</span>
        </div>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          {{ isFacial ? 'EXPLORAÇÃO DE BIOMETRIA FACIAL, IDENTIFICAÇÃO E FAIXA ETÁRIA' : 'EXPLORAÇÃO DE CARDS, SNAPSHOTS E METADADOS VIA FILTROS DINÂMICOS' }}
        </span>
      </div>
    </div>

    <!-- Barra de Filtros (Específica para Facial ou Objetos) -->
    <AnalyticFacialFilterBar
      v-if="isFacial"
      :events="pluginEvents"
      :selected-camera-id="expCameraId"
      :selected-age-group="expAgeGroup"
      :search-known-face="expSearchFace"
      :selected-face-type="expFaceType"
      :start-date="expStartDate"
      :start-time="expStartTime"
      :end-date="expEndDate"
      :end-time="expEndTime"
      @update:selected-camera-id="expCameraId = $event"
      @update:selected-age-group="expAgeGroup = $event"
      @update:search-known-face="expSearchFace = $event"
      @update:selected-face-type="expFaceType = $event"
      @update:start-date="expStartDate = $event"
      @update:start-time="expStartTime = $event"
      @update:end-date="expEndDate = $event"
      @update:end-time="expEndTime = $event"
      @reset="expCameraId = 'ALL'; expAgeGroup = 'ALL'; expSearchFace = ''; expFaceType = 'ALL'; expStartDate = ''; expStartTime = ''; expEndDate = ''; expEndTime = ''"
      @export="handleExport"
    />
    <AnalyticFilterBar
      v-else
      :events="pluginEvents"
      :selected-camera-id="expCameraId"
      :selected-object="expObject"
      :start-date="expStartDate"
      :start-time="expStartTime"
      :end-date="expEndDate"
      :end-time="expEndTime"
      @update:selected-camera-id="expCameraId = $event"
      @update:selected-object="expObject = $event"
      @update:start-date="expStartDate = $event"
      @update:start-time="expStartTime = $event"
      @update:end-date="expEndDate = $event"
      @update:end-time="expEndTime = $event"
      @reset="expCameraId = 'ALL'; expObject = 'ALL'; expStartDate = ''; expStartTime = ''; expEndDate = ''; expEndTime = ''"
      @export="handleExport"
    />

    <!-- Grid de Cards Puros (Facial vs Objetos) -->
    <AnalyticFacialEventsGrid v-if="isFacial" :events="filteredEvents" @inspect="activeEventModal = $event" />
    <AnalyticEventsGrid v-else :events="filteredEvents" @inspect="activeEventModal = $event" />

    <!-- Modal com Payload JSON Formatado -->
    <AnalyticEventDetailModal :is-open="!!activeEventModal" :event="activeEventModal" @close="activeEventModal = null" />
  </div>
</template>
