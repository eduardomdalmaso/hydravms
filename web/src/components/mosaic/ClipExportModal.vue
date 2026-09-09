<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { CameraStreamInfo } from '../../types/mosaic'

const props = defineProps<{ isOpen: boolean; camera: CameraStreamInfo; startMs: number; endMs: number }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'export', config: any): void }>()

const toIso = (ms: number) => {
  const d = new Date(ms), p = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

const startTimeStr = ref(toIso(props.startMs || Date.now() - 300000)), endTimeStr = ref(toIso(props.endMs || Date.now()))
const format = ref<'mp4' | 'mkv' | 'avi'>('mp4'), exportMode = ref<'stream-copy' | 'nvenc'>('stream-copy')
const addWatermark = ref(true), addSha256 = ref(true), isExporting = ref(false)

watch(() => props.isOpen, (open) => {
  if (open) { startTimeStr.value = toIso(props.startMs || Date.now() - 300000); endTimeStr.value = toIso(props.endMs || Date.now()); isExporting.value = false }
})

const durationSec = computed(() => Math.max(0, Math.round((new Date(endTimeStr.value).getTime() - new Date(startTimeStr.value).getTime()) / 1000)))
const durationFormatted = computed(() => {
  const s = durationSec.value; if (s <= 0) return '0s (0s)'
  const y = Math.floor(s / 31536000), mo = Math.floor((s % 31536000) / 2592000)
  const d = Math.floor((s % 2592000) / 86400), h = Math.floor((s % 86400) / 3600)
  const m = Math.floor((s % 3600) / 60), sec = s % 60
  const p: string[] = []
  if (y) p.push(`${y}a`); if (mo) p.push(`${mo}mês${mo > 1 ? 'es' : ''}`); if (d) p.push(`${d}d`)
  if (h) p.push(`${h}h`); if (m || h || d) p.push(`${m}m`); p.push(`${sec}s`)
  return `${p.join(' ')} (${s}s)`
})
const estimatedSize = computed(() => {
  const mb = (durationSec.value * 4) / 8
  if (mb >= 1048576) return `~${(mb / 1048576).toFixed(2)} TB`
  if (mb >= 1024) return `~${(mb / 1024).toFixed(2)} GB`
  return mb >= 1 ? `~${mb.toFixed(1)} MB` : `~${(mb * 1024).toFixed(0)} KB`
})

const handleDownload = () => {
  isExporting.value = true
  setTimeout(() => {
    isExporting.value = false
    emit('export', { start: startTimeStr.value, end: endTimeStr.value, format: format.value, mode: exportMode.value })
    const link = document.createElement('a'); link.href = '#'; link.setAttribute('download', `evidencia_${props.camera.id}_${Date.now()}.${format.value}`); document.body.appendChild(link)
    emit('close')
  }, 1200)
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 520px; width: 520px;">
      <div class="vms-modal-header vms-flex-between">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">// EXPORTAR TRECHO FORENSE</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 0.75rem; padding: 0.85rem 1.1rem;">
        <div style="background: #0b0e14; border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.4rem 0.65rem; display: flex; justify-content: space-between; align-items: center;">
          <span class="vms-font-semibold vms-text-xs" style="color: #fff;">{{ camera.name }}</span>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan);">{{ camera.resolution }} // {{ camera.codec || 'H.265' }}</span>
        </div>
        <div class="vms-flex-row" style="gap: 0.65rem;">
          <div class="vms-form-group" style="flex: 1;"><label class="vms-label">Horário de Início</label><input v-model="startTimeStr" type="datetime-local" step="1" class="vms-auth-input vms-text-mono" style="font-size: 11px; padding: 3px 6px;" /></div>
          <div class="vms-form-group" style="flex: 1;"><label class="vms-label">Horário de Fim</label><input v-model="endTimeStr" type="datetime-local" step="1" class="vms-auth-input vms-text-mono" style="font-size: 11px; padding: 3px 6px;" /></div>
        </div>
        <div class="vms-flex-row" style="gap: 0.65rem;">
          <div class="vms-form-group" style="flex: 1;"><label class="vms-label">Formato</label><select v-model="format" class="vms-auth-input"><option value="mp4">MP4 (Universal)</option><option value="mkv">MKV (Matroska)</option><option value="avi">AVI (Forense)</option></select></div>
          <div class="vms-form-group" style="flex: 1;"><label class="vms-label">Modo Render</label><select v-model="exportMode" class="vms-auth-input"><option value="stream-copy">[STREAM-COPY] Instantâneo</option><option value="nvenc">[NVENC] GPU Reencode</option></select></div>
        </div>
        <div class="vms-flex-between" style="background: rgba(255, 94, 58, 0.06); border: 1px solid rgba(255, 94, 58, 0.2); border-radius: 6px; padding: 0.4rem 0.65rem;">
          <label style="display: flex; align-items: center; gap: 0.35rem; font-size: 10px; color: #cbd5e1; cursor: pointer;"><input v-model="addWatermark" type="checkbox" /> Marca d'água Timestamp</label>
          <label style="display: flex; align-items: center; gap: 0.35rem; font-size: 10px; color: #cbd5e1; cursor: pointer;"><input v-model="addSha256" type="checkbox" /> Hash Forense SHA-256</label>
        </div>
        <div class="vms-flex-between" style="font-family: var(--vms-font-jetbrains); font-size: 9.5px; color: var(--vms-text-dim);">
          <span>DURAÇÃO: <strong style="color: #fff;">{{ durationFormatted }}</strong></span>
          <span>ESTIMATIVA: <strong style="color: var(--vms-neu-accent-green);">{{ estimatedSize }}</strong></span>
        </div>
      </div>
      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid var(--vms-border); padding: 0.65rem 1.1rem;">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">FECHAR</button>
        <button class="vms-btn vms-btn-primary" :disabled="isExporting || durationSec <= 0" @click="handleDownload">
          <span v-if="isExporting">GERANDO {{ format.toUpperCase() }}...</span>
          <span v-else>BAIXAR</span>
        </button>
      </div>
    </div>
  </div>
</template>
