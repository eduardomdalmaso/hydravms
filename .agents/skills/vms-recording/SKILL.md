---
name: vms-recording
description: Guidelines and patterns for MinIO S3 object storage integration, multi-drive pools, recording profiles (Continuous, Motion, AI Event, Hybrid), pre-event ring buffers, auto-purge, and storage circuit breakers in HydraVMS.
---

# 📹 MinIO S3 Storage, Recording Profiles & Housekeeping Skill (HydraVMS)

This skill covers the storage pipeline, recording modes, segment naming conventions, multi-drive expansion, retention policies, and capacity watermarks.

---

## 1. Recording Profiles

1. **Continuous (24/7):** Continuous fMP4 chunks (60s to 300s) stored in `/continuous/`.
2. **Motion-Triggered (VMD):** Pixel-based motion detection with 5s pre-buffer + 10s post-buffer stored in `/motion/`.
3. **AI Event-Triggered:** YOLO inference trigger (Person, Vehicle, LPR, Face, Intrusion) with 5s RAM pre-buffer + 15-30s clip + snapshot + metadata JSON stored in `/events/`.
4. **Hybrid / Scheduled:** Time-table switching (e.g. Continuous during day hours + AI Event-only at night) and dual-streaming (5 FPS sub-stream vs 30 FPS main-stream on alarm).

---

## 2. MinIO Object Key Hierarchy

```text
hydravms-events/
└── {tenant_id}/{YYYY}/{MM}/{DD}/{camera_id}/{event_type}/{event_id}/
    ├── snapshot_full.jpg
    ├── crop_target.jpg
    ├── clip_15s.mp4
    └── metadata.json

hydravms-recordings/
└── {tenant_id}/{camera_id}/
    ├── continuous/{YYYY}/{MM}/{DD}/{YYYYMMDD_HHMMSS}_{duration_s}.mp4
    └── motion/{YYYY}/{MM}/{DD}/{YYYYMMDD_HHMMSS}_{duration_s}.mp4
```

---

## 3. Storage Watermarks & Circuit Breaker

| Disk Usage | System Action |
| :--- | :--- |
| **< 80%** | Standard operation & TTL retention purge. |
| **80% - 90%** | Warning alert on UI; start early continuous video purge (> 15 days). |
| **90% - 94.9%** | Aggressive continuous video purge; preserve motion & event clips. |
| **95% - 97.9%** | **Circuit Breaker:** Pause continuous recording; preserve Motion and AI Event recordings only. |
| **>= 98%** | **Hard Lock:** Block heavy video I/O. Snapshot & metadata only. Preserve OS/DB inodes. |
