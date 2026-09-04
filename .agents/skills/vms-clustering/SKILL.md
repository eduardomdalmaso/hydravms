---
name: vms-clustering
description: Distributed multi-node clustering, horizontal scale-out between servers, inter-node communication protocols (gRPC, NATS JetStream Mesh, WebRTC), node discovery, and load balancing for HydraVMS.
---

# 🌐 Distributed Multi-Node Clustering & Scaling Skill (HydraVMS)

This skill provides architectural patterns, protocols, and scaling playbooks for distributing HydraVMS across multiple physical servers, GPU clusters, and edge ingestion nodes.

---

## 1. Cluster Node Roles & Topology

```text
                                [ CLIENT TRAFFIC / WEB UI ]
                                             │
                                             ▼ (HTTPS / WSS)
                 ┌────────────────────────────────────────────────────────┐
                 │       CONTROL PLANE NODES (HydraVMS Core / API)        │
                 │   • REST API Gateway & Swagger (JWT Auth, Multi-Tenant)│
                 │   • WebSocket Gateway (Live Alarms & Telemetry Push)   │
                 │   • PostgreSQL Coordinator (Metadata, Rules, RBAC)     │
                 └───────────────────────────┬────────────────────────────┘
                                             │
                 ┌───────────────────────────┴────────────────────────────┐
                 │     INTER-NODE MESSAGING MESH (NATS JetStream Cluster) │
                 │      (Topics: `cluster.nodes.*`, `cameras.*.frames`)   │
                 └───────┬───────────────────┬────────────────────┬───────┘
                         │                   │                    │
                         ▼                   ▼                    ▼
      ┌─────────────────────────┐ ┌─────────────────────┐ ┌───────────────┐
      │   EDGE INGEST NODES     │ │  GPU WORKER NODES   │ │ MINIO S3 POOL │
      │   (HydraStream Ingest)  │ │ (YOLO TensorRT/SAHI)│ │ (Distributed) │
      │ • RTSP/H.264 Ingestion  │ │ • Motion-gated infer│ │ • Multi-drive │
      │ • Motion-Gating Filter  │ │ • RTX 5090 Workers  │ │   Erasure Code│
      │ • WebRTC WHIP/WHEP Gate │ │ • Batch Slice NMM   │ │ • Presigned   │
      │ • MinIO fMP4 Chunk Rec  │ │ • 18,500+ FPS Power │ │   Video URLs  │
      └─────────────────────────┘ └─────────────────────┘ └───────────────┘
```

---

## 2. Protocolos de Comunicação Inter-Nodos

| Propósito | Protocolo Selecionado | Justificativa Técnica |
| :--- | :--- | :--- |
| **Transmissão de Frames Local (Same Host)** | **POSIX `/dev/shm` + CUDA IPC** | Latência sub-microssegundo, Zero-Copy a **8.46 GB/s**. |
| **Streaming de Frames Distribuído (Rede)** | **gRPC (HTTP/2 + Protobuf)** | Compressão binária ultrarrápida com streaming bidirecional entre nós de ingestão e nós GPU. |
| **Barramento de Eventos & Telemetria** | **NATS JetStream Mesh** | Cluster nativo em Go, failover automático, suporte a > 10M msgs/s com roteamento por assunto (`cameras.<tenant>.<cam_id>.events`). |
| **Descoberta & Heartbeat de Nós** | **NATS Key-Value (NATS KV)** | Cada nó registra um TTL de 5s (`cluster.nodes.<node_id>`). Se o nó travar, o cluster detecta em 5s e redistribui as câmeras. |
| **Vídeo ao Vivo para Clientes Web** | **WebRTC (WHIP / WHEP via Pion Go)** | Latência ultrabaixa (< 200ms) direto para o player Vue 3 no navegador. |

---

## 3. Playbook de Adição Rápida de Servidores (Scale-Out em 60s)

### A. Adicionar um Novo Nó de Ingestão (Edge Ingest Node):
1. Inicie o binário `hydrastream` com `--cluster-nats-url="nats://cluster-core:4222"` e `--node-role=ingest`.
2. O nó se auto-registra no NATS KV.
3. O Control Plane aloca automaticamente o próximo lote de câmeras RTSP para esse nó.

### B. Adicionar um Novo Servidor GPU (RTX 5090 Worker Node):
1. Inicie o worker com `--cluster-nats-url="nats://cluster-core:4222"` e `--node-role=gpu_worker`.
2. O worker subscreve ao pool de tarefas de inferência (`queue: inference-workers`).
3. O NATS balanceia a carga de fatias SAHI e frames com movimento entre todas as GPUs do cluster.
