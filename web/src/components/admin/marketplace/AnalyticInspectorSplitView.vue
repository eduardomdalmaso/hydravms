<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import type { AnalyticInstance } from '../../../types/marketplace'
import AnalyticLiveStreamViewer from './AnalyticLiveStreamViewer.vue'
import AnalyticParametersPanel from './AnalyticParametersPanel.vue'
import AnalyticSnapshotFeed from './AnalyticSnapshotFeed.vue'

const props = defineProps<{ instance: AnalyticInstance }>()
const emit = defineEmits<{
  (e: 'saved', inst: AnalyticInstance): void
  (e: 'delete', id: string): void
  (e: 'close'): void
}>()

const localInst = ref<AnalyticInstance>({ ...props.instance })
const formActive = ref(props.instance.is_active)
const activeSeconds = ref(1254) // Initial simulated active uptime

let timerInterval: number | null = null

const formattedUptime = computed(() => {
  const sec = activeSeconds.value
  const years = Math.floor(sec / (365 * 24 * 3600))
  const days = Math.floor((sec % (365 * 24 * 3600)) / (24 * 3600))
  const hours = Math.floor((sec % (24 * 3600)) / 3600)
  const minutes = Math.floor((sec % 3600) / 60)
  const seconds = sec % 60

  const yStr = String(years).padStart(2, '0') + 'a'
  const dStr = String(days).padStart(3, '0') + 'd'
  const hStr = String(hours).padStart(2, '0') + 'h'
  const mStr = String(minutes).padStart(2, '0') + 'm'
  const sStr = String(seconds).padStart(2, '0') + 's'

  return `${yStr} ${dStr} ${hStr} ${mStr} ${sStr}`
})

const startTimer = () => {
  if (timerInterval) clearInterval(timerInterval)
  timerInterval = window.setInterval(() => {
    if (formActive.value) activeSeconds.value++
  }, 1000)
}

watch(() => props.instance, (newVal) => {
  localInst.value = { ...newVal }
  formActive.value = newVal.is_active
})

const toggleActiveState = () => {
  formActive.value = !formActive.value
  localInst.value.is_active = formActive.value
  emit('saved', { ...localInst.value })
}

const handleSave = () => {
  emit('saved', { ...localInst.value, is_active: formActive.value })
}

onMounted(() => { startTimer() })
onUnmounted(() => { if (timerInterval) clearInterval(timerInterval) })
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem; width: 100%;">
    <!-- Top Header Bar with Active Uptime Timer -->
    <div class="vms-card vms-flex-between" style="padding: 0.75rem 1.25rem; align-items: center; background: #0c0f17; border: 1px solid var(--vms-border); border-radius: 8px;">
      <div class="vms-flex-row" style="gap: 1.25rem; align-items: center; flex-wrap: wrap;">
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px 8px; font-size: 11px;" @click="emit('close')">
          ← VOLTAR
        </button>
        <div class="vms-flex-col" style="gap: 2px;">
          <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
            <span class="vms-text-sm vms-font-semibold" style="color: var(--vms-neu-accent-orange);">
              {{ localInst.name.toUpperCase() }}
            </span>
            <!-- Live Uptime Timer -->
            <div class="vms-flex-row" style="gap: 6px; align-items: center; background: rgba(0, 240, 255, 0.08); border: 1px solid rgba(0, 240, 255, 0.25); padding: 2px 8px; border-radius: 4px;">
              <span style="width: 6px; height: 6px; border-radius: 50%;" :style="{ background: formActive ? '#00ff9d' : '#ff5e3a', boxShadow: formActive ? '0 0 6px #00ff9d' : 'none' }" />
              <span class="vms-text-mono vms-text-2xs" style="color: #00f0ff; font-weight: bold; letter-spacing: 0.5px;">
                TEMPO ATIVO: {{ formattedUptime }}
              </span>
            </div>
          </div>
          <span class="vms-text-mono vms-text-2xs vms-text-dim">CÂMERA: {{ localInst.camera_name }} // PLUGIN: {{ localInst.plugin_name }}</span>
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-sm" :class="formActive ? 'vms-btn-secondary' : 'vms-btn-primary'" style="font-weight: bold; padding: 5px 12px;" @click="toggleActiveState">
          {{ formActive ? 'PAUSAR' : 'RETOMAR' }}
        </button>
        <button class="vms-btn vms-btn-danger vms-btn-sm" style="padding: 5px 10px;" @click="emit('delete', localInst.id)">
          EXCLUIR
        </button>
        <button :disabled="formActive" class="vms-btn vms-btn-primary vms-btn-sm" :style="{ opacity: formActive ? 0.35 : 1, cursor: formActive ? 'not-allowed' : 'pointer' }" style="font-weight: bold; padding: 5px 14px;" @click="handleSave">
          SALVAR
        </button>
      </div>
    </div>

    <!-- Main Split Columns -->
    <div class="vms-flex-row" style="gap: 1rem; align-items: flex-start;">
      <!-- Left Column: Realtime Stream with BBox Overlay & Parameters Panel -->
      <div class="vms-flex-col" style="flex: 1.35; gap: 1rem; min-width: 0;">
        <AnalyticLiveStreamViewer
          :camera-id="localInst.camera_id"
          :camera-name="localInst.camera_name"
          :is-active="formActive"
          :zones="localInst.zones"
          :fps="localInst.fps_rate"
        />
        <AnalyticParametersPanel
          :instance="localInst"
          :disabled="formActive"
          @update="(updated) => localInst = updated"
        />
      </div>

      <!-- Right Column: Recent Snapshots Feed -->
      <div style="flex: 0.85; min-width: 320px;">
        <AnalyticSnapshotFeed
          :camera-id="localInst.camera_id"
          :camera-name="localInst.camera_name"
          :analytic-name="localInst.name"
        />
      </div>
    </div>
  </div>
</template>
