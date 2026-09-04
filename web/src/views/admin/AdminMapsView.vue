<script setup lang="ts">
import { ref } from 'vue'

const maps = ref([
  { id: 'map_01', name: 'Planta Baixa - Galpao 01', file: 'galpao1_planta.png', cameras: 8, resolution: '2048x1536' },
  { id: 'map_02', name: 'Perimetro & Acessos Externos', file: 'perimetro_master.svg', cameras: 14, resolution: '3840x2160' },
  { id: 'map_03', name: 'Predio Administrativo - 1 Andar', file: 'predio_adm_01.png', cameras: 6, resolution: '1920x1080' }
])

const isModalOpen = ref(false)
const newMapName = ref('')
const newMapFile = ref('')

const handleAddMap = () => {
  if (!newMapName.value) return
  maps.value.push({
    id: `map_0${maps.value.length + 1}`,
    name: newMapName.value,
    file: newMapFile.value || 'mapa_novo.png',
    cameras: 0,
    resolution: '1920x1080'
  })
  isModalOpen.value = false
  newMapName.value = ''
  newMapFile.value = ''
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE PLANTAS BAIXAS & MAPAS</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">PAGINA DE GESTAO DE PLANTAS COM GEO-LOCALIZACAO DE CAMERAS</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVA PLANTA BAIXA]</button>
    </div>

    <!-- Maps list -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead>
          <tr><th>ID</th><th>NOME DA PLANTA</th><th>ARQUIVO</th><th>RESOLUCAO</th><th>CAMERAS VINCULADAS</th><th>ACOES</th></tr>
        </thead>
        <tbody>
          <tr v-for="m in maps" :key="m.id">
            <td class="vms-text-mono">{{ m.id }}</td>
            <td class="vms-font-semibold" style="color: #fff;">{{ m.name }}</td>
            <td class="vms-text-mono vms-text-2xs vms-text-dim">{{ m.file }}</td>
            <td class="vms-text-mono vms-text-xs">{{ m.resolution }}</td>
            <td><span class="vms-badge vms-channel-pill-telegram">{{ m.cameras }} CAMERAS</span></td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR PINOS]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Novo Mapa -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CADASTRAR NOVA PLANTA BAIXA</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="isModalOpen = false">[X]</button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group">
            <label class="vms-label">Nome da Planta / Setor</label>
            <input v-model="newMapName" class="vms-auth-input" placeholder="Ex: Galpao Logistica Setor B" />
          </div>
          <div class="vms-form-group">
            <label class="vms-label">Nome do Arquivo (PNG/SVG/JPEG)</label>
            <input v-model="newMapFile" class="vms-auth-input" placeholder="Ex: galpao_setor_b.png" />
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddMap">SALVAR PLANTA</button>
        </div>
      </div>
    </div>
  </div>
</template>
