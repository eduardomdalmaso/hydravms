<script setup lang="ts">
import type { PluginManifest } from '../../../types/marketplace'

defineProps<{ plugin: PluginManifest }>()
const emit = defineEmits<{
  (e: 'install', id: string): void
  (e: 'uninstall', id: string): void
  (e: 'toggle', id: string): void
  (e: 'details', plugin: PluginManifest): void
}>()
</script>

<template>
  <div class="vms-plugin-card" :class="{ installed: plugin.is_installed }">
    <div class="vms-flex-col" style="gap: 0.5rem;">
      <div class="vms-flex-between" style="align-items: flex-start; gap: 0.5rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">{{ plugin.name }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">v{{ plugin.version }} // {{ plugin.author }}</span>
        </div>
        <span class="vms-badge" :class="plugin.is_installed ? (plugin.status === 'running' ? 'vms-badge-green' : 'vms-badge-orange') : 'vms-badge-blue'" style="font-size: 10px;">
          {{ plugin.is_installed ? (plugin.status === 'running' ? '[INSTALADO]' : '[PAUSADO]') : '[DISPONÍVEL]' }}
        </span>
      </div>

      <p class="vms-text-xs vms-text-dim" style="margin: 0; line-height: 1.4; height: 36px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;">
        {{ plugin.description }}
      </p>

      <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap; margin-top: 0.25rem;">
        <span class="vms-badge vms-badge-secondary" style="font-size: 9px; font-family: var(--vms-font-jetbrains);">{{ plugin.hardware_req }}</span>
        <span v-if="plugin.is_installed" class="vms-badge vms-badge-orange" style="font-size: 9px;">MÓDULO ATIVO NO MENU</span>
      </div>
    </div>

    <!-- Actions Area -->
    <div class="vms-flex-between" style="align-items: center; border-top: 1px solid var(--vms-border); padding-top: 0.65rem;">
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('details', plugin)">
        MANIFESTO
      </button>

      <div class="vms-flex-row" style="gap: 0.35rem;">
        <template v-if="plugin.is_installed">
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('toggle', plugin.id)">
            {{ plugin.status === 'running' ? 'PAUSAR' : 'ATIVAR' }}
          </button>
          <button class="vms-btn vms-btn-danger vms-btn-sm" style="font-size: 10px; padding: 4px 8px;" @click="emit('uninstall', plugin.id)">
            DESINSTALAR
          </button>
        </template>
        <template v-else>
          <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-size: 11px; padding: 4px 14px;" :disabled="plugin.status === 'updating'" @click="emit('install', plugin.id)">
            {{ plugin.status === 'updating' ? 'BAIXANDO...' : 'INSTALAR' }}
          </button>
        </template>
      </div>
    </div>
  </div>
</template>
