<script setup lang="ts">
const props = defineProps<{
  latitude: number
  longitude: number
  locationName: string
}>()
const emit = defineEmits<{
  (e: 'update:latitude', val: number): void
  (e: 'update:longitude', val: number): void
  (e: 'update:locationName', val: string): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="min-height: 360px; justify-content: space-between; gap: 0.85rem;">
    <div class="vms-flex-col" style="gap: 0.85rem;">
      <div class="vms-form-group">
      <label class="vms-label">Identificação / Setor no Mapa</label>
      <input 
        :value="locationName" 
        class="vms-auth-input" 
        placeholder="Ex: Entrada Principal - Guarita Norte (Setor A)"
        @input="emit('update:locationName', ($event.target as HTMLInputElement).value)" 
      />
    </div>

    <div class="vms-flex-row" style="gap: 0.6rem;">
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Latitude (GPS)</label>
        <input 
          :value="latitude" 
          type="number" 
          step="0.000001" 
          class="vms-auth-input vms-text-mono" 
          @input="emit('update:latitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" 
        />
      </div>
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Longitude (GPS)</label>
        <input 
          :value="longitude" 
          type="number" 
          step="0.000001" 
          class="vms-auth-input vms-text-mono" 
          @input="emit('update:longitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" 
        />
      </div>
    </div>

    <!-- HUD Mini Geo Radar Preview -->
    <div style="background: #04060a; border: 1px dashed rgba(0, 240, 255, 0.3); border-radius: 4px; padding: 0.75rem; display: flex; align-items: center; justify-content: space-between;">
      <div class="vms-flex-row" style="gap: 0.6rem; align-items: center;">
        <div style="width: 12px; height: 12px; border-radius: 50%; background: var(--vms-neu-accent-cyan); box-shadow: 0 0 8px var(--vms-neu-accent-cyan);"></div>
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">PIN GEO-ESPACIAL VINCULADO</span>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan);">{{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</span>
        </div>
      </div>
      <span class="vms-badge" style="background: rgba(0, 240, 255, 0.1); color: var(--vms-neu-accent-cyan); font-size: 9px;">
        MAPA VMS SINCRONIZADO
      </span>
    </div>
    </div>
  </div>
</template>
