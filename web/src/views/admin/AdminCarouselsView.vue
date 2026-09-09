<script setup lang="ts">
import { useDesktopRondas } from '../../composables/useDesktopRondas'
import RondaBreadcrumb from '../../components/admin/rondas/RondaBreadcrumb.vue'
import RondaFolderCard from '../../components/admin/rondas/RondaFolderCard.vue'
import RondaAppCard from '../../components/admin/rondas/RondaAppCard.vue'
import RondaInspectorSplitView from '../../components/admin/rondas/RondaInspectorSplitView.vue'
import TreeContextMenu, { type ContextMenuTarget } from '../../components/admin/TreeContextMenu.vue'
import CreateRondaFolderModal from '../../components/admin/rondas/CreateRondaFolderModal.vue'
import RondaWizardModal from '../../components/admin/rondas/RondaWizardModal.vue'
import ConfirmDeleteFolderModal from '../../components/admin/ConfirmDeleteFolderModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedRonda, notification, contextMenu,
  isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalRondas, displayedFolders,
  displayedRondas, openContextMenu, handleDragStart, handleDropOnFolder, moveRondaToFolder,
  requestDeleteFolder, confirmDeleteFolder, deleteRondaById, handleSaveFolder, handleSaveRonda, showNotification
} = useDesktopRondas()

const handleContextAction = (action: string, target: ContextMenuTarget, extra?: any) => {
  if (action === 'open-folder' && target.id) currentFolderId.value = target.id
  else if (action === 'create-folder') isFolderModalOpen.value = true
  else if (action === 'create-stream') { if (target.type === 'folder' && target.id) currentFolderId.value = target.id; isWizardOpen.value = true }
  else if (action === 'inspect-stream' && target.id) {
    let r = displayedRondas.value.find(i => i.id === target.id)
    if (!r) { for (const f of folders.value) { r = f.rondas.find(i => i.id === target.id); if (r) break } }
    if (r) selectedRonda.value = r
  }
  else if (action === 'delete-folder' && target.id) requestDeleteFolder(target.id)
  else if (action === 'delete-stream' && target.id) deleteRondaById(target.id)
  else if (action === 'move-stream' && target.id) moveRondaToFolder(target.id, extra)
  else if (action === 'test-stream' && target.id) showNotification(`[RONDA TEST] "${target.name || ''}" // SEQUENCIAMENTO OK`)
}

const handleToggleLock = () => {
  if (!selectedRonda.value) return
  selectedRonda.value.is_locked = !selectedRonda.value.is_locked
  showNotification(selectedRonda.value.is_locked ? `[TRAVADO COM CADEADO] Ronda "${selectedRonda.value.name}" bloqueada` : `[DESTRAVADO] Ronda "${selectedRonda.value.name}" liberada`)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE RONDAS & CARROSSEL</h3>
        <span v-if="selectedRonda" class="vms-text-mono vms-text-2xs vms-text-dim">INSPECAO // {{ selectedRonda.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">ORGANIZACAO EM PASTAS // EMPRESAS & CLIENTES FINAIS</span>
      </div>
      <div v-if="!selectedRonda" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar rondas..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Pasta (Empresa/Cliente)" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Nova Ronda" @click="isWizardOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg></button>
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

    <RondaBreadcrumb :current-folder="currentFolder" :selected-ronda="selectedRonda" :total-folders="folders.length" :total-rondas="totalRondas" :current-folder-rondas-count="currentFolder?.rondas.length" @back="selectedRonda ? (selectedRonda = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedRonda = null" @navigate-folder="selectedRonda = null" @save="selectedRonda ? handleSaveRonda(selectedRonda) : undefined" @toggle-lock="handleToggleLock" />

    <div class="vms-desktop-container">
      <RondaInspectorSplitView v-if="selectedRonda" :ronda="selectedRonda" @saved="(msg) => showNotification(msg)" />
      <div v-else class="vms-desktop-canvas" @contextmenu.prevent="openContextMenu($event, { type: 'canvas' })">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <RondaFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-ronda="handleDropOnFolder" @context="(ev, fold) => openContextMenu(ev, { type: 'folder', id: fold.id, name: fold.name })" />
          </div>
          <div v-if="displayedRondas.length > 0" class="vms-flex-col" style="gap: 0.4rem; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <div class="vms-desktop-section-title">// RONDAS GLOBAIS NA RAIZ (ARRASTE PARA UMA EMPRESA/CLIENTE)</div>
            <div class="vms-desktop-grid">
              <RondaAppCard v-for="r in displayedRondas" :key="r.id" :ronda="r" :is-selected="false" @dragstart="handleDragStart" @select="(ron) => selectedRonda = ron" @context="(ev, ron) => openContextMenu(ev, { type: 'stream', id: ron.id, name: ron.name, currentFolderId: currentFolderId })" />
            </div>
          </div>
        </div>
        <div v-else class="vms-desktop-grid">
          <div v-if="displayedRondas.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// PASTA VAZIA (CLIQUE EM [+] OU BOTAO DIREITO PARA CADASTRAR RONDA)</div>
          <RondaAppCard v-for="r in displayedRondas" :key="r.id" :ronda="r" :is-selected="false" @dragstart="handleDragStart" @select="(ron) => selectedRonda = ron" @context="(ev, ron) => openContextMenu(ev, { type: 'stream', id: ron.id, name: ron.name, currentFolderId: currentFolderId })" />
        </div>
      </div>
    </div>

    <TreeContextMenu :is-open="contextMenu.isOpen" :x="contextMenu.x" :y="contextMenu.y" :target="contextMenu.target" :folders="folders as any" @close="contextMenu.isOpen = false" @action="handleContextAction" />
    <CreateRondaFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @save="handleSaveFolder" />
    <RondaWizardModal :is-open="isWizardOpen" :target-folder-id="currentFolderId || undefined" :folders="folders" @close="isWizardOpen = false" @save="handleSaveRonda" />
    <ConfirmDeleteFolderModal :is-open="isConfirmDeleteOpen" :folder-name="folderToDelete?.name || ''" :item-count="folderToDelete?.itemCount || 0" item-type="ronda(s)" @close="isConfirmDeleteOpen = false" @confirm="confirmDeleteFolder" />
  </div>
</template>
