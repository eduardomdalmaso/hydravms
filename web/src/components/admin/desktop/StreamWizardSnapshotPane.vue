<script setup lang="ts">
defineProps<{
  isTesting: boolean
  hasSnapshot: boolean
  snapshotUrl?: string
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
      <span v-if="hasSnapshot" class="vms-badge" style="background: rgba(0, 255, 157, 0.15); color: var(--vms-neu-accent-green); border: 1px solid rgba(0, 255, 157, 0.4); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        FLUXO ESTABELECIDO
      </span>
      <span v-else-if="isTesting" class="vms-badge" style="background: rgba(0, 240, 255, 0.15); color: var(--vms-neu-accent-cyan); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        BUSCANDO SNAPSHOT...
      </span>
      <span v-else class="vms-badge" style="background: rgba(255, 255, 255, 0.06); color: var(--vms-text-dim); font-size: 9px; white-space: nowrap; flex-shrink: 0;">
        AGUARDANDO PROBE
      </span>
    </div>

    <!-- Snapshot Screen / Canvas with enlarged 16:9 frame -->
    <div style="width: 100%; height: 245px; background: #020305; border: 1px dashed rgba(255, 255, 255, 0.15); border-radius: 3px; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden;">
      <img v-if="hasSnapshot" :src="snapshotUrl || 'https://images.unsplash.com/photo-1557597774-9d273605dfa9?w=800&q=80'" alt="Snapshot" style="width: 100%; height: 100%; object-fit: cover;" />
      
      <div v-if="isTesting" class="vms-flex-col" style="align-items: center; gap: 0.5rem; z-index: 2;">
        <div style="width: 28px; height: 28px; border: 2px solid var(--vms-neu-accent-cyan); border-top-color: transparent; border-radius: 50%; animation: spin 0.8s linear infinite;"></div>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan); font-weight: 600;">NEGOCIANDO SDP & NAL UNITS...</span>
      </div>

      <div v-else-if="!hasSnapshot" class="vms-flex-col" style="align-items: center; gap: 0.4rem; text-align: center; padding: 0.75rem;">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#555" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">Clique em "Buscar Snapshot" para testar o sinal</span>
      </div>

      <div v-if="hasSnapshot" style="position: absolute; top: 6px; left: 6px; font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-green); background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;">
        REC // 30 FPS
      </div>
      <div v-if="hasSnapshot" style="position: absolute; bottom: 6px; right: 6px; font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-cyan); background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;">
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
        <strong style="color: #fff;">{{ hasSnapshot ? `${detectedResolution} (${detectedFps} FPS)` : '--' }}</strong>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes spin { to { transform: rotate(360deg); } }
</style>
