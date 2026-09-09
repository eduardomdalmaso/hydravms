-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v2
-- Unified Folders, Mosaic Layouts, Interactive Maps & Virtual Tours
-- PostgreSQL Dialect with Multi-Tenant Scoping, Indexes, and RLS Policies
-- ============================================================================

-- 1. Universal Folders Table (Multi-Module Hierarchy)
CREATE TABLE IF NOT EXISTS folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    module VARCHAR(32) NOT NULL, -- 'cameras', 'layouts', 'maps', 'tours', 'workflows', 'users'
    parent_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    color_hex VARCHAR(9) DEFAULT '#ff5e3a',
    icon VARCHAR(32) DEFAULT 'folder',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_tenant_module_parent_name UNIQUE (tenant_id, module, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid), name)
);

CREATE INDEX IF NOT EXISTS idx_folders_tenant_module ON folders (tenant_id, module, parent_id);
CREATE INDEX IF NOT EXISTS idx_folders_parent ON folders (parent_id);

-- 2. Link Cameras to Unified Folders
ALTER TABLE cameras ADD COLUMN IF NOT EXISTS folder_id UUID REFERENCES folders(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_cameras_folder ON cameras (tenant_id, folder_id);

-- 3. Mosaic Layouts (Saved Camera Grids & Layout Workspaces)
CREATE TABLE IF NOT EXISTS mosaic_layouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(128) NOT NULL,
    grid_type VARCHAR(32) NOT NULL DEFAULT '2x2', -- '1x1', '2x2', '3x3', '4x4', 'custom', 'hero_bottom'
    is_shared BOOLEAN NOT NULL DEFAULT TRUE,      -- Visible to all operators in tenant
    slots_config JSONB NOT NULL DEFAULT '[]'::jsonb, -- Array of { slot_index: 0, camera_id: 'cam_01', ptz_lock: false }
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mosaic_layouts_tenant ON mosaic_layouts (tenant_id, folder_id);
CREATE INDEX IF NOT EXISTS idx_mosaic_layouts_user ON mosaic_layouts (tenant_id, user_id);

-- 4. Interactive Maps & Floorplans (GIS & Image Schematics)
CREATE TABLE IF NOT EXISTS interactive_maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    name VARCHAR(128) NOT NULL,
    map_type VARCHAR(32) NOT NULL DEFAULT 'image', -- 'image' (Floorplan PNG/SVG) or 'gis' (OpenStreetMap / Leaflet tiles)
    image_s3_bucket VARCHAR(64) DEFAULT 'hydravms-maps',
    image_s3_key TEXT,
    center_lat DOUBLE PRECISION,
    center_lng DOUBLE PRECISION,
    zoom_level INT DEFAULT 16,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_interactive_maps_tenant ON interactive_maps (tenant_id, folder_id);

-- 5. Map Camera & Sensor Pins (Position, FOV Cone & Status)
CREATE TABLE IF NOT EXISTS map_pins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    map_id UUID NOT NULL REFERENCES interactive_maps(id) ON DELETE CASCADE,
    camera_id VARCHAR(64),
    pin_type VARCHAR(32) NOT NULL DEFAULT 'camera', -- 'camera', 'sensor', 'alarm_zone', 'access_gate'
    position_x_pct NUMERIC(5,2),                    -- Relative percentage (0.00% to 100.00%) for image maps
    position_y_pct NUMERIC(5,2),                    -- Relative percentage (0.00% to 100.00%) for image maps
    gps_lat DOUBLE PRECISION,                       -- For GIS maps
    gps_lng DOUBLE PRECISION,                       -- For GIS maps
    fov_angle_deg INT DEFAULT 90,                   -- Field of view opening angle (e.g. 90deg, 360deg)
    fov_direction_deg INT DEFAULT 0,                -- Direction angle (0-360deg)
    fov_radius_px INT DEFAULT 80,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_map_pins_map ON map_pins (tenant_id, map_id);
CREATE INDEX IF NOT EXISTS idx_map_pins_camera ON map_pins (tenant_id, camera_id);

-- 6. Virtual Tours (Automated Guard Patrol Sequences)
CREATE TABLE IF NOT EXISTS virtual_tours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    loop_mode BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_virtual_tours_tenant ON virtual_tours (tenant_id, folder_id);

-- 7. Tour Steps (Sequence of Camera Views & Dwell Times)
CREATE TABLE IF NOT EXISTS tour_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tour_id UUID NOT NULL REFERENCES virtual_tours(id) ON DELETE CASCADE,
    step_order INT NOT NULL DEFAULT 0,
    camera_id VARCHAR(64),
    layout_id UUID REFERENCES mosaic_layouts(id) ON DELETE SET NULL,
    dwell_time_seconds INT NOT NULL DEFAULT 10,
    ptz_preset_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tour_steps_tour ON tour_steps (tour_id, step_order);

-- 8. Row-Level Security (RLS) Policies for Zero-Leakage Multi-Tenancy
ALTER TABLE folders ENABLE ROW LEVEL SECURITY;
ALTER TABLE folders FORCE ROW LEVEL SECURITY;
CREATE POLICY folders_tenant_isolation ON folders
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE mosaic_layouts ENABLE ROW LEVEL SECURITY;
ALTER TABLE mosaic_layouts FORCE ROW LEVEL SECURITY;
CREATE POLICY mosaic_layouts_tenant_isolation ON mosaic_layouts
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE interactive_maps ENABLE ROW LEVEL SECURITY;
ALTER TABLE interactive_maps FORCE ROW LEVEL SECURITY;
CREATE POLICY interactive_maps_tenant_isolation ON interactive_maps
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE map_pins ENABLE ROW LEVEL SECURITY;
ALTER TABLE map_pins FORCE ROW LEVEL SECURITY;
CREATE POLICY map_pins_tenant_isolation ON map_pins
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);

ALTER TABLE virtual_tours ENABLE ROW LEVEL SECURITY;
ALTER TABLE virtual_tours FORCE ROW LEVEL SECURITY;
CREATE POLICY virtual_tours_tenant_isolation ON virtual_tours
    FOR ALL
    USING (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid);
