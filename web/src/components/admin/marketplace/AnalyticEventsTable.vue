<script setup lang="ts">
import type { AnalyticEventRecord } from '../../../types/marketplace'

defineProps<{ events: AnalyticEventRecord[] }>()
const emit = defineEmits<{ (e: 'inspect', event: AnalyticEventRecord): void }>()
</script>

<template>
  <div class="vms-card" style="padding: 0; overflow: hidden;">
    <div style="overflow-x: auto; max-height: 580px;">
      <table class="vms-table" style="width: 100%; font-size: 11px;">
        <thead>
          <tr>
            <th>HORÁRIO (UTC)</th>
            <th>CÂMERA</th>
            <th>ANALÍTICO</th>
            <th>IDENTIFICAÇÃO</th>
            <th>CONFIANÇA</th>
            <th>SEVERIDADE</th>
            <th>DETALHES</th>
            <th style="text-align: right;">AÇÃO</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="events.length === 0">
            <td colspan="8" class="vms-text-mono vms-text-dim" style="text-align: center; padding: 2.5rem;">
              // NENHUM EVENTO ENCONTRADO COM OS FILTROS APLICADOS
            </td>
          </tr>
          <tr v-for="evt in events" :key="evt.id" style="cursor: pointer;" @click="emit('inspect', evt)">
            <td class="vms-text-mono" style="color: #94a3b8;">{{ evt.timestamp.replace('T', ' ').replace('Z', '') }}</td>
            <td class="vms-font-semibold" style="color: #ffffff;">{{ evt.camera_name }}</td>
            <td style="color: var(--vms-neu-accent-orange);">{{ evt.plugin_name }}</td>
            <td>
              <span class="vms-badge vms-badge-secondary" style="font-family: var(--vms-font-jetbrains); font-size: 11px;">
                {{ evt.object_label }}
              </span>
            </td>
            <td class="vms-text-mono" style="color: #00f0ff; font-weight: bold;">
              {{ Math.round(evt.confidence * 100) }}%
            </td>
            <td>
              <span
                class="vms-badge"
                :class="evt.severity === 'critical' ? 'vms-badge-danger' : (evt.severity === 'warning' ? 'vms-badge-orange' : 'vms-badge-blue')"
                style="font-size: 10px;"
              >
                [{{ evt.severity.toUpperCase() }}]
              </span>
            </td>
            <td class="vms-text-dim" style="max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
              {{ evt.details }}
            </td>
            <td style="text-align: right;">
              <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 2px 8px;" @click.stop="emit('inspect', evt)">
                VER DETALHES
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
