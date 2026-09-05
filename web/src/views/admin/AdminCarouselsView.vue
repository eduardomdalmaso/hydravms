<script setup lang="ts">
import { ref } from 'vue'

const carousels = ref([
  { id: 'car_01', name: 'Ronda Principal - Entradas', interval: 10, cameras: 4, transition: 'CORTE SECO', status: 'ATIVO' },
  { id: 'car_02', name: 'Ronda Noturna - Perimetro', interval: 15, cameras: 8, transition: 'CROSSFADE', status: 'ATIVO' },
  { id: 'car_03', name: 'Ronda Docas & Expedicao', interval: 8, cameras: 3, transition: 'CORTE SECO', status: 'PAUSADO' }
])

const isModalOpen = ref(false)
const newCarouselName = ref('')
const newInterval = ref(10)

const handleAddCarousel = () => {
  if (!newCarouselName.value) return
  carousels.value.push({
    id: `car_0${carousels.value.length + 1}`,
    name: newCarouselName.value,
    interval: newInterval.value || 10,
    cameras: 2,
    transition: 'CORTE SECO',
    status: 'ATIVO'
  })
  isModalOpen.value = false
  newCarouselName.value = ''
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE RONDAS & CARROSSEL</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PAGINA DE CONFIGURACAO DE SEQUENCIAMENTO AUTOMATICO DE CAMERAS</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVA RONDA]</button>
    </div>

    <!-- Carousels list -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead>
          <tr><th>ID</th><th>NOME DA RONDA</th><th>TEMPO POR CAMERA</th><th>TOTAL CAMERAS</th><th>TRANSICAO</th><th>STATUS</th><th>ACOES</th></tr>
        </thead>
        <tbody>
          <tr v-for="c in carousels" :key="c.id">
            <td class="vms-text-mono">{{ c.id }}</td>
            <td class="vms-font-semibold" style="color: #fff;">{{ c.name }}</td>
            <td><span class="vms-badge vms-badge-warning">{{ c.interval }}s</span></td>
            <td class="vms-text-mono vms-text-xs">{{ c.cameras }} CAMERAS</td>
            <td class="vms-text-mono vms-text-2xs vms-text-dim">{{ c.transition }}</td>
            <td style="text-align: center;"><span class="vms-status-led" :class="c.status === 'ATIVO' ? 'online' : 'warning'" :title="c.status"></span></td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Nova Ronda -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CRIAR NOVA RONDA VIRTUAL</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="isModalOpen = false">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group">
            <label class="vms-label">Nome da Ronda / Sequencia</label>
            <input v-model="newCarouselName" class="vms-auth-input" placeholder="Ex: Ronda Portarias Bloco A" />
          </div>
          <div class="vms-form-group">
            <label class="vms-label">Tempo de Exibicao por Camera (segundos)</label>
            <input v-model.number="newInterval" type="number" min="3" max="120" class="vms-auth-input" />
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddCarousel">SALVAR</button>
        </div>
      </div>
    </div>
  </div>
</template>
