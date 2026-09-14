# 🔒 Diretrizes Invioláveis de Segurança & Hardening (Hydra Ecosystem)

Este documento estabelece as diretrizes invioláveis de segurança, autenticação obrigatória, isolamento multi-tenant, autorização no backend (RBAC), prevenção a IDOR, saneamento de segredos e imunidade a injeções para o ecossistema Hydra.

---

## 1. Banco com Tranca (Isolamento Multi-Tenant Estrito)
- **Tenant Scoping Obrigatório:** Toda e qualquer consulta (`SELECT`), mutação (`INSERT`, `UPDATE`, `DELETE`) e agregação no banco de dados DEVE ser filtrada pelo `tenant_id` autenticado via `middleware.GetTenantID(r.Context())`.
- **Proibição de Bypass de Tenant:** Nunca utilize cláusulas permissivas como `OR tenant_id = '00000000-0000-0000-0000-000000000001'` em consultas de operadores ou inquilinos regulares.
- **Consultas Globais Protegidas:** Rotas que listam dados de múltiplos inquilinos (ex: superadmin) devem validar explicitamente se o usuário possui `role == 'superadmin'` e pertence ao tenant raiz.

---

## 2. Permissão no Servidor (Backend RBAC Enforcement)
- **Nunca Confiar Apenas na UI:** Ocultar elementos no frontend (`v-if="isAdmin"`) é apenas conveniência visual. Toda operação administrativa de escrita ou consulta sensível DEVE ser validada no backend com `middleware.RequireRole("admin", "superadmin")`.
- **Rejeição Automática (403 Forbidden):** Requisições com tokens de `operator` ou `viewer` que tentem acessar endpoints de configuração, gestão de usuários, storage, nós ou auditoria devem ser bloqueadas no gateway HTTP.

---

## 3. Prevenção a IDOR (Insecure Direct Object Reference)
- **Validação de Posse em Recursos por ID:** Ao buscar, alterar ou deletar qualquer recurso por ID na URL (`/api/v1/{resource}/{id}`), a consulta SQL DEVE incluir `WHERE id = $1 AND tenant_id = $2`.
- **Presigned URLs Seguras:** Presigned URLs para S3/MinIO devem validar permissão e escopo de tenant antes de assinar links de download.

---

## 4. Chaves & Segredos (Zero Hardcoded Secrets)
- **Proibição de Fallbacks Estáticos:** Nunca utilize strings estáticas em código para `JWT_SECRET`, senhas de banco ou chaves de API.
- **Validação de Startup:** Se `JWT_SECRET` não for configurado no ambiente, o sistema deve gerar dinamicamente uma chave criptográfica forte aleatória via `crypto/rand` com alerta explícito de segurança.
- **Preservação de Senhas:** Rotinas de inicialização (`EnsureAdminUser`) NUNCA devem sobrescrever a senha do administrador se a conta já existir no banco (`ON CONFLICT (id) DO NOTHING`).

---

## 5. Sanitização de Inputs, LFI & Mitigação de SSRF
- **Prevenção a Arbitrary File Read (Path Traversal):** Antes de usar `http.ServeFile` ou abrir caminhos de disco fornecidos por usuários, sanitize o caminho com `filepath.Clean(path)` e valide que o caminho canônico esteja estritamente restrito ao diretório permitido de gravações (`storage/`). Bloqueie qualquer acesso a `/etc`, `/proc`, `/sys`, `/root`, etc., com `403 Forbidden`.
- **Prevenção a SSRF:** Endpoints de probing ou conexão de rede devem rejeitar endereços de metadados em nuvem (`169.254.169.254`, `0.0.0.0`, `255.255.255.255`) e validar faixas de portas permitidas.
- **Prepared Statements:** 100% das operações SQL devem utilizar placeholders parametrizados (`$1, $2...`) mitigando SQL Injection.
- **XSS Protegido:** Nunca utilize `v-html`, `innerHTML` ou concatenação de HTML com entradas de usuários.
