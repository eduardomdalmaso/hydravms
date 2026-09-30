#!/usr/bin/env bash
# ==============================================================================
# Hydra Ecosystem - 1-Click Turnkey Bootstrapper & PM2 Orchestrator (Linux)
# ==============================================================================
set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$ROOT_DIR")"

echo -e "\e[1;36m========================================================\e[0m"
echo -e "\e[1;36m   HYDRA ECOSYSTEM // AUTO-BOOTSTRAP & PM2 LAUNCHER     \e[0m"
echo -e "\e[1;36m========================================================\e[0m"

# 1. HydraVMS
echo -e "\n\e[1;33m[1/5] Compilando HydraVMS Backend...\e[0m"
if [ ! -f "$ROOT_DIR/bin/hydravms" ]; then
    mkdir -p "$ROOT_DIR/bin"
    (cd "$ROOT_DIR" && go build -o bin/hydravms ./cmd/hydravms)
fi

# 2. HydraStream
echo -e "\n\e[1;33m[2/5] Compilando HydraStream (Data Plane RTSP)...\e[0m"
if [ -d "$PARENT_DIR/HydraStream" ] && [ ! -f "$PARENT_DIR/HydraStream/bin/hydrastream" ]; then
    mkdir -p "$PARENT_DIR/HydraStream/bin"
    (cd "$PARENT_DIR/HydraStream" && go build -o bin/hydrastream ./cmd/hydrastream)
fi

# 3. HydraForge
echo -e "\n\e[1;33m[3/5] Compilando HydraForge (Laboratório IA)...\e[0m"
if [ -d "$PARENT_DIR/HydraForge" ] && [ ! -f "$PARENT_DIR/HydraForge/bin/hydraforge" ]; then
    mkdir -p "$PARENT_DIR/HydraForge/bin"
    (cd "$PARENT_DIR/HydraForge" && go build -o bin/hydraforge ./cmd/hydraforge)
fi

# 4. Frontend Dependencies
echo -e "\n\e[1;33m[4/5] Instalando dependências do Frontend...\e[0m"
if [ ! -d "$ROOT_DIR/web/node_modules" ]; then
    (cd "$ROOT_DIR/web" && npm install)
fi

# 5. Launch with PM2
echo -e "\n\e[1;33m[5/5] Inicializando serviços no PM2...\e[0m"
(cd "$ROOT_DIR" && pm2 delete all 2>/dev/null || true && pm2 start ecosystem.config.js && pm2 save)

echo -e "\n\e[1;32m========================================================\e[0m"
echo -e "\e[1;32m   ECOSSISTEMA HYDRA OPERACIONAL // TURNKEY READY       \e[0m"
echo -e "\e[1;32m========================================================\e[0m"
echo "  [0] HydraVMS API:    http://localhost:8083"
echo "  [1] HydraVMS Web UI: http://localhost:5173"
echo "  [2] HydraStream:     http://localhost:8080"
echo "  [3] HydraForge:      http://localhost:8081"
echo -e "\e[1;32m========================================================\e[0m\n"
