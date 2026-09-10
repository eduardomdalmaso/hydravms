<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FolderNode, StreamItem } from '../../types/streamTree'
import { probeCameraSnapshot } from '../../services/adminApi'; import { generateChannelPaths } from '../../constants/rtspPresets'
import StreamWizardNetworkPane from './desktop/StreamWizardNetworkPane.vue'; import StreamWizardChannelsPane, { type ChannelItem } from './desktop/StreamWizardChannelsPane.vue'
import StreamWizardDevicePane from './desktop/StreamWizardDevicePane.vue'; import StreamWizardGeoPane from './desktop/StreamWizardGeoPane.vue'


const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: FolderNode[]; initialData?: Partial<StreamItem> }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', stream: Partial<StreamItem>, folderId: string): void }>()
const step = ref(1), isTesting = ref(false), hasSnapshot = ref(false), snapshotUrl = ref<string | undefined>(undefined), authRequired = ref(false)
const detectedCodec = ref('H.265'), detectedRes = ref('1920x1080 Full HD'), detectedFps = ref(30), latency = ref(1)
const channels = ref<ChannelItem[]>([{ id: 1, name: 'Canal 01 (Principal)', path: '/stream1', subPath: '/stream2', status: 'online' }])
const form = ref({
  name: '', protocol: 'RTSP' as 'RTSP' | 'RTMP' | 'ONVIF', url: 'rtsp://192.168.1.100:554/live', path: '/live', subUrl: '',
  ip: '192.168.1.100', port: 554, user: '', pass: '', streamKey: 'stream_alpha_01', folderId: '',
  brand: 'Intelbras', model: 'VIP 3230 B', serialNumber: 'SN-94820194812', firmware: 'V5.5.80', macAddress: '3C:52:A1:8B:4F:10',
  latitude: -23.550520, longitude: -46.633308, locationName: ''
})

watch(() => props.isOpen, (open) => {
  if (open) {
    step.value = 1; isTesting.value = false; hasSnapshot.value = false; snapshotUrl.value = undefined; authRequired.value = false
    const init = props.initialData; form.value.name = init?.name || ''; form.value.protocol = init?.protocol || 'RTSP'
    form.value.url = init?.url || 'rtsp://192.168.1.100:554/live'; form.value.path = init?.url?.includes('/stream') ? '/stream1' : '/live'
    form.value.ip = init?.ip || (form.value.protocol === 'ONVIF' ? '192.168.1.145' : '192.168.1.100'); form.value.port = init?.port || (form.value.protocol === 'ONVIF' ? 80 : 554)
    form.value.folderId = props.targetFolderId || (props.folders[0]?.id || ''); form.value.locationName = init?.name ? `${init.name} - Setor Monitorado` : ''
    const chPaths = generateChannelPaths(form.value.path, 1)
    channels.value = [{ id: 1, name: `${form.value.name || 'Canal 01'} (Principal)`, path: chPaths.main, subPath: chPaths.sub, status: 'online' }]
    if (init) fetchSnapshot()
  }
})

const fetchSnapshot = async (cb?: () => void) => {
  isTesting.value = true
  try {
    let targetUrl = form.value.url; const auth = form.value.user ? `${form.value.user}:${form.value.pass}@` : ''
    if (form.value.protocol === 'ONVIF' || !targetUrl.includes(form.value.ip)) targetUrl = `rtsp://${auth}${form.value.ip}:554/stream1`
    else if (auth && !targetUrl.includes('@')) targetUrl = targetUrl.replace('rtsp://', `rtsp://${auth}`)
    const res = await probeCameraSnapshot({ ip: form.value.ip, port: form.value.port, user: form.value.user, password: form.value.pass, url: targetUrl, protocol: form.value.protocol })
    hasSnapshot.value = true; authRequired.value = !!res.auth_required; snapshotUrl.value = res.snapshot_url || undefined
    detectedCodec.value = res.codec || (form.value.protocol === 'ONVIF' ? 'H.265 (HEVC)' : 'H.264'); detectedRes.value = res.resolution || '1920x1080 Full HD'; detectedFps.value = res.fps || 30; latency.value = res.latency_ms || 1
    if (res.manufacturer) form.value.brand = res.manufacturer; if (res.model) form.value.model = res.model; if (res.firmware) form.value.firmware = res.firmware; if (res.serial_number) form.value.serialNumber = res.serial_number
    if (res.rtsp_url) {
      form.value.url = res.rtsp_url; const p = res.rtsp_url.split('/'); if (p.length > 3) form.value.path = '/' + p.slice(3).join('/')
      const ch = generateChannelPaths(form.value.path, 1); channels.value = [{ id: 1, name: `${form.value.name || 'Canal 01'} (Principal)`, path: ch.main, subPath: ch.sub, status: 'online' }]
    }
    if (cb) cb()
  } finally { isTesting.value = false }
}

const handleNext = () => {
  if (!form.value.name.trim()) form.value.name = `Camera ${form.value.protocol} ${form.value.ip}`
  if (step.value === 1) {
    const ch = generateChannelPaths(form.value.path || '/stream1', 1); channels.value = [{ id: 1, name: `${form.value.name} (Principal)`, path: ch.main, subPath: ch.sub, status: 'online' }]
    if (!hasSnapshot.value) { fetchSnapshot(() => { step.value++ }); return }
  }
  step.value++
}

const finish = () => {
  if (!form.value.name.trim()) form.value.name = `Camera ${form.value.protocol} ${form.value.ip}`
  let finalUrl = form.value.protocol === 'RTMP' ? `rtmp://localhost:1935/live/${form.value.streamKey}` : form.value.url
  if (form.value.user && !finalUrl.includes('@')) finalUrl = finalUrl.replace('rtsp://', `rtsp://${form.value.user}:${form.value.pass}@`)
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
    <div class="vms-modal-dialog" style="max-width: 1060px; width: 1060px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3">NOVO FLUXO // ETAPA {{ step }} DE 4</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" @click="emit('close')"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
      </div>
      <div class="vms-flex-between" style="padding: 0.5rem 1.25rem; background: #0c0f14; border-bottom: 1px solid var(--vms-border);">
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 1 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">1. REDE & SNAPSHOT</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 2 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">2. CANAIS</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 3 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">3. DISPOSITIVO</span>
        <span class="vms-text-mono vms-text-2xs" :style="{ color: step >= 4 ? 'var(--vms-neu-accent-orange)' : 'var(--vms-text-dim)' }">4. GEOLOCALIZACAO</span>
      </div>
      <div class="vms-modal-body" style="padding: 1.25rem; min-height: 420px; overflow-x: hidden; box-sizing: border-box;">
        <StreamWizardNetworkPane v-if="step === 1" :form="form" :folders="folders" :is-testing="isTesting" :has-snapshot="hasSnapshot" :snapshot-url="snapshotUrl" :auth-required="authRequired" :detected-codec="detectedCodec" :detected-resolution="detectedRes" :detected-fps="detectedFps" :latency-ms="latency" @test="fetchSnapshot" @reset-snapshot="hasSnapshot = false; snapshotUrl = undefined; authRequired = false" />
        <StreamWizardChannelsPane v-else-if="step === 2" v-model:channels="channels" :base-path="form.path" :is-onvif="form.protocol === 'ONVIF'" />
        <StreamWizardDevicePane v-else-if="step === 3" v-model:brand="form.brand" v-model:model="form.model" v-model:serial-number="form.serialNumber" v-model:firmware="form.firmware" v-model:mac-address="form.macAddress" />
        <StreamWizardGeoPane v-else-if="step === 4" v-model:latitude="form.latitude" v-model:longitude="form.longitude" v-model:location-name="form.locationName" />
      </div>
      <div class="vms-modal-footer">
        <button v-if="step > 1" class="vms-btn vms-btn-secondary" @click="step--">ANTERIOR</button>
        <button v-if="step < 4" class="vms-btn vms-btn-primary" :disabled="isTesting" @click="handleNext"><span v-if="isTesting">VALIDANDO...</span><span v-else>PROXIMO</span></button>
        <button v-else class="vms-btn vms-btn-primary" @click="finish">SALVAR</button>
      </div>
    </div>
  </div>
</template>
