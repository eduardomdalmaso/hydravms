import { ref, onMounted, onUnmounted } from 'vue'

export interface MonitorDispatchEvent {
  type: 'DISPATCH_LAYOUT' | 'PING_MONITORS'
  layoutId?: string
  targetMonitor: number
}

const CHANNEL_NAME = 'hydravms_multimonitor_bus'

export function useMultiMonitor() {
  const isPopout = ref(false)
  const activeMonitorNumber = ref(1)
  let channel: BroadcastChannel | null = null

  const checkPopoutMode = () => {
    if (typeof window === 'undefined') return
    const hash = window.location.hash || ''
    const search = window.location.search || ''
    const isPop = hash.includes('popout=true') || search.includes('popout=true')
    isPopout.value = isPop
    const match = hash.match(/monitor=(\d+)/) || search.match(/monitor=(\d+)/)
    if (match && match[1]) {
      activeMonitorNumber.value = parseInt(match[1], 10)
    } else {
      activeMonitorNumber.value = isPop ? 2 : 1
    }
  }

  const openInPopout = (layoutId?: string, monitorNum: number = 2) => {
    if (typeof window === 'undefined') return
    const url = `${window.location.origin}${window.location.pathname}#vms?popout=true&monitor=${monitorNum}&layout=${layoutId || ''}`
    const win = window.open(url, `vms_screen_${monitorNum}`, 'width=1366,height=768,menubar=no,toolbar=no,location=no,status=no')
    if (win) win.focus()
  }

  const dispatchToMonitor = (layoutId: string, monitorNum: number) => {
    if (!channel && typeof BroadcastChannel !== 'undefined') {
      channel = new BroadcastChannel(CHANNEL_NAME)
    }
    channel?.postMessage({
      type: 'DISPATCH_LAYOUT',
      layoutId,
      targetMonitor: monitorNum
    } as MonitorDispatchEvent)
  }

  const initChannelListener = (onReceiveLayout: (layoutId: string) => void) => {
    if (typeof BroadcastChannel === 'undefined') return
    channel = new BroadcastChannel(CHANNEL_NAME)
    channel.onmessage = (e: MessageEvent<MonitorDispatchEvent>) => {
      const data = e.data
      if (data && data.type === 'DISPATCH_LAYOUT') {
        if (data.targetMonitor === activeMonitorNumber.value || (isPopout.value && data.targetMonitor === 0)) {
          if (data.layoutId) onReceiveLayout(data.layoutId)
        }
      }
    }
  }

  onMounted(() => {
    checkPopoutMode()
  })

  onUnmounted(() => {
    if (channel) channel.close()
  })

  return {
    isPopout,
    activeMonitorNumber,
    openInPopout,
    dispatchToMonitor,
    initChannelListener
  }
}
