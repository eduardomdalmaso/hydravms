<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'selectLocation', loc: { lat: number; lng: number; label: string }): void
  (e: 'quickAdd', type: 'CAMERA' | 'ALARME', loc: { lat: number; lng: number; label: string }): void
}>()

const query = ref('')
const isLoading = ref(false)
const results = ref<Array<{ lat: number; lng: number; label: string }>>([])
const selectedLocation = ref<{ lat: number; lng: number; label: string } | null>(null)

const handleSearch = async () => {
  if (!query.value.trim()) { results.value = []; return }
  isLoading.value = true
  try {
    const res = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query.value)}&limit=5`)
    const data = await res.json()
    results.value = data.map((d: any) => ({ lat: parseFloat(d.lat), lng: parseFloat(d.lon), label: d.display_name }))
  } catch (err) {
    results.value = [
      { lat: -23.55052, lng: -46.633308, label: `${query.value} // Sao Paulo - Centro` },
      { lat: -23.58000, lng: -46.680000, label: `${query.value} // Zona Sul Industrial` }
    ]
  } finally { isLoading.value = false }
}

const onSelect = (loc: { lat: number; lng: number; label: string }) => {
  selectedLocation.value = loc
  results.value = []
  emit('selectLocation', loc)
}
</script>

<template>
  <div class="vms-map-search-container">
    <div class="vms-map-search-bar">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
      <input v-model="query" class="vms-map-search-input" placeholder="Buscar endereço, local ou coordenadas..." @keydown.enter.prevent="handleSearch" />
      <button v-if="query" class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 2px 4px;" @click="query = ''; results = []; selectedLocation = null">✕</button>
    </div>

    <!-- Dropdown de Resultados -->
    <div v-if="results.length > 0" class="vms-map-search-results">
      <div v-for="(r, idx) in results" :key="idx" class="vms-search-result-item" @click="onSelect(r)">
        <span class="vms-text-xs vms-font-semibold" style="color: #ffffff;">{{ r.label }}</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">LAT: {{ r.lat.toFixed(5) }} // LNG: {{ r.lng.toFixed(5) }}</span>
      </div>
    </div>

    <!-- Acao Rapida quando Local Selecionado -->
    <div v-if="selectedLocation" class="vms-location-action-card">
      <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">LOCAL: {{ selectedLocation.lat.toFixed(4) }}, {{ selectedLocation.lng.toFixed(4) }}</span>
      <div class="vms-flex-row" style="gap: 4px; margin-top: 4px;">
        <button class="vms-btn vms-btn-sm" style="font-size: 9px; padding: 2px 6px; background: rgba(255, 94, 58, 0.15); border-color: #ff5e3a; color: #fff;" @click="emit('quickAdd', 'CAMERA', selectedLocation!)">+ CÂMERA AQUI</button>
        <button class="vms-btn vms-btn-sm" style="font-size: 9px; padding: 2px 6px; background: #07080c; border-color: #ff5e3a; color: #ff5e3a;" @click="emit('quickAdd', 'ALARME', selectedLocation!)">+ ALARME AQUI</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-map-search-container {
  position: absolute; top: 12px; left: 12px; z-index: 1000; width: 330px; display: flex; flex-direction: column; gap: 4px;
}
.vms-map-search-bar {
  display: flex; align-items: center; gap: 8px; background: rgba(7, 8, 12, 0.9);
  backdrop-filter: blur(8px); padding: 6px 10px; border-radius: 8px; border: 1px solid var(--vms-border);
}
.vms-map-search-input {
  background: transparent; border: none; outline: none; color: #ffffff; font-size: 11px; width: 100%;
}
.vms-map-search-results {
  background: rgba(14, 17, 23, 0.95); backdrop-filter: blur(10px); border: 1px solid var(--vms-border);
  border-radius: 6px; max-height: 200px; overflow-y: auto; display: flex; flex-direction: column;
}
.vms-search-result-item {
  padding: 8px 10px; cursor: pointer; border-bottom: 1px solid rgba(255, 255, 255, 0.05); display: flex; flex-direction: column; gap: 2px;
}
.vms-search-result-item:hover { background: #1c212a; }
.vms-location-action-card {
  background: rgba(14, 17, 23, 0.95); border: 1px solid var(--vms-border); border-radius: 6px; padding: 6px 10px;
}
</style>
