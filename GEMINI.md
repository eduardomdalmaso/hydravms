# HydraVMS Project Guidelines (GEMINI.md)

Este documento define a arquitetura de software, **protocolos de ingestão de câmeras (ONVIF WS-Discovery, RTSP, RTMP)**, **motor de geração de relatórios (PDF, Excel, CSV) & consultas analíticas otimizadas**, **regras de segurança e acesso estritamente via Token**, **escalabilidade em cluster distribuído**, **comunicação inter-servidores (gRPC / NATS / WebRTC)**, **especificação REST API & Swagger/OpenAPI 3.0**, **padrões visuais de UI/UX das páginas**, padrão Medallion de dados, **arquitetura de banco de dados relacional SQL (PostgreSQL / SQLite WAL)**, pipelines de IA com **SAHI**, **Inferência Acionada por Movimento (Motion-Gated)**, **modos de gravação**, política de armazenamento MinIO S3, isolamento multi-tenant, **autopreservação de disco (Anti-Crash Guard)**, exportação ultrarrápida e protocolo de auto-reparo do ecossistema **HydraVMS** (Video Management System, Orquestrador de Eventos de IA e Gravação).

---

## 1. Protocolos de Ingestão de Câmeras (ONVIF, RTSP & RTMP)

As diretrizes detalhadas encontram-se em [`.agents/skills/vms-camera-ingest/SKILL.md`](file:///home/hades/Documents/HydraVMS/.agents/skills/vms-camera-ingest/SKILL.md):

```text
[ Câmeras IP ONVIF ] ──► [ WS-Discovery Probe ] ──► Auto-Fetch Main/Sub RTSP URIs
[ Câmeras RTSP ]     ──► [ HydraStream Engine ] ──► TCP Interleaved Zero-Copy RingBuffer
[ Drones / OBS / App]──► [ RTMP Server :1935 ]  ──► Validação de Stream Key & Pipeline
```

### Detalhes dos Protocolos Suportados:
1. **ONVIF (Profile S/T/G):**
   - **Auto-Descoberta:** Envio de probe multicast UDP para `239.255.255.250:3702` para listar câmeras na rede local sem digitação manual de IPs.
   - **Extração Automática:** Obtém perfis Main (1080p/4K @ 30 FPS) e Sub-stream (480p @ 10 FPS) com digest authentication e controle PTZ.
2. **RTSP (RFC 2326):**
   - Ingestão contínua em modo **TCP Interleaved RTP** via HydraStream, eliminando artefatos e perda de pacotes de redes instáveis.
3. **RTMP Push Server (Porta 1935):**
   - Permite que drones de segurança, bodycams ou encoders móveis transmitam diretamente para o HydraVMS através de `rtmp://vms.domain.com/live/{tenant_id}/{camera_id}?key={stream_key}` com autenticação instantânea de chave.

---

## 2. Motor de Relatórios & Consultas Analíticas Otimizadas

As diretrizes detalhadas encontram-se em [`.agents/skills/vms-reports-analytics/SKILL.md`](file:///home/hades/Documents/HydraVMS/.agents/skills/vms-reports-analytics/SKILL.md):

### 2.1. Geração Assíncrona de Relatórios:
- **PDF Executivo / Forense:** Relatórios completos com gráficos de tendências, tabelas de alarmes, thumbnails de snapshots embutidos e links auditáveis para download de clipes MP4 de evidência.
- **Excel (XLSX) & CSV Streaming:** Exportação em lote de centenas de milhares de linhas utilizando streaming com `excelize` em Go, consumindo menos de 20MB de RAM.
- **Armazenamento Seguro:** Arquivos gerados são salvos no bucket MinIO `hydravms-reports/` e disponibilizados via Presigned URLs com expiração de 24 horas.

### 2.2. Consultas Analíticas Otimizadas (< 10ms):
- **Materialized Views:** Rollups de agregação horária e diária (`mv_hourly_event_stats`) pré-computados para carregamento instantâneo de dashboards.
- **Índices Parciais Seletivos:** Índices dedicados para eventos não reconhecidos (`status = 'new'`) e eventos de alta severidade (`severity = 'critical'`).
- **Métricas de Performance Operacional:** Cálculo automático de MTTA (*Mean Time to Acknowledge*), MTTR (*Mean Time to Resolve*), ranking de câmeras com mais incidentes e taxa de falsos positivos.

---

## 3. Acesso à API Estritamente via Token (Token-Only Enforcement)

As diretrizes detalhadas encontram-se em [`.agents/rules/security.md`](file:///home/hades/Documents/HydraVMS/.agents/rules/security.md):

- **Proibição de Rotas Públicas Anônimas:** Todos os endpoints REST (`/api/v1/*`), feeds WebRTC e conexões WebSocket (`/ws/v1/*`) EXIGEM autenticação obrigatória via Token.
- **Única Rota Sem Token:** `POST /api/v1/auth/login` (emissão inicial de token).
- **Tipos de Tokens Aceitos:** User JWT Bearer Token e API Keys / Service Tokens (M2M) com escopos e expiração.

---

## 4. Escalabilidade em Cluster & Comunicação Inter-Servidores

As diretrizes detalhadas encontram-se em [`.agents/skills/vms-clustering/SKILL.md`](file:///home/hades/Documents/HydraVMS/.agents/skills/vms-clustering/SKILL.md):

- **Topologia do Cluster:** Control Plane Nodes (API/Auth), Edge Ingest Nodes (RTSP/ONVIF/RTMP), GPU Worker Nodes (YOLO/SAHI na RTX 5090) e MinIO S3 Distributed Pool.
- **Matriz de Comunicação:** POSIX `/dev/shm` local, **gRPC (HTTP/2 + Protobuf)** entre servidores, **NATS JetStream Mesh** para eventos (> 10M msgs/s) e **WebRTC (WHIP/WHEP)** para clientes web (< 200ms).

---

## 5. API REST & Documentação Swagger / OpenAPI 3.0

As diretrizes detalhadas encontram-se em [`.agents/skills/vms-api-rest/SKILL.md`](file:///home/hades/Documents/HydraVMS/.agents/skills/vms-api-rest/SKILL.md) e a especificação completa em [`api/openapi.yaml`](file:///home/hades/Documents/HydraVMS/api/openapi.yaml):

- **Swagger UI Interativo:** `GET /swagger/index.html` embutido no binário Go.
- **Padrão de Resposta de Erro:** RFC 7807 Problem Details.

---

## 6. Padrões Visuais de UI/UX das Páginas (Clean Slate Dark)

As diretrizes detalhadas encontram-se em [`.agents/rules/ui-standards.md`](file:///home/hades/Documents/HydraVMS/.agents/rules/ui-standards.md):

- **Tokens de Cores:** Fundo `#0f172a`, Superfícies `#1e293b`, Inputs `#334155`, Primário `#3b82f6`, Sucesso `#10b981`, Alerta `#f59e0b`, Perigo `#ef4444`.
- **Tipografia:** `'Inter'` / `'Geist'` para textos/títulos e `'JetBrains Mono'` para números, FPS e timestamps.
- **Header:** 56px com Breadcrumb, Telemetria compacta e Perfil.
- **Sidebar:** 240px (expandida) / 64px (recolhida) com indicador de rota ativa com borda azul.
- **Modularidade:** **Máximo de 100 linhas por arquivo** (`.vue`, `.ts`, `.css`).
- **Nomenclatura Concisa & Botões:** Rótulos atômicos (`GRAVANDO`, `ANALITICOS`, `CODEC`, `COMPRESSAO`, `RESOLUCAO`, `FPS`) e texto de botão de confirmação sempre **`SALVAR`** (detalhes em [`.agents/rules/ui-telemetry-naming.md`](.agents/rules/ui-telemetry-naming.md)).

---

## 7. Arquitetura de Banco de Dados Relacional SQL (PostgreSQL 15+ / SQLite WAL)

Consulte o DDL completo em [`migrations/000001_initial_schema.up.sql`](file:///home/hades/Documents/HydraVMS/migrations/000001_initial_schema.up.sql):

- **Tenants, RBAC & Tokens:** `tenants`, `users`, `user_sessions`, `api_tokens`.
- **Câmeras (ONVIF/RTSP/RTMP):** `cameras`, `camera_recording_profiles`, `camera_zones`.
- **IA & Alarmes:** `ai_rules`, `events` (Gold Layer), `event_audit_actions`.
- **Gravações, Relatórios & Storage:** `recordings`, `report_jobs`, `storage_pools`, `storage_purge_logs`, `video_exports`.
- **Auditoria:** `audit_logs`.

---

## 8. Pipeline de IA com SAHI & Motion-Gated Inference

1. **Motion-Gated Pre-Filter (CPU):** Compara pixels em grayscale 160x120 (< 0.1ms) e descarta frames sem variação, poupando **80% a 95% de GPU**.
2. **SAHI (Slicing Aided Hyper Inference):** Fatiamento em janelas 640x640 (overlap 20%) nas ROIs ativas para detecção precisa de alvos pequenos em 1080p/4K com Non-Maximum Merging (NMM).

---

## 9. Modos de Gravação & Autopreservação de Armazenamento

- **Perfis de Gravação:** Contínua (24/7), Movimento (VMD), Evento de IA (Smart Alarm) e Híbrido.
- **Mecanismo Anti-Crash (StorageGuard):** Reserva inviolável de 5% de disco; purges em cascata aos 80%, 90%, 95% (Circuit Breaker) e 98% (Hard Lock), protegendo provas judiciais (`is_pinned = true`).
- **Exportação Rápida:** Zero-Reencode Stream-Copy (1 hora em 1-3s) e gravação de watermark via NVIDIA NVENC a 800+ FPS.

---

## 10. Arquitetura de Plugins Dinâmicos & Marketplace (Estilo Vezha / Scrypted)

Para permitir a evolução de analíticos sem interromper o ecossistema ou re-compilar o binário principal Go:
- **Manifesto `plugin.json`:** Define metadados, permissões, schema de configuração e telas customizadas.
- **Isolamento de Processos (Sandboxing):** Plugins rodam como subprocessos isolados comunicando via IPC (NATS JetStream e /dev/shm para Zero-Copy frames).
- **Hot-Reload & Blue/Green Updates:** Atualizações de analíticos sobem nova versão em paralelo, transferem subscrições e finalizam a versão antiga com zero descarte de frames.
- **Marketplace Multi-Tenant:** Catálogo de analíticos (LPR, Facial, EPI, Fogo/Fumaça) com ativação e parametrização independente por tenant.


---

## 11. Motor de Workflows de Notificação & Integrações (Telegram, WebSockets, Webhooks)

Para automação inteligente de despacho de alertas em tempo real:
- **Pipeline Visual (Trigger ➔ Condições/Filtros ➔ Ações/Destinos):** Permite aos operadores configurar para onde cada tipo de evento de IA deve ser enviado.
- **Canal Nativo Telegram Bot:** Disparo instantâneo com texto formatado (câmera, tipo de evento, timestamp, confiança) + foto/snapshot cropado da cena via `sendPhoto`.
- **Canais Múltiplos Simultâneos:** Disparo em paralelo para WebSockets (operadores online), HTTP Webhooks (sistemas externos de segurança/SIEM), Email e MQTT.
- **Anti-Spam Cooldown & Circuit Breaker:** Janela de silenciamento configurável (ex: 30s) para evitar disparos repetidos da mesma detecção consecutiva.


---

## 12. Arquitetura Híbrida, Cluster de Nós & Playback Sincronizado

### 1. Protocolo de Live Streaming de Ultrabaixa Latência
- **Padrão:** WebRTC (WHEP - WebRTC HTTP Egress Protocol) com latência `< 300ms`.
- **Fallback Automático:** WebSocket-fMP4 (MSE) para redes corporativas com bloqueio UDP/STUN.

### 2. Controles PTZ Inteligentes (Renderização Condicional)
- Controles de Pan/Tilt/Zoom (Virtual Joystick, Presets com thumbnails e Click-to-Center) são renderizados **exclusivamente** em câmeras com suporte ONVIF PTZ ativo (`camera.has_ptz == true`), mantendo a interface limpa em câmeras estáticas.

### 3. Timeline de Gravações Multi-Câmera com Master Clock UTC
- Barra de tempo unificada que comanda até 16 câmeras simultaneamente sincronizadas no mesmo milissegundo UTC.
- Botão *"Desacoplar Câmera"* para investigação individual em slot específico sem perder a referência do grid.

### 4. Operação Híbrida & Cluster Multi-Nó (Edge & Cloud)
- Suporte a servidores Edge locais e nós em nuvem com tela dedicada de configuração de peering e comunicação gRPC / NATS mesh.
- Buffer local inteligente: se a conexão com a nuvem cair, o Edge grava e executa IA 100% offline, sincronizando metadados Gold e snapshots ao restabelecer o link.

### 5. Conformidade Forense & Privacidade (LGPD)
- Exportação de MP4 com marca d'água de timestamp e Hash SHA-256 criptográfico contra adulteração.
- Registro imutável de auditoria (`audit_logs`) para toda visualização de stream ao vivo e download de evidências.
- Máscaras de privacidade dinâmicas por câmera (desfoque de faces e placas de terceiros).


---

## 13. Armazenamento Hierárquico, Buffer NVMe & Tratamento de H.265 (HEVC)

### 1. Desacoplamento em 4 Planos (Split-Plane)
- **Data Plane (HydraStream):** Ingestão e WebRTC WHEP em nós leves de CPU.
- **Analytics Plane (HydraForge):** TensorRT + SAHI em nós com GPU RTX 5090.
- **Storage Tier (StorageGuard):** Buffer NVMe de escrita rápida + Spillover para HDs / MinIO S3.
- **Control Plane (HydraVMS):** API REST, PostgreSQL, RBAC, Workflows e UI Vue 3.

### 2. Busca Rápida Global de Gravações (Direct-to-Source)
- Índice espaço-temporal centralizado (< 5ms) + Presigned URLs diretas ao storage de origem (Zero-Hop).
- Sprite Sheets WebP para pré-visualização instantânea na timeline (*scrubbing*) e HTTP 206 Range Requests para play imediato (< 200ms).

### 3. Gestão Dinâmica de Discos NVMe/HD via Painel Web
- Adição de novos discos em tempo real sem reiniciar o VMS (`storage_pools`).
- Seletor com recomendação inteligente (*self-advisory*):
  - `[HOT_VIDEO_BUFFER]`: XFS otimizado para blocos de vídeo de 64MB.
  - `[EVENTS_DATABASE]`: EXT4 dedicado para PostgreSQL Tablespace e SQLite WAL.
  - `[AI_SNAPSHOTS_CACHE]`: Armazenamento de fotos de alertas e sprites.
  - `[WARM_LONG_TERM]`: HDs mecânicos de longo prazo.
- Dreno contínuo (*Spillover Daemon*) com válvulas de nível de água (70%, 85%, 95%).

### 4. Estratégia de Vídeo H.265 (HEVC) & Dual-Stream
- Gravação e IA em H.265 original (economia de 50% de espaço); decodificação por GPU via NVDEC.
- Dual-Stream: Main Stream (H.265 4K) para IA/Gravação; Sub Stream (H.264) para grade multi-view no navegador.
- WebCodecs no cliente + JIT Transcode com NVIDIA NVENC sob demanda para navegadores legados (somente enquanto a aba estiver aberta).


---

## 14. Mosaico de Câmeras de Alta Performance & Playback Inferior

### 1. Otimização de Performance para 16+ Streams no Navegador
- **Comutação Automática de Sub-Stream:** Grades multi-câmera (2x2, 3x3, 4x4) utilizam o Sub-Stream leve (H.264 640x360 @ 15 FPS) da própria câmera, mantendo o consumo de CPU/GPU do navegador abaixo de 15%.
- **Hero / Single-View Switching:** Ao expandir uma câmera para 1x1, o stream comuta de forma transparente para o Main-Stream (Full-HD ou 4K).
- **Pausa Fora da Tela (IntersectionObserver):** Streams de abas inativas ou fora da viewport são pausados para economizar banda de rede e GPU.

### 2. Gaveta de Gravações na Parte Inferior (Bottom Playback Drawer)
- **Inspeção Instantânea em 1-Clique:** Clicar em qualquer slot de câmera do mosaico abre a gaveta inferior com a timeline sem perder o contexto visual das outras câmeras.
- **Barra Temporal Multicamadas:** Destaque em cores para Gravação Contínua (azul), Detecção de Movimento (amarelo) e Alarmes Críticos de IA (vermelho).
- **Controles de Avanço & Velocidade:** Botões de salto de ±10s, avanço quadro a quadro (Frame Step ±1) e velocidades estritamente limitadas para o ambiente web em **`1x`, `2x` e `3x`** (garantindo estabilidade de buffer e zero drop frames no navegador).
- **Ações de Saída:**
  - *"Tirar Gravação"*: Fecha a gaveta inferior e retorna a visualização ao vivo imediata.
  - *"Baixar Trecho (Clip)"*: Exporta o trecho selecionado em MP4 com Hash forense SHA-256.

## 15. Catálogo Completo de Skills do HydraVMS

| Skill | Finalidade |
| :--- | :--- |
| **`vms-mosaic-playback`** | Mosaico de alta performance (1x1 a 4x4), gaveta de timeline inferior e controle de velocidade (1x, 2x e 3x). |
| **`vms-workflows`** | Motor de automação de eventos, setup de canais (Telegram, WebSockets, Webhooks) e anti-spam. |
| **`vms-plugins`** | Arquitetura de plugins dinâmicos, hot-reload, manifesto `plugin.json` e marketplace. |
| **`vms-camera-ingest`** | Ingestão multi-protocolo, auto-descoberta ONVIF WS-Discovery, RTSP interleaved e servidor RTMP push. |
| **`vms-reports-analytics`** | Geração assíncrona de relatórios (PDF/Excel/CSV), agregações analíticas rápidas (<10ms) e métricas SLA. |
| **`vms-api-rest`** | Padrões REST, autenticação estrita por Token, Swagger/OpenAPI 3.0, códigos HTTP, RFC 7807 e WebSockets. |
| **`vms-clustering`** | Arquitetura multi-nó, escala horizontal, gRPC, NATS mesh, WebRTC e auto-discovery de servidores. |
| **`vms-database`** | Schema SQL relacional, tokens de API, índices compostos, consultas tenant-scoped e migrations. |
| **`vms-ui`** | Design system Clean Slate Dark, tokens, tipografia, layout e componentes modulares (< 100 linhas). |
| **`vms-events`** | Regras de eventos de IA, polígonos normalizados, triggers e ciclo de vida. |
| **`vms-recording`** | Perfis de gravação, MinIO S3 e Presigned URLs. |
| **`vms-storage-guard`** | Autopreservação de disco, watermarks, circuit breaker e anti-crash. |
| **`yolo-sahi`** | Fatiamento SAHI, batch inference na RTX 5090 e Non-Maximum Merging (NMM). |
| **`validate-project`** | Executa auditoria de conformidade DDD, contagem de linhas web e testes unitários. |

---

## 16. Protocolo de Auto-Reparo & Validação Pós-Prompt (Self-Healing Rules)

```bash
# 1. Auditoria de Limite de Linhas Web (< 100 linhas por arquivo)
wc -l web/src/**/*.vue web/src/**/*.ts web/src/**/*.css web/src/**/*.js 2>/dev/null | awk '$1 > 100 { print "VIOLATION: " $2 " has " $1 " lines (>100)" }'

# 2. Auditoria de Pureza DDD (Nenhum import de infra no domínio)
grep -rnE "(net/http|database/sql)" internal/domain/ && echo "VIOLAÇÃO DDD: Remova infraestrutura do domínio!"

# 3. Compilação e Testes Automatizados
make test && make build
```
