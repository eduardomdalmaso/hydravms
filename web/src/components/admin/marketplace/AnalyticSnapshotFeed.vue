<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'
import { useEventBus, initEventSocket } from '../../../services/eventSocket'

const props = defineProps<{
  cameraId: string
  cameraName: string
  analyticName: string
}>()

interface LocalSnapshot {
  id: string
  time: string
  label: string
  conf: number
  zone: string
  url: string
}

const snapshots = ref<LocalSnapshot[]>([])
const { subscribe } = useEventBus()
let unsubscribe: (() => void) | null = null

onMounted(() => {
  initEventSocket()
  const baseImg = getCameraSnapshotUrl(props.cameraId, false)
  
  unsubscribe = subscribe((evt) => {
    const isTargetCam = evt.subject === props.cameraId || evt.data?.camera_id === props.cameraId
    if (isTargetCam) {
      const snapUrl = evt.data?.snapshot_url || baseImg
      const label = (evt.data?.object_label || evt.data?.class || 'DETECÇÃO').toUpperCase()
      const conf = Math.round((evt.data?.confidence || 0.95) * 100)
      const zone = evt.data?.zone_name || evt.data?.rule_name || 'ZONA ATIVA'
      const time = new Date(evt.time || Date.now()).toLocaleTimeString()

      snapshots.value.unshift({
        id: evt.id || `snap-${Date.now()}`,
        time,
        label,
        conf,
        zone,
        url: snapUrl
      })

      if (snapshots.value.length > 30) snapshots.value.pop()
    }
  })
})

onUnmounted(() => {
  if (unsubscribe) unsubscribe()
})
</script>

<template>
  <div class="vms-card vms-flex-col" style="flex: 1; height: 100%; max-height: calc(100vh - 210px); display: flex; flex-direction: column; background: #0b0e14; border: 1px solid var(--vms-border); border-radius: 8px; overflow: hidden;">
    <div class="vms-flex-between" style="padding: 0.75rem 1rem; border-bottom: 1px solid var(--vms-border); background: rgba(255,255,255,0.02);">
      <span class="vms-text-mono vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">SNAPSHOTS // DETECÇÕES</span>
      <span class="vms-badge vms-badge-green" style="font-size: 9px;">{{ snapshots.length }} EVENTOS</span>
    </div>

    <div class="vms-flex-col" style="padding: 0.75rem; gap: 0.65rem; overflow-y: auto; flex: 1;">
      <div v-if="snapshots.length === 0" class="vms-text-mono vms-text-xs vms-text-dim" style="text-align: center; padding: 2rem 0;">
        // AGUARDANDO DETECÇÕES REAIS...
      </div>

      <div
        v-for="s in snapshots"
        :key="s.id"
        class="vms-card"
        style="padding: 0.5rem; background: #121824; border: 1px solid rgba(255,255,255,0.06); border-radius: 6px; display: flex; gap: 0.65rem; align-items: center;"
      >
        <div style="width: 72px; height: 50px; border-radius: 4px; overflow: hidden; background: #000; flex-shrink: 0; border: 1px solid rgba(0, 240, 255, 0.3);">
          <img :src="s.url" style="width: 100%; height: 100%; object-fit: cover;" alt="Snapshot" />
        </div>

        <div class="vms-flex-col" style="gap: 2px; flex: 1; min-width: 0;">
          <div class="vms-flex-between" style="align-items: center;">
            <span class="vms-badge vms-badge-orange" style="font-size: 8.5px; font-weight: bold;">[{{ s.label }}] {{ s.conf }}%</span>
            <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ s.time }}</span>
          </div>
          <span class="vms-text-mono vms-text-2xs" style="color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ s.zone }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
