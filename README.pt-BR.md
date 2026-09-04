# ⚡ HydraVMS — Sistema de Gerenciamento de Vídeo Corporativo de Alta Performance

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org/)
[![WebRTC](https://img.shields.io/badge/Streaming-WebRTC%20%2F%20WHEP-orange.svg)](https://webrtc.org/)
[![YOLO](https://img.shields.io/badge/Analíticos%20IA-YOLO%20%2B%20SAHI-yellow.svg)](https://github.com/ultralytics/ultralytics)

**HydraVMS** é um Sistema de Gerenciamento de Vídeo (VMS) cloud-native de última geração projetado para monitoramento ao vivo de baixíssima latência, ingestão acelerada por hardware zero-copy, reprodução síncrona multi-câmera em linha do tempo e analíticos de visão computacional em tempo real (YOLO + SAHI).

---

## 🚀 Principais Destaques & Recursos

* **⚡ Streaming de Baixíssima Latência (<300ms):** WebRTC (WHEP) nativo com fallback automático para WebSocket fMP4 para grades densas.
* **⏱️ Reprodução Síncrona Multi-Câmera:** Linha do tempo interativa em 2 níveis (régua superior multi-dias + régua inferior 24h em milissegundos) com relógio UTC sincronizado, scrubbing instantâneo, velocidades de reprodução (`0.5X`–`3X`) e recorte de intervalos com barras arrastáveis.
* **🧠 Analíticos de IA com YOLO & SAHI em Tempo Real:** Pipeline acelerado por GPU NVIDIA RTX via TensorRT e CUDA IPC para detecção e rastreamento sub-milissegundo.
* **🚨 HUD de Alarmes e Sensores Integrado:** Mapeamento em tempo real de sensores vinculados às câmeras com LEDs de 3 estados (🟢 Normal / 🟡 Alerta / 🔴 Offline) e contadores de disparos acumulados no grid.
* **💾 Armazenamento em Camadas (Tiered Storage):** Buffer NVMe de alta velocidade com transbordo automático em background para HDDs e storage S3/MinIO.
* **🗺️ Plantas Baixas e Rondas Virtuais:** Mapas 2D interativos e rondas programadas com carrossel automático.
* **⚡ Pipeline H.265 / HEVC Zero-Transcode:** Gravação direta sem re-encodificação com decodificação WebCodecs no cliente e aceleração NVENC sob demanda.
* **🛡️ Segurança Corporativa & Workflows:** Autenticação RBAC com tokens, API REST OpenAPI 3.0 e automação visual de disparos (Telegram, Webhooks e WebSockets).

---

## 📄 Licença

Distribuído sob a licença MIT.
