<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { PluginManifest, AnalyticInstance } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'

const props = defineProps<{ installedPlugins: PluginManifest[]; initialPluginId?: string | null }>()
const emit = defineEmits<{ (e: 'save', instance: AnalyticInstance): void }>()

const selectedPluginId = ref(props.initialPluginId || (props.installedPlugins[0]?.id || ''))
const instanceName = ref(''), selectedCamera = ref('')
const streamType = ref<'main_1080p' | 'sub_stream'>('main_1080p')
const hardwareTarget = ref<'rtx_5090_cuda' | 'cpu_shm'>('rtx_5090_cuda')
const confidence = ref(0.70), sahiEnabled = ref(true)
const cameras = ref<{ id: string; name: string }[]>([])

watch(() => props.initialPluginId, (newId) => { if (newId) selectedPluginId.value = newId })
const activePlugin = computed(() => props.installedPlugins.find(p => p.id === selectedPluginId.value))

onMounted(async () => {
  const cams = await fetchCameras()
  cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
  if (cameras.value.length > 0 && !selectedCamera.value) selectedCamera.value = cameras.value[0].id
})

const handleSave = () => {
  if (!activePlugin.value) return
  const camObj = cameras.value.find(c => c.id === selectedCamera.value)
  const newInst: AnalyticInstance = {
    id: `inst-${Date.now()}`, plugin_id: activePlugin.value.id, plugin_name: activePlugin.value.name,
    name: instanceName.value.trim() || `${activePlugin.value.name} // ${camObj?.name || 'Câmera'}`,
    camera_id: selectedCamera.value || 'cam_stream', camera_name: camObj?.name || 'Câmera',
    stream_type: streamType.value, hardware_target: hardwareTarget.value,
    confidence_threshold: confidence.value, roi_mode: 'full_frame',
    specific_params: { sahi_enabled: sahiEnabled.value },
    is_active: true, fps_rate: 30.0, detections_count: 0, created_at: new Date().toISOString()
  }
  emit('save', newInst); instanceName.value = ''
}
</script>

<template>
  <div class="vms-card" style="padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem;">
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <span class="vms-text-sm vms-font-semibold" style="color: var(--vms-neu-accent-orange);">NOVO ANALÍTICO // INSTANCIAÇÃO</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">GPU CUDA ACCEL</span>
    </div>

    <div class="vms-flex-col" style="gap: 0.85rem;">
      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">ANALÍTICO / PLUGIN:</label>
        <select v-model="selectedPluginId" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;"><option v-for="p in installedPlugins" :key="p.id" :value="p.id">[{{ p.name }}] // v{{ p.version }}</option></select>
      </div>

      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">NOME DA INSTÂNCIA (OPCIONAL):</label>
        <input v-model="instanceName" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;" placeholder="Ex: LPR Entrada Faixa 1" />
      </div>

      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">CÂMERA DE VÍDEO:</label>
        <select v-model="selectedCamera" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;">
          <option v-if="cameras.length === 0" value="">[NENHUMA CÂMERA DISPONÍVEL]</option>
          <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
        </select>
      </div>

      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">FLUXO DE VÍDEO:</label>
          <select v-model="streamType" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="main_1080p">Main Stream (1080p/4K)</option><option value="sub_stream">Sub Stream (480p Leve)</option></select>
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">EXECUÇÃO:</label>
          <select v-model="hardwareTarget" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="rtx_5090_cuda">NVIDIA RTX 5090</option><option value="cpu_shm">CPU SHM Zero-Copy</option></select>
        </div>
      </div>

      <div class="vms-flex-col" style="gap: 0.25rem;">
        <div class="vms-flex-between">
          <label class="vms-text-xs vms-font-semibold">LIMIAR DE CONFIANÇA:</label>
          <span class="vms-text-mono vms-text-xs" style="color: var(--vms-neu-accent-orange);">{{ Math.round(confidence * 100) }}%</span>
        </div>
        <input v-model.number="confidence" type="range" min="0.1" max="1.0" step="0.05" style="accent-color: var(--vms-neu-accent-orange);" />
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center; margin-top: 0.25rem;">
        <input id="sahi-check" v-model="sahiEnabled" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" />
        <label for="sahi-check" class="vms-text-xs vms-font-semibold" style="cursor: pointer;">Habilitar SAHI (Slicing Aided Hyper Inference)</label>
      </div>

      <button class="vms-btn vms-btn-primary" style="margin-top: 0.5rem; font-weight: bold;" @click="handleSave">SALVAR</button>
    </div>
  </div>
</template>
