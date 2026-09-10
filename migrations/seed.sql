-- ============================================================================
-- HYDRAVMS DATABASE SEED DATA
-- Default Tenant, Super Admin User, Control Plane Node & Official Plugins
-- (ZERO MOCK DATA)
-- ============================================================================

-- 1. Default Master Tenant
INSERT INTO tenants (id, slug, name, plan, max_cameras, max_retention_days, max_storage_bytes, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'master',
    'Hydra Master Operations',
    'enterprise',
    64,
    90,
    10995116277760,
    TRUE
) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 2. Default Super Admin User (password: admin123 -> bcrypt hash)
INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'Cyber Admin',
    'admin@hydravms.io',
    '$2a$10$7EqJtq98hPqEX7fNZaFWoOZh9g.iI0g.pX6QyHkLq0/3rYxZ/hD8W',
    'super_admin',
    TRUE
) ON CONFLICT (id) DO NOTHING;

-- 3. Default Control Plane Cluster Node
INSERT INTO cluster_nodes (id, tenant_id, node_name, node_role, ip_address, grpc_port, webrtc_port, http_port, status)
VALUES (
    '60000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'HYDRA-VMS-NODE-0',
    'control_plane',
    '127.0.0.1',
    50051,
    8889,
    8083,
    'online'
) ON CONFLICT (id) DO UPDATE SET node_name = EXCLUDED.node_name;

-- 4. Official Plugins Catalogue (Marketplace Definitions)
INSERT INTO plugins (id, name, version, author, category, runtime, entrypoint, is_official)
VALUES 
    ('plugin-yolo-lpr', 'LPR / Reconhecimento de Placas', '2.4.0', 'Hydra Intelligence', 'analytics', 'python3', 'lpr.py', TRUE),
    ('plugin-facial-rec', 'Reconhecimento Facial Forense', '1.9.5', 'Hydra Intelligence', 'analytics', 'python3', 'face.py', TRUE),
    ('plugin-ppe-detection', 'Detecção de EPI (Segurança do Trabalho)', '3.1.0', 'Hydra Intelligence', 'safety', 'python3', 'ppe.py', TRUE),
    ('plugin-fire-smoke', 'Detecção Precoce de Fogo e Fumaça', '2.0.1', 'Hydra Intelligence', 'safety', 'python3', 'fire.py', TRUE)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 5. System Boot Audit Log
INSERT INTO audit_logs (id, tenant_id, user_id, ip_address, action, entity_type, entity_id, payload_json)
VALUES 
    ('90000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '127.0.0.1', 'SYSTEM_BOOT', 'system', 'core', '{"message": "HydraVMS Control Plane inicializado com sucesso", "version": "1.0.0"}')
ON CONFLICT (id) DO NOTHING;
