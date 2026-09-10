<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { AnalyticInstance, AnalyticFolderNode, PluginManifest } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'

const props = defineProps<{
  isOpen: boolean; plugin: PluginManifest; folders: AnalyticFolderNode[]; currentFolderId?: string | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', payload: { inst: AnalyticInstance; targetFolderId: string | null }): void
}>()

const name = ref(''), camera = ref(''), stream = ref<'main_1080p' | 'sub_stream'>('main_1080p')
const hardware = ref<'rtx_5090_cuda' | 'cpu_shm'>('rtx_5090_cuda')
const confidence = ref(0.70), sahi = ref(true)
const selectedFolder = ref<string | null>(props.currentFolderId || null)
const cameras = ref<{ id: string; name: string }[]>([])

onMounted(async () => {
  const cams = await fetchCameras()
  cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
  if (cameras.value.length > 0) camera.value = cameras.value[0].id
})

const handleConfirm = () => {
  const camObj = cameras.value.find(c => c.id === camera.value)
  const newInst: AnalyticInstance = {
    id: `inst-${Date.now()}`, plugin_id: props.plugin.id, plugin_name: props.plugin.name,
    name: name.value.trim() || `${props.plugin.name} // ${camObj?.name || 'Câmera'}`,
    camera_id: camera.value || 'cam_stream', camera_name: camObj?.name || 'Câmera',
    stream_type: stream.value, hardware_target: hardware.value,
    confidence_threshold: confidence.value, roi_mode: 'full_frame',
    specific_params: { sahi_enabled: sahi.value },
    is_active: true, fps_rate: 30.0, detections_count: 0, created_at: new Date().toISOString()
  }
  emit('save', { inst: newInst, targetFolderId: selectedFolder.value })
  name.value = ''
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content" style="max-width: 500px; width: 100%;">
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">NOVO ANALÍTICO // {{ plugin.name }}</h4>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">WIZARD DE INSTANCIAÇÃO</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 0.85rem; padding: 1.25rem 0;">
        <div class="vms-flex-col" style="gap: 0.25rem;">
          <label class="vms-text-xs vms-font-semibold">NOME DO ANALÍTICO:</label>
          <input v-model="name" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;" placeholder="Ex: LPR Pista 01..." />
        </div>

        <div class="vms-flex-col" style="gap: 0.25rem;">
          <label class="vms-text-xs vms-font-semibold">CÂMERA DE VÍDEO:</label>
          <select v-model="camera" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;">
            <option v-if="cameras.length === 0" value="">[NENHUMA CÂMERA DISPONÍVEL]</option>
            <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>

        <div class="vms-flex-col" style="gap: 0.25rem;">
          <label class="vms-text-xs vms-font-semibold">PASTA DE DESTINO:</label>
          <select v-model="selectedFolder" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;">
            <option :value="null">[RAIZ DO DESKTOP]</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
          </select>
        </div>

        <div class="vms-flex-row" style="gap: 0.75rem;">
          <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
            <label class="vms-text-xs vms-font-semibold">STREAM:</label>
            <select v-model="stream" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="main_1080p">Main Stream</option><option value="sub_stream">Sub Stream</option></select>
          </div>
          <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
            <label class="vms-text-xs vms-font-semibold">TARGET:</label>
            <select v-model="hardware" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="rtx_5090_cuda">RTX 5090</option><option value="cpu_shm">CPU SHM</option></select>
          </div>
        </div>
      </div>

      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.85rem;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="handleConfirm">SALVAR</button>
      </div>
    </div>
  </div>
</template>
