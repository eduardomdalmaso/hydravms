<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AlarmItem, AlarmFolderNode, AlarmSensorType } from '../../../types/alarmTree'

const props = defineProps<{
  isOpen: boolean
  targetFolderId?: string
  folders: AlarmFolderNode[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', alarm: Partial<AlarmItem>, folderId?: string): void
}>()

const name = ref('')
const zone = ref('')
const sensorType = ref<AlarmSensorType>('IVS')
const sensitivity = ref(80)
const selectedFolderId = ref(props.targetFolderId || '')
const linkedCamera = ref('')

watch(() => props.isOpen, (open) => {
  if (open) {
    name.value = ''
    zone.value = ''
    sensorType.value = 'IVS'
    sensitivity.value = 80
    selectedFolderId.value = props.targetFolderId || ''
    linkedCamera.value = ''
  }
})

const capitalize = (s: string) => s.trim() ? s.trim().charAt(0).toUpperCase() + s.trim().slice(1) : ''

const handleSave = () => {
  if (!name.value.trim()) return
  emit('save', {
    name: capitalize(name.value),
    zone: zone.value.trim() || 'Zona Geral',
    type: sensorType.value,
    sensitivity: sensitivity.value,
    linkedCameraName: linkedCamera.value
  }, selectedFolderId.value || undefined)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="max-width: 520px; width: 95%;">
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRAR SENSOR DE ALARME</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-flex-col" style="gap: 0.85rem; padding: 0.5rem 0;">
        <div class="vms-form-group">
          <label class="vms-label">Nome do Sensor (Primeira Letra Maiuscula)</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: Sensor Perimetro Norte" autofocus />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Zona / Localizacao</label>
          <input v-model="zone" class="vms-auth-input" placeholder="Ex: Zona 01 // Muro Norte" />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Tipo de Sensor</label>
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <button v-for="t in (['IVS', 'PIR', 'MAG', 'SMK'] as const)" :key="t" class="vms-btn vms-btn-sm" :class="sensorType === t ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="sensorType = t">{{ t }}</button>
          </div>
        </div>

        <div class="vms-form-group">
          <label class="vms-label">Pasta / Zona de Destino</label>
          <select v-model="selectedFolderId" class="vms-auth-input">
            <option value="">[RAIZ] Sem Pasta Definida</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">[ZONA] {{ f.name }}</option>
          </select>
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
