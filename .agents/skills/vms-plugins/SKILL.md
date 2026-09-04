---
name: vms-plugins
description: Architecture, manifest specification, lifecycle management, hot-reload, and marketplace distribution for HydraVMS analytics and device integration plugins.
---

# HydraVMS Plugin System & Marketplace Architecture (vms-plugins)

Este skill define a especificação técnica para desenvolvimento, empacotamento, instalação, execução isolada e atualização sem downtime de **Plugins de Analíticos e Integrações** no HydraVMS, inspirado no modelo de extensibilidade de plataformas como Vezha, Scrypted e Homebridge.

---

## 1. Princípios do Sistema de Plugins

1. **Process Isolation (Sandboxing):** Cada plugin roda como um subprocesso independente ou container isolado (IPC via NATS JetStream e gRPC), garantindo que falhas em um analítico não afetem o Control Plane do VMS.
2. **Dynamic Hot-Reload (Zero Downtime):** Plugins podem ser instalados, iniciados, parados, reconfigurados e atualizados dinamicamente sem reiniciar o servidor principal Go ou pipelines RTSP.
3. **Multi-Tenant Scoping:** Plugins são registrados globalmente no catálogo/marketplace, mas ativados e configurados com credenciais e parâmetros isolados por tenant.
4. **Hardware Acceleration Pass-Through:** Suporte transparente a `/dev/shm` (Zero-Copy frame sharing) e device nodes de GPU (`/dev/nvidia*`, `/dev/dri`) para analíticos neurais de alta performance.

---

## 2. Especificação do Manifesto (`plugin.json`)

Todo plugin deve conter um arquivo `plugin.json` na raiz de seu pacote:

```json
{
  "id": "hydra-analytics-lpr",
  "name": "LPR / ALPR Vehicle Identification",
  "version": "1.4.2",
  "author": "Hydra Vision Lab",
  "description": "Reconhecimento automático de placas veiculares brasileiras (Mercosul e antigas) com SAHI e YOLO11.",
  "category": "analytics",
  "runtime": "python3",
  "entrypoint": "main.py",
  "min_vms_version": "1.0.0",
  "permissions": [
    "stream:read",
    "events:publish",
    "shm:read",
    "gpu:cuda"
  ],
  "config_schema": {
    "type": "object",
    "properties": {
      "confidence_threshold": {
        "type": "number",
        "title": "Limiar de Confiança",
        "default": 0.65,
        "minimum": 0.1,
        "maximum": 1.0
      },
      "sahi_enabled": {
        "type": "boolean",
        "title": "Habilitar SAHI (Placas Distantes)",
        "default": true
      },
      "region": {
        "type": "string",
        "title": "Região de Placas",
        "enum": ["mercosul", "usa", "eu"],
        "default": "mercosul"
      }
    },
    "required": ["confidence_threshold"]
  },
  "ui": {
    "dashboard_card": true,
    "live_reticle": "lpr_box",
    "settings_component": "LprSettingsView"
  }
}
```

---

## 3. Ciclo de Vida do Plugin (Plugin Lifecycle)

O `PluginManager` do HydraVMS gerencia os seguintes estados:

```text
[ MARKETPLACE / LOCAL PKG ]
          │ (Download / Validate Signature)
          ▼
   [ INSTALLED ] ──► (Configure Params)
          │
          ▼ (Start Request)
   [ INITIALIZING ] ──► (Verify GPU / Alloc SHM)
          │
          ▼
     [ RUNNING ] ◄──► [ NATS JetStream (Frames In / Events Out) ]
          │
          ├──► (Config Change) ──► [ HOT-RELOADING ] ──► [ RUNNING ]
          ├──► (Stop Request)   ──► [ STOPPED ]
          └──► (Crash/Fault)    ──► [ CRASH_LOOP ] ──► (Circuit Breaker)
```

### Comandos de Controle via NATS:
- `hydra.plugins.<tenant_id>.<plugin_id>.start`
- `hydra.plugins.<tenant_id>.<plugin_id>.stop`
- `hydra.plugins.<tenant_id>.<plugin_id>.reload_config`
- `hydra.plugins.<tenant_id>.<plugin_id>.health`

---

## 4. Atualização Blue/Green de Plugins Sem Perda de Frames

1. **Stage New Version:** O VMS baixa e extrai `v1.4.3` em `/var/lib/hydravms/plugins/<id>/1.4.3`.
2. **Warm-up:** O subprocesso da nova versão é inicializado, carrega pesos no CUDA/VRAM e conecta no NATS.
3. **Drain & Swap:** O VMS redireciona a subscrição de frames RTSP/SHM para a nova versão.
4. **Graceful Terminate:** A versão antiga recebe sinal `SIGTERM`, processa os frames pendentes no buffer e encerra.

---


---

## 5. Estratégia de Rollback & Auto-Recuperação (Self-Healing)

Se uma nova versão de plugin falhar (ex: bug no script, incompatibilidade CUDA ou crash repetitivo):

1. **Retenção de Versões (Multi-Version Directory):**
   - Diretórios anteriores são preservados: `/var/lib/hydravms/plugins/<id>/1.4.1` e `1.4.2`.
2. **Rollback Manual em 1-Clique:**
   - O operador seleciona a versão anterior via painel web ou API (`POST /api/v1/plugins/<id>/rollback`).
   - O `PluginManager` mata o processo atual e sobe a versão estável em menos de 1 segundo.
3. **Auto-Rollback por Circuit Breaker (Automático):**
   - Se a nova versão falhar 3 vezes em 30 segundos (ou entrar em `CRASH_LOOP`), o VMS dispara o **Auto-Rollback** automático para a `last_stable_version`.
   - Um alerta crítico é emitido na UI/audit log com os logs de erro da versão com falha.
4. **Config Rollback Protection:**
   - Snapshots do `config_values` são vinculados a cada versão para evitar que parâmetros novos quebrem a versão restaurada.

## 6. Estrutura do Marketplace

O HydraVMS possui um catálogo integrado com suporte a:
- **Repositório Oficial Remoto:** Busca plugins assinados digitalmente (SHA256 + ECDSA).
- **Side-loading Local:** Upload de arquivos `.hpk` (Hydra Plugin Package - tar.gz assinado) pelo painel web.
- **Auto-Update Policy:** Atualizações automáticas de patches de segurança ou aprovação manual pelo administrador do tenant.
