<script setup lang="ts">
defineProps<{
  isTesting: boolean
  hasSnapshot: boolean
  snapshotUrl?: string
  authRequired?: boolean
  detectedCodec: string
  detectedResolution: string
  detectedFps: number
  latencyMs: number
}>()
</script>

<template>
  <div class="vms-flex-col" style="width: 490px; flex-shrink: 0; height: 380px; background: #07090e; border: 1px solid rgba(0, 240, 255, 0.25); border-radius: 4px; padding: 0.85rem; justify-content: space-between; box-sizing: border-box;">
    <!-- Header with guaranteed no-wrap title -->
    <div class="vms-flex-between" style="border-bottom: 1px solid rgba(255, 255, 255, 0.08); padding-bottom: 0.4rem; gap: 0.75rem; align-items: center;">
      <span class="vms-text-mono" style="color: var(--vms-neu-accent-cyan); font-size: 11px; font-weight: 700; white-space: nowrap; letter-spacing: 0.5px;">
        // PROBE & SNAPSHOT LIVE
      </span>
      <span v-if="authRequired" class="vms-badge" style="background: rgba(255, 94, 58, 0.15); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        AUTENTICAÇÃO NECESSÁRIA (401)
      </span>
      <span v-else-if="hasSnapshot && snapshotUrl" class="vms-badge" style="background: rgba(0, 255, 157, 0.15); color: var(--vms-neu-accent-green); border: 1px solid rgba(0, 255, 157, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        FLUXO ESTABELECIDO
      </span>
      <span v-else class="vms-badge" style="background: rgba(255, 255, 255, 0.06); color: var(--vms-text-dim); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        AGUARDANDO TESTE
      </span>
    </div>

    <!-- Snapshot Screen / Canvas with enlarged 16:9 frame -->
    <div style="width: 100%; height: 245px; background: #020305; border: 1px dashed rgba(255, 255, 255, 0.15); border-radius: 3px; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden;">
      <div v-if="isTesting" class="loader"></div>

      <img v-else-if="hasSnapshot && snapshotUrl" :src="snapshotUrl" alt="Snapshot" style="width: 100%; height: 100%; object-fit: contain;" />
      
      <div v-else class="vms-flex-col" style="align-items: center; justify-content: center;">
        <svg width="52" height="52" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
      </div>

      <div v-if="hasSnapshot && snapshotUrl && detectedFps > 0 && !isTesting" style="position: absolute; top: 6px; left: 6px; font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-green); background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;">
        LIVE // {{ detectedFps }} FPS
      </div>
      <div v-if="latencyMs > 0 && !isTesting" style="position: absolute; bottom: 6px; right: 6px; font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-cyan); background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;">
        LATENCIA: {{ latencyMs }}ms
      </div>
    </div>


    <!-- Auto-Detected Telemetry Box -->
    <div class="vms-flex-col" style="gap: 3px; background: rgba(0, 240, 255, 0.05); border: 1px solid rgba(0, 240, 255, 0.15); padding: 0.45rem 0.65rem; border-radius: 3px; font-family: var(--vms-font-jetbrains); font-size: 10px;">
      <div class="vms-flex-between">
        <span style="color: var(--vms-text-dim);">CODEC NAL:</span>
        <strong style="color: #fff;">{{ hasSnapshot ? detectedCodec : '--' }}</strong>
      </div>
      <div class="vms-flex-between">
        <span style="color: var(--vms-text-dim);">RESOLUCAO & TAXA:</span>
        <strong style="color: #fff;">{{ hasSnapshot ? (detectedFps > 0 ? `${detectedResolution} (${detectedFps} FPS)` : detectedResolution) : '--' }}</strong>
      </div>
    </div>
  </div>
</template>

