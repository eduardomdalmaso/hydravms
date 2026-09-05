<script setup lang="ts">
import type { AlarmItem } from '../../../types/alarmTree'
import AlarmDetailLeftPane from './AlarmDetailLeftPane.vue'
import AlarmZoneRightPane from './AlarmZoneRightPane.vue'

defineProps<{ alarm: AlarmItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
</script>

<template>
  <div class="vms-desktop-split-layout">
    <AlarmDetailLeftPane :alarm="alarm" />
    <AlarmZoneRightPane :alarm="alarm" @saved="(msg) => emit('saved', msg)" />
  </div>
</template>
