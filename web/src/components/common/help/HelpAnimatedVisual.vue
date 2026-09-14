<script setup lang="ts">
defineProps<{ type: string }>()
</script>

<template>
  <div class="vms-help-visual-container">
    <!-- 1. STREAMS SCAN RADAR ANIMATION -->
    <div v-if="type === 'streams_scan'" class="vms-radar-box">
      <div class="vms-radar-sweep" />
      <div class="vms-radar-dot dot-1" />
      <div class="vms-radar-dot dot-2" />
      <div class="vms-radar-dot dot-3" />
      <div class="vms-radar-hud">
        <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">[ONVIF SCAN] 3 CÂMERAS DETECTADAS</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">TCP 554 // RTSP INTERLEAVED // 30 FPS</span>
      </div>
    </div>

    <!-- 2. GRID SLOTS MULTI-VIEW ANIMATION -->
    <div v-else-if="type === 'grid_slots'" class="vms-grid-demo">
      <div class="vms-grid-cell active"><span class="vms-cell-label">CAM 01 // PORTARIA</span><div class="vms-live-pulse" /></div>
      <div class="vms-grid-cell active"><span class="vms-cell-label">CAM 02 // ESTACIONAMENTO</span><div class="vms-live-pulse" /></div>
      <div class="vms-grid-cell active"><span class="vms-cell-label">CAM 03 // DOCA</span><div class="vms-live-pulse" /></div>
      <div class="vms-grid-cell active"><span class="vms-cell-label">CAM 04 // PERÍMETRO</span><div class="vms-live-pulse" /></div>
    </div>

    <!-- 3. RBAC TREE ACCESS CONTROL ANIMATION -->
    <div v-else-if="type === 'rbac_tree'" class="vms-rbac-demo">
      <div class="vms-node master">ADMIN MASTER // ACESSO GLOBAL</div>
      <div class="vms-beam" />
      <div class="vms-flex-row" style="gap: 0.5rem; width: 100%; justify-content: center;">
        <div class="vms-node tenant">CLIENTE ALPHA // TENANT</div>
        <div class="vms-node tenant">CLIENTE BETA // TENANT</div>
      </div>
    </div>

    <!-- 4. BRAND PALETTE ANIMATION -->
    <div v-else-if="type === 'brand_palette'" class="vms-brand-demo">
      <div class="vms-brand-tab">
        <div class="vms-brand-dot" />
        <span class="vms-text-mono vms-text-2xs" style="color: #fff;">MINHA_EMPRESA_VMS</span>
      </div>
      <div class="vms-brand-header">
        <span class="vms-text-xs vms-font-bold" style="color: var(--vms-neu-accent-orange);">LOGO PERSONALIZADO // WHITE-LABEL</span>
      </div>
    </div>

    <!-- 5. FLOW / WORKFLOW / CLUSTER FALLBACK -->
    <div v-else class="vms-flow-demo">
      <div class="vms-flow-card">EVENTO IA // INTRUSÃO</div>
      <div class="vms-arrow">➔</div>
      <div class="vms-flow-card action">BOT TELEGRAM (FOTO)</div>
    </div>
  </div>
</template>

<style scoped>
.vms-help-visual-container {
  width: 100%; height: 130px; background: #07090e; border: 1px solid var(--vms-border);
  border-radius: 8px; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden;
}
.vms-radar-box { position: relative; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; }
.vms-radar-sweep {
  position: absolute; width: 160px; height: 160px; border-radius: 50%; border: 1px solid rgba(0, 255, 157, 0.2);
  background: conic-gradient(from 0deg, rgba(0, 255, 157, 0.35) 0deg, transparent 90deg);
  animation: radar-spin 3s linear infinite;
}
.vms-radar-dot { position: absolute; width: 8px; height: 8px; border-radius: 50%; background: #00ff9d; box-shadow: 0 0 8px #00ff9d; }
.dot-1 { top: 30px; left: 60px; } .dot-2 { bottom: 35px; right: 70px; } .dot-3 { top: 40px; right: 100px; }
.vms-radar-hud { position: absolute; bottom: 8px; left: 12px; display: flex; flex-direction: column; gap: 2px; }
@keyframes radar-spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.vms-grid-demo { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; width: 85%; height: 85%; }
.vms-grid-cell { background: #121620; border: 1px solid rgba(255,94,58,0.3); border-radius: 4px; position: relative; display: flex; align-items: flex-end; padding: 4px; }
.vms-cell-label { font-size: 8px; font-family: var(--vms-font-mono); color: #fff; }
.vms-live-pulse { position: absolute; top: 4px; right: 4px; width: 6px; height: 6px; border-radius: 50%; background: #00ff9d; box-shadow: 0 0 6px #00ff9d; animation: pulse 1.5s infinite; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.3; } }

.vms-rbac-demo { display: flex; flex-direction: column; align-items: center; gap: 8px; width: 90%; }
.vms-node { font-size: 9px; font-family: var(--vms-font-mono); padding: 4px 8px; border-radius: 4px; border: 1px solid var(--vms-border); background: #141824; color: #fff; }
.vms-node.master { border-color: var(--vms-neu-accent-orange); color: var(--vms-neu-accent-orange); }
.vms-beam { width: 2px; height: 12px; background: var(--vms-neu-accent-orange); }

.vms-brand-demo { display: flex; flex-direction: column; gap: 8px; width: 85%; }
.vms-brand-tab { background: #181d28; padding: 4px 10px; border-radius: 4px 4px 0 0; width: fit-content; display: flex; align-items: center; gap: 6px; border-top: 2px solid var(--vms-neu-accent-orange); }
.vms-brand-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--vms-neu-accent-orange); }
.vms-brand-header { background: #11141c; padding: 8px; border-radius: 4px; border: 1px dashed var(--vms-border); text-align: center; }

.vms-flow-demo { display: flex; align-items: center; justify-content: center; gap: 12px; width: 85%; }
.vms-flow-card { background: #131722; border: 1px solid var(--vms-border); border-radius: 4px; padding: 6px 12px; font-size: 9px; font-family: var(--vms-font-mono); color: #fff; }
.vms-flow-card.action { border-color: #00ff9d; color: #00ff9d; }
.vms-arrow { color: var(--vms-neu-accent-orange); font-size: 14px; animation: slide-arrow 1.2s infinite; }
@keyframes slide-arrow { 0%, 100% { transform: translateX(0); } 50% { transform: translateX(4px); } }
</style>
