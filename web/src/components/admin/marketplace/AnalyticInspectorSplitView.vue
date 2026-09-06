<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AnalyticInstance } from '../../../types/marketplace'

const props = defineProps<{ instance: AnalyticInstance }>()
const emit = defineEmits<{
  (e: 'saved', inst: AnalyticInstance): void
  (e: 'delete', id: string): void
  (e: 'close'): void
}>()

const formName = ref(props.instance.name), formCamera = ref(props.instance.camera_id)
const formStream = ref(props.instance.stream_type), formHardware = ref(props.instance.hardware_target)
const formConfidence = ref(props.instance.confidence_threshold), formSahi = ref(!!props.instance.specific_params?.sahi_enabled)
const formActive = ref(props.instance.is_active)

watch(() => props.instance, (i) => {
  formName.value = i.name; formCamera.value = i.camera_id; formStream.value = i.stream_type
  formHardware.value = i.hardware_target; formConfidence.value = i.confidence_threshold
  formSahi.value = !!i.specific_params?.sahi_enabled; formActive.value = i.is_active
})

const cameras = [
  { id: 'cam_portaria_01', name: 'CAM_01 // Portaria Principal' },
  { id: 'cam_garagem_02', name: 'CAM_02 // Garagem Subsolo' },
  { id: 'cam_hall_03', name: 'CAM_03 // Catracas Recepção' },
  { id: 'cam_galpao_04', name: 'CAM_04 // Galpão Logística B' }
]

const handleSave = () => {
  const camObj = cameras.find(c => c.id === formCamera.value)
  const updated: AnalyticInstance = {
    ...props.instance, name: formName.value, camera_id: formCamera.value,
    camera_name: camObj?.name || props.instance.camera_name, stream_type: formStream.value,
    hardware_target: formHardware.value, confidence_threshold: formConfidence.value,
    is_active: formActive.value, specific_params: { ...props.instance.specific_params, sahi_enabled: formSahi.value }
  }
  emit('saved', updated)
}
</script>

<template>
  <div class="vms-card" style="padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; max-width: 680px; margin: 0 auto; width: 100%;">
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-sm vms-font-semibold" style="color: var(--vms-neu-accent-orange);">INSPEÇÃO // {{ formName }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ID: {{ instance.id }} // {{ instance.plugin_name }}</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <button class="vms-btn vms-btn-danger vms-btn-sm" @click="emit('delete', instance.id)">EXCLUIR</button>
        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="handleSave">SALVAR</button>
      </div>
    </div>

    <div class="vms-flex-col" style="gap: 0.75rem;">
      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">NOME DO ANALÍTICO:</label>
        <input v-model="formName" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;" />
      </div>

      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">CÂMERA ASSOCIADA:</label>
        <select v-model="formCamera" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;"><option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option></select>
      </div>

      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">STREAM:</label>
          <select v-model="formStream" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="main_1080p">Main (1080p/4K)</option><option value="sub_stream">Sub (480p)</option></select>
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">TARGET:</label>
          <select v-model="formHardware" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="rtx_5090_cuda">RTX 5090</option><option value="cpu_shm">CPU SHM</option></select>
        </div>
      </div>


      <div class="vms-flex-col" style="gap: 0.25rem;">
        <div class="vms-flex-between">
          <label class="vms-text-xs vms-font-semibold">LIMIAR DE CONFIANÇA:</label>
          <span class="vms-text-mono vms-text-xs" style="color: var(--vms-neu-accent-orange);">{{ Math.round(formConfidence * 100) }}%</span>
        </div>
        <input v-model.number="formConfidence" type="range" min="0.1" max="1.0" step="0.05" style="accent-color: var(--vms-neu-accent-orange);" />
      </div>

      <div class="vms-flex-between" style="padding-top: 0.25rem;">
        <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.4rem; align-items: center; cursor: pointer;">
          <input v-model="formSahi" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" /> Habilitar SAHI
        </label>
        <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.4rem; align-items: center; cursor: pointer;">
          <input v-model="formActive" type="checkbox" style="accent-color: #00ff9d;" /> Instância Ativa
        </label>
      </div>

    </div>
  </div>
</template>
