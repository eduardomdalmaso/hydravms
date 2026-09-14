<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { pageHelpCatalog } from '../../../data/pageHelpData'
import type { PageHelpGuide } from '../../../types/helpTypes'
import HelpModalSidebar from './HelpModalSidebar.vue'
import HelpAnimatedVisual from './HelpAnimatedVisual.vue'

const props = defineProps<{
  isOpen: boolean
  currentPageId: string
}>()

const emit = defineEmits<{ (e: 'close'): void }>()

const activeGuide = computed<PageHelpGuide>(() => {
  return pageHelpCatalog[props.currentPageId] || pageHelpCatalog['video_streams']
})

const activeSectionId = ref<string>('overview')

watch(() => props.currentPageId, () => {
  activeSectionId.value = activeGuide.value.sections[0]?.id || 'overview'
})

const currentSection = computed(() => {
  return activeGuide.value.sections.find(s => s.id === activeSectionId.value) || activeGuide.value.sections[0]
})

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.isOpen) emit('close')
}

onMounted(() => window.addEventListener('keydown', handleKeydown))
onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-help-dialog">
      <!-- Top Header -->
      <div class="vms-modal-header" style="background: #11141d; padding: 0.65rem 1rem;">
        <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
          <div style="width: 22px; height: 22px; border-radius: 50%; background: rgba(255, 94, 58, 0.2); border: 1px solid #ff5e3a; display: flex; align-items: center; justify-content: center; font-weight: bold; color: #ff5e3a; font-size: 11px;">?</div>
          <div class="vms-flex-col" style="gap: 1px;">
            <span class="vms-font-bold vms-text-xs" style="color: #fff;">{{ activeGuide.title }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ activeGuide.subtitle }}</span>
          </div>
        </div>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar (ESC)" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <!-- Main Body: Sidebar + Animated Content -->
      <div class="vms-help-body">
        <HelpModalSidebar :guide="activeGuide" :active-section-id="activeSectionId" @select-section="activeSectionId = $event" />

        <!-- Right Content & Visuals -->
        <div class="vms-help-content">
          <!-- Animated Visual Area -->
          <HelpAnimatedVisual v-if="currentSection" :type="currentSection.animationType" />

          <!-- Description & Steps -->
          <div v-if="currentSection" class="vms-flex-col" style="gap: 0.75rem;">
            <p class="vms-text-xs" style="color: #d1d5db; line-height: 1.4; margin: 0;">
              {{ currentSection.description }}
            </p>

            <!-- Steps List -->
            <div class="vms-flex-col" style="gap: 0.35rem;">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">// COMO OPERAR // PASSO A PASSO</span>
              <div v-for="(step, idx) in currentSection.steps" :key="idx" class="vms-step-row">
                <span class="vms-step-badge">{{ idx + 1 }}</span>
                <span class="vms-text-xs" style="color: #fff;">{{ step }}</span>
              </div>
            </div>

            <!-- Practical Example Card -->
            <div class="vms-example-card">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-green);">
                [{{ currentSection.example.label }}]
              </span>
              <div class="vms-text-xs" style="color: #9ca3af;">
                <strong style="color: #fff;">Cenário:</strong> {{ currentSection.example.scenario }}
              </div>
              <div class="vms-text-xs" style="color: #9ca3af;">
                <strong style="color: #00ff9d;">Resultado:</strong> {{ currentSection.example.result }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-help-dialog {
  background: #0e1118; border: 1px solid var(--vms-border); border-radius: 10px;
  width: 780px; max-width: 95vw; height: 500px; max-height: 90vh;
  display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 16px 40px rgba(0,0,0,0.7);
  animation: help-pop 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes help-pop { from { opacity: 0; transform: scale(0.96) translateY(8px); } to { opacity: 1; transform: scale(1) translateY(0); } }

.vms-help-body { display: flex; flex: 1; overflow: hidden; }
.vms-help-content { flex: 1; padding: 1rem 1.25rem; overflow-y: auto; display: flex; flex-direction: column; gap: 0.85rem; }
.vms-step-row { display: flex; align-items: flex-start; gap: 0.6rem; background: rgba(255,255,255,0.02); padding: 0.35rem 0.55rem; border-radius: 4px; }
.vms-step-badge {
  background: rgba(255,94,58,0.2); border: 1px solid var(--vms-neu-accent-orange); color: var(--vms-neu-accent-orange);
  font-family: var(--vms-font-mono); font-size: 10px; font-weight: bold; width: 16px; height: 16px;
  border-radius: 50%; display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.vms-example-card {
  background: rgba(0, 255, 157, 0.04); border: 1px dashed rgba(0, 255, 157, 0.3); border-radius: 6px;
  padding: 0.65rem 0.85rem; display: flex; flex-direction: column; gap: 0.3rem;
}
</style>
