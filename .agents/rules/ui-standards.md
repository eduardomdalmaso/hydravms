# 🎨 Padrão Visual, Design System & UI/UX das Páginas (HydraVMS)

Este documento estabelece o design system oficial, hierarquia de páginas, tipografia, cabeçalhos, barra lateral (sidebar) e sistema de alertas para a interface web em **Vue 3 + TypeScript + Vite**.

---

## 1. Paleta de Cores e Tokens Oficiais (Clean Slate Dark)

| Token CSS | Hex / Valor | Função Semântica |
| :--- | :--- | :--- |
| `--vms-bg-base` | `#090d16` / `#0f172a` | Fundo principal da aplicação (Slate 900) |
| `--vms-bg-surface` | `#1e293b` | Superfície da Sidebar, Header, Cards e Painéis |
| `--vms-bg-elevated` | `#334155` | Elementos elevados, Inputs, Dropdowns, Toolbars |
| `--vms-border` | `#334155` | Bordas sutis de containers e divisores |
| `--vms-text-main` | `#f8fafc` | Títulos, valores e textos principais (Slate 50) |
| `--vms-text-muted` | `#94a3b8` | Subtítulos, timestamps, labels secundários |
| `--vms-primary` | `#ff5e3a` / `#3b82f6` | Ações principais e botões |

---

## 2. Padrão Universal de Indicadores de Status (Neon Dot Sem Texto)

⚠️ **Regra Inviolável de Status:** Todo indicador de status em tabelas, cards e sidebars deve ser renderizado **exclusivamente como uma bolinha neon (`.vms-status-led`)** com tooltip no hover, **NUNCA com texto ao lado** (ex: proibido `[ONLINE]`, `[ATIVO]`, etc.):
- **🟢 Online / Armado:** `#00ff9d` (Verde neon com glow)
- **🟡 Desarmado:** `#fcee0a` (Amarelo cyber neon com glow)
- **🔴 Offline / Inativo:** `#ff003c` (Vermelho magenta neon com glow)
- **🔵 Gravação Ativa:** `#00f0ff` (Azul neon com pulse animation)

---

## 3. Tipografia & Hierarquia de Títulos
- **Fontes:** Títulos/Textos: `'Roboto'` / `'Inter'`, Telemetria/FPS/Timestamps: `'JetBrains Mono'`.
- **H1:** `20px` - `24px` | **H2:** `16px` - `18px` | **H3:** `14px` - `15px` | **Body:** `12px` - `13px`.

---

## 4. Padrão da Barra Lateral (Sidebar Navigation)
- **Dimensões:** Expandida `250px`, Recolhida `10px` / `64px` com toggle pill.
- **Accordion:** Cabeçalhos `[TAG]` em Roboto laranja, itens em background `#1e2126` com border-left ativa laranja.

---

## 5. Padrão de Alertas, Toasts e Modais
- **Proibição Absoluta de Janelas Nativas:** NUNCA usar `window.alert()`, `window.confirm()` ou `window.prompt()`.
- **Toast Notifications:** 4s com auto-dismiss no canto superior direito.
- **Modais (`BaseModal.vue`):** Backdrop escuro com blur e botões de ação estilizados.

---

## 6. Regra Inviolável de Modularidade Frontend (< 100 Linhas)
Nenhum arquivo `.vue`, `.ts`, `.js` ou `.css` em `web/` pode ultrapassar **100 linhas**. Componentes devem ser fracionados em subcomponentes atômicos.

---

## 7. Regra de Nomenclatura Concisa & Botão SALVAR
- **Zero Prolixidade:** Nomes atômicos e objetivos (`GRAVANDO`, `ANALITICOS`, `CODEC`, `COMPRESSAO`, `RESOLUCAO`, `FPS`).
- **Botões de Ação:** O texto de botões de confirmação é estritamente **`SALVAR`**.
- Detalhes completos em [`.agents/rules/ui-telemetry-naming.md`](./ui-telemetry-naming.md).

