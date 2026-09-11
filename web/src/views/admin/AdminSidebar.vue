<script setup lang="ts">
import { ref, h } from "vue"
import { useMarketplace } from "../../composables/useMarketplace"
import { useI18n } from "../../composables/useI18n"
export type AdminPageId = string

defineProps<{ activePage: AdminPageId }>()
const emit = defineEmits<{ (e: "selectPage", page: AdminPageId): void }>()
const isCollapsed = ref(false), isConfigOpen = ref(true), isMarketOpen = ref(true), isSysOpen = ref(true)
const openPlugins = ref<Record<string, boolean>>({})
const { installedPlugins } = useMarketplace()
const { t } = useI18n()

const makeIcon = (d: string) => () => h("svg", { width: 12, height: 12, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", "stroke-width": 2, "stroke-linecap": "round", "stroke-linejoin": "round", style: { flexShrink: 0 } }, [h("path", { d })])

const configItems = [
  { id: "video_streams", icon: makeIcon("M23 7l-7 5 7 5V7z M1 5h15v14H1z") },
  { id: "alarms", icon: makeIcon("M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9 M13.73 21a2 2 0 0 1-3.46 0") },
  { id: "workflows", icon: makeIcon("M22 11.08V12a10 10 0 1 1-5.93-9.14 M22 4L12 14.01l-3-3") },
  { id: "users", icon: makeIcon("M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2 M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8z") },
  { id: "layouts", icon: makeIcon("M3 3h7v7H3z M14 3h7v7h-7z M14 14h7v7h-7z M3 14h7v7H3z") },
  { id: "carousels", icon: makeIcon("M23 4v6h-6 M1 20v-6h6 M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15") },
  { id: "maps", icon: makeIcon("M1 6v16l7-4 8 4 7-4V2l-7 4-8-4-7 4z M8 2v16 M16 6v16") }
]
const storageItems = [
  { id: "storage", icon: makeIcon("M4 6h16M4 12h16M4 18h16") },
  { id: "performance", icon: makeIcon("M13 2L3 14h9l-1 8 10-12h-9l1-8z") },
  { id: "logs", icon: makeIcon("M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6h6 M16 13H8 M16 17H8 M10 9H8") }
]
const iconMarket = makeIcon("M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z M3 6h18 M16 10a4 4 0 0 1-8 0")
const iconLayers = makeIcon("M12 2L2 7l10 5 10-5-10-5z M2 17l10 5 10-5 M2 12l10 5 10-5")
const iconActivity = makeIcon("M22 12h-4l-3 9L9 3l-3 9H2")

const togglePlugin = (id: string) => { openPlugins.value[id] = openPlugins.value[id] === undefined ? false : !openPlugins.value[id] }
const isPluginOpen = (id: string) => openPlugins.value[id] !== false
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed"><div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div></div>
    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; justify-content: space-between;">
      <!-- Seção Superior com Scroll: Analíticos & Plugins Instalados -->
      <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column;">
        <template v-for="p in installedPlugins" :key="p.id">
          <div class="vms-accordion-header" :class="{ active: isPluginOpen(p.id) }" @click="togglePlugin(p.id)">
            <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange); max-width: 190px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">[{{ p.name.toUpperCase() }}]</span>
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" :style="{ transform: isPluginOpen(p.id) ? 'rotate(180deg)' : 'rotate(0deg)' }"><polyline points="6 9 12 15 18 9" /></svg>
          </div>
          <div v-if="isPluginOpen(p.id)" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
            <div class="vms-asset-item" :class="{ active: activePage === `plugin_${p.id}_instances` }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', `plugin_${p.id}_instances`)">
              <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;"><component :is="iconLayers" /><span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ t('analytics') }}</span></div>
            </div>
            <div class="vms-asset-item" :class="{ active: activePage === `plugin_${p.id}_events` }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', `plugin_${p.id}_events`)">
              <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;"><component :is="iconActivity" /><span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ t('events') }}</span></div>
            </div>
          </div>
        </template>
      </div>

      <!-- Seção Fixa Inferior: Configuração, Sistema e Marketplace -->
      <div style="flex-shrink: 0; border-top: 1px solid rgba(255, 255, 255, 0.12); background: #07090e; display: flex; flex-direction: column;">
        <!-- [CONFIGURAÇÃO] -->
        <div class="vms-accordion-header" :class="{ active: isConfigOpen }" @click="isConfigOpen = !isConfigOpen">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">{{ t('configuration') }}</span>
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" :style="{ transform: isConfigOpen ? 'rotate(180deg)' : 'rotate(0deg)' }"><polyline points="6 9 12 15 18 9" /></svg>
        </div>
        <div v-if="isConfigOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
          <div v-for="item in configItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', item.id)">
            <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;"><component :is="item.icon" /><span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ t(item.id) }}</span></div>
          </div>
        </div>

        <!-- [SISTEMA] -->
        <div class="vms-accordion-header" :class="{ active: isSysOpen }" @click="isSysOpen = !isSysOpen">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">{{ t('system') }}</span>
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" :style="{ transform: isSysOpen ? 'rotate(180deg)' : 'rotate(0deg)' }"><polyline points="6 9 12 15 18 9" /></svg>
        </div>
        <div v-if="isSysOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
          <div v-for="item in storageItems" :key="item.id" class="vms-asset-item" :class="{ active: activePage === item.id }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', item.id)">
            <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;"><component :is="item.icon" /><span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ t(item.id) }}</span></div>
          </div>
        </div>

        <!-- [MARKETPLACE] -->
        <div class="vms-accordion-header" :class="{ active: isMarketOpen }" @click="isMarketOpen = !isMarketOpen">
          <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">{{ t('marketplace') }}</span>
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" :style="{ transform: isMarketOpen ? 'rotate(180deg)' : 'rotate(0deg)' }"><polyline points="6 9 12 15 18 9" /></svg>
        </div>
        <div v-if="isMarketOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14;">
          <div class="vms-asset-item" :class="{ active: activePage === 'marketplace' }" style="padding: 0.4rem 0.55rem;" @click="emit('selectPage', 'marketplace')">
            <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;"><component :is="iconMarket" /><span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ t('marketplace') }}</span></div>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
