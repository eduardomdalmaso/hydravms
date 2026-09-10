-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v4
-- Performance Indexes, JSONB GIN Path Ops & Security Hardening
-- PostgreSQL 16 Dialect
-- ============================================================================

-- 1. Composite Tree Traversal Index for Hierarchical Folders CTE
CREATE INDEX IF NOT EXISTS idx_folders_tree_traversal 
ON folders (tenant_id, module, parent_id, sort_order);

-- 2. Add dynamic JSONB metadata column for AI detections if not present
ALTER TABLE events ADD COLUMN IF NOT EXISTS detection_metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

-- 3. High-Speed JSONB Indexes using jsonb_path_ops for YOLO BBoxes & Metadata
CREATE INDEX IF NOT EXISTS idx_events_bbox_fast 
ON events USING GIN (bbox_normalized jsonb_path_ops);

CREATE INDEX IF NOT EXISTS idx_events_detection_metadata_fast 
ON events USING GIN (detection_metadata jsonb_path_ops);

-- 4. Partial Index for Critical Active Alarms (<2ms Query Response)
CREATE INDEX IF NOT EXISTS idx_events_active_critical 
ON events (tenant_id, created_at DESC)
WHERE status = 'new' AND severity IN ('high', 'critical');

-- 5. Dedicated Application Role (Least Privilege Security Enforcement)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'hydra_app') THEN
        CREATE ROLE hydra_app WITH LOGIN PASSWORD 'hydra_secure_app_pass_2026' NOSUPERUSER NOBYPASSRLS NOCREATEROLE NOCREATEDB;
    END IF;
END
$$;

GRANT CONNECT ON DATABASE hydravms TO hydra_app;
GRANT USAGE ON SCHEMA public TO hydra_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO hydra_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO hydra_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO hydra_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO hydra_app;
