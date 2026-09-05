<script setup lang="ts">
import { ref, watch } from 'vue'
import type { GridLayout } from '../../../types/mosaic'
import type { EnterpriseLayoutItem, LayoutFolderNode, LayoutSlotItem } from '../../../types/layoutTree'

const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: LayoutFolderNode[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', layout: EnterpriseLayoutItem): void }>()

const name = ref(''), grid = ref<GridLayout>('2x2'), folderId = ref(''), isLocked = ref(true)
const selectedUsers = ref<string[]>(['usr_03', 'usr_04'])
const slots = ref<LayoutSlotItem[]>([
  { slotIndex: 0, cameraId: 'cam_01', cameraName: 'Portaria Principal' },
  { slotIndex: 1, cameraId: 'cam_02', cameraName: 'Estacionamento' }
])

const availableCameras = [
  { id: 'cam_01', name: 'Portaria Principal' },
  { id: 'cam_02', name: 'Estacionamento' },
  { id: 'cam_03', name: 'Docas' },
  { id: 'cam_04', name: 'Perímetro' },
  { id: 'cam_05', name: 'NOC' }
]

watch(() => props.isOpen, (open) => {
  if (open) {
    name.value = 'Novo Layout'
    grid.value = '2x2'
    folderId.value = props.targetFolderId || (props.folders[0]?.id || '')
    isLocked.value = true
  }
})

const handleSave = () => {
  if (!name.value.trim()) return
  const targetF = props.folders.find(f => f.id === folderId.value)
  const newLayout: EnterpriseLayoutItem = {
    id: `lay_${Date.now()}`,
    name: name.value.trim(),
    grid: grid.value,
    companyScope: targetF ? targetF.name : 'GLOBAL // RAIZ',
    folderId: folderId.value || undefined,
    is_locked: isLocked.value,
    created_by: 'adminMaster',
    createdAt: '2026-09-05',
    targetScope: 'specific_users',
    allowedUserIds: [...selectedUsers.value],
    slots: [...slots.value]
  }
  emit('save', newLayout)
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 520px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">CRIAR LAYOUT // EMPRESA & CLIENTE</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 0.85rem; max-height: 480px; overflow-y: auto;">
        <div class="vms-form-group">
          <label class="vms-label">NOME DO LAYOUT</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: Portaria & Acessos 2x2" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">EMPRESA / CLIENTE FINAL (PASTA DESTINO)</label>
          <select v-model="folderId" class="vms-auth-input" style="font-size: 11px;">
            <option value="">[RAIZ // SEM PASTA]</option>
            <option v-for="f in folders" :key="f.id" :value="f.id">{{ f.name }}</option>
          </select>
        </div>
        <div class="vms-form-group">
          <label class="vms-label">FORMATO DE GRADE</label>
          <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap;">
            <button v-for="g in (['1x1', '1x2', '2x2', '1+5', '3x3', '4x4'] as GridLayout[])" :key="g"
              class="vms-btn vms-btn-sm" :class="grid === g ? 'vms-btn-primary' : 'vms-btn-secondary'"
              style="font-size: 10px; padding: 3px 8px;" @click="grid = g">{{ g }}</button>
          </div>
        </div>
        <!-- Lock option -->
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
