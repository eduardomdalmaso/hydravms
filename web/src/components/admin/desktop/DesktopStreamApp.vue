<script setup lang="ts">
import type { StreamItem } from '../../../types/streamTree'

defineProps<{
  stream: StreamItem
  isSelected: boolean
}>()

const emit = defineEmits<{
  (e: 'select', stream: StreamItem): void
  (e: 'dragstart', stream: StreamItem): void
  (e: 'context', event: MouseEvent, stream: StreamItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ active: isSelected }"
    draggable="true"
    @dragstart="emit('dragstart', stream)"
    @click="emit('select', stream)"
    @contextmenu.prevent="emit('context', $event, stream)"
  >
    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="24" height="24" viewBox="0 0 576 512" fill="#ff5e3a">
        <path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ stream.name }}
    </span>
  </div>
</template>
