<script setup lang="ts">
import type { AlarmFolderNode, AlarmItem } from '../../../types/alarmTree'

defineProps<{
  currentFolder?: AlarmFolderNode | null
  selectedAlarm?: AlarmItem | null
  totalFolders: number
  totalAlarms: number
  currentFolderAlarmsCount?: number
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'navigate-root'): void
  (e: 'navigate-folder'): void
}>()
</script>

<template>
  <div class="vms-flex-between" style="padding: 0.5rem 0.75rem; background: #16191f; border-radius: 6px; border: 1px solid var(--vms-border);">
    <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
      <button v-if="currentFolder || selectedAlarm" class="vms-btn vms-btn-secondary vms-btn-sm" title="Voltar" @click="emit('back')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
        <span>VOLTAR</span>
      </button>

      <span
        class="vms-text-mono vms-text-xs"
        :style="{ color: !currentFolder && !selectedAlarm ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }"
        style="cursor: pointer;"
        @click="emit('navigate-root')"
      >
        [TOPOLOGIA DE ZONAS]
      </span>

      <template v-if="currentFolder">
        <span class="vms-text-dim">/</span>
        <span
          class="vms-font-bold"
          :style="{ color: !selectedAlarm ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)', cursor: selectedAlarm ? 'pointer' : 'default' }"
          style="font-size: 12px;"
          @click="selectedAlarm ? emit('navigate-folder') : undefined"
        >
          {{ currentFolder.name }}
        </span>
      </template>

      <template v-if="selectedAlarm">
        <span class="vms-text-dim">/</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 12px;">
          {{ selectedAlarm.name }}
        </span>
      </template>
    </div>

    <span class="vms-text-mono vms-text-2xs vms-text-dim">
      {{ selectedAlarm ? '[INSPECAO // ATIVA]' : currentFolder ? `${currentFolderAlarmsCount ?? 0} SENSORES` : `${totalFolders} ZONAS // ${totalAlarms} SENSORES` }}
    </span>
  </div>
</template>
