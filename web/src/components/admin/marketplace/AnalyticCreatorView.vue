<script setup lang="ts">
import type { PluginManifest, AnalyticInstance } from '../../../types/marketplace'
import AnalyticInstanceForm from './AnalyticInstanceForm.vue'
import AnalyticInstancesList from './AnalyticInstancesList.vue'

defineProps<{
  installedPlugins: PluginManifest[]
  instances: AnalyticInstance[]
  initialPluginId?: string | null
}>()

const emit = defineEmits<{
  (e: 'save', instance: AnalyticInstance): void
  (e: 'toggle', id: string): void
  (e: 'delete', id: string): void
  (e: 'explore', pluginId: string): void
  (e: 'goToCatalog'): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <div v-if="installedPlugins.length === 0" class="vms-card" style="padding: 2.5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem;">
      <span class="vms-text-sm vms-text-dim">// NENHUM ANALÍTICO INSTALADO NO MOMENTO</span>
      <span class="vms-text-xs vms-text-dim">Instale um analítico no Catálogo para poder criar instâncias vinculadas a câmeras.</span>
      <button class="vms-btn vms-btn-primary" @click="emit('goToCatalog')">
        IR PARA O CATÁLOGO
      </button>
    </div>

    <div v-else class="vms-analytic-creator-layout">
      <AnalyticInstanceForm
        :installed-plugins="installedPlugins"
        :initial-plugin-id="initialPluginId"
        @save="emit('save', $event)"
      />

      <AnalyticInstancesList
        :instances="instances"
        @toggle="emit('toggle', $event)"
        @delete="emit('delete', $event)"
        @explore="emit('explore', $event)"
      />
    </div>
  </div>
</template>
