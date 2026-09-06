<script setup lang="ts">
import { computed } from 'vue'
import { useMarketplace } from '../../composables/useMarketplace'
import { useDesktopAnalytics } from '../../composables/useDesktopAnalytics'
import AnalyticBreadcrumb from '../../components/admin/marketplace/AnalyticBreadcrumb.vue'
import AnalyticFolderCard from '../../components/admin/marketplace/AnalyticFolderCard.vue'
import AnalyticAppCard from '../../components/admin/marketplace/AnalyticAppCard.vue'
import AnalyticInspectorSplitView from '../../components/admin/marketplace/AnalyticInspectorSplitView.vue'
import CreateAnalyticFolderModal from '../../components/admin/marketplace/CreateAnalyticFolderModal.vue'
import AnalyticWizardModal from '../../components/admin/marketplace/AnalyticWizardModal.vue'

const props = defineProps<{ pluginId: string }>()
const { toast, getPluginById } = useMarketplace()
const plugin = computed(() => getPluginById(props.pluginId))

const {
  searchQuery, currentFolderId, currentFolder, selectedInstance, isFolderModalOpen, isWizardOpen,
  totalInstances, displayedFolders, displayedInstances, pluginFolders, handleDragStart,
  handleDropOnFolder, handleCreateFolder, handleCreateInstance, handleSaveInstance, handleDeleteInstance
} = useDesktopAnalytics(props.pluginId)
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Top Bar -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">{{ plugin?.name.toUpperCase() }} // GESTÃO DE ANALÍTICOS</h3>
        <span v-if="selectedInstance" class="vms-text-mono vms-text-2xs vms-text-dim">INSPEÇÃO // {{ selectedInstance.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">ORGANIZAÇÃO EM PASTAS // GERENCIAMENTO DE INSTÂNCIAS</span>
      </div>
      <div v-if="!selectedInstance" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar analíticos..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Pasta" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Novo Analítico" @click="isWizardOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></button>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="vms-toast">
      <div v-if="toast" class="vms-toast-notification" title="Clique para fechar" @click="toast = null"><span>{{ toast }}</span></div>
    </Transition>

    <AnalyticBreadcrumb
      :current-folder="currentFolder" :selected-instance="selectedInstance" :total-folders="pluginFolders.length"
      :total-instances="totalInstances" :folder-instances-count="currentFolder?.instances.length"
      @back="selectedInstance ? (selectedInstance = null) : (currentFolderId = null)"
      @navigate-root="currentFolderId = null; selectedInstance = null" @save="selectedInstance ? handleSaveInstance(selectedInstance) : undefined"
    />

    <!-- Desktop Canvas Container -->
    <div class="vms-desktop-container">
      <AnalyticInspectorSplitView v-if="selectedInstance" :instance="selectedInstance" @saved="handleSaveInstance" @delete="handleDeleteInstance" @close="selectedInstance = null" />

      <div v-else class="vms-desktop-canvas">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 1.25rem;">
          <div class="vms-desktop-grid">
            <AnalyticFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-instance="handleDropOnFolder" />
          </div>
          <div v-if="displayedInstances.length > 0" class="vms-flex-col" style="gap: 0.5rem; border-top: 1px solid var(--vms-border); padding-top: 1rem;">
            <div class="vms-desktop-section-title">// ANALÍTICOS NA RAIZ (ARRASTE PARA UMA PASTA)</div>
            <div class="vms-desktop-grid">
              <AnalyticAppCard v-for="i in displayedInstances" :key="i.id" :instance="i" @dragstart="handleDragStart" @select="(inst) => selectedInstance = inst" />
            </div>
          </div>
        </div>

        <div v-else class="vms-desktop-grid">
          <div v-if="displayedInstances.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// PASTA VAZIA (CLIQUE EM [+] PARA CRIAR ANALÍTICO)</div>
          <AnalyticAppCard v-for="i in displayedInstances" :key="i.id" :instance="i" @dragstart="handleDragStart" @select="(inst) => selectedInstance = inst" />
        </div>
      </div>
    </div>

    <!-- Modais -->
    <CreateAnalyticFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @create="handleCreateFolder" />
    <AnalyticWizardModal v-if="plugin && isWizardOpen" :is-open="isWizardOpen" :plugin="plugin" :folders="pluginFolders" :current-folder-id="currentFolderId" @close="isWizardOpen = false" @save="handleCreateInstance($event.inst, $event.targetFolderId)" />
  </div>
</template>
