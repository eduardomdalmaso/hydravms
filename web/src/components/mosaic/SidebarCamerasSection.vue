<script setup lang="ts">
import { ref, computed } from "vue"
import { useI18n } from "../../composables/useI18n"
import type { CameraStreamInfo, StreamProtocol } from "../../types/mosaic"

const props = defineProps<{ cameras: CameraStreamInfo[]; isOpen: boolean }>()
const emit = defineEmits<{
  (e: "toggle"): void
  (e: "selectCamera", camera: CameraStreamInfo): void
}>()
const { t } = useI18n()

const isSearchOpen = ref(false)
const searchQuery = ref("")
const selectedProtocol = ref<StreamProtocol | "ALL">("ALL")

const filteredCameras = computed(() => {
  return props.cameras.filter(c => {
    const matchesName = c.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesProto = selectedProtocol.value === "ALL" || c.protocol === selectedProtocol.value
    return matchesName && matchesProto
  })
})
</script>

<template>
  <div>
    <!-- Clickable Header Row with Orange Counter -->
    <div class="vms-accordion-header" :class="{ active: isOpen }" title="Clique na linha para expandir ou recuar" @click="emit('toggle')">
      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
        <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[{{ t('cameras').toUpperCase() }}]</span>
        <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.35); font-size: 8.5px; font-weight: 700; padding: 1px 5px;">
          {{ filteredCameras.length }}
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <button class="vms-sidebar-icon-btn" :title="t('search_cam')" @click.stop="isSearchOpen = !isSearchOpen">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div v-if="isOpen && isSearchOpen" style="padding: 0.4rem 0.5rem; background: #11141b; border-bottom: 1px solid var(--vms-border); display: flex; flex-direction: column; gap: 0.3rem;">
      <input v-model="searchQuery" class="vms-auth-input" style="padding: 0.25rem 0.5rem; font-size: 11px; height: 26px;" :placeholder="t('search_cam')" autofocus />
      <div class="vms-flex-row" style="gap: 0.25rem;">
        <button v-for="p in (['ALL', 'RTSP', 'RTMP', 'ONVIF'] as const)" :key="p" class="vms-btn vms-btn-sm" :class="selectedProtocol === p ? 'vms-btn-primary' : 'vms-btn-ghost'" style="font-size: 9px; padding: 1px 4px;" @click="selectedProtocol = p">
          {{ p === 'ALL' ? t('all') : p }}
        </button>
      </div>
    </div>

    <!-- Camera List with LED Status & White Protocols in Roboto -->
    <div v-if="isOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14; max-height: 200px; overflow-y: auto;">
      <div v-for="cam in filteredCameras" :key="cam.id" class="vms-asset-item" style="padding: 0.4rem 0.55rem;" @click="emit('selectCamera', cam)">
        <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 120px;">
          {{ cam.name }}
        </span>
        <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(255, 255, 255, 0.08); color: #ffffff; border: 1px solid rgba(255, 255, 255, 0.15); font-family: var(--vms-font-roboto); font-weight: 500; font-size: 8px; padding: 1px 4px;">
            {{ cam.protocol || 'RTSP' }}
          </span>
          <!-- LED Indicator (Green: Online, Red: Offline) -->
          <span class="vms-status-led" :class="cam.status === 'offline' ? 'offline' : 'online'" :title="cam.status === 'offline' ? 'Offline' : 'Online'"></span>
        </div>
      </div>
    </div>
  </div>
</template>
