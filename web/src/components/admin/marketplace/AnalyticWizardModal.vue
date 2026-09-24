<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { AnalyticInstance, AnalyticFolderNode, PluginManifest, AnalyticMode, Point2D, Line2D } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'
import { useMarketplace } from '../../../composables/useMarketplace'
import AnalyticWizardStep1 from './AnalyticWizardStep1.vue'
import AnalyticWizardStep2 from './AnalyticWizardStep2.vue'

const props = defineProps<{
  isOpen: boolean; plugin: PluginManifest; folders: AnalyticFolderNode[]; currentFolderId?: string | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', payload: { inst: AnalyticInstance; targetFolderId: string | null }): void
}>()

const { gpuDetected } = useMarketplace()
const step = ref<1 | 2>(1)
const name = ref(''), camera = ref(''), stream = ref<'main_1080p' | 'sub_stream'>('main_1080p')
const hardware = ref<'rtx_5090_cuda' | 'cpu_shm'>(gpuDetected.value ? 'rtx_5090_cuda' : 'cpu_shm')
const mode = ref<AnalyticMode>('intrusion')
const fps = ref(15), motionGated = ref(true), autoSahi = ref(true)
const scheduleStart = ref('18:00'), scheduleEnd = ref('06:00'), selectedDays = ref([1,2,3,4,5])
const selectedFolder = ref<string | null>(props.currentFolderId || null)
const selectedClasses = ref<string[]>(['person', 'car', 'motorcycle', 'cell_phone'])
const cameras = ref<{ id: string; name: string }[]>([])

const points = ref<Point2D[]>([
  { x: 0.20, y: 0.15 }, { x: 0.80, y: 0.15 }, { x: 0.90, y: 0.50 },
  { x: 0.80, y: 0.85 }, { x: 0.20, y: 0.85 }, { x: 0.10, y: 0.50 }
])
const line = ref<Line2D>({ p1: { x: 0.20, y: 0.50 }, p2: { x: 0.80, y: 0.50 } })

onMounted(async () => {
  const cams = await fetchCameras()
  cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
  if (cameras.value.length > 0) camera.value = cameras.value[0].id
})

const snapshotUrl = computed(() => camera.value ? `http://localhost:8083/api/v1/cameras/${camera.value}/snapshot` : '')

const toggleClass = (id: string) => {
  const idx = selectedClasses.value.indexOf(id)
  if (idx >= 0) selectedClasses.value.splice(idx, 1)
  else selectedClasses.value.push(id)
}

const toggleDay = (d: number) => {
  const idx = selectedDays.value.indexOf(d)
  if (idx >= 0) selectedDays.value.splice(idx, 1)
  else selectedDays.value.push(d)
}

const handleSave = () => {
  const camObj = cameras.value.find(c => c.id === camera.value)
  const newInst: AnalyticInstance = {
    id: `inst-${Date.now()}`, plugin_id: props.plugin.id, plugin_name: props.plugin.name,
    name: name.value.trim() || `${props.plugin.name} // ${camObj?.name || 'Câmera'}`,
    camera_id: camera.value || 'cam_stream', camera_name: camObj?.name || 'Câmera',
    stream_type: stream.value, hardware_target: hardware.value,
    confidence_threshold: 0.70, roi_mode: mode.value === 'counting' ? 'counting_line' : 'custom_polygon',
    fps_rate: fps.value, motion_gated: motionGated.value, auto_sahi: autoSahi.value,
    specific_params: { mode: mode.value, classes: selectedClasses.value },
    zones: [{
      id: `zone-${Date.now()}`, name: `ZONA 1 // ${mode.value.toUpperCase()}`,
      mode: mode.value, target_classes: selectedClasses.value, polygon: points.value, line: line.value,
      schedule: { days: selectedDays.value, start_time: scheduleStart.value, end_time: scheduleEnd.value }
    }],
    is_active: true, detections_count: 0, created_at: new Date().toISOString()
  }
  emit('save', { inst: newInst, targetFolderId: selectedFolder.value })
  step.value = 1; name.value = ''
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content vms-fixed-wizard">
      <!-- Header -->
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">NOVO ANALÍTICO // {{ plugin.name }}</h4>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">
            PASSO {{ step }} DE 2 // {{ step === 1 ? 'GEOMETRIA & OBJETOS' : 'PERFORMANCE & AGENDAMENTO' }}
          </span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <!-- Content -->
      <div class="vms-wizard-body">
        <AnalyticWizardStep1
          v-if="step === 1"
          :name="name" :camera="camera" :cameras="cameras" :mode="mode"
          :selected-classes="selectedClasses" :points="points" :line="line" :snapshot-url="snapshotUrl"
          @update:name="name = $event" @update:camera="camera = $event" @update:mode="mode = $event"
          @toggle:class="toggleClass" @update:points="points = $event" @update:line="line = $event"
        />
        <AnalyticWizardStep2
          v-else
          :fps="fps" :motion-gated="motionGated" :auto-sahi="autoSahi"
          :schedule-start="scheduleStart" :schedule-end="scheduleEnd" :selected-days="selectedDays"
          :hardware="hardware" :gpu-detected="gpuDetected"
          @update:fps="fps = $event" @update:motion-gated="motionGated = $event"
          @update:auto-sahi="autoSahi = $event" @update:schedule-start="scheduleStart = $event"
          @update:schedule-end="scheduleEnd = $event" @toggle:day="toggleDay" @update:hardware="hardware = $event"
        />
      </div>

      <!-- Footer Navigation -->
      <div class="vms-modal-footer vms-flex-between">
        <button v-if="step === 1" class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">CANCELAR</button>
        <button v-else class="vms-btn vms-btn-secondary vms-btn-sm" @click="step = 1">VOLTAR</button>
        <button v-if="step === 1" class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="step = 2">AVANÇAR</button>
        <button v-else class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-fixed-wizard { width: 780px; max-width: 95vw; height: 580px; display: flex; flex-direction: column; }
.vms-wizard-body { flex: 1; overflow-y: auto; padding: 0.85rem 0; }
</style>
