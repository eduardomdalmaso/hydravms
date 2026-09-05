<script setup lang="ts">
import { ref, computed } from 'vue'
import type { EnterpriseLayoutItem } from '../../../types/layoutTree'
import type { GridLayout } from '../../../types/mosaic'
import LayoutCompanyUsersPicker from './LayoutCompanyUsersPicker.vue'

const props = defineProps<{ layout: EnterpriseLayoutItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const localName = ref(props.layout.name)
const localGrid = ref<GridLayout>(props.layout.grid)
const localScope = ref(props.layout.companyScope)
const localLocked = ref(props.layout.is_locked)
const localUsers = ref<string[]>([...props.layout.allowedUserIds])

const isCompany = computed(() => localScope.value.toUpperCase().startsWith('EMPRESA'))

const gridOptions: GridLayout[] = ['1x1', '1x2', '2x2', '1+5', '3x3', '4x4']
const scopeOptions = [
  'EMPRESA // ALPHA SEGURANCA',
  'EMPRESA // BETA VIGILANCIA',
  'CLIENTE // CONDOMINIO VISTA VERDE',
  'CLIENTE // EDIFICIO HORIZONTE'
]

const handleSelectGrid = (g: GridLayout) => { localGrid.value = g; props.layout.grid = g }
const toggleUser = (id: string) => {
  const idx = localUsers.value.indexOf(id)
  if (idx >= 0) localUsers.value.splice(idx, 1); else localUsers.value.push(id)
}

const handleSave = () => {
  props.layout.name = localName.value
  props.layout.grid = localGrid.value
  props.layout.companyScope = localScope.value
  props.layout.is_locked = localLocked.value
  props.layout.allowedUserIds = [...localUsers.value]
  emit('saved', `[LAYOUT] "${props.layout.name}" atualizado com sucesso`)
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
      <label class="vms-label">SELECAO DE EMPRESA OU CLIENTE</label>
      <select v-model="localScope" class="vms-auth-input" style="font-size: 11px; padding: 4px 6px;">
        <option v-for="opt in scopeOptions" :key="opt" :value="opt">[{{ opt }}]</option>
      </select>
    </div>

    <!-- Se for Empresa, abre os Usuários com Permissão -->
    <LayoutCompanyUsersPicker v-if="isCompany" :company-name="localScope" :allowed-user-ids="localUsers" @toggle-user="toggleUser" />

    <!-- Trava com Cadeado -->
    <div class="vms-flex-between" style="background: #07080c !important; padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border); cursor: pointer;" @click="localLocked = !localLocked">
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" :stroke="localLocked ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)'" stroke-width="2.5"><rect x="3" y="11" width="18" height="11" rx="2"/><path :d="localLocked ? 'M7 11V7a5 5 0 0 1 10 0v4' : 'M7 11V7a5 5 0 0 1 9.9-1'"/></svg>
        <span class="vms-text-xs vms-font-bold" :style="{ color: localLocked ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">{{ localLocked ? '[TRAVADO COM CADEADO]' : '[DESTRAVADO]' }}</span>
      </div>
      <span class="vms-badge" :class="localLocked ? 'vms-badge-danger' : 'vms-badge-neutral'" style="font-size: 8px;">{{ localLocked ? '[LOCK]' : '[LIVRE]' }}</span>
    </div>
  </div>
</template>
