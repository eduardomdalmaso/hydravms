// HydraVMS - Camera Recording Profiles & Segment Query API (< 100 lines)
import type { RecordingProfile } from '../types/recordingSchedule'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8083'

export async function fetchRemoteRecordingProfiles(cameraId: string): Promise<RecordingProfile[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cameras/${encodeURIComponent(cameraId)}/recording-profiles`, {
      headers: { 'Accept': 'application/json' },
      signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) return []
    const data = await res.json()
    if (!Array.isArray(data.profiles)) return []
    return data.profiles.map((p: any, idx: number) => ({
      id: p.id || `REC_0${idx + 1}`,
      name: p.name || 'Perfil de Gravacao',
      mode: p.mode || 'continuous',
      isActive: p.is_active !== false,
      preBuffer: p.pre_buffer_s || 5,
      postBuffer: p.post_buffer_s || 15,
      schedule: typeof p.schedule_json === 'string' ? JSON.parse(p.schedule_json) : (p.schedule_json || [])
    }))
  } catch {
    return []
  }
}

export async function saveRemoteRecordingProfile(cameraId: string, profile: RecordingProfile): Promise<boolean> {
  try {
    const payload = {
      id: profile.id.includes('-') ? profile.id : undefined,
      name: profile.name,
      mode: profile.mode,
      pre_buffer_s: profile.preBuffer,
      post_buffer_s: profile.postBuffer,
      is_active: profile.isActive,
      schedule_json: JSON.stringify(profile.schedule)
    }
    const res = await fetch(`${API_BASE}/api/v1/cameras/${encodeURIComponent(cameraId)}/recording-profiles`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch {
    return false
  }
}

export async function deleteRemoteRecordingProfile(cameraId: string, profileId: string): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/cameras/${encodeURIComponent(cameraId)}/recording-profiles/${encodeURIComponent(profileId)}`, {
      method: 'DELETE',
      signal: AbortSignal.timeout(3000)
    })
    return res.ok
  } catch {
    return false
  }
}

export interface RemoteRecordingSegment {
  id: string
  start_time: string
  end_time: string
  recording_mode: string
  duration_seconds: number
}

export async function fetchRemoteRecordings(cameraId: string, start?: string, end?: string): Promise<RemoteRecordingSegment[]> {
  try {
    const q = new URLSearchParams()
    if (start) q.set('start', start)
    if (end) q.set('end', end)
    const res = await fetch(`${API_BASE}/api/v1/cameras/${encodeURIComponent(cameraId)}/recordings?${q.toString()}`, {
      headers: { 'Accept': 'application/json' },
      signal: AbortSignal.timeout(3000)
    })
    if (!res.ok) return []
    const data = await res.json()
    return Array.isArray(data.recordings) ? data.recordings : []
  } catch {
    return []
  }
}
