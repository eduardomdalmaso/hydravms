<script setup lang="ts">
import type { PageHelpGuide } from '../../../types/helpTypes'

defineProps<{
  guide: PageHelpGuide
  activeSectionId: string
}>()

const emit = defineEmits<{ (e: 'selectSection', sectionId: string): void }>()
</script>

<template>
  <div class="vms-help-sidebar">
    <div class="vms-flex-col" style="gap: 4px; padding: 0.75rem 0.85rem; border-bottom: 1px solid var(--vms-border);">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">
        [{{ guide.badge }} // GUIA HUD]
      </span>
      <span class="vms-font-bold vms-text-xs" style="color: #fff; line-height: 1.2;">
        {{ guide.title }}
      </span>
    </div>

    <div class="vms-flex-col" style="gap: 0.35rem; padding: 0.6rem; flex: 1; overflow-y: auto;">
      <button
        v-for="sec in guide.sections"
        :key="sec.id"
        class="vms-help-tab-btn"
        :class="{ active: activeSectionId === sec.id }"
        @click="emit('selectSection', sec.id)"
      >
        <div class="vms-flex-col" style="gap: 2px; text-align: left;">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">{{ sec.title }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">[{{ sec.tag }}]</span>
        </div>
        <span class="vms-tab-indicator" />
      </button>
    </div>

    <div class="vms-flex-col" style="padding: 0.75rem; border-top: 1px solid var(--vms-border); background: rgba(0,0,0,0.2);">
      <span class="vms-text-mono vms-text-2xs vms-text-dim">HYDRA CORE VMS</span>
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-green);">BUILD 2026 // ONLINE</span>
    </div>
  </div>
</template>

<style scoped>
.vms-help-sidebar {
  width: 210px; background: #0b0e14; border-right: 1px solid var(--vms-border);
  display: flex; flex-direction: column; height: 100%; flex-shrink: 0;
}
.vms-help-tab-btn {
  display: flex; justify-content: space-between; align-items: center; padding: 0.5rem 0.65rem;
  background: transparent; border: 1px solid transparent; border-radius: 6px; cursor: pointer;
  transition: all 0.2s ease;
}
.vms-help-tab-btn:hover { background: rgba(255, 255, 255, 0.04); border-color: rgba(255, 255, 255, 0.08); }
.vms-help-tab-btn.active {
  background: rgba(255, 94, 58, 0.12); border-color: rgba(255, 94, 58, 0.35);
}
.vms-tab-indicator {
  width: 4px; height: 16px; border-radius: 2px; background: transparent; transition: background 0.2s;
}
.vms-help-tab-btn.active .vms-tab-indicator {
  background: var(--vms-neu-accent-orange); box-shadow: 0 0 8px var(--vms-neu-accent-orange);
}
</style>
