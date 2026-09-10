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
    'HYDRA-VMS-NODE-0',
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

-- 8. Additional Users
INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active)
VALUES 
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Operador NOC', 'operador@hydravms.io', '$2a$10$7EqJtq98hPqEX7fNZaFWoOZh9g.iI0g.pX6QyHkLq0/3rYxZ/hD8W', 'operator', TRUE),
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', 'Auditor Forense', 'auditor@hydravms.io', '$2a$10$7EqJtq98hPqEX7fNZaFWoOZh9g.iI0g.pX6QyHkLq0/3rYxZ/hD8W', 'auditor', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 9. Additional Cluster Nodes
INSERT INTO cluster_nodes (id, tenant_id, node_name, node_role, ip_address, grpc_port, webrtc_port, http_port, status)
VALUES 
    ('60000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '[NODE] HYDRA-EDGE-INGEST-01', 'edge_ingest', '192.168.1.10', 50051, 8889, 8080, 'online'),
    ('60000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '[NODE] HYDRA-GPU-RTX5090', 'gpu_worker', '192.168.1.20', 50051, 8889, 8081, 'online')
ON CONFLICT (id) DO UPDATE SET node_name = EXCLUDED.node_name;

-- 10. Interactive Maps
INSERT INTO interactive_maps (id, tenant_id, folder_id, name, map_type, image_s3_bucket, image_s3_key, is_active)
VALUES (
    '70000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '30000000-0000-0000-0000-000000000001',
    'Planta Geral - Galpão Principal',
    'image',
    'hydravms-maps',
    'floorplans/planta_galpao_01.png',
    TRUE
) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO map_pins (id, tenant_id, map_id, camera_id, pin_type, position_x_pct, position_y_pct, fov_angle_deg, fov_direction_deg)
VALUES 
    ('71000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', 'cam_entrance_01', 'camera', 15.50, 42.00, 90, 45),
    ('71000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', 'cam_perimeter_02', 'camera', 82.00, 20.00, 110, 180),
    ('71000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', 'cam_server_04', 'camera', 48.00, 75.00, 75, 270)
ON CONFLICT (id) DO NOTHING;

-- 11. Virtual Tours
INSERT INTO virtual_tours (id, tenant_id, folder_id, name, description, is_active, loop_mode)
VALUES (
    '80000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '40000000-0000-0000-0000-000000000001',
    'Ronda Perimetral Noturna',
    'Varredura perimétrica dos portões de acesso e sala técnica',
    TRUE,
    TRUE
) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO tour_steps (id, tour_id, step_order, camera_id, dwell_time_seconds)
VALUES 
    ('81000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', 1, 'cam_entrance_01', 10),
    ('81000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', 2, 'cam_perimeter_02', 12),
    ('81000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000001', 3, 'cam_parking_03', 10),
    ('81000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000001', 4, 'cam_server_04', 15)
ON CONFLICT (id) DO NOTHING;

-- 12. Plugins Marketplace
INSERT INTO plugins (id, name, version, author, category, runtime, entrypoint, is_official)
VALUES 
    ('plugin-yolo-lpr', 'LPR / Reconhecimento de Placas', '2.4.0', 'Hydra Intelligence', 'analytics', 'python3', 'lpr.py', TRUE),
    ('plugin-facial-rec', 'Reconhecimento Facial Forense', '1.9.5', 'Hydra Intelligence', 'analytics', 'python3', 'face.py', TRUE),
    ('plugin-ppe-detection', 'Detecção de EPI (Segurança do Trabalho)', '3.1.0', 'Hydra Intelligence', 'safety', 'python3', 'ppe.py', TRUE),
    ('plugin-fire-smoke', 'Detecção Precoce de Fogo e Fumaça', '2.0.1', 'Hydra Intelligence', 'safety', 'python3', 'fire.py', TRUE)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 13. System Audit Logs
INSERT INTO audit_logs (id, tenant_id, user_id, ip_address, action, entity_type, entity_id, payload_json)
VALUES 
    ('90000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '127.0.0.1', 'SYSTEM_BOOT', 'system', 'core', '{"message": "HydraVMS Control Plane inicializado com sucesso", "version": "1.0.0"}'),
    ('90000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '127.0.0.1', 'CAMERA_SYNC', 'camera', 'cam_entrance_01', '{"message": "Fluxo RTSP TCP Interleaved estabelecido com 30 FPS", "codec": "H.265"}'),
    ('90000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '127.0.0.1', 'STORAGE_CHECK', 'storage_pool', 'pool_nvme_01', '{"message": "Buffer NVMe XFS montado com 953 GB livres", "watermark": "healthy"}')
ON CONFLICT (id) DO NOTHING;


