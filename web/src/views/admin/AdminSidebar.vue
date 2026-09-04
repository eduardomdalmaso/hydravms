<script setup lang="ts">
import { ref } from "vue"

export type AdminPageId = "video_streams" | "alarms" | "users" | "layouts" | "carousels" | "maps" | "workflows" | "storage"

defineProps<{ activePage: AdminPageId }>()
const emit = defineEmits<{ (e: "selectPage", page: AdminPageId): void }>()

const isCollapsed = ref(false)
const isConfigOpen = ref(true)

const configItems: { id: AdminPageId; label: string }[] = [
  { id: "video_streams", label: "FLUXO DE VIDEO" },
  { id: "alarms", label: "ALARMES" },
  { id: "users", label: "USUARIO" },
  { id: "layouts", label: "LAYOUTS" },
  { id: "carousels", label: "RONDAS" },
  { id: "maps", label: "MAPAS" }
]

const mainItems: { id: AdminPageId; label: string }[] = [
  { id: "workflows", label: "WORKFLOWS & ALARMES" },
  { id: "storage", label: "STORAGE & DISCOS" }
]
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>

    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <!-- Section Header: CONFIG -->
      <div class="vms-accordion-header" :class="{ active: isConfigOpen }" title="Clique para expandir/recuar" @click="isConfigOpen = !isConfigOpen">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[CONFIG]</span>
        </div>
      </div>

      <!-- Config Items List -->
      <div v-if="isConfigOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem;">
        <div
          v-for="item in configItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }"
          :style="{ borderColor: activePage === item.id ? 'rgba(255,94,58,0.45)' : 'var(--vms-border)', backgroundColor: activePage === item.id ? 'rgba(255,94,58,0.12)' : 'var(--vms-neu-bg)' }"
          @click="emit('selectPage', item.id)"
        >
          <span class="vms-text-xs" :style="{ color: activePage === item.id ? 'var(--vms-neu-accent-orange)' : '#ffffff', fontFamily: 'var(--vms-font-roboto)', fontWeight: activePage === item.id ? '700' : '500' }">
            [{{ item.label }}]
          </span>
        </div>
      </div>

      <!-- Standalone Modules -->
      <div style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; border-top: 1px solid rgba(255,255,255,0.06); margin-top: 0.2rem;">
        <div
          v-for="item in mainItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }"
          :style="{ borderColor: activePage === item.id ? 'rgba(255,94,58,0.45)' : 'var(--vms-border)', backgroundColor: activePage === item.id ? 'rgba(255,94,58,0.12)' : 'var(--vms-neu-bg)' }"
          @click="emit('selectPage', item.id)"
        >
          <span class="vms-text-xs" :style="{ color: activePage === item.id ? 'var(--vms-neu-accent-orange)' : '#ffffff', fontFamily: 'var(--vms-font-roboto)', fontWeight: activePage === item.id ? '700' : '500' }">
            [{{ item.label }}]
          </span>
        </div>
      </div>
    </div>
  </aside>
</template>
