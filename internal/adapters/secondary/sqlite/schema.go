package sqlite

import "database/sql"

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tenants (
		id TEXT PRIMARY KEY,
		slug TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		plan TEXT NOT NULL DEFAULT 'enterprise',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'admin',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		last_login_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS folders (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		module TEXT NOT NULL DEFAULT 'cameras',
		parent_id TEXT,
		name TEXT NOT NULL,
		color_hex TEXT NOT NULL DEFAULT '#3b82f6',
		icon TEXT NOT NULL DEFAULT 'folder',
		sort_order INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS cameras (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		name TEXT NOT NULL,
		protocol TEXT NOT NULL DEFAULT 'RTSP',
		rtsp_url TEXT,
		sub_stream_url TEXT,
		onvif_ip TEXT,
		onvif_port INTEGER DEFAULT 80,
		onvif_user TEXT,
		onvif_pass TEXT,
		rtmp_stream_key TEXT,
		location TEXT,
		status TEXT NOT NULL DEFAULT 'ONLINE',
		resolution TEXT NOT NULL DEFAULT '1920x1080',
		fps REAL NOT NULL DEFAULT 30.0,
		bitrate_kbps REAL NOT NULL DEFAULT 4000.0,
		codec TEXT NOT NULL DEFAULT 'H.264',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		folder_id TEXT,
		assigned_node_id TEXT,
		is_ptz BOOLEAN NOT NULL DEFAULT 0,
		is_audio_enabled BOOLEAN NOT NULL DEFAULT 0,
		is_recording BOOLEAN NOT NULL DEFAULT 1,
		record_mode TEXT NOT NULL DEFAULT 'continuous',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS storage_pools (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		name TEXT NOT NULL,
		mount_path TEXT NOT NULL UNIQUE,
		total_bytes INTEGER NOT NULL,
		used_bytes INTEGER NOT NULL,
		watermark_warning_pct REAL NOT NULL DEFAULT 85.0,
		watermark_critical_pct REAL NOT NULL DEFAULT 95.0,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS recordings (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		storage_pool_id TEXT,
		file_path TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		duration_seconds INTEGER NOT NULL,
		file_size_bytes INTEGER NOT NULL,
		record_mode TEXT NOT NULL DEFAULT 'continuous',
		is_pinned BOOLEAN NOT NULL DEFAULT 0,
		checksum_sha256 TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (camera_id) REFERENCES cameras(id) ON DELETE CASCADE,
		FOREIGN KEY (storage_pool_id) REFERENCES storage_pools(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS camera_recording_profiles (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT 'Perfil Principal',
		mode TEXT NOT NULL DEFAULT 'continuous',
		stream_type TEXT NOT NULL DEFAULT 'main',
		segment_duration_s INTEGER NOT NULL DEFAULT 60,
		pre_buffer_s INTEGER NOT NULL DEFAULT 5,
		post_buffer_s INTEGER NOT NULL DEFAULT 15,
		retention_days INTEGER NOT NULL DEFAULT 30,
		schedule_json TEXT NOT NULL DEFAULT '[]',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (camera_id) REFERENCES cameras(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		rule_id TEXT,
		zone_id TEXT,
		event_type TEXT NOT NULL,
		severity TEXT NOT NULL DEFAULT 'medium',
		status TEXT NOT NULL DEFAULT 'new',
		triggered_at DATETIME NOT NULL,
		resolved_at DATETIME,
		resolved_by_user_id TEXT,
		object_class TEXT,
		confidence REAL NOT NULL DEFAULT 0.0,
		bbox_normalized TEXT NOT NULL DEFAULT '{}',
		tracking_id INTEGER,
		snapshot_s3_key TEXT,
		crop_s3_key TEXT,
		clip_s3_key TEXT,
		is_pinned BOOLEAN NOT NULL DEFAULT 0,
		notes TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (camera_id) REFERENCES cameras(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id TEXT NOT NULL,
		user_id TEXT,
		ip_address TEXT,
		action TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id TEXT,
		payload_json TEXT NOT NULL DEFAULT '{}',
		created_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS cluster_nodes (
		id TEXT PRIMARY KEY,
		node_name TEXT NOT NULL,
		node_role TEXT NOT NULL DEFAULT 'EDGE_INGEST',
		ip_address TEXT NOT NULL,
		grpc_port INTEGER NOT NULL DEFAULT 50051,
		webrtc_port INTEGER NOT NULL DEFAULT 8889,
		http_port INTEGER NOT NULL DEFAULT 8083,
		gpu_device_info TEXT,
		cpu_usage_pct REAL NOT NULL DEFAULT 0.0,
		ram_usage_pct REAL NOT NULL DEFAULT 0.0,
		gpu_usage_pct REAL NOT NULL DEFAULT 0.0,
		vram_used_mb INTEGER NOT NULL DEFAULT 0,
		active_streams_count INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'ONLINE',
		last_heartbeat_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS plugins (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		version TEXT NOT NULL,
		author TEXT,
		category TEXT NOT NULL DEFAULT 'analytics',
		runtime TEXT NOT NULL DEFAULT 'python3',
		entrypoint TEXT NOT NULL,
		min_vms_version TEXT NOT NULL DEFAULT '1.0.0',
		permissions TEXT NOT NULL DEFAULT '[]',
		config_schema TEXT NOT NULL DEFAULT '{}',
		ui_schema TEXT NOT NULL DEFAULT '{}',
		is_official BOOLEAN NOT NULL DEFAULT 1,
		is_deprecated BOOLEAN NOT NULL DEFAULT 0,
		package_url TEXT,
		package_checksum TEXT,
		hardware_req TEXT,
		description TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tenant_plugins (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		plugin_id TEXT NOT NULL,
		installed_version TEXT NOT NULL,
		previous_version TEXT,
		last_stable_version TEXT,
		is_enabled BOOLEAN NOT NULL DEFAULT 1,
		status TEXT NOT NULL DEFAULT 'running',
		pid INTEGER,
		assigned_gpu_device TEXT DEFAULT '0',
		config_values TEXT NOT NULL DEFAULT '{}',
		last_health_check DATETIME,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (plugin_id) REFERENCES plugins(id) ON DELETE CASCADE,
		UNIQUE (tenant_id, plugin_id)
	);

	CREATE TABLE IF NOT EXISTS layouts (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		folder_id TEXT,
		user_id TEXT,
		name TEXT NOT NULL,
		grid_type TEXT NOT NULL DEFAULT '2x2',
		is_locked BOOLEAN NOT NULL DEFAULT 0,
		is_shared BOOLEAN NOT NULL DEFAULT 1,
		target_monitor INTEGER NOT NULL DEFAULT 0,
		allowed_user_ids TEXT NOT NULL DEFAULT '[]',
		slots_config TEXT NOT NULL DEFAULT '[]',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
		FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_cameras_tenant ON cameras(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_folders_tenant ON folders(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_layouts_tenant ON layouts(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_events_camera_time ON events(tenant_id, camera_id, triggered_at);
	CREATE INDEX IF NOT EXISTS idx_recordings_range ON recordings(tenant_id, camera_id, start_time, end_time);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON audit_logs(tenant_id, created_at);
	CREATE INDEX IF NOT EXISTS idx_tenant_plugins_tenant ON tenant_plugins(tenant_id, is_enabled);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	_, _ = db.Exec("ALTER TABLE camera_recording_profiles ADD COLUMN stream_type TEXT NOT NULL DEFAULT 'main'")
	return nil
}
