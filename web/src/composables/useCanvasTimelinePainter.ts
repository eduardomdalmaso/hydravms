export function drawCanvasTimeline(
  ctx: CanvasRenderingContext2D,
  w: number,
  h: number,
  currentTime: number,
  zoomMinutes: number,
  isExportMode: boolean,
  exportStart?: number,
  exportEnd?: number,
  recordedRanges: { start: number; end: number }[] = []
) {
  // Dark Background (Represents GAPs / No Recording)
  ctx.fillStyle = '#14171c'; ctx.fillRect(0, 0, w, h)
  const windowSpan = zoomMinutes * 60000
  const windowStart = currentTime - windowSpan / 2

  // Render Real Recording Segments from Database (Solid Orange)
  ctx.fillStyle = isExportMode ? 'rgba(255, 94, 58, 0.18)' : '#ff5e3a'
  for (const r of recordedRanges) {
    if (r.end < windowStart || r.start > windowStart + windowSpan) continue
    const clampStart = Math.max(r.start, windowStart)
    const clampEnd = Math.min(r.end, windowStart + windowSpan)
    const x = ((clampStart - windowStart) / windowSpan) * w
    const segmentWidth = ((clampEnd - clampStart) / windowSpan) * w
    if (segmentWidth > 0) {
      ctx.fillRect(x, 12, segmentWidth, h - 30)
    }
  }

  // Export Mode Orange Dual Selection Region & Bars
  if (isExportMode && exportStart && exportEnd) {
    const x1 = ((exportStart - windowStart) / windowSpan) * w
    const x2 = ((exportEnd - windowStart) / windowSpan) * w
    const minX = Math.min(x1, x2), spanW = Math.abs(x2 - x1)

    // Orange translucent selection highlight
    ctx.fillStyle = 'rgba(255, 94, 58, 0.32)'
    ctx.fillRect(minX, 0, spanW, h)

    // Start Bar (Orange)
    ctx.fillStyle = '#ff5e3a'; ctx.strokeStyle = '#ff5e3a'; ctx.lineWidth = 3
    ctx.beginPath(); ctx.moveTo(x1, 0); ctx.lineTo(x1, h); ctx.stroke()
    ctx.fillRect(x1 - 18, 0, 36, 12); ctx.fillStyle = '#000'; ctx.font = '700 8px Roboto'
    ctx.fillText('◄ INÍCIO', x1 - 15, 9)

    // End Bar (Orange)
    ctx.fillStyle = '#ff5e3a'; ctx.strokeStyle = '#ff5e3a'; ctx.lineWidth = 3
    ctx.beginPath(); ctx.moveTo(x2, 0); ctx.lineTo(x2, h); ctx.stroke()
    ctx.fillRect(x2 - 14, 0, 28, 12); ctx.fillStyle = '#000'; ctx.font = '700 8px Roboto'
    ctx.fillText('FIM ►', x2 - 11, 9)
  }

  // Ticks and numbers stay 100% full-contrast crisp white in Roboto without dimming
  ctx.fillStyle = '#ffffff'; ctx.font = '500 10px Roboto, sans-serif'
  const tickCount = zoomMinutes >= 720 ? 12 : zoomMinutes <= 3 ? 5 : 6
  let lastStr = ''
  for (let t = 0; t <= tickCount; t++) {
    const tickX = (t / tickCount) * w
    const d = new Date(windowStart + (t / tickCount) * windowSpan)
    const timeStr = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    if (timeStr === lastStr) continue
    lastStr = timeStr
    ctx.fillRect(tickX, h - 14, 1, 6)
    ctx.fillText(timeStr, Math.min(w - 40, Math.max(2, tickX - 14)), h - 2)
  }

  if (!isExportMode) {
    ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 2; ctx.beginPath()
    ctx.moveTo(w / 2, 0); ctx.lineTo(w / 2, h); ctx.stroke()
  }
}
