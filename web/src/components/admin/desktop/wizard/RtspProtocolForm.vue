<script setup lang="ts">
import { buildRtspUrl, parseRtspUrl } from '../../../../constants/rtspPresets'

interface Props { url?: string; ip?: string; port?: number; path?: string; user?: string; pass?: string }
const props = withDefaults(defineProps<Props>(), { url: '', ip: '', port: 554, path: '', user: '', pass: '' })
const emit = defineEmits<{
  (e: 'update:url', val: string): void; (e: 'update:ip', val: string): void; (e: 'update:port', val: number): void
  (e: 'update:path', val: string): void; (e: 'update:user', val: string): void; (e: 'update:pass', val: string): void
}>()

const sync = (ipVal = props.ip, portVal = props.port, pathVal = props.path, userVal = props.user, passVal = props.pass) => {
  emit('update:url', buildRtspUrl(undefined, ipVal, portVal, pathVal, userVal, passVal))
}

const onUrlInput = (e: Event) => {
  const v = (e.target as HTMLInputElement).value; emit('update:url', v)
  const parsed = parseRtspUrl(v)
  if (parsed) {
    if (parsed.ip) emit('update:ip', parsed.ip); if (parsed.port) emit('update:port', parsed.port)
    if (parsed.path) emit('update:path', parsed.path); if (parsed.user) emit('update:user', parsed.user)
    if (parsed.pass) emit('update:pass', parsed.pass)
  }
}

const onPath = (e: Event) => {
  const v = (e.target as HTMLInputElement).value; emit('update:path', v)
  sync(props.ip, props.port, v, props.user, props.pass)
}

const onIp = (e: Event) => { const v = (e.target as HTMLInputElement).value; emit('update:ip', v); sync(v, props.port, props.path, props.user, props.pass) }
const onPort = (e: Event) => { const v = parseInt((e.target as HTMLInputElement).value) || 554; emit('update:port', v); sync(props.ip, v, props.path, props.user, props.pass) }
const onUser = (e: Event) => { const v = (e.target as HTMLInputElement).value; emit('update:user', v); sync(props.ip, props.port, props.path, v, props.pass) }
const onPass = (e: Event) => { const v = (e.target as HTMLInputElement).value; emit('update:pass', v); sync(props.ip, props.port, props.path, props.user, v) }
</script>

<template>
  <div class="vms-flex-col" style="min-height: 180px; gap: 0.55rem; box-sizing: border-box;">
    <div class="vms-form-group" style="width: 100%;">
      <label class="vms-label">URL RTSP Principal</label>
      <input :value="url" class="vms-auth-input vms-text-mono" placeholder="rtsp://admin:pass@192.168.1.100:554/live" @input="onUrlInput" />
    </div>
    <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
      <div class="vms-form-group" style="flex: 1.1 1 130px; min-width: 110px;">
        <label class="vms-label">IP da Câmera</label>
        <input :value="ip" class="vms-auth-input" placeholder="192.168.1.100" @input="onIp" />
      </div>
      <div class="vms-form-group" style="flex: 1.2 1 130px; min-width: 110px;">
        <label class="vms-label">Path do Streaming</label>
        <input :value="path" class="vms-auth-input vms-text-mono" placeholder="/live" @input="onPath" />
      </div>
      <div class="vms-form-group" style="width: 80px; flex-shrink: 0;">
        <label class="vms-label">Porta</label>
        <input :value="port" type="number" class="vms-auth-input" @input="onPort" />
      </div>
    </div>
    <div class="vms-flex-row" style="gap: 0.5rem; flex-wrap: wrap;">
      <div class="vms-form-group" style="flex: 1 1 130px; min-width: 110px;">
        <label class="vms-label">Usuário</label>
        <input :value="user" class="vms-auth-input" placeholder="admin" @input="onUser" />
      </div>
      <div class="vms-form-group" style="flex: 1 1 130px; min-width: 110px;">
        <label class="vms-label">Senha</label>
        <input :value="pass" type="password" class="vms-auth-input" placeholder="••••••••" @input="onPass" />
      </div>
    </div>
  </div>
</template>
