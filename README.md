# ⚡ HydraVMS — Enterprise High-Performance Video Management System

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![NATS](https://img.shields.io/badge/NATS-JetStream-27AAE1?logo=nats.io&logoColor=white)](https://nats.io/)
[![MinIO](https://img.shields.io/badge/MinIO-S3%20Object%20Storage-C72C48?logo=minio&logoColor=white)](https://min.io/)
[![WebRTC](https://img.shields.io/badge/Streaming-WebRTC%20%2F%20WHEP-orange.svg)](https://webrtc.org/)
[![YOLO](https://img.shields.io/badge/AI%20Analytics-YOLO%20%2B%20SAHI-yellow.svg)](https://github.com/ultralytics/ultralytics)

[**English**] | [**Português do Brasil**](README.pt-BR.md)

**HydraVMS** is a next-generation, cloud-native Video Management System (VMS) engineered for ultra-low-latency live monitoring, zero-copy hardware-accelerated ingest, synchronized multi-camera timeline playback, tiered S3/MinIO storage, and real-time AI computer vision analytics (YOLO + SAHI).

---

## 🌐 Ecosystem Port Map

| Service / Container | Port(s) | Protocol / Description |
| :--- | :--- | :--- |
| **`hydra-vms` (Frontend Web)** | `5173` | Vue 3 + Vite Cyberpunk High-Tech HUD |
| **`hydra-vms-api` (Control Plane)** | `8083` | Go REST API & WebSocket Real-time Gateway |
| **`hydra_postgres` (Podman)** | `5432` | PostgreSQL 16 Relational DB (RLS Multi-Tenancy) |
| **`hydra_nats` (Podman)** | `4222`, `8222` | NATS JetStream Event Mesh & Telemetry |
| **`hydra_minio` (Podman)** | `9000`, `9001` | S3 API (:9000) & Web Console (:9001) |
| **`hydra-stream` (Ingest Engine)** | `8080` | Zero-Copy SHM Video Multiplexer |
| **MediaMTX (RTSP / WebRTC)** | `8554`, `8889` | RTSP Ingest (:8554) & WebRTC WHEP (:8889) |
| **`hydra-forge` (AI Studio)** | `8081` | YOLO Training Studio & TensorRT Compiler |
| **`hydra-vault` (Curator)** | `8082` | Dataset Active Learning & SAM 2 Auto-Labeling |

---

## 🚀 Quickstart & Bootstrapping Guide

### Step 1: Start Infrastructure Containers (Podman / Docker)

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

### Step 2: Run Database Migrations & Initial Seed

```bash
cd /home/hades/Documents/HydraVMS

# Apply schema migrations v1 to v5
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000001_initial_schema.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000002_folders_layouts_maps_tours.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000003_cluster_nodes_and_ptz.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000004_performance_indexes_and_jsonb_tuning.up.sql
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/000005_storage_pools_enhancement.up.sql

# Insert initial seed data (tenant, users, folders, cameras, storage pools)
podman exec -i hydra_postgres psql -U postgres -d hydravms < migrations/seed.sql
```

### Step 3: Build and Run Backend & Frontend via PM2

```bash
# Build Go Backend
cd /home/hades/Documents/HydraVMS
go build -o bin/hydravms ./cmd/hydravms

# Build Vue 3 Frontend
cd web
npm install
npm run build

# Start / Restart all services with PM2
pm2 restart all # or pm2 start ecosystem.config.js
```

---

## 🏛️ Architecture Highlights

1. **Hexagonal Architecture (Ports & Adapters):** Package `internal/domain` has **0 imports** of database or HTTP infrastructure.
2. **Multi-Tenancy with RLS:** Native PostgreSQL `ROW LEVEL SECURITY` guarantees zero cross-tenant leakage.
3. **High-Throughput Storage Tiering:** NVMe Hot Buffer (`/dev/shm`) ➔ MinIO S3 Long-term Storage.
4. **Zero-Hop WebRTC & H.265 Playback:** WebCodecs API in browser + JIT NVENC hardware transcoding on NVIDIA RTX 5090.
5. **Event Mesh:** NATS JetStream pub/sub bridge to real-time WebSockets and Telegram bots.

---

## 📄 License
Licensed under the MIT License.
