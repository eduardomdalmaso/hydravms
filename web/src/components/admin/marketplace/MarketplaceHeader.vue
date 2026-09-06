<script setup lang="ts">
import type { MarketplaceTab } from '../../../types/marketplace'

defineProps<{
  currentTab: MarketplaceTab
  installedCount: number
  totalPlugins: number
  activeInstancesCount: number
  totalEventsCount: number
}>()

const emit = defineEmits<{ (e: 'changeTab', tab: MarketplaceTab): void }>()
</script>

<template>
  <div class="vms-flex-between" style="align-items: center; flex-wrap: wrap; gap: 1rem;">
    <div class="vms-flex-col" style="gap: 2px;">
      <div class="vms-flex-row" style="align-items: center; gap: 0.65rem;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange); margin: 0;">MARKETPLACE DE ANALÍTICOS</h3>
        <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ installedCount }}/{{ totalPlugins }} ATIVOS</span>
      </div>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">
        LOJA DE EXTENSÕES // GERENCIAMENTO DE INSTÂNCIAS // EXPLORADOR DE TELEMETRIA FORENSE
      </span>
    </div>

    <!-- Navigation Tabs -->
    <div class="vms-marketplace-tabs" style="border-bottom: none; padding-bottom: 0;">
      <button
        class="vms-marketplace-tab-btn"
        :class="{ active: currentTab === 'catalog' }"
        @click="emit('changeTab', 'catalog')"
      >
        [CATÁLOGO // LOJA]
      </button>
      <button
        class="vms-marketplace-tab-btn"
        :class="{ active: currentTab === 'create' }"
        @click="emit('changeTab', 'create')"
      >
        [CRIAR ANALÍTICO ({{ activeInstancesCount }})]
      </button>
      <button
        class="vms-marketplace-tab-btn"
        :class="{ active: currentTab === 'explorer' }"
        @click="emit('changeTab', 'explorer')"
      >
        [EXPLORAR DADOS ({{ totalEventsCount }})]
      </button>
    </div>
  </div>
</template>
