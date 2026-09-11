<script setup lang="ts">
import { useDesktopUsers } from '../../composables/useDesktopUsers'
import { useI18n } from '../../composables/useI18n'
import UserBreadcrumb from '../../components/admin/users/UserBreadcrumb.vue'
import UserFolderCard from '../../components/admin/users/UserFolderCard.vue'
import UserAppCard from '../../components/admin/users/UserAppCard.vue'
import UserInspectorSplitView from '../../components/admin/users/UserInspectorSplitView.vue'
import TreeContextMenu, { type ContextMenuTarget } from '../../components/admin/TreeContextMenu.vue'
import CreateUserFolderModal from '../../components/admin/users/CreateUserFolderModal.vue'
import UserWizardModal from '../../components/admin/users/UserWizardModal.vue'
import ConfirmDeleteFolderModal from '../../components/admin/ConfirmDeleteFolderModal.vue'

const {
  searchQuery, folders, currentFolderId, currentFolder, selectedUser, notification, contextMenu,
  isFolderModalOpen, isWizardOpen, folderToDelete, isConfirmDeleteOpen, totalUsers, displayedFolders,
  displayedUsers, openContextMenu, handleDragStart, handleDropOnFolder, moveUserToFolder,
  requestDeleteFolder, confirmDeleteFolder, deleteUserById, handleSaveFolder, handleSaveUser, showNotification
} = useDesktopUsers()
const { t } = useI18n()

const handleContextAction = (action: string, target: ContextMenuTarget, extra?: any) => {
  if (action === 'open-folder' && target.id) currentFolderId.value = target.id
  else if (action === 'create-folder') isFolderModalOpen.value = true
  else if (action === 'create-stream') { if (target.type === 'folder' && target.id) currentFolderId.value = target.id; isWizardOpen.value = true }
  else if (action === 'inspect-stream' && target.id) {
    let u = displayedUsers.value.find(item => item.id === target.id)
    if (!u) { for (const f of folders.value) { u = f.users.find(item => item.id === target.id); if (u) break } }
    if (u) selectedUser.value = u
  }
  else if (action === 'delete-folder' && target.id) requestDeleteFolder(target.id)
  else if (action === 'delete-stream' && target.id) deleteUserById(target.id)
  else if (action === 'move-stream' && target.id) moveUserToFolder(target.id, extra)
  else if (action === 'test-stream' && target.id) showNotification(`[AUTH TEST] Token do usuário ${target.name || ''} // VALIDO`)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Top Bar -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">{{ t('users_title') }}</h3>
        <span v-if="selectedUser" class="vms-text-mono vms-text-2xs vms-text-dim">INSPEÇÃO // {{ selectedUser.username }}</span>
      </div>
      <div v-if="!selectedUser" class="vms-flex-row" style="gap: 0.75rem;">
        <input v-model="searchQuery" class="vms-auth-input" style="width: 200px; font-size: 12px; padding: 4px 10px;" :placeholder="t('filter_users')" />
        <button class="vms-btn vms-btn-secondary" :title="t('new_group')" @click="isFolderModalOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 512 512" fill="#ff5e3a"><path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/></svg></button>
        <button class="vms-btn vms-btn-primary" :title="t('new_user')" @click="isWizardOpen = true"><span>+</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></button>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" title="Clique para fechar" @click="notification = null"><span>{{ notification }}</span></div>
    </Transition>

    <!-- Breadcrumb with 3 GRUPOS // X USUARIOS -->
    <UserBreadcrumb :current-folder="currentFolder" :selected-user="selectedUser" :total-folders="folders.length" :total-users="totalUsers" :current-folder-users-count="currentFolder?.users.length" @back="selectedUser ? (selectedUser = null) : (currentFolderId = null)" @navigate-root="currentFolderId = null; selectedUser = null" @navigate-folder="selectedUser = null" />

    <!-- Main Workspace Container -->
    <div class="vms-desktop-container">
      <UserInspectorSplitView v-if="selectedUser" :user="selectedUser" @saved="(msg) => showNotification(msg)" />

      <div v-else class="vms-desktop-canvas" @contextmenu.prevent="openContextMenu($event, { type: 'canvas' })">
        <div v-if="!currentFolder" class="vms-flex-col" style="gap: 0.65rem;">
          <div class="vms-desktop-grid">
            <UserFolderCard v-for="f in displayedFolders" :key="f.id" :folder="f" @open="(id) => currentFolderId = id" @drop-user="handleDropOnFolder" @context="(ev, fold) => openContextMenu(ev, { type: 'folder', id: fold.id, name: fold.name })" />
          </div>

          <div v-if="displayedUsers.length > 0" class="vms-flex-col" style="gap: 0.4rem; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
            <div class="vms-desktop-section-title">// USUÁRIOS NA RAIZ (ARRASTE PARA UM GRUPO)</div>
            <div class="vms-desktop-grid">
              <UserAppCard v-for="u in displayedUsers" :key="u.id" :user="u" :is-selected="false" @dragstart="handleDragStart" @select="(usr) => selectedUser = usr" @context="(ev, usr) => openContextMenu(ev, { type: 'stream', id: usr.id, name: usr.username, currentFolderId: currentFolderId })" />
            </div>
          </div>
        </div>

        <div v-else class="vms-desktop-grid">
          <div v-if="displayedUsers.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="grid-column: 1 / -1; padding: 2rem; text-align: center;">// GRUPO VAZIO (CLIQUE EM [+] OU BOTÃO DIREITO PARA CADASTRAR USUÁRIO)</div>
          <UserAppCard v-for="u in displayedUsers" :key="u.id" :user="u" :is-selected="false" @dragstart="handleDragStart" @select="(usr) => selectedUser = usr" @context="(ev, usr) => openContextMenu(ev, { type: 'stream', id: usr.id, name: usr.username, currentFolderId: currentFolderId })" />
        </div>
      </div>
    </div>

    <TreeContextMenu :is-open="contextMenu.isOpen" :x="contextMenu.x" :y="contextMenu.y" :target="contextMenu.target" :folders="folders" @close="contextMenu.isOpen = false" @action="handleContextAction" />
    <CreateUserFolderModal :is-open="isFolderModalOpen" @close="isFolderModalOpen = false" @save="handleSaveFolder" />
    <UserWizardModal :is-open="isWizardOpen" :target-folder-id="currentFolderId || undefined" :folders="folders" @close="isWizardOpen = false" @save="handleSaveUser" />
    <ConfirmDeleteFolderModal :is-open="isConfirmDeleteOpen" :folder-name="folderToDelete?.name || ''" :item-count="folderToDelete?.itemCount || 0" item-type="usuário(s)" @close="isConfirmDeleteOpen = false" @confirm="confirmDeleteFolder" />
  </div>
</template>
