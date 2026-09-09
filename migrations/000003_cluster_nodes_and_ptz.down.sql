-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v3 DOWN
-- Revert Cluster Nodes and PTZ
-- ============================================================================

DROP POLICY IF EXISTS ptz_patrols_tenant_isolation ON ptz_patrols;
DROP POLICY IF EXISTS ptz_presets_tenant_isolation ON ptz_presets;
DROP POLICY IF EXISTS cluster_nodes_tenant_isolation ON cluster_nodes;

DROP TABLE IF EXISTS ptz_patrol_points CASCADE;
DROP TABLE IF EXISTS ptz_patrols CASCADE;
DROP TABLE IF EXISTS ptz_presets CASCADE;

ALTER TABLE cameras DROP COLUMN IF EXISTS assigned_node_id;

DROP TABLE IF EXISTS cluster_nodes CASCADE;
