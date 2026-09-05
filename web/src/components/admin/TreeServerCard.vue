<script setup lang="ts">
import { ref } from 'vue'

export interface CameraNode {
  id: string
  name: string
  ip: string
  codec: string
  res: string
  fps: number
  bitrate: string
  status: 'online' | 'offline' | 'recording'
}

export interface SwitchNode {
  id: string
  name: string
  ip: string
  cameras: CameraNode[]
}

export interface ServerNode {
  id: string
  name: string
  ip: string
  location: string
  switches: SwitchNode[]
}

const props = defineProps<{ server: ServerNode }>()
const emit = defineEmits<{ (e: 'test-stream', camId: string): void }>()
const isExpanded = ref(true)
const switchOpen = ref<Record<string, boolean>>({ sw_01: true, sw_02: true, sw_03: true })

const toggleSwitch = (swId: string) => {
  switchOpen.value[swId] = !switchOpen.value[swId]
}

const totalCams = () => props.server.switches.reduce((acc, sw) => acc + sw.cameras.length, 0)
</script>

<template>
  <div class="vms-tree-server-card">
    <div class="vms-tree-server-header" @click="isExpanded = !isExpanded">
      <div class="vms-flex-row" style="gap: 0.75rem;">
        <span class="vms-text-dim" style="font-size: 11px;">{{ isExpanded ? '▼' : '►' }}</span>
        <span class="vms-font-bold" style="color: #fff;">{{ server.name }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// {{ server.ip }} ({{ server.location }})</span>
      </div>
      <div class="vms-flex-row" style="gap: 1rem;">
        <span class="vms-badge vms-badge-info">{{ totalCams() }} FLUXOS</span>
        <span class="vms-status-led online" title="Edge Server Online"></span>
      </div>
    </div>

    <div v-if="isExpanded" style="padding: 0.5rem 0.75rem;">
      <div v-for="sw in server.switches" :key="sw.id" class="vms-tree-switch-group">
        <div class="vms-tree-switch-header" @click="toggleSwitch(sw.id)" style="cursor: pointer;">
          <div class="vms-flex-row" style="gap: 0.5rem;">
            <span class="vms-text-dim" style="font-size: 10px;">{{ switchOpen[sw.id] ? '▼' : '►' }}</span>
            <span class="vms-text-mono vms-font-semibold" style="color: var(--vms-neu-accent-orange); font-size: 11px;">{{ sw.name }}</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ sw.ip }}</span>
          </div>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ sw.cameras.length }} PORTAS POE</span>
        </div>

        <div v-if="switchOpen[sw.id]" class="vms-flex-col" style="gap: 0.35rem; margin-top: 0.35rem;">
          <div v-for="cam in sw.cameras" :key="cam.id" class="vms-tree-camera-row">
            <div class="vms-flex-row" style="gap: 0.75rem;">
              <span class="vms-status-led" :class="cam.status" :title="cam.status === 'recording' ? 'Gravacao Ativa' : (cam.status === 'offline' ? 'Offline' : 'Online')"></span>
              <span class="vms-font-medium" style="color: #fff; font-size: 12px;">{{ cam.name }}</span>
              <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ cam.ip }}</span>
            </div>
            <div class="vms-flex-row" style="gap: 0.65rem;">
              <span class="vms-badge vms-badge-info" style="font-size: 9px;">{{ cam.codec }}</span>
              <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ cam.res }} @ {{ cam.fps }}FPS</span>
              <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-blue);">{{ cam.bitrate }}</span>
              <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 10px; padding: 2px 6px;" @click.stop="emit('test-stream', cam.id)">[TESTAR]</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
