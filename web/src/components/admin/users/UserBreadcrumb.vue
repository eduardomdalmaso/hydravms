<script setup lang="ts">
import type { UserFolderNode, UserItem } from '../../../types/userTree'

defineProps<{
  currentFolder?: UserFolderNode | null
  selectedUser?: UserItem | null
  totalFolders: number
  totalUsers: number
  currentFolderUsersCount?: number
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
      <button v-if="currentFolder || selectedUser" class="vms-btn vms-btn-secondary vms-btn-sm" title="Voltar" @click="emit('back')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
        <span>VOLTAR</span>
      </button>

      <span
        class="vms-text-mono vms-text-xs"
        :style="{ color: !currentFolder && !selectedUser ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }"
        style="cursor: pointer;"
        @click="emit('navigate-root')"
      >
        [GRUPOS & DEPARTAMENTOS]
      </span>

      <template v-if="currentFolder">
        <span class="vms-text-dim">/</span>
        <span
          class="vms-font-bold"
          :style="{ color: !selectedUser ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)', cursor: selectedUser ? 'pointer' : 'default' }"
          style="font-size: 12px;"
          @click="selectedUser ? emit('navigate-folder') : undefined"
        >
          {{ currentFolder.name }}
        </span>
      </template>

      <template v-if="selectedUser">
        <span class="vms-text-dim">/</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 12px;">
          {{ selectedUser.username }}
        </span>
      </template>
    </div>

    <span class="vms-text-mono vms-text-2xs vms-text-dim">
      {{ selectedUser ? '[INSPECAO // ATIVA]' : currentFolder ? `${currentFolderUsersCount ?? 0} USUARIOS` : `${totalFolders} GRUPOS // ${totalUsers} USUARIOS` }}
    </span>
  </div>
</template>
