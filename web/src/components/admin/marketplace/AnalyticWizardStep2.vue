<script setup lang="ts">
defineProps<{
  fps: number
  motionGated: boolean
  autoSahi: boolean
  scheduleStart: string
  scheduleEnd: string
  selectedDays: number[]
  hardware: 'rtx_5090_cuda' | 'cpu_shm'
  gpuDetected?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:fps', val: number): void
  (e: 'update:motionGated', val: boolean): void
  (e: 'update:autoSahi', val: boolean): void
  (e: 'update:scheduleStart', val: string): void
  (e: 'update:scheduleEnd', val: string): void
  (e: 'toggle:day', day: number): void
  (e: 'update:hardware', val: 'rtx_5090_cuda' | 'cpu_shm'): void
}>()

const days = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.15rem;">
    <!-- Slider de FPS -->
    <div class="vms-flex-col" style="gap: 0.4rem;">
      <div class="vms-flex-between" style="align-items: center;">
        <label class="vms-text-xs vms-font-semibold">TAXA DE PROCESSAMENTO (FPS):</label>
        <span class="vms-badge vms-badge-orange" style="font-size: 11px;">{{ fps }} FPS</span>
      </div>
      <input :value="fps" type="range" min="1" max="50" step="1" style="accent-color: var(--vms-neu-accent-orange); width: 100%; cursor: pointer;" @input="emit('update:fps', Number(($event.target as HTMLInputElement).value))" />
    </div>

    <!-- Modos de Otimização e Detecção -->
    <div class="vms-flex-col" style="gap: 0.65rem; background: rgba(0,0,0,0.25); padding: 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.5rem; align-items: center; cursor: pointer;">
        <input :checked="motionGated" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" @change="emit('update:motionGated', ($event.target as HTMLInputElement).checked)" /> Detecção de movimento apenas
      </label>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="margin-left: 1.5rem;">
        {{ motionGated ? 'Pesa menos na GPU/CPU ignorando frames sem alteração.' : 'Analisa frame a frame continuamente 24/7.' }}
      </span>

      <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.5rem; align-items: center; cursor: pointer; margin-top: 0.35rem;">
        <input :checked="autoSahi" type="checkbox" style="accent-color: #00ff9d;" @change="emit('update:autoSahi', ($event.target as HTMLInputElement).checked)" /> Auto-SAHI Adaptativo
      </label>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="margin-left: 1.5rem;">Fatia automaticamente o frame para detectar alvos pequenos/distantes com alta precisão.</span>
    </div>

    <!-- Horário de Ativação -->
    <div class="vms-flex-col" style="gap: 0.4rem;">
      <label class="vms-text-xs vms-font-semibold">HORÁRIO DE ATIVAÇÃO:</label>
      <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
        <input :value="scheduleStart" type="time" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; width: 110px;" @input="emit('update:scheduleStart', ($event.target as HTMLInputElement).value)" />
        <span class="vms-text-xs vms-text-dim">até</span>
        <input :value="scheduleEnd" type="time" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; width: 110px;" @input="emit('update:scheduleEnd', ($event.target as HTMLInputElement).value)" />
      </div>
      <div class="vms-flex-row" style="gap: 0.35rem; margin-top: 0.3rem;">
        <span v-for="(dName, dIdx) in days" :key="dIdx" class="vms-chip" :class="{ active: selectedDays.includes(dIdx) }" @click="emit('toggle:day', dIdx)">{{ dName }}</span>
      </div>
    </div>

    <!-- Hardware Target Dinâmico -->
    <div class="vms-flex-row" style="gap: 0.75rem;">
      <div class="vms-flex-col" style="gap: 0.2rem; flex: 1;">
        <div class="vms-flex-between">
          <label class="vms-text-xs vms-font-semibold">TARGET DE HARDWARE:</label>
          <span v-if="!gpuDetected" class="vms-badge vms-badge-secondary" style="font-size: 9px; color: #ff003c;">[GPU INDISPONÍVEL]</span>
        </div>
        <select :value="hardware" class="vms-auth-input" style="padding: 5px 8px; font-size: 11px;" @change="emit('update:hardware', ($event.target as HTMLSelectElement).value as any)">
          <option value="rtx_5090_cuda" :disabled="!gpuDetected">GPU RTX 5090 (CUDA / TensorRT)</option>
          <option value="cpu_shm">CPU SHM (ONNX Runtime)</option>
        </select>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-chip { font-size: 10px; font-family: var(--vms-font-mono); padding: 3px 8px; border-radius: 4px; background: rgba(255,255,255,0.06); color: #8b94a0; cursor: pointer; border: 1px solid transparent; }
.vms-chip.active { background: rgba(255,94,58,0.2); color: #ff5e3a; border-color: rgba(255,94,58,0.4); font-weight: bold; }
</style>
