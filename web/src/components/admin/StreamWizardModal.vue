<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { FolderNode, StreamItem } from '../../types/streamTree'
import { probeCameraSnapshot } from '../../services/adminApi'
import { generateChannelPaths, parseRtspUrl, deriveSubStreamUrl, detectBrandFromUrl } from '../../constants/rtspPresets'
import StreamWizardNetworkPane from './desktop/StreamWizardNetworkPane.vue'
import StreamWizardChannelsPane, { type ChannelItem } from './desktop/StreamWizardChannelsPane.vue'
import StreamWizardDevicePane from './desktop/StreamWizardDevicePane.vue'
import StreamWizardGeoPane from './desktop/StreamWizardGeoPane.vue'

const props = defineProps<{ isOpen: boolean; targetFolderId?: string; folders: FolderNode[]; initialData?: Partial<StreamItem> }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', stream: Partial<StreamItem>, folderId: string): void }>()
const step = ref(1), isTesting = ref(false), hasSnapshot = ref(false), snapshotUrl = ref<string | undefined>(undefined), authRequired = ref(false), testError = ref<string | null>(null), saveAlertError = ref<string | null>(null)
const detectedCodec = ref('H.265'), detectedRes = ref('1920x1080 Full HD'), detectedFps = ref(30), latency = ref(1)
const channels = ref<ChannelItem[]>([{ id: 1, name: 'Canal 01 (Principal)', path: '/stream1', subPath: '/stream2', status: 'online' }])
const form = ref({
  name: '', protocol: 'RTSP' as 'RTSP' | 'RTMP' | 'ONVIF' | 'LOOP', url: '', path: '', subUrl: '',
  ip: '', port: 554, user: '', pass: '', streamKey: '', folderId: '',
  brand: '', model: '', serialNumber: '', firmware: '', macAddress: '',
  latitude: -23.550520, longitude: -46.633308, locationName: '', timezone: 'America/Sao_Paulo'
})

const wizardSteps = computed(() => {
  if (form.value.protocol === 'LOOP') return [{ key: 'network', title: '1. ARQUIVO' }, { key: 'geo', title: '2. GEOLOCALIZACAO' }]
  if (form.value.protocol === 'RTMP') return [{ key: 'network', title: '1. REDE' }, { key: 'device', title: '2. DISPOSITIVO' }, { key: 'geo', title: '3. GEOLOCALIZACAO' }]
  return [{ key: 'network', title: '1. REDE' }, { key: 'channels', title: '2. CANAIS' }, { key: 'device', title: '3. DISPOSITIVO' }, { key: 'geo', title: '4. GEOLOCALIZACAO' }]
})
const currentStepKey = computed(() => wizardSteps.value[step.value - 1]?.key || 'network')
const totalSteps = computed(() => wizardSteps.value.length)

const updateChannelsFromNetwork = () => {
  const p = form.value.path || (form.value.url ? parseRtspUrl(form.value.url)?.path : '') || '/stream1'
  const subDerived = form.value.url ? (parseRtspUrl(deriveSubStreamUrl(form.value.url))?.path || '') : ''
  const ch = generateChannelPaths(p, 1)
  channels.value = [{ id: 1, name: `${form.value.name || 'Canal 01'} (Principal)`, path: ch.main, subPath: subDerived || ch.sub, status: 'online' }]
}

watch(() => props.isOpen, (open) => {
  if (open) {
    step.value = 1; isTesting.value = false; hasSnapshot.value = false; snapshotUrl.value = undefined; authRequired.value = false; testError.value = null; saveAlertError.value = null
    const init = props.initialData as any
    const parsed = (init?.url && init.protocol !== 'LOOP' && init.protocol !== 'RTMP') ? parseRtspUrl(init.url) : null
    Object.assign(form.value, {
      name: init?.name || '', protocol: (init?.protocol as any) || 'RTSP', url: init?.url || '',
      path: init?.path || parsed?.path || '', ip: init?.ip || parsed?.ip || '',
      port: init?.port || parsed?.port || (init?.protocol === 'ONVIF' ? 80 : 554),
      user: init?.user || parsed?.user || '', pass: init?.pass || parsed?.pass || '',
      streamKey: init?.streamKey || '', brand: init?.brand || '', model: init?.model || '',
      serialNumber: init?.serialNumber || '', firmware: init?.firmware || '', macAddress: init?.macAddress || '',
      latitude: typeof init?.latitude === 'number' ? init.latitude : -23.550520,
      longitude: typeof init?.longitude === 'number' ? init.longitude : -46.633308,
      locationName: init?.locationName || '', timezone: init?.timezone || 'America/Sao_Paulo',
      folderId: props.targetFolderId || init?.folderId || (props.folders[0]?.id || '')
    })
    const detected = detectBrandFromUrl(form.value.url || form.value.path)
    if (detected.brand && !form.value.brand) {
      form.value.brand = detected.brand; if (detected.model && !form.value.model) form.value.model = detected.model
    }
    updateChannelsFromNetwork()
  }
})

const goToStep = (targetStep: number) => {
  if (targetStep < 1 || targetStep > totalSteps.value) return
  saveAlertError.value = null
  if (currentStepKey.value === 'network') updateChannelsFromNetwork()
  step.value = targetStep
}

const fetchSnapshot = async (): Promise<boolean> => {
  isTesting.value = true; testError.value = null
  try {
    let targetUrl = form.value.url
    const auth = form.value.user ? `${form.value.user}:${form.value.pass}@` : ''
    if (form.value.protocol === 'ONVIF') targetUrl = `rtsp://${auth}${form.value.ip}:554/stream1`
    else if (form.value.protocol === 'RTSP') {
      if (!targetUrl && form.value.ip) targetUrl = `rtsp://${auth}${form.value.ip}:${form.value.port || 554}${form.value.path || '/live'}`
      else if (auth && !targetUrl.includes('@')) targetUrl = targetUrl.replace('rtsp://', `rtsp://${auth}`)
    }
    const streamId = form.value.name ? form.value.name.toLowerCase().replace(/[^a-z0-9_-]/g, '_') : ''
    const res = await probeCameraSnapshot({ ip: form.value.ip, port: form.value.port, user: form.value.user, password: form.value.pass, url: targetUrl, protocol: form.value.protocol, stream_id: streamId })
    if (res && res.online && !res.auth_required) {
      hasSnapshot.value = true; authRequired.value = false; snapshotUrl.value = res.snapshot_url || undefined; testError.value = null
      detectedCodec.value = res.codec || (form.value.protocol === 'ONVIF' ? 'H.265 (HEVC)' : 'H.264'); detectedRes.value = res.resolution || '1920x1080 Full HD'; detectedFps.value = res.fps || 30; latency.value = res.latency_ms || 1
      if (res.manufacturer && !form.value.brand) form.value.brand = res.manufacturer
      if (res.model && !form.value.model) form.value.model = res.model
      if (res.firmware && !form.value.firmware) form.value.firmware = res.firmware
      if (res.serial_number && !form.value.serialNumber) form.value.serialNumber = res.serial_number
      if (res.rtsp_url && form.value.protocol !== 'LOOP') {
        form.value.url = res.rtsp_url; const p = res.rtsp_url.split('/'); if (p.length > 3) form.value.path = '/' + p.slice(3).join('/')
        updateChannelsFromNetwork()
      }
      return true
    } else if (res && res.auth_required) {
      hasSnapshot.value = false; authRequired.value = true; snapshotUrl.value = undefined; testError.value = res.error || 'Autenticação necessária (401).'
      return false
    } else {
      hasSnapshot.value = false; authRequired.value = false; snapshotUrl.value = undefined; testError.value = res?.error || `Falha de conexão com ${form.value.ip || form.value.url}.`
      return false
    }
  } catch (err: any) {
    hasSnapshot.value = false; authRequired.value = false; snapshotUrl.value = undefined; testError.value = `Erro de teste: ${err?.message || 'Servidor inacessível'}`
    return false
  } finally { isTesting.value = false }
}

const handleNext = async () => {
  saveAlertError.value = null
  if (!form.value.name.trim()) form.value.name = form.value.protocol === 'LOOP' ? `Fluxo Loop ${form.value.url.split('/').pop() || 'Video'}` : `Camera ${form.value.protocol} ${form.value.ip}`
  if (currentStepKey.value === 'network') updateChannelsFromNetwork()
  step.value++
  if (currentStepKey.value === 'geo' && !hasSnapshot.value && !isTesting.value) fetchSnapshot()
}

const finish = async () => {
  if (isTesting.value) return
  saveAlertError.value = null

  if (!hasSnapshot.value || !snapshotUrl.value) {
    const ok = await fetchSnapshot()
    if (!ok || !snapshotUrl.value) {
      saveAlertError.value = testError.value || 'Não foi possível conectar com o fluxo principal.'
      return
    }
  }

  let finalUrl = form.value.protocol === 'RTMP' ? `rtmp://localhost:1935/live/${form.value.streamKey}` : (form.value.protocol === 'LOOP' ? (form.value.url.startsWith('file://') ? form.value.url : `file://${form.value.url}`) : form.value.url)
  if (form.value.user && !finalUrl.includes('@') && form.value.protocol !== 'LOOP') finalUrl = finalUrl.replace('rtsp://', `rtsp://${form.value.user}:${form.value.pass}@`)

  let finalSubUrl = ''
  if (form.value.protocol === 'RTSP' || form.value.protocol === 'ONVIF') {
    const ch1 = channels.value[0]
    if (ch1 && ch1.subPath && form.value.ip) {
      const auth = form.value.user ? (form.value.pass ? `${form.value.user}:${form.value.pass}@` : `${form.value.user}@`) : ''
      finalSubUrl = `rtsp://${auth}${form.value.ip}:${form.value.port || 554}${ch1.subPath.startsWith('/') ? ch1.subPath : '/' + ch1.subPath}`
    } else {
      finalSubUrl = deriveSubStreamUrl(finalUrl)
    }

    isTesting.value = true
    try {
      const subProbe = await probeCameraSnapshot({
        ip: form.value.ip, port: form.value.port, user: form.value.user, password: form.value.pass,
        url: finalSubUrl, protocol: form.value.protocol, stream_id: `${form.value.name || 'cam'}_sub`
      })
      if (!subProbe || !subProbe.online || subProbe.auth_required) {
        saveAlertError.value = `Sub-stream (${finalSubUrl}) inacessível. O sistema só registra fluxos com conexão ativa validada.`
        isTesting.value = false
        return
      }
    } catch (err: any) {
      saveAlertError.value = `Falha na validação da Sub-Stream: ${err?.message || 'Sem conexão'}`
      isTesting.value = false
      return
    } finally { isTesting.value = false }
  }

  emit('save', {
    id: (props.initialData as any)?.id, name: form.value.name, protocol: form.value.protocol, url: finalUrl,
    subUrl: finalSubUrl || undefined, ip: form.value.ip, port: form.value.port,
    codec: detectedCodec.value.includes('H.265') ? 'H.265' : 'H.264', resolution: detectedRes.value || '1080P', fps: detectedFps.value, bitrate: '4.0 Mbps',
    recordMode: (props.initialData as any)?.recordMode || 'disabled', status: 'online', has_ptz: form.value.protocol === 'ONVIF',
    brand: form.value.brand, model: form.value.model, serialNumber: form.value.serialNumber, firmware: form.value.firmware, macAddress: form.value.macAddress,
    latitude: form.value.latitude, longitude: form.value.longitude, locationName: form.value.locationName, timezone: form.value.timezone, snapshotUrl: snapshotUrl.value
  }, form.value.folderId)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="vms-modal-backdrop">
    <div class="vms-modal-dialog" style="width: 800px; max-width: 95vw; height: 590px; max-height: 88vh; display: flex; flex-direction: column; box-sizing: border-box;">
      <div class="vms-modal-header" style="flex-shrink: 0;">
        <h3 class="vms-h3">{{ props.initialData?.id ? 'ALTERAR CADASTRO' : 'NOVO FLUXO' }} // ETAPA {{ step }} DE {{ totalSteps }}</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" @click="emit('close')"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
      </div>

      <!-- Abas Clicáveis com Navegação Direta -->
      <div class="vms-flex-between" style="flex-shrink: 0; padding: 0.35rem 1.25rem; background: #0c0f14; border-bottom: 1px solid var(--vms-border); flex-wrap: wrap; gap: 0.4rem;">
        <button v-for="(st, idx) in wizardSteps" :key="st.key" type="button" class="vms-btn vms-btn-ghost vms-text-mono vms-text-2xs" :style="{ color: step === (idx + 1) ? 'var(--vms-neu-accent-orange)' : (props.initialData?.id || step > (idx + 1) ? '#ffffff' : 'var(--vms-text-dim)'), borderBottom: step === (idx + 1) ? '2px solid var(--vms-neu-accent-orange)' : '2px solid transparent', borderRadius: '0', padding: '4px 8px', cursor: 'pointer' }" @click="goToStep(idx + 1)">
          {{ st.title }}
        </button>
      </div>

      <!-- Corpo com Altura Fixa -->
      <div class="vms-modal-body" style="flex: 1; min-height: 0; overflow-y: auto; overflow-x: hidden; padding: 1.15rem; box-sizing: border-box;">
        <StreamWizardNetworkPane v-if="currentStepKey === 'network'" :form="form" :folders="folders" @reset-snapshot="hasSnapshot = false; snapshotUrl = undefined; authRequired = false; testError = null" />
        <StreamWizardChannelsPane v-else-if="currentStepKey === 'channels'" v-model:channels="channels" :base-path="form.path" :is-onvif="form.protocol === 'ONVIF'" :form="form" />
        <StreamWizardDevicePane v-else-if="currentStepKey === 'device'" v-model:brand="form.brand" v-model:model="form.model" v-model:serial-number="form.serialNumber" v-model:firmware="form.firmware" v-model:mac-address="form.macAddress" />
        <StreamWizardGeoPane v-else-if="currentStepKey === 'geo'" v-model:latitude="form.latitude" v-model:longitude="form.longitude" v-model:location-name="form.locationName" v-model:timezone="form.timezone" />
      </div>

      <!-- Banner de Erro -->
      <div v-if="saveAlertError" style="flex-shrink: 0; margin: 0 1.15rem 0.5rem; padding: 0.55rem 0.9rem; background: rgba(255, 0, 60, 0.15); border: 1px solid #ff003c; border-radius: 4px; display: flex; align-items: center; justify-content: space-between; gap: 0.5rem;">
        <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff003c" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <span class="vms-text-mono vms-text-2xs" style="color: #ff8899; font-weight: 600;">{{ saveAlertError }}</span>
        </div>
        <button class="vms-btn vms-btn-ghost vms-btn-xs" style="color: #fff; padding: 2px 6px;" @click="saveAlertError = null">FECHAR</button>
      </div>

      <!-- Rodapé Fixo -->
      <div class="vms-modal-footer" style="flex-shrink: 0;">
        <button v-if="step > 1" class="vms-btn vms-btn-secondary" :disabled="isTesting" @click="goToStep(step - 1)">ANTERIOR</button>
        <button v-if="step < totalSteps" class="vms-btn vms-btn-primary" :disabled="isTesting" @click="handleNext">PROXIMO</button>
        <button v-else class="vms-btn vms-btn-primary" :disabled="isTesting" @click="finish"><span v-if="isTesting">VALIDANDO FLUXOS...</span><span v-else>SALVAR</span></button>
      </div>
    </div>
  </div>
</template>
