<script setup lang="ts">
import { ref } from 'vue'

const pools = ref([
  { id: 'p1', name: 'NVMe Hot Tier 01', path: '/var/lib/hydravms/hot-01', tier: 'HOT', usedGb: 340, totalGb: 1920, status: 'ONLINE' },
  { id: 'p2', name: 'HDD Warm Archive 01', path: '/var/lib/hydravms/warm-01', tier: 'WARM', usedGb: 8400, totalGb: 18000, status: 'ONLINE' },
  { id: 'p3', name: 'PostgreSQL Events DB', path: '/var/lib/hydravms/db-events', tier: 'DATABASE', usedGb: 48, totalGb: 500, status: 'ONLINE' }
])

const handleAddDisk = () => {
  window.alert('Iniciando assistente de formatacao e montagem de novo disco...')
}
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1.5rem;">
    <div class="vms-page-header">
      <div class="vms-page-title-group">
        <h1 class="vms-h1">DISCOS & POOLS DE ARMAZENAMENTO</h1>
        <span class="vms-text-sm vms-text-muted">Gestao de buffer quente NVMe, expurgo automatico StorageGuard e HDs de longo prazo.</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="handleAddDisk">[+ ADICIONAR NOVO DISCO]</button>
    </div>

    <div class="vms-grid-container" style="grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.25rem;">
      <div v-for="p in pools" :key="p.id" class="vms-card vms-flex-col" style="padding: 1.25rem; gap: 0.75rem;">
        <div class="vms-flex-between">
          <span class="vms-badge" :class="p.tier === 'HOT' ? 'vms-channel-pill-telegram' : 'vms-badge-info'">
            [{{ p.tier }}]
          </span>
          <span class="vms-status-led" :class="p.status === 'ONLINE' ? 'online' : 'offline'" :title="p.status"></span>
        </div>
        <h3 class="vms-h3" style="color: #fff;">{{ p.name }}</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ p.path }}</span>
        
        <!-- Usage Bar -->
        <div class="vms-flex-col" style="gap: 0.25rem; margin-top: 0.5rem;">
          <div class="vms-flex-between vms-text-xs vms-text-muted">
            <span>Uso: {{ p.usedGb }} GB</span>
            <span>Total: {{ p.totalGb }} GB ({{ Math.round((p.usedGb / p.totalGb) * 100) }}%)</span>
          </div>
          <div class="vms-neu-track">
            <div class="vms-neu-track-fill" :style="{ width: `${(p.usedGb / p.totalGb) * 100}%` }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
