<script setup lang="ts">
import type { AlarmItem } from '../../../types/alarmTree'

defineProps<{
  alarm: AlarmItem
  isSelected: boolean
}>()

const emit = defineEmits<{
  (e: 'select', alarm: AlarmItem): void
  (e: 'dragstart', alarm: AlarmItem): void
  (e: 'context', event: MouseEvent, alarm: AlarmItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ active: isSelected }"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragstart', alarm)"
    @click="emit('select', alarm)"
    @contextmenu.prevent="emit('context', $event, alarm)"
  >
    <!-- LED Status Dot -->
    <span
      style="position: absolute; top: 8px; left: 8px; width: 7px; height: 7px; border-radius: 50%;"
      :style="{
        background: alarm.status === 'online' ? '#00ff9d' : alarm.status === 'unarmed' ? '#fcee0a' : '#ff003c',
        boxShadow: alarm.status === 'online' ? '0 0 6px #00ff9d' : alarm.status === 'unarmed' ? '0 0 6px #fcee0a' : '0 0 6px #ff003c'
      }"
      :title="alarm.status === 'online' ? 'ARMADO' : alarm.status === 'unarmed' ? 'DESARMADO' : 'OFFLINE'"
    />

    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
        <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ alarm.name }}
    </span>
  </div>
</template>
