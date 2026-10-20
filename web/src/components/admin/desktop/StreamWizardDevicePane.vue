<script setup lang="ts">
defineProps<{
  brand: string
  model: string
  serialNumber: string
  firmware: string
  macAddress: string
}>()

const emit = defineEmits<{
  (e: 'update:brand', val: string): void
  (e: 'update:model', val: string): void
  (e: 'update:serialNumber', val: string): void
  (e: 'update:firmware', val: string): void
  (e: 'update:macAddress', val: string): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.85rem; width: 100%; box-sizing: border-box;">
    <div class="vms-flex-between" style="align-items: center; border-bottom: 1px solid rgba(255, 255, 255, 0.08); padding-bottom: 0.4rem;">
      <div>
        <h4 class="vms-h4" style="margin: 0; color: #fff;">IDENTIFICAÇÃO DO DISPOSITIVO // METADADOS TÉCNICOS</h4>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">Preenchimento automático via ONVIF Probe ou assinatura de fluxo RTSP</span>
      </div>
      <span class="vms-badge" :style="{ background: brand ? 'rgba(0, 255, 157, 0.1)' : 'rgba(255, 94, 58, 0.1)', color: brand ? 'var(--vms-neu-accent-green)' : 'var(--vms-neu-accent-orange)', fontSize: '9px' }">
        {{ brand ? '[DETECTADO: ' + brand + ']' : '[AGUARDANDO PROBE]' }}
      </span>
    </div>

    <div class="vms-flex-row" style="gap: 0.6rem;">
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Fabricante / Marca</label>
        <input 
          :value="brand" 
          list="device-brand-options"
          class="vms-auth-input" 
          placeholder="Ex: Hikvision, Intelbras, Tapo" 
          @input="emit('update:brand', ($event.target as HTMLInputElement).value)" 
        />
        <datalist id="device-brand-options">
          <option value="Hikvision" />
          <option value="Intelbras" />
          <option value="Dahua" />
          <option value="TP-Link Tapo" />
          <option value="Reolink" />
          <option value="Axis Communications" />
          <option value="Hanwha Vision" />
          <option value="Uniview" />
          <option value="Tecvoz" />
          <option value="Generic / Outro" />
        </datalist>
      </div>
      <div class="vms-form-group" style="flex: 1.2;">
        <label class="vms-label">Modelo do Dispositivo</label>
        <input :value="model" class="vms-auth-input" placeholder="Ex: DS-2CD Series ou VIP 3230 B" @input="emit('update:model', ($event.target as HTMLInputElement).value)" />
      </div>
    </div>

    <div class="vms-flex-row" style="gap: 0.6rem;">
      <div class="vms-form-group" style="flex: 1.2;">
        <label class="vms-label">Número de Série (SN Forense)</label>
        <input :value="serialNumber" class="vms-auth-input vms-text-mono" placeholder="Ex: DS-2CD2043G2-I20210518" @input="emit('update:serialNumber', ($event.target as HTMLInputElement).value)" />
      </div>
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Versão de Firmware</label>
        <input :value="firmware" class="vms-auth-input vms-text-mono" placeholder="Ex: V5.5.80_Build210518" @input="emit('update:firmware', ($event.target as HTMLInputElement).value)" />
      </div>
    </div>

    <div class="vms-form-group" style="margin: 0;">
      <label class="vms-label">Endereço Físico (MAC Address)</label>
      <input :value="macAddress" class="vms-auth-input vms-text-mono" placeholder="Ex: 3C:52:A1:8B:4F:10" @input="emit('update:macAddress', ($event.target as HTMLInputElement).value)" />
    </div>

    <div style="background: rgba(0, 255, 157, 0.05); border: 1px solid rgba(0, 255, 157, 0.2); padding: 0.45rem 0.75rem; border-radius: 4px; font-size: 11px; display: flex; justify-content: space-between; align-items: center;">
      <span style="color: var(--vms-text-dim);">AUDITORIA FORENSE & LGPD:</span>
      <span style="color: var(--vms-neu-accent-green); font-family: var(--vms-font-jetbrains); font-size: 10px;">METADADOS VINCULADOS AO CERTIFICADO DE EXPORTAÇÃO</span>
    </div>
  </div>
</template>
