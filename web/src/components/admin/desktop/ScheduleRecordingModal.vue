<script setup lang="ts">
import { ref, watch } from 'vue'
import type { RecordingProfile, RecordingMode } from '../../../types/recordingSchedule'
import { createDefaultSchedule } from '../../../types/recordingSchedule'
import WeeklyScheduleGrid from './WeeklyScheduleGrid.vue'

const props = defineProps<{ isOpen: boolean; profile?: RecordingProfile | null }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', profile: RecordingProfile): void }>()

const form = ref<RecordingProfile>({
  id: '', name: '', mode: 'continuous', isActive: true, schedule: createDefaultSchedule(true), preBuffer: 5, postBuffer: 15
})

watch(() => props.isOpen, (open) => {
  if (!open) return
  if (props.profile) form.value = { ...props.profile, schedule: props.profile.schedule.map(r => [...r]) }
  else form.value = { id: `REC_0${Math.floor(Math.random() * 90) + 10}`, name: 'Perfil de Gravacao', mode: 'continuous', isActive: true, schedule: createDefaultSchedule(true), preBuffer: 5, postBuffer: 15 }
})

watch(() => form.value, (newVal) => {
  if (props.profile) {
    props.profile.mode = newVal.mode; props.profile.name = newVal.name
    props.profile.preBuffer = newVal.preBuffer; props.profile.postBuffer = newVal.postBuffer
    props.profile.schedule = newVal.schedule
  }
}, { deep: true })

const handleIntegerInput = (field: 'preBuffer' | 'postBuffer', ev: Event) => {
  const val = (ev.target as HTMLInputElement).value.replace(/\D/g, '')
  form.value[field] = val === '' ? 0 : parseInt(val, 10)
}

const handleSave = () => {
  if (!form.value.name) form.value.name = 'Perfil de Gravacao'
  emit('save', { ...form.value })
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-card" style="max-width: 640px; width: 95%; max-height: 90vh; overflow-y: auto; display: flex; flex-direction: column; gap: 1rem;">
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">CONFIGURACAO DE HORARIOS DE GRAVACAO</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// ID: {{ form.id }} (DEFINA MODOS E MATRIZ TEMPORAL)</span>
        </div>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <!-- Mode Selector -->
      <div class="vms-form-group">
        <label class="vms-label">Tipo de Gravacao</label>
        <select v-model="form.mode" class="vms-auth-input">
          <option value="continuous">[CONTINUA] Gravacao 24/7 Ininterrupta</option>
          <option value="motion">[MOVIMENTO] Deteccao por Movimento (VMD)</option>
          <option value="ai_event">[EVENTO IA] Acionamento por Smart Alarme</option>
        </select>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Nome do Perfil</label>
        <input v-model="form.name" class="vms-auth-input" placeholder="Ex: Horario Noturno 22h-06h" />
      </div>

      <!-- Weekly Schedule Grid -->
      <WeeklyScheduleGrid v-model:schedule="form.schedule" :mode="form.mode" />

      <!-- Advanced Settings (Pre-recording Time & Delay Time Direct Integer Inputs) -->
      <div class="vms-flex-col" style="gap: 0.5rem; background: rgba(255, 255, 255, 0.02); border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.75rem;">
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">// CONFIGURACOES AVANCADAS (BUFFERS)</span>
        <div class="vms-flex-row" style="gap: 0.75rem;">
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Pre-recording Time (s)</label>
            <input type="text" inputmode="numeric" :value="form.preBuffer" class="vms-auth-input" placeholder="Ex: 5" @input="handleIntegerInput('preBuffer', $event)" />
          </div>
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label">Delay Time / Pos (s)</label>
            <input type="text" inputmode="numeric" :value="form.postBuffer" class="vms-auth-input" placeholder="Ex: 15" @input="handleIntegerInput('postBuffer', $event)" />
          </div>
        </div>
      </div>

      <div class="vms-modal-footer" style="padding: 0.75rem 0 0 0; margin-top: 0.5rem; border-top: 1px solid var(--vms-border); display: flex; justify-content: flex-end; gap: 0.75rem; background: transparent;">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
