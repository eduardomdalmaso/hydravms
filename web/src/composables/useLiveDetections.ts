import { ref, computed } from 'vue'
import { useAnalyticsAlerts } from './useAnalyticsAlerts'
import { useMarketplace } from './useMarketplace'
import { getGloballyActiveAnalyticInstances } from './useDesktopAnalytics'
import { createRemoteEvent } from '../services/api'
import type { AnalyticEventRecord, AnalyticInstance, Point2D } from '../types/marketplace'

export interface BoundingBox {
  id: string | number
  box: [number, number, number, number] // [x_pct, y_pct, w_pct, h_pct]
  color?: string
  category?: string
  confidence?: number
}

const activeCameraDetections = ref<Record<string, BoundingBox[]>>({})
const isRunning = ref(false)
let pollTimer: any = null
let lastEventTime = 0

function isPointInPolygon(px: number, py: number, vs: Point2D[]): boolean {
  if (!vs || vs.length < 3) return false
  const poly = vs.map(p => ({
    x: p.x <= 1.0 ? p.x * 100 : p.x,
    y: p.y <= 1.0 ? p.y * 100 : p.y
  }))
  let inside = false
  for (let i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    const xi = poly[i].x, yi = poly[i].y
    const xj = poly[j].x, yj = poly[j].y
    const intersect = ((yi > py) !== (yj > py)) && (px < ((xj - xi) * (py - yi)) / (yj - yi) + xi)
    if (intersect) inside = !inside
  }
  return inside
}

function getActiveInstances(camId: string): AnalyticInstance[] {
  const liveActive = getGloballyActiveAnalyticInstances()
  if (liveActive.length > 0) return liveActive.filter(i => i.is_active && (!i.camera_id || i.camera_id === camId || i.camera_id === 'cam_01' || i.camera_id === 'c113'))
  try {
    if (typeof localStorage !== 'undefined') {
      const root = JSON.parse(localStorage.getItem('hydra_vms_analytic_instances') || '[]')
      const flds = JSON.parse(localStorage.getItem('hydra_vms_analytic_folders') || '[]')
      const active: AnalyticInstance[] = Array.isArray(root) ? root.filter((i: any) => i.is_active) : []
      if (Array.isArray(flds)) flds.forEach((f: any) => { if (Array.isArray(f.instances)) active.push(...f.instances.filter((i: any) => i.is_active)) })
      return active.filter(i => i.is_active && (!i.camera_id || i.camera_id === camId || i.camera_id === 'cam_01' || i.camera_id === 'c113'))
    }
  } catch {}
  return []
}

export function useLiveDetections() {
  const { alerts } = useAnalyticsAlerts()
  const { events } = useMarketplace()

  const setDetections = (cameraId: string, boxes: BoundingBox[]) => { activeCameraDetections.value[cameraId] = boxes }
  const getDetections = (cameraId: string) => computed(() => activeCameraDetections.value[cameraId] || [])

  const pushAIEvent = (
    cameraId: string, cameraName: string, pluginId: string, pluginName: string,
    type: 'INTRUSÃO', confidence: number, bbox: [number, number, number, number]
  ) => {
    const timestamp = new Date().toLocaleTimeString(), isoTime = new Date().toISOString()
    const snapUrl = `http://localhost:8080/api/v1/streams/${cameraId || 'cam_01'}/snapshot`

    alerts.value.unshift({
      id: `ai_alert_${Date.now()}_${Math.random().toString(36).substr(2, 4)}`,
      camera_id: cameraId, camera_name: cameraName, event_type: `IA // ${type}`,
      severity: 'critical', confidence, timestamp, snapshot_url: snapUrl, bbox, is_acknowledged: false
    })
    if (alerts.value.length > 50) alerts.value.pop()

    const newEvt: AnalyticEventRecord = {
      id: `evt-${Date.now()}-${Math.random().toString(36).substr(2, 4)}`,
      plugin_id: pluginId || 'object_detection_sota', plugin_name: pluginName || 'YOLO26m Object Detection',
      camera_id: cameraId, camera_name: cameraName, event_type: type, severity: 'critical',
      object_label: 'pessoa', confidence, timestamp: isoTime,
      details: `Intrusão de pessoa detectada na câmera ${cameraName} com confiança de ${(confidence * 100).toFixed(0)}%`,
      snapshot_url: snapUrl, bbox, raw_payload: { model: 'yolo26m', bbox, confidence }
    }
    events.value.unshift(newEvt)
    if (events.value.length > 100) events.value.pop()

    createRemoteEvent({
      camera_id: cameraId, event_type: 'ai.detection.zone_intrusion', severity: 'critical',
      object_class: 'pessoa', confidence,
      bbox_normalized: { x_center: (bbox[0] + bbox[2] / 2) / 100, y_center: (bbox[1] + bbox[3] / 2) / 100, width: bbox[2] / 100, height: bbox[3] / 100 },
      snapshot_url: snapUrl, notes: `Intrusão de pessoa na câmera ${cameraName} (confiança: ${(confidence * 100).toFixed(0)}%)`
    }).catch(() => {})
  }

  const fetchStreamDetections = async (camId: string, camName: string) => {
    try {
      const res = await fetch(`http://localhost:8081/api/v1/inference/detections?camera_id=${encodeURIComponent(camId)}`)
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

        // Strict pause check: if no active instance exists, do not dispatch any event
        const activeInsts = getActiveInstances(camId)
        if (!activeInsts || activeInsts.length === 0 || boxes.length === 0) return

        const now = Date.now()
        if (now - lastEventTime < 4000) return

        for (const inst of activeInsts) {
          if (!inst.is_active) continue

          const threshold = inst.confidence_threshold || 0.60
          const zone = inst.zones?.[0]
          const poly = zone?.polygon || inst.specific_params?.polygon

          for (const b of boxes) {
            if ((b.confidence || 0) < threshold) continue

            // Bottom-center foot coordinate (percentage 0..100)
            const feetX = b.box[0] + b.box[2] / 2
            const feetY = b.box[1] + b.box[3]

            const isInside = poly && poly.length >= 3 ? isPointInPolygon(feetX, feetY, poly) : (inst.roi_mode === 'full_frame')
            if (isInside) {
              lastEventTime = now
              pushAIEvent(camId, camName, inst.plugin_id, inst.plugin_name, 'INTRUSÃO', b.confidence || 0.85, b.box)
              return
            }
          }
        }
      }
    } catch {}
  }

  const startLiveAnalyticsPipeline = () => {
    if (isRunning.value) return
    isRunning.value = true
    const poll = async () => { await fetchStreamDetections('cam_01', 'C113') }
    poll()
    pollTimer = setInterval(poll, 120)
  }

  const stopLiveAnalyticsPipeline = () => {
    if (pollTimer) clearInterval(pollTimer)
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
