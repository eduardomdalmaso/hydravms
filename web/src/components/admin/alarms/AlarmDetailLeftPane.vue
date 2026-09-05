<script setup lang="ts">
import type { AlarmItem } from '../../../types/alarmTree'

defineProps<{ alarm: AlarmItem }>()
</script>

<template>
  <div class="vms-split-pane">
    <!-- Header -->
    <div class="vms-split-header">
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <div style="width: 36px; height: 36px; border-radius: 8px; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center;">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" /><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
          </svg>
        </div>
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-font-bold" style="color: #fff; font-size: 13px;">{{ alarm.id.toUpperCase() }} // {{ alarm.name }}</span>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">ZONA // {{ alarm.zone.toUpperCase() }}</span>
        </div>
      </div>
    </div>

    <!-- Telemetry & Sensor Parameters Grid (Preenchendo todo o espaço da div) -->
    <div class="vms-telemetry-grid" style="flex: 1; align-content: start;">
      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">TIPO DE SENSOR</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">[{{ alarm.type }}] SENSOR</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ZONA / LOCAL</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">{{ alarm.zone }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">STATUS ATUAL</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">
          {{ alarm.status === 'online' ? 'ARMADO (SUPERVISIONADO)' : (alarm.status === 'alert' ? 'DISPARADO (ALERTA)' : 'DESARMADO') }}
        </span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">SENSIBILIDADE</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">{{ alarm.sensitivity }}%</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">CAMERA VINCULADA</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-text-regular);">{{ alarm.linkedCameraName || 'NENHUMA' }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ULTIMO DISPARO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-text-main);">{{ alarm.lastTrigger || 'SEM DISPAROS RECENTES' }}</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">MODO RETARDO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-text-main);">IMEDIATO (0s)</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">CANAL I/O RELAY</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">RELAY_01 (NORMAL ABERTO)</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">BARRAMENTO / PROTOCOLO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">MODBUS TCP // IP ONVIF</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ANTI-SABOTAGEM (TAMPER)</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #00ff9d;">HABILITADO (SUPERVISIONADO)</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">NIVEL DE ALIMENTACAO</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">13.8V DC // NORMAL</span>
      </div>

      <div class="vms-telemetry-card">
        <span class="vms-text-dim vms-text-2xs">ESTADO DO LACO (EOL)</span>
        <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: #fff;">FECHADO (SUPERVISIONADO)</span>
      </div>
    </div>
  </div>
</template>
