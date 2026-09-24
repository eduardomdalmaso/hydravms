<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { AnalyticInstance, AnalyticFolderNode, PluginManifest, ZoneConfig, Point2D, Line2D } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useMarketplace } from '../../../composables/useMarketplace'
import AnalyticRoiCanvas from './AnalyticRoiCanvas.vue'
import AnalyticZoneAccordionItem from './AnalyticZoneAccordionItem.vue'

const props = defineProps<{ isOpen: boolean; plugin: PluginManifest; folders: AnalyticFolderNode[]; currentFolderId?: string | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', payload: { inst: AnalyticInstance; targetFolderId: string | null }): void }>()

const { gpuDetected } = useMarketplace()
const zoneColors = ['#ff5e3a', '#00f0ff', '#fcee0a', '#00ff9d', '#e024c3', '#a855f7']
const activeTab = ref<'rules' | 'perf'>('rules')
const name = ref(''), camera = ref(''), fps = ref(15), motionGated = ref(true)
const hardware = ref<'rtx_5090_cuda' | 'cpu_shm'>(gpuDetected.value ? 'rtx_5090_cuda' : 'cpu_shm')
const selectedFolder = ref<string | null>(props.currentFolderId || null)
const cameras = ref<{ id: string; name: string }[]>([]), activeZoneId = ref<string>('')
const zones = ref<ZoneConfig[]>([{
  id: `zone-${Date.now()}-1`, name: 'ZONA 1', mode: 'intrusion', target_classes: ['person', 'car', 'motorcycle', 'cell_phone'],
  polygon: [{ x: 0.15, y: 0.15 }, { x: 0.85, y: 0.15 }, { x: 0.90, y: 0.50 }, { x: 0.85, y: 0.85 }, { x: 0.15, y: 0.85 }, { x: 0.10, y: 0.50 }],
  line: { p1: { x: 0.20, y: 0.50 }, p2: { x: 0.80, y: 0.50 } }, schedules: [{ id: 's1', days: [1,2,3,4,5], start_time: '00:00', end_time: '23:59' }]
}])

onMounted(async () => {
  const cams = await fetchCameras(); cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
  if (cameras.value.length > 0) camera.value = cameras.value[0].id
  if (zones.value.length > 0) activeZoneId.value = zones.value[0].id
})

const snapshotUrl = computed(() => camera.value ? getCameraSnapshotUrl(camera.value, false) : '')

const addZone = () => {
  const count = zones.value.length + 1
  const newZ: ZoneConfig = {
    id: `zone-${Date.now()}-${count}`, name: `ZONA ${count}`, mode: 'intrusion', target_classes: ['person', 'car'],
    polygon: [{ x: 0.25, y: 0.25 }, { x: 0.75, y: 0.25 }, { x: 0.80, y: 0.55 }, { x: 0.75, y: 0.75 }, { x: 0.25, y: 0.75 }, { x: 0.20, y: 0.55 }],
    line: { p1: { x: 0.25, y: 0.50 }, p2: { x: 0.75, y: 0.50 } }, schedules: [{ id: `s-${Date.now()}`, days: [1,2,3,4,5], start_time: '00:00', end_time: '23:59' }]
  }
  zones.value.push(newZ); activeZoneId.value = newZ.id
}
const deleteZone = (idx: number) => { if (zones.value.length > 1) { zones.value.splice(idx, 1); activeZoneId.value = zones.value[0].id } }
const addSchedule = (zId: string) => { const z = zones.value.find(item => item.id === zId); if (z) { if (!z.schedules) z.schedules = []; z.schedules.push({ id: `s-${Date.now()}`, days: [1,2,3,4,5], start_time: '18:00', end_time: '06:00' }) } }

const handleSave = () => {
  const camObj = cameras.value.find(c => c.id === camera.value)
  const newInst: AnalyticInstance = {
    id: `inst-${Date.now()}`, plugin_id: props.plugin.id, plugin_name: props.plugin.name,
    name: name.value.trim() || `${props.plugin.name} // ${camObj?.name || 'Câmera'}`,
    camera_id: camera.value || 'cam_stream', camera_name: camObj?.name || 'Câmera',
    stream_type: 'main_1080p', hardware_target: hardware.value, confidence_threshold: 0.70,
    roi_mode: 'custom_polygon', fps_rate: fps.value, motion_gated: motionGated.value, auto_sahi: true,
    specific_params: { zones_count: zones.value.length }, zones: zones.value, is_active: true, detections_count: 0, created_at: new Date().toISOString()
  }
  emit('save', { inst: newInst, targetFolderId: selectedFolder.value })
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-solid-panel">
      <!-- Header -->
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">NOVO ANALÍTICO // {{ plugin.name }}</h4>
          <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ zones.length }} {{ zones.length === 1 ? 'ZONA' : 'ZONAS' }}</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
      </div>

      <div class="vms-wizard-split-body">
        <!-- Esquerda: Snapshot & Canvas ROI -->
        <div class="vms-wizard-left-pane">
          <AnalyticRoiCanvas
            :snapshot-url="snapshotUrl" :zones="zones" :active-zone-id="activeZoneId" :zone-colors="zoneColors"
            @update:points="zones.find(z => z.id === $event.zoneId)!.polygon = $event.points"
            @update:line="zones.find(z => z.id === $event.zoneId)!.line = $event.line"
            @select:zone="activeZoneId = $event"
          />
        </div>

        <!-- Direita: Abas [Zonas & Regras] vs [Performance & Hardware] -->
        <div class="vms-wizard-right-pane">
          <div class="vms-flex-row" style="gap: 0.35rem; border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 0.4rem;">
            <button class="vms-btn vms-btn-sm" :class="activeTab === 'rules' ? 'vms-btn-primary' : 'vms-btn-ghost'" style="font-size: 10px; padding: 3px 8px;" @click="activeTab = 'rules'">1. REGRAS & ZONAS</button>
            <button class="vms-btn vms-btn-sm" :class="activeTab === 'perf' ? 'vms-btn-primary' : 'vms-btn-ghost'" style="font-size: 10px; padding: 3px 8px;" @click="activeTab = 'perf'">2. PERFORMANCE & HARDWARE</button>
          </div>

          <!-- Aba 1: Regras & Zonas -->
          <div v-if="activeTab === 'rules'" class="vms-flex-col" style="gap: 0.55rem; flex: 1;">
            <div class="vms-flex-row" style="gap: 0.5rem;">
              <input :value="name" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px; flex: 1.3;" placeholder="Nome do Analítico..." @input="name = ($event.target as HTMLInputElement).value" />
              <select :value="camera" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px; flex: 1;" @change="camera = ($event.target as HTMLSelectElement).value"><option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option></select>
            </div>
            <div class="vms-flex-between" style="align-items: center; margin-top: 0.2rem;">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">ZONAS DE VÍDEO (ACCORDION)</span>
              <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 2px 7px;" @click="addZone">+ ADICIONAR ZONA</button>
            </div>
            <div class="vms-flex-col" style="gap: 0.45rem;">
              <AnalyticZoneAccordionItem
                v-for="(z, zIdx) in zones" :key="z.id" :zone="z" :color="zoneColors[zIdx % zoneColors.length]"
                :is-open="activeZoneId === z.id" :can-delete="zones.length > 1"
                @toggle="activeZoneId = activeZoneId === z.id ? '' : z.id" @delete="deleteZone(zIdx)"
                @update:name="z.name = $event" @update:mode="z.mode = $event"
                @toggle:class="z.target_classes.includes($event) ? z.target_classes.splice(z.target_classes.indexOf($event), 1) : z.target_classes.push($event)"
                @add:schedule="addSchedule(z.id)" @delete:schedule="z.schedules?.splice($event, 1)"
                @update:schedule-time="z.schedules![$event.idx][$event.field] = $event.val"
                @toggle:schedule-day="z.schedules![$event.idx].days.includes($event.day) ? z.schedules![$event.idx].days.splice(z.schedules![$event.idx].days.indexOf($event.day), 1) : z.schedules![$event.idx].days.push($event.day)"
              />
            </div>
          </div>

          <!-- Aba 2: Performance, Hardware & Motion Gated -->
          <div v-else class="vms-flex-col" style="gap: 0.85rem; padding-top: 0.35rem;">
            <div class="vms-flex-col" style="gap: 0.35rem;">
              <div class="vms-flex-between"><span class="vms-text-mono vms-text-xs vms-font-semibold">TAXA DE PROCESSAMENTO (FPS):</span><span class="vms-badge vms-badge-orange">{{ fps }} FPS</span></div>
              <input :value="fps" type="range" min="1" max="50" step="1" style="accent-color: var(--vms-neu-accent-orange); width: 100%;" @input="fps = Number(($event.target as HTMLInputElement).value)" />
            </div>
            <div class="vms-flex-col" style="gap: 0.35rem;">
              <div class="vms-flex-between"><span class="vms-text-mono vms-text-xs vms-font-semibold">TARGET DE HARDWARE:</span><span v-if="!gpuDetected" class="vms-badge vms-badge-secondary" style="font-size: 9px; color: #ff003c;">[GPU INDISPONÍVEL]</span></div>
              <div class="vms-flex-row" style="gap: 0.5rem;">
                <button class="vms-btn vms-btn-sm" :class="hardware === 'cpu_shm' ? 'vms-btn-primary' : 'vms-btn-secondary'" style="flex: 1; padding: 6px;" @click="hardware = 'cpu_shm'">CPU</button>
                <button class="vms-btn vms-btn-sm" :class="hardware === 'rtx_5090_cuda' ? 'vms-btn-primary' : 'vms-btn-secondary'" :disabled="!gpuDetected" style="flex: 1; padding: 6px;" @click="hardware = 'rtx_5090_cuda'">GPU</button>
              </div>
            </div>
            <div class="vms-flex-col" style="gap: 0.4rem; background: rgba(0,0,0,0.25); padding: 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
              <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.5rem; align-items: center; cursor: pointer;">
                <input :checked="motionGated" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" @change="motionGated = ($event.target as HTMLInputElement).checked" />
                <span>Detecção de movimento apenas</span>
              </label>
              <span class="vms-text-mono vms-text-2xs vms-text-dim" style="margin-left: 1.4rem;">Pesa menos na GPU/CPU ignorando frames estáticos sem variação.</span>
            </div>
          </div>
        </div>
      </div>
      <div class="vms-modal-footer vms-flex-between">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="handleSave">SALVAR</button>
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
