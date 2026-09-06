import { ref, computed } from 'vue'
import type { StoragePoolItem, UnallocatedDiskDevice, NewStoragePayload, StorageRole } from '../types/storagePool'
import { initialStoragePools, detectedLocalDisks } from '../data/mockStoragePools'

export function useStorageManager() {
  const pools = ref<StoragePoolItem[]>([...initialStoragePools])
  const unallocatedDisks = ref<UnallocatedDiskDevice[]>([...detectedLocalDisks])
  const filterRole = ref<StorageRole | 'ALL'>('ALL')
  const searchQuery = ref('')
  const isAddModalOpen = ref(false)
  const notification = ref<string | null>(null)

  const showNotification = (msg: string) => {
    notification.value = msg
    setTimeout(() => { if (notification.value === msg) notification.value = null }, 4000)
  }

  const totalCapacityGb = computed(() => pools.value.reduce((acc, p) => acc + p.totalGb, 0))
  const totalUsedGb = computed(() => pools.value.reduce((acc, p) => acc + p.usedGb, 0))
  const overallPercentage = computed(() => totalCapacityGb.value ? Math.round((totalUsedGb.value / totalCapacityGb.value) * 100) : 0)

  const hotPool = computed(() => pools.value.find(p => p.role === 'HOT_BUFFER'))
  const hotUsagePercent = computed(() => hotPool.value ? Math.round((hotPool.value.usedGb / hotPool.value.totalGb) * 100) : 0)

  const filteredPools = computed(() => {
    return pools.value.filter(p => {
      const matchRole = filterRole.value === 'ALL' || p.role === filterRole.value
      const query = searchQuery.value.toLowerCase()
      const matchQuery = !query || p.name.toLowerCase().includes(query) || p.nodeOrServer.toLowerCase().includes(query) || p.pathOrEndpoint.toLowerCase().includes(query)
      return matchRole && matchQuery
    })
  })

  const addPool = (payload: NewStoragePayload) => {
    const newId = `sp_${Date.now()}`
    const pool: StoragePoolItem = {
      id: newId, ...payload, usedGb: 0, status: 'ONLINE', isSpilloverActive: payload.role === 'HOT_BUFFER'
    }
    pools.value.push(pool)
    unallocatedDisks.value = unallocatedDisks.value.filter(d => d.devicePath !== payload.pathOrEndpoint)
    isAddModalOpen.value = false
    if (payload.role === 'WARM_ARCHIVE') {
      showNotification(`[ALERTA DE ARQUITETURA] Gravacao ativada em "${payload.name}". Requer disco NVMe de Buffer dedicado; o disco do sistema opera apenas Snapshots.`)
    } else {
      showNotification(`[STORAGE] Pool "${payload.name}" ativado com sucesso`)
    }
  }

  const removePool = (id: string) => {
    const pool = pools.value.find(p => p.id === id)
    if (!pool) return
    if (pool.usedGb > 0 || pool.role === 'HOT_BUFFER' || pool.role === 'DATABASE') {
      showNotification(`[BLOQUEIO] O pool "${pool.name}" contem ${pool.usedGb} GB de gravacoes e nao pode ser desanexado`)
      return
    }
    pools.value = pools.value.filter(p => p.id !== id)
    showNotification(`[STORAGE] Pool "${pool.name}" desanexado do VMS`)
  }

  const triggerSpilloverDrain = () => {
    if (!hotPool.value) return
    hotPool.value.usedGb = Math.max(80, hotPool.value.usedGb - 200)
    showNotification('[SPILLOVER] Dreno manual do buffer NVMe para HDs concluido // 200 GB transferidos')
  }

  const triggerManualPurge = (id: string) => {
    const p = pools.value.find(item => item.id === id)
    if (!p) return
    p.usedGb = Math.max(10, Math.round(p.usedGb * 0.8))
    showNotification(`[STORAGEGUARD] Expurgo de blocos nao-pinados executado no pool "${p.name}"`)
  }

  return {
    pools, unallocatedDisks, filterRole, filterType: filterRole, searchQuery, isAddModalOpen, notification,
    totalCapacityGb, totalUsedGb, overallPercentage, hotPool, hotUsagePercent, filteredPools,
    addPool, removePool, triggerSpilloverDrain, triggerManualPurge, showNotification
  }
}
