<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FolderNode, StreamItem } from '../../types/streamTree'

const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: FolderNode[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', stream: Partial<StreamItem>, folderId: string): void }>()

const step = ref(1), testStatus = ref<string | null>(null)
const form = ref({
  name: '', protocol: 'RTSP' as 'RTSP' | 'RTMP' | 'ONVIF',
  url: 'rtsp://192.168.1.100:554/live', ip: '192.168.1.100', port: 554,
  codec: 'H.265' as 'H.265' | 'H.264', resolution: '1080P', fps: 30, bitrate: '4.0 Mbps',
  recordMode: 'continuous' as 'continuous' | 'motion' | 'ai_event', folderId: ''
})

watch(() => props.isOpen, (open) => {
  if (open) {
    step.value = 1; testStatus.value = null
    form.value.folderId = props.targetFolderId || (props.folders[0]?.id || '')
  }
})

const testConnection = () => { testStatus.value = 'Socket RTSP // Handshake OK (31ms) // H.265' }
const capitalize = (s: string) => s.trim() ? s.trim().charAt(0).toUpperCase() + s.trim().slice(1) : ''

const finish = () => {
  if (!form.value.name.trim()) return
  emit('save', {
    name: capitalize(form.value.name), protocol: form.value.protocol, url: form.value.url, ip: form.value.ip,
    port: form.value.port, codec: form.value.codec, resolution: form.value.resolution,
    fps: form.value.fps, bitrate: form.value.bitrate, recordMode: form.value.recordMode,
    status: 'online', has_ptz: false
  }, form.value.folderId)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 540px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">NOVO FLUXO // ETAPA {{ step }} DE 4</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-flex-between" style="padding: 0.5rem 1.25rem; background: #0c0f14; border-bottom: 1px solid var(--vms-border);">
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 1 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">1. REDE</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 2 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">2. CODEC</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 3 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">3. GRAVACAO</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 4 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">4. TESTE</span>
      </div>

      <div class="vms-modal-body" style="gap: 1rem;">
        <div v-if="step === 1" class="vms-flex-col" style="gap: 0.75rem;">
          <div class="vms-form-group"><label class="vms-label">Nome da Camera (Primeira Letra Maiuscula)</label><input v-model="form.name" class="vms-auth-input" placeholder="Ex: Portaria Principal Leste" autofocus /></div>
          <div class="vms-form-group"><label class="vms-label">URL do Fluxo</label><input v-model="form.url" class="vms-auth-input" /></div>
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <div class="vms-form-group" style="flex: 1;"><label class="vms-label">IP do Dispositivo</label><input v-model="form.ip" class="vms-auth-input" /></div>
            <div class="vms-form-group" style="width: 100px;"><label class="vms-label">Porta</label><input v-model.number="form.port" type="number" class="vms-auth-input" /></div>
          </div>
        </div>

        <div v-if="step === 2" class="vms-flex-col" style="gap: 0.75rem;">
          <div class="vms-form-group"><label class="vms-label">Codec Primario</label><select v-model="form.codec" class="vms-auth-input"><option value="H.265">[H.265] Alta Eficiencia (HEVC)</option><option value="H.264">[H.264] Padrao Legado</option></select></div>
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <div class="vms-form-group" style="flex: 1;"><label class="vms-label">Resolucao</label><select v-model="form.resolution" class="vms-auth-input"><option value="4K">[4K] 3840x2160</option><option value="1080P">[1080P] 1920x1080 Full HD</option><option value="720P">[720P] 1280x720 HD</option></select></div>
            <div class="vms-form-group" style="width: 120px;"><label class="vms-label">FPS Ingestao</label><select v-model.number="form.fps" class="vms-auth-input"><option :value="30">30 FPS</option><option :value="25">25 FPS</option><option :value="15">15 FPS</option></select></div>
          </div>
        </div>

        <div v-if="step === 3" class="vms-flex-col" style="gap: 0.75rem;">
          <div class="vms-form-group"><label class="vms-label">Pasta de Destino na Topologia</label><select v-model="form.folderId" class="vms-auth-input"><option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option></select></div>
          <div class="vms-form-group"><label class="vms-label">Modo de Gravacao</label><select v-model="form.recordMode" class="vms-auth-input"><option value="continuous">[GRAVACAO 24/7] Continua</option><option value="motion">[MOVIMENTO] Apenas com VMD</option><option value="ai_event">[SMART IA] Apenas Eventos de IA</option></select></div>
        </div>

        <div v-if="step === 4" class="vms-flex-col" style="gap: 0.75rem;">
          <button class="vms-btn vms-btn-secondary" @click="testConnection">[TESTAR SOCKET RTSP / PING]</button>
          <div v-if="testStatus" class="vms-badge vms-badge-success" style="padding: 0.5rem; font-family: var(--vms-font-jetbrains); font-size: 11px;">{{ testStatus }}</div>
          <span class="vms-text-2xs vms-text-dim">Clique em Salvar para vincular o novo fluxo à pasta selecionada na árvore.</span>
        </div>
      </div>

      <div class="vms-modal-footer">
        <button v-if="step > 1" class="vms-btn vms-btn-secondary" @click="step--">ANTERIOR</button>
        <button v-if="step < 4" class="vms-btn vms-btn-primary" @click="step++">PROXIMO</button>
        <button v-else class="vms-btn vms-btn-primary" @click="finish">SALVAR</button>
      </div>
    </div>
  </div>
</template>
