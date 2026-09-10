<script setup lang="ts">
import { ref } from "vue"

const props = defineProps<{ cameraId: string; alarmInfo?: { name: string; zone: string; triggers: number; status: "online" | "alert" | "offline" } | null }>()
const isAlarmOpen = ref(false)
</script>

<template>
  <div v-if="alarmInfo" class="vms-gear-wrapper" style="position: relative;" @mouseenter="isAlarmOpen = true" @mouseleave="isAlarmOpen = false">
    <button class="vms-gear-btn" :title="`Alarme: ${alarmInfo.name}`" @click.stop="isAlarmOpen = !isAlarmOpen">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path><path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
    </button>
    <div v-if="isAlarmOpen" class="vms-alarm-slot-popover">
      <div class="vms-flex-between" style="gap: 0.5rem; align-items: center; border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 0.2rem;">
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-weight: 700;">[ALARME VINCULADO]</span>
        <span class="vms-status-led" :class="alarmInfo.status"></span>
      </div>
      <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500;">
        {{ alarmInfo.name }}
      </span>
      <span class="vms-text-2xs" style="color: #cbd5e1; font-family: var(--vms-font-roboto);">
        {{ alarmInfo.zone }}
      </span>
      <div class="vms-flex-between" style="gap: 0.35rem; align-items: center; background: rgba(255,94,58,0.12); padding: 2px 5px; border-radius: 4px; border: 1px solid rgba(255,94,58,0.25);">
        <span class="vms-text-2xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-size: 10px;">DISPAROS:</span>
        <span class="vms-badge" style="background: var(--vms-neu-accent-orange); color: #000; font-size: 9px; font-weight: 800; padding: 0 4px;">
          {{ alarmInfo.triggers }} ACUMULADOS
        </span>
      </div>
    </div>
  </div>
</template>
