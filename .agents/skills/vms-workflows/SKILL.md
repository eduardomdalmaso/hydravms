---
name: vms-workflows
description: Event notification workflow engine, Telegram bot integration, webhooks, rate limiting, and automated action pipelines for HydraVMS.
---

# HydraVMS Event Notification & Automation Workflows (vms-workflows)

Este skill define o funcionamento do motor de automação de eventos e envio de notificações em múltiplos canais (Telegram, WebSockets, Webhooks, Email, MQTT).

---

## 1. Arquitetura do Pipeline de Workflows

```text
[ Evento Camada Gold (NATS JetStream) ]
                 │
                 ▼
     [ 1. Trigger Evaluator ] ── (Filtra por Câmera, Tipo de Evento, Severidade)
                 │
                 ▼
     [ 2. Cooldown & Anti-Spam ] ── (Verifica janela de tempo, ex: max 1/30s)
                 │
                 ▼
     [ 3. Action Dispatcher (Parallel Workers) ]
        ├──► [ Telegram Bot ] (Texto formatado + Snapshot JPEG)
        ├──► [ WebSocket ] (Push em tempo real para operadores)
        ├──► [ HTTP Webhook ] (JSON assinado com HMAC-SHA256)
        └──► [ MQTT Broker ] (Disparo de sirene/relé IoT)
```

---

## 2. Integração Nativa do Telegram

### Setup do Canal Telegram:
- **`bot_token`:** Token gerado via `@BotFather`.
- **`chat_id`:** ID do grupo ou usuário (ex: `-100192837465`).
- **`include_snapshot`:** Booleano (baixa frame cropado do MinIO via Presigned URL e envia como `multipart/form-data` no endpoint `sendPhoto`).

### Template da Mensagem Telegram:
```text
🚨 ALERTA DE SEGURANÇA // HYDRA VMS
━━━━━━━━━━━━━━━━━━━━
📹 Câmera: Portaria Principal [CAM_01]
🎯 Evento: Invasão de Perímetro (Pessoa)
📊 Confiança: 94.2%
⏰ Horário: 2026-09-03 19:45:00 UTC
━━━━━━━━━━━━━━━━━━━━
🔗 Visualizar no VMS: https://vms.corp/events/view/evt_9a8b7c
```

---

## 3. Prevenção de Spam (Cooldown & Debounce)

- **`cooldown_seconds` (Padrão: 30s):** Se múltiplos frames consecutivos detectarem a mesma pessoa na mesma zona, apenas o primeiro frame gera notificação imediata.
- **Circuit Breaker de Notificação:** Se um webhook externo falhar 5 vezes seguidas (timeout ou HTTP 5xx), ele entra em quarentena temporária para não consumir recursos.
