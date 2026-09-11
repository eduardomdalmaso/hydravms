<script setup lang="ts">
import { ref, computed } from "vue"
import { useI18n } from "../../composables/useI18n"
import type { CustomLayout } from "../../types/mosaic"
import GridIconPreview from "./GridIconPreview.vue"

const props = defineProps<{ layouts: CustomLayout[]; activeLayoutId: string; isOpen: boolean }>()
const emit = defineEmits<{
  (e: "toggle"): void
  (e: "selectLayout", id: string): void
  (e: "newLayout"): void
  (e: "contextMenu", payload: { event: MouseEvent; layout?: CustomLayout; isHeader?: boolean }): void
}>()
const { t } = useI18n()

const isSearchOpen = ref(false)
const searchQuery = ref("")

const filteredLayouts = computed(() => {
  if (!searchQuery.value) return props.layouts
  return props.layouts.filter(l => l.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})
</script>

<template>
  <div>
    <!-- Clickable Header Row with Orange Counter -->
    <div class="vms-accordion-header" :class="{ active: isOpen }" title="Clique na linha para expandir ou recuar" @click="emit('toggle')">
      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
        <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[{{ t('layouts').toUpperCase() }}]</span>
        <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.35); font-size: 8.5px; font-weight: 700; padding: 1px 5px;">
          {{ filteredLayouts.length }}
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <button class="vms-sidebar-icon-btn" :title="t('search_lay')" @click.stop="isSearchOpen = !isSearchOpen">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
        <button class="vms-sidebar-icon-btn" title="Criar novo layout" @click.stop="emit('newLayout')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search Input for Layouts -->
    <div v-if="isOpen && isSearchOpen" style="padding: 0.4rem 0.5rem; background: #11141b; border-bottom: 1px solid var(--vms-border);">
      <input v-model="searchQuery" class="vms-auth-input" style="padding: 0.25rem 0.5rem; font-size: 11px; height: 26px;" :placeholder="t('search_lay')" autofocus />
    </div>

    <!-- Scrollable Layout List with Clean Grid Icon and Name -->
    <div v-if="isOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14; max-height: 200px; overflow-y: auto;">
      <div
        v-for="lay in filteredLayouts"
        :key="lay.id"
        class="vms-asset-item"
        :class="{ active: lay.id === activeLayoutId }"
        style="padding: 0.4rem 0.55rem;"
        :title="`${lay.name} [Grade ${lay.grid.toUpperCase()}]`"
        @click="emit('selectLayout', lay.id)"
        @contextmenu.prevent="emit('contextMenu', { event: $event, layout: lay, isHeader: false })"
      >
        <div class="vms-flex-row" style="gap: 0.4rem; min-width: 0; flex: 1; align-items: center;">
          <GridIconPreview :grid="lay.grid" />
          <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ lay.name }}</span>
          <svg v-if="lay.is_system" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0;">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>
