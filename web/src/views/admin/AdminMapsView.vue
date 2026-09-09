<script setup lang="ts">
import { useDesktopMaps } from '../../composables/useDesktopMaps'
import MapBreadcrumb from '../../components/admin/maps/MapBreadcrumb.vue'
import MapFolderCard from '../../components/admin/maps/MapFolderCard.vue'
import MapItemCard from '../../components/admin/maps/MapItemCard.vue'
import MapEditorView from '../../components/admin/maps/MapEditorView.vue'
import MapNewFolderModal from '../../components/admin/maps/MapNewFolderModal.vue'
import MapNewModal from '../../components/admin/maps/MapNewModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedMap, notification,
  isFolderModalOpen, isNewMapModalOpen, totalMaps, displayedFolders, displayedMaps,
  showNotification, handleDragStart, handleDropOnFolder, handleCreateFolder,
  handleCreateMap, handleSaveMap, handleToggleLock
} = useDesktopMaps()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CADASTRO DE PLANTAS BAIXAS & MAPAS</h3>
        <span v-if="selectedMap" class="vms-text-mono vms-text-2xs vms-text-dim">EDITOR // {{ selectedMap.name }}</span>
        <span v-else class="vms-text-mono vms-text-2xs vms-text-dim">ORGANIZACAO EM PASTAS // MAPAS COLORIDOS OPENSOURCE & PLANTAS</span>
      </div>
      <div v-if="!selectedMap" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar mapas..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Pasta (Empresa/Cliente)" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Novo Mapa / Planta" @click="isNewMapModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"/><line x1="8" y1="2" x2="8" y2="18"/><line x1="16" y1="6" x2="16" y2="22"/></svg></button>
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

    <MapBreadcrumb :current-folder="currentFolder" :selected-map="selectedMap" :total-folders="folders.length" :total-maps="totalMaps" :current-folder-maps-count="currentFolder?.maps.length" @back="selectedMap ? (selectedMap = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedMap = null" @navigate-folder="selectedMap = null" @save="selectedMap ? handleSaveMap(selectedMap) : undefined" @toggle-lock="handleToggleLock" />

    <div class="vms-desktop-container">
      <MapEditorView v-if="selectedMap" :map-item="selectedMap" @notify="(msg) => showNotification(msg)" />
      <div v-else class="vms-desktop-canvas">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <MapFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-map="handleDropOnFolder" />
          </div>
          <div v-if="displayedMaps.length > 0" class="vms-desktop-grid" style="border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <MapItemCard v-for="m in displayedMaps" :key="m.id" :map-item="m" @open="(item) => selectedMap = item" @drag-start="handleDragStart" />
          </div>
        </div>
        <div v-else class="vms-desktop-grid">
          <div v-if="displayedMaps.length === 0" class="vms-text-dim vms-text-sm" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">Nenhum mapa nesta pasta.</div>
          <MapItemCard v-for="m in displayedMaps" :key="m.id" :map-item="m" @open="(item) => selectedMap = item" @drag-start="handleDragStart" />
        </div>
      </div>
    </div>

    <!-- Modais -->
    <MapNewFolderModal v-if="isFolderModalOpen" @close="isFolderModalOpen = false" @create="handleCreateFolder" />
    <MapNewModal v-if="isNewMapModalOpen" :default-scope="currentFolder?.name" @close="isNewMapModalOpen = false" @create="handleCreateMap" />
  </div>
</template>
