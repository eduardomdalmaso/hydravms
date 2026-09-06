<script setup lang="ts">
import { ref, watch } from 'vue'
import type { EnterpriseRondaItem, RondaFolderNode } from '../../../types/rondaTree'

const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: RondaFolderNode[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', ronda: EnterpriseRondaItem): void }>()

const name = ref('')
const folderId = ref(props.targetFolderId || '')
const transition = ref<'CORTE SECO' | 'CROSSFADE' | 'FADE PRETO'>('CORTE SECO')
const defaultInterval = ref(10)
const isLocked = ref(false)

watch(() => props.isOpen, (open) => {
  if (open) {
    name.value = ''
    folderId.value = props.targetFolderId || ''
    transition.value = 'CORTE SECO'
    defaultInterval.value = 10
    isLocked.value = false
  }
})

const handleSave = () => {
  if (!name.value.trim()) return
  const folder = props.folders.find(f => f.id === folderId.value)
  const newRonda: EnterpriseRondaItem = {
    id: `ronda_${Date.now()}`,
    name: name.value.trim(),
    companyScope: folder ? folder.name : 'GLOBAL // TODAS AS UNIDADES',
    folderId: folderId.value || undefined,
    is_locked: isLocked.value,
    status: 'ATIVO',
    transition: transition.value,
    allowedUserIds: ['usr_02', 'usr_03'],
    createdAt: new Date().toISOString().slice(0, 16).replace('T', ' '),
    streams: [
      { id: `rs_${Date.now()}_1`, cameraId: 'cam_01', cameraName: 'Portaria Principal (Entrada)', resolution: '1080P', fps: 30, intervalSeconds: defaultInterval.value, orderIndex: 0 },
      { id: `rs_${Date.now()}_2`, cameraId: 'cam_02', cameraName: 'Estacionamento Visitantes', resolution: '1080P', fps: 25, intervalSeconds: defaultInterval.value, orderIndex: 1 }
    ]
  }
  emit('save', newRonda)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 480px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CRIAR NOVA RONDA VIRTUAL</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 1rem;">
        <div class="vms-form-group">
          <label class="vms-label">NOME DA RONDA</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: Ronda Portaria & Estacionamento" autofocus />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">EMPRESA / CLIENTE FINAL (PASTA DESTINO)</label>
          <select v-model="folderId" class="vms-auth-input" style="font-size: 11px;">
            <option value="">[RAIZ // GLOBAL]</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
          </select>
        </div>
        <div class="vms-form-group">
          <label class="vms-label">TIPO DE TRANSICAO</label>
          <div class="vms-flex-row" style="gap: 0.4rem;">
            <button v-for="t in (['CORTE SECO', 'CROSSFADE', 'FADE PRETO'] as const)" :key="t"
              class="vms-btn vms-btn-sm" :class="transition === t ? 'vms-btn-primary' : 'vms-btn-secondary'"
              style="font-size: 9.5px; padding: 3px 8px;" @click="transition = t">{{ t }}</button>
          </div>
        </div>
        <div class="vms-flex-between" style="background: #07080c; padding: 0.6rem 0.8rem; border-radius: 6px; border: 1px solid var(--vms-border); cursor: pointer;" @click="isLocked = !isLocked">
          <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            <span class="vms-text-xs vms-font-bold" style="color: #fff;">TRAVAR COM CADEADO (IMUTAVEL)</span>
          </div>
          <span class="vms-badge" :class="isLocked ? 'vms-badge-danger' : 'vms-badge-neutral'">{{ isLocked ? '[ATIVO]' : '[INATIVO]' }}</span>
        </div>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSave">SALVAR</button>
      </div>
    </div>
  </div>
</template>
