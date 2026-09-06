<script setup lang="ts">
import type { EnterpriseWorkflowItem } from '../../../types/workflowTree'

defineProps<{ workflow: EnterpriseWorkflowItem }>()
const emit = defineEmits<{
  (e: 'open', w: EnterpriseWorkflowItem): void
  (e: 'dragStart', w: EnterpriseWorkflowItem): void
}>()
</script>

<template>
  <div
    class="vms-desktop-app-card"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragStart', workflow)"
    @click="emit('open', workflow)"
  >
    <!-- Lock Indicator -->
    <div v-if="workflow.is_locked" style="position: absolute; top: 6px; right: 6px;" title="[TRAVADO COM CADEADO]">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="#ff5e3a">
        <path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/>
      </svg>
    </div>

    <!-- 44x44 Orange Icon Container (Standard) -->
    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35);">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="6" height="6" rx="1"/>
        <rect x="15" y="3" width="6" height="6" rx="1"/>
        <rect x="9" y="15" width="6" height="6" rx="1"/>
        <path d="M6 9v3a1 1 0 0 0 1 1h10a1 1 0 0 0 1-1V9"/>
        <path d="M12 13v2"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ workflow.name }}
    </span>

    <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; justify-content: center; flex-wrap: wrap;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
        [FLOW]
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8.5px;">
        // {{ workflow.nodes.length }} BLOCOS
      </span>
    </div>
  </div>
</template>
