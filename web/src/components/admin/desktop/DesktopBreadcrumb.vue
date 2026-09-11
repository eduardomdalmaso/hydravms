<script setup lang="ts">
import type { FolderNode, StreamItem } from '../../../types/streamTree'
import { useI18n } from '../../../composables/useI18n'

defineProps<{
  currentFolder?: FolderNode | null
  selectedStream?: StreamItem | null
  totalFolders: number
  totalStreams: number
  currentFolderStreamsCount?: number
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'navigate-root'): void
  (e: 'navigate-folder'): void
}>()

const { t } = useI18n()
</script>

<template>
  <div class="vms-flex-between" style="padding: 0.5rem 0.75rem; background: #16191f; border-radius: 6px; border: 1px solid var(--vms-border);">
    <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
      <button
        v-if="currentFolder || selectedStream"
        class="vms-btn vms-btn-secondary vms-btn-sm"
        :title="t('back')"
        @click="emit('back')"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
        <span>{{ t('back') }}</span>
      </button>

      <span
        class="vms-text-mono vms-text-xs"
        :style="{ color: !currentFolder && !selectedStream ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }"
        style="cursor: pointer;"
        @click="emit('navigate-root')"
      >
        {{ t('root_topology') }}
      </span>

      <template v-if="currentFolder">
        <span class="vms-text-dim">/</span>
        <span
          class="vms-font-bold"
          :style="{ color: !selectedStream ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)', cursor: selectedStream ? 'pointer' : 'default' }"
          style="font-size: 12px;"
          @click="selectedStream ? emit('navigate-folder') : undefined"
        >
          {{ currentFolder.name }}
        </span>
      </template>

      <template v-if="selectedStream">
        <span class="vms-text-dim">/</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 12px;">
          {{ selectedStream.name }}
        </span>
      </template>
    </div>

    <span class="vms-text-mono vms-text-2xs vms-text-dim">
      {{ selectedStream ? t('inspect_active') : currentFolder ? `${currentFolderStreamsCount ?? 0} ${t('streams')}` : `${totalFolders} ${t('folders')} // ${totalStreams} ${t('streams')}` }}
    </span>
  </div>
</template>
