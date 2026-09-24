<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { AnalyticInstance } from '../../../types/marketplace'
import { fetchCameras } from '../../../services/api'

const props = defineProps<{
  instance: AnalyticInstance
  disabled: boolean
}>()

const emit = defineEmits<{
  (e: 'update', inst: AnalyticInstance): void
}>()

const cameras = ref<{ id: string; name: string }[]>([])
const formName = ref(props.instance.name)
const formCamera = ref(props.instance.camera_id)
const formFps = ref(props.instance.fps_rate || 15)
const formHardware = ref(props.instance.hardware_target || 'rtx_5090_cuda')
const formStream = ref(props.instance.stream_type || 'main_1080p')
const formMotion = ref(props.instance.motion_gated !== false)

onMounted(async () => {
  const cams = await fetchCameras()
  cameras.value = cams.map(c => ({ id: c.id, name: c.name }))
})

watch([formName, formCamera, formFps, formHardware, formStream, formMotion], () => {
  const camObj = cameras.value.find(c => c.id === formCamera.value)
  emit('update', {
    ...props.instance,
    name: formName.value,
    camera_id: formCamera.value,
    camera_name: camObj?.name || props.instance.camera_name,
    fps_rate: formFps.value,
    hardware_target: formHardware.value,
    stream_type: formStream.value,
    motion_gated: formMotion.value
  })
})
</script>

<template>
  <div class="vms-card vms-flex-col" style="padding: 1rem; gap: 0.85rem; background: #0b0e14; border: 1px solid var(--vms-border); border-radius: 8px;">
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem; align-items: center;">
      <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">// PARÂMETROS DO ANALÍTICO</span>
      <span v-if="disabled" class="vms-badge vms-badge-orange" style="font-size: 9px;">[PAUSE PARA EDITAR]</span>
    </div>

    <fieldset :disabled="disabled" class="vms-flex-col" style="gap: 0.75rem; border: none; padding: 0; margin: 0;">
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1.2;">
          <label class="vms-text-xs vms-font-semibold">NOME DO ANALÍTICO:</label>
          <input v-model="formName" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;" />
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1.2;">
          <label class="vms-text-xs vms-font-semibold">CÂMERA ASSOCIADA:</label>
          <select v-model="formCamera" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;">
            <option v-for="c in cameras" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 0.8;">
          <label class="vms-text-xs vms-font-semibold">TAXA (FPS):</label>
          <input v-model.number="formFps" type="number" min="1" max="50" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;" />
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">STREAM:</label>
          <select v-model="formStream" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;">
            <option value="main_1080p">Main (1080p/4K)</option>
            <option value="sub_stream">Sub (480p)</option>
          </select>
        </div>

        <div class="vms-flex-col" style="gap: 0.25rem; flex: 1;">
          <label class="vms-text-xs vms-font-semibold">HARDWARE TARGET:</label>
          <div class="vms-flex-row" style="gap: 4px;">
            <button
              type="button"
              class="vms-btn"
              :class="formHardware === 'cpu_shm' ? 'vms-btn-primary' : 'vms-btn-secondary'"
              style="flex: 1; padding: 4px 6px; font-size: 10px;"
              @click="formHardware = 'cpu_shm'"
            >
              CPU
            </button>
            <button
              type="button"
              class="vms-btn"
              :class="formHardware === 'rtx_5090_cuda' ? 'vms-btn-primary' : 'vms-btn-secondary'"
              style="flex: 1; padding: 4px 6px; font-size: 10px;"
              @click="formHardware = 'rtx_5090_cuda'"
            >
              GPU
            </button>
          </div>
        </div>

        <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.4rem; align-items: center; cursor: pointer; flex: 1.2; padding-top: 14px;">
          <input v-model="formMotion" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" />
          Detecção de movimento apenas
        </label>
      </div>
    </fieldset>
  </div>
</template>
