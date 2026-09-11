<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from '../../composables/useI18n'
import type { CarouselConfig } from '../../types/mosaic'

const props = defineProps<{ carousels: CarouselConfig[]; isOpen: boolean }>()
const emit = defineEmits<{
  (e: 'toggle'): void
  (e: 'selectCarousel', carousel: CarouselConfig): void
  (e: 'openModal'): void
}>()
const { t } = useI18n()

const isSearchOpen = ref(false)
const searchQuery = ref('')

const filteredCarousels = computed(() => {
  if (!searchQuery.value) return props.carousels
  return props.carousels.filter(c => c.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})
</script>

<template>
  <div>
    <!-- Clickable Header Row with Orange Counter -->
    <div class="vms-accordion-header" :class="{ active: isOpen }" title="Clique na linha para expandir ou recuar" @click="emit('toggle')">
      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
        <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[{{ t('carousels').toUpperCase() }}]</span>
        <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.35); font-size: 8.5px; font-weight: 700; padding: 1px 5px;">
          {{ filteredCarousels.length }}
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <button class="vms-sidebar-icon-btn" :title="t('search_rnd')" @click.stop="isSearchOpen = !isSearchOpen">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
        <button class="vms-sidebar-icon-btn" title="Criar nova ronda" @click.stop="emit('openModal')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search Bar for Carousels -->
    <div v-if="isOpen && isSearchOpen" style="padding: 0.4rem 0.5rem; background: #11141b; border-bottom: 1px solid var(--vms-border);">
      <input v-model="searchQuery" class="vms-auth-input" style="padding: 0.25rem 0.5rem; font-size: 11px; height: 26px;" :placeholder="t('search_rnd')" autofocus />
    </div>

    <!-- Scrollable Carousels List with Clean Name -->
    <div v-if="isOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14; max-height: 200px; overflow-y: auto;">
      <div v-for="c in filteredCarousels" :key="c.id" class="vms-asset-item" style="padding: 0.4rem 0.55rem;" @click="emit('selectCarousel', c)">
        <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ c.name }}</span>
      </div>
    </div>
  </div>
</template>
