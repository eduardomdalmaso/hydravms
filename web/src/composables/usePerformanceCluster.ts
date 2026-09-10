import { ref, onMounted, onUnmounted } from 'vue'
import type { ServerNodeItem } from '../types/performanceCluster'
import { initialServerNodes } from '../data/mockPerformanceNodes'
import { fetchLiveClusterNodes, createRemoteClusterNode, deleteRemoteClusterNode } from '../services/adminApi'

export function usePerformanceCluster() {
  const servers = ref<ServerNodeItem[]>([...initialServerNodes])
  const isRefreshing = ref(false), countdown = ref(10)
  let timerInterval: ReturnType<typeof setInterval> | null = null

  const loadNodes = async () => {
    const live = await fetchLiveClusterNodes()
    servers.value = live
  }

  const refreshMetrics = () => {
    isRefreshing.value = true
    setTimeout(() => {
      servers.value.forEach(s => {
        s.cpuPercent = Math.min(95, Math.max(15, Math.round(s.cpuPercent + (Math.random() * 8 - 4))))
        s.gpus.forEach(g => {
          g.computePercent = Math.min(98, Math.max(20, Math.round(g.computePercent + (Math.random() * 6 - 3))))
          g.vramUsedGb = parseFloat(Math.min(g.vramTotalGb, Math.max(4, g.vramUsedGb + (Math.random() * 1.6 - 0.8))).toFixed(1))
        })
      })
      isRefreshing.value = false; countdown.value = 10
    }, 300)
  }

  const addNode = async (data: { node_name: string; node_role: string; ip_address: string; http_port: number; grpc_port: number; webrtc_port: number; gpu_device_info?: string }) => {
    const ok = await createRemoteClusterNode(data)
    if (ok) await loadNodes()
    return ok
  }

  const removeNode = async (id: string) => {
    const ok = await deleteRemoteClusterNode(id)
    if (ok) servers.value = servers.value.filter(s => s.id !== id)
    return ok
  }

  onMounted(() => {
    loadNodes()
    timerInterval = setInterval(() => { if (countdown.value > 1) countdown.value--; else refreshMetrics() }, 1000)
  })

  onUnmounted(() => { if (timerInterval) clearInterval(timerInterval) })

  return { servers, isRefreshing, countdown, refreshMetrics, loadNodes, addNode, removeNode }
}
