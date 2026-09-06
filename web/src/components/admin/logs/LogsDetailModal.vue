<script setup lang="ts">
import type { LogEntry } from '../../../types/systemLogs'

defineProps<{ log: LogEntry }>()
const emit = defineEmits<{ (e: 'close'): void }>()
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 520px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">REGISTRO FORENSE // DETALHE DO LOG</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">✕</button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 12px;">
        <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 8px;">
          <div>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">IDENTIFICADOR</span>
            <div class="vms-text-mono vms-text-xs" style="color: #ff5e3a; font-weight: 700;">{{ log.id }}</div>
          </div>
          <div style="text-align: right;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">TIMESTAMP UTC</span>
            <div class="vms-text-mono vms-text-xs" style="color: #ffffff;">{{ log.timestamp }}</div>
          </div>
        </div>

        <div class="vms-grid-container" style="grid-template-columns: 1fr 1fr; gap: 8px;">
          <div>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">CATEGORIA // NIVEL</span>
            <div class="vms-text-xs" style="color: #ffffff; font-weight: 600;">{{ log.category }} // [{{ log.level }}]</div>
          </div>
          <div>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">TENANT / ESCOPO</span>
            <div class="vms-text-xs" style="color: #ffffff; font-weight: 600;">{{ log.tenantName }}</div>
          </div>
          <div>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">ATOR / ORIGEM</span>
            <div class="vms-text-xs" style="color: #ffffff; font-weight: 600;">{{ log.actor }}</div>
          </div>
          <div>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">IP DE ORIGEM</span>
            <div class="vms-text-mono vms-text-xs" style="color: #ffffff;">{{ log.ipAddress }}</div>
          </div>
        </div>

        <div>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">ACAO EXECUTADA // ALVO</span>
          <div class="vms-text-mono vms-text-xs" style="color: var(--vms-neu-accent-orange); font-weight: 700; margin-top: 2px;">
            {{ log.action }} // {{ log.target }}
          </div>
        </div>

        <div>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">DETALHAMENTO OPERACIONAL</span>
          <div class="vms-card" style="padding: 10px; background: #07080c; border: 1px solid var(--vms-border); margin-top: 4px; font-size: 0.8rem; color: #ffffff; line-height: 1.4;">
            {{ log.details }}
          </div>
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-primary" @click="emit('close')">FECHAR</button>
      </div>
    </div>
  </div>
</template>
