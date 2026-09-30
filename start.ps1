# ==============================================================================
# Hydra Ecosystem - 1-Click Turnkey Bootstrapper & PM2 Orchestrator
# Builds missing binaries, installs frontend deps, and launches all 4 planes.
# ==============================================================================

$ErrorActionPreference = "Stop"

Write-Host "========================================================" -ForegroundColor Cyan
Write-Host "   HYDRA ECOSYSTEM // AUTO-BOOTSTRAP & PM2 LAUNCHER     " -ForegroundColor Cyan
Write-Host "========================================================" -ForegroundColor Cyan

$RootDir = $PSScriptRoot
$ParentDir = Split-Path $RootDir -Parent

# 1. Build HydraVMS
Write-Host "`n[1/5] Verificando HydraVMS Backend..." -ForegroundColor Yellow
if (-not (Test-Path "$RootDir\hydravms.exe")) {
    Write-Host "  -> Compilando hydravms.exe..." -ForegroundColor Gray
    Push-Location $RootDir
    go build -o hydravms.exe ./cmd/hydravms
    Pop-Location
    Write-Host "  -> hydravms.exe compilado com sucesso!" -ForegroundColor Green
} else {
    Write-Host "  -> hydravms.exe pronto." -ForegroundColor Green
}

# 2. Build HydraStream
$StreamDir = Join-Path $ParentDir "HydraStream"
Write-Host "`n[2/5] Verificando HydraStream (Data Plane RTSP)..." -ForegroundColor Yellow
if (Test-Path $StreamDir) {
    if (-not (Test-Path "$StreamDir\hydrastream.exe")) {
        Write-Host "  -> Compilando hydrastream.exe..." -ForegroundColor Gray
        Push-Location $StreamDir
        go build -o hydrastream.exe ./cmd/hydrastream
        Pop-Location
        Write-Host "  -> hydrastream.exe compilado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "  -> hydrastream.exe pronto." -ForegroundColor Green
    }
}

# 3. Build HydraForge
$ForgeDir = Join-Path $ParentDir "HydraForge"
Write-Host "`n[3/5] Verificando HydraForge (Laboratório de IA & Inferência)..." -ForegroundColor Yellow
if (Test-Path $ForgeDir) {
    if (-not (Test-Path "$ForgeDir\hydraforge.exe")) {
        Write-Host "  -> Compilando hydraforge.exe..." -ForegroundColor Gray
        Push-Location $ForgeDir
        go build -o hydraforge.exe ./cmd/hydraforge
        Pop-Location
        Write-Host "  -> hydraforge.exe compilado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "  -> hydraforge.exe pronto." -ForegroundColor Green
    }
}

# 4. FFmpeg Tooling Check
Write-Host "`n[4/6] Verificando FFmpeg & Codecs de Vídeo..." -ForegroundColor Yellow
$ffmpegCmd = Get-Command ffmpeg -ErrorAction SilentlyContinue
if (-not $ffmpegCmd) {
    $wingetFfmpeg = Get-ChildItem -Path "$env:LOCALAPPDATA\Microsoft\WinGet\Packages" -Filter "ffmpeg.exe" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty FullName
    if (-not $wingetFfmpeg) {
        Write-Host "  -> FFmpeg não encontrado. Auto-instalando via winget..." -ForegroundColor Gray
        winget install Gyan.FFmpeg --accept-source-agreements --accept-package-agreements | Out-Null
        Write-Host "  -> FFmpeg instalado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "  -> FFmpeg localizado em: $wingetFfmpeg" -ForegroundColor Green
    }
} else {
    Write-Host "  -> FFmpeg pronto no PATH." -ForegroundColor Green
}

# 5. Frontend Dependencies
Write-Host "`n[5/6] Verificando dependências do Frontend..." -ForegroundColor Yellow
if (-not (Test-Path "$RootDir\web\node_modules")) {
    Write-Host "  -> Instalando dependências npm no frontend..." -ForegroundColor Gray
    Push-Location "$RootDir\web"
    npm install
    Pop-Location
    Write-Host "  -> Dependências do frontend instaladas!" -ForegroundColor Green
} else {
    Write-Host "  -> node_modules do frontend pronto." -ForegroundColor Green
}

# 6. PM2 Check and Launch
Write-Host "`n[6/6] Inicializando serviços no PM2..." -ForegroundColor Yellow
Push-Location $RootDir
pm2 delete all 2>$null | Out-Null
pm2 start ecosystem.config.js
pm2 save | Out-Null
Pop-Location

Write-Host "`n========================================================" -ForegroundColor Green
Write-Host "   ECOSSISTEMA HYDRA PRONTO E OPERACIONAL EM SEGUNDOS   " -ForegroundColor Green
Write-Host "========================================================" -ForegroundColor Green
Write-Host "  [0] HydraVMS API:    http://localhost:8083" -ForegroundColor White
Write-Host "  [1] HydraVMS Web UI: http://localhost:5173" -ForegroundColor White
Write-Host "  [2] HydraStream:     http://localhost:8080" -ForegroundColor White
Write-Host "  [3] HydraForge:      http://localhost:8081" -ForegroundColor White
Write-Host "========================================================`n" -ForegroundColor Green
