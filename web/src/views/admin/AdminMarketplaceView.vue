<script setup lang="ts">
import { ref } from 'vue'
import { useMarketplace } from '../../composables/useMarketplace'
import type { PluginManifest } from '../../types/marketplace'
import PluginCatalogGrid from '../../components/admin/marketplace/PluginCatalogGrid.vue'
import PluginDetailsModal from '../../components/admin/marketplace/PluginDetailsModal.vue'

const {
  plugins, searchQuery, selectedCategory, statusFilter, toast,
  installedPlugins, filteredPlugins, installPlugin, uninstallPlugin, togglePlugin
} = useMarketplace()

const detailsPlugin = ref<PluginManifest | null>(null)
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.15rem;">
    <!-- Toast Feedback Notification -->
    <Transition name="vms-toast">
      <div v-if="toast" class="vms-toast-notification" title="Clique para fechar" @click="toast = null"><span>{{ toast }}</span></div>
    </Transition>

    <!-- Header da Loja de Extensões -->
    <div class="vms-flex-between" style="align-items: center; flex-wrap: wrap; gap: 1rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <div class="vms-flex-row" style="align-items: center; gap: 0.65rem;">
          <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange); margin: 0;">MARKETPLACE // LOJA DE ANALÍTICOS</h3>
          <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ installedPlugins.length }}/{{ plugins.length }} INSTALADOS</span>
        </div>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">
          INSTALE ANALÍTICOS DE IA PARA ATIVAR SEUS MÓDULOS DE GESTÃO E EVENTOS NO MENU LATERAL
        </span>
      </div>
    </div>

    <!-- Catálogo e Instalação -->
    <PluginCatalogGrid
      :plugins="filteredPlugins" :search-query="searchQuery" :selected-category="selectedCategory" :status-filter="statusFilter"
      @update:search-query="searchQuery = $event" @update:selected-category="selectedCategory = $event" @update:status-filter="statusFilter = $event"
      @install="installPlugin" @uninstall="uninstallPlugin" @toggle="togglePlugin" @details="detailsPlugin = $event"
    />

    <!-- Modal de Detalhes Técnicos do Plugin -->
    <PluginDetailsModal :is-open="!!detailsPlugin" :plugin="detailsPlugin" @close="detailsPlugin = null" />
  </div>
</template>
