<script setup lang="ts">
import { ref } from 'vue'
import type { AnalyticFolderNode } from '../../../types/marketplace'

const props = defineProps<{ folder: AnalyticFolderNode }>()
const emit = defineEmits<{
  (e: 'open', id: string): void
  (e: 'drop-instance', folderId: string): void
}>()

const isDragOver = ref(false)

const onDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragOver.value = false
  emit('drop-instance', props.folder.id)
}
</script>

<template>
  <div
    class="vms-desktop-folder-card"
    :class="{ 'drag-over': isDragOver }"
    style="position: relative;"
    @click="emit('open', folder.id)"
    @dragover.prevent="isDragOver = true"
    @dragleave="isDragOver = false"
    @drop="onDrop"
  >
    <!-- Border Counter Badge Always Orange -->
    <span style="position: absolute; top: 6px; right: 6px; font-family: var(--vms-font-jetbrains); font-size: 9px; font-weight: 700; color: #ff5e3a; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 4px; padding: 1px 5px; box-shadow: 0 0 6px rgba(255, 94, 58, 0.25);">
      {{ folder.instances.length }}
    </span>

    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="24" height="24" viewBox="0 0 512 512" fill="#ff5e3a">
        <path d="M64 480H448c35.3 0 64-28.7 64-64V160c0-35.3-28.7-64-64-64H288c-10.1 0-19.6-4.7-25.6-12.8L243.2 57.6C231.1 41.5 212.1 32 192 32H64C28.7 32 0 60.7 0 96V416c0 35.3 28.7 64 64 64z"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ folder.name }}
    </span>
  </div>
</template>
