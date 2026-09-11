<script setup lang="ts">
import { ref } from 'vue'
import { useOnvifDiscovery } from '../../../composables/useOnvifDiscovery'
import { useI18n } from '../../../composables/useI18n'
import type { DiscoveredOnvifCamera } from '../../../types/onvifDiscovery'
import OnvifDiscoveryCard from './OnvifDiscoveryCard.vue'

const emit = defineEmits<{ (e: 'importOnvif', streamData: any, folderId?: string): void }>()
const { isScanning, filteredDevices, availableCount, lastScanTime, searchFilter, scanNetwork, convertToStreamItem, markAsImported } = useOnvifDiscovery()
const { t } = useI18n()
const isCollapsed = ref(false)

const handleImport = (cam: DiscoveredOnvifCamera) => {
  const streamData = convertToStreamItem(cam)
  markAsImported(cam.id)
  emit('importOnvif', streamData)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.75rem; background: #0e1117; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-md); padding: 0.85rem 1rem;">
    <div class="vms-flex-between" style="align-items: center; border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <div class="vms-flex-row" style="gap: 0.6rem; align-items: center;">
        <span class="vms-badge" style="background: rgba(255, 94, 58, 0.12); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.3); font-size: 10px; font-weight: 700;">
          {{ t('onvif_radar') }}
        </span>
        <h4 class="vms-h4" style="color: #fff; margin: 0; font-size: 13px; font-family: var(--vms-font-roboto); font-weight: 600;">
          {{ t('onvif_detected') }}
        </h4>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">
          ({{ availableCount }} {{ t('new_available') }} // {{ t('last_scan') }}: {{ lastScanTime || t('never') }})
        </span>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <input 
          v-model="searchFilter" 
          class="vms-auth-input" 
          style="width: 150px; font-size: 11px; padding: 2px 8px; height: 26px;" 
          :placeholder="t('filter_ip_name')" 
        />
        <button 
          class="vms-btn vms-btn-secondary vms-btn-sm" 
          :disabled="isScanning"
          style="height: 26px; padding: 0 10px; font-size: 10px;"
          @click="scanNetwork"
        >
          <span v-if="isScanning" style="color: var(--vms-neu-accent-cyan);">{{ t('scanning') }}</span>
          <span v-else>↻ {{ t('scan_network') }}</span>
        </button>
        <button 
          class="vms-btn vms-btn-ghost vms-btn-sm" 
          style="height: 26px; padding: 0 6px;"
          @click="isCollapsed = !isCollapsed"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" :style="{ transform: isCollapsed ? 'rotate(180deg)' : 'rotate(0deg)' }">
            <polyline points="6 9 12 15 18 9" />
          </svg>
        </button>
      </div>
    </div>

    <div v-if="!isCollapsed" class="vms-flex-col" style="gap: 0.75rem;">
      <div v-if="filteredDevices.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="padding: 1rem; text-align: center;">
        {{ t('no_onvif_found') }}
      </div>
      <div v-else style="display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 0.75rem;">
        <OnvifDiscoveryCard 
          v-for="cam in filteredDevices" 
          :key="cam.id" 
          :camera="cam" 
          @import="handleImport" 
        />
      </div>
    </div>
  </div>
</template>
