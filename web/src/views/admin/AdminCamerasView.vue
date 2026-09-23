<script setup lang="ts">
import { useDesktopTree } from '../../composables/useDesktopTree'
import { useI18n } from '../../composables/useI18n'
import { useAdminCameraActions } from '../../composables/useAdminCameraActions'
import DesktopBreadcrumb from '../../components/admin/desktop/DesktopBreadcrumb.vue'
import DesktopFolderCard from '../../components/admin/desktop/DesktopFolderCard.vue'
import DesktopStreamApp from '../../components/admin/desktop/DesktopStreamApp.vue'
import StreamInspectorSplitView from '../../components/admin/desktop/StreamInspectorSplitView.vue'
import OnvifDiscoverySection from '../../components/admin/desktop/OnvifDiscoverySection.vue'
import TreeContextMenu from '../../components/admin/TreeContextMenu.vue'
import CreateFolderModal from '../../components/admin/CreateFolderModal.vue'
import StreamWizardModal from '../../components/admin/StreamWizardModal.vue'
import ConfirmDeleteFolderModal from '../../components/admin/ConfirmDeleteFolderModal.vue'
import ConfirmDeleteStreamModal from '../../components/admin/ConfirmDeleteStreamModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedStream, notification, contextMenu,
  isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalStreams, displayedFolders,
  displayedStreams, openContextMenu, handleDragStart, handleDropOnFolder, moveStreamToFolder,
  requestDeleteFolder, confirmDeleteFolder, deleteStreamById, handleSaveFolder, handleSaveStream, showNotification
} = useDesktopTree()

const { t } = useI18n()

const {
  wizardInitialData, fileInputRef, isConfirmDeleteStreamOpen, streamToDelete,
  handleExportStreams, triggerFileInput, handleImportFile, handleOpenNewWizard,
  onSaveStream, confirmDeleteStream, handleContextAction
} = useAdminCameraActions(
  folders, displayedStreams, currentFolderId, selectedStream,
  isFolderModalOpen, isWizardOpen, showNotification, handleSaveStream,
  deleteStreamById, requestDeleteFolder, moveStreamToFolder
)
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Top Bar -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">{{ t('desktop_title') }}</h3>
        <span v-if="selectedStream" class="vms-text-mono vms-text-2xs vms-text-dim">INSPEÇÃO // {{ selectedStream.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">{{ t('desktop_sub') }}</span>
      </div>
      <div v-if="!selectedStream" class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 170px; height: 32px; font-size: 12px; padding: 4px 10px; box-sizing: border-box;" :placeholder="t('filter_placeholder')" />
        <button class="vms-btn vms-btn-secondary" style="height: 32px; padding: 0 10px; box-sizing: border-box;" :title="t('new_folder')" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-secondary" style="height: 32px; padding: 0 10px; box-sizing: border-box;" :title="t('new_stream')" @click="handleOpenNewWizard()"><span>+</span><svg width="14" height="14" viewBox="0 0 576 512" fill="#ff5e3a"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg></button>
        <button class="vms-btn vms-btn-secondary" style="height: 32px; padding: 0 10px; box-sizing: border-box;" title="Exportar Lista de Fluxos (Planilha CSV / Excel)" @click="handleExportStreams">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
        </button>
        <button class="vms-btn vms-btn-secondary" style="height: 32px; padding: 0 10px; box-sizing: border-box;" title="Importar Lista de Fluxos (CSV, TXT ou JSON)" @click="triggerFileInput">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        </button>
        <input ref="fileInputRef" type="file" accept=".csv, .json, .txt" style="display: none;" @change="handleImportFile" />
      </div>
    </div>

    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" @click="notification = null"><span>{{ notification }}</span></div>
    </Transition>

    <!-- Bloco Exclusivo ONVIF WS-Discovery -->
    <OnvifDiscoverySection v-if="!selectedStream" @import-onvif="handleOpenNewWizard" />

    <DesktopBreadcrumb :current-folder="currentFolder" :selected-stream="selectedStream" :total-folders="folders.length" :total-streams="totalStreams" :current-folder-streams-count="currentFolder?.streams.length" @back="selectedStream ? (selectedStream = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedStream = null" @navigate-folder="selectedStream = null" />

    <div class="vms-desktop-container">
      <StreamInspectorSplitView v-if="selectedStream" :stream="selectedStream" @saved="(msg) => showNotification(msg)" />
      <div v-else class="vms-desktop-canvas" @contextmenu.prevent="openContextMenu($event, { type: 'canvas' })">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <DesktopFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-stream="handleDropOnFolder" @context="(ev, fold) => openContextMenu(ev, { type: 'folder', id: fold.id, name: fold.name })" />
          </div>
          <div v-if="displayedStreams.length > 0" class="vms-flex-col" style="gap: 0.4rem; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <div class="vms-desktop-section-title">// FLUXOS NA RAIZ (ARRASTE PARA UMA PASTA)</div>
            <div class="vms-desktop-grid">
              <DesktopStreamApp v-for="s in displayedStreams" :key="s.id" :stream="s" :is-selected="false" @dragstart="handleDragStart" @select="(stream) => selectedStream = stream" @context="(ev, stream) => openContextMenu(ev, { type: 'stream', id: stream.id, name: stream.name, currentFolderId: currentFolderId })" />
            </div>
          </div>
        </div>
        <div v-else class="vms-desktop-grid">
          <div v-if="displayedStreams.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// PASTA VAZIA (CLIQUE EM [+] OU BOTÃO DIREITO PARA CADASTRAR)</div>
          <DesktopStreamApp v-for="s in displayedStreams" :key="s.id" :stream="s" :is-selected="false" @dragstart="handleDragStart" @select="(stream) => selectedStream = stream" @context="(ev, stream) => openContextMenu(ev, { type: 'stream', id: stream.id, name: stream.name, currentFolderId: currentFolderId })" />
        </div>
      </div>
    </div>

    <TreeContextMenu :is-open="contextMenu.isOpen" :x="contextMenu.x" :y="contextMenu.y" :target="contextMenu.target" :folders="folders" @close="contextMenu.isOpen = false" @action="handleContextAction" />
    <CreateFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @save="handleSaveFolder" />
    <StreamWizardModal :is-open="isWizardOpen" :initial-data="wizardInitialData" :target-folder-id="currentFolderId || undefined" :folders="folders" @close="isWizardOpen = false" @save="onSaveStream" />
    <ConfirmDeleteFolderModal :is-open="isConfirmDeleteOpen" :folder-name="folderToDelete?.name || ''" :item-count="folderToDelete?.itemCount || 0" item-type="fluxo(s)" @close="isConfirmDeleteOpen = false" @confirm="confirmDeleteFolder" />
    <ConfirmDeleteStreamModal :is-open="isConfirmDeleteStreamOpen" :stream-name="streamToDelete?.name || ''" :stream-id="streamToDelete?.id || ''" @close="isConfirmDeleteStreamOpen = false" @confirm="confirmDeleteStream" />
  </div>
</template>
