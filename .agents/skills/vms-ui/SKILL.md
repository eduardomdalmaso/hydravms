---
name: vms-ui
description: Design system guidelines, typography, modular CSS architecture (<100 lines/file), responsive layouts, and UI component standards for HydraVMS.
---

# HydraVMS Frontend & UI Design System (vms-ui)

Este skill define a arquitetura de estilos compartilhados, hierarquia visual, tipografia e responsividade para todas as telas do HydraVMS.

---

## 1. Arquitetura Modular de CSS (`web/src/assets/css/`)

Todos os elementos da interface compartilham a mesma biblioteca modular de CSS, onde **nenhum arquivo ultrapassa 100 linhas**:

| Arquivo CSS | Escopo & Responsabilidade |
| :--- | :--- |
| **`variables.css`** | Paleta Clean Slate Dark (`#090d16`, `#0f172a`, `#1e293b`), tokens semânticos, sombras e z-index. |
| **`typography.css`** | Fontes Google (`Inter` e `JetBrains Mono`), escala de tamanhos (`11px` a `24px`), títulos H1-H4. |
| **`reset.css`** | Reset moderno, `box-sizing: border-box` e scrollbars minimalistas escuras. |
| **`layout.css`** | Grid principal da aplicação (`vms-app-layout`, `vms-content-area`, containers flexíveis). |
| **`sidebar.css`** | Barra de navegação lateral (estados 240px expandido / 64px compacto, links ativos). |
| **`header.css`** | Cabeçalho superior (56px, breadcrumbs, badges de telemetria e perfil do usuário). |
| **`buttons.css`** | Botões padronizados (`primary`, `secondary`, `danger`, `ghost`, `icon`, `sm`, `lg`). |
| **`forms.css`** | Inputs, selects (sem emojis), checkboxes, toggles e validações. |
| **`cards.css`** | Cards de superfície, containers elevados e slots HUD de vídeo de câmeras (16:9). |
| **`badges.css`** | Badges de status (`online`, `offline`, `recording` pulsante, `warning`, `info`). |
| **`tables.css`** | Tabelas de dados com scroll horizontal e paginação. |
| **`modals.css`** | Diálogos modais com backdrop blur escuro e animações suaves de entrada. |
| **`toasts.css`** | Notificações flutuantes e banners de alerta em tempo real. |
| **`responsive.css`** | Breakpoints responsivos para auto-ajuste em telas pequenas, tablets e mobile. |
| **`neumorphism.css`** | Estilos Neumorphism Dark (Soft UI, botões circulares 3D, sombras convex/concave e destaque laranja). |
| **`main.css`** | Agregador mestre importando todos os módulos em ordem estrita de cascata. |

---

## 2. Padrão de Auto-Ajuste & Responsividade (Small Screens)

- **Desktop (> 1024px):** Layout completo com Sidebar fixa de 240px e telemetria no Header.
- **Tablets / Telas Médias (<= 1024px):** Sidebar compacta de 200px, telemetria recolhida para dropdown.
- **Mobile / Telas Pequenas (<= 768px):** 
  - Sidebar vira gaveta (drawer) off-canvas flutuante com backdrop escuro.
  - Cabeçalhos de página empilham verticalmente.
  - Modais ocupam largura total com margem de segurança de 8px.
  - Tabelas ganham scroll horizontal com cabeçalho sticky.
- **Telas Ultra-Pequenas (<= 480px):**
  - Header reduz para 50px de altura.
  - Padding do container reduzido para 8px.
  - Botões ajustados para toque tátil com fonte de 12px.

---

## 3. Checklist de Auditoria Frontend

```bash
# Auditoria de Linhas (Regra estrita de < 100 linhas por arquivo)
wc -l web/src/assets/css/*.css | awk '$1 > 100 { print "VIOLATION: " $2 " has " $1 " lines" }'
```
