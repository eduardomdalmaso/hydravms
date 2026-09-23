# ⚡ HydraVMS — Sistema de Gerenciamento de Vídeo de Alta Performance

[![Licença](https://img.shields.io/badge/licenca-MIT-blue.svg)](LICENSE)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![SQLite WAL](https://img.shields.io/badge/SQLite-WAL%20Zero--Config-003B57?logo=sqlite&logoColor=white)](https://sqlite.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16%20Enterprise-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![NATS](https://img.shields.io/badge/NATS-JetStream-27AAE1?logo=nats.io&logoColor=white)](https://nats.io/)
[![MinIO](https://img.shields.io/badge/MinIO-S3%20Object%20Storage-C72C48?logo=minio&logoColor=white)](https://min.io/)
[![WebRTC](https://img.shields.io/badge/Streaming-WebRTC%20%2F%20WHEP-orange.svg)](https://webrtc.org/)
[![YOLO](https://img.shields.io/badge/IA%20Analiticos-YOLO%20%2B%20SAHI-yellow.svg)](https://github.com/ultralytics/ultralytics)

[**English**](README.md) | [**Português do Brasil**]

O **HydraVMS** é um Sistema de Gerenciamento de Vídeo (VMS) cloud-native de última geração, projetado para monitoramento ao vivo de ultra-baixa latência, ingestão Zero-Copy acelerada por hardware, timeline sincronizada multi-câmera, armazenamento hierárquico MinIO S3 e analíticos de visão computacional em tempo real (YOLO + SAHI na NVIDIA RTX 5090).

---

## 💾 Motor de Banco de Dados Relacional (Arquitetura Dual-Engine)

O HydraVMS possui uma camada agnóstica em Arquitetura Hexagonal suportando dois modos de execução nativos:

1. **SQLite em Modo WAL (`hydravms.db` — Padrão Zero-Config):**
   * **Zero instalação externa:** Motor SQLite em Go puro compilado diretamente dentro do executável.
   * **Inicialização Instantânea:** Cria tabelas, índices e o usuário administrador padrão (`admin` / `admin`) automaticamente no primeiro boot.
   * **Alta Performance:** O modo *Write-Ahead Logging* (WAL) garante leituras simultâneas sem bloqueio e alta vazão de escrita para NVRs Edge e estações locais de trabalho.
2. **PostgreSQL 16+ (Modo Corporativo em Cluster):**
   * Ativado automaticamente ao definir a variável `DATABASE_URL` ou `DB_DRIVER=postgres`.
   * Isolamento multi-inquilino com *Row-Level Security* (RLS), particionamento de tabelas para bilhões de eventos e escalabilidade distribuída.

---

## 🌐 Mapa de Portas do Ecossistema

| Serviço / Container | Porta(s) | Protocolo / Descrição |
| :--- | :--- | :--- |
| **`hydra-vms` (Frontend Web)** | `5173` | HUD Cyberpunk em Vue 3 + Vite |
| **`hydra-vms-api` (Control Plane)** | `8083` | API REST em Go & Gateway WebSocket em tempo real |
| **`hydravms.db` (Embutido Padrão)** | — | Banco Relacional SQLite WAL Zero-Config |
| **`hydra_postgres` (Podman Opcional)**| `5432` | Banco Relacional PostgreSQL 16 Corporativo |
| **`hydra_nats` (Podman / Standalone)**| `4222`, `8222` | Barramento de Eventos e Telemetria NATS JetStream |
| **`hydra_minio` (Podman / S3)** | `9000`, `9001` | S3 API (:9000) & Console Web (:9001) |
| **`hydra-stream` (Motor Ingestão)** | `8080` | Multiplexador de Vídeo Zero-Copy (/dev/shm) |
| **MediaMTX (RTSP / WebRTC)** | `8554`, `8889` | Ingestão RTSP (:8554) & WebRTC WHEP (:8889) |
| **`hydra-forge` (Estúdio IA)** | `8081` | Estúdio de Treinamento YOLO & Compilador TensorRT |
| **`hydra-vault` (Cofre de Dados)** | `8082` | Curadoria Ativa de Datasets & Auto-Rotulagem SAM 2 |

---

## 🚀 Guia de Inicialização Rápida

### Modo A: Modo Standalone Zero-Dependências (Padrão com SQLite WAL)

Você **NÃO** precisa instalar PostgreSQL ou Docker para rodar e testar o HydraVMS localmente:

```bash
# 1. Iniciar o Backend Go diretamente (cria o hydravms.db automaticamente)
go run ./cmd/hydravms

# 2. Iniciar o Frontend em outro terminal
cd web
npm install
npm run dev
```

Acesse `http://localhost:5173` (ou `http://localhost:8083`) e faça login com:
* **Usuário / E-mail:** `admin` ou `admin@hydravms.io`
* **Senha:** `admin`

---

### Modo B: Stack Corporativa em Containers (PostgreSQL 16 + NATS + MinIO)

Para clusters de produção com serviços dedicados:

```bash
# 1. PostgreSQL 16
podman run -d --name hydra_postgres --restart always \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=hydravms \
  -p 5432:5432 \
  docker.io/library/postgres:16-alpine

# 2. NATS JetStream
podman run -d --name hydra_nats --restart always \
  -p 4222:4222 -p 8222:8222 \
  docker.io/library/nats:latest -js -m 8222

# 3. MinIO S3 Object Storage
podman run -d --name hydra_minio --restart always \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -p 9000:9000 -p 9001:9001 \
  -v hydra_minio_data:/data \
  docker.io/minio/minio:latest server /data --console-address ":9001"
```

### Passo 2: Executar Migrações e Seed Inicial

```bash
cd /home/hades/Documents/HydraVMS

# Aplicar migrações de schema v1 a v5
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000001_initial_schema.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000002_folders_layouts_maps_tours.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000003_cluster_nodes_and_ptz.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000004_performance_indexes_and_jsonb_tuning.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000005_storage_pools_enhancement.up.sql

# Inserir dados de teste e estrutura inicial
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/seed.sql
```

### Passo 3: Compilar e Executar via PM2

```bash
# Compilar Backend em Go
cd /home/hades/Documents/HydraVMS
go build -o bin/hydravms ./cmd/hydravms

# Compilar Frontend em Vue 3
cd web
npm install
npm run build

# Reiniciar todos os serviços com PM2
pm2 restart all
```

---

## 🏛️ Destaques da Arquitetura

1. **Arquitetura Hexagonal (Ports & Adapters):** O pacote `internal/domain` tem **0 imports** de infraestrutura ou banco de dados.
2. **Multi-Tenancy com RLS:** Isolamento direto no kernel do PostgreSQL através de `ROW LEVEL SECURITY`.
3. **Armazenamento em Camadas (Tiering):** Buffer Rápido NVMe/RAM (`/dev/shm`) ➔ MinIO S3 de Longo Prazo.
4. **Reprodução Zero-Hop de H.265:** WebCodecs API no navegador + JIT NVENC na RTX 5090 quando necessário.
5. **Barramento NATS JetStream:** Notificações em tempo real conectadas com WebSockets e Bots do Telegram.

---

## 📄 Licença
Distribuído sob a licença MIT.
