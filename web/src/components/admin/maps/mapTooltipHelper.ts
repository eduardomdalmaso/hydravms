export function getCameraTooltipHtml(name: string, lat: number, lng: number): string {
  return `
    <div class="vms-map-hover-card camera-card">
      <div class="vms-hover-snapshot-box">
        <svg width="28" height="28" viewBox="0 0 576 512" fill="#ff5e3a" style="opacity:0.9"><path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/></svg>
        <span class="vms-snap-badge snap-live">[LIVE 30 FPS]</span>
        <span class="vms-snap-badge snap-res">1080P // H.265</span>
      </div>
      <div class="vms-hover-info">
        <div class="vms-hover-title">${name}</div>
        <div class="vms-hover-mono">GEO: ${lat.toFixed(5)}, ${lng.toFixed(5)}</div>
        <div class="vms-hover-status">TRANSMISSAO ESTAVEL // 29ms</div>
      </div>
    </div>
  `
}

export function getAlarmTooltipHtml(name: string, lat: number, lng: number): string {
  return `
    <div class="vms-map-hover-card alarm-card">
      <div class="vms-hover-alarm-header">
        <div class="vms-hover-title">${name}</div>
        <div class="vms-hover-alarm-badge">[3 DISPAROS RECENTES]</div>
      </div>
      <div class="vms-hover-trigger-list">
        <div class="vms-trigger-item">
          <span class="vms-trigger-time">14:32:10</span>
          <span class="vms-trigger-desc">VIOLACAO DE PERIMETRO</span>
        </div>
        <div class="vms-trigger-item">
          <span class="vms-trigger-time">12:15:04</span>
          <span class="vms-trigger-desc">DETECCAO DE MOVIMENTO</span>
        </div>
        <div class="vms-trigger-item">
          <span class="vms-trigger-time">09:40:22</span>
          <span class="vms-trigger-desc">TESTE DE ZONA MANUAL</span>
        </div>
      </div>
      <div class="vms-hover-mono">GEO: ${lat.toFixed(5)}, ${lng.toFixed(5)}</div>
    </div>
  `
}
