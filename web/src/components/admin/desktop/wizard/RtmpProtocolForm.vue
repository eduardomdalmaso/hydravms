<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  streamKey?: string
}
const props = withDefaults(defineProps<Props>(), {
  streamKey: 'stream_alpha_01'
})
const emit = defineEmits<{ (e: 'update:streamKey', val: string): void }>()

const publishUrl = computed(() => `rtmp://${window.location.hostname || 'localhost'}:1935/live/${props.streamKey || 'stream_key'}`)

const generateKey = () => {
  const rand = 'stream_' + Math.random().toString(36).substring(2, 9)
  emit('update:streamKey', rand)
}

const copyUrl = () => {
  navigator.clipboard.writeText(publishUrl.value)
}
</script>

<template>
  <div class="vms-flex-col" style="min-height: 180px; gap: 0.55rem; justify-content: space-between;">
    <div class="vms-form-group">
      <label class="vms-label">Chave de Transmissão (Stream Key)</label>
      <div class="vms-flex-row" style="gap: 0.4rem;">
        <input 
          :value="streamKey" 
          class="vms-auth-input vms-text-mono" 
          style="flex: 1;" 
          placeholder="Ex: drone_alpha_01" 
          @input="emit('update:streamKey', ($event.target as HTMLInputElement).value)" 
        />
        <button 
          class="vms-btn vms-btn-secondary" 
          type="button"
          style="padding: 0 12px; font-size: 11px; white-space: nowrap;" 
          @click="generateKey"
        >
          GERAR CHAVE
        </button>
      </div>
    </div>

    <div class="vms-form-group">
      <label class="vms-label">URL de Publicação RTMP (Copie para o Drone / OBS / App)</label>
      <div class="vms-flex-row" style="gap: 0.4rem;">
        <input 
          :value="publishUrl" 
          readonly 
          class="vms-auth-input vms-text-mono" 
          style="flex: 1; background: rgba(0,0,0,0.5); color: var(--vms-neu-accent-cyan);" 
        />
        <button 
          class="vms-btn vms-btn-secondary" 
          type="button"
          style="padding: 0 12px; font-size: 11px; white-space: nowrap;" 
          title="Copiar URL" 
          @click="copyUrl"
        >
          COPIAR
        </button>
      </div>
    </div>

    <div style="background: rgba(0, 240, 255, 0.08); color: var(--vms-neu-accent-cyan); border: 1px solid rgba(0, 240, 255, 0.2); font-size: 11px; padding: 0.45rem 0.65rem; border-radius: 4px; line-height: 1.4;">
      <strong style="color: #fff;">INGESTAO PUSH:</strong> Transmita o sinal para a porta :1935 do servidor usando a URL acima.
    </div>
  </div>
</template>
