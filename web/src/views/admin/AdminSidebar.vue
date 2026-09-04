<script setup lang="ts">
import { ref } from "vue"

export type AdminPageId = "cameras" | "layouts" | "maps" | "carousels" | "users" | "workflows" | "storage"

defineProps<{ activePage: AdminPageId }>()
const emit = defineEmits<{ (e: "selectPage", page: AdminPageId): void }>()

const isCollapsed = ref(false)

const navItems: { id: AdminPageId; label: string; count?: string }[] = [
  { id: "cameras", label: "CAMERAS & FLUXOS", count: "5" },
  { id: "layouts", label: "LAYOUTS DO SISTEMA", count: "8" },
  { id: "maps", label: "MAPAS & PLANTAS", count: "2" },
  { id: "carousels", label: "RONDAS AUTOMATICAS", count: "1" },
  { id: "users", label: "USUARIOS & RBAC", count: "3" },
  { id: "workflows", label: "WORKFLOWS & ALARMES", count: "4" },
  { id: "storage", label: "STORAGE & DISCOS", count: "NVMe" }
]
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <!-- Clickable Interactive Border Line -->
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>

    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <!-- Section Header -->
      <div class="vms-accordion-header active" style="border-left: 3px solid var(--vms-neu-accent-orange);">
        <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[ADMIN CENTER]</span>
          <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.35); font-size: 8.5px; font-weight: 700; padding: 1px 5px;">
            7 MODULOS
          </span>
        </div>
      </div>

      <!-- Navigation List in VMS Sidebar Style -->
      <div style="padding: 0.5rem 0.4rem; display: flex; flex-direction: column; gap: 0.35rem;">
        <div
          v-for="item in navItems"
          :key="item.id"
          class="vms-asset-item"
          :class="{ active: activePage === item.id }"
          :style="{
            borderColor: activePage === item.id ? 'rgba(255,94,58,0.45)' : 'var(--vms-border)',
            backgroundColor: activePage === item.id ? 'rgba(255,94,58,0.12)' : 'var(--vms-neu-bg)'
          }"
          @click="emit('selectPage', item.id)"
        >
          <div class="vms-flex-row" style="gap: 0.45rem; align-items: center; min-width: 0;">
            <span
              class="vms-text-xs"
              :style="{
                color: activePage === item.id ? 'var(--vms-neu-accent-orange)' : '#ffffff',
                fontFamily: 'var(--vms-font-roboto)',
                fontWeight: activePage === item.id ? '700' : '500',
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis'
              }"
            >
              [{{ item.label }}]
            </span>
          </div>
          <span v-if="item.count" class="vms-badge" style="background: rgba(255,255,255,0.06); color: #cbd5e1; font-family: var(--vms-font-roboto); font-size: 8px; padding: 1px 4px;">
            {{ item.count }}
          </span>
        </div>
      </div>
    </div>
  </aside>
</template>
