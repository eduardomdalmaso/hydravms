<script setup lang="ts">
import { ref } from 'vue'
import type { StreamItem } from '../../../types/streamTree'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'

const props = defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'test', id: string): void
  (e: 'delete', id: string): void
}>()

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
  <div class="vms-desktop-inspector">
    <!-- Header -->
    <div class="vms-flex-between" style="padding: 0.75rem 1rem; border-bottom: 1px solid var(--vms-border); background: #16191f;">
      <div class="vms-flex-row" style="gap: 0.5rem;">
        <span class="vms-status-led" :class="stream.status" :title="stream.status"></span>
        <span class="vms-font-bold" style="color: #fff; font-size: 13px;">{{ stream.name }}</span>
      </div>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar Painel" @click="emit('close')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <!-- Inspector Content -->
    <div class="vms-flex-col" style="padding: 1rem; gap: 1rem;">
      <!-- Mini Preview -->
      <div style="height: 140px; background: #000000; border: 1px solid var(--vms-border); border-radius: 6px; display: flex; flex-direction: column; justify-content: center; align-items: center; gap: 0.5rem; position: relative; overflow: hidden;">
        <img v-if="stream.snapshotUrl" :src="`${getCameraSnapshotUrl(stream.id)}&k=${refreshKey}`" alt="Snapshot" style="width: 100%; height: 100%; object-fit: contain; display: block;" />
        <template v-else>
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="rgba(255, 255, 255, 0.3)" stroke-width="1.5"><path d="m22 8-6 4 6 4V8Z"/><rect width="14" height="12" x="2" y="6" rx="2" ry="2"/></svg>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-regular);">// {{ stream.resolution }} @ {{ stream.fps }} FPS</span>
        </template>
        <span class="vms-badge vms-badge-orange" style="position: absolute; top: 6px; left: 6px; font-size: 8px;">{{ stream.codec }}</span>
        <button style="position: absolute; top: 6px; right: 6px; width: 24px; height: 24px; border-radius: 4px; background: rgba(14, 17, 23, 0.85); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center; cursor: pointer;" title="Atualizar Snapshot" @click.stop="refreshSnapshot">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round" :style="{ transform: isRefreshing ? 'rotate(360deg)' : 'none', transition: 'transform 0.5s ease' }">
            <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
          </svg>
        </button>
      </div>


      <!-- Specs Card -->
      <div class="vms-flex-col" style="gap: 0.5rem; background: rgba(255, 255, 255, 0.02); padding: 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
        <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">PROTOCOLO</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-main);">{{ stream.protocol }}</span></div>
        <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">ENDEREÇO IP</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-main);">{{ stream.ip }}:{{ stream.port }}</span></div>
        <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">BITRATE ESTIMADO</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-regular);">{{ stream.bitrate }}</span></div>
        <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">MODO GRAVAÇÃO</span><span class="vms-badge vms-badge-orange" style="font-size: 9px;">{{ stream.recordMode.toUpperCase() }}</span></div>
        <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">CONTROLE PTZ</span><span class="vms-text-mono vms-text-xs" :style="{ color: stream.has_ptz ? 'var(--vms-status-online)' : 'var(--vms-text-dim)' }">{{ stream.has_ptz ? 'HABILITADO' : 'DESATIVADO' }}</span></div>
      </div>

      <!-- Actions -->
      <div class="vms-flex-col" style="gap: 0.5rem; margin-top: auto;">
        <button class="vms-btn vms-btn-primary" style="font-size: 11px;" @click="emit('test', stream.id)">[TESTAR HANDSHAKE RTSP]</button>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="color: #ff5e3a; font-size: 11px; display: inline-flex; align-items: center; justify-content: center; gap: 6px;" @click="emit('delete', stream.id)">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
          <span>[REMOVER FLUXO]</span>
        </button>
      </div>
    </div>
  </div>
</template>
