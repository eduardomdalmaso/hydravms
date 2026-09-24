<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { AnalyticInstance } from '../../../types/marketplace'

const props = defineProps<{
  instance: AnalyticInstance
  isSelected?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', inst: AnalyticInstance): void
  (e: 'dragstart', inst: AnalyticInstance): void
  (e: 'pause', inst: AnalyticInstance): void
  (e: 'edit', inst: AnalyticInstance): void
  (e: 'delete', id: string): void
}>()

const showMenu = ref(false)
const menuPos = ref({ x: 0, y: 0 })

const handleContextMenu = (e: MouseEvent) => {
  e.preventDefault()
  e.stopPropagation()
  showMenu.value = true
  menuPos.value = { x: e.clientX, y: e.clientY }
}

const closeMenu = () => { showMenu.value = false }

onMounted(() => { window.addEventListener('click', closeMenu) })
onUnmounted(() => { window.removeEventListener('click', closeMenu) })
</script>

<template>
  <div
    class="vms-desktop-app-card"
    :class="{ selected: isSelected }"
    draggable="true"
    style="position: relative;"
    @dragstart="emit('dragstart', instance)"
    @click="emit('select', instance)"
    @contextmenu="handleContextMenu"
  >
    <!-- Status Indicator Pill -->
    <div
      style="position: absolute; top: 6px; right: 6px; width: 8px; height: 8px; border-radius: 50%;"
      :style="{ background: instance.is_active ? '#00ff9d' : '#ff5e3a', boxShadow: instance.is_active ? '0 0 6px #00ff9d' : 'none' }"
      :title="instance.is_active ? '[STATUS: ATIVO]' : '[STATUS: PAUSADO]'"
    />

    <!-- 44x44 Orange Icon Container -->
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

    <!-- Context Menu Cyberpunk Modal -->
    <Teleport to="body">
      <div
        v-if="showMenu"
        class="vms-card"
        :style="{ position: 'fixed', top: `${menuPos.y}px`, left: `${menuPos.x}px`, zIndex: 99999 }"
        style="padding: 4px; background: #0b0e14; border: 1px solid var(--vms-neu-accent-orange); box-shadow: 0 8px 24px rgba(0,0,0,0.8); min-width: 140px; display: flex; flex-direction: column; gap: 2px;"
        @click.stop
      >
        <button
          class="vms-btn vms-btn-ghost vms-text-2xs vms-text-mono"
          style="justify-content: flex-start; padding: 6px 10px; font-size: 10px; text-align: left;"
          @click="emit('pause', instance); closeMenu()"
        >
          {{ instance.is_active ? '[PAUSAR]' : '[RETOMAR]' }}
        </button>
        <button
          class="vms-btn vms-btn-ghost vms-text-2xs vms-text-mono"
          style="justify-content: flex-start; padding: 6px 10px; font-size: 10px; text-align: left;"
          @click="emit('edit', instance); closeMenu()"
        >
          [EDITAR // ABRIR]
        </button>
        <div style="height: 1px; background: rgba(255,255,255,0.08); margin: 2px 0;" />
        <button
          class="vms-btn vms-btn-ghost vms-text-2xs vms-text-mono"
          style="justify-content: flex-start; padding: 6px 10px; font-size: 10px; text-align: left; color: #ff003c;"
          @click="emit('delete', instance.id); closeMenu()"
        >
          [EXCLUIR]
        </button>
      </div>
    </Teleport>
  </div>
</template>

