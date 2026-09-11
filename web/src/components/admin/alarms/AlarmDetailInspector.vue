<script setup lang="ts">
import type { AlarmItem } from '../../../types/alarmTree'

defineProps<{ alarm: AlarmItem }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'test', id: string): void
  (e: 'toggleStatus', id: string): void
  (e: 'delete', id: string): void
}>()
</script>

<template>
  <div class="vms-desktop-inspector">
    <!-- Header -->
    <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.75rem;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">INSPEÇÃO DE SENSOR // HUD</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ID: {{ alarm.id }}</span>
      </div>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <!-- Inspector Icon Header -->
    <div class="vms-flex-col" style="align-items: center; gap: 0.5rem; padding: 1rem 0; border-bottom: 1px solid var(--vms-border);">
      <div style="width: 56px; height: 56px; border-radius: 12px; background: rgba(255, 94, 58, 0.15); border: 1px solid rgba(255, 94, 58, 0.4); display: flex; align-items: center; justify-content: center; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" /><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
        </svg>
      </div>
      <span class="vms-font-bold" style="color: #fff; font-size: 14px; text-align: center;">{{ alarm.name }}</span>
      <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.4); font-size: 10px;">
        [{{ alarm.type }}] SENSOR
      </span>
    </div>

    <!-- Technical Details -->
    <div class="vms-flex-col" style="gap: 0.5rem; background: rgba(255, 255, 255, 0.02); padding: 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">ZONA / LOCAL</span><span class="vms-text-mono vms-text-xs" style="color: #fff;">{{ alarm.zone }}</span></div>
      <div class="vms-flex-between">
        <span class="vms-text-dim vms-text-2xs">STATUS ATUAL</span>
        <div class="vms-flex-row" style="gap: 6px; align-items: center;">
          <span class="vms-status-led" :class="alarm.status"></span>
          <span class="vms-text-mono vms-text-xs" style="color: #fff;">
            {{ alarm.status === 'online' ? 'ARMADO' : (alarm.status === 'alert' ? 'DESARMADO' : 'OFFLINE') }}
          </span>
        </div>
      </div>
      <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">SENSIBILIDADE</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-neu-accent-orange);">{{ alarm.sensitivity }}%</span></div>
      <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">CÂMERA VINCULADA</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-regular);">{{ alarm.linkedCameraName || 'NENHUMA' }}</span></div>
      <div class="vms-flex-between"><span class="vms-text-dim vms-text-2xs">ÚLTIMO DISPARO</span><span class="vms-text-mono vms-text-xs" style="color: var(--vms-text-main);">{{ alarm.lastTrigger || 'N/A' }}</span></div>
    </div>

    <!-- Actions -->
    <div class="vms-flex-col" style="gap: 0.5rem; margin-top: auto;">
      <button class="vms-btn vms-btn-primary" style="font-size: 11px;" @click="emit('test', alarm.id)">[TESTAR DISPARO DE ALARME]</button>
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="color: #ff5e3a; font-size: 11px; display: inline-flex; align-items: center; justify-content: center; gap: 6px;" @click="emit('delete', alarm.id)">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/></svg>
        <span>[REMOVER SENSOR]</span>
      </button>
    </div>
  </div>
</template>
