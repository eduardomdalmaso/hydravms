<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { AnalyticInstance, AnalyticFolderNode, PluginManifest, ZoneConfig } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useMarketplace } from '../../../composables/useMarketplace'
import AnalyticRoiCanvas from './AnalyticRoiCanvas.vue'
import AnalyticZoneAccordionItem from './AnalyticZoneAccordionItem.vue'
import AnalyticFacialParamsPanel from './AnalyticFacialParamsPanel.vue'
import AnalyticPerfHardwarePanel from './AnalyticPerfHardwarePanel.vue'

const props = defineProps<{ isOpen: boolean; plugin: PluginManifest; folders: AnalyticFolderNode[]; currentFolderId?: string | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', payload: { inst: AnalyticInstance; targetFolderId: string | null }): void }>()

const { gpuDetected } = useMarketplace()
const zoneColors = ['#ff5e3a', '#00f0ff', '#fcee0a', '#00ff9d', '#e024c3', '#a855f7']
const activeTab = ref<'rules' | 'perf'>('rules'), facialSubTab = ref<'params' | 'zones'>('params')
const name = ref(''), camera = ref(''), fps = ref(15), motionGated = ref(true)
const hardware = ref<'rtx_5090_cuda' | 'cpu_shm'>(gpuDetected.value ? 'rtx_5090_cuda' : 'cpu_shm')
const selectedFolder = ref<string | null>(props.currentFolderId || null), cameras = ref<{ id: string; name: string }[]>([]), activeZoneId = ref<string>('')

// Parâmetros faciais específicos (estilo Vezha)
const minHeight = ref(40), maxHeight = ref<number | string>(0), minFaces = ref(25)
const blurThreshold = ref(60), recognitionThreshold = ref(75)
const detectorSensitivity = ref<'baixo' | 'normal' | 'alto'>('normal'), detectGenderAge = ref(false)
const isFacial = computed(() => props.plugin?.id === 'facial_recognition_sota' || props.plugin?.name?.toLowerCase().includes('facial'))

// Polígono padrão cobrindo 100% do snapshot para todos os analíticos
const getFullSnapshotPolygon = () => [{ x: 0.01, y: 0.01 }, { x: 0.99, y: 0.01 }, { x: 0.99, y: 0.99 }, { x: 0.01, y: 0.99 }]
const zones = ref<ZoneConfig[]>([{
  id: `zone-${Date.now()}-1`, name: 'ZONA 1', mode: isFacial.value ? 'capture' : 'intrusion',
  target_classes: isFacial.value ? [] : ['person'], polygon: getFullSnapshotPolygon(),
  line: { p1: { x: 0.10, y: 0.50 }, p2: { x: 0.90, y: 0.50 } }, schedules: [{ id: 's1', days: [0, 1, 2, 3, 4, 5, 6], start_time: '00:00', end_time: '23:59' }]
}])
const isNameValid = computed(() => name.value.trim().length > 0)

onMounted(async () => {
  const cams = await fetchCameras(); cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
  if (cameras.value.length > 0) camera.value = cameras.value[0].id
  if (zones.value.length > 0) activeZoneId.value = zones.value[0].id
})

const snapshotUrl = computed(() => camera.value ? getCameraSnapshotUrl(camera.value, false) : '')
const addZone = () => {
  const count = zones.value.length + 1
  const newZ: ZoneConfig = {
    id: `zone-${Date.now()}-${count}`, name: `ZONA ${count}`, mode: isFacial.value ? 'capture' : 'intrusion',
    target_classes: isFacial.value ? [] : ['person'], polygon: getFullSnapshotPolygon(),
    line: { p1: { x: 0.10, y: 0.50 }, p2: { x: 0.90, y: 0.50 } }, schedules: [{ id: `s-${Date.now()}`, days: [0, 1, 2, 3, 4, 5, 6], start_time: '00:00', end_time: '23:59' }]
  }
  zones.value.push(newZ); activeZoneId.value = newZ.id
}
const deleteZone = (idx: number) => { if (zones.value.length > 1) { zones.value.splice(idx, 1); activeZoneId.value = zones.value[0].id } }
const addSchedule = (zId: string) => { const z = zones.value.find(item => item.id === zId); if (z) { if (!z.schedules) z.schedules = []; z.schedules.push({ id: `s-${Date.now()}`, days: [], start_time: '00:00', end_time: '23:59' }) } }

const handleSave = () => {
  if (!isNameValid.value || activeTab.value !== 'perf') return
  const camObj = cameras.value.find(c => c.id === camera.value)
  const specificParams = isFacial.value ? {
    min_height: Number(minHeight.value) || 40, max_height: Number(maxHeight.value) || 0,
    min_faces: Number(minFaces.value) || 25, blur_threshold: Number(blurThreshold.value) || 60,
    recognition_threshold: (Number(recognitionThreshold.value) || 75) / 100.0,
    detector_sensitivity: detectorSensitivity.value, detect_gender_age: detectGenderAge.value,
    zones_count: zones.value.length
  } : { zones_count: zones.value.length }

  const newInst: AnalyticInstance = {
    id: `inst-${Date.now()}`, plugin_id: props.plugin.id, plugin_name: props.plugin.name,
    name: name.value.trim(), camera_id: camera.value || 'cam_stream', camera_name: camObj?.name || 'Câmera',
    stream_type: 'main_1080p', hardware_target: hardware.value, confidence_threshold: 0.70,
    roi_mode: 'custom_polygon', fps_rate: fps.value, motion_gated: motionGated.value, auto_sahi: true,
    specific_params: specificParams, zones: zones.value, is_active: true, detections_count: 0, created_at: new Date().toISOString()
  }
  emit('save', { inst: newInst, targetFolderId: selectedFolder.value })
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop">
    <div class="vms-modal-solid-panel">
      <!-- Header -->
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">NOVO ANALÍTICO // {{ plugin.name }}</h4>
          <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ isFacial ? 'BIOMETRIA FACIAL' : `${zones.length} ZONA(S)` }}</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
      </div>

      <div class="vms-wizard-split-body">
        <!-- Esquerda: Snapshot & Canvas ROI (100% full snapshot por padrão) -->
        <div class="vms-wizard-left-pane">
          <AnalyticRoiCanvas
            :snapshot-url="snapshotUrl" :zones="zones" :active-zone-id="activeZoneId" :zone-colors="zoneColors"
            @update:points="zones.find(z => z.id === $event.zoneId)!.polygon = $event.points"
            @update:line="zones.find(z => z.id === $event.zoneId)!.line = $event.line"
            @select:zone="activeZoneId = $event"
          />
        </div>

        <!-- Direita: Abas [Regras/Biometria] vs [Performance & Hardware] -->
        <div class="vms-wizard-right-pane">
          <div class="vms-flex-row" style="gap: 0.35rem; border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 0.4rem;">
            <button class="vms-btn vms-btn-sm" :class="activeTab === 'rules' ? 'vms-btn-primary' : 'vms-btn-ghost'" style="font-size: 10px; padding: 3px 8px;" @click="activeTab = 'rules'">
              {{ isFacial ? '1. BIOMETRIA & ZONAS' : '1. REGRAS & ZONAS' }}
            </button>
            <button class="vms-btn vms-btn-sm" :class="activeTab === 'perf' ? 'vms-btn-primary' : 'vms-btn-ghost'" :disabled="!isNameValid" :style="{ opacity: !isNameValid ? 0.4 : 1, cursor: !isNameValid ? 'not-allowed' : 'pointer' }" style="font-size: 10px; padding: 3px 8px;" @click="isNameValid && (activeTab = 'perf')">
              2. PERFORMANCE & HARDWARE
            </button>
          </div>

          <!-- Aba 1: Regras / Biometria -->
          <div v-if="activeTab === 'rules'" class="vms-flex-col" style="gap: 0.5rem; flex: 1;">
            <div class="vms-flex-row" style="gap: 0.5rem;">
              <input :value="name" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px; flex: 1.3;" placeholder="Nome do Analítico..." @input="name = ($event.target as HTMLInputElement).value" />
              <select :value="camera" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px; flex: 1;" @change="camera = ($event.target as HTMLSelectElement).value">
                <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>

            <!-- Seletor de Sub-Abas no Facial: Parâmetros vs Zonas de Captura -->
            <div v-if="isFacial" class="vms-flex-row" style="gap: 0.35rem;">
              <button class="vms-btn vms-btn-sm" :class="facialSubTab === 'params' ? 'vms-btn-secondary' : 'vms-btn-ghost'" style="flex: 1; padding: 3px; font-size: 9.5px; border: 1px solid var(--vms-border);" @click="facialSubTab = 'params'">
                [1. PARÂMETROS BIOMÉTRICOS]
              </button>
              <button class="vms-btn vms-btn-sm" :class="facialSubTab === 'zones' ? 'vms-btn-secondary' : 'vms-btn-ghost'" style="flex: 1; padding: 3px; font-size: 9.5px; border: 1px solid var(--vms-border);" @click="facialSubTab = 'zones'">
                [2. ZONAS DE CAPTURA ({{ zones.length }})]
              </button>
            </div>

            <!-- Painel Biométrico Facial -->
            <AnalyticFacialParamsPanel
              v-if="isFacial && facialSubTab === 'params'"
              v-model:min-height="minHeight" v-model:max-height="maxHeight"
              v-model:min-faces="minFaces" v-model:blur-threshold="blurThreshold"
              v-model:recognition-threshold="recognitionThreshold"
              v-model:detector-sensitivity="detectorSensitivity"
              v-model:detect-gender-age="detectGenderAge"
            />

            <!-- Zonas de Vídeo (Para analítico geral ou sub-aba de zonas do facial) -->
            <template v-if="!isFacial || facialSubTab === 'zones'">
              <div class="vms-flex-between" style="align-items: center; margin-top: 0.15rem;">
                <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">ZONAS DE VÍDEO (ACCORDION)</span>
                <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 2px 7px;" @click="addZone">+ ADICIONAR ZONA</button>
              </div>
              <div class="vms-flex-col" style="gap: 0.45rem;">
                <AnalyticZoneAccordionItem
                  v-for="(z, zIdx) in zones" :key="z.id" :zone="z" :color="zoneColors[zIdx % zoneColors.length]"
                  :is-open="activeZoneId === z.id" :can-delete="zones.length > 1" :is-facial="isFacial"
                  @toggle="activeZoneId = activeZoneId === z.id ? '' : z.id" @delete="deleteZone(zIdx)"
                  @update:name="z.name = $event" @update:mode="z.mode = $event"
                  @toggle:class="z.target_classes.includes($event) ? z.target_classes.splice(z.target_classes.indexOf($event), 1) : z.target_classes.push($event)"
                  @add:schedule="addSchedule(z.id)" @delete:schedule="z.schedules?.splice($event, 1)"
                  @update:schedule-time="z.schedules![$event.idx][$event.field] = $event.val"
                  @toggle:schedule-day="z.schedules![$event.idx].days.includes($event.day) ? z.schedules![$event.idx].days.splice(z.schedules![$event.idx].days.indexOf($event.day), 1) : z.schedules![$event.idx].days.push($event.day)"
                />
              </div>
            </template>
          </div>

          <!-- Aba 2: Performance, Hardware & Motion Gated -->
          <AnalyticPerfHardwarePanel
            v-else
            v-model:fps="fps" v-model:hardware="hardware" v-model:motion-gated="motionGated" :gpu-detected="gpuDetected"
          />
        </div>
      </div>

      <!-- Footer com Fluxo Obrigatório -->
      <div class="vms-modal-footer vms-flex-between">
        <button v-if="activeTab === 'rules'" class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">CANCELAR</button>
        <button v-else class="vms-btn vms-btn-secondary vms-btn-sm" @click="activeTab = 'rules'">← VOLTAR</button>

        <button v-if="activeTab === 'rules'" :disabled="!isNameValid" class="vms-btn vms-btn-primary vms-btn-sm" :style="{ opacity: !isNameValid ? 0.4 : 1, cursor: !isNameValid ? 'not-allowed' : 'pointer' }" style="font-weight: bold;" @click="isNameValid && (activeTab = 'perf')">AVANÇAR →</button>
        <button v-else :disabled="!isNameValid" class="vms-btn vms-btn-primary vms-btn-sm" :style="{ opacity: !isNameValid ? 0.4 : 1, cursor: !isNameValid ? 'not-allowed' : 'pointer' }" style="font-weight: bold;" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-modal-solid-panel { width: 980px; max-width: 96vw; height: 610px; max-height: 90vh; background: #0c0f17; border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 8px; box-shadow: 0 20px 50px rgba(0,0,0,0.95), 0 0 20px rgba(255, 94, 58, 0.12); display: flex; flex-direction: column; overflow: hidden; }
.vms-wizard-split-body { flex: 1; display: flex; overflow: hidden; }
.vms-wizard-left-pane { flex: 1.05; padding: 0.85rem; border-right: 1px solid rgba(255,255,255,0.06); display: flex; flex-direction: column; }
.vms-wizard-right-pane { flex: 1; padding: 0.85rem; overflow-y: auto; display: flex; flex-direction: column; gap: 0.65rem; background: #080a10; }
</style>
