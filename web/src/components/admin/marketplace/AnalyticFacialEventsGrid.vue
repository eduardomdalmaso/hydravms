<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { AnalyticEventRecord } from '../../../types/marketplace'

const props = defineProps<{ events: AnalyticEventRecord[] }>()
const emit = defineEmits<{ (e: 'inspect', event: AnalyticEventRecord): void }>()

const pageSize = ref(12), currentPage = ref(1)
const viewMode = ref<Record<string, 'crop' | 'full'>>({})
const totalPages = computed(() => Math.max(1, Math.ceil(props.events.length / pageSize.value)))

watch([() => props.events.length, pageSize], () => { if (currentPage.value > totalPages.value) currentPage.value = 1 })

const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return props.events.slice(start, start + pageSize.value)
})

const formatTimestamp = (iso: string) => {
  if (!iso) return '--:--:--'
  try { return new Date(iso).toLocaleString('pt-BR', { timeZone: 'UTC', dateStyle: 'short', timeStyle: 'medium' }) } catch { return iso }
}
</script>

<template>
  <div v-if="events.length === 0" class="vms-text-mono vms-text-sm vms-text-dim" style="text-align: center; padding: 4rem 2rem; border: 1px dashed rgba(255,255,255,0.08); border-radius: 6px; background: rgba(0,0,0,0.2);">
    <div class="vms-flex-col" style="gap: 0.5rem; align-items: center;">
      <span class="vms-text-sm vms-font-semibold" style="color: #ffffff;">// NENHUMA BIOMETRIA FACIAL ENCONTRADA</span>
      <span class="vms-text-2xs vms-text-dim">Nenhum evento corresponde aos filtros de câmera, face conhecida, idade ou período.</span>
    </div>
  </div>

  <div v-else class="vms-flex-col" style="gap: 1rem;">
    <!-- Grid de Cards Biométricos -->
    <div class="vms-facial-events-grid">
      <div v-for="evt in paginatedEvents" :key="evt.id" class="vms-face-card" @click="emit('inspect', evt)">
        <div class="vms-face-crop-preview">
          <!-- Imagem: Cropped ou Full Frame -->
          <img
            v-if="viewMode[evt.id] === 'full' ? (evt.snapshot_url || evt.crop_url) : (evt.crop_url || evt.snapshot_url)"
            :src="viewMode[evt.id] === 'full' ? (evt.snapshot_url || evt.crop_url) : (evt.crop_url || evt.snapshot_url)"
            class="vms-face-thumb-img"
            :alt="viewMode[evt.id] === 'full' ? 'Full Frame' : 'Cropped Face'"
          />
          <div v-else class="vms-face-thumb-placeholder">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--vms-neu-accent-orange)" stroke-width="1.5"><circle cx="12" cy="12" r="9"/><circle cx="12" cy="10" r="3"/><path d="M7 18.5c1.2-2 3-3 5-3s3.8 1 5 3"/></svg>
            <span class="vms-text-mono vms-text-2xs" style="color: #64748b;">SEM IMAGEM</span>
          </div>

          <!-- Bounding Box Reticle SÓLIDO E LARANJA (Apenas no modo FULL FRAME) -->
          <svg v-if="evt.bbox && viewMode[evt.id] === 'full'" viewBox="0 0 100 100" preserveAspectRatio="none" style="position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; z-index: 5;">
            <rect :x="evt.bbox[0]" :y="evt.bbox[1]" :width="evt.bbox[2]" :height="evt.bbox[3]" fill="none" stroke="#ff5e3a" stroke-width="1.8" />
          </svg>

          <!-- Toggle Crop / Full Frame -->
          <div class="vms-flex-row" style="position: absolute; top: 6px; right: 6px; gap: 2px; z-index: 10;" @click.stop>
            <button class="vms-chip" :class="{ active: (viewMode[evt.id] || 'crop') === 'crop' }" style="font-size: 8px; padding: 1px 4px;" title="Ver recorte focado na face" @click="viewMode[evt.id] = 'crop'">CROP</button>
            <button class="vms-chip" :class="{ active: viewMode[evt.id] === 'full' }" style="font-size: 8px; padding: 1px 4px;" title="Ver foto de cena completa (Full Frame)" @click="viewMode[evt.id] = 'full'">FULL</button>
          </div>
        </div>

        <!-- Identificação e Idade / Gênero -->
        <div class="vms-flex-col" style="gap: 3px; padding: 0.2rem 0;">
          <div class="vms-flex-between" style="align-items: center; gap: 4px;">
            <span class="vms-text-xs vms-font-semibold" :style="{ color: (evt.face_status === 'Reconhecido' || evt.is_known) ? '#00ff9d' : '#ff5e3a' }" style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
              {{ (evt.face_status === 'Reconhecido' || evt.is_known) ? (evt.person_name || 'RECONHECIDO') : 'DESCONHECIDO' }}
            </span>
            <span class="vms-badge" :style="{ fontSize: '8px', padding: '1px 5px', fontWeight: 'bold', background: (evt.face_status === 'Reconhecido' || evt.is_known ? 'rgba(0,255,157,0.15)' : 'rgba(255,94,58,0.15)'), color: (evt.face_status === 'Reconhecido' || evt.is_known ? '#00ff9d' : '#ff5e3a'), border: `1px solid ${evt.face_status === 'Reconhecido' || evt.is_known ? '#00ff9d55' : '#ff5e3a55'}` }">
              {{ (evt.face_status === 'Reconhecido' || evt.is_known) ? 'RECONHECIDO' : 'DESCONHECIDO' }}
            </span>
          </div>
          <div class="vms-flex-between" style="align-items: center;">
            <span class="vms-text-mono vms-text-2xs" style="color: #cbd5e1;">
              {{ (evt.gender ? evt.gender.toUpperCase() : 'GÊNERO N/D') }} // {{ evt.age ? `~${evt.age} ANOS` : 'IDADE N/D' }}
            </span>
          </div>
          <span class="vms-text-2xs vms-text-dim" style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ evt.camera_name }}
          </span>
        </div>

        <!-- Footer com Data/Hora e Inspecionar -->
        <div class="vms-flex-between" style="border-top: 1px solid var(--vms-border); padding-top: 0.45rem; margin-top: auto; align-items: center;">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ formatTimestamp(evt.timestamp) }}</span>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 2px 6px;">INSPECIONAR</button>
        </div>
      </div>
    </div>

    <!-- Paginação -->
    <div class="vms-flex-between" style="align-items: center; border-top: 1px solid rgba(255,255,255,0.08); padding-top: 0.75rem; flex-wrap: wrap; gap: 0.75rem;">
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">ITENS / PÁGINA:</span>
        <button v-for="size in [12, 24, 48]" :key="size" class="vms-btn vms-btn-sm" :class="pageSize === size ? 'vms-btn-primary' : 'vms-btn-secondary'" style="font-size: 9.5px; padding: 3px 7px;" @click="pageSize = size">{{ size }}</button>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 4px 8px;" :disabled="currentPage <= 1" @click="currentPage--">ANTERIOR</button>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">PÁGINA {{ currentPage }} DE {{ totalPages }} // TOTAL: {{ events.length }}</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 9.5px; padding: 4px 8px;" :disabled="currentPage >= totalPages" @click="currentPage++">PRÓXIMO</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-facial-events-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 0.85rem; }
.vms-face-card { background: var(--vms-surface); border: 1px solid var(--vms-border); border-radius: 6px; padding: 0.65rem; display: flex; flex-direction: column; gap: 0.4rem; cursor: pointer; transition: transform 0.15s, border-color 0.15s; }
.vms-face-card:hover { transform: translateY(-2px); border-color: rgba(255, 94, 58, 0.4); }
.vms-face-crop-preview { width: 100%; aspect-ratio: 4 / 3; background: #060910; border-radius: 4px; overflow: hidden; position: relative; display: flex; align-items: center; justify-content: center; }
.vms-face-thumb-img { width: 100%; height: 100%; object-fit: cover; display: block; }
.vms-face-thumb-placeholder { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; width: 100%; height: 100%; }
</style>
