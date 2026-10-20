<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = withDefaults(defineProps<{
  latitude: number
  longitude: number
  locationName: string
  timezone?: string
}>(), {
  timezone: 'America/Sao_Paulo'
})

const emit = defineEmits<{
  (e: 'update:latitude', val: number): void
  (e: 'update:longitude', val: number): void
  (e: 'update:locationName', val: string): void
  (e: 'update:timezone', val: string): void
}>()

const selectedRegion = ref('south_america')

const TIMEZONES_BY_REGION: Record<string, { label: string; timezones: { id: string; label: string; offset: string }[] }> = {
  south_america: {
    label: 'América do Sul (Brasil)',
    timezones: [
      { id: 'America/Sao_Paulo', label: 'Brasília / SP / RJ / Sul / NE / CO', offset: 'UTC-03:00' },
      { id: 'America/Manaus', label: 'Amazonas (Manaus)', offset: 'UTC-04:00' },
      { id: 'America/Cuiaba', label: 'Mato Grosso / MS (Cuiabá)', offset: 'UTC-04:00' },
      { id: 'America/Porto_Velho', label: 'Rondônia (Porto Velho)', offset: 'UTC-04:00' },
      { id: 'America/Boa_Vista', label: 'Roraima (Boa Vista)', offset: 'UTC-04:00' },
      { id: 'America/Rio_Branco', label: 'Acre (Rio Branco)', offset: 'UTC-05:00' },
      { id: 'America/Noronha', label: 'Fernando de Noronha', offset: 'UTC-02:00' },
      { id: 'America/Buenos_Aires', label: 'Argentina (Buenos Aires)', offset: 'UTC-03:00' },
      { id: 'America/Santiago', label: 'Chile (Santiago)', offset: 'UTC-03:00' },
      { id: 'America/Bogota', label: 'Colômbia (Bogotá)', offset: 'UTC-05:00' },
    ]
  },
  north_america: {
    label: 'América do Norte',
    timezones: [
      { id: 'America/New_York', label: 'Leste (NY / Miami) [EST/EDT]', offset: 'UTC-05:00' },
      { id: 'America/Chicago', label: 'Central (Chicago) [CST/CDT]', offset: 'UTC-06:00' },
      { id: 'America/Denver', label: 'Montanhas (Denver) [MST/MDT]', offset: 'UTC-07:00' },
      { id: 'America/Los_Angeles', label: 'Pacífico (LA / SF) [PST/PDT]', offset: 'UTC-08:00' },
      { id: 'America/Mexico_City', label: 'México (CDMX)', offset: 'UTC-06:00' },
    ]
  },
  europe: {
    label: 'Europa',
    timezones: [
      { id: 'Europe/London', label: 'Reino Unido / Irlanda [WET/GMT]', offset: 'UTC+00:00' },
      { id: 'Europe/Paris', label: 'França / Alemanha / Itália [CET]', offset: 'UTC+01:00' },
      { id: 'Europe/Madrid', label: 'Espanha (Madri) [CET]', offset: 'UTC+01:00' },
      { id: 'Europe/Lisbon', label: 'Portugal (Lisboa) [WET]', offset: 'UTC+00:00' },
    ]
  },
  utc: {
    label: 'Universal (UTC)',
    timezones: [
      { id: 'UTC', label: 'UTC (Tempo Universal Coordenado)', offset: 'UTC+00:00' },
    ]
  }
}

const availableTimezones = computed(() => TIMEZONES_BY_REGION[selectedRegion.value]?.timezones || TIMEZONES_BY_REGION.south_america.timezones)

watch(() => props.timezone, (tz) => {
  if (!tz) return
  for (const [region, data] of Object.entries(TIMEZONES_BY_REGION)) {
    if (data.timezones.some(t => t.id === tz)) {
      selectedRegion.value = region
      break
    }
  }
}, { immediate: true })

const onRegionChange = (e: Event) => {
  const reg = (e.target as HTMLSelectElement).value
  selectedRegion.value = reg
  const firstTz = TIMEZONES_BY_REGION[reg]?.timezones[0]?.id || 'America/Sao_Paulo'
  emit('update:timezone', firstTz)
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.85rem; width: 100%; box-sizing: border-box;">
    <!-- Fuso Horário Section -->
    <div class="vms-flex-row" style="gap: 0.75rem; width: 100%; flex-wrap: wrap;">
      <div class="vms-form-group" style="flex: 1; min-width: 170px;">
        <label class="vms-label">Região / Continente</label>
        <select :value="selectedRegion" class="vms-select" style="height: 38px; width: 100%;" @change="onRegionChange">
          <option v-for="(val, key) in TIMEZONES_BY_REGION" :key="key" :value="key">{{ val.label }}</option>
        </select>
      </div>
      <div class="vms-form-group" style="flex: 1.5; min-width: 220px;">
        <label class="vms-label">Estado / Fuso Horário (Gravação & Eventos)</label>
        <select :value="timezone || 'America/Sao_Paulo'" class="vms-select" style="height: 38px; width: 100%;" @change="emit('update:timezone', ($event.target as HTMLSelectElement).value)">
          <option v-for="tz in availableTimezones" :key="tz.id" :value="tz.id">[{{ tz.offset }}] {{ tz.label }}</option>
        </select>
      </div>
    </div>

    <!-- Coordenadas Inputs -->
    <div class="vms-flex-row" style="gap: 0.75rem; width: 100%; flex-wrap: wrap;">
      <div class="vms-form-group" style="flex: 1; min-width: 180px;">
        <label class="vms-label">Latitude (GPS)</label>
        <input :value="latitude" type="number" step="0.000001" class="vms-auth-input vms-text-mono" placeholder="-23.550520" @input="emit('update:latitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" />
      </div>
      <div class="vms-form-group" style="flex: 1; min-width: 180px;">
        <label class="vms-label">Longitude (GPS)</label>
        <input :value="longitude" type="number" step="0.000001" class="vms-auth-input vms-text-mono" placeholder="-46.633308" @input="emit('update:longitude', parseFloat(($event.target as HTMLInputElement).value) || 0)" />
      </div>
    </div>

    <!-- Setor / Observação -->
    <div class="vms-form-group" style="width: 100%; margin: 0;">
      <label class="vms-label">Setor / Observação no Mapa (Opcional)</label>
      <input :value="locationName" class="vms-auth-input" placeholder="Ex: Portaria Leste ou Setor A" @input="emit('update:locationName', ($event.target as HTMLInputElement).value)" />
    </div>

    <!-- HUD Status -->
    <div style="background: rgba(0, 240, 255, 0.04); border: 1px solid rgba(0, 240, 255, 0.15); border-radius: 4px; padding: 0.5rem 0.75rem; display: flex; align-items: center; justify-content: space-between;">
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <div style="width: 8px; height: 8px; border-radius: 50%; background: var(--vms-neu-accent-cyan);"></div>
        <span class="vms-text-mono vms-text-2xs" style="color: #fff;">FUSO ATIVO: <strong style="color: var(--vms-neu-accent-cyan);">{{ timezone || 'America/Sao_Paulo' }}</strong></span>
      </div>
      <span class="vms-badge" style="background: rgba(0, 240, 255, 0.1); color: var(--vms-neu-accent-cyan); font-size: 9px;">TIMESTAMPS SINCRONIZADOS</span>
    </div>
  </div>
</template>
