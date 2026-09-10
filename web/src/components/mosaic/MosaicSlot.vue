<script setup lang="ts">
import type { CameraStreamInfo } from '../../types/mosaic'

const props = defineProps<{
  camera?: CameraStreamInfo
  isActive?: boolean
  isHero?: boolean
}>()

const emit = defineEmits<{ (e: 'select', cam: CameraStreamInfo): void }>()
</script>

<template>
  <div
    class="vms-slot"
    :class="{ active: isActive, 'vms-slot-hero': isHero }"
    @click="camera && emit('select', camera)"
  >
    <template v-if="camera">
      <div class="vms-slot-video" style="background: #000; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; overflow: hidden; position: relative;">
        <img
          :src="`http://localhost:8080/api/v1/streams/${camera.id}/mjpeg`"
          alt="Camera Stream"
          style="width: 100%; height: 100%; object-fit: cover; display: block;"
        />
      </div>
      <div class="vms-slot-hud">
        <div class="vms-flex-row" style="gap: 0.5rem;">
          <span class="vms-badge" style="background: rgba(15, 23, 42, 0.85); color: #fff; font-size: 11px;">
            {{ camera.name }}
          </span>
          <span v-if="camera.has_ptz" class="vms-badge vms-badge-info" style="font-size: 10px;">[PTZ]</span>
        </div>
        <div class="vms-flex-between">
          <span class="vms-text-mono vms-text-xs vms-text-dim" style="background: rgba(0,0,0,0.6); padding: 2px 4px; border-radius: 2px;">
            [REC AUTO]
          </span>
          <span class="vms-badge vms-badge-recording" style="font-size: 10px;">[REC]</span>
        </div>
      </div>
    </template>
    <template v-else>
      <span class="vms-text-xs vms-text-dim">[SLOT VAZIO]</span>
    </template>
  </div>
</template>
