<script setup lang="ts">
import { ref } from 'vue'
import type { UserModulePermission } from '../../../types/userTree'

const props = defineProps<{ module: UserModulePermission }>()
const emit = defineEmits<{ (e: 'toggle', msg: string): void }>()
const isExpanded = ref(false)

const handleMasterToggle = () => {
  props.module.isEnabled = !props.module.isEnabled
  props.module.permissions.forEach(p => { p.isEnabled = props.module.isEnabled })
  emit('toggle', props.module.isEnabled ? `Módulo "${props.module.name}" LIBERADO.` : `Módulo "${props.module.name}" BLOQUEADO.`)
}

const handleSubToggle = (permId: string) => {
  const perm = props.module.permissions.find(p => p.id === permId)
  if (!perm) return
  perm.isEnabled = !perm.isEnabled
  if (perm.isEnabled) props.module.isEnabled = true
  emit('toggle', `Permissão "${perm.name}" ${perm.isEnabled ? 'ATIVADA' : 'DESATIVADA'}.`)
}
</script>

<template>
  <div style="border: 1px solid var(--vms-border); border-radius: 6px; background: #0c0e14; overflow: hidden;">
    <!-- Module Bar Header -->
    <div class="vms-flex-between" style="padding: 0.6rem 0.8rem; background: #13161c; border-bottom: 1px solid rgba(255, 255, 255, 0.05);">
      <div class="vms-flex-row" style="gap: 0.6rem; align-items: center; min-width: 0;">
        <input type="checkbox" class="vms-checkbox" :checked="module.isEnabled" @change="handleMasterToggle" />
        <div class="vms-flex-col" style="gap: 2px; min-width: 0; cursor: pointer;" @click="isExpanded = !isExpanded">
          <span class="vms-font-bold" style="color: #fff; font-size: 11px;">{{ module.name }}</span>
          <span class="vms-text-dim vms-text-2xs" style="font-size: 10px;">{{ module.description }}</span>
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <span class="vms-badge" :class="module.isEnabled ? 'vms-badge-info' : 'vms-badge-warning'" style="font-size: 8px;">
          {{ module.isEnabled ? '[LIBERADO]' : '[BLOQUEADO]' }}
        </span>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 4px; font-size: 10px;" title="Detalhes do Escopo" @click="isExpanded = !isExpanded">
          <span>{{ isExpanded ? '▲' : '▼' }}</span>
        </button>
      </div>
    </div>

    <!-- Sub-permissions and Granular Scope Matrix -->
    <div v-if="isExpanded" style="padding: 0.6rem 0.8rem; display: flex; flex-direction: column; gap: 0.45rem; background: rgba(0, 0, 0, 0.2);">
      <div v-if="module.scopeTarget" class="vms-flex-between" style="padding-bottom: 0.4rem; border-bottom: 1px dashed var(--vms-border);">
        <span class="vms-text-dim vms-text-2xs">ALVOS / ESCOPO ATRIBUIDO:</span>
        <span class="vms-text-mono vms-text-2xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">
          {{ module.scopeTarget }}
        </span>
      </div>

      <div v-for="p in module.permissions" :key="p.id" class="vms-flex-between" style="padding: 0.25rem 0.4rem; border-radius: 4px; background: rgba(255, 255, 255, 0.015);">
        <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">{{ p.name }}</span>
        <input type="checkbox" class="vms-checkbox" :checked="p.isEnabled" @change="handleSubToggle(p.id)" />
      </div>
    </div>
  </div>
</template>
