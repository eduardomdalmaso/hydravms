<script setup lang="ts">
import type { MapFolderNode, EnterpriseMapItem } from '../../../types/mapTree'
import { useI18n } from '../../../composables/useI18n'

defineProps<{
  currentFolder: MapFolderNode | null
  selectedMap: EnterpriseMapItem | null
  totalFolders: number
  totalMaps: number
  currentFolderMapsCount?: number
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'navigateRoot'): void
  (e: 'navigateFolder'): void
  (e: 'save'): void
  (e: 'toggleLock'): void
}>()

const { t } = useI18n()
</script>

<template>
  <div class="vms-flex-between vms-breadcrumb-container">
    <div class="vms-breadcrumb-trail">
      <button v-if="selectedMap || currentFolder" class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 6px; font-size: 11px;" @click="emit('back')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        <span>{{ t('back') }}</span>
      </button>

      <span class="vms-breadcrumb-item" :class="{ 'vms-breadcrumb-active': !currentFolder && !selectedMap }" @click="emit('navigateRoot')">
        {{ t('maps_root') }}
      </span>

      <template v-if="currentFolder">
        <span class="vms-text-dim">/</span>
        <span class="vms-breadcrumb-item" :class="{ 'vms-breadcrumb-active': currentFolder && !selectedMap }" @click="emit('navigateFolder')">
          {{ currentFolder.name }}
        </span>
      </template>

      <template v-if="selectedMap">
        <span class="vms-text-dim">/</span>
        <span class="vms-breadcrumb-item vms-breadcrumb-active" style="color: var(--vms-neu-accent-orange);">
          {{ selectedMap.name }}
        </span>
      </template>
    </div>

    <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
      <template v-if="!selectedMap">
        <span v-if="!currentFolder" class="vms-badge vms-badge-neutral" style="font-size: 10px;">
          {{ totalFolders }} {{ t('folders') }} // {{ totalMaps }} {{ t('maps_label') }}
        </span>
        <span v-else class="vms-badge vms-badge-neutral" style="font-size: 10px;">
          {{ currentFolderMapsCount || 0 }} {{ t('maps_label') }}
        </span>
      </template>

      <template v-else>
        <span class="vms-badge" style="font-size: 10px; background: #14171d; border: 1px solid var(--vms-border); color: var(--vms-neu-accent-orange);">
          {{ selectedMap.markers.length }} {{ t('pins_fixed') }}
        </span>

        <button
          type="button"
          class="vms-btn vms-btn-sm"
          :class="selectedMap.is_locked ? 'vms-btn-primary' : 'vms-btn-secondary'"
          :style="{
            padding: '4px 8px',
            background: selectedMap.is_locked ? 'rgba(255, 94, 58, 0.2)' : '#07080c',
            borderColor: selectedMap.is_locked ? '#ff5e3a' : 'var(--vms-border)'
          }"
          :title="selectedMap.is_locked ? t('locked_badge') : t('unlocked_badge')"
          @click="emit('toggleLock')"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" :fill="selectedMap.is_locked ? '#ff5e3a' : 'none'" stroke="#ff5e3a" stroke-width="1.5">
            <path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/>
          </svg>
        </button>

        <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-weight: bold;" @click="emit('save')">
          {{ t('save') }}
        </button>
      </template>
    </div>
  </div>
</template>
