<script setup lang="ts">
import { ref, h } from "vue"
export type AdminPageId = "video_streams" | "alarms" | "workflows" | "users" | "layouts" | "carousels" | "maps" | "storage" | "performance" | "logs"

defineProps<{ activePage: AdminPageId }>()
const emit = defineEmits<{ (e: "selectPage", page: AdminPageId): void }>()
const isCollapsed = ref(false), isConfigOpen = ref(true), isSystemOpen = ref(true)

const makeIcon = (d: string) => () => h("svg", { width: 12, height: 12, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", "stroke-width": 2, "stroke-linecap": "round", "stroke-linejoin": "round", style: { flexShrink: 0 } }, [h("path", { d })])

const configItems = [
  { id: "video_streams" as AdminPageId, label: "Fluxo de Vídeo", icon: makeIcon("M23 7l-7 5 7 5V7z M1 5h15v14H1z") },
  { id: "alarms" as AdminPageId, label: "Alarmes", icon: makeIcon("M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9 M13.73 21a2 2 0 0 1-3.46 0") },
  { id: "workflows" as AdminPageId, label: "Workflows", icon: makeIcon("M22 11.08V12a10 10 0 1 1-5.93-9.14 M22 4L12 14.01l-3-3") },
  { id: "users" as AdminPageId, label: "Usuários", icon: makeIcon("M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2 M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8z") },
  { id: "layouts" as AdminPageId, label: "Layouts", icon: makeIcon("M3 3h7v7H3z M14 3h7v7h-7z M14 14h7v7h-7z M3 14h7v7H3z") },
  { id: "carousels" as AdminPageId, label: "Rondas", icon: makeIcon("M23 4v6h-6 M1 20v-6h6 M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15") },
  { id: "maps" as AdminPageId, label: "Mapas & Plantas", icon: makeIcon("M1 6v16l7-4 8 4 7-4V2l-7 4-8-4-7 4z M8 2v16 M16 6v16") }
]
const storageItems = [
  { id: "storage" as AdminPageId, label: "Storage & Discos", icon: makeIcon("M4 6h16M4 12h16M4 18h16") },
  { id: "performance" as AdminPageId, label: "Performance", icon: makeIcon("M13 2L3 14h9l-1 8 10-12h-9l1-8z") },
  { id: "logs" as AdminPageId, label: "Logs & Auditoria", icon: makeIcon("M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6h6 M16 13H8 M16 17H8 M10 9H8") }
]
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>
    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <!-- [CONFIGURAÇÃO] Section -->
      <div class="vms-accordion-header" :class="{ active: isConfigOpen }" title="Clique para expandir ou recuar" @click="isConfigOpen = !isConfigOpen">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[CONFIGURAÇÃO]</span>
        </div>
        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" :style="{ transform: isConfigOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s ease' }">
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </div>
      <div v-if="isConfigOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
        <div v-for="item in configItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', item.id)">
          <div class="vms-flex-row" style="gap: 0.45rem; align-items: center; min-width: 0;">
            <component :is="item.icon" />
            <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ item.label }}</span>
          </div>
        </div>
      </div>

      <!-- [SISTEMA] Section -->
      <div class="vms-accordion-header" :class="{ active: isSystemOpen }" title="Clique para expandir ou recuar" @click="isSystemOpen = !isSystemOpen">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[SISTEMA]</span>
        </div>
        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" :style="{ transform: isSystemOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s ease' }">
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </div>
      <div v-if="isSystemOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
        <div v-for="item in storageItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', item.id)">
          <div class="vms-flex-row" style="gap: 0.45rem; align-items: center; min-width: 0;">
            <component :is="item.icon" />
            <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ item.label }}</span>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
