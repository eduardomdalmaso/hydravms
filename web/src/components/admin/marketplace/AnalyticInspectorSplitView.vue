<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { AnalyticInstance } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'
import { formatPluginId } from '../../../utils/idFormatter'

const props = defineProps<{ instance: AnalyticInstance }>()
const emit = defineEmits<{
  (e: 'saved', inst: AnalyticInstance): void
  (e: 'delete', id: string): void
  (e: 'close'): void
}>()

const formName = ref(props.instance.name), formCamera = ref(props.instance.camera_id)
const formStream = ref(props.instance.stream_type), formHardware = ref(props.instance.hardware_target)
const formConfidence = ref(props.instance.confidence_threshold), formSahi = ref(!!props.instance.auto_sahi)
const formFps = ref(props.instance.fps_rate || 15), formMotion = ref(props.instance.motion_gated !== false)
const formActive = ref(props.instance.is_active)
const cameras = ref<{ id: string; name: string }[]>([])

onMounted(async () => {
  const cams = await fetchCameras()
  cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
})

watch(() => props.instance, (i) => {
  formName.value = i.name; formCamera.value = i.camera_id; formStream.value = i.stream_type
  formHardware.value = i.hardware_target; formConfidence.value = i.confidence_threshold
  formSahi.value = !!i.auto_sahi; formFps.value = i.fps_rate || 15
  formMotion.value = i.motion_gated !== false; formActive.value = i.is_active
})

const toggleActiveState = () => {
  formActive.value = !formActive.value
  const updated: AnalyticInstance = { ...props.instance, is_active: formActive.value }
  emit('saved', updated)
}

const handleSave = () => {
  if (props.instance.is_active) return
  const camObj = cameras.value.find(c => c.id === formCamera.value)
  const updated: AnalyticInstance = {
    ...props.instance, name: formName.value, camera_id: formCamera.value,
    camera_name: camObj?.name || props.instance.camera_name, stream_type: formStream.value,
    hardware_target: formHardware.value, confidence_threshold: formConfidence.value,
    fps_rate: formFps.value, motion_gated: formMotion.value, auto_sahi: formSahi.value,
    is_active: formActive.value
  }
  emit('saved', updated)
}
</script>

<template>
  <div class="vms-card" style="padding: 1.25rem; display: flex; flex-direction: column; gap: 1rem; max-width: 680px; margin: 0 auto; width: 100%;">
    <!-- Header -->
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem; align-items: center;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-sm vms-font-semibold" style="color: var(--vms-neu-accent-orange);">INSPEÇÃO // {{ formName }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ formatPluginId(instance.id) }} // {{ instance.plugin_name }}</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-sm" :class="formActive ? 'vms-btn-secondary' : 'vms-btn-primary'" style="font-weight: bold;" @click="toggleActiveState">
          {{ formActive ? 'PAUSAR' : 'RETOMAR' }}
        </button>
        <button class="vms-btn vms-btn-danger vms-btn-sm" @click="emit('delete', instance.id)">EXCLUIR</button>
        <button :disabled="formActive" class="vms-btn vms-btn-primary vms-btn-sm" :style="{ opacity: formActive ? 0.4 : 1, cursor: formActive ? 'not-allowed' : 'pointer' }" style="font-weight: bold;" @click="handleSave">SALVAR</button>
      </div>
    </div>

    <!-- Banner de Bloqueio em Execução -->
    <div v-if="formActive" class="vms-text-mono vms-text-2xs" style="padding: 0.6rem 0.85rem; border-radius: 4px; background: rgba(252, 238, 10, 0.1); border: 1px solid rgba(252, 238, 10, 0.3); color: #fcee0a;">
      ⚠️ [AVISO] ANALÍTICO EM EXECUÇÃO: Clique em <strong>PAUSAR</strong> no topo para liberar a edição de parâmetros e geometria.
    </div>

    <!-- Campos de Formulário (Desabilitados se ativo) -->
    <fieldset :disabled="formActive" class="vms-flex-col" style="gap: 0.75rem; border: none; padding: 0; margin: 0;">
      <div class="vms-flex-col" style="gap: 0.25rem;">
        <label class="vms-text-xs vms-font-semibold">NOME DO ANALÍTICO:</label>
        <input v-model="formName" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;" />
      </div>

      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1.2;">
          <label class="vms-text-xs vms-font-semibold">CÂMERA ASSOCIADA:</label>
          <select v-model="formCamera" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;">
            <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 0.8;">
          <label class="vms-text-xs vms-font-semibold">TAXA (FPS):</label>
          <input v-model.number="formFps" type="number" min="1" max="50" class="vms-auth-input" style="padding: 6px 10px; font-size: 11px;" />
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">STREAM:</label>
          <select v-model="formStream" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="main_1080p">Main (1080p/4K)</option><option value="sub_stream">Sub (480p)</option></select>
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">TARGET:</label>
          <select v-model="formHardware" class="vms-auth-input" style="padding: 6px 8px; font-size: 11px;"><option value="rtx_5090_cuda">GPU RTX 5090</option><option value="cpu_shm">CPU SHM</option></select>
        </div>
      </div>

      <div class="vms-flex-between" style="padding-top: 0.25rem;">
        <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.4rem; align-items: center; cursor: pointer;">
          <input v-model="formMotion" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" /> Detecção de movimento apenas
        </label>
        <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.4rem; align-items: center; cursor: pointer;">
          <input v-model="formSahi" type="checkbox" style="accent-color: #00ff9d;" /> Auto-SAHI Adaptativo
        </label>
      </div>
    </fieldset>
  </div>
</template>
