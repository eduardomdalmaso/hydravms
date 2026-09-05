<script setup lang="ts">
import type { EnterpriseLayoutItem } from '../../../types/layoutTree'

defineProps<{ layout: EnterpriseLayoutItem; isSelected: boolean }>()
const emit = defineEmits<{
  (e: 'select', layout: EnterpriseLayoutItem): void
  (e: 'dragstart', layout: EnterpriseLayoutItem): void
  (e: 'context', event: MouseEvent, layout: EnterpriseLayoutItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ active: isSelected }"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragstart', layout)"
    @click="emit('select', layout)"
    @contextmenu.prevent="emit('context', $event, layout)"
  >
    <!-- Locked Padlock Icon if layout is locked -->
    <div v-if="layout.is_locked" style="position: absolute; top: 6px; right: 6px;" title="[TRAVADO] Somente leitura para clientes/operadores">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
      </svg>
    </div>

    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center;">
      <!-- Grid icon -->
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M3 15h18"/><path d="M9 3v18"/><path d="M15 3v18"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ layout.name }}
    </span>

    <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; justify-content: center; flex-wrap: wrap;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
        [{{ layout.grid }}]
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8.5px;">
        // {{ layout.slots.filter(s => !!s.cameraId).length }} CAMS
      </span>
    </div>
  </div>
</template>
