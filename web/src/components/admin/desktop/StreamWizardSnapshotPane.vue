<script setup lang="ts">
defineProps<{
  isTesting: boolean
  hasSnapshot: boolean
  snapshotUrl?: string
  authRequired?: boolean
  testError?: string | null
  detectedCodec: string
  detectedResolution: string
  detectedFps: number
  latencyMs: number
}>()
</script>

<template>
  <div class="vms-flex-col" style="flex: 1 1 360px; min-width: 280px; max-width: 100%; min-height: 380px; background: #07090e; border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.85rem; justify-content: space-between; box-sizing: border-box;">
    <!-- Header -->
    <div class="vms-flex-between" style="border-bottom: 1px solid rgba(255, 255, 255, 0.08); padding-bottom: 0.4rem; gap: 0.5rem; align-items: center; flex-wrap: wrap;">
      <span class="vms-text-mono" style="color: #ffffff; font-size: 11px; font-weight: 700; white-space: nowrap; letter-spacing: 0.5px;">
        SNAPSHOT DO FLUXO
      </span>
      <span v-if="isTesting" class="vms-badge" style="background: rgba(0, 240, 255, 0.15); color: var(--vms-neu-accent-cyan); border: 1px solid rgba(0, 240, 255, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        TESTANDO CONEXÃO...
      </span>
      <span v-else-if="authRequired" class="vms-badge" style="background: rgba(255, 94, 58, 0.15); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        AUTENTICAÇÃO NECESSÁRIA (401)
      </span>
      <span v-else-if="hasSnapshot" class="vms-badge" style="background: rgba(0, 255, 157, 0.15); color: var(--vms-neu-accent-green); border: 1px solid rgba(0, 255, 157, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        ONLINE
      </span>
      <span v-else-if="testError" class="vms-badge" style="background: rgba(255, 0, 60, 0.15); color: #ff003c; border: 1px solid rgba(255, 0, 60, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        OFFLINE
      </span>
      <span v-else class="vms-badge" style="background: rgba(255, 255, 255, 0.06); color: var(--vms-text-dim); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        AGUARDANDO TESTE
      </span>
    </div>

    <!-- Snapshot Screen / Canvas Frame -->
    <div style="width: 100%; height: 245px; background: #000000; border: 1px solid var(--vms-border); border-radius: 6px; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; padding: 1rem; box-sizing: border-box;">
      <div v-if="isTesting" class="vms-flex-col" style="align-items: center; gap: 0.75rem;">
        <div class="loader"></div>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan);">VERIFICANDO FLUXO...</span>
      </div>

      <img v-else-if="snapshotUrl" :src="snapshotUrl" alt="Snapshot" style="width: 100%; height: 100%; object-fit: contain;" />

      <div v-else-if="testError" class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.45rem; text-align: center; max-width: 90%;">
        <svg width="48" height="48" viewBox="0 0 576 512" fill="#ff003c"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        <span class="vms-text-mono vms-text-xs" style="color: #ff003c; font-weight: 700;">FALHA DE CONEXÃO</span>
        <span class="vms-text-mono vms-text-2xs" style="color: #ff8899; word-break: break-word;">{{ testError }}</span>
      </div>

      <div v-else class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.5rem; text-align: center;">
        <svg width="48" height="48" viewBox="0 0 576 512" :fill="hasSnapshot ? '#00ff9d' : '#ff5e3a'"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        <span v-if="hasSnapshot" class="vms-text-mono vms-text-xs" style="color: #00ff9d; font-weight: 600;">SINAL ONLINE DETECTADO</span>
        <span v-else class="vms-text-dim vms-text-2xs">CLIQUE EM "TESTAR" PARA VALIDAR A CONEXÃO</span>
      </div>
    </div>

    <!-- Auto-Detected Telemetry Box -->
    <div class="vms-flex-col" style="gap: 4px; background: rgba(255, 255, 255, 0.03); border: 1px solid var(--vms-border); padding: 0.5rem 0.75rem; border-radius: 4px; font-family: var(--vms-font-jetbrains); font-size: 11px;">
      <div class="vms-flex-between">
        <span style="color: var(--vms-text-dim);">CODEC NAL:</span>
        <strong style="color: #fff;">{{ hasSnapshot ? detectedCodec : '--' }}</strong>
      </div>
      <div class="vms-flex-between">
        <span style="color: var(--vms-text-dim);">RESOLUÇÃO:</span>
        <strong style="color: #fff;">{{ hasSnapshot ? (detectedFps > 0 ? `${detectedResolution} (${detectedFps} FPS)` : detectedResolution) : '--' }}</strong>
      </div>
    </div>
  </div>
</template>

