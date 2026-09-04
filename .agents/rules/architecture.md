# 🏛️ Regras de Arquitetura & Engenharia (HydraVMS)

Este documento estabelece os padrões arquiteturais obrigatórios para o backend (Go), o ecossistema de dados, ingestão de streaming e armazenamento distribuído do HydraVMS.

---

## 1. Arquitetura Hexagonal (Ports & Adapters) & DDD

A estrutura de código deve respeitar rigorosamente a separação de camadas:

1. **`internal/domain/` (Domínio Puro):**
   - Contém entidades puras de negócio (`tenant`, `camera`, `event`, `recording`, `zone`, `ai_rule`, `storage_pool`, `workflow`).
   - **REGRA INQUEBRÁVEL:** NUNCA deve importar bibliotecas de infraestrutura como `net/http`, `database/sql`, `github.com/jackc/pgx`, `github.com/gin-gonic`, etc.
2. **`internal/ports/` (Contratos de Portas):**
   - Define interfaces puras de entrada (Driving) e saída (Driven): repositórios, clientes de mensageria, storage MinIO e telemetria.
3. **`internal/application/` (Casos de Uso & Orquestração):**
   - Coordena o fluxo entre as portas e entidades, regras de negócio e os estágios do pipeline Medallion.
4. **`internal/adapters/` (Infraestrutura):**
   - `primary/`: Handlers HTTP REST, Middlewares de Contexto, Gateways WebSocket e WHEP.
   - `secondary/`: Implementações concretas de PostgreSQL (pgx pool), MinIO S3 SDK, NATS JetStream, NVML e NVENC.

---

## 2. Desacoplamento em 4 Planos (Split-Plane Architecture)

Para garantir escalabilidade independente de custos e resiliência:

1. **Data Plane (HydraStream Nodes):** Ingestão RTSP/RTMP, demuxing TCP, entrega WebRTC (WHEP) e buffer Zero-Copy `/dev/shm` a 8.46 GB/s. Roda em nós de CPU leves e baratos.
2. **Analytics Plane (HydraForge / Workers):** Motion-Gated Pre-Filter (VMD CPU), SAHI slicing e inferência YOLO TensorRT. Roda exclusivamente onde houver GPU (RTX 5090).
3. **Storage & Cache Tier (StorageGuard / MinIO):** Buffer de escrita NVMe Quente (*Write-Back*) com *Spillover Daemon* para HDDs mecânicos e buckets S3 distribuídos.
4. **Control Plane (HydraVMS Management):** Banco de dados multi-tenant, autenticação estrita por Token, catálogo de metadados, Workflows de Alertas e UI Vue 3.

---

## 3. Estratégia de Busca Rápida Global de Gravações (Direct-to-Source)

Para recuperar gravações históricas em `< 300ms` em qualquer lugar do mundo:
- **Índice Espaço-Temporal Centralizado:** Catálogo leve no PostgreSQL com respostas em `< 5ms` (evita `ListObjects` custoso no S3).
- **Streaming Direto Zero-Hop:** O servidor central emite *Presigned URLs* seguras apontando direto para o storage de origem da filial, sem tráfego duplo no servidor de aplicação.
- **Sprite Sheets de Timeline:** O player baixa miniaturas WebP compactas (~15 KB por minuto) para pré-visualização instantânea ao arrastar o mouse (*scrubbing*).
- **HTTP 206 Partial Content:** Suporte a *Range Requests* para iniciar o playback em `< 200ms`.

---

## 4. Hierarquia de Armazenamento & Gestão Dinâmica de Discos

- **Buffer Quente (Write-Back Cache NVMe):** Absorve rajadas de escrita de milhares de câmeras. O *Spillover Daemon* transfere blocos de 60s para os HDs em segundo plano.
- **Adição Dinâmica de Discos via Painel Web:** Cadastro de novos SSDs/HDs em tempo real na tabela `storage_pools` sem reiniciar o VMS.
- **Seletor de Finalidade de Disco:** `[HOT_VIDEO_BUFFER]` (XFS 64MB allocsize), `[EVENTS_DATABASE]` (EXT4 com PostgreSQL Tablespace), `[AI_SNAPSHOTS_CACHE]` e `[WARM_LONG_TERM]`.
- **Válvulas de Nível de Água (Watermarks):** Dreno acelerado aos 85% e gravação direta sequencial em RAM/HD aos 95% para garantir 0% de perda de frames.

---

## 5. Estratégia de Vídeo H.265 (HEVC) & Dual-Stream

- **Gravação e IA em H.265 Puro:** Vídeo original gravado sem transcode no MinIO S3 (economia de 50% de disco); decodificação direta na GPU via NVIDIA NVDEC para IA.
- **Dual-Stream Nativo:** Main Stream H.265 4K para IA/Gravação; Sub Stream H.264 para grades multi-câmera no navegador (universal e sem carga no servidor).
- **WebCodecs & JIT Transcode:** Player web usa aceleração de hardware do cliente; se o navegador for legado, ativa transcodificação Just-In-Time via NVIDIA NVENC apenas enquanto a aba estiver aberta.

---

## 6. Pipeline Medallion de IA (Bronze ➔ Silver ➔ Gold)

- **Bronze Layer (Ingestão & Telemetria Bruta):** Processa frames em `/dev/shm` e NATS com VMD Pre-Filter descartando 90% dos frames estáticos antes da GPU.
- **Silver Layer (Filtragem Espacial & SAHI):** Aplica SAHI em lotes, Non-Maximum Merging (NMM) e testes Point-in-Polygon em zonas de segurança.
- **Gold Layer (Eventos de Negócio & Despacho):** Persiste eventos acionáveis no banco, grava snapshots no MinIO e dispara workflows para Telegram, WebSockets e Webhooks.
