<script setup lang="ts">
import { useDesktopWorkflows } from "../../composables/useDesktopWorkflows"
import WorkflowBreadcrumb from "../../components/admin/workflows/WorkflowBreadcrumb.vue"
import WorkflowFolderCard from "../../components/admin/workflows/WorkflowFolderCard.vue"
import WorkflowAppCard from "../../components/admin/workflows/WorkflowAppCard.vue"
import WorkflowFlowEditorView from "../../components/admin/workflows/WorkflowFlowEditorView.vue"
import WorkflowNewFolderModal from "../../components/admin/workflows/WorkflowNewFolderModal.vue"
import WorkflowNewModal from "../../components/admin/workflows/WorkflowNewModal.vue"

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedWorkflow, notification,
  isFolderModalOpen, isNewWorkflowModalOpen, totalWorkflows, displayedFolders, displayedWorkflows,
  showNotification, handleDragStart, handleDropOnFolder, handleCreateFolder,
  handleCreateWorkflow, handleSaveWorkflow, handleToggleLock
} = useDesktopWorkflows()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">WORKFLOWS & INTEGRAÇÕES</h3>
        <span v-if="selectedWorkflow" class="vms-text-mono vms-text-2xs vms-text-dim">EDITOR FLOW // {{ selectedWorkflow.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">ORGANIZACAO EM PASTAS // MOTOR DE INTEGRACAO VISUAL FLOW</span>
      </div>
      <div v-if="!selectedWorkflow" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar fluxos..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Pasta (Empresa/Cliente)" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Novo Aplicativo (Workflow Flow)" @click="isNewWorkflowModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg></button>
      </div>
    </div>

    <!-- Toast Notification no canto inferior direito -->
    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" @click="notification = null">
        <svg v-if="notification.includes('[TRAVADO')" width="14" height="14" viewBox="0 0 24 24" fill="#ff5e3a"><path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/></svg>
        <svg v-else-if="notification.includes('[DESTRAVADO')" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 9.9-1"/></svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>
        <span>{{ notification }}</span>
      </div>
    </Transition>

    <WorkflowBreadcrumb :current-folder="currentFolder" :selected-workflow="selectedWorkflow" :total-folders="folders.length" :total-workflows="totalWorkflows" :current-folder-workflows-count="currentFolder?.workflows.length" @back="selectedWorkflow ? (selectedWorkflow = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedWorkflow = null" @navigate-folder="selectedWorkflow = null" @save="selectedWorkflow ? handleSaveWorkflow(selectedWorkflow) : undefined" @toggle-lock="handleToggleLock" />

    <div class="vms-desktop-container">
      <WorkflowFlowEditorView v-if="selectedWorkflow" :workflow="selectedWorkflow" @notify="(msg) => showNotification(msg)" />
      <div v-else class="vms-desktop-canvas">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 1.25rem;">
          <div class="vms-desktop-grid">
            <WorkflowFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-workflow="handleDropOnFolder" />
          </div>
          <div v-if="displayedWorkflows.length > 0" class="vms-desktop-grid" style="border-top: 1px solid var(--vms-border); padding-top: 1.25rem;">
            <WorkflowAppCard v-for="w in displayedWorkflows" :key="w.id" :workflow="w" @open="(item) => selectedWorkflow = item" @drag-start="handleDragStart" />
          </div>
        </div>
        <div v-else class="vms-desktop-grid">
          <div v-if="displayedWorkflows.length === 0" class="vms-text-dim vms-text-sm" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">Nenhum workflow nesta pasta.</div>
          <WorkflowAppCard v-for="w in displayedWorkflows" :key="w.id" :workflow="w" @open="(item) => selectedWorkflow = item" @drag-start="handleDragStart" />
        </div>
      </div>
    </div>

    <!-- Modais -->
    <WorkflowNewFolderModal v-if="isFolderModalOpen" @close="isFolderModalOpen = false" @create="handleCreateFolder" />
    <WorkflowNewModal v-if="isNewWorkflowModalOpen" :folders="folders" :current-folder-id="currentFolderId" @close="isNewWorkflowModalOpen = false" @create="handleCreateWorkflow" />
  </div>
</template>
