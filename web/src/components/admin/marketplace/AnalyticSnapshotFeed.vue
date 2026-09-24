<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCameraSnapshotUrl } from '../../../utils/streamUrls'

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

onMounted(() => {
  const baseImg = getCameraSnapshotUrl(props.cameraId, false)
  const now = new Date()
  snapshots.value = [
    { id: '1', time: new Date(now.getTime() - 12000).toLocaleTimeString(), label: 'PESSOA', conf: 96, zone: 'ZONA 01 // ENTRADA', url: baseImg },
    { id: '2', time: new Date(now.getTime() - 45000).toLocaleTimeString(), label: 'VEÍCULO', conf: 91, zone: 'ZONA 02 // ESTACIONAMENTO', url: baseImg },
    { id: '3', time: new Date(now.getTime() - 110000).toLocaleTimeString(), label: 'PESSOA', conf: 88, zone: 'ZONA 01 // ENTRADA', url: baseImg }
  ]
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
        // AGUARDANDO DETECÇÕES...
      </div>

      <div
        v-for="s in snapshots"
        :key="s.id"
        class="vms-card"
        style="padding: 0.5rem; background: #121824; border: 1px solid rgba(255,255,255,0.06); border-radius: 6px; display: flex; gap: 0.65rem; align-items: center;"
      >
        <!-- Thumbnail with crop -->
        <div style="width: 72px; height: 50px; border-radius: 4px; overflow: hidden; background: #000; flex-shrink: 0; border: 1px solid rgba(0, 240, 255, 0.3);">
          <img :src="s.url" style="width: 100%; height: 100%; object-fit: cover;" alt="Snapshot" />
        </div>

        <!-- Meta -->
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
