<script setup lang="ts">
defineProps<{
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
  <div class="vms-flex-col" style="gap: 1.25rem; width: 100%; box-sizing: border-box;">
    <!-- Header Section -->
    <div class="vms-flex-between" style="align-items: center; border-bottom: 1px solid rgba(255, 255, 255, 0.08); padding-bottom: 0.6rem;">
      <div>
        <h4 class="vms-h4" style="margin: 0; color: #fff;">COORDENADAS GEOGRÁFICAS (GPS)</h4>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">Posicionamento geo-espacial para Mapas & Plantas</span>
      </div>
      <span class="vms-badge" style="background: rgba(0, 240, 255, 0.1); color: var(--vms-neu-accent-cyan); font-size: 9px;">
        GEOLOCALIZAÇÃO
      </span>
    </div>

    <!-- Coordenadas Inputs -->
    <div class="vms-flex-row" style="gap: 1rem; width: 100%; flex-wrap: wrap;">
      <div class="vms-form-group" style="flex: 1; min-width: 200px;">
        <label class="vms-label">Latitude (GPS)</label>
        <input 
          :value="latitude" 
          type="number" 
          step="0.000001" 
          class="vms-auth-input vms-text-mono" 
          placeholder="-23.550520"
          @input="emit('update:latitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" 
        />
      </div>
      <div class="vms-form-group" style="flex: 1; min-width: 200px;">
        <label class="vms-label">Longitude (GPS)</label>
        <input 
          :value="longitude" 
          type="number" 
          step="0.000001" 
          class="vms-auth-input vms-text-mono" 
          placeholder="-46.633308"
          @input="emit('update:longitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" 
        />
      </div>
    </div>

    <!-- HUD Mini Geo Radar Preview -->
    <div style="background: #04060a; border: 1px dashed rgba(0, 240, 255, 0.3); border-radius: 4px; padding: 1rem; display: flex; align-items: center; justify-content: space-between;">
      <div class="vms-flex-row" style="gap: 0.75rem; align-items: center;">
        <div style="width: 12px; height: 12px; border-radius: 50%; background: var(--vms-neu-accent-cyan); box-shadow: 0 0 8px var(--vms-neu-accent-cyan);"></div>
        <div class="vms-flex-col" style="gap: 2px;">
          <span class="vms-text-xs vms-font-semibold" style="color: #fff;">PIN GEO-ESPACIAL VINCULADO</span>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-cyan);">{{ (latitude || 0).toFixed(6) }}, {{ (longitude || 0).toFixed(6) }}</span>
        </div>
      </div>
      <span class="vms-badge" style="background: rgba(0, 240, 255, 0.1); color: var(--vms-neu-accent-cyan); font-size: 9px;">
        MAPA SINCRONIZADO
      </span>
    </div>

    <!-- Setor / Observação -->
    <div class="vms-form-group" style="width: 100%;">
      <label class="vms-label">Setor / Observação no Mapa (Opcional)</label>
      <input 
        :value="locationName" 
        class="vms-auth-input" 
        placeholder="Ex: Guarita Norte ou Setor A"
        @input="emit('update:locationName', ($event.target as HTMLInputElement).value)" 
      />
    </div>
  </div>
</template>
