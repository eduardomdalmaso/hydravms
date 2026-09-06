<script setup lang="ts">
import { computed } from 'vue'
import type { EnterpriseRondaItem } from '../../../types/rondaTree'

const props = defineProps<{ ronda: EnterpriseRondaItem; isSelected: boolean }>()
const emit = defineEmits<{
  (e: 'select', ronda: EnterpriseRondaItem): void
  (e: 'dragstart', ronda: EnterpriseRondaItem): void
  (e: 'context', event: MouseEvent, ronda: EnterpriseRondaItem): void
}>()

const totalDuration = computed(() =>
  props.ronda.streams.reduce((acc, s) => acc + (s.intervalSeconds || 0), 0)
)
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ active: isSelected }"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragstart', ronda)"
    @click="emit('select', ronda)"
    @contextmenu.prevent="emit('context', $event, ronda)"
  >
    <div v-if="ronda.is_locked" style="position: absolute; top: 6px; right: 6px;" title="[TRAVADO]">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
    </div>

    <div style="width: 44px; height: 44px; border-radius: 10px; background: rgba(255, 94, 58, 0.12); border: 1px solid rgba(255, 94, 58, 0.35); display: flex; align-items: center; justify-content: center;">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
      </svg>
    </div>

    <span class="vms-font-medium" style="color: #fff; font-size: 11px; line-height: 1.2; word-break: break-word; max-width: 100%;">
      {{ ronda.name }}
    </span>

    <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; justify-content: center; flex-wrap: wrap;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); font-size: 9px;">
        [{{ totalDuration }}s]
      </span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim" style="font-size: 8.5px;">
        // {{ ronda.streams.length }} FLUXOS
      </span>
      <span class="vms-status-led" :class="ronda.status === 'ATIVO' ? 'online' : 'warning'" style="width: 6px; height: 6px;"></span>
    </div>
  </div>
</template>
