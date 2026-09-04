<script setup lang="ts">
const props = defineProps<{
  minConfidence: number
  cooldownSeconds: number
  selectedEventType: string
  selectedSeverity: string
}>()

const emit = defineEmits<{
  (e: 'update:minConfidence', val: number): void
  (e: 'update:cooldownSeconds', val: number): void
  (e: 'update:selectedEventType', val: string): void
  (e: 'update:selectedSeverity', val: string): void
}>()
</script>

<template>
  <div class="vms-flex-col">
    <div class="vms-form-group">
      <label class="vms-label">Tipo de Evento Analítico</label>
      <select
        :value="selectedEventType"
        class="vms-select"
        @change="emit('update:selectedEventType', ($event.target as HTMLSelectElement).value)"
      >
        <option value="person_intrusion">[AI] Invasão de Perímetro // Pessoa</option>
        <option value="vehicle_loitering">[AI] Permanência / Ronda // Veículo</option>
        <option value="alpr_match">[AI] Placa Veicular (ALPR Match)</option>
        <option value="fire_smoke">[AI] Fogo e Fumaça (Crítico)</option>
      </select>
    </div>
    <div class="vms-form-group">
      <label class="vms-label">Severidade Mínima</label>
      <select
        :value="selectedSeverity"
        class="vms-select"
        @change="emit('update:selectedSeverity', ($event.target as HTMLSelectElement).value)"
      >
        <option value="critical">Crítico (Alerta Vermelho Imediato)</option>
        <option value="warning">Atenção (Alerta Amarelo)</option>
        <option value="info">Informativo (Todos os Eventos)</option>
      </select>
    </div>
    <div class="vms-flex-row" style="gap: 1rem;">
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Confiança Mínima: {{ Math.round(minConfidence * 100) }}%</label>
        <input
          type="range"
          min="0.5"
          max="0.99"
          step="0.05"
          :value="minConfidence"
          class="vms-input"
          @input="emit('update:minConfidence', parseFloat(($event.target as HTMLInputElement).value))"
        />
      </div>
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Anti-Spam Cooldown (segundos)</label>
        <input
          type="number"
          min="5"
          max="3600"
          :value="cooldownSeconds"
          class="vms-input"
          @input="emit('update:cooldownSeconds', parseInt(($event.target as HTMLInputElement).value) || 30)"
        />
      </div>
    </div>
  </div>
</template>
