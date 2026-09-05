# 🏷️ Padrão de Nomenclatura, Rótulos e Telemetria UI (HydraVMS)

Este documento estabelece a regra oficial de concisão para rótulos de telemetria, títulos de cards, botões de ação e campos da interface em todo o ecossistema **HydraVMS**.

---

## 1. Regra de Ouro: Nomenclatura Atômica e Concisa (Zero Prolixidade)

Rótulos de dados, telemetria e cards devem ser **estritamente atômicos e diretos**, eliminando palavras desnecessárias, adjetivos ou particípios redundantes.

### Dicionário Oficial de Rótulos de Telemetria:
| ❌ Proibido (Prolixo / Composto) | ✅ Obrigatório (Conciso / HUD) | Descrição / Exemplo de Valor |
| :--- | :--- | :--- |
| `GRAVANDO ATIVAMENTE` | `GRAVANDO` | `SIM (HABILITADO)` / `NAO (DESLIGADO)` |
| `ANALITICOS VINCULADOS` | `ANALITICOS` | `2` |
| `ALARMES & SENSORES VINCULADOS` | `ALARMES` ou `ALARMES & SENSORES` | `1` |
| `CODEC & COMPRESSAO` | `CODEC` e `COMPRESSAO` (Separados) | `H.265` e `HEVC` |
| `RESOLUCAO NATIVA` | `RESOLUCAO` e `FPS` (Separados) | `1080P` e `30 FPS` |
| `TAXA DE FLUXO / BITRATE` | `BITRATE` ou `TAXA DE FLUXO` | `4.0 Mbps` |

---

## 2. Padrão de Nomenclatura de Pastas e Itens (Caixa Alta / Capitalização)

- **Nome de Pastas / Zonas:** **SEMPRE EM CAIXA ALTA (UPPERCASE)**.
  - Exemplo: `PORTARIA & ACESSOS PRINCIPAIS`, `GALPAO DE CARGAS & DOCAS`, `DATACENTER & CPD`, `PERIMETRO EXTERNO`.
- **Nome de Fluxo e Alarme:** **SEMPRE PRIMEIRA LETRA MAIÚSCULA (Title Case / Capitalize)**.
  - Exemplo: `Portaria Principal (Entrada)`, `Estacionamento Visitantes`, `Sensor Muro Norte`, `Barreira Leste IVS`.

---

## 3. Padrão de Botões de Confirmação e Salvamento

O texto dos botões de salvamento em formulários, modais, painéis e assistentes deve ser **estritamente `SALVAR`**, nunca adicionando o nome da entidade ou ações compostas:
- ❌ Proibido: `SALVAR HORARIOS`, `SALVAR CAMERA`, `SALVAR ALARME`, `SALVAR RONDA`, `SALVAR PLANTA`, `SALVAR E ATIVAR FLUXO`, `SALVAR ALTERACOES`.
- ✅ Obrigatório: **`SALVAR`**.

---

## 4. Ícones de Ação Padronizados (Laranja `#ff5e3a`)

- **Edição:** Flaticon `2355330` (Bloco de notas com lápis em `#ff5e3a`).
- **Exclusão:** Flaticon `6997199` (Lixeira clássica com tampa e nervuras em `#ff5e3a`).
- **Dimensão:** Bounding box fixo de `28px x 28px` com SVG de `16px x 16px`.

---

## 5. Paleta de Cores e Tipografia
- **3 Cores Primárias:** Branco (informações/textos), Laranja `#ff5e3a` (ações/HUD/alertas) e Preto/Escuro `#07080c` / `#14171d` (superfícies/fundos).
- **2 Fontes Oficiais:** `'Roboto'` (interface) e `'JetBrains Mono'` (métricas, códigos, IDs e telemetria).
