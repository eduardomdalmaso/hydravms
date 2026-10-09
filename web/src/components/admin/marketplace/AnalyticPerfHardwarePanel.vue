<script setup lang="ts">
defineProps<{
  fps: number
  hardware: 'rtx_5090_cuda' | 'cpu_shm'
  motionGated: boolean
  gpuDetected: boolean
}>()

const emit = defineEmits<{
  (e: 'update:fps', val: number): void
  (e: 'update:hardware', val: 'rtx_5090_cuda' | 'cpu_shm'): void
  (e: 'update:motionGated', val: boolean): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.85rem; padding-top: 0.35rem;">
    <!-- FPS Rate -->
    <div class="vms-flex-col" style="gap: 0.35rem;">
      <div class="vms-flex-between">
        <span class="vms-text-mono vms-text-xs vms-font-semibold">TAXA DE PROCESSAMENTO (FPS):</span>
        <span class="vms-badge vms-badge-orange">{{ fps }} FPS</span>
      </div>
      <input
        :value="fps"
        type="range"
        min="1"
        max="50"
        step="1"
        style="accent-color: var(--vms-neu-accent-orange); width: 100%;"
        @input="emit('update:fps', Number(($event.target as HTMLInputElement).value))"
      />
    </div>

    <!-- Hardware Target -->
    <div class="vms-flex-col" style="gap: 0.35rem;">
      <div class="vms-flex-between" style="align-items: center;">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-mono vms-text-xs vms-font-semibold">TARGET DE HARDWARE:</span>
          <span class="vms-info-tip" title="GPU: Aceleração NVIDIA via NVDEC (decodificação de vídeo na VRAM), CUDA / TensorRT (inferência de IA) e NVENC (gravação/exportação). CPU: Zero-Copy (/dev/shm).">?</span>
        </div>
        <span v-if="!gpuDetected" class="vms-badge vms-badge-secondary" style="font-size: 9px; color: #ff003c;">
          [GPU INDISPONÍVEL]
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <button
          class="vms-btn vms-btn-sm"
          :class="hardware === 'cpu_shm' ? 'vms-btn-primary' : 'vms-btn-secondary'"
          style="flex: 1; padding: 6px;"
          @click="emit('update:hardware', 'cpu_shm')"
        >
          CPU
        </button>
        <button
          class="vms-btn vms-btn-sm"
          :class="hardware === 'rtx_5090_cuda' ? 'vms-btn-primary' : 'vms-btn-secondary'"
          :disabled="!gpuDetected"
          style="flex: 1; padding: 6px;"
          @click="emit('update:hardware', 'rtx_5090_cuda')"
        >
          GPU
        </button>
      </div>
    </div>

    <!-- Motion Gated -->
    <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 0.65rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.5rem; align-items: center; cursor: pointer; color: #fff;">
        <input
          :checked="motionGated"
          type="checkbox"
          style="accent-color: var(--vms-neu-accent-orange);"
          @change="emit('update:motionGated', ($event.target as HTMLInputElement).checked)"
        />
        <span>Detecção de movimento apenas</span>
      </label>
      <span class="vms-info-tip" title="Pesa menos na GPU/CPU ignorando frames estáticos sem variação.">?</span>
    </div>
  </div>
</template>

<style scoped>
.vms-info-tip { display: inline-flex; align-items: center; justify-content: center; width: 14px; height: 14px; border-radius: 50%; background: rgba(255, 255, 255, 0.08); color: #8b94a0; font-size: 10px; font-family: var(--vms-font-jetbrains); font-weight: 700; cursor: help; transition: all 0.2s ease; user-select: none; }
.vms-info-tip:hover { background: var(--vms-neu-accent-orange); color: #ffffff; transform: scale(1.1); }
</style>
