<script setup lang="ts">
import type { LayoutFolderNode, EnterpriseLayoutItem } from '../../../types/layoutTree'

defineProps<{
  currentFolder?: LayoutFolderNode | null
  selectedLayout?: EnterpriseLayoutItem | null
  totalFolders: number
  totalLayouts: number
  currentFolderLayoutsCount?: number
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'navigate-root'): void
  (e: 'navigate-folder'): void
  (e: 'save'): void
  (e: 'toggleLock'): void
}>()
</script>

<template>
  <div class="vms-flex-between" style="padding: 0.5rem 0.75rem; background: #16191f; border-radius: 6px; border: 1px solid var(--vms-border); align-items: center;">
    <!-- Left path breadcrumbs -->
    <div class="vms-flex-row" style="gap: 0.75rem; align-items: center; min-width: 0;">
      <button v-if="currentFolder || selectedLayout" class="vms-btn vms-btn-secondary vms-btn-sm" title="Voltar" @click="emit('back')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
        <span>VOLTAR</span>
      </button>

      <span
        class="vms-text-mono vms-text-xs"
        :style="{ color: !currentFolder && !selectedLayout ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }"
        style="cursor: pointer;"
        @click="emit('navigate-root')"
      >
        [EMPRESAS & CLIENTES]
      </span>

      <template v-if="currentFolder">
        <span class="vms-text-dim">/</span>
        <span
          class="vms-font-bold"
          :style="{ color: !selectedLayout ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)', cursor: selectedLayout ? 'pointer' : 'default' }"
          style="font-size: 12px;"
          @click="selectedLayout ? emit('navigate-folder') : undefined"
        >
          {{ currentFolder.name }}
        </span>
      </template>

      <template v-if="selectedLayout">
        <span class="vms-text-dim">/</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {{ selectedLayout.name }}
        </span>
      </template>
    </div>

    <!-- Right side: Counters + Padlock & SALVAR button at far right -->
    <div class="vms-flex-row" style="gap: 0.5rem; align-items: center; flex-shrink: 0;">
      <span class="vms-text-mono vms-text-2xs vms-text-dim">
        {{ selectedLayout ? '[INSPEÇÃO // GRADE ATIVA]' : currentFolder ? `${currentFolderLayoutsCount ?? 0} LAYOUTS` : `${totalFolders} PASTAS // ${totalLayouts} LAYOUTS` }}
      </span>

      <template v-if="selectedLayout">
        <button
          type="button"
          class="vms-btn vms-btn-sm"
          :class="selectedLayout.is_locked ? 'vms-btn-primary' : 'vms-btn-secondary'"
          :style="{
            padding: '4px 8px',
            background: selectedLayout.is_locked ? 'rgba(255, 94, 58, 0.2)' : '#07080c',
            borderColor: selectedLayout.is_locked ? '#ff5e3a' : 'var(--vms-border)'
          }"
          :title="selectedLayout.is_locked ? '[TRAVADO COM CADEADO] (Clique para destravar)' : '[DESTRAVADO] (Clique para travar com cadeado)'"
          @click="emit('toggleLock')"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" :fill="selectedLayout.is_locked ? '#ff5e3a' : 'none'" stroke="#ff5e3a" stroke-width="1.5">
            <path d="M12 2C9.24 2 7 4.24 7 7v3H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2h-1V7c0-2.76-2.24-5-5-5zm-3 5c0-1.66 1.34-3 3-3s3 1.34 3 3v3H9V7zm3 7a1.5 1.5 0 0 1 1 1.37V17a1 1 0 1 1-2 0v-1.63A1.5 1.5 0 0 1 12 14z"/>
          </svg>
        </button>

        <button
          class="vms-btn vms-btn-primary vms-btn-sm"
          style="padding: 4px 16px; font-weight: 700; box-shadow: 0 0 12px rgba(255, 94, 58, 0.4);"
          @click="emit('save')"
        >
          SALVAR
        </button>
      </template>
    </div>
  </div>
</template>
