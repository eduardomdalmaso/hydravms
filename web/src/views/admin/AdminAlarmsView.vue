<script setup lang="ts">
import { useDesktopAlarms } from '../../composables/useDesktopAlarms'
import AlarmBreadcrumb from '../../components/admin/alarms/AlarmBreadcrumb.vue'
import AlarmFolderCard from '../../components/admin/alarms/AlarmFolderCard.vue'
import AlarmAppCard from '../../components/admin/alarms/AlarmAppCard.vue'
import AlarmInspectorSplitView from '../../components/admin/alarms/AlarmInspectorSplitView.vue'
import TreeContextMenu, { type ContextMenuTarget } from '../../components/admin/TreeContextMenu.vue'
import CreateAlarmFolderModal from '../../components/admin/alarms/CreateAlarmFolderModal.vue'
import AlarmWizardModal from '../../components/admin/alarms/AlarmWizardModal.vue'
import ConfirmDeleteFolderModal from '../../components/admin/ConfirmDeleteFolderModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedAlarm, notification, contextMenu,
  isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalAlarms, displayedFolders,
  displayedAlarms, openContextMenu, handleDragStart, handleDropOnFolder, moveAlarmToFolder,
  requestDeleteFolder, confirmDeleteFolder, deleteAlarmById, handleSaveFolder, handleSaveAlarm, showNotification
} = useDesktopAlarms()

const handleContextAction = (action: string, target: ContextMenuTarget, extra?: any) => {
  if (action === 'open-folder' && target.id) currentFolderId.value = target.id
  else if (action === 'create-folder') isFolderModalOpen.value = true
  else if (action === 'create-stream') { if (target.type === 'folder' && target.id) currentFolderId.value = target.id; isWizardOpen.value = true }
  else if (action === 'inspect-stream' && target.id) {
    let a = displayedAlarms.value.find(item => item.id === target.id)
    if (!a) { for (const f of folders.value) { a = f.alarms.find(item => item.id === target.id); if (a) break } }
    if (a) selectedAlarm.value = a
  }
  else if (action === 'delete-folder' && target.id) requestDeleteFolder(target.id)
  else if (action === 'delete-stream' && target.id) deleteAlarmById(target.id)
  else if (action === 'move-stream' && target.id) moveAlarmToFolder(target.id, extra)
  else if (action === 'test-stream' && target.id) showNotification(`[TRIGGER TEST] Sensor ${target.name || ''} // DISPARO OK`)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Top Bar -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CONFIGURAÇÃO DE ALARMES & SENSORES</h3>
        <span v-if="selectedAlarm" class="vms-text-mono vms-text-2xs vms-text-dim">INSPEÇÃO // {{ selectedAlarm.name }}</span>
      </div>
      <div v-if="!selectedAlarm" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" placeholder="Filtrar sensores..." />
        <button class="vms-btn vms-btn-secondary" title="Nova Zona / Pasta" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" title="Novo Sensor de Alarme" @click="isWizardOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" /><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" /></svg></button>
      </div>
    </div>

    <!-- Floating Toast Notification -->
    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" title="Clique para fechar" @click="notification = null"><span>{{ notification }}</span></div>
    </Transition>

    <!-- Breadcrumb with 3 ZONAS // X SENSORES -->
    <AlarmBreadcrumb :current-folder="currentFolder" :selected-alarm="selectedAlarm" :total-folders="folders.length" :total-alarms="totalAlarms" :current-folder-alarms-count="currentFolder?.alarms.length" @back="selectedAlarm ? (selectedAlarm = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedAlarm = null" @navigate-folder="selectedAlarm = null" />

    <!-- Main Workspace Container -->
    <div class="vms-desktop-container">
      <AlarmInspectorSplitView v-if="selectedAlarm" :alarm="selectedAlarm" @saved="(msg) => showNotification(msg)" />

      <div v-else class="vms-desktop-canvas" @contextmenu.prevent="openContextMenu($event, { type: 'canvas' })">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <AlarmFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-alarm="handleDropOnFolder" @context="(ev, fold) => openContextMenu(ev, { type: 'folder', id: fold.id, name: fold.name })" />
          </div>

          <div v-if="displayedAlarms.length > 0" class="vms-flex-col" style="gap: 0.4rem; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <div class="vms-desktop-section-title">// SENSORES NA RAIZ (ARRASTE PARA UMA ZONA)</div>
            <div class="vms-desktop-grid">
              <AlarmAppCard v-for="a in displayedAlarms" :key="a.id" :alarm="a" :is-selected="false" @dragstart="handleDragStart" @select="(alarm) => selectedAlarm = alarm" @context="(ev, alarm) => openContextMenu(ev, { type: 'stream', id: alarm.id, name: alarm.name, currentFolderId: currentFolderId })" />
            </div>
          </div>
        </div>

        <div v-else class="vms-desktop-grid">
          <div v-if="displayedAlarms.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// ZONA VAZIA (CLIQUE EM [+] OU BOTÃO DIREITO PARA CADASTRAR SENSOR)</div>
          <AlarmAppCard v-for="a in displayedAlarms" :key="a.id" :alarm="a" :is-selected="false" @dragstart="handleDragStart" @select="(alarm) => selectedAlarm = alarm" @context="(ev, alarm) => openContextMenu(ev, { type: 'stream', id: alarm.id, name: alarm.name, currentFolderId: currentFolderId })" />
        </div>
      </div>
    </div>

    <TreeContextMenu :is-open="contextMenu.isOpen" :x="contextMenu.x" :y="contextMenu.y" :target="contextMenu.target" :folders="folders" @close="contextMenu.isOpen = false" @action="handleContextAction" />
    <CreateAlarmFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @save="handleSaveFolder" />
    <AlarmWizardModal :is-open="isWizardOpen" :target-folder-id="currentFolderId || undefined" :folders="folders" @close="isWizardOpen = false" @save="handleSaveAlarm" />
    <ConfirmDeleteFolderModal :is-open="isConfirmDeleteOpen" :folder-name="folderToDelete?.name || ''" :item-count="folderToDelete?.itemCount || 0" item-type="sensor(es)" @close="isConfirmDeleteOpen = false" @confirm="confirmDeleteFolder" />
  </div>
</template>
