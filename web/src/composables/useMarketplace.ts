import { ref, computed } from 'vue'
import type { PluginManifest, AnalyticInstance, AnalyticEventRecord } from '../types/marketplace'
import { fetchPlugins, installRemotePlugin, uninstallRemotePlugin, toggleRemotePlugin } from '../services/api'

// Module-level singletons so all views and sidebar share the exact same state
const plugins = ref<PluginManifest[]>([])
const instances = ref<AnalyticInstance[]>([])
const events = ref<AnalyticEventRecord[]>([])
const toast = ref<string | null>(null)
const isLoaded = ref(false)

export function useMarketplace() {
  const searchQuery = ref(''), statusFilter = ref<'ALL' | 'installed' | 'available'>('ALL')
  const showToast = (msg: string) => { toast.value = msg; setTimeout(() => { toast.value = null }, 3500) }
  const installedPlugins = computed(() => plugins.value.filter(p => p.is_installed))

  const loadPlugins = async () => {
    const list = await fetchPlugins()
    if (list.length > 0) {
      plugins.value = list
    }
    isLoaded.value = true
  }

  if (!isLoaded.value && typeof window !== 'undefined') {
    loadPlugins()
  }

  const filteredPlugins = computed(() => plugins.value.filter(p => {
    const mStat = statusFilter.value === 'ALL' || (statusFilter.value === 'installed' ? p.is_installed : !p.is_installed)
    const mQ = !searchQuery.value || p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || (p.description && p.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
    return mStat && mQ
  }))

  const installPlugin = async (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    p.status = 'updating'
    showToast(`[DOWNLOAD] Baixando pacote ${p.name}...`)
    const ok = await installRemotePlugin(id)
    if (ok) {
      p.is_installed = true
      p.status = 'running'
      showToast(`[INSTALADO] Módulo ${p.name} ativado no sistema!`)
    } else {
      p.status = 'available'
      showToast(`[ERRO] Falha ao instalar analítico ${p.name}`)
    }
    await loadPlugins()
  }

  const uninstallPlugin = async (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    const ok = await uninstallRemotePlugin(id)
    if (ok) {
      p.is_installed = false
      p.status = 'available'
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

  const deleteInstance = (id: string) => {
    instances.value = instances.value.filter(x => x.id !== id)
    showToast(`[EXCLUÍDO] Instância removida`)
  }

  const getPluginById = (id: string) => plugins.value.find(p => p.id === id)
  const getInstancesByPlugin = (pluginId: string) => instances.value.filter(i => i.plugin_id === pluginId)
  const getEventsByPlugin = (pluginId: string) => events.value.filter(e => e.plugin_id === pluginId)

  return {
    plugins, instances, events, searchQuery, statusFilter, toast,
    installedPlugins, filteredPlugins, showToast, loadPlugins, installPlugin, uninstallPlugin, togglePlugin,
    createInstance, toggleInstance, deleteInstance, getPluginById, getInstancesByPlugin, getEventsByPlugin
  }
}
