<script setup lang="ts">
import type { FlowNodeItem } from '../../../types/workflowTree'

defineProps<{ node: FlowNodeItem; isLocked?: boolean }>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 8px;">
    <!-- 1. DETECCAO IA -->
    <template v-if="node.kind === 'AI_DETECT'">
      <div class="vms-form-group">
        <label class="vms-label">CLASSE DE OBJETO ALVO</label>
        <select v-model="node.config.targetClass" class="vms-auth-input" :disabled="isLocked">
          <option value="PESSOA">[PESSOA]</option>
          <option value="VEICULO">[VEICULO]</option>
          <option value="PLACA_LPR">[PLACA LPR]</option>
          <option value="EPI_CAPACETE">[EPI CAPACETE]</option>
          <option value="FOGO_FUMACA">[FOGO E FUMACA]</option>
        </select>
      </div>
    </template>

    <!-- 2. CONDICOES LOGICAS E FLUXO EM Y -->
    <template v-else-if="node.kind === 'BRANCH_OR'">
      <div class="vms-form-group">
        <label class="vms-label">LOGICA DE JUNCAO</label>
        <select v-model="node.config.logic" class="vms-auth-input" :disabled="isLocked">
          <option value="OR">[OU // QUALQUER ENTRADA ATIVA]</option>
        </select>
      </div>
      <div style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 8px;">
        <span class="vms-text-dim vms-text-2xs">// Interliga 2 alertas ou deteccoes para 1 saida comum.</span>
      </div>
    </template>

    <template v-else-if="node.kind === 'LOGIC_AND'">
      <div class="vms-form-group">
        <label class="vms-label">JANELA TEMPORAL ({{ node.config.windowSeconds || 5 }}s)</label>
        <input v-model.number="node.config.windowSeconds" type="range" min="1" max="30" class="vms-auth-input" :disabled="isLocked" />
      </div>
      <div style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 8px;">
        <span class="vms-text-dim vms-text-2xs">// Dispara apenas se os dois eventos ocorrerem simultaneamente.</span>
      </div>
    </template>

    <template v-else-if="node.kind === 'SPLIT_Y'">
      <div class="vms-form-group">
        <label class="vms-label">MODO DE BIFURCACAO</label>
        <select v-model="node.config.mode" class="vms-auth-input" :disabled="isLocked">
          <option value="PARALELO">[PARALELO // DISPARAR AMBAS]</option>
        </select>
      </div>
      <div style="background: #07080c; border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 8px;">
        <span class="vms-text-dim vms-text-2xs">// Bifurca 1 evento de entrada para 2 acoes simultaneas.</span>
      </div>
    </template>

    <template v-else-if="node.kind === 'SCHEDULE'">
      <div class="vms-form-group">
        <label class="vms-label">JANELA DE HORARIO</label>
        <select v-model="node.config.timeWindow" class="vms-auth-input" :disabled="isLocked">
          <option value="NOTURNO // 22H AS 06H">[NOTURNO // 22H AS 06H]</option>
          <option value="COMERCIAL // 08H AS 18H">[COMERCIAL // 08H AS 18H]</option>
        </select>
      </div>
    </template>

    <template v-else-if="node.kind === 'ANTI_SPAM'">
      <div class="vms-form-group">
        <label class="vms-label">COOLDOWN ({{ node.config.cooldownSeconds || 30 }}s)</label>
        <input v-model.number="node.config.cooldownSeconds" type="range" min="5" max="300" step="5" class="vms-auth-input" :disabled="isLocked" />
      </div>
    </template>

    <!-- 3. ACOES -->
    <template v-else-if="node.kind === 'TELEGRAM'">
      <div class="vms-form-group">
        <label class="vms-label">CHAT ID / CANAL</label>
        <input v-model="node.config.chatId" class="vms-auth-input" placeholder="-1001928374..." :disabled="isLocked" />
      </div>
    </template>

    <template v-else-if="node.kind === 'WEBHOOK'">
      <div class="vms-form-group">
        <label class="vms-label">ENDPOINT URL</label>
        <input v-model="node.config.endpointUrl" class="vms-auth-input" placeholder="https://siem.corp.com/alerts" :disabled="isLocked" />
      </div>
    </template>
  </div>
</template>
