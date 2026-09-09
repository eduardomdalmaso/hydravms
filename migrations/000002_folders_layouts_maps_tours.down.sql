-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v2 DOWN
-- Revert Folders, Layouts, Maps and Tours
-- ============================================================================

DROP POLICY IF EXISTS virtual_tours_tenant_isolation ON virtual_tours;
DROP POLICY IF EXISTS map_pins_tenant_isolation ON map_pins;
DROP POLICY IF EXISTS interactive_maps_tenant_isolation ON interactive_maps;
DROP POLICY IF EXISTS mosaic_layouts_tenant_isolation ON mosaic_layouts;
DROP POLICY IF EXISTS folders_tenant_isolation ON folders;

DROP TABLE IF EXISTS tour_steps CASCADE;
DROP TABLE IF EXISTS virtual_tours CASCADE;
DROP TABLE IF EXISTS map_pins CASCADE;
DROP TABLE IF EXISTS interactive_maps CASCADE;
DROP TABLE IF EXISTS mosaic_layouts CASCADE;

ALTER TABLE cameras DROP COLUMN IF EXISTS folder_id;

DROP TABLE IF EXISTS folders CASCADE;
