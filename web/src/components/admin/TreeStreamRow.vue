<script setup lang="ts">
import type { StreamItem } from '../../types/streamTree'

defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{
  (e: 'test-stream', id: string): void
  (e: 'stream-context', event: MouseEvent, stream: StreamItem): void
}>()
</script>

<template>
  <div class="vms-tree-camera-row" @contextmenu.prevent="emit('stream-context', $event, stream)">
    <div class="vms-flex-row" style="gap: 0.75rem;">
      <span class="vms-status-led" :class="stream.status" :title="stream.status === 'recording' ? 'Gravacao Ativa' : (stream.status === 'offline' ? 'Offline' : 'Online')"></span>
      <span class="vms-font-medium" style="color: #fff; font-size: 12px;">{{ stream.name }}</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ stream.ip }}</span>
    </div>
    <div class="vms-flex-row" style="gap: 0.65rem;">
      <span class="vms-badge vms-badge-info" style="font-size: 9px;">{{ stream.codec }}</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ stream.resolution }} @ {{ stream.fps }}FPS</span>
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-blue);">{{ stream.bitrate }}</span>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 2px 6px;" @click.stop="emit('test-stream', stream.id)">[TESTAR]</button>
    </div>
  </div>
</template>
