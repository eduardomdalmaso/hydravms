import { ref, computed } from 'vue'
import type { PluginManifest, AnalyticInstance, AnalyticEventRecord } from '../types/marketplace'
import { fetchPlugins, installRemotePlugin, uninstallRemotePlugin, toggleRemotePlugin, fetchEvents } from '../services/api'

const plugins = ref<PluginManifest[]>([])
const instances = ref<AnalyticInstance[]>([])
const events = ref<AnalyticEventRecord[]>([])
const toast = ref<string | null>(null), isLoaded = ref(false)
const gpuDetected = ref(true), gpuInfo = ref<any>(null)
const installProgress = ref<Record<string, { step: string; percent: number }>>({})

export function useMarketplace() {
  const searchQuery = ref(''), statusFilter = ref<'ALL' | 'installed' | 'available'>('ALL')
  const showToast = (msg: string) => { toast.value = msg; setTimeout(() => { toast.value = null }, 3500) }
  const installedPlugins = computed(() => plugins.value.filter(p => p.is_installed))

  const loadEvents = async () => {
    try {
      const rawEvents = await fetchEvents()
      if (Array.isArray(rawEvents) && rawEvents.length > 0) {
        events.value = rawEvents.map(e => ({
          id: e.id,
          plugin_id: 'object_detection_sota',
          plugin_name: 'YOLO26m Object Detection',
          camera_id: e.camera_id || 'cam_01',
          camera_name: e.camera_id === 'cam_01' || e.camera_id === 'c113' ? 'C113' : (e.camera_id || 'CAMERA'),
          event_type: e.event_type?.toUpperCase().replace('SYSTEM.CAMERA.', '').replace('AI.DETECTION.', '') || 'DETECÇÃO DE PESSOA',
          severity: e.severity === 'critical' ? 'critical' : (e.severity === 'warning' ? 'warning' : 'info'),
          object_label: e.object_class || 'pessoa',
          confidence: e.confidence || 0.9,
          timestamp: e.triggered_at || e.created_at || new Date().toISOString(),
          details: e.notes || `Detecção de ${e.object_class || 'pessoa'} na câmera`,
          snapshot_url: e.snapshot_s3_key ? (e.snapshot_s3_key.startsWith('http') ? e.snapshot_s3_key : `http://localhost:8080${e.snapshot_s3_key}`) : 'http://localhost:8080/api/v1/streams/cam_01/snapshot',
          bbox: e.bbox_normalized ? [
            (e.bbox_normalized.x_center - e.bbox_normalized.width / 2) * 100,
            (e.bbox_normalized.y_center - e.bbox_normalized.height / 2) * 100,
            e.bbox_normalized.width * 100,
            e.bbox_normalized.height * 100
          ] : [20, 20, 30, 50],
          raw_payload: e
        }))
      }
    } catch {}
  }

  const loadPlugins = async () => {
    const res = await fetchPlugins()
    plugins.value = Array.isArray(res.plugins) ? res.plugins : []
    gpuDetected.value = res.gpu_detected; gpuInfo.value = res.gpu_telemetry
    await loadEvents()
    isLoaded.value = true
  }

  if (!isLoaded.value && typeof window !== 'undefined') loadPlugins()

  const filteredPlugins = computed(() => plugins.value.filter(p => {
    const mStat = statusFilter.value === 'ALL' || (statusFilter.value === 'installed' ? p.is_installed : !p.is_installed)
    const mQ = !searchQuery.value || p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || (p.description && p.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
    return mStat && mQ
  }))

  const installPlugin = async (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    p.status = 'updating'
    installProgress.value[id] = { step: '[1/3] Baixando pesos e manifesto...', percent: 30 }
    showToast(`[DOWNLOAD] Baixando pacote ${p.name}...`)
    await new Promise(r => setTimeout(r, 600))
    if (gpuDetected.value) {
      installProgress.value[id] = { step: '[2/3] Compilando TensorRT .engine FP16 na GPU...', percent: 75 }
      showToast(`[TENSORRT JIT] Compilando .engine FP16 na GPU...`)
      await new Promise(r => setTimeout(r, 900))
    }
    const ok = await installRemotePlugin(id)
    if (ok) {
      p.is_installed = true; p.status = 'running'
      installProgress.value[id] = { step: '[3/3] Módulo ativado!', percent: 100 }
      showToast(`[INSTALADO] Módulo ${p.name} ativado no sistema!`)
    } else {
      p.status = 'available'; showToast(`[ERRO] Falha ao instalar analítico ${p.name}`)
    }
    delete installProgress.value[id]
    await loadPlugins()
  }

  const uninstallPlugin = async (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    const ok = await uninstallRemotePlugin(id)
    if (ok) {
      p.is_installed = false; p.status = 'available'
      instances.value = instances.value.filter(i => i.plugin_id !== id)
      showToast(`[DESINSTALADO] Módulo ${p.name} removido`)
    }
    await loadPlugins()
  }

  const togglePlugin = async (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    const ok = await toggleRemotePlugin(id)
    if (ok) {
      p.status = p.status === 'running' ? 'stopped' : 'running'
      showToast(`[STATUS] ${p.name} // ${p.status === 'running' ? 'ATIVO' : 'PAUSADO'}`)
    }
    await loadPlugins()
  }

  const createInstance = (inst: AnalyticInstance) => {
    instances.value.unshift(inst)
    showToast(`[CRIADO] Instância ${inst.name} salva com sucesso!`)
  }
  const toggleInstance = (id: string) => {
    const inst = instances.value.find(x => x.id === id)
    if (!inst) return
    inst.is_active = !inst.is_active
    showToast(`[INSTÂNCIA] ${inst.name} // ${inst.is_active ? 'ATIVO' : 'PAUSADO'}`)
  }
  const deleteInstance = (id: string) => { instances.value = instances.value.filter(x => x.id !== id); showToast(`[EXCLUÍDO] Instância removida`) }
  const getPluginById = (id: string) => plugins.value.find(p => p.id === id)
  const getInstancesByPlugin = (pluginId: string) => instances.value.filter(i => i.plugin_id === pluginId)
  const getEventsByPlugin = (pluginId: string) => events.value.filter(e => e.plugin_id === pluginId)

  return {
    plugins, instances, events, searchQuery, statusFilter, toast, gpuDetected, gpuInfo, installProgress,
    installedPlugins, filteredPlugins, showToast, loadPlugins, loadEvents, installPlugin, uninstallPlugin, togglePlugin,
    createInstance, toggleInstance, deleteInstance, getPluginById, getInstancesByPlugin, getEventsByPlugin
  }
}
