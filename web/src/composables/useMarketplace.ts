import { ref, computed } from 'vue'
import type { PluginManifest, AnalyticInstance, AnalyticEventRecord } from '../types/marketplace'
import { mockMarketplacePlugins } from '../data/mockMarketplacePlugins'
import { mockMarketplaceInstances } from '../data/mockMarketplaceInstances'
import { mockMarketplaceEvents } from '../data/mockMarketplaceEvents'

// Module-level singletons so all views and sidebar share the exact same state
const plugins = ref<PluginManifest[]>([...mockMarketplacePlugins])
const instances = ref<AnalyticInstance[]>([...mockMarketplaceInstances])
const events = ref<AnalyticEventRecord[]>([...mockMarketplaceEvents])
const toast = ref<string | null>(null)

export function useMarketplace() {
  const searchQuery = ref(''), selectedCategory = ref('ALL'), statusFilter = ref<'ALL' | 'installed' | 'available'>('ALL')
  const showToast = (msg: string) => { toast.value = msg; setTimeout(() => { toast.value = null }, 3500) }
  const installedPlugins = computed(() => plugins.value.filter(p => p.is_installed))

  const filteredPlugins = computed(() => plugins.value.filter(p => {
    const mCat = selectedCategory.value === 'ALL' || p.category === selectedCategory.value
    const mStat = statusFilter.value === 'ALL' || (statusFilter.value === 'installed' ? p.is_installed : !p.is_installed)
    const mQ = !searchQuery.value || p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    return mCat && mStat && mQ
  }))

  const installPlugin = (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    p.status = 'updating'; showToast(`[DOWNLOAD] Baixando pacote ${p.name}...`)
    setTimeout(() => { p.is_installed = true; p.status = 'running'; showToast(`[INSTALADO] Módulo ${p.name} adicionado ao menu lateral!`) }, 800)
  }

  const uninstallPlugin = (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    p.is_installed = false; p.status = 'available'; p.instances_count = 0
    instances.value = instances.value.filter(i => i.plugin_id !== id)
    showToast(`[DESINSTALADO] Módulo ${p.name} removido do menu lateral`)
  }

  const togglePlugin = (id: string) => {
    const p = plugins.value.find(x => x.id === id)
    if (!p) return
    p.status = p.status === 'running' ? 'stopped' : 'running'
    showToast(`[STATUS] ${p.name} // ${p.status === 'running' ? 'ATIVO' : 'PAUSADO'}`)
  }

  const createInstance = (inst: AnalyticInstance) => {
    instances.value.unshift(inst)
    const p = plugins.value.find(x => x.id === inst.plugin_id)
    if (p) p.instances_count++
    showToast(`[CRIADO] Instância ${inst.name} salva com sucesso!`)
  }

  const toggleInstance = (id: string) => {
    const inst = instances.value.find(x => x.id === id)
    if (!inst) return
    inst.is_active = !inst.is_active
    showToast(`[INSTÂNCIA] ${inst.name} // ${inst.is_active ? 'ATIVO' : 'PAUSADO'}`)
  }

  const deleteInstance = (id: string) => {
    const inst = instances.value.find(x => x.id === id)
    if (!inst) return
    const p = plugins.value.find(x => x.id === inst.plugin_id)
    if (p && p.instances_count > 0) p.instances_count--
    instances.value = instances.value.filter(x => x.id !== id)
    showToast(`[EXCLUÍDO] Instância ${inst.name} removida`)
  }

  const getPluginById = (id: string) => plugins.value.find(p => p.id === id)
  const getInstancesByPlugin = (pluginId: string) => instances.value.filter(i => i.plugin_id === pluginId)
  const getEventsByPlugin = (pluginId: string) => events.value.filter(e => e.plugin_id === pluginId)

  return {
    plugins, instances, events, searchQuery, selectedCategory, statusFilter, toast,
    installedPlugins, filteredPlugins, showToast, installPlugin, uninstallPlugin, togglePlugin,
    createInstance, toggleInstance, deleteInstance, getPluginById, getInstancesByPlugin, getEventsByPlugin
  }
}
