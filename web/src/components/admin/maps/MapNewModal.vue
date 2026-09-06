<script setup lang="ts">
import { ref } from 'vue'
import type { EnterpriseMapItem, MapKind } from '../../../types/mapTree'

const props = defineProps<{ defaultScope?: string }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'create', mapItem: EnterpriseMapItem): void
}>()

const mapName = ref('')
const mapKind = ref<MapKind>('MAP_OPENSOURCE')
const scope = ref(props.defaultScope || 'EMPRESA // ALPHA SEGURANCA')
const isLocked = ref(false)
const floorplanFile = ref('planta_setor_novo.png')

const handleSubmit = () => {
  if (!mapName.value.trim()) return
  const newMap: EnterpriseMapItem = {
    id: `map_${Date.now()}`,
    name: mapName.value.trim(),
    type: mapKind.value,
    companyScope: scope.value,
    is_locked: isLocked.value,
    initialCenter: [-23.55052, -46.633308],
    initialZoom: 16,
    floorplanFile: mapKind.value === 'PLANTA_BAIXA' ? floorplanFile.value : undefined,
    markers: [],
    createdAt: new Date().toISOString().split('T')[0]
  }
  emit('create', newMap)
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">NOVO MAPA OU PLANTA BAIXA</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="vms-modal-body">
        <div class="vms-form-group">
          <label class="vms-label">TIPO DE VISUALIZACAO</label>
          <div class="vms-flex-row" style="gap: 1rem;">
            <label class="vms-flex-row" style="gap: 0.4rem; align-items: center; cursor: pointer;">
              <input v-model="mapKind" type="radio" value="MAP_OPENSOURCE" />
              <span class="vms-text-xs vms-font-bold" style="color: #ffffff;">MAPA COLORIDO OPENSOURCE</span>
            </label>
            <label class="vms-flex-row" style="gap: 0.4rem; align-items: center; cursor: pointer;">
              <input v-model="mapKind" type="radio" value="PLANTA_BAIXA" />
              <span class="vms-text-xs vms-font-bold" style="color: #ff5e3a;">PLANTA BAIXA ARQUITETONICA</span>
            </label>
          </div>
        </div>

        <div class="vms-form-group">
          <label class="vms-label">NOME DO MAPA / SETOR</label>
          <input v-model="mapName" class="vms-auth-input" placeholder="Ex: Perimetro Fabril Norte & Docas" autofocus />
        </div>

        <div v-if="mapKind === 'PLANTA_BAIXA'" class="vms-form-group">
          <label class="vms-label">ARQUIVO DA PLANTA (PNG/SVG/JPEG)</label>
          <input v-model="floorplanFile" class="vms-auth-input" placeholder="Ex: planta_interna_docas.svg" />
        </div>
      </div>

      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" :disabled="!mapName.trim()" @click="handleSubmit">CRIAR MAPA</button>
      </div>
    </div>
  </div>
</template>
