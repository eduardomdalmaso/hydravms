# ⚡ HydraVMS — Enterprise High-Performance Video Management System

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org/)
[![WebRTC](https://img.shields.io/badge/Streaming-WebRTC%20%2F%20WHEP-orange.svg)](https://webrtc.org/)
[![YOLO](https://img.shields.io/badge/AI%20Analytics-YOLO%20%2B%20SAHI-yellow.svg)](https://github.com/ultralytics/ultralytics)

**HydraVMS** is a next-generation, cloud-native Video Management System (VMS) engineered for ultra-low-latency live monitoring, zero-copy hardware-accelerated ingest, synchronized multi-camera timeline playback, and real-time AI computer vision analytics (YOLO + SAHI).

Built with a high-density, dark neumorphic Cyberpunk HUD interface, HydraVMS provides security operators and enterprise facilities with unmatched responsiveness, instant forensic playback, and automated workflow triggers.

---

## 🚀 Key Highlights & Capabilities

* **⚡ Ultra-Low Latency Streaming (<300ms):** Pure WebRTC (WHEP) live streaming with automatic WebSocket fMP4 fallback for seamless multi-grid monitoring.
* **⏱️ Synchronized Multi-Camera Playback:** Two-tier interactive canvas timeline (top horizontal multi-day strip + bottom 24h millisecond ruler) with synchronized UTC clock, instant frame scrubbing, variable speed playback (`0.5X`–`3X`), and draggable video interval clipping.
* **🧠 Real-Time YOLO & SAHI AI Analytics:** Native hardware-accelerated pipeline running on NVIDIA RTX GPUs via TensorRT and CUDA IPC, enabling sub-millisecond detection, object tracking, and instant alarm triggers.
* **🚨 Integrated Alarm & Sensor HUD:** Real-time sensor mapping on camera streams with 3-state LED telemetry (🟢 Normal / 🟡 Alert / 🔴 Offline) and accumulated trigger counters directly in the grid.
* **💾 Multi-Tier Storage Architecture:** High-speed NVMe write-back buffer with automated background tiered spillover to HDDs and S3/MinIO cloud object storage.
* **🗺️ Interactive Maps & Virtual Guard Patrols:** Low-latency 2D facility floorplans with dynamic camera reticles and programmable automated carousel patrols.
* **⚡ Zero-Transcode H.265 / HEVC Pipeline:** Pure direct stream storage with client-side WebCodecs decoding and on-demand JIT NVENC acceleration.
* **🛡️ Enterprise Security & Workflows:** Token-enforced RBAC, OpenAPI 3.0 REST API, and visual trigger automation dispatching instant Telegram, Webhook, and WebSocket alerts.

---

## 🏗️ Architecture Ecosystem

```
[ IP Cameras / RTSP Streams ]
          │
          ▼
[ HydraStream Data Plane ] ──(Zero-Copy /dev/shm)──► [ HydraForge AI Studio ]
          │                                                   │
          ▼                                                   ▼
[ Tiered NVMe / HDD / S3 ]                             [ YOLO TensorRT Engine ]
          │                                                   │
          └─────────────────────► [ HydraVMS ] ◄──────────────┘
                                  (Control Plane & Web HUD)
```

---

## 💻 Tech Stack

* **Frontend:** Vue 3, Vite, TypeScript, Canvas 2D Timeline Engine, Dark Neumorphic Design System (Modular architecture, <100 lines per component).
* **Backend & API:** Go (Hexagonal Architecture / Ports & Adapters), OpenAPI 3.0, SQLite WAL / PostgreSQL.
* **Media & Ingest:** HydraStream Data Plane, WebRTC WHEP, RTSP/ONVIF/RTMP demuxing, WebCodecs.
* **AI & Acceleration:** PyTorch, Ultralytics YOLO, TensorRT, NVIDIA CUDA 13.x / RTX 5090.

---

## 📦 Getting Started

### Prerequisites
* Node.js 20+ & npm
* Go 1.22+ (for control plane API)

### Running the Web HUD
```bash
cd web
npm install
npm run dev
```

### Production Build
```bash
cd web
npm run build
```

---

## 📄 License

This project is licensed under the MIT License.
