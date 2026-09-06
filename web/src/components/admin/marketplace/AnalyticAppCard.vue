<script setup lang="ts">
import type { AnalyticInstance } from '../../../types/marketplace'

defineProps<{
  instance: AnalyticInstance
  isSelected?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', inst: AnalyticInstance): void
  (e: 'dragstart', inst: AnalyticInstance): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ selected: isSelected }"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragstart', instance)"
    @click="emit('select', instance)"
  >
    <!-- Status Indicator Pill -->
    <div
      style="position: absolute; top: 6px; right: 6px; width: 8px; height: 8px; border-radius: 50%;"
      :style="{ background: instance.is_active ? '#00ff9d' : '#ff5e3a', boxShadow: instance.is_active ? '0 0 6px #00ff9d' : 'none' }"
      :title="instance.is_active ? '[STATUS: ATIVO]' : '[STATUS: PAUSADO]'"
    />

    <!-- 44x44 Orange Icon Container (Standard) -->
    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="12 2 2 7 12 12 22 7 12 2"/>
        <polyline points="2 17 12 22 22 17"/>
        <polyline points="2 12 12 17 22 12"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ instance.name }}
    </span>

    <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; justify-content: center; flex-wrap: wrap;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
        {{ instance.is_active ? '[ATIVO]' : '[PAUSADO]' }}
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8.5px;">
        // {{ instance.camera_name.split('//')[0]?.trim() }}
      </span>
    </div>
  </div>
</template>
