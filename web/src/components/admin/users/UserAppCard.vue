<script setup lang="ts">
import type { UserItem } from '../../../types/userTree'

defineProps<{
  user: UserItem
  isSelected: boolean
}>()

const emit = defineEmits<{
  (e: 'select', user: UserItem): void
  (e: 'dragstart', user: UserItem): void
  (e: 'context', event: MouseEvent, user: UserItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ active: isSelected }"
    draggable="true"
    @dragstart="emit('dragstart', user)"
    @click="emit('select', user)"
    @contextmenu.prevent="emit('context', $event, user)"
  >
    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
        <circle cx="12" cy="7" r="4"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ user.username }}
    </span>
    <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
      [{{ user.role.replace('_', ' ').toUpperCase() }}]
    </span>
  </div>
</template>
