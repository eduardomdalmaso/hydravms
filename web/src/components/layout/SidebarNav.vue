<script setup lang="ts">
import { useI18n } from '../../composables/useI18n'

const props = defineProps<{
  activeTab: string
  username: string
}>()

const emit = defineEmits<{
  (e: 'navigate', tab: string): void
  (e: 'logout'): void
}>()

const { t } = useI18n()

const navItems = [
  { id: 'mosaic', key: 'mosaic_live' },
  { id: 'layout_designer', key: 'create_layouts' },
  { id: 'workflows', key: 'workflows_alarms' },
  { id: 'storage', key: 'storage_disks' }
]
</script>

<template>
  <aside class="vms-sidebar" style="background: var(--vms-neu-surface); border-right: 1px solid var(--vms-border); box-shadow: 6px 0 16px rgba(0,0,0,0.4);">
    <div class="vms-sidebar-brand" style="border-bottom: 1px solid var(--vms-border); padding: 1.25rem 1rem;">
      <div class="vms-flex-col" style="gap: 0.25rem;">
        <span style="color: #ffffff; font-family: var(--vms-font-roboto); font-size: 15px; font-weight: 800; letter-spacing: 0.8px;">HYDRA VMS</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">MASTER CONSOLE</span>
      </div>
    </div>
    <nav class="vms-sidebar-nav" style="padding: 1rem 0.75rem; gap: 0.5rem;">
      <button
        v-for="item in navItems"
        :key="item.id"
        class="vms-btn vms-btn-sm"
        :style="{
          width: '100%',
          justifyContent: 'flex-start',
          backgroundColor: activeTab === item.id ? 'var(--vms-neu-bg)' : 'transparent',
          color: activeTab === item.id ? 'var(--vms-neu-accent-orange)' : '#ffffff',
          fontFamily: 'var(--vms-font-roboto)',
          fontWeight: '500',
          boxShadow: activeTab === item.id ? 'var(--vms-neu-inset-dark), var(--vms-neu-inset-light)' : 'none',
          border: activeTab === item.id ? '1px solid rgba(255,94,58,0.25)' : '1px solid transparent',
          padding: '0.65rem 0.85rem'
        }"
        @click="emit('navigate', item.id)"
      >
        {{ t(item.key) }}
      </button>
    </nav>
    <div class="vms-sidebar-footer" style="padding: 1rem; border-top: 1px solid var(--vms-border); gap: 0.75rem;">
      <div class="vms-flex-col" style="gap: 0.125rem;">
        <span class="vms-text-2xs vms-text-muted" style="font-family: var(--vms-font-roboto);">{{ t('active_user') }}</span>
        <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ username }}</span>
      </div>
      <button class="vms-btn vms-btn-secondary vms-btn-sm" style="width: 100%; font-family: var(--vms-font-roboto);" @click="emit('logout')">
        {{ t('disconnect') }}
      </button>
    </div>
  </aside>
</template>
