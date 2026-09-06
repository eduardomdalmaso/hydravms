<script setup lang="ts">
defineProps<{
  activeTab: 'CAMERAS' | 'ALARMES' | null
  camerasCount: number
  alarmsCount: number
}>()

const emit = defineEmits<{
  (e: 'toggle', tab: 'CAMERAS' | 'ALARMES'): void
}>()
</script>

<template>
  <div class="vms-map-top-right-tabs">
    <button
      class="vms-map-tab-btn"
      :class="{ 'vms-map-tab-active': activeTab === 'CAMERAS' }"
      @click="emit('toggle', 'CAMERAS')"
    >
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M23 7l-7 5 7 5V7z"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
      </svg>
      <span>CAMERAS</span>
      <span class="vms-tab-counter">{{ camerasCount }}</span>
    </button>

    <button
      class="vms-map-tab-btn"
      :class="{ 'vms-map-tab-active': activeTab === 'ALARMES' }"
      @click="emit('toggle', 'ALARMES')"
    >
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/>
      </svg>
      <span>ALARMES</span>
      <span class="vms-tab-counter">{{ alarmsCount }}</span>
    </button>
  </div>
</template>

<style scoped>
.vms-map-top-right-tabs {
  position: absolute; top: 12px; right: 12px; z-index: 1000;
  display: flex; gap: 6px; background: rgba(7, 8, 12, 0.88);
  backdrop-filter: blur(8px); padding: 4px; border-radius: 8px;
  border: 1px solid var(--vms-border); box-shadow: 0 4px 20px rgba(0, 0, 0, 0.6);
}
.vms-map-tab-btn {
  display: flex; align-items: center; gap: 6px; padding: 6px 12px;
  background: #14171d; border: 1px solid var(--vms-border); border-radius: 6px;
  color: var(--vms-text-dim); font-family: var(--vms-font-roboto); font-size: 11px;
  font-weight: 700; cursor: pointer; transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.vms-map-tab-btn:hover { color: #ffffff; border-color: var(--vms-neu-accent-orange); }
.vms-map-tab-active {
  background: rgba(255, 94, 58, 0.15); border-color: #ff5e3a; color: #ff5e3a;
  box-shadow: 0 0 10px rgba(255, 94, 58, 0.25);
}
.vms-tab-counter {
  font-family: var(--vms-font-mono); font-size: 10px; background: #07080c;
  padding: 1px 6px; border-radius: 4px; border: 1px solid var(--vms-border);
}
</style>
