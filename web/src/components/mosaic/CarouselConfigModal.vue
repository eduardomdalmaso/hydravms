<script setup lang="ts">
import { ref } from 'vue'
import type { CarouselConfig } from '../../types/mosaic'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', config: CarouselConfig): void
}>()

const name = ref('Ronda Principal - Entradas')
const interval = ref(10)
const selectedCams = ref<string[]>(['cam_01', 'cam_02'])

const handleSave = () => {
  if (selectedCams.value.length === 0) return
  emit('created', {
    id: `car-${Date.now()}`,
    name: name.value,
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
        <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="emit('close')">[FECHAR]</button>
      </div>
      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">Nome do Carrossel</label>
          <input v-model="name" class="vms-auth-input" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Intervalo de Troca (Segundos)</label>
          <input v-model="interval" type="number" min="3" max="120" class="vms-auth-input" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">Cameras Participantes da Ronda</label>
          <div class="vms-flex-col" style="gap: 0.5rem;">
            <label class="vms-checkbox-label">
              <input type="checkbox" checked class="vms-checkbox" />
              <span>[CAM_01] Portaria Principal</span>
            </label>
            <label class="vms-checkbox-label">
              <input type="checkbox" checked class="vms-checkbox" />
              <span>[CAM_02] Estacionamento Visitantes</span>
            </label>
          </div>
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">ADICIONAR AO MOSAICO</button>
      </div>
    </div>
  </div>
</template>
