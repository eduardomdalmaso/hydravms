import { ref, computed } from 'vue'
import { useAnalyticsAlerts } from './useAnalyticsAlerts'
import { useMarketplace } from './useMarketplace'
import type { AnalyticEventRecord } from '../types/marketplace'

export interface BoundingBox {
  id: string | number
  box: [number, number, number, number] // [x_pct, y_pct, w_pct, h_pct]
  color?: string
  category?: string
  confidence?: number
}

const activeCameraDetections = ref<Record<string, BoundingBox[]>>({})
const isRunning = ref(false)
let inferenceTimer: any = null
let lastEventTime = 0

export function useLiveDetections() {
  const { alerts } = useAnalyticsAlerts()
  const { events, instances } = useMarketplace()

  const setDetections = (cameraId: string, boxes: BoundingBox[]) => {
    activeCameraDetections.value[cameraId] = boxes
  }

  const getDetections = (cameraId: string) => {
    return computed(() => activeCameraDetections.value[cameraId] || [])
  }

  const pushAIEvent = (
    cameraId: string,
    cameraName: string,
    pluginId: string,
    pluginName: string,
    type: 'INTRUSÃO' | 'MULTIDÃO' | 'DETECÇÃO DE PESSOA' | 'VEÍCULO',
    confidence: number,
    bbox: [number, number, number, number]
  ) => {
    const timestamp = new Date().toLocaleTimeString()
    const isoTime = new Date().toISOString()
    const sevMap: Record<string, 'info' | 'warning' | 'critical'> = {
      'INTRUSÃO': 'critical',
      'MULTIDÃO': 'warning',
      'DETECÇÃO DE PESSOA': 'info',
      'VEÍCULO': 'info'
    }

    alerts.value.unshift({
      id: `ai_alert_${Date.now()}_${Math.random().toString(36).substr(2, 4)}`,
      camera_id: cameraId,
      camera_name: cameraName,
      event_type: `IA // ${type}`,
      severity: sevMap[type] || 'info',
      confidence,
      timestamp,
      snapshot_url: '',
      is_acknowledged: false
    })
    if (alerts.value.length > 50) alerts.value.pop()

    const newEvt: AnalyticEventRecord = {
      id: `evt-${Date.now()}-${Math.random().toString(36).substr(2, 4)}`,
      plugin_id: pluginId || 'object_detection_sota',
      plugin_name: pluginName || 'YOLO26m Object Detection',
      camera_id: cameraId,
      camera_name: cameraName,
      event_type: type,
      severity: sevMap[type] || 'info',
      object_label: type === 'VEÍCULO' ? 'veículo' : 'pessoa',
      confidence,
      timestamp: isoTime,
      details: `${type} detectada na câmera ${cameraName} com confiança de ${(confidence * 100).toFixed(0)}%`,
      snapshot_url: '',
      bbox,
      raw_payload: { model: 'yolo26m', bbox, confidence }
    }
    events.value.unshift(newEvt)
    if (events.value.length > 100) events.value.pop()
  }

  const fetchRealCameraInference = async (camId: string, camName: string, conf: number) => {
    try {
      const res = await fetch('http://localhost:8081/api/v1/inference/predict', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: 'yolo26m.pt',
          source: camId,
          classes: 'person',
          conf: conf || 0.25,
          device: '0'
        })
      })
      if (!res.ok) return
      const data = await res.json()
      if (data && Array.isArray(data.detections)) {
        const boxes: BoundingBox[] = data.detections
          .filter((d: any) => d.class_name === 'person' || d.label === 'person' || d.label === 'pessoa')
          .map((d: any) => ({
            id: d.id || d.track_id || Math.random(),
            box: d.box as [number, number, number, number],
            color: '#ff5e3a',
            category: 'person',
            confidence: d.confidence || d.conf
          }))
        activeCameraDetections.value[camId] = boxes

        // Dispatch real AI event on person detection (with cooldown)
        const persons = boxes.filter(b => b.category === 'person')
        const now = Date.now()
        if (persons.length > 0 && now - lastEventTime > 4000) {
          lastEventTime = now
          const type = persons.length >= 2 ? 'MULTIDÃO' : (persons[0].confidence && persons[0].confidence > 0.85 ? 'INTRUSÃO' : 'DETECÇÃO DE PESSOA')
          pushAIEvent(camId, camName, 'object_detection_sota', 'YOLO26m Object Detection', type, persons[0].confidence || 0.85, persons[0].box)
        }
      }
    } catch {}
  }

  const startLiveAnalyticsPipeline = () => {
    if (isRunning.value) return
    isRunning.value = true

    // Real inference loop at regular interval (every 800ms)
    const runTick = async () => {
      const activeInst = instances.value.filter(i => i.is_active)
      if (activeInst.length > 0) {
        for (const inst of activeInst) {
          await fetchRealCameraInference(inst.camera_id || 'cam_01', inst.camera_name || 'C113', inst.confidence_threshold || 0.25)
        }
      } else {
        await fetchRealCameraInference('cam_01', 'C113', 0.25)
      }
    }

    runTick()
    inferenceTimer = setInterval(runTick, 800)
  }

  const stopLiveAnalyticsPipeline = () => {
    if (inferenceTimer) clearInterval(inferenceTimer)
    isRunning.value = false
  }

  return {
    activeCameraDetections,
    setDetections,
    getDetections,
    pushAIEvent,
    startLiveAnalyticsPipeline,
    stopLiveAnalyticsPipeline
  }
}
