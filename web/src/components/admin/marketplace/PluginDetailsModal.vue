<script setup lang="ts">
import type { PluginManifest } from '../../../types/marketplace'

defineProps<{
  isOpen: boolean
  plugin: PluginManifest | null
}>()

const emit = defineEmits<{ (e: 'close'): void }>()
</script>

<template>
  <div v-if="isOpen && plugin" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content" style="max-width: 580px; width: 100%;">
      <div class="vms-modal-header vms-flex-between">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff;">{{ plugin.name }}</h4>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">MANIFESTO TÉCNICO // plugin.json</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body vms-flex-col" style="gap: 1rem; padding: 1.25rem 0;">
        <div class="vms-flex-col" style="gap: 0.35rem;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">DESCRIÇÃO DO PACOTE</span>
          <p class="vms-text-xs vms-text-dim" style="margin: 0; line-height: 1.4;">{{ plugin.description }}</p>
        </div>

        <div class="vms-flex-row" style="gap: 1rem; flex-wrap: wrap;">
          <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 140px;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">VERSÃO:</span>
            <span class="vms-text-xs vms-font-semibold">v{{ plugin.version }}</span>
          </div>
          <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 140px;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">AUTOR / FORNECEDOR:</span>
            <span class="vms-text-xs vms-font-semibold">{{ plugin.author }}</span>
          </div>
          <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 140px;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">RUNTIME:</span>
            <span class="vms-badge vms-badge-secondary" style="font-size: 10px; width: fit-content;">{{ plugin.runtime }}</span>
          </div>
          <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 140px;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">ACELERAÇÃO:</span>
            <span class="vms-badge vms-badge-orange" style="font-size: 10px; width: fit-content;">{{ plugin.hardware_req }}</span>
          </div>
        </div>

        <div class="vms-flex-col" style="gap: 0.35rem;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">PERMISSÕES REQUERIDAS (SANDBOX IPC)</span>
          <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
            <span v-for="perm in plugin.permissions" :key="perm" class="vms-badge vms-badge-secondary" style="font-size: 10px;">
              {{ perm }}
            </span>
          </div>
        </div>

        <div class="vms-flex-col" style="gap: 0.35rem;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">CONFIGURAÇÃO PADRÃO (JSON)</span>
          <pre class="vms-text-mono vms-text-2xs" style="background: #080a0f; padding: 0.75rem; border-radius: 4px; border: 1px solid var(--vms-border); overflow-x: auto; color: #a5f3fc; margin: 0;">{{ JSON.stringify(plugin.default_config, null, 2) }}</pre>
        </div>
      </div>

      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.85rem;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ASSINATURA ECDSA // VALIDADE: OFICIAL</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">FECHAR</button>
      </div>
    </div>
  </div>
</template>
