<script setup lang="ts">
import { ref, h } from "vue"

export type AdminPageId = "video_streams" | "alarms" | "users" | "layouts" | "carousels" | "maps" | "workflows" | "storage"

defineProps<{ activePage: AdminPageId }>()
const emit = defineEmits<{ (e: "selectPage", page: AdminPageId): void }>()

const isCollapsed = ref(false)
const isConfigOpen = ref(true)

const makeIcon = (d: string) => () => h("svg", { width: 12, height: 12, viewBox: "0 0 24 24", fill: "none", stroke: "#ff5e3a", "stroke-width": 2, "stroke-linecap": "round", "stroke-linejoin": "round", style: { flexShrink: 0 } }, [h("path", { d })])

const configItems = [
  { id: "video_streams" as AdminPageId, label: "Fluxo de Vídeo", icon: makeIcon("M23 7l-7 5 7 5V7z M1 5h15v14H1z") },
  { id: "alarms" as AdminPageId, label: "Alarmes", icon: makeIcon("M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9 M13.73 21a2 2 0 0 1-3.46 0") },
  { id: "users" as AdminPageId, label: "Usuário", icon: makeIcon("M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2 M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8z") },
  { id: "layouts" as AdminPageId, label: "Layouts", icon: makeIcon("M3 3h7v7H3z M14 3h7v7h-7z M14 14h7v7h-7z M3 14h7v7H3z") },
  { id: "carousels" as AdminPageId, label: "Rondas", icon: makeIcon("M23 4v6h-6 M1 20v-6h6 M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15") },
  { id: "maps" as AdminPageId, label: "Mapas", icon: makeIcon("M1 6v16l7-4 8 4 7-4V2l-7 4-8-4-7 4z M8 2v16 M16 6v16") }
]
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>

    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <!-- [CONFIG] Section Header (No topo do Sidebar) -->
      <div class="vms-accordion-header" :class="{ active: isConfigOpen }" title="Clique para expandir ou recuar" @click="isConfigOpen = !isConfigOpen">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[CONFIG]</span>
        </div>
        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" :style="{ transform: isConfigOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s ease' }">
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </div>

      <!-- Config Items List -->
      <div v-if="isConfigOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
        <div
          v-for="item in configItems" :key="item.id"
          class="vms-asset-item" :class="{ active: activePage === item.id }"
          style="padding: 0.4rem 0.55rem;"
          @click="emit('selectPage', item.id)"
        >
          <div class="vms-flex-row" style="gap: 0.45rem; align-items: center; min-width: 0;">
            <component :is="item.icon" />
            <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500;">
              {{ item.label }}
            </span>
          </div>
          <span v-if="activePage === item.id" class="vms-status-led online"></span>
        </div>
      </div>

      <!-- [WORKFLOWS & ALARMES] Standalone Header -->
      <div
        class="vms-accordion-header" :class="{ active: activePage === 'workflows' }"
        style="border-top: 1px solid var(--vms-border);" title="Workflows e Automações" @click="emit('selectPage', 'workflows')"
      >
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" :style="{ color: activePage === 'workflows' ? 'var(--vms-neu-accent-orange)' : '#ffffff' }">[WORKFLOWS & ALARMES]</span>
        </div>
      </div>

      <!-- [STORAGE & DISCOS] Standalone Header -->
      <div
        class="vms-accordion-header" :class="{ active: activePage === 'storage' }"
        style="border-top: 1px solid var(--vms-border);" title="Storage e Discos" @click="emit('selectPage', 'storage')"
      >
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" :style="{ color: activePage === 'storage' ? 'var(--vms-neu-accent-orange)' : '#ffffff' }">[STORAGE & DISCOS]</span>
        </div>
      </div>
    </div>
  </aside>
</template>
