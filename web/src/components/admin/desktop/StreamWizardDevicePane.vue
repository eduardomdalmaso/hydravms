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
  <div class="vms-flex-col" style="min-height: 380px; justify-content: space-between; gap: 0.85rem;">
    <div class="vms-flex-col" style="gap: 0.85rem;">
      <div class="vms-flex-between" style="align-items: center; border-bottom: 1px solid rgba(255, 255, 255, 0.08); padding-bottom: 0.4rem;">
        <div>
          <h4 class="vms-h4" style="margin: 0; color: #fff;">IDENTIFICAÇÃO DO HARDWARE // DISPOSITIVO</h4>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">Metadados técnicos do fabricante, modelo e número de série forense</span>
        </div>
        <span class="vms-badge" style="background: rgba(0, 240, 255, 0.1); color: var(--vms-neu-accent-cyan); font-size: 9px;">
          INVENTARIO VMS ATIVO
        </span>
      </div>

      <div class="vms-flex-row" style="gap: 0.6rem;">
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label">Fabricante / Marca</label>
          <select :value="brand" class="vms-auth-input" @change="emit('update:brand', ($event.target as HTMLSelectElement).value)">
            <option value="Intelbras">Intelbras</option>
            <option value="Hikvision">Hikvision</option>
            <option value="Dahua">Dahua</option>
            <option value="Axis Communications">Axis Communications</option>
            <option value="Hanwha Vision">Hanwha Vision</option>
            <option value="Uniview">Uniview</option>
            <option value="Generic / Outro">Generic / Outro</option>
          </select>
        </div>
        <div class="vms-form-group" style="flex: 1.2;">
          <label class="vms-label">Modelo do Dispositivo</label>
          <input :value="model" class="vms-auth-input" placeholder="Ex: VIP 3230 B ou DS-2CD2043G2" @input="emit('update:model', ($event.target as HTMLInputElement).value)" />
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.6rem;">
        <div class="vms-form-group" style="flex: 1.2;">
          <label class="vms-label">Número de Série (SN)</label>
          <input :value="serialNumber" class="vms-auth-input vms-text-mono" placeholder="Ex: SN-94820194812" @input="emit('update:serialNumber', ($event.target as HTMLInputElement).value)" />
        </div>
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label">Versão de Firmware</label>
          <input :value="firmware" class="vms-auth-input vms-text-mono" placeholder="Ex: V5.5.80_Build210518" @input="emit('update:firmware', ($event.target as HTMLInputElement).value)" />
        </div>
      </div>

      <div class="vms-form-group">
        <label class="vms-label">Endereço Físico (MAC Address)</label>
        <input :value="macAddress" class="vms-auth-input vms-text-mono" placeholder="Ex: 3C:52:A1:8B:4F:10" @input="emit('update:macAddress', ($event.target as HTMLInputElement).value)" />
      </div>
    </div>

    <div style="background: rgba(0, 255, 157, 0.06); border: 1px solid rgba(0, 255, 157, 0.2); padding: 0.45rem 0.75rem; border-radius: 4px; font-size: 11px; display: flex; justify-content: space-between; align-items: center;">
      <span style="color: var(--vms-text-dim);">AUDITORIA FORENSE & LGPD:</span>
      <span style="color: var(--vms-neu-accent-green); font-family: var(--vms-font-jetbrains); font-size: 10px;">METADADOS VINCULADOS AO CERTIFICADO DE EXPORTACAO</span>
    </div>
  </div>
</template>
