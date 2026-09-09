<script setup lang="ts">
interface Props {
  ip?: string; port?: number; user?: string; pass?: string
}
const props = withDefaults(defineProps<Props>(), {
  ip: '192.168.1.145', port: 80, user: 'admin', pass: ''
})

const emit = defineEmits<{
  (e: 'update:ip', val: string): void
  (e: 'update:port', val: number): void
  (e: 'update:user', val: string): void
  (e: 'update:pass', val: string): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="min-height: 180px; gap: 0.55rem; justify-content: space-between;">
    <div class="vms-flex-row" style="gap: 0.5rem;">
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">IP da Câmera (ONVIF Device)</label>
        <input :value="ip" class="vms-auth-input" placeholder="192.168.1.145" @input="emit('update:ip', ($event.target as HTMLInputElement).value)" />
      </div>
      <div class="vms-form-group" style="width: 100px;">
        <label class="vms-label">Porta ONVIF</label>
        <input :value="port" type="number" class="vms-auth-input" @input="emit('update:port', parseInt(($event.target as HTMLInputElement).value) || 80)" />
      </div>
    </div>
    <div class="vms-flex-row" style="gap: 0.5rem;">
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Usuário ONVIF</label>
        <input :value="user" class="vms-auth-input" placeholder="admin" @input="emit('update:user', ($event.target as HTMLInputElement).value)" />
      </div>
      <div class="vms-form-group" style="flex: 1;">
        <label class="vms-label">Senha ONVIF (Digest Auth)</label>
        <input :value="pass" type="password" class="vms-auth-input" placeholder="••••••••" @input="emit('update:pass', ($event.target as HTMLInputElement).value)" />
      </div>
    </div>
    <div style="background: rgba(0, 255, 157, 0.08); color: var(--vms-neu-accent-green); border: 1px solid rgba(0, 255, 157, 0.2); font-size: 11px; padding: 0.45rem 0.65rem; border-radius: 4px; line-height: 1.4;">
      <strong style="color: #fff;">AUTO-DISCOVERY:</strong> Perfis Main/Sub e PTZ são extraídos automaticamente via SOAP no teste.
    </div>
  </div>
</template>
