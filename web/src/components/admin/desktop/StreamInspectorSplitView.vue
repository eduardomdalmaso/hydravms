<script setup lang="ts">
import type { StreamItem } from '../../../types/streamTree'
import StreamDetailLeftPane from './StreamDetailLeftPane.vue'
import StreamRecordingRightPane from './StreamRecordingRightPane.vue'

defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
</script>

<template>
  <div class="vms-desktop-split-layout">
    <StreamDetailLeftPane :stream="stream" />
    <StreamRecordingRightPane :stream="stream" @saved="(msg) => emit('saved', msg)" />
  </div>
</template>
