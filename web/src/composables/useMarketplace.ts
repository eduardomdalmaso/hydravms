import { ref, computed } from 'vue'
import type { PluginManifest, AnalyticInstance, AnalyticEventRecord } from '../types/marketplace'
import { fetchPlugins, installRemotePlugin, uninstallRemotePlugin, toggleRemotePlugin } from '../services/api'

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

  const loadPlugins = async () => {
    const res = await fetchPlugins()
    plugins.value = Array.isArray(res.plugins) ? res.plugins : []
    gpuDetected.value = res.gpu_detected; gpuInfo.value = res.gpu_telemetry
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

  const createInstance = (inst: AnalyticInstance) => { instances.value.unshift(inst); showToast(`[CRIADO] Instância ${inst.name} salva com sucesso!`) }
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
    installedPlugins, filteredPlugins, showToast, loadPlugins, installPlugin, uninstallPlugin, togglePlugin,
    createInstance, toggleInstance, deleteInstance, getPluginById, getInstancesByPlugin, getEventsByPlugin
  }
}
