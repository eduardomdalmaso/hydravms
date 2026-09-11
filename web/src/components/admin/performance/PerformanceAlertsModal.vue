<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', config: { cpuLimit: number; vramMinPercent: number; ramLimit: number }): void
}>()

const cpuLimit = ref(80)
const vramMinPercent = ref(10)
const ramLimit = ref(85)
const errorMsg = ref<string | null>(null)

const enforcePercent = (val: any, min = 1, max = 99, fallback = 50) => {
  const num = parseInt(val)
  if (isNaN(num)) return fallback
  return Math.min(max, Math.max(min, num))
}

const onInput = (target: 'cpu' | 'vram' | 'ram', e: Event) => {
  errorMsg.value = null
  const input = e.target as HTMLInputElement
  const clamped = enforcePercent(input.value)
  if (target === 'cpu') cpuLimit.value = clamped
  else if (target === 'vram') vramMinPercent.value = clamped
  else ramLimit.value = clamped
  input.value = String(clamped)
}

const handleSave = () => {
  if (cpuLimit.value < 1 || cpuLimit.value > 99 ||
      vramMinPercent.value < 1 || vramMinPercent.value > 99 ||
      ramLimit.value < 1 || ramLimit.value > 99) {
    errorMsg.value = 'Os valores devem ser porcentagens válidas entre 1% e 99%.'
    return
  }
  emit('save', {
    cpuLimit: enforcePercent(cpuLimit.value, 1, 99, 80),
    vramMinPercent: enforcePercent(vramMinPercent.value, 1, 99, 10),
    ramLimit: enforcePercent(ramLimit.value, 1, 99, 85)
  })
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 460px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CONFIGURAÇÃO DE ALERTAS DE PERFORMANCE</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">✕</button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 12px;">
        <div v-if="errorMsg" class="vms-badge vms-badge-orange" style="padding: 6px 10px; font-weight: bold;">
          {{ errorMsg }}
        </div>

        <div class="vms-form-group">
          <label class="vms-label">LIMITE CRÍTICO DE CPU // ALERTA QUANDO EXCEDER (1% A 99%)</label>
          <input :value="cpuLimit" type="number" min="1" max="99" class="vms-auth-input" @input="onInput('cpu', $event)" />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">LIMITE DE VRAM GPU // ALERTA QUANDO LIVRE MENOS DE (1% A 99%)</label>
          <input :value="vramMinPercent" type="number" min="1" max="99" class="vms-auth-input" @input="onInput('vram', $event)" />
        </div>

        <div class="vms-form-group">
          <label class="vms-label">LIMITE DE MEMÓRIA RAM // ALERTA QUANDO EXCEDER (1% A 99%)</label>
          <input :value="ramLimit" type="number" min="1" max="99" class="vms-auth-input" @input="onInput('ram', $event)" />
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
