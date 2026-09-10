-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v5
-- Storage Pools Enhancement for MinIO S3, NVMe Buffers & NAS
-- PostgreSQL 16 Dialect
-- ============================================================================

ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS source_type VARCHAR(32) NOT NULL DEFAULT 'LOCAL_DISK';
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS role VARCHAR(32) NOT NULL DEFAULT 'WARM_ARCHIVE';
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS node_or_server VARCHAR(128) NOT NULL DEFAULT 'MAQUINA LOCAL // NO 01';
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS path_or_endpoint TEXT;
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS filesystem VARCHAR(32) NOT NULL DEFAULT 'XFS';
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS retention_days INT NOT NULL DEFAULT 30;
ALTER TABLE storage_pools ADD COLUMN IF NOT EXISTS is_spillover_active BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE storage_pools ALTER COLUMN mount_path DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_storage_pools_role ON storage_pools (role, is_active);
CREATE INDEX IF NOT EXISTS idx_storage_pools_source_type ON storage_pools (source_type);
