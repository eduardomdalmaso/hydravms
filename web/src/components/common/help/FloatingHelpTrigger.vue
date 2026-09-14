<script setup lang="ts">
defineProps<{ isOpen: boolean }>()
const emit = defineEmits<{ (e: 'toggle'): void }>()
</script>

<template>
  <div class="vms-floating-help-wrapper">
    <button
      class="vms-floating-help-btn"
      :class="{ active: isOpen }"
      title="Guia da Tela & Ajuda Rápida"
      @click="emit('toggle')"
    >
      <div class="vms-btn-pulse" />
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10" />
        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
        <line x1="12" y1="17" x2="12.01" y2="17" />
      </svg>
      <span class="vms-help-pill-label">AJUDA</span>
    </button>
  </div>
</template>

<style scoped>
.vms-floating-help-wrapper {
  position: fixed; bottom: 1.25rem; right: 1.25rem; z-index: 999;
}
.vms-floating-help-btn {
  display: flex; align-items: center; gap: 6px; padding: 6px 12px;
  background: #121620; border: 1px solid rgba(255, 94, 58, 0.4); border-radius: 20px;
  color: #fff; cursor: pointer; box-shadow: 0 4px 16px rgba(0,0,0,0.5);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); position: relative;
}
.vms-floating-help-btn:hover {
  background: #1a202c; border-color: var(--vms-neu-accent-orange); transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 94, 58, 0.25);
}
.vms-floating-help-btn.active {
  background: var(--vms-neu-accent-orange); color: #fff; border-color: #fff;
}
.vms-btn-pulse {
  position: absolute; top: 0; left: 0; right: 0; bottom: 0; border-radius: 20px;
  border: 1px solid var(--vms-neu-accent-orange); animation: btn-glow 2.5s infinite; pointer-events: none;
}
.vms-help-pill-label {
  font-family: var(--vms-font-mono); font-size: 10px; font-weight: 700; letter-spacing: 0.5px;
}
@keyframes btn-glow {
  0% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.08); opacity: 0; }
  100% { transform: scale(1); opacity: 0; }
}
</style>
