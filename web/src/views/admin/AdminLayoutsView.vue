<script setup lang="ts">
import { ref } from 'vue'

const layouts = ref([
  { id: 'lay_01', name: 'Padrao 2x2 - Perimetro', grid: '2x2', enabled_for: 'TODOS OS OPERADORES', cameras: 4 },
  { id: 'lay_02', name: 'Visao Destaque 1+5', grid: '1+5', enabled_for: 'TODOS OS OPERADORES', cameras: 6 },
  { id: 'lay_03', name: 'Monitoramento 3x3', grid: '3x3', enabled_for: 'SUPERVISORES', cameras: 9 },
  { id: 'lay_04', name: 'Mural 4x4 - Galpao Geral', grid: '4x4', enabled_for: 'ADMINISTRADORES', cameras: 16 }
])

const isModalOpen = ref(false)
const newLayoutName = ref('')
const selectedGrid = ref('2x2')

const handleAddLayout = () => {
  if (!newLayoutName.value) return
  layouts.value.push({
    id: `lay_0${layouts.value.length + 1}`,
    name: newLayoutName.value,
    grid: selectedGrid.value,
    enabled_for: 'TODOS OS OPERADORES',
    cameras: selectedGrid.value === '1x1' ? 1 : selectedGrid.value === '2x2' ? 4 : 9
  })
  isModalOpen.value = false
  newLayoutName.value = ''
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE LAYOUTS DE GRADE</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PAGINA DE COMPOSICAO DE TELAS MULTI-VIEW</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVO LAYOUT]</button>
    </div>

    <!-- Layouts list -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead>
          <tr><th>ID</th><th>NOME DO LAYOUT</th><th>GRADE</th><th>TOTAL CANAIS</th><th>HABILITADO PARA</th><th>ACOES</th></tr>
        </thead>
        <tbody>
          <tr v-for="l in layouts" :key="l.id">
            <td class="vms-text-mono">{{ l.id }}</td>
            <td class="vms-font-semibold" style="color: #fff;">{{ l.name }}</td>
            <td><span class="vms-badge vms-badge-info">[{{ l.grid }}]</span></td>
            <td class="vms-text-mono vms-text-xs">{{ l.cameras }} SLOTS</td>
            <td class="vms-text-mono vms-text-2xs vms-text-dim">{{ l.enabled_for }}</td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR GRADE]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Novo Layout -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CRIAR NOVO LAYOUT DE GRADE</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="isModalOpen = false">[X]</button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group">
            <label class="vms-label">Nome do Layout</label>
            <input v-model="newLayoutName" class="vms-auth-input" placeholder="Ex: Portaria + Acesso Principal" />
          </div>
          <div class="vms-form-group">
            <label class="vms-label">Formato de Grade</label>
            <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
              <button v-for="g in ['1x1', '1x2', '2x2', '1+5', '3x3', '4x4']" :key="g"
                class="vms-btn vms-btn-sm" :class="selectedGrid === g ? 'vms-btn-primary' : 'vms-btn-secondary'"
                @click="selectedGrid = g">{{ g }}</button>
            </div>
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddLayout">SALVAR LAYOUT</button>
        </div>
      </div>
    </div>
  </div>
</template>
