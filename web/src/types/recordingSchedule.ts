export type RecordingMode = 'continuous' | 'motion' | 'ai_event'

export interface RecordingProfile {
  id: string
  name: string
  mode: RecordingMode
  isActive: boolean
  schedule: boolean[][] // 7 days (0=DOM..6=SAB) x 24 hours (0..23)
  preBuffer: number
  postBuffer: number
}

export function createDefaultSchedule(filled = true): boolean[][] {
  return Array.from({ length: 7 }, () => Array(24).fill(filled))
}
