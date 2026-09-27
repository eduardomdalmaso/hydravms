<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { GridLayout } from '../../types/mosaic'
import { fetchLiveLayouts, createRemoteLayout } from '../../services/adminApi'

const emit = defineEmits<{ (e: 'applyLayout', layout: GridLayout): void }>()

const layoutName = ref('')
const selectedGrid = ref<GridLayout>('2x2')
const savedLayouts = ref<{ id: string; name: string; grid: GridLayout; camCount: number }[]>([])

const gridTemplates: { id: GridLayout; label: string; slots: number }[] = [
  { id: '1x1', label: '1x1 (Hero)', slots: 1 },
  { id: '1x2', label: '1x2 (2 Cameras)', slots: 2 },
  { id: '2x2', label: '2x2 (4 Cameras)', slots: 4 },
  { id: '3x3', label: '3x3 (9 Cameras)', slots: 9 },
  { id: '4x4', label: '4x4 (16 Cameras)', slots: 16 },
  { id: '5x5', label: '5x5 (25 Cameras)', slots: 25 },
  { id: '6x6', label: '6x6 (36 Cameras)', slots: 36 },
  { id: '7x7', label: '7x7 (49 Cameras)', slots: 49 },
  { id: '8x8', label: '8x8 (64 Cameras)', slots: 64 },
  { id: '10x10', label: '10x10 (100 Cameras)', slots: 100 },
  { id: '1+5', label: '1+5 (1 Destaque + 5)', slots: 6 },
  { id: '1+8', label: '1+8 (1 Destaque + 8)', slots: 9 },
  { id: '1+7', label: '1+7 (1 Destaque + 7)', slots: 8 }
]

const loadLayouts = async () => {
  const layouts = await fetchLiveLayouts()
  savedLayouts.value = layouts.map(l => ({
    id: l.id,
    name: l.name,
    grid: l.grid,
    camCount: l.slots?.length || 0
  }))
}

onMounted(loadLayouts)

const handleSaveLayout = async () => {
  if (!layoutName.value.trim()) return
  await createRemoteLayout({
    name: layoutName.value.trim(),
    grid: selectedGrid.value,
    is_locked: false,
    allowedUserIds: [],
    slots: []
  })
  layoutName.value = ''
  await loadLayouts()
}
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1.5rem;">
    <div class="vms-page-header">
      <div class="vms-page-title-group">
        <h1 class="vms-h1">DESIGNER & CRIADOR DE LAYOUTS</h1>
        <span class="vms-text-sm vms-text-muted">Monte grades personalizadas de mosaico para operacoes de monitoramento.</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="handleSaveLayout">SALVAR</button>
    </div>

    <div class="vms-grid-container" style="grid-template-columns: 1fr 1fr; gap: 1.25rem;">
      <!-- Grid Template Picker -->
      <div class="vms-card vms-flex-col" style="gap: 1rem; padding: 1.25rem;">
        <h3 class="vms-h3">1. DEFINIR GRADE DO MOSAICO</h3>
        <div class="vms-form-group">
          <label class="vms-label">Nome do Layout</label>
          <input v-model="layoutName" class="vms-auth-input" placeholder="Ex: Portaria & Docas" />
        </div>
        <div class="vms-grid-container" style="grid-template-columns: repeat(2, 1fr); gap: 0.75rem;">
          <div
            v-for="tmpl in gridTemplates"
            :key="tmpl.id"
            class="vms-card vms-card-elevated"
            :style="{
              padding: '0.85rem',
              cursor: 'pointer',
              borderColor: selectedGrid === tmpl.id ? 'var(--vms-primary)' : 'transparent',
              boxShadow: selectedGrid === tmpl.id ? 'var(--vms-primary-glow)' : 'var(--vms-shadow-convex-sm)'
            }"
            @click="selectedGrid = tmpl.id"
          >
            <span class="vms-text-sm vms-font-semibold" :style="{ color: selectedGrid === tmpl.id ? 'var(--vms-primary)' : '#fff' }">
              {{ tmpl.label }}
            </span>
          </div>
        </div>
      </div>

      <!-- Saved Layouts List -->
      <div class="vms-card vms-flex-col" style="gap: 1rem; padding: 1.25rem;">
        <h3 class="vms-h3">2. LAYOUTS SALVOS</h3>
        <div class="vms-flex-col" style="gap: 0.5rem; overflow-y: auto; max-height: 280px;">
          <div
            v-for="lay in savedLayouts"
            :key="lay.id"
            class="vms-card vms-card-elevated vms-flex-between"
            style="padding: 0.75rem 1rem;"
          >
            <div class="vms-flex-col" style="gap: 0.125rem;">
              <span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ lay.name }}</span>
              <span class="vms-text-mono vms-text-2xs vms-text-dim">GRADE: {{ lay.grid }} // {{ lay.camCount }} SLOTS</span>
            </div>
            <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('applyLayout', lay.grid)">
              [CARREGAR]
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
