<script setup lang="ts">
import { ref } from 'vue'
import type { StreamItem } from '../../../types/streamTree'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'

const props = defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{ (e: 'close'): void }>()
const refreshKey = ref(Date.now()), isRefreshing = ref(false)

const refreshSnapshot = async () => {
  if (isRefreshing.value) return
  isRefreshing.value = true
  try {
    await fetch(getCameraSnapshotUrl(props.stream.id, true))
    refreshKey.value = Date.now()
  } catch {}
  finally { setTimeout(() => { isRefreshing.value = false }, 500) }
}
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
    <div style="width: 100%; aspect-ratio: 16 / 9; max-height: 480px; background: #000000; border: 1px solid var(--vms-border); border-radius: 8px; display: flex; align-items: center; justify-content: center; overflow: hidden; position: relative;">
      <img v-if="stream.snapshotUrl" :src="`${getCameraSnapshotUrl(stream.id)}&k=${refreshKey}`" alt="Snapshot" style="width: 100%; height: 100%; object-fit: contain; display: block;" />
      <svg v-else width="56" height="56" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
      <button style="position: absolute; top: 8px; right: 8px; width: 30px; height: 30px; border-radius: 6px; background: rgba(14, 17, 23, 0.85); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all 0.2s;" title="Capturar Novo Snapshot do Fluxo" @click.stop="refreshSnapshot">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round" :style="{ transform: isRefreshing ? 'rotate(360deg)' : 'none', transition: 'transform 0.5s ease' }">
          <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
        </svg>
      </button>
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
        <span class="vms-text-mono vms-text-xs vms-font-semibold" :style="{ color: stream.recordMode && stream.recordMode !== 'disabled' ? 'var(--vms-neu-accent-green)' : 'var(--vms-text-dim)' }">
          {{ stream.recordMode && stream.recordMode !== 'disabled' ? 'SIM' : 'NAO' }}
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">TAXA DE FLUXO / BITRATE</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.bitrate }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ANALITICOS</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.analyticsCount || 0 }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ALARMES & SENSORES</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ stream.alarmsCount || 0 }}</span>
      </div>
    </div>
  </div>
</template>
