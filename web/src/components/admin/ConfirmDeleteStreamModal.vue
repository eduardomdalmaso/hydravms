<script setup lang="ts">
defineProps<{
  isOpen: boolean
  streamName: string
  streamId: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm'): void
}>()
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click="emit('close')">
    <div class="vms-modal-card" style="max-width: 440px;" @click.stop>
      <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="color: var(--vms-neu-accent-magenta, #ff003c);">TEM CERTEZA DA ACAO?</h4>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// AVISO DE EXCLUSAO DE FLUXO</span>
        </div>
        <button class="vms-btn vms-btn-ghost" @click="emit('close')">X</button>
      </div>

      <div class="vms-flex-col" style="gap: 1rem; margin-top: 1rem;">
        <div style="background: rgba(255, 0, 60, 0.08); border: 1px solid rgba(255, 0, 60, 0.3); border-radius: 6px; padding: 0.85rem;">
          <p style="color: #fff; font-size: 13px; line-height: 1.4; margin: 0;">
            Deseja realmente remover o fluxo de camera <strong style="color: #ff003c;">"{{ streamName }}"</strong> (ID: {{ streamId }})?
          </p>
          <p class="vms-text-dim" style="font-size: 11px; margin-top: 0.5rem; margin-bottom: 0;">
            Ao confirmar, a ingestao e as regras de analiticos vinculadas a esta camera serao desativadas.
          </p>
        </div>

        <div class="vms-flex-row" style="justify-content: flex-end; gap: 0.75rem; margin-top: 0.5rem;">
          <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
          <button class="vms-btn" style="background: #ff003c; color: #fff; border: 1px solid #ff003c;" @click="emit('confirm')">SALVAR</button>
        </div>
      </div>
    </div>
  </div>
</template>
