<script setup lang="ts">
import { ref } from "vue"
import type { WorkspaceSlot, CameraStreamInfo } from "../../types/mosaic"
import SlotLinkedAlarm from "./SlotLinkedAlarm.vue"

const props = defineProps<{ slot: WorkspaceSlot; isActive?: boolean; isHero?: boolean }>()
const emit = defineEmits<{ (e: "selectCamera", cam: CameraStreamInfo): void; (e: "clear", slotIndex: number): void }>()

const isGearOpen = ref(false)
const decoderMode = ref<"H264" | "MSE">("H264")

const handleSlotClick = () => {
  if (props.slot.type === "camera" && props.slot.data) {
    emit("selectCamera", props.slot.data as CameraStreamInfo)
  }
}
</script>

<template>
  <div class="vms-slot" :class="{ active: isActive, 'vms-slot-hero': isHero }" draggable="true" @click="handleSlotClick">
    <button v-if="slot.data" class="vms-slot-close-btn" title="Fechar Slot" @click.stop="emit('clear', slot.slot_index)">
      <svg width="9" height="9" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
    </button>

    <!-- CAMERA TYPE -->
    <template v-if="slot.type === 'camera' && slot.data">
      <div class="vms-slot-video" style="background: #05070a; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; position: relative;">
        <div v-if="(slot.data as CameraStreamInfo).status === 'offline'" class="vms-offline-sphere-container" title="Camera Offline - Reconectando...">
          <div class="vms-ubuntu-spinner"></div>
        </div>
      </div>

      <!-- HUD on Hover with Gear, Linked Alarm Bell & Decoder Toggle -->
      <div class="vms-slot-hud">
        <div class="vms-flex-row" style="gap: 0.35rem; margin-top: 1.1rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.9); color: #fff; font-size: 11px;">
            {{ (slot.data as CameraStreamInfo).name }}
          </span>

          <div class="vms-gear-wrapper" style="position: relative;" @mouseenter="isGearOpen = true" @mouseleave="isGearOpen = false">
            <button class="vms-gear-btn" title="Informacoes tecnicas" @click.stop="isGearOpen = !isGearOpen">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
              </svg>
            </button>
            <div v-if="isGearOpen" class="vms-gear-popover">
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">CODEC: H.265</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">RESOLUCAO: {{ (slot.data as CameraStreamInfo).resolution }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">MEDIA FPS: {{ (slot.data as CameraStreamInfo).fps }} FPS</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">DECODER: {{ decoderMode }}</span>
              <span class="vms-text-mono vms-text-2xs" style="color: #94a3b8;">PROTOCOLO: {{ (slot.data as CameraStreamInfo).protocol || 'RTSP' }}</span>
            </div>
          </div>

          <!-- Linked Alarm Bell Popover Component -->
          <SlotLinkedAlarm :cameraId="(slot.data as CameraStreamInfo).id" />

          <button class="vms-decoder-pill-btn" title="Alternar decodificador" @click.stop="decoderMode = decoderMode === 'H264' ? 'MSE' : 'H264'">
            [{{ decoderMode }}]
          </button>
        </div>
      </div>
    </template>

    <!-- MAP TYPE -->
    <template v-else-if="slot.type === 'map' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #111419;">
        <span class="vms-badge vms-channel-pill-telegram">[MAPA INTERATIVO]</span>
        <span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span>
      </div>
    </template>

    <!-- CAROUSEL TYPE -->
    <template v-else-if="slot.type === 'carousel' && slot.data">
      <div class="vms-flex-col" style="align-items: center; justify-content: center; gap: 0.4rem; width: 100%; height: 100%; background: #0c0f14;">
        <span class="vms-badge vms-badge-warning">[RONDA ATIVA // {{ (slot.data as any).interval_seconds }}S]</span>
        <span class="vms-text-sm vms-font-semibold" style="color: #fff;">{{ (slot.data as any).name }}</span>
      </div>
    </template>

    <!-- EMPTY SLOT -->
    <template v-else>
      <div style="display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; user-select: none;">
        <span class="vms-text-mono vms-text-xs vms-text-dim">[ + ]</span>
      </div>
    </template>
  </div>
</template>
