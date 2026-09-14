// HydraVMS - Technical ID & Numeric Formatter (< 50 lines)

export function formatNumericId(index: number, prefix: string = ''): string {
  const numStr = String(index).padStart(2, '0')
  return prefix ? `${prefix}${numStr}` : `#${numStr}`
}

export function formatChannelId(indexOrId: number | string): string {
  if (typeof indexOrId === 'number') return `#${String(indexOrId).padStart(2, '0')}`
  const numMatch = indexOrId.match(/\d+/)
  if (numMatch) return `#${String(parseInt(numMatch[0], 10)).padStart(2, '0')}`
  return '#01'
}

export function formatProfileId(indexOrId: number | string): string {
  if (typeof indexOrId === 'number') return `#${String(indexOrId).padStart(2, '0')}`
  const numMatch = indexOrId.match(/\d+/)
  if (numMatch) return `#${String(parseInt(numMatch[0], 10)).padStart(2, '0')}`
  return '#01'
}

export function formatUserId(indexOrId: number | string): string {
  if (typeof indexOrId === 'number') return `#${String(indexOrId).padStart(2, '0')}`
  const numMatch = indexOrId.match(/\d+/)
  if (numMatch) return `#${String(parseInt(numMatch[0], 10)).padStart(2, '0')}`
  return '#01'
}

export function formatPluginId(indexOrId: number | string): string {
  if (typeof indexOrId === 'number') return `#${String(indexOrId).padStart(2, '0')}`
  const numMatch = indexOrId.match(/\d+/)
  if (numMatch) return `#${String(parseInt(numMatch[0], 10)).padStart(2, '0')}`
  return '#01'
}

export function formatPoolId(indexOrId: number | string): string {
  if (typeof indexOrId === 'number') return `#${String(indexOrId).padStart(2, '0')}`
  const numMatch = indexOrId.match(/\d+/)
  if (numMatch) return `#${String(parseInt(numMatch[0], 10)).padStart(2, '0')}`
  return '#01'
}

