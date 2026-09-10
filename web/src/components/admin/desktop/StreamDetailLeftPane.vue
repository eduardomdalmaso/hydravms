<script setup lang="ts">
import type { StreamItem } from '../../../types/streamTree'

defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{ (e: 'close'): void }>()
</script>

<template>
  <div class="vms-split-pane">
    <!-- Header -->
    <div class="vms-split-header">
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div style="width: 36px; height: 36px; border-radius: 8px; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center;">
          <svg width="18" height="18" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        </div>
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: #fff; font-size: 13px;">{{ stream.id.toUpperCase() }} // {{ stream.name }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ stream.url }}</span>
        </div>
      </div>
    </div>

    <!-- Snapshot / Preview Box -->
    <div style="width: 100%; aspect-ratio: 16 / 9; max-height: 480px; background: #07080c; border: 1px solid var(--vms-border); border-radius: 8px; display: flex; align-items: center; justify-content: center; overflow: hidden; position: relative;">
      <img v-if="stream.snapshotUrl" :src="stream.snapshotUrl" alt="Snapshot" style="width: 100%; height: 100%; object-fit: cover;" />
      <svg v-else width="56" height="56" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
      <div v-if="stream.snapshotUrl && stream.fps" style="position: absolute; top: 6px; left: 6px; font-family: var(--vms-font-jetbrains); font-size: 10px; color: var(--vms-neu-accent-green); background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;">
        REC // {{ stream.fps }} FPS
      </div>
    </div>


    <!-- Telemetry Information Grid -->
    <div class="vms-telemetry-grid">
      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">CODEC</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.codec }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">COMPRESSAO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.codec === 'H.265' ? 'HEVC' : 'AVC' }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">RESOLUCAO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.resolution }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">FPS</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.fps }} FPS</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">GRAVANDO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">
          {{ stream.recordMode !== 'disabled' ? 'SIM (HABILITADO)' : 'NAO (DESLIGADO)' }}
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">TAXA DE FLUXO / BITRATE</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.bitrate }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ANALITICOS</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">2</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ALARMES & SENSORES</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">1</span>
      </div>
    </div>
  </div>
</template>
