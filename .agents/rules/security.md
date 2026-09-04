# 🔒 Regras Estritas de Segurança (HydraVMS)

Este documento estabelece as diretrizes invioláveis de segurança, autenticação obrigatória via token, autorização multi-tenant e proteção de mídia para o HydraVMS.

---

## 1. Acesso à API Exclusivamente via Token (Token-Only Enforcement)

- **Proibição de Rotas Públicas Anônimas:** Todos os endpoints REST (`/api/v1/*`), feeds WebRTC e conexões WebSocket (`/ws/v1/*`) EXIGEM autenticação obrigatória via Token.
- **Única Exceção:** A rota `POST /api/v1/auth/login` (utilizada para emissão do primeiro token).
- **Tipos de Tokens Suportados:**
  1. **User JWT Bearer Token:** Utilizado por operadores e usuários da interface web (`Authorization: Bearer <jwt_token>`).
  2. **API Keys / Service Tokens (PAT):** Utilizados para integrações M2M (Machine-to-Machine), workers Python de IA e automações externas (`Authorization: Bearer vms_live_...` ou header `X-API-Key: vms_live_...`).
- **Bloqueio Automático (401 Unauthorized):** Qualquer requisição sem token válido ou com token expirado/revogado é rejeitada imediatamente no middleware HTTP antes de atingir qualquer controller ou camada de aplicação.

---

## 2. Isolamento Multi-Tenant Estrito (Data Leak Prevention)

- **Princípio de Tenant Scoping:** Toda e qualquer consulta (SELECT), mutação (INSERT, UPDATE, DELETE) e busca no banco de dados DEVE ser escopada pelo `tenant_id` injetado pelo middleware de autenticação a partir das claims seguras do Token.
- **Proibição de Queries sem Tenant:** Nenhuma rota acessível por usuário comum pode executar queries sem a cláusula `WHERE tenant_id = $1`.
- **Validação de Pertinência de Recursos:** Antes de alterar ou deletar câmeras, regras, zonas ou eventos, o backend deve verificar se o ID do recurso pertence estritamente ao `tenant_id` da sessão ativa.

---

## 3. Proteção de Mídia e Armazenamento (MinIO S3)

- **Bucket Não Público:** Os buckets do MinIO (`hydravms-events` e `hydravms-recordings`) NUNCA devem ter políticas de acesso público anônimo habilitadas.
- **Presigned URLs Exclusivas:** O acesso a vídeos contínuos (`.mp4`), clipes de eventos e snapshots (`.jpg`) é concedido unicamente via **Presigned URLs com TTL curto (15 a 30 minutos)** geradas sob demanda pelo backend após autenticação do usuário.
- **Hierarquia com Prefixo de Tenant:** Todo objeto gravado no MinIO deve obrigatoriamente possuir a chave prefixada com `{tenant_id}/` (ex: `tenant_123/2026/09/03/...`).

---

## 4. Autenticação, Tokens e RBAC

- **Algoritmo de Criptografia:** Senhas devem ser hasheadas com `bcrypt` (custo mínimo 12) ou `Argon2id`.
- **Tokens JWT:** Assinados com `HMAC-SHA256` ou `Ed25519` com expiração de curta duração (15 a 60 minutos) e suporte a rotação de refresh tokens.
- **Perfis de Acesso (RBAC):**
  - `admin`: Controle total do tenant (câmeras, usuários, regras de IA, retenção e discos).
  - `operator`: Monitoramento ao vivo, visualização de eventos, reconhecimento de alarmes e exportação de clipes.
  - `viewer`: Apenas visualização de live view e playback básico (sem permissão de exportar ou alterar regras).

---

## 5. Proteção de Rede & Prevenção de Ataques (SSRF / Injeção)

- **Sanitização de URLs RTSP (Anti-SSRF):** Antes de conectar a um stream RTSP fornecido pelo usuário, validar o formato da URI e bloquear endereços locais perigosos.
- **Prepared Statements Obrigatórios:** Todas as operações no PostgreSQL devem utilizar placeholders parametrizados (`$1, $2...` via `pgx`) para mitigar 100% de risco de SQL Injection.
- **Rate Limiting:** Rotas de login e conexões de WebSocket protegidas contra ataques de força bruta.

---

## 6. Auditoria de Segurança Imutável

- Toda ação administrativa e consultas críticas geram registros imediatos na tabela `audit_logs` com IP de origem, timestamp, `user_id` e payload.
