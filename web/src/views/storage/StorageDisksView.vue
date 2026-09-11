<script setup lang="ts">
import { computed } from 'vue'
import { useStorageManager } from '../../composables/useStorageManager'
import StorageTelemetryHeader from '../../components/admin/storage/StorageTelemetryHeader.vue'
import StorageTabsFilter from '../../components/admin/storage/StorageTabsFilter.vue'
import StoragePoolCard from '../../components/admin/storage/StoragePoolCard.vue'
import StorageAddModal from '../../components/admin/storage/StorageAddModal.vue'

const {
  pools, unallocatedDisks, filterRole, searchQuery, isAddModalOpen, notification,
  totalCapacityGb, totalUsedGb, overallPercentage, hotUsagePercent, filteredPools,
  addPool, removePool, triggerSpilloverDrain
} = useStorageManager()

const counts = computed(() => ({
  all: pools.value.length,
  buffer: pools.value.filter(p => p.role === 'HOT_BUFFER').length,
  recordings: pools.value.filter(p => p.role === 'WARM_ARCHIVE').length,
  snapshots: pools.value.filter(p => p.role === 'SNAPSHOTS').length,
  database: pools.value.filter(p => p.role === 'DATABASE').length
}))
</script>

<template>
  <div class="vms-content-area vms-flex-col" style="gap: 1.25rem;">
    <!-- Toast Feedback Notification -->
    <div v-if="notification" class="vms-badge vms-badge-orange" style="padding: 8px 14px; font-weight: bold; border-radius: 4px;">
      {{ notification }}
    </div>

    <!-- Telemetria Geral de Capacidade & Botoes Globais -->
    <StorageTelemetryHeader
      :totalCapacityGb="totalCapacityGb"
      :totalUsedGb="totalUsedGb"
      :overallPercentage="overallPercentage"
      :hotUsagePercent="hotUsagePercent"
      @triggerDrain="triggerSpilloverDrain"
      @openAddModal="isAddModalOpen = true"
    />

    <!-- Abas de Filtro: Todos, Buffer, Gravacoes, Snapshots, Banco de Dados -->
    <StorageTabsFilter
      :activeFilter="filterRole"
      :searchQuery="searchQuery"
      :counts="counts"
      @update:filter="filterRole = $event"
      @update:search="searchQuery = $event"
    />

    <!-- Grid de Storages & Discos -->
    <div v-if="filteredPools.length > 0" class="vms-grid-container" style="grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.25rem;">
      <StoragePoolCard
        v-for="p in filteredPools"
        :key="p.id"
        :pool="p"
        @remove="removePool"
      />
    </div>

    <!-- Empty State -->
    <div v-else class="vms-card vms-flex-col vms-flex-center" style="padding: 3rem; text-align: center; gap: 0.5rem;">
      <span class="vms-text-muted vms-text-mono vms-text-sm">[NENHUM STORAGE LOCALIZADO COM OS FILTROS APLICADOS]</span>
      <span class="vms-text-2xs vms-text-dim">Clique no ícone de storage no cabeçalho para registrar um disco físico local ou servidor de rede.</span>
    </div>

    <!-- Modal de Cadastro Unificado -->
    <StorageAddModal
      v-if="isAddModalOpen"
      :unallocatedDisks="unallocatedDisks"
      @close="isAddModalOpen = false"
      @save="addPool"
    />
  </div>
</template>
