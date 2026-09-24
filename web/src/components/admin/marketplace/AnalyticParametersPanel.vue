<script setup lang="ts">
import type { AnalyticInstance } from '../../../types/marketplace'

defineProps<{
  instance: AnalyticInstance
}>()
</script>

<template>
  <div class="vms-card vms-flex-col" style="padding: 0.85rem 1rem; gap: 0.75rem; background: #0b0e14; border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 8px;">
    <div class="vms-flex-between" style="border-bottom: 1px solid rgba(255, 255, 255, 0.06); padding-bottom: 0.4rem; align-items: center;">
      <span class="vms-text-mono vms-text-2xs vms-font-semibold" style="color: #8b94a0; letter-spacing: 0.5px;">// PARÂMETROS DE EXECUÇÃO (SOMENTE LEITURA)</span>
      <span class="vms-badge vms-badge-secondary" style="font-size: 8.5px; color: #64748b; background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.08);">[FIXADO]</span>
    </div>

    <!-- Read-only HUD Parameters Grid in Muted Grey Tones -->
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(170px, 1fr)); gap: 0.6rem;">
      <!-- Camera -->
      <div class="vms-param-box">
        <span class="vms-param-label">CÂMERA:</span>
        <span class="vms-param-value">{{ instance.camera_name }}</span>
      </div>

      <!-- FPS -->
      <div class="vms-param-box">
        <span class="vms-param-label">TAXA DE INFERÊNCIA:</span>
        <span class="vms-param-value">{{ instance.fps_rate || 15 }} FPS</span>
      </div>

      <!-- Hardware -->
      <div class="vms-param-box">
        <span class="vms-param-label">HARDWARE TARGET:</span>
        <span class="vms-param-value">{{ instance.hardware_target === 'rtx_5090_cuda' ? 'GPU RTX 5090 (CUDA)' : 'CPU (ZERO-COPY SHM)' }}</span>
      </div>

      <!-- Stream -->
      <div class="vms-param-box">
        <span class="vms-param-label">FLUXO DE VÍDEO:</span>
        <span class="vms-param-value">{{ instance.stream_type === 'sub_stream' ? 'Sub-Stream (480p)' : 'Main Stream (1080p/4K)' }}</span>
      </div>

      <!-- Motion Gated -->
      <div class="vms-param-box">
        <span class="vms-param-label">FILTRO DE MOVIMENTO:</span>
        <span class="vms-param-value">{{ instance.motion_gated !== false ? 'Motion-Gated [ATIVO]' : 'Contínuo [DESATIVADO]' }}</span>
      </div>

      <!-- Zones Count -->
      <div class="vms-param-box">
        <span class="vms-param-label">ZONAS ATIVAS:</span>
        <span class="vms-param-value">{{ (instance.zones || []).length > 0 ? `${instance.zones?.length} Zona(s)` : 'Frame Completo' }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-param-box {
  background: #0e1118;
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 5px;
  padding: 0.45rem 0.65rem;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.vms-param-label {
  font-family: var(--vms-font-mono);
  font-size: 8.5px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.4px;
}
.vms-param-value {
  font-family: var(--vms-font-mono);
  font-size: 11px;
  color: #cbd5e1;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
