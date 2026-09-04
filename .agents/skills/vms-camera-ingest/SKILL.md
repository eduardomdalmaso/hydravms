---
name: vms-camera-ingest
description: Multi-protocol ingest (RTSP, RTMP, ONVIF WS-Discovery), H.265/HEVC handling, dual-stream routing, WebCodecs client playback, and JIT NVENC hardware transcoding.
---

# HydraVMS Camera Ingestion & H.265 Handling (vms-camera-ingest)

Este skill define a ingestão de vídeo multi-protocolo, o aproveitamento de fluxos duplos (Dual-Stream) e o pipeline de entrega inteligente para codecs H.264 e H.265 (HEVC).

---

## 1. Ingestão Multi-Protocolo & Descoberta

1. **ONVIF WS-Discovery:** Probe multicast UDP (`239.255.255.250:3702`) para auto-descoberta de câmeras IP na rede local em < 2 segundos.
2. **RTSP Demuxer (RFC 2326):** Conexões interleaved sobre TCP para evitar perda de pacotes e corrupção de frames.
3. **Servidor RTMP Push Embutido:** Porta 1935 para receber streams push de drones, bodycams e apps mobile.

---

## 2. Estratégia de Vídeo H.265 (HEVC) & Dual-Stream

```text
               [ Câmera IP (Main: H.265 4K / Sub: H.264) ]
                                   │
                                   ▼
┌────────────────────────────────────────────────────────────────────────┐
│ 1. GRAVAÇÃO & IA: H.265 PURO (Zero Transcode = Economia de 50% de Disco)│
│ - Grava o H.265 original no MinIO S3 (Zero-Reencode Stream-Copy)       │
│ - O Analítico (HydraForge) decodifica H.265 direto na GPU via NVDEC    │
└──────────────────────────────────┬─────────────────────────────────────┘
                                   │
                                   ▼
┌────────────────────────────────────────────────────────────────────────┐
│ 2. ENTREGA WEB INTELIGENTE (Client-Side Detection)                     │
│ ├─► Navegador Moderno (Chrome/Edge/Safari/Mobile): WebRTC H.265 nativo │
│ ├─► Grade Multi-Câmera: Exibe Sub-Stream H.264 leve e universal        │
│ └─► Navegador Legado (sem H.265): JIT Transcode com NVIDIA NVENC       │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. WebCodecs & Transcodificação Sob Demanda (JIT NVENC)

- **Client-Side Hardware Decoding:** O frontend Vue 3 consulta `VideoDecoder.isConfigSupported()` para tocar H.265 direto via GPU do cliente sem custo para o servidor.
- **JIT Transcode com NVENC:** Se um navegador legado solicitar uma câmera H.265 em tela cheia, o HydraStream ativa uma sessão de transcode transitória via hardware NVIDIA NVENC a 800+ FPS, desativando-a assim que a aba do navegador for fechada.
