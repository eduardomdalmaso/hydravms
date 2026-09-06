<script setup lang="ts">
import { ref } from 'vue'
import type { RondaFolderNode } from '../../../types/rondaTree'

const props = defineProps<{ folder: RondaFolderNode }>()
const emit = defineEmits<{
  (e: 'open', id: string): void
  (e: 'context', event: MouseEvent, folder: RondaFolderNode): void
  (e: 'drop-ronda', folderId: string): void
}>()

const isDragOver = ref(false)
const onDrop = (e: DragEvent) => {
  e.preventDefault(); isDragOver.value = false; emit('drop-ronda', props.folder.id)
}
</script>

<template>
  <div
    class="vms-desktop-folder-card"
    :class="{ 'drag-over': isDragOver }"
    style="position: relative;"
    @click="emit('open', folder.id)"
    @contextmenu.prevent="emit('context', $event, folder)"
    @dragover.prevent="isDragOver = true"
    @dragleave="isDragOver = false"
    @drop="onDrop"
  >
    <span style="position: absolute; top: 6px; right: 6px; font-family: var(--vms-font-jetbrains); font-size: 9px; font-weight: 700; color: #ff5e3a; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 4px; padding: 1px 5px;">
      {{ folder.rondas.length }}
    </span>

    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center;">
      <svg width="24" height="24" viewBox="0 0 512 512" fill="#ff5e3a">
        <path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ folder.name }}
    </span>
  </div>
</template>
