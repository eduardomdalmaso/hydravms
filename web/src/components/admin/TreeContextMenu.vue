<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '../../composables/useI18n'

export interface ContextMenuTarget {
  type: 'canvas' | 'folder' | 'stream'
  id?: string
  name?: string
  currentFolderId?: string | null
}

export interface SimpleFolderOption {
  id: string
  name: string
}

defineProps<{
  isOpen: boolean
  x: number
  y: number
  target: ContextMenuTarget
  folders?: SimpleFolderOption[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'action', action: string, target: ContextMenuTarget, extra?: any): void
}>()

const { t } = useI18n()
const isSubmenuOpen = ref(false)

const handleAction = (action: string, target: ContextMenuTarget, extra?: any) => {
  emit('action', action, target, extra)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-context-menu-backdrop" @click="emit('close')" @contextmenu.prevent="emit('close')">
    <div class="vms-context-menu" :style="{ left: `${x}px`, top: `${y}px` }">
      <div v-if="target.type === 'canvas'" class="vms-flex-col">
        <button class="vms-context-item" @click="handleAction('create-folder', target)"><span>{{ t('new_folder') }}</span></button>
        <button class="vms-context-item" @click="handleAction('create-stream', target)"><span>{{ t('new_stream') }}</span></button>
      </div>

      <div v-else-if="target.type === 'folder'" class="vms-flex-col">
        <div class="vms-context-menu-section-title">// PASTA: {{ target.name }}</div>
        <button class="vms-context-item" @click="handleAction('open-folder', target)"><span>[ABRIR PASTA]</span></button>
        <button class="vms-context-item" @click="handleAction('create-stream', target)"><span>[NOVO FLUXO NESTA PASTA]</span></button>
        <div class="vms-context-divider"></div>
        <button class="vms-context-item danger" @click="handleAction('delete-folder', target)">[DELETAR]</button>
      </div>

      <div v-else-if="target.type === 'stream'" class="vms-flex-col">
        <div class="vms-context-menu-section-title">// FLUXO: {{ target.name }}</div>
        <button class="vms-context-item" @click="handleAction('inspect-stream', target)">[INSPECIONAR / DETALHES]</button>
        <button class="vms-context-item" @click="handleAction('test-stream', target)">[TESTAR CONEXAO RTSP]</button>
        
        <!-- Enviar para pasta com submenu lateral aberto no hover -->
        <div class="vms-context-divider"></div>
        <div class="vms-context-submenu-wrapper" style="position: relative;" @mouseenter="isSubmenuOpen = true" @mouseleave="isSubmenuOpen = false">
          <button class="vms-context-item" style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
            <span>ENVIAR PARA PASTA</span>
            <span style="color: var(--vms-neu-accent-orange); font-size: 10px; margin-left: 8px;">▶</span>
          </button>

          <div v-if="isSubmenuOpen && folders && folders.length > 0" class="vms-context-menu" style="position: absolute; left: calc(100% + 2px); top: 0; min-width: 170px; z-index: 10000; box-shadow: 0 8px 24px rgba(0,0,0,0.8); border: 1px solid rgba(255, 94, 58, 0.4);">
            <div class="vms-context-menu-section-title">// SELECIONE A PASTA:</div>
            <button v-if="target.currentFolderId" class="vms-context-item" @click="handleAction('move-stream', target, null)"><span>[RAIZ] Sem Pasta</span></button>
            <button v-for="f in (folders || []).filter(fol => fol.id !== target.currentFolderId)" :key="f.id" class="vms-context-item" @click="handleAction('move-stream', target, f.id)">
              <span>[PASTA] {{ f.name }}</span>
            </button>
          </div>
        </div>

        <div class="vms-context-divider"></div>
        <button class="vms-context-item danger" @click="handleAction('delete-stream', target)">[{{ t('delete_stream') }}]</button>
      </div>
    </div>
  </div>
</template>
