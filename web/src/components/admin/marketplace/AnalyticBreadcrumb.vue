<script setup lang="ts">
import type { AnalyticFolderNode, AnalyticInstance } from '../../../types/marketplace'
import { useI18n } from '../../../composables/useI18n'

defineProps<{
  currentFolder: AnalyticFolderNode | null
  selectedInstance: AnalyticInstance | null
  totalFolders: number
  totalInstances: number
  folderInstancesCount?: number
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'navigateRoot'): void
  (e: 'save'): void
}>()

const { t } = useI18n()
</script>

<template>
  <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.3); padding: 6px 12px; border-radius: 6px; border: 1px solid var(--vms-border);">
    <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
      <button
        v-if="currentFolder || selectedInstance"
        class="vms-btn vms-btn-ghost vms-btn-sm"
        style="padding: 2px 6px; font-size: 11px;"
        @click="emit('back')"
      >
        ◀ {{ t('back') }}
      </button>

      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center; font-size: 11px; font-family: var(--vms-font-jetbrains);">
        <span
          style="cursor: pointer; color: var(--vms-neu-accent-orange); font-weight: 700;"
          @click="emit('navigateRoot')"
        >
          {{ t('root') }}
        </span>
        <template v-if="currentFolder">
          <span style="color: #64748b;">/</span>
          <span style="color: #ffffff; font-weight: 600;">{{ currentFolder.name }}</span>
        </template>
        <template v-if="selectedInstance">
          <span style="color: #64748b;">/</span>
          <span style="color: #00f0ff; font-weight: 600;">{{ selectedInstance.name }}</span>
        </template>
      </div>
    </div>

    <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
      <span class="vms-text-mono vms-text-2xs vms-text-dim">
        <template v-if="currentFolder">
          {{ folderInstancesCount ?? 0 }} {{ t('analytics_label') }} {{ t('in_folder') }}
        </template>
        <template v-else>
          {{ totalFolders }} {{ t('folders') }} // {{ totalInstances }} {{ t('analytics_label') }}
        </template>
      </span>
      <button
        v-if="selectedInstance"
        class="vms-btn vms-btn-primary vms-btn-sm"
        style="font-size: 10px; padding: 3px 12px; font-weight: bold;"
        @click="emit('save')"
      >
        {{ t('save') }}
      </button>
    </div>
  </div>
</template>
