<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchSampleVideos } from '../../../../services/adminApi'

const props = defineProps<{ url: string }>()
const emit = defineEmits<{ (e: 'update:url', val: string): void }>()

const sampleList = ref<{ name: string; path: string }[]>([
  { name: 'estacionamento_shopping.mp4', path: '/home/hades/Documents/HydraStream/samples/estacionamento_shopping.mp4' },
  { name: 'portaria_predial.mp4', path: '/home/hades/Documents/HydraStream/samples/portaria_predial.mp4' },
  { name: 'avenida_transito.mp4', path: '/home/hades/Documents/HydraStream/samples/avenida_transito.mp4' }
])

onMounted(async () => {
  const loaded = await fetchSampleVideos()
  if (loaded && loaded.length > 0) sampleList.value = loaded
})

const onSelectSample = (event: Event) => {
  const target = event.target as HTMLSelectElement
  if (target.value) emit('update:url', target.value)
}
</script>

<template>
  <div class="vms-flex-col" style="min-height: 180px; gap: 0.55rem; justify-content: space-between;">
    <div class="vms-form-group">
      <label class="vms-label">Arquivo de Amostra Pré-Gravado (Loop de Teste)</label>
      <select class="vms-auth-input" style="height: 38px;" @change="onSelectSample">
        <option value="">[SELECIONE UM VÍDEO DA BIBLIOTECA]</option>
        <option v-for="s in sampleList" :key="s.name" :value="s.name" :selected="url === s.name || url === s.path">
          {{ s.name }}
        </option>
      </select>
    </div>

    <div class="vms-form-group">
      <label class="vms-label">Caminho do Arquivo (.mp4 / .mkv / .avi)</label>
      <input 
        :value="url" 
        class="vms-auth-input vms-text-mono" 
        placeholder="Ex: estacionamento_shopping.mp4 ou /storage/videos/camera1.mp4" 
        @input="emit('update:url', ($event.target as HTMLInputElement).value)" 
      />
    </div>

    <div style="background: rgba(0, 240, 255, 0.08); color: var(--vms-neu-accent-cyan); border: 1px solid rgba(0, 240, 255, 0.2); font-size: 11px; padding: 0.45rem 0.65rem; border-radius: 4px; line-height: 1.4;">
      <strong style="color: #fff;">FLUXO VIRTUAL EM LOOP:</strong> O HydraStream reproduzirá este vídeo em repetição contínua 24/7 na velocidade real do arquivo.
    </div>
  </div>
</template>
