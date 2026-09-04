<script setup lang="ts">
import { ref } from 'vue'
import AdminCamerasView from './AdminCamerasView.vue'
import AdminLayoutsView from './AdminLayoutsView.vue'
import AdminMapsView from './AdminMapsView.vue'
import AdminCarouselsView from './AdminCarouselsView.vue'
import AdminUsersView from './AdminUsersView.vue'
import WorkflowListView from '../workflows/WorkflowListView.vue'
import StorageDisksView from '../storage/StorageDisksView.vue'

type AdminPageId = 'cameras' | 'layouts' | 'maps' | 'carousels' | 'users' | 'workflows' | 'storage'

const activePage = ref<AdminPageId>('cameras')

const adminPages: { id: AdminPageId; label: string }[] = [
  { id: 'cameras', label: '[CAMERAS]' },
  { id: 'layouts', label: '[LAYOUTS]' },
  { id: 'maps', label: '[MAPAS & PLANTAS]' },
  { id: 'carousels', label: '[RONDAS]' },
  { id: 'users', label: '[USUARIOS & RBAC]' },
  { id: 'workflows', label: '[WORKFLOWS IA]' },
  { id: 'storage', label: '[STORAGE & RETENCAO]' }
]
</script>

<template>
  <div class="vms-admin-container" style="display: flex; flex-direction: column; flex: 1; height: 100%; overflow: hidden; background: var(--vms-neu-bg);">
    <!-- Top Horizontal Page Navigation Bar -->
    <div style="display: flex; align-items: center; justify-content: space-between; padding: 0.5rem 1.25rem; background: var(--vms-neu-surface); border-bottom: 1px solid var(--vms-border); box-shadow: 0 4px 12px rgba(0,0,0,0.3); z-index: 10;">
      <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
        <button
          v-for="page in adminPages"
          :key="page.id"
          class="vms-btn vms-btn-sm"
          :style="{
            backgroundColor: activePage === page.id ? 'var(--vms-neu-bg)' : 'transparent',
            color: activePage === page.id ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-regular)',
            boxShadow: activePage === page.id ? 'var(--vms-neu-inset-dark), var(--vms-neu-inset-light)' : 'none',
            border: activePage === page.id ? '1px solid rgba(255,94,58,0.3)' : '1px solid transparent',
            fontWeight: activePage === page.id ? '600' : '400',
            padding: '0.4rem 0.85rem'
          }"
          @click="activePage = page.id"
        >
          {{ page.label }}
        </button>
      </div>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">HYDRA VMS // ADMIN CENTER</span>
    </div>

    <!-- Full-Width Page View Host -->
    <main class="vms-admin-content" style="flex: 1; padding: 1.5rem 2rem; overflow-y: auto;">
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
