export function drawCanvasTimeline(
  ctx: CanvasRenderingContext2D, w: number, h: number, currentTime: number,
  zoomMinutes: number, isExportMode: boolean, exportStart?: number, exportEnd?: number,
  recordedRanges: { start: number; end: number }[] = []
) {
  const now = Date.now(), windowSpan = zoomMinutes * 60000, windowStart = currentTime - windowSpan / 2
  const windowEnd = windowStart + windowSpan, nowX = ((now - windowStart) / windowSpan) * w

  // 1. Clean Slate Dark Canvas Track Background (Past vs Future)
  ctx.fillStyle = '#0e1117'; ctx.fillRect(0, 0, w, h)
  if (nowX < w) {
    ctx.fillStyle = '#07090d'; ctx.fillRect(Math.max(0, nowX), 0, w - Math.max(0, nowX), h)
  }

  // 2. Real Recorded Segments from Database Only (HUD Orange #ff5e3a)
  ctx.fillStyle = isExportMode ? 'rgba(255, 94, 58, 0.35)' : '#ff5e3a'
  for (const r of recordedRanges) {
    if (r.end < windowStart || r.start > windowEnd) continue
    const segX1 = ((Math.max(r.start, windowStart) - windowStart) / windowSpan) * w
    const segX2 = ((Math.min(r.end, windowEnd) - windowStart) / windowSpan) * w
    if (segX2 > segX1) {
      ctx.fillRect(segX1, 10, Math.max(2, segX2 - segX1), h - 28)
    }
  }

  // 3. Export Range Selection Brackets
  if (isExportMode && exportStart && exportEnd) {
    const x1 = ((exportStart - windowStart) / windowSpan) * w, x2 = ((exportEnd - windowStart) / windowSpan) * w
    ctx.fillStyle = 'rgba(255, 94, 58, 0.25)'; ctx.fillRect(Math.min(x1, x2), 0, Math.abs(x2 - x1), h)
    ctx.fillStyle = '#ff5e3a'; ctx.strokeStyle = '#ff5e3a'; ctx.lineWidth = 2
    ctx.beginPath(); ctx.moveTo(x1, 0); ctx.lineTo(x1, h); ctx.moveTo(x2, 0); ctx.lineTo(x2, h); ctx.stroke()
  }

  // 4. Time Ticks and Numbers
  ctx.font = '500 10px Roboto, sans-serif'
  const ticks = zoomMinutes >= 720 ? 12 : zoomMinutes <= 5 ? 5 : 8
  for (let t = 0; t <= ticks; t++) {
    const tx = (t / ticks) * w, d = new Date(windowStart + (t / ticks) * windowSpan)
    const lbl = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    ctx.fillStyle = tx > nowX ? 'rgba(255, 255, 255, 0.2)' : '#8b94a0'
    ctx.fillRect(tx, h - 14, 1, 5)
    ctx.fillText(lbl, Math.min(w - 36, Math.max(4, tx - 13)), h - 3)
  }

  // 5. Real-Time Edge Line (at NOW)
  if (nowX >= 0 && nowX <= w) {
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)'; ctx.lineWidth = 1; ctx.beginPath()
    ctx.moveTo(nowX, 0); ctx.lineTo(nowX, h); ctx.stroke()
  }

  // 6. Playhead Scrubber with Current Timestamp Badge
  if (!isExportMode) {
    const curX = ((currentTime - windowStart) / windowSpan) * w
    ctx.strokeStyle = '#ff5e3a'; ctx.lineWidth = 2; ctx.beginPath()
    ctx.moveTo(curX, 0); ctx.lineTo(curX, h); ctx.stroke()

    // Scrubber Timestamp Pill
    const curD = new Date(currentTime)
    const curTimeStr = `${String(curD.getHours()).padStart(2, '0')}:${String(curD.getMinutes()).padStart(2, '0')}:${String(curD.getSeconds()).padStart(2, '0')}`
    ctx.fillStyle = '#ff5e3a'; ctx.fillRect(Math.max(2, Math.min(w - 52, curX - 25)), 0, 50, 12)
    ctx.fillStyle = '#000000'; ctx.font = '700 8px JetBrains Mono'
    ctx.fillText(curTimeStr, Math.max(6, Math.min(w - 48, curX - 21)), 9)
  }
}
