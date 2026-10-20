<script setup lang="ts">
import { ref, computed } from 'vue'
import { generateChannelPaths, deriveSubStreamUrl } from '../../../constants/rtspPresets'
import { probeCameraSnapshot } from '../../../services/adminApi'
import { copyToClipboard } from '../../../utils/streamUrls'

export interface ChannelItem {
  id: number; name: string; path: string; subPath?: string; status: 'online' | 'offline' | 'testing' | 'pending'
  mainStatus?: 'online' | 'offline' | 'testing' | 'pending'; subStatus?: 'online' | 'offline' | 'testing' | 'pending'
}

const props = defineProps<{
  channels: ChannelItem[]; basePath: string; isOnvif?: boolean
  form: { ip: string; port: number; user: string; pass: string; protocol: string; url: string }
}>()

const emit = defineEmits<{ (e: 'update:channels', val: ChannelItem[]): void }>()
const isTestingAll = ref(false), copiedNotice = ref<string | null>(null)

const totalStreams = computed(() => props.channels.reduce((acc, ch) => acc + (ch.path ? 1 : 0) + (ch.subPath ? 1 : 0), 0))
const activeStreams = computed(() => props.channels.reduce((acc, ch) => {
  const m = (ch.mainStatus === 'online' || (!ch.mainStatus && ch.status === 'online')) ? 1 : 0
  const s = (ch.subStatus === 'online' || (!ch.subStatus && ch.status === 'online')) ? 1 : 0
  return acc + m + s
}, 0))

const addChannel = () => {
  const nextId = props.channels.length + 1
  const paths = generateChannelPaths(props.basePath || '/stream1', nextId)
  emit('update:channels', [...props.channels, {
    id: nextId, name: `Canal ${nextId.toString().padStart(2, '0')}`,
    path: paths.main, subPath: paths.sub, status: 'pending'
  }])
}

const removeChannel = (idx: number) => {
  if (props.channels.length <= 1) return
  emit('update:channels', props.channels.filter((_, i) => i !== idx))
}

const quickBatchAdd = (count: number) => {
  const list: ChannelItem[] = []
  for (let i = 1; i <= count; i++) {
    const paths = generateChannelPaths(props.basePath || '/stream1', i)
    list.push({ id: i, name: `Canal ${i.toString().padStart(2, '0')}`, path: paths.main, subPath: paths.sub, status: 'pending' })
  }
  emit('update:channels', list)
}

const getChannelUrl = (ch: ChannelItem, isSub = false) => {
  const mainPath = ch.path || '/live'
  const subPath = ch.subPath || (props.basePath ? generateChannelPaths(props.basePath, ch.id).sub : '')
  const auth = props.form.user ? (props.form.pass ? `${props.form.user}:${props.form.pass}@` : `${props.form.user}@`) : ''
  const host = props.form.ip || '127.0.0.1', port = props.form.port || 554
  const main = `rtsp://${auth}${host}:${port}${mainPath.startsWith('/') ? mainPath : '/' + mainPath}`
  return isSub ? (subPath ? `rtsp://${auth}${host}:${port}${subPath.startsWith('/') ? subPath : '/' + subPath}` : deriveSubStreamUrl(main)) : main
}

const copyChannelUrl = async (ch: ChannelItem, isSub = false) => {
  if (await copyToClipboard(getChannelUrl(ch, isSub))) {
    copiedNotice.value = `CH #${ch.id} ${isSub ? 'SUB' : 'MAIN'} COPIADO!`
    setTimeout(() => { copiedNotice.value = null }, 1800)
  }
}

const testSingle = async (ch: ChannelItem) => {
  ch.status = 'testing'; ch.mainStatus = 'testing'; ch.subStatus = 'testing'
  try {
    const mainUrl = getChannelUrl(ch, false), subUrl = getChannelUrl(ch, true)
    const [resMain, resSub] = await Promise.all([
      probeCameraSnapshot({ ip: props.form.ip, port: props.form.port, user: props.form.user, password: props.form.pass, url: mainUrl, protocol: 'RTSP', stream_id: `probe_ch_${ch.id}_main` }),
      probeCameraSnapshot({ ip: props.form.ip, port: props.form.port, user: props.form.user, password: props.form.pass, url: subUrl, protocol: 'RTSP', stream_id: `probe_ch_${ch.id}_sub` })
    ])

    ch.mainStatus = (resMain?.online && !resMain.auth_required) ? 'online' : 'offline'
    ch.subStatus = (resSub?.online && !resSub.auth_required) ? 'online' : 'offline'
    ch.status = (ch.mainStatus === 'online' && ch.subStatus === 'online') ? 'online' : (ch.mainStatus === 'online' || ch.subStatus === 'online' ? 'online' : 'offline')
  } catch {
    ch.status = 'offline'; ch.mainStatus = 'offline'; ch.subStatus = 'offline'
  }
}

const testAll = async () => {
  isTestingAll.value = true
  await Promise.all(props.channels.map(ch => testSingle(ch)))
  isTestingAll.value = false
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.65rem; width: 100%; box-sizing: border-box;">
    <div class="vms-flex-between" style="align-items: center;">
      <div>
        <h4 class="vms-h4" style="margin: 0; color: #fff;">TOPOLOGIA DE CANAIS // VALIDAÇÃO MAIN & SUB</h4>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">
          <strong v-if="isOnvif" style="color: var(--vms-neu-accent-green);">[ONVIF AUTO-DISCOVERY]</strong>
          <span v-else>Teste real RFC 2326 TCP em ambos os fluxos (Main + Sub)</span>
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.35rem;">
        <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(4)">+4</button>
        <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(8)">+8</button>
        <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(16)">+16</button>
        <button class="vms-btn vms-btn-primary" style="font-size: 10px; padding: 3px 8px;" type="button" @click="addChannel">+ CANAL</button>
        <button class="vms-btn vms-btn-primary" style="font-size: 10px; padding: 3px 8px; background: var(--vms-neu-accent-cyan); color: #000;" type="button" :disabled="isTestingAll" @click="testAll">
          {{ isTestingAll ? 'TESTANDO...' : 'TESTAR TODOS' }}
        </button>
      </div>
    </div>

    <!-- Channels List Container -->
    <div style="max-height: 290px; overflow-y: auto; display: flex; flex-direction: column; gap: 0.35rem; padding-right: 4px;">
      <div v-for="(ch, idx) in channels" :key="ch.id" style="display: flex; align-items: center; gap: 0.5rem; background: #07090e; border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 4px; padding: 0.4rem 0.6rem;">
        <div style="width: 22px; font-family: var(--vms-font-jetbrains); font-size: 11px; font-weight: 700; color: var(--vms-neu-accent-cyan); cursor: pointer;" title="Clique para copiar RTSP deste Canal" @click="copyChannelUrl(ch)">#{{ ch.id.toString().padStart(2, '0') }}</div>
        <div class="vms-form-group" style="flex: 1.1; margin: 0;"><input v-model="ch.name" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;" placeholder="Nome Canal" /></div>
        <div class="vms-form-group" style="flex: 1.4; margin: 0;"><input v-model="ch.path" class="vms-auth-input vms-text-mono" style="padding: 4px 8px; font-size: 11px; color: var(--vms-neu-accent-cyan); cursor: pointer;" title="Clique para copiar URL RTSP Main" placeholder="Path Main (ex: /401)" @click.ctrl="copyChannelUrl(ch, false)" /></div>
        <div class="vms-form-group" style="flex: 1.4; margin: 0;"><input v-model="ch.subPath" class="vms-auth-input vms-text-mono" style="padding: 4px 8px; font-size: 11px; color: var(--vms-text-dim); cursor: pointer;" title="Clique para copiar URL RTSP Sub" placeholder="Path Sub (ex: /402)" @click.ctrl="copyChannelUrl(ch, true)" /></div>
        <div style="min-width: 145px; display: flex; align-items: center; justify-content: flex-end; gap: 0.3rem; flex-shrink: 0;">
          <template v-if="ch.status === 'testing'">
            <span class="vms-badge" style="color: var(--vms-neu-accent-cyan); font-size: 9px;">TESTANDO...</span>
          </template>
          <template v-else-if="ch.mainStatus || ch.subStatus">
            <span class="vms-badge" :style="{ color: ch.mainStatus === 'online' ? 'var(--vms-neu-accent-green)' : 'var(--vms-neu-accent-red)', background: ch.mainStatus === 'online' ? 'rgba(0,255,157,0.1)' : 'rgba(255,0,60,0.1)', border: '1px solid currentColor', fontSize: '9px', padding: '1px 4px' }">M:{{ ch.mainStatus === 'online' ? 'OK' : 'OFF' }}</span>
            <span class="vms-badge" :style="{ color: ch.subStatus === 'online' ? 'var(--vms-neu-accent-cyan)' : 'var(--vms-neu-accent-red)', background: ch.subStatus === 'online' ? 'rgba(0,240,255,0.1)' : 'rgba(255,0,60,0.1)', border: '1px solid currentColor', fontSize: '9px', padding: '1px 4px' }">S:{{ ch.subStatus === 'online' ? 'OK' : 'OFF' }}</span>
            <button type="button" class="vms-btn vms-btn-ghost" style="padding: 1px 3px; font-size: 10px;" title="Copiar RTSP Main" @click="copyChannelUrl(ch, false)">📋</button>
            <button type="button" class="vms-btn vms-btn-ghost" style="padding: 1px 4px; font-size: 9px;" title="Retestar" @click="testSingle(ch)">↻</button>
          </template>
          <template v-else>
            <span v-if="ch.status === 'online'" class="vms-badge" style="color: var(--vms-neu-accent-green); background: rgba(0,255,157,0.1); border: 1px solid rgba(0,255,157,0.3); font-size: 9px;">[ONLINE]</span>
            <button v-else type="button" class="vms-btn vms-btn-secondary" style="font-size: 9px; padding: 2px 6px;" @click="testSingle(ch)">TESTAR</button>
            <button type="button" class="vms-btn vms-btn-ghost" style="padding: 1px 3px; font-size: 10px;" title="Copiar RTSP Main" @click="copyChannelUrl(ch, false)">📋</button>
          </template>
        </div>
        <button v-if="channels.length > 1" class="vms-btn vms-btn-ghost" style="color: var(--vms-neu-accent-red); padding: 2px 5px;" type="button" title="Remover" @click="removeChannel(idx)">✕</button>
      </div>
    </div>

    <div style="background: rgba(0, 240, 255, 0.04); border: 1px solid rgba(0, 240, 255, 0.15); padding: 0.35rem 0.65rem; border-radius: 4px; font-size: 11px; display: flex; justify-content: space-between; align-items: center;">
      <span style="color: var(--vms-text-dim); font-family: var(--vms-font-jetbrains);">
        FLUXOS ATIVOS: <strong style="color: #fff;">{{ activeStreams }} / {{ totalStreams }}</strong>
        <span style="margin-left: 0.4rem; font-size: 10px; color: var(--vms-text-dim);">({{ channels.length }} {{ channels.length === 1 ? 'CANAL' : 'CANAIS' }} // MAIN &amp; SUB)</span>
      </span>
      <span v-if="copiedNotice" class="vms-badge" style="color: var(--vms-neu-accent-green); background: rgba(0,255,157,0.15); border: 1px solid rgba(0,255,157,0.4); font-size: 10px; padding: 1px 6px;">{{ copiedNotice }}</span>
      <span v-else style="color: var(--vms-neu-accent-green); font-size: 10px; font-family: var(--vms-font-jetbrains);">PROBE REAL RFC 2326 TCP ATIVO</span>
    </div>
  </div>
</template>
