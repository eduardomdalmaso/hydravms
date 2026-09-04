# 🎨 Padrão Visual, Design System & UI/UX das Páginas (HydraVMS)

Este documento estabelece o design system oficial, hierarquia de páginas, tipografia, cabeçalhos, barra lateral (sidebar) e sistema de alertas para a interface web em **Vue 3 + TypeScript + Vite**.

---

## 1. Paleta de Cores e Tokens Oficiais (Clean Slate Dark)

A interface deve adotar um tema escuro corporativo, moderno e de alto contraste funcional:

| Token CSS | Hex / Valor | Função Semântica |
| :--- | :--- | :--- |
| `--vms-bg-base` | `#090d16` ou `#0f172a` | Fundo principal da aplicação (Slate 900) |
| `--vms-bg-surface` | `#1e293b` | Superfície da Sidebar, Header, Cards e Painéis (Slate 800) |
| `--vms-bg-elevated` | `#334155` | Elementos elevados, Inputs, Dropdowns, Toolbars (Slate 700) |
| `--vms-border` | `#334155` | Bordas sutis de containers e divisores |
| `--vms-border-focus` | `#3b82f6` | Anel de foco e seleção ativa de câmera |
| `--vms-text-main` | `#f8fafc` | Títulos, valores e textos principais (Slate 50) |
| `--vms-text-muted` | `#94a3b8` | Subtítulos, timestamps, labels secundários (Slate 400) |
| `--vms-primary` | `#3b82f6` | Ações principais, botões de ação e abas ativas (Blue 500) |
| `--vms-success` | `#10b981` | Status Online, gravação ativa e confirmações (Emerald 500) |
| `--vms-warning` | `#f59e0b` | Alertas de atenção, uso moderado de disco (Amber 500) |
| `--vms-danger` | `#ef4444` | Violação de zonas, erros críticos e ações destrutivas (Red 500) |

---

## 2. Tipografia & Hierarquia de Títulos

A tipografia deve ser limpa, sem serifa e com fonte monospace para telemetria numérica:

- **Famílias de Fontes:**
  - Títulos e Textos Gerais: `'Inter'`, `'Geist'`, `'system-ui'`, `-apple-system`, `sans-serif`.
  - Telemetria, FPS, Coordenadas BBox e Timestamps: `'JetBrains Mono'`, `'Fira Code'`, `monospace`.
- **Hierarquia de Títulos (Headings):**
  - **H1 (Título da Página):** `font-size: 20px` a `24px`, `font-weight: 600`, cor `--vms-text-main`, acompanhado de breadcrumb ou badge de status.
  - **H2 (Título de Seção / Painel):** `font-size: 16px` a `18px`, `font-weight: 600`, cor `#e2e8f0` (Slate 200).
  - **H3 (Título de Card / Modal):** `font-size: 14px` a `15px`, `font-weight: 500`, cor `#cbd5e1` (Slate 300).
  - **Body / Textos Gerais:** `font-size: 13px` a `14px`, `line-height: 1.5`, cor `--vms-text-main`.
  - **Labels / Badges / Metadados:** `font-size: 11px` a `12px`, `font-weight: 500`, cor `--vms-text-muted`.

---

## 3. Padrão do Cabeçalho Superior (Header / Navbar)

- **Dimensões:** Altura fixa de `56px` a `60px`, `width: 100%`, `position: sticky; top: 0; z-index: 40;`.
- **Estilo:** Background `--vms-bg-surface`, borda inferior de `1px solid var(--vms-border)`.
- **Composição Visual (Grid em 3 Zonas):**
  - **Esquerda:** Logo minimalista do HydraVMS + Breadcrumb da página ativa (ex: `Câmeras / Portaria Principal`).
  - **Centro:** Mini-widget de Telemetria de Sistema em tempo real (CPU %, GPU %, Storage Guard Status `[ 47% OK ]`).
  - **Direita:** Seletor de Tenant Ativo (para Master Admin) + Botão de Notificações com sino e badge numérico + Avatar / Menu de Perfil do Usuário com Logout.

---

## 4. Padrão da Barra Lateral (Sidebar Navigation)

- **Dimensões:**
  - Expandida: Largura fixa de `220px` a `240px`.
  - Recolhida (Modo Compacto): Largura fixa de `64px` (apenas ícones centralizados com tooltips).
- **Estilo:** Background `--vms-bg-surface`, borda direita de `1px solid var(--vms-border)`.
- **Itens de Menu Padronizados:**
  1. `[ Live Monitor ]` (Grade de câmeras multi-view 1x1, 2x2, 3x3, 4x4)
  2. `[ Events & AI Alerts ]` (Feed em tempo real, filtros, visualizador de clips)
  3. `[ Recordings & Playback ]` (Timeline contínua/movimento e busca rápida)
  4. `[ Cameras & Streams ]` (Cadastro, status RTSP e resolução)
  5. `[ AI Rules & Zones ]` (Desenhador de polígonos e configuração SAHI)
  6. `[ Storage & Disks ]` (Pools de HDs, políticas de retenção e logs de expurgo)
  7. `[ Settings & Users ]` (RBAC, tenants e preferências)
- **Indicador de Item Ativo:**
  - Borda esquerda de `3px solid var(--vms-primary)`.
  - Fundo sutilmente elevado (`rgba(59, 130, 246, 0.1)`).
  - Texto e ícone iluminados na cor branca (`#ffffff`).

---

## 5. Padrão de Alertas, Toasts e Modais

### ⚠️ Proibição Absoluta de Janelas Nativas do Navegador:
- **NUNCA** utilizar `window.alert()`, `window.confirm()` ou `window.prompt()`. Todos os diálogos devem ser componentes Vue 3 estilizados.

### Componentes de Notificação e Diálogo:
1. **Toast Notifications (Flutuantes no Canto Superior Direito):**
   - Duração: 4 segundos (com barra de progresso suave de auto-dismiss).
   - Variantes:
     - `Info`: Borda azul (`#3b82f6`), ícone de informação.
     - `Success`: Borda verde (`#10b981`), ícone de check.
     - `Warning`: Borda amarela (`#f59e0b`), ícone de exclamação.
     - `Danger`: Borda vermelha (`#ef4444`), ícone de alerta crítico.
2. **Live AI Alarm Banner / Toast:**
   - Exibição instantânea no topo do painel quando um evento da camada Gold é recebido via WebSocket.
   - Contém thumbnail do snapshot cropado, nome da câmera, classe detectada e botão *"Visualizar Câmera"*.
3. **Modais de Diálogo (`BaseModal.vue`):**
   - Backdrop semi-transparente escuro (`rgba(9, 13, 22, 0.75)` com `backdrop-filter: blur(4px)`).
   - Caixa centralizada com bordas arredondadas (`border-radius: 8px`), cabeçalho com título claro, corpo do formulário e rodapé com botões *"Cancelar"* (neutro) e *"Confirmar / Salvar"* (primário/perigo).

---

## 6. Regra Inviolável de Modularidade Frontend (< 100 Linhas)

- **Limite Estrito:** Nenhum arquivo `.vue`, `.ts`, `.js` ou `.css` na pasta `web/` pode ultrapassar **100 linhas**.
- **Técnicas de Decomposição:**
  - Extração de lógica para Composables (`useCameraStream.ts`, `useTelemetry.ts`, `useEvents.ts`).
  - Divisão de componentes em subcomponentes atômicos (`CameraHeader.vue`, `CameraPlayer.vue`, `CameraFooter.vue`).
  - Separação de CSS em arquivos modulares na pasta `assets/css/`.
