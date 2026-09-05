<script setup lang="ts">
import { ref } from 'vue'
import type { FolderNode, StreamItem } from '../../types/streamTree'
import TreeStreamRow from './TreeStreamRow.vue'

const props = defineProps<{ folder: FolderNode }>()
const emit = defineEmits<{
  (e: 'remove-folder', folderId: string): void
  (e: 'test-stream', streamId: string): void
  (e: 'folder-context', event: MouseEvent, folder: FolderNode): void
  (e: 'stream-context', event: MouseEvent, stream: StreamItem): void
}>()

const isExpanded = ref(props.folder.isExpanded ?? true)
</script>

<template>
  <div class="vms-tree-server-card">
    <!-- Folder Header -->
    <div
      class="vms-tree-server-header"
      @click="isExpanded = !isExpanded"
      @contextmenu.prevent="emit('folder-context', $event, folder)"
    >
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <span class="vms-text-dim" style="font-size: 11px;">{{ isExpanded ? '▼' : '►' }}</span>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--vms-neu-accent-orange)" stroke-width="2"><path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"/></svg>
        <span class="vms-font-bold" style="color: #fff; font-size: 13px;">{{ folder.name }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// {{ folder.streams.length }} FLUXOS</span>
      </div>

      <div class="vms-flex-row" style="gap: 0.65rem;" @click.stop>
        <span class="vms-badge vms-badge-info">{{ folder.streams.length }} FLUXOS</span>
        <button
          class="vms-btn vms-btn-ghost vms-btn-sm"
          style="font-size: 10px; color: #ff5e3a; padding: 3px 6px; display: inline-flex; align-items: center;"
          title="Excluir Pasta"
          @click="emit('remove-folder', folder.id)"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
        </button>
      </div>
    </div>

    <!-- Streams within folder -->
    <div v-if="isExpanded" class="vms-flex-col" style="padding: 0.75rem 1rem; gap: 0.45rem;">
      <div v-if="folder.streams.length === 0" class="vms-text-mono vms-text-2xs vms-text-dim" style="padding: 0.5rem; text-align: center;">
        // PASTA VAZIA (CLIQUE COM BOTAO DIREITO PARA ADICIONAR FLUXO)
      </div>
      <TreeStreamRow
        v-for="stream in folder.streams"
        :key="stream.id"
        :stream="stream"
        @test-stream="(id) => emit('test-stream', id)"
        @stream-context="(ev, s) => emit('stream-context', ev, s)"
      />
    </div>
  </div>
</template>
