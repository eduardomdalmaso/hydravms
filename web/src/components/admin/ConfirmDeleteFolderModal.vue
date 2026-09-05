<script setup lang="ts">
defineProps<{
  isOpen: boolean
  folderName: string
  itemCount: number
  itemType?: string
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
          <h4 class="vms-h4" style="color: var(--vms-neu-accent-orange);">TEM CERTEZA DA ACAO?</h4>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// AVISO DE EXCLUSAO DE PASTA</span>
        </div>
        <button class="vms-btn vms-btn-ghost" @click="emit('close')">X</button>
      </div>

      <div class="vms-flex-col" style="gap: 1rem; margin-top: 1rem;">
        <div style="background: rgba(255, 94, 58, 0.1); border: 1px solid rgba(255, 94, 58, 0.35); border-radius: 6px; padding: 0.85rem;">
          <p style="color: #fff; font-size: 13px; line-height: 1.4; margin: 0;">
            A pasta <strong style="color: #ff5e3a;">"{{ folderName }}"</strong> contem <strong style="color: #ff5e3a;">{{ itemCount }} {{ itemType || 'fluxo(s)' }}</strong> internamente.
          </p>
          <p class="vms-text-dim" style="font-size: 11px; margin-top: 0.5rem; margin-bottom: 0;">
            Tem certeza da acao? Ao confirmar, a pasta sera excluida e todos os itens internos serao preservados e movidos para a raiz.
          </p>
        </div>

        <div class="vms-flex-row" style="justify-content: flex-end; gap: 0.75rem; margin-top: 0.5rem;">
          <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="emit('confirm')">SALVAR</button>
        </div>
      </div>
    </div>
  </div>
</template>
