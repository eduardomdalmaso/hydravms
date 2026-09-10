ALTER TABLE camera_recording_profiles ADD COLUMN IF NOT EXISTS name VARCHAR(128) NOT NULL DEFAULT 'Perfil Principal';
ALTER TABLE camera_recording_profiles DROP CONSTRAINT IF EXISTS uq_camera_recording_profile;
ALTER TABLE camera_recording_profiles ADD CONSTRAINT uq_camera_recording_profile UNIQUE (tenant_id, camera_id, id);
