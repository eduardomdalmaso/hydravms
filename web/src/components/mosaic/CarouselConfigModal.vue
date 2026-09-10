<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { CarouselConfig } from '../../types/mosaic'
import { fetchCameras } from '../../services/api'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', config: CarouselConfig): void
}>()

const name = ref('')
const interval = ref(10)
const selectedCams = ref<string[]>([])
const availableCams = ref<{ id: string; name: string }[]>([])

onMounted(async () => {
  const cams = await fetchCameras()
  availableCams.value = cams.map(c => ({ id: c.id, name: c.name }))
})

const toggleCam = (id: string) => {
  if (selectedCams.value.includes(id)) {
    selectedCams.value = selectedCams.value.filter(c => c !== id)
  } else {
    selectedCams.value.push(id)
  }
}

const handleSave = () => {
  if (selectedCams.value.length === 0 || !name.value.trim()) return
  emit('created', {
    id: `car-${Date.now()}`,
    name: name.value.trim(),
    camera_ids: selectedCams.value,
    interval_seconds: interval.value
  })
  emit('close')
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CONFIGURAR CARROSSEL // RONDA VIRTUAL</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome do Carrossel</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: Ronda Geral Entradas" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Intervalo de Troca (Segundos)</label>
          <input v-model.number="interval" type="number" min="3" max="120" class="vms-auth-input" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Cameras Participantes da Ronda</label>
          <div v-if="availableCams.length === 0" class="vms-text-dim vms-text-xs" style="padding: 0.5rem 0;">
            Nenhuma câmera disponível
          </div>
          <div v-else class="vms-flex-col" style="gap: 0.5rem; max-height: 140px; overflow-y: auto;">
            <label v-for="c in availableCams" :key="c.id" class="vms-checkbox-label">
              <input type="checkbox" :checked="selectedCams.includes(c.id)" class="vms-checkbox" @change="toggleCam(c.id)" />
              <span>[{{ c.id.toUpperCase() }}] {{ c.name }}</span>
            </label>
          </div>
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" :disabled="selectedCams.length === 0 || !name.trim()" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
