<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { EnterpriseLayoutItem } from '../../../types/layoutTree'
import type { GridLayout } from '../../../types/mosaic'
import LayoutCompanyUsersPicker from './LayoutCompanyUsersPicker.vue'
import { fetchFolders } from '../../../services/api'

const props = defineProps<{ layout: EnterpriseLayoutItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const localName = ref(props.layout.name)
const localGrid = ref<GridLayout>(props.layout.grid)
const localFolderId = ref(props.layout.folderId || '')
const localLocked = ref(props.layout.is_locked)
const localTargetMonitor = ref<number>(props.layout.target_monitor || 0)
const localUsers = ref<string[]>([...props.layout.allowedUserIds])
const availableFolders = ref<{ id: string; name: string }[]>([])

onMounted(async () => {
  const dbF = await fetchFolders('layouts')
  availableFolders.value = dbF || []
})

watch(() => props.layout.is_locked, (v) => localLocked.value = v)
watch(localLocked, (v) => props.layout.is_locked = v)
watch(() => props.layout.target_monitor, (v) => localTargetMonitor.value = v || 0)
watch(localTargetMonitor, (v) => props.layout.target_monitor = v)

const gridOptions: GridLayout[] = [
  '1x1', '1x2', '2x2', '3x3', '4x4', '5x5', '6x6', '7x7', '8x8', '10x10', '1+5', '1+8', '1+7'
]

const handleSelectGrid = (g: GridLayout) => { localGrid.value = g; props.layout.grid = g }
const toggleUser = (id: string) => {
  const idx = localUsers.value.indexOf(id)
  if (idx >= 0) localUsers.value.splice(idx, 1); else localUsers.value.push(id)
}

const handleFolderChange = () => {
  props.layout.folderId = localFolderId.value || undefined
  const target = availableFolders.value.find(f => f.id === localFolderId.value)
  props.layout.companyScope = target ? target.name : 'GLOBAL // RAIZ'
}

const handleSave = () => {
  props.layout.name = localName.value
  props.layout.grid = localGrid.value
  props.layout.is_locked = localLocked.value
  props.layout.target_monitor = localTargetMonitor.value
  props.layout.allowedUserIds = [...localUsers.value]
  handleFolderChange()
  emit('saved', `[LAYOUT] "${props.layout.name}" atualizado`)
}
const toggleLock = () => {
  localLocked.value = !localLocked.value; props.layout.is_locked = localLocked.value
  emit('saved', localLocked.value ? `[TRAVADO COM CADEADO] Grade "${props.layout.name}" bloqueada` : `[DESTRAVADO] Grade "${props.layout.name}" liberada`)
}
</script>

<template>
  <div class="vms-split-pane">
    <div class="vms-split-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">CONFIGURACOES DO LAYOUT</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// PARAMETROS, EMPRESA & ACESSO</span>
      </div>
      <span class="vms-badge" style="background: #14171d; color: #ff5e3a !important; border: 1px solid var(--vms-border);">[{{ localGrid }}]</span>
    </div>

    <div class="vms-form-group">
      <label class="vms-label">NOME DO LAYOUT</label>
      <input v-model="localName" class="vms-auth-input" style="font-size: 11px; padding: 4px 8px;" />
    </div>

    <!-- Alterar Formato da Grade -->
    <div class="vms-form-group">
      <label class="vms-label">FORMATO DA GRADE</label>
      <div class="vms-flex-row" style="gap: 0.35rem; flex-wrap: wrap;">
        <button v-for="g in gridOptions" :key="g" class="vms-btn vms-btn-sm" :class="localGrid === g ? 'vms-btn-primary' : 'vms-btn-secondary'" style="font-size: 9.5px; padding: 2px 7px;" @click="handleSelectGrid(g)">
          {{ g }}
        </button>
      </div>
    </div>

    <!-- Seleção da Empresa ou Cliente -->
    <div class="vms-form-group">
      <label class="vms-label">EMPRESA / CLIENTE FINAL (PASTA)</label>
      <select v-model="localFolderId" class="vms-auth-input" style="font-size: 11px; padding: 4px 6px;" @change="handleFolderChange">
        <option value="">[RAIZ // SEM PASTA]</option>
        <option v-for="f in availableFolders" :key="f.id" :value="f.id">{{ f.name }}</option>
      </select>
    </div>

    <!-- Monitor Alvo / Destino Multi-Tela -->
    <div class="vms-form-group">
      <label class="vms-label">MONITOR ALVO (MULTI-TELA / VIDEO WALL)</label>
      <select v-model="localTargetMonitor" class="vms-auth-input" style="font-size: 11px; padding: 4px 6px;">
        <option :value="0">[LIVRE] // QUALQUER MONITOR ATIVO</option>
        <option :value="1">[MONITOR 01] // TELA PRINCIPAL (CONTROL ROOM)</option>
        <option :value="2">[MONITOR 02] // MONITOR SECUNDÁRIO (POP-OUT)</option>
        <option :value="3">[MONITOR 03] // VIDEO WALL 01 (PERÍMETRO)</option>
        <option :value="4">[MONITOR 04] // VIDEO WALL 02 (ALARMES & IA)</option>
        <option :value="5">[MONITOR 05] // VIDEO WALL 03 (SUPERVISÃO)</option>
        <option :value="6">[MONITOR 06] // VIDEO WALL 04 (DOCAS & CARGA)</option>
        <option :value="7">[MONITOR 07] // VIDEO WALL 05 (ESTACIONAMENTO)</option>
        <option :value="8">[MONITOR 08] // VIDEO WALL 06 (PORTARIA)</option>
      </select>
    </div>

    <!-- Usuários com Permissão -->
    <LayoutCompanyUsersPicker :company-name="layout.companyScope || 'EMPRESA'" :allowed-user-ids="localUsers" @toggle-user="toggleUser" />

    <!-- Trava com Cadeado -->
    <div class="vms-flex-between" style="background: #07080c !important; padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border); cursor: pointer;" @click="toggleLock">
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <svg width="14" height="14" viewBox="0 0 24 24" :fill="localLocked ? '#ff5e3a' : 'none'" stroke="#ff5e3a" stroke-width="1.5">
          <path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/>
        </svg>
        <span class="vms-text-xs vms-font-bold" :style="{ color: localLocked ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">{{ localLocked ? '[TRAVADO COM CADEADO]' : '[DESTRAVADO]' }}</span>
      </div>
      <span class="vms-badge" :class="localLocked ? 'vms-badge-danger' : 'vms-badge-neutral'" style="font-size: 8px;">{{ localLocked ? '[LOCK]' : '[LIVRE]' }}</span>
    </div>
  </div>
</template>
