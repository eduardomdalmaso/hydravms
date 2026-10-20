---
name: vms-ecosystem-bootstrap
description: Procedimentos e arquitetura para orquestrar e subir o ecossistema completo (HydraStream, HydraVMS e ReductorPrompt) de forma genérica no Windows e Linux via PM2 e scripts.
---

# Hydra Suite Ecosystem Bootstrap (vms-ecosystem-bootstrap)

Este skill define o padrão oficial para inicialização, orquestração e gerenciamento de processos do ecossistema Hydra (**HydraStream**, **HydraVMS** e **ReductorPrompt**) tanto em ambientes **Windows** quanto **Linux**.

---

## 1. Topologia de Serviços & Mapeamento de Portas

| Serviço | Plano / Papel | Runtime | Porta Local | Protocolo |
| :--- | :--- | :--- | :---: | :--- |
| **`hydra-stream`** | Data Plane / Ingestão Zero-Copy | Go (`hydrastream`) + Rust (`hydra-engine`) | `:8080`<br>`:8085`<br>`:1935` | HTTP REST API / Snapshots<br>WSP (WebSocket WebCodecs)<br>RTMP Push Server |
| **`hydra-vms-api`** | Control Plane / Orquestrador | Go (`hydravms`) | `:8083` | HTTP REST API & WebSockets (JWT Token-Only) |
| **`hydra-vms`** | Frontend / Painel VMS | Node.js / Vite (Vue 3) | `:5173` | HTTP Web SPA |
| **`reductor-api`** | Knowledge Engine / RAG / MCP | Python / FastAPI (Uvicorn) | `:8002` | HTTP REST & MCP Call Tool Hub |
| **`reductor-web`** | Frontend / Painel Reductor | Node.js / Vite (React) | `:5174` | HTTP Web SPA |
| **`hydra-forge`** | Analytics Plane (Opcional) | Go / Python / TensorRT | `:8081` | Treinamento YOLO na RTX 5090 |

---

## 2. Padrão de Configuração Multiplataforma (PM2)

Para evitar vazamento de caminhos absolutos e manter o repositório limpo:
1. **Templates no Git:** Cada projeto mantém um `ecosystem.config.example.js` (ou `.cjs`) comitado.
2. **Configuração Local:** O arquivo real `ecosystem.config.js` é gerado localmente e ignorado pelo `.gitignore`.
3. **Resolução Dinâmica:** Caminhos de binários e executáveis adaptam-se automaticamente ao `process.platform`:

```javascript
const path = require("path");
const isWin = process.platform === "win32";

module.exports = {
  apps: [
    {
      name: "hydra-stream",
      cwd: __dirname,
      script: isWin ? "./hydrastream.exe" : "./bin/hydrastream",
      interpreter: "none",
      autorestart: true,
      max_restarts: 10,
      restart_delay: 2000,
      env: {
        PORT: "8080",
        JWT_SECRET: "hydravms-enterprise-secure-jwt-signing-key-32b-secret",
      },
    },
  ],
};
```

---

## 3. Comandos de Inicialização Rápida

### A. Subida Completa via PM2 (Recomendado)

#### No Windows (PowerShell):
```powershell
# 1. Subir Data Plane (HydraStream)
cd C:\Users\hades\Documents\HydraStream
pm2 start ecosystem.config.js

# 2. Subir Control Plane (HydraVMS Backend + Frontend)
cd C:\Users\hades\Documents\hydravms
pm2 start ecosystem.config.js

# 3. Subir Knowledge & Tool Hub (ReductorPrompt API + Web)
cd C:\Users\hades\Documents\reductor-prompt
$env:PM2_TARGET="reductor"; pm2 start ecosystem.config.cjs

# 4. Verificar status geral
pm2 list
```

#### No Linux (Bash):
```bash
# 1. Subir Data Plane (HydraStream)
cd ~/Documents/HydraStream && pm2 start ecosystem.config.js

# 2. Subir Control Plane (HydraVMS)
cd ~/Documents/hydravms && pm2 start ecosystem.config.js

# 3. Subir Knowledge & Tool Hub (ReductorPrompt)
cd ~/Documents/reductor-prompt && PM2_TARGET=reductor pm2 start ecosystem.config.cjs

# 4. Verificar status geral
pm2 list
```

---

## 4. Gestão e Ciclo de Vida do Cluster PM2

| Ação | Comando |
| :--- | :--- |
| **Listar processos ativos** | `pm2 list` ou `pm2 status` |
| **Ver logs consolidados** | `pm2 logs` |
| **Ver logs de serviço específico** | `pm2 logs hydra-stream --lines 50` |
| **Reiniciar todos os serviços** | `pm2 restart all` |
| **Parar todos os serviços** | `pm2 stop all` |
| **Salvar estado para reboot** | `pm2 save` |

---

## 5. Diretrizes Invioláveis para Windows (Window Suppression)

No Windows, todo processo de console disparado por código (ex: `hydrastream` iniciando `hydra-engine.exe` ou `ffmpeg.exe`) **DEVE** utilizar o helper `proc.SetHideWindow(cmd)` para evitar pop-up de terminais indesejados na tela:

```go
// Exemplo em Go no HydraStream
cmd := exec.Command(binPath, "server", "--port", port, "--rtmp-port", rtmpPort)
proc.SetHideWindow(cmd) // Aplica CREATE_NO_WINDOW (0x08000000) e HideWindow = true
cmd.Start()
```

---

## 6. Checklist de Validação Pós-Boot

Após iniciar todos os serviços, valide a saúde dos endpoints:

```powershell
# 1. HydraStream Data Plane
(Invoke-WebRequest -Uri "http://127.0.0.1:8080" -UseBasicParsing).StatusCode # Retorna 200

# 2. HydraVMS Control Plane (Segurança Estrita Token-Only)
# Deve retornar 401 Unauthorized para requests sem token (comportamento correto)
(Invoke-WebRequest -Uri "http://127.0.0.1:8083/health" -UseBasicParsing).StatusCode

# 3. ReductorPrompt Knowledge API
(Invoke-WebRequest -Uri "http://127.0.0.1:8002/health" -UseBasicParsing).StatusCode # Retorna 200

# 4. Frontends Web
(Invoke-WebRequest -Uri "http://127.0.0.1:5173" -UseBasicParsing).StatusCode # HydraVMS Web 200
(Invoke-WebRequest -Uri "http://127.0.0.1:5174" -UseBasicParsing).StatusCode # Reductor Web 200
```
