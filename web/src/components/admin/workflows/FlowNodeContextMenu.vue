<script setup lang="ts">
defineProps<{
  x: number
  y: number
  nodeTitle?: string
  isLocked?: boolean
}>()

const emit = defineEmits<{
  (e: 'delete'): void
  (e: 'close'): void
}>()
</script>

<template>
  <div class="vms-node-context-backdrop" @click="emit('close')" @contextmenu.prevent="emit('close')">
    <div
      class="vms-node-context-menu"
      :style="{ left: `${x}px`, top: `${y}px` }"
      @click.stop
    >
      <button
        class="vms-context-icon-btn"
        :disabled="isLocked"
        title="Deletar"
        @click="emit('delete')"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
          <polyline points="3 6 5 6 21 6"/>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.vms-node-context-backdrop {
  position: fixed; inset: 0; z-index: 9999; background: transparent;
}
.vms-node-context-menu {
  position: absolute; background: #0e1117; border: 1px solid var(--vms-border);
  border-radius: 8px; padding: 4px; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.7);
  display: flex; align-items: center; justify-content: center;
}
.vms-context-icon-btn {
  display: flex; align-items: center; justify-content: center; width: 32px; height: 32px;
  background: #14171d; border: 1px solid var(--vms-border); border-radius: 6px;
  cursor: pointer; transition: all 0.15s ease;
}
.vms-context-icon-btn:hover:not(:disabled) {
  background: rgba(255, 94, 58, 0.2); border-color: #ff5e3a; transform: scale(1.05);
}
.vms-context-icon-btn:disabled {
  opacity: 0.4; cursor: not-allowed;
}
</style>
