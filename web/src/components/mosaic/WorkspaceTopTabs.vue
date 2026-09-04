<script setup lang="ts">
import type { WorkspaceTab } from '../../types/mosaic'

defineProps<{ tabs: WorkspaceTab[]; activeTabId: string }>()
const emit = defineEmits<{
  (e: 'selectTab', id: string): void
  (e: 'closeTab', id: string): void
  (e: 'saveTab', id: string): void
}>()
</script>

<template>
  <div class="vms-workspace-top-tabs">
    <!-- Horizontal List of Open Layout Tabs -->
    <div class="vms-tabs-scroll-area">
      <div
        v-for="tab in tabs"
        :key="tab.id"
        class="vms-workspace-tab"
        :class="{ active: tab.id === activeTabId, unsaved: tab.is_temporary || !tab.is_saved }"
        :title="`${tab.name} [Grade ${tab.grid.toUpperCase()}]`"
        @click="emit('selectTab', tab.id)"
      >
        <span class="vms-tab-title">{{ tab.name }}</span>

        <!-- Padlock icon for Admin/System layout tabs -->
        <svg
          v-if="tab.is_system"
          width="10"
          height="10"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#ff5e3a"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          style="flex-shrink: 0;"
        >
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
          <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
        </svg>

        <!-- Save Icon for Temporary/Unsaved Operator Layouts -->
        <button
          v-if="tab.is_temporary || !tab.is_saved"
          class="vms-tab-icon-action-btn"
          title="Salvar layout no seu perfil"
          @click.stop="emit('saveTab', tab.id)"
        >
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
            <polyline points="17 21 17 13 7 13 7 21"></polyline>
            <polyline points="7 3 7 8 15 8"></polyline>
          </svg>
        </button>

        <!-- Standard Ubuntu-Style Circular Close Button -->
        <button
          class="vms-ubuntu-close-btn"
          style="position: static; width: 15px; height: 15px; opacity: 1; pointer-events: auto; padding: 0; flex-shrink: 0;"
          title="Fechar aba"
          @click.stop="emit('closeTab', tab.id)"
        >
          <svg width="7" height="7" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <path d="M2 2L10 10M10 2L2 10" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>
