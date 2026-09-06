<script setup lang="ts">
import type { LogEntry } from '../../../types/systemLogs'

defineProps<{ logs: LogEntry[] }>()
const emit = defineEmits<{ (e: 'selectLog', log: LogEntry): void }>()
</script>

<template>
  <div class="vms-table-wrapper" style="overflow-x: auto; background: #0c0e14; border: 1px solid var(--vms-border); border-radius: var(--vms-radius-lg);">
    <table class="vms-table" style="width: 100%; border-collapse: collapse; text-align: left;">
      <thead>
        <tr style="border-bottom: 1px solid var(--vms-border); background: #07080c;">
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">TIMESTAMP</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">NIVEL</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">CATEGORIA</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">ATOR / ORIGEM</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">TENANT</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">ACAO</th>
          <th class="vms-th" style="padding: 10px 12px; font-size: 0.7rem; color: #ff5e3a;">DETALHES & IP</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="log in logs"
          :key="log.id"
          class="vms-log-row"
          style="border-bottom: 1px solid rgba(255, 255, 255, 0.05); cursor: pointer;"
          @click="emit('selectLog', log)"
        >
          <td class="vms-td" style="padding: 9px 12px; font-family: var(--vms-font-jetbrains); font-size: 0.72rem; color: #ffffff; white-space: nowrap;">
            {{ log.timestamp }}
          </td>
          <td class="vms-td" style="padding: 9px 12px; white-space: nowrap;">
            <span
              class="vms-badge"
              :class="{
                'vms-badge-orange': log.level === 'CRITICAL',
                'vms-badge-outline': log.level === 'WARNING',
                'vms-badge-dim': log.level === 'INFO'
              }"
              style="font-size: 0.65rem; font-weight: 700; padding: 2px 6px;"
            >
              [{{ log.level }}]
            </span>
          </td>
          <td class="vms-td" style="padding: 9px 12px; white-space: nowrap;">
            <span class="vms-text-mono vms-text-2xs" :style="{ color: log.category === 'SYSTEM' ? '#ff5e3a' : '#ffffff' }">
              [{{ log.category === 'SYSTEM' ? 'SISTEMA' : 'AUDITORIA' }}]
            </span>
          </td>
          <td class="vms-td" style="padding: 9px 12px; font-size: 0.75rem; color: #ffffff; font-weight: 500;">
            {{ log.actor }}
          </td>
          <td class="vms-td" style="padding: 9px 12px; font-size: 0.7rem; color: rgba(255, 255, 255, 0.7); white-space: nowrap;">
            {{ log.tenantName }}
          </td>
          <td class="vms-td" style="padding: 9px 12px; font-family: var(--vms-font-jetbrains); font-size: 0.7rem; color: var(--vms-neu-accent-orange); font-weight: 700; white-space: nowrap;">
            {{ log.action }}
          </td>
          <td class="vms-td" style="padding: 9px 12px; font-size: 0.75rem; color: #ffffff; max-width: 380px;">
            <div>{{ log.details }}</div>
            <div class="vms-text-mono vms-text-2xs vms-text-dim" style="margin-top: 2px;">IP: {{ log.ipAddress }} // ALVO: {{ log.target }}</div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.vms-log-row:hover { background: rgba(255, 94, 58, 0.05); }
.vms-badge-dim { background: rgba(255, 255, 255, 0.05); color: rgba(255, 255, 255, 0.8); border: 1px solid rgba(255, 255, 255, 0.15); border-radius: 3px; }
.vms-badge-outline { background: transparent; color: #ff5e3a; border: 1px solid #ff5e3a; border-radius: 3px; }
</style>
