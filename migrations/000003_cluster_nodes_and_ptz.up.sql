-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v3
-- Multi-Node Clustering (Edge Ingest, GPU Workers, Control Plane) & ONVIF PTZ
-- PostgreSQL Dialect with Multi-Tenant Scoping, Indexes, and RLS Policies
-- ============================================================================

-- 1. Cluster Nodes (Distributed Multi-Server Topology)
CREATE TABLE IF NOT EXISTS cluster_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE, -- NULL for shared cluster infra, or scoped to dedicated tenant
    node_name VARCHAR(128) NOT NULL,
    node_role VARCHAR(32) NOT NULL DEFAULT 'edge_ingest', -- 'control_plane', 'edge_ingest', 'gpu_worker'
    ip_address VARCHAR(45) NOT NULL,
    grpc_port INT NOT NULL DEFAULT 50051,
    webrtc_port INT NOT NULL DEFAULT 8889,
    http_port INT NOT NULL DEFAULT 8080,
    gpu_device_info VARCHAR(128),                         -- e.g. "NVIDIA GeForce RTX 5090 (32GB VRAM)"
    cuda_compute_capability VARCHAR(16),                  -- e.g. "sm_120"
    cpu_usage_pct NUMERIC(5,2) DEFAULT 0.0,
    ram_usage_pct NUMERIC(5,2) DEFAULT 0.0,
    gpu_usage_pct NUMERIC(5,2) DEFAULT 0.0,
    vram_used_mb BIGINT DEFAULT 0,
    active_streams_count INT NOT NULL DEFAULT 0,
    max_streams_capacity INT NOT NULL DEFAULT 64,
    status VARCHAR(32) NOT NULL DEFAULT 'online',         -- 'online', 'offline', 'draining', 'error'
    last_heartbeat_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cluster_nodes_role ON cluster_nodes (node_role, status);
CREATE INDEX IF NOT EXISTS idx_cluster_nodes_heartbeat ON cluster_nodes (last_heartbeat_at DESC);
CREATE INDEX IF NOT EXISTS idx_cluster_nodes_tenant ON cluster_nodes (tenant_id) WHERE tenant_id IS NOT NULL;

-- 2. Link Cameras to Assigned Ingest & Inference Cluster Node
ALTER TABLE cameras ADD COLUMN IF NOT EXISTS assigned_node_id UUID REFERENCES cluster_nodes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_cameras_assigned_node ON cameras (assigned_node_id);

-- 3. ONVIF PTZ Presets (Saved Positions with Visual Thumbnails)
CREATE TABLE IF NOT EXISTS ptz_presets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    preset_token VARCHAR(64) NOT NULL,                    -- ONVIF PresetToken identifier on physical camera
    name VARCHAR(128) NOT NULL,
    pan_coord DOUBLE PRECISION,                           -- Absolute Normalized Pan (-1.0 to 1.0)
    tilt_coord DOUBLE PRECISION,                          -- Absolute Normalized Tilt (-1.0 to 1.0)
    zoom_coord DOUBLE PRECISION,                          -- Absolute Normalized Zoom (0.0 to 1.0)
    thumbnail_s3_key TEXT,                                -- Snapshot visual preview of preset view
    is_home_preset BOOLEAN NOT NULL DEFAULT FALSE,        -- Return to this position on idle
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ptz_preset_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE,
    CONSTRAINT uq_camera_preset_token UNIQUE (tenant_id, camera_id, preset_token)
);

CREATE INDEX IF NOT EXISTS idx_ptz_presets_camera ON ptz_presets (tenant_id, camera_id);

-- 4. ONVIF PTZ Patrols (Hardware Guard Tours)
CREATE TABLE IF NOT EXISTS ptz_patrols (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    loop_mode BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ptz_patrol_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ptz_patrols_camera ON ptz_patrols (tenant_id, camera_id);

-- 5. PTZ Patrol Points (Sequence of Presets, Dwell Times and Speeds)
CREATE TABLE IF NOT EXISTS ptz_patrol_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patrol_id UUID NOT NULL REFERENCES ptz_patrols(id) ON DELETE CASCADE,
    preset_id UUID NOT NULL REFERENCES ptz_presets(id) ON DELETE CASCADE,
    step_order INT NOT NULL DEFAULT 0,
    dwell_time_seconds INT NOT NULL DEFAULT 15,
    transition_speed_pct INT NOT NULL DEFAULT 100,        -- 1% to 100% Pan/Tilt movement speed
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ptz_patrol_points_patrol ON ptz_patrol_points (patrol_id, step_order);

-- 6. Row-Level Security (RLS) Policies
ALTER TABLE cluster_nodes ENABLE ROW LEVEL SECURITY;
ALTER TABLE cluster_nodes FORCE ROW LEVEL SECURITY;
CREATE POLICY cluster_nodes_tenant_isolation ON cluster_nodes
    FOR ALL
    USING (tenant_id IS NULL OR tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id IS NULL OR tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE ptz_presets ENABLE ROW LEVEL SECURITY;
ALTER TABLE ptz_presets FORCE ROW LEVEL SECURITY;
CREATE POLICY ptz_presets_tenant_isolation ON ptz_presets
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE ptz_patrols ENABLE ROW LEVEL SECURITY;
ALTER TABLE ptz_patrols FORCE ROW LEVEL SECURITY;
CREATE POLICY ptz_patrols_tenant_isolation ON ptz_patrols
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);
