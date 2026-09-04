<script setup lang="ts">
import { ref } from "vue"
import AdminSidebar, { type AdminPageId } from "./AdminSidebar.vue"
import AdminCamerasView from "./AdminCamerasView.vue"
import AdminLayoutsView from "./AdminLayoutsView.vue"
import AdminMapsView from "./AdminMapsView.vue"
import AdminCarouselsView from "./AdminCarouselsView.vue"
import AdminUsersView from "./AdminUsersView.vue"
import WorkflowListView from "../workflows/WorkflowListView.vue"
import StorageDisksView from "../storage/StorageDisksView.vue"

const activePage = ref<AdminPageId>("cameras")
</script>

<template>
  <div class="vms-admin-container" style="display: flex; flex-direction: row; flex: 1; height: 100%; overflow: hidden; background: var(--vms-neu-bg);">
    <!-- Cyberpunk Neumorphic Sidebar for Admin -->
    <AdminSidebar :activePage="activePage" @selectPage="activePage = $event" />

    <!-- Full-Width Page View Host -->
    <main class="vms-admin-content" style="flex: 1; min-width: 0; padding: 1.5rem 2rem; overflow-y: auto;">
      <AdminCamerasView v-if="activePage === 'cameras'" />
      <AdminLayoutsView v-else-if="activePage === 'layouts'" />
      <AdminMapsView v-else-if="activePage === 'maps'" />
      <AdminCarouselsView v-else-if="activePage === 'carousels'" />
      <AdminUsersView v-else-if="activePage === 'users'" />
      <WorkflowListView v-else-if="activePage === 'workflows'" />
      <StorageDisksView v-else-if="activePage === 'storage'" />
    </main>
  </div>
</template>
