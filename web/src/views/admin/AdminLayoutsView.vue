<script setup lang="ts">
import { useDesktopLayouts } from '../../composables/useDesktopLayouts'
import LayoutBreadcrumb from '../../components/admin/layouts/LayoutBreadcrumb.vue'
import LayoutFolderCard from '../../components/admin/layouts/LayoutFolderCard.vue'
import LayoutAppCard from '../../components/admin/layouts/LayoutAppCard.vue'
import LayoutInspectorSplitView from '../../components/admin/layouts/LayoutInspectorSplitView.vue'
import TreeContextMenu, { type ContextMenuTarget } from '../../components/admin/TreeContextMenu.vue'
import CreateLayoutFolderModal from '../../components/admin/layouts/CreateLayoutFolderModal.vue'
import LayoutWizardModal from '../../components/admin/layouts/LayoutWizardModal.vue'
import ConfirmDeleteFolderModal from '../../components/admin/ConfirmDeleteFolderModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedLayout, notification, contextMenu,
  isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalLayouts, displayedFolders,
  displayedLayouts, openContextMenu, handleDragStart, handleDropOnFolder, moveLayoutToFolder,
  requestDeleteFolder, confirmDeleteFolder, deleteLayoutById, handleSaveFolder, handleSaveLayout, showNotification
} = useDesktopLayouts()

const handleContextAction = (action: string, target: ContextMenuTarget, extra?: any) => {
  if (action === 'open-folder' && target.id) currentFolderId.value = target.id
  else if (action === 'create-folder') isFolderModalOpen.value = true
  else if (action === 'create-stream') { if (target.type === 'folder' && target.id) currentFolderId.value = target.id; isWizardOpen.value = true }
  else if (action === 'inspect-stream' && target.id) {
    let l = displayedLayouts.value.find(i => i.id === target.id)
    if (!l) { for (const f of folders.value) { l = f.layouts.find(i => i.id === target.id); if (l) break } }
    if (l) selectedLayout.value = l
  }
  else if (action === 'delete-folder' && target.id) requestDeleteFolder(target.id)
  else if (action === 'delete-stream' && target.id) deleteLayoutById(target.id)
  else if (action === 'move-stream' && target.id) moveLayoutToFolder(target.id, extra)
  else if (action === 'test-stream' && target.id) showNotification(`[LAYOUT TEST] Grade "${target.name || ''}" // VALIDA`)
}

const handleToggleLock = () => {
  if (!selectedLayout.value) return
  selectedLayout.value.is_locked = !selectedLayout.value.is_locked
  showNotification(selectedLayout.value.is_locked ? `[TRAVADO COM CADEADO] Grade "${selectedLayout.value.name}" bloqueada` : `[DESTRAVADO] Grade "${selectedLayout.value.name}" liberada`)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE LAYOUTS DE GRADE</h3>
        <span v-if="selectedLayout" class="vms-text-mono vms-text-2xs vms-text-dim">INSPECAO // {{ selectedLayout.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">ORGANIZACAO EM PASTAS // EMPRESAS & CLIENTES FINAIS</span>
      </div>
      <div v-if="!selectedLayout" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar layouts..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Pasta (Empresa/Cliente)" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Novo Layout" @click="isWizardOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M9 3v18"/></svg></button>
      </div>
    </div>

    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" @click="notification = null">
        <svg v-if="notification.includes('[TRAVADO')" width="14" height="14" viewBox="0 0 24 24" fill="#ff5e3a"><path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/></svg>
        <svg v-else-if="notification.includes('[DESTRAVADO')" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#00ff9d" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 9.9-1"/></svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#00f0ff" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>
        <span>{{ notification }}</span>
      </div>
    </Transition>

    <LayoutBreadcrumb :current-folder="currentFolder" :selected-layout="selectedLayout" :total-folders="folders.length" :total-layouts="totalLayouts" :current-folder-layouts-count="currentFolder?.layouts.length" @back="selectedLayout ? (selectedLayout = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedLayout = null" @navigate-folder="selectedLayout = null" @save="selectedLayout ? handleSaveLayout(selectedLayout) : undefined" @toggle-lock="handleToggleLock" />

    <div class="vms-desktop-container">
      <LayoutInspectorSplitView v-if="selectedLayout" :layout="selectedLayout" @saved="(msg) => showNotification(msg)" />
      <div v-else class="vms-desktop-canvas" @contextmenu.prevent="openContextMenu($event, { type: 'canvas' })">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <LayoutFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-layout="handleDropOnFolder" @context="(ev, fold) => openContextMenu(ev, { type: 'folder', id: fold.id, name: fold.name })" />
          </div>
          <div v-if="displayedLayouts.length > 0" class="vms-flex-col" style="gap: 0.4rem; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <div class="vms-desktop-section-title">// LAYOUTS GLOBAIS NA RAIZ (ARRASTE PARA UMA EMPRESA/CLIENTE)</div>
            <div class="vms-desktop-grid">
              <LayoutAppCard v-for="l in displayedLayouts" :key="l.id" :layout="l" :is-selected="false" @dragstart="handleDragStart" @select="(lay) => selectedLayout = lay" @context="(ev, lay) => openContextMenu(ev, { type: 'stream', id: lay.id, name: lay.name, currentFolderId: currentFolderId })" />
            </div>
          </div>
        </div>
        <div v-else class="vms-desktop-grid">
          <div v-if="displayedLayouts.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// PASTA VAZIA (CLIQUE EM [+] OU BOTAO DIREITO PARA CADASTRAR LAYOUT)</div>
          <LayoutAppCard v-for="l in displayedLayouts" :key="l.id" :layout="l" :is-selected="false" @dragstart="handleDragStart" @select="(lay) => selectedLayout = lay" @context="(ev, lay) => openContextMenu(ev, { type: 'stream', id: lay.id, name: lay.name, currentFolderId: currentFolderId })" />
        </div>
      </div>
    </div>

    <TreeContextMenu :is-open="contextMenu.isOpen" :x="contextMenu.x" :y="contextMenu.y" :target="contextMenu.target" :folders="folders as any" @close="contextMenu.isOpen = false" @action="handleContextAction" />
    <CreateLayoutFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @save="handleSaveFolder" />
    <LayoutWizardModal :is-open="isWizardOpen" :target-folder-id="currentFolderId || undefined" :folders="folders" @close="isWizardOpen = false" @save="handleSaveLayout" />
    <ConfirmDeleteFolderModal :is-open="isConfirmDeleteOpen" :folder-name="folderToDelete?.name || ''" :item-count="folderToDelete?.itemCount || 0" item-type="layout(s)" @close="isConfirmDeleteOpen = false" @confirm="confirmDeleteFolder" />
  </div>
</template>
