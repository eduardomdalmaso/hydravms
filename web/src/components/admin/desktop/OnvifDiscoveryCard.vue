<script setup lang="ts">
import type { DiscoveredOnvifCamera } from '../../../types/onvifDiscovery'
import { useI18n } from '../../../composables/useI18n'

defineProps<{ camera: DiscoveredOnvifCamera }>()
const emit = defineEmits<{ (e: 'import', camera: DiscoveredOnvifCamera): void }>()
const { t } = useI18n()
</script>

<template>
  <div class="vms-asset-card" style="background: #14171d; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-sm); padding: 0.85rem; display: flex; flex-direction: column; justify-content: space-between; gap: 0.6rem; box-sizing: border-box;">
    <!-- Header: Name, Model & Clean ONVIF Tag -->
    <div class="vms-flex-between" style="align-items: flex-start; gap: 0.5rem;">
      <div class="vms-flex-col" style="gap: 2px; min-width: 0; flex: 1;">
        <span class="vms-text-xs vms-font-semibold" style="color: #ffffff; font-family: var(--vms-font-roboto); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {{ camera.name }}
        </span>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {{ camera.manufacturer }} // {{ camera.model }}
        </span>
      </div>
      <div class="vms-badge" style="background: rgba(255, 255, 255, 0.04); color: #8b94a0; border: 1px solid rgba(255, 255, 255, 0.08); font-size: 9px; padding: 2px 6px; flex-shrink: 0;">
        <span style="width: 5px; height: 5px; border-radius: 50%; background: #00ff9d; display: inline-block; margin-right: 4px;"></span>
        PROFILE S/T
      </div>
    </div>

    <!-- Metadata Grid: IP, MAC and PTZ Tag -->
    <div style="display: flex; justify-content: space-between; align-items: center; font-family: var(--vms-font-jetbrains); font-size: 10px; color: #8b94a0; background: #0b0e14; padding: 0.35rem 0.55rem; border-radius: 4px; border: 1px solid rgba(255, 255, 255, 0.04);">
      <span style="white-space: nowrap;">IP: <strong style="color: #fff;">{{ camera.ip }}:{{ camera.port }}</strong></span>
      <span style="white-space: nowrap;">MAC: {{ camera.macAddress }}</span>
      <span v-if="camera.hasPtz" style="color: var(--vms-neu-accent-yellow); font-weight: 600; font-size: 9px; white-space: nowrap;">[PTZ]</span>
    </div>

    <!-- Footer: Detection count and Action Button -->
    <div class="vms-flex-between" style="align-items: center; padding-top: 0.15rem;">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim); white-space: nowrap;">
        {{ camera.profiles.length }} {{ camera.profiles.length > 1 ? t('rtsp_profiles') : t('rtsp_profile') }}
      </span>
      <button 
        class="vms-btn vms-btn-primary vms-btn-sm" 
        style="padding: 0 12px; font-size: 10px; height: 26px; white-space: nowrap;"
        @click="emit('import', camera)"
      >
        <span>+ {{ t('import_stream') }}</span>
      </button>
    </div>
  </div>
</template>
