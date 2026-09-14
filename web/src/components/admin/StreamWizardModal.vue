<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { FolderNode, StreamItem } from '../../types/streamTree'
import { probeCameraSnapshot } from '../../services/adminApi'; import { generateChannelPaths } from '../../constants/rtspPresets'
import StreamWizardNetworkPane from './desktop/StreamWizardNetworkPane.vue'; import StreamWizardChannelsPane, { type ChannelItem } from './desktop/StreamWizardChannelsPane.vue'
import StreamWizardDevicePane from './desktop/StreamWizardDevicePane.vue'; import StreamWizardGeoPane from './desktop/StreamWizardGeoPane.vue'

const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: FolderNode[]; initialData?: Partial<StreamItem> }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', stream: Partial<StreamItem>, folderId: string): void }>()
const step = ref(1), isTesting = ref(false), hasSnapshot = ref(false), snapshotUrl = ref<string | undefined>(undefined), authRequired = ref(false), testError = ref<string | null>(null)
const detectedCodec = ref('H.265'), detectedRes = ref('1920x1080 Full HD'), detectedFps = ref(30), latency = ref(1)
const channels = ref<ChannelItem[]>([{ id: 1, name: 'Canal 01 (Principal)', path: '/stream1', subPath: '/stream2', status: 'online' }])
const form = ref({
  name: '', protocol: 'RTSP' as 'RTSP' | 'RTMP' | 'ONVIF' | 'LOOP', url: 'rtsp://192.168.1.100:554/live', path: '/live', subUrl: '',
  ip: '192.168.1.100', port: 554, user: '', pass: '', streamKey: 'stream_alpha_01', folderId: '',
  brand: 'Intelbras', model: 'VIP 3230 B', serialNumber: 'SN-94820194812', firmware: 'V5.5.80', macAddress: '3C:52:A1:8B:4F:10',
  latitude: -23.550520, longitude: -46.633308, locationName: ''
})

const wizardSteps = computed(() => {
  if (form.value.protocol === 'LOOP') return [{ key: 'network', title: '1. ARQUIVO & SNAPSHOT' }, { key: 'geo', title: '2. GEOLOCALIZACAO' }]
  if (form.value.protocol === 'RTMP') return [{ key: 'network', title: '1. REDE & SNAPSHOT' }, { key: 'device', title: '2. DISPOSITIVO' }, { key: 'geo', title: '3. GEOLOCALIZACAO' }]
  return [{ key: 'network', title: '1. REDE & SNAPSHOT' }, { key: 'channels', title: '2. CANAIS' }, { key: 'device', title: '3. DISPOSITIVO' }, { key: 'geo', title: '4. GEOLOCALIZACAO' }]
})
const currentStepKey = computed(() => wizardSteps.value[step.value - 1]?.key || 'network')
const totalSteps = computed(() => wizardSteps.value.length)

watch(() => props.isOpen, (open) => {
  if (open) {
    step.value = 1; isTesting.value = false; hasSnapshot.value = false; snapshotUrl.value = undefined; authRequired.value = false; testError.value = null
    const init = props.initialData; form.value.name = init?.name || ''; form.value.protocol = (init?.protocol as any) || 'RTSP'
    form.value.url = init?.url || 'rtsp://192.168.1.100:554/live'; form.value.path = init?.url?.includes('/stream') ? '/stream1' : '/live'
    form.value.ip = init?.ip || (form.value.protocol === 'ONVIF' ? '192.168.1.145' : '192.168.1.100'); form.value.port = init?.port || (form.value.protocol === 'ONVIF' ? 80 : 554)
    form.value.folderId = props.targetFolderId || (props.folders[0]?.id || ''); form.value.locationName = init?.name ? `${init.name} - Setor Monitorado` : ''
    const chPaths = generateChannelPaths(form.value.path, 1)
    channels.value = [{ id: 1, name: `${form.value.name || 'Canal 01'} (Principal)`, path: chPaths.main, subPath: chPaths.sub, status: 'online' }]
    if (init) fetchSnapshot()
  }
})

const fetchSnapshot = async (): Promise<boolean> => {
  isTesting.value = true; testError.value = null
  try {
    let targetUrl = form.value.url; const auth = form.value.user ? `${form.value.user}:${form.value.pass}@` : ''
    if (form.value.protocol === 'ONVIF' || (!targetUrl.includes(form.value.ip) && form.value.protocol !== 'LOOP')) targetUrl = `rtsp://${auth}${form.value.ip}:554/stream1`
    else if (auth && !targetUrl.includes('@') && form.value.protocol !== 'LOOP') targetUrl = targetUrl.replace('rtsp://', `rtsp://${auth}`)
    const res = await probeCameraSnapshot({ ip: form.value.ip, port: form.value.port, user: form.value.user, password: form.value.pass, url: targetUrl, protocol: form.value.protocol })
    if (res && res.online) {
      hasSnapshot.value = true; authRequired.value = !!res.auth_required; snapshotUrl.value = res.snapshot_url || undefined
      testError.value = res.auth_required ? 'Autenticação necessária (401). Informe Usuário e Senha.' : null
      detectedCodec.value = res.codec || (form.value.protocol === 'ONVIF' ? 'H.265 (HEVC)' : 'H.264'); detectedRes.value = res.resolution || '1920x1080 Full HD'; detectedFps.value = res.fps || 30; latency.value = res.latency_ms || 1
      if (res.manufacturer) form.value.brand = res.manufacturer; if (res.model) form.value.model = res.model; if (res.firmware) form.value.firmware = res.firmware; if (res.serial_number) form.value.serialNumber = res.serial_number
      if (res.rtsp_url && form.value.protocol !== 'LOOP') {
        form.value.url = res.rtsp_url; const p = res.rtsp_url.split('/'); if (p.length > 3) form.value.path = '/' + p.slice(3).join('/')
        const ch = generateChannelPaths(form.value.path, 1); channels.value = [{ id: 1, name: `${form.value.name || 'Canal 01'} (Principal)`, path: ch.main, subPath: ch.sub, status: 'online' }]
      }
      return true
    } else {
      hasSnapshot.value = false; authRequired.value = false; snapshotUrl.value = undefined
      testError.value = res?.error || `Falha de conexão com ${form.value.ip || form.value.url}. Endereço inacessível.`
      return false
    }
  } catch (err: any) {
    hasSnapshot.value = false; authRequired.value = false; snapshotUrl.value = undefined
    testError.value = `Erro de teste: ${err?.message || 'Servidor inacessível'}`
    return false
  } finally { isTesting.value = false }
}

const handleNext = async () => {
  if (!form.value.name.trim()) form.value.name = form.value.protocol === 'LOOP' ? `Fluxo Loop ${form.value.url.split('/').pop() || 'Video'}` : `Camera ${form.value.protocol} ${form.value.ip}`
  if (currentStepKey.value === 'network') {
    const ch = generateChannelPaths(form.value.path || '/stream1', 1); channels.value = [{ id: 1, name: `${form.value.name} (Principal)`, path: ch.main, subPath: ch.sub, status: 'online' }]
    if (!hasSnapshot.value) { const ok = await fetchSnapshot(); if (!ok) return }
  }
  step.value++
}

const finish = () => {
  if (!form.value.name.trim()) form.value.name = form.value.protocol === 'LOOP' ? `Fluxo Loop ${form.value.url.split('/').pop() || 'Video'}` : `Camera ${form.value.protocol} ${form.value.ip}`
  let finalUrl = form.value.protocol === 'RTMP' ? `rtmp://localhost:1935/live/${form.value.streamKey}` : (form.value.protocol === 'LOOP' ? (form.value.url.startsWith('file://') ? form.value.url : `file://${form.value.url}`) : form.value.url)
  if (form.value.user && !finalUrl.includes('@') && form.value.protocol !== 'LOOP') finalUrl = finalUrl.replace('rtsp://', `rtsp://${form.value.user}:${form.value.pass}@`)
  emit('save', {
    name: form.value.name, protocol: form.value.protocol, url: finalUrl, ip: form.value.ip, port: form.value.port,
    codec: detectedCodec.value.includes('H.265') ? 'H.265' : 'H.264', resolution: '1080P', fps: detectedFps.value,
    bitrate: '4.0 Mbps', recordMode: 'disabled', status: 'online', has_ptz: form.value.protocol === 'ONVIF',
    latitude: form.value.latitude, longitude: form.value.longitude, locationName: form.value.locationName, snapshotUrl: snapshotUrl.value
  }, form.value.folderId); emit('close')
}
</script>
<template>
  <div v-if="isOpen" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 980px; width: 95vw; box-sizing: border-box;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">NOVO FLUXO // ETAPA {{ step }} DE {{ totalSteps }}</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" @click="emit('close')"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
      </div>
      <div class="vms-flex-between" style="padding: 0.5rem 1.25rem; background: #0c0f14; border-bottom: 1px solid var(--vms-border); flex-wrap: wrap; gap: 0.4rem;">
        <span v-for="(st, idx) in wizardSteps" :key="st.key" class="vms-text-mono vms-text-2xs" :style="{ color: step >= (idx + 1) ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">{{ st.title }}</span>
      </div>
      <div class="vms-modal-body" style="padding: 1.15rem; min-height: 400px; max-height: calc(85vh - 120px); overflow-y: auto; overflow-x: hidden; box-sizing: border-box;">
        <StreamWizardNetworkPane v-if="currentStepKey === 'network'" :form="form" :folders="folders" :is-testing="isTesting" :has-snapshot="hasSnapshot" :snapshot-url="snapshotUrl" :auth-required="authRequired" :test-error="testError" :detected-codec="detectedCodec" :detected-resolution="detectedRes" :detected-fps="detectedFps" :latency-ms="latency" @test="fetchSnapshot" @reset-snapshot="hasSnapshot = false; snapshotUrl = undefined; authRequired = false; testError = null" />
        <StreamWizardChannelsPane v-else-if="currentStepKey === 'channels'" v-model:channels="channels" :base-path="form.path" :is-onvif="form.protocol === 'ONVIF'" />
        <StreamWizardDevicePane v-else-if="currentStepKey === 'device'" v-model:brand="form.brand" v-model:model="form.model" v-model:serial-number="form.serialNumber" v-model:firmware="form.firmware" v-model:mac-address="form.macAddress" />
        <StreamWizardGeoPane v-else-if="currentStepKey === 'geo'" v-model:latitude="form.latitude" v-model:longitude="form.longitude" v-model:location-name="form.locationName" />
      </div>
      <div class="vms-modal-footer">
        <button v-if="step > 1" class="vms-btn vms-btn-secondary" @click="step--">ANTERIOR</button>
        <button v-if="step < totalSteps" class="vms-btn vms-btn-primary" :disabled="isTesting" @click="handleNext"><span v-if="isTesting">VALIDANDO...</span><span v-else>PROXIMO</span></button>
        <button v-else class="vms-btn vms-btn-primary" @click="finish">SALVAR</button>
      </div>
    </div>
  </div>
</template>
