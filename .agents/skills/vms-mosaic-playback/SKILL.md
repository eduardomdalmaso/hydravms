---
name: vms-mosaic-playback
description: High-performance multi-view camera mosaic, dynamic sub-stream switching, instant bottom playback drawer, multi-speed playback (1x, 2x, 3x), and timeline scrubbing.
---

# HydraVMS High-Performance Mosaic & Timeline Playback (vms-mosaic-playback)

Este skill define a arquitetura para renderização fluida de grades de câmeras (1x1 a 4x4) e o mecanismo de inspeção instantânea de gravações com gaveta inferior (*Bottom Playback Drawer*).

---

## 1. Estratégias de Performance para 16+ Câmeras na Web

1. **Auto Sub-Stream Switching:**
   - Em grades multi-câmera (2x2, 3x3, 4x4), o frontend requisita o **Sub Stream leve** (H.264 640x360 @ 15 FPS), mantendo o uso de CPU/GPU do navegador abaixo de 15%.
   - Ao clicar em uma câmera ou maximizar para 1x1, o player comuta instantaneamente para o **Main Stream (Full-HD / 4K)**.
2. **Desativação Fora da Tela (IntersectionObserver):**
   - Slots de vídeo que não estão visíveis na tela têm seu canal WebRTC pausado para economizar banda e decodificação.
3. **Decodificação por Hardware no Cliente:**
   - Aproveitamento do WebCodecs para decodificar múltiplos streams simultâneos direto na GPU do computador do operador.

---

## 2. Gaveta de Gravações na Parte Inferior (Bottom Playback Drawer)

- **Abertura em 1-Clique:** Clicar em qualquer slot de câmera abre a gaveta inferior sem perder a visão do mosaico.
- **Barra Temporal com Segmentos Coloridos:**
  - Azul: Gravação Contínua.
  - Amarelo: Evento de Movimento (VMD).
  - Vermelho: Alerta Crítico de IA (Invasão de Perímetro, ALPR, Fogo).
- **Controles de Reprodução Estritamente Otimizados para Web:**
  - Play / Pause, Pulo de ±10 segundos, Passo a passo quadro a quadro (Frame Step ±1).
  - Multiplicadores de Velocidade: **`1x` (Normal), `2x` e `3x` (Avanço Rápido Estável)** — limitados estritamente até 3x para evitar engasgos de buffer e *drop frames* no navegador.
  - Botão *"Tirar Gravação"*: Fecha a gaveta e restaura o mosaico ao modo 100% Ao Vivo.
  - Botão *"Baixar Trecho (Clip)"*: Exporta o trecho selecionado em MP4 com hash forense SHA-256.
