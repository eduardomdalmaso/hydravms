-- ============================================================================
-- HYDRAVMS DATABASE SEED DATA
-- Default Tenant, Admin User, Folders, Nodes, Cameras, Layouts
-- ============================================================================

-- 1. Default Tenant
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

-- 3. Default Root & Child Folders for Cameras
INSERT INTO folders (id, tenant_id, module, parent_id, name, color_hex, icon, sort_order)
VALUES 
    ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'cameras', NULL, 'Perímetro Externo', '#ff5e3a', 'folder', 1),
    ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'cameras', NULL, 'Área Interna', '#00f0ff', 'folder', 2),
    ('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'cameras', '10000000-0000-0000-0000-000000000001', 'Portões & Acesso', '#ff5e3a', 'folder', 1),
    ('10000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', 'cameras', '10000000-0000-0000-0000-000000000002', 'Datacenter / Servidores', '#00ff9d', 'folder', 1)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 4. Default Folders for Layouts, Maps & Tours
INSERT INTO folders (id, tenant_id, module, parent_id, name, color_hex, icon, sort_order)
VALUES 
    ('20000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'layouts', NULL, 'Mosaicos Operacionais', '#ff5e3a', 'folder', 1),
    ('30000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'maps', NULL, 'Plantas Baixas', '#00f0ff', 'folder', 1),
    ('40000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'tours', NULL, 'Rondas Automáticas', '#fcee0a', 'folder', 1)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 5. Default Cluster Node
INSERT INTO cluster_nodes (id, tenant_id, node_name, node_role, ip_address, grpc_port, webrtc_port, http_port, status)
VALUES (
    '60000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '[NODE] HYDRA-CORE-MASTER',
    'control_plane',
    '127.0.0.1',
    50051,
    8889,
    8083,
    'online'
) ON CONFLICT (id) DO UPDATE SET node_name = EXCLUDED.node_name;

-- 6. Default Cameras
INSERT INTO cameras (id, tenant_id, name, protocol, rtsp_url, sub_stream_url, location, status, resolution, fps, bitrate_kbps, is_active, folder_id, assigned_node_id)
VALUES 
    ('cam_entrance_01', '00000000-0000-0000-0000-000000000001', '[STREAM] CAM_01 // ENTRADA PRINCIPAL', 'rtsp', 'rtsp://localhost:8554/cam_entrance_01', 'rtsp://localhost:8554/cam_entrance_01_sub', 'Guarita Principal', 'online', '1920x1080', 30.0, 4096, TRUE, '10000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000001'),
    ('cam_perimeter_02', '00000000-0000-0000-0000-000000000001', '[STREAM] CAM_02 // PERIMETRO NORTE', 'rtsp', 'rtsp://localhost:8554/cam_perimeter_02', 'rtsp://localhost:8554/cam_perimeter_02_sub', 'Muro Norte', 'online', '1920x1080', 30.0, 3072, TRUE, '10000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001'),
    ('cam_parking_03', '00000000-0000-0000-0000-000000000001', '[STREAM] CAM_03 // ESTACIONAMENTO', 'rtsp', 'rtsp://localhost:8554/cam_parking_03', 'rtsp://localhost:8554/cam_parking_03_sub', 'Pátio de Veículos', 'online', '1920x1080', 30.0, 3072, TRUE, '10000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001'),
    ('cam_server_04', '00000000-0000-0000-0000-000000000001', '[STREAM] CAM_04 // SALA DE SERVIDORES', 'rtsp', 'rtsp://localhost:8554/cam_server_04', 'rtsp://localhost:8554/cam_server_04_sub', 'Rack 04 - Core NOC', 'online', '1920x1080', 30.0, 2048, TRUE, '10000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000001')
ON CONFLICT (tenant_id, id) DO UPDATE SET name = EXCLUDED.name;

-- 7. Default Mosaic Layout
INSERT INTO mosaic_layouts (id, tenant_id, folder_id, user_id, name, grid_type, is_shared, slots_config)
VALUES (
    '50000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000002',
    'Mosaico 2x2 Geral',
    '2x2',
    TRUE,
    '[{"slot_index": 0, "camera_id": "cam_entrance_01"}, {"slot_index": 1, "camera_id": "cam_perimeter_02"}, {"slot_index": 2, "camera_id": "cam_parking_03"}, {"slot_index": 3, "camera_id": "cam_server_04"}]'::jsonb
) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;
