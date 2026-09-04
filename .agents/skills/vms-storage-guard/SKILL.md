---
name: vms-storage-guard
description: Dynamic NVMe hot buffer management, continuous spillover daemon, disk role allocation, multi-watermark safety valves, and forensic retention.
---

# HydraVMS StorageGuard & Tiered Storage Engine (vms-storage-guard)

Este skill define o funcionamento do motor de armazenamento multinível, buffer de alta velocidade em NVMe, dreno contínuo (*Spillover Daemon*) e adição dinâmica de discos no HydraVMS.

---

## 1. Arquitetura de Armazenamento Hierárquico (Tiering)

```text
[ Ingestão de 1.000+ Câmeras ]
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│ TIER QUENTE: POOL DE SSDs NVMe (/var/lib/hydravms/hot-01)   │
│ - Formatação XFS (allocsize=64M, noatime, nodiratime)       │
│ - Absorve rajadas de escrita com latência zero              │
└─────────────┬───────────────────────────────────────────────┘
              │ (Spillover Daemon: Esvaziamento Contínuo de Blocos de 60s)
              ▼
┌─────────────────────────────────────────────────────────────┐
│ TIER FRIO: POOL DE HDs MECÂNICOS / S3 (MinIO Object Storage)│
│ - Arquivo permanente de 30 a 90 dias                        │
│ - Gravação em blocos grandes sem travar o braço mecânico    │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Adição Dinâmica de Discos via Painel Web (`storage_pools`)

O administrador pode espetar novos SSDs/HDs e ativá-los no painel web sem reiniciar o VMS:

| Função Selecionada | Sistema de Arquivos | Ponto de Montagem | Uso no Sistema |
| :--- | :--- | :--- | :--- |
| **`[HOT_VIDEO_BUFFER]`** | XFS (64MB blocks) | `/var/lib/hydravms/storage/hot-XX` | Buffer de alta velocidade para gravação de vídeo |
| **`[EVENTS_DATABASE]`** | EXT4 (writeback) | `/var/lib/hydravms/storage/db-events` | PostgreSQL Tablespace e SQLite WAL |
| **`[AI_SNAPSHOTS_CACHE]`**| XFS / EXT4 | `/var/lib/hydravms/storage/snapshots` | Fotos cropadas de alertas e sprite sheets de busca |
| **`[WARM_LONG_TERM]`** | XFS | `/var/lib/hydravms/storage/warm-XX` | HDs mecânicos para retenção de 30 a 90 dias |

---

## 3. Válvulas de Nível de Água (Watermark Safety Valves)

- **< 70% de Uso do NVMe:** Dreno normal em segundo plano para o HD.
- **85% de Uso do NVMe:** Dreno acelerado multi-thread (4x mais rápido).
- **95% de Uso do NVMe (Válvula de Emergência):** O tráfego novo grava temporariamente direto no buffer sequencial de RAM/HD até o SSD esfriar (**0% de perda de frames**).
- **Proteção de Provas:** Arquivos com `is_pinned = true` são blindados contra expurgo.
