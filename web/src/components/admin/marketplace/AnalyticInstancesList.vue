<script setup lang="ts">
import type { AnalyticInstance } from '../../../types/marketplace'

defineProps<{ instances: AnalyticInstance[] }>()
const emit = defineEmits<{
  (e: 'toggle', id: string): void
  (e: 'delete', id: string): void
  (e: 'explore', pluginId: string): void
}>()
</script>

<template>
  <div class="vms-card" style="padding: 1.25rem; display: flex; flex-direction: column; gap: 0.85rem;">
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
      <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">INSTÂNCIAS CONFIGURADAS ({{ instances.length }})</span>
      <span class="vms-text-mono vms-text-2xs vms-text-dim">STATUS // FPS // DETECÇÕES</span>
    </div>

    <div v-if="instances.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="text-align: center; padding: 3rem;">
      // NENHUMA INSTÂNCIA CRIADA. UTILIZE O FORMULÁRIO AO LADO PARA ATIVAR UM ANALÍTICO.
    </div>

    <div v-else class="vms-flex-col" style="gap: 0.65rem; max-height: 540px; overflow-y: auto;">
      <div
        v-for="inst in instances"
        :key="inst.id"
        class="vms-card"
        style="padding: 0.85rem; background: var(--vms-neu-bg); border: 1px solid var(--vms-border); display: flex; flex-direction: column; gap: 0.5rem;"
      >
        <div class="vms-flex-between" style="align-items: flex-start;">
          <div class="vms-flex-col" style="gap: 2px;">
            <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">{{ inst.name }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ inst.camera_name }} // {{ inst.plugin_name }}</span>
          </div>
          <span class="vms-badge" :class="inst.is_active ? 'vms-badge-green' : 'vms-badge-orange'" style="font-size: 10px;">
            {{ inst.is_active ? '[ATIVO]' : '[PAUSADO]' }}
          </span>
        </div>

        <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 4px 8px; border-radius: 4px;">
          <div class="vms-flex-row" style="gap: 0.85rem; font-size: 11px; font-family: var(--vms-font-jetbrains);">
            <span>TAXA: <strong style="color: #00f0ff;">{{ inst.fps_rate }} FPS</strong></span>
            <span>DETECÇÕES: <strong style="color: #00ff9d;">{{ inst.detections_count }}</strong></span>
            <span>TARGET: <strong style="color: #ff5e3a;">{{ inst.hardware_target === 'rtx_5090_cuda' ? 'RTX 5090' : 'CPU SHM' }}</strong></span>
          </div>

          <div class="vms-flex-row" style="gap: 0.4rem;">
            <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 3px 8px;" @click="emit('explore', inst.plugin_id)">
              DADOS
            </button>
            <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 3px 8px;" @click="emit('toggle', inst.id)">
              {{ inst.is_active ? 'PAUSAR' : 'ATIVAR' }}
            </button>
            <button class="vms-btn vms-btn-danger vms-btn-sm" style="font-size: 10px; padding: 3px 8px;" @click="emit('delete', inst.id)">
              EXCLUIR
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
