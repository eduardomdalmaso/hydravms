import { ref, computed, onMounted } from 'vue'
import type { StoragePoolItem, UnallocatedDiskDevice, NewStoragePayload, StorageRole } from '../types/storagePool'
import { fetchStoragePools, fetchDetectedDisks, createRemoteStoragePool, deleteRemoteStoragePool, triggerRemoteDrain } from '../services/storageApi'

export function useStorageManager() {
  const pools = ref<StoragePoolItem[]>([])
  const unallocatedDisks = ref<UnallocatedDiskDevice[]>([])
  const filterRole = ref<StorageRole | 'ALL'>('ALL'), searchQuery = ref(''), isAddModalOpen = ref(false)
  const notification = ref<string | null>(null)

  const showNotification = (msg: string) => {
    notification.value = msg
    setTimeout(() => { if (notification.value === msg) notification.value = null }, 4000)
  }

  const loadPools = async () => {
    pools.value = await fetchStoragePools([])
    unallocatedDisks.value = await fetchDetectedDisks()
  }
  onMounted(loadPools)

  const totalCapacityGb = computed(() => pools.value.reduce((acc, p) => acc + p.totalGb, 0))
  const totalUsedGb = computed(() => pools.value.reduce((acc, p) => acc + p.usedGb, 0))
  const overallPercentage = computed(() => totalCapacityGb.value ? Math.round((totalUsedGb.value / totalCapacityGb.value) * 100) : 0)

  const hotPool = computed(() => pools.value.find(p => p.role === 'HOT_BUFFER'))
  const hotUsagePercent = computed(() => hotPool.value ? Math.round((hotPool.value.usedGb / hotPool.value.totalGb) * 100) : 0)

  const filteredPools = computed(() => {
    return pools.value.filter(p => {
      const matchRole = filterRole.value === 'ALL' || p.role === filterRole.value
      const query = searchQuery.value.toLowerCase()
      return matchRole && (!query || p.name.toLowerCase().includes(query) || p.nodeOrServer.toLowerCase().includes(query) || p.pathOrEndpoint.toLowerCase().includes(query))
    })
  })

  const addPool = async (payload: NewStoragePayload) => {
    isAddModalOpen.value = false
    await createRemoteStoragePool(payload)
    await loadPools()
    showNotification(payload.role === 'WARM_ARCHIVE'
      ? `[ALERTA] Gravação ativada em "${payload.name}". Requer disco NVMe de Buffer dedicado.`
      : `[STORAGE] Pool "${payload.name}" ativado com sucesso`)
  }

  const removePool = async (id: string) => {
    const pool = pools.value.find(p => p.id === id)
    if (!pool) return
    await deleteRemoteStoragePool(id)
    await loadPools()
    showNotification(`[STORAGE] Pool "${pool.name}" desanexado`)
  }

  const triggerSpilloverDrain = async () => {
    const msg = await triggerRemoteDrain()
    if (hotPool.value) hotPool.value.usedGb = Math.max(80, hotPool.value.usedGb - 200)
    showNotification(`[SPILLOVER] ${msg}`)
  }

  return {
    pools, unallocatedDisks, filterRole, filterType: filterRole, searchQuery, isAddModalOpen, notification,
    totalCapacityGb, totalUsedGb, overallPercentage, hotPool, hotUsagePercent, filteredPools,
    addPool, removePool, triggerSpilloverDrain, showNotification, reload: loadPools
  }
}
