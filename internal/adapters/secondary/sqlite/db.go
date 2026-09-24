package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// OpenDB opens a SQLite database in WAL mode and initializes tables and seed data.
func OpenDB(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create db dir: %w", err)
		}
	}

	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// SQLite WAL supports concurrent readers with a single writer
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	if err := seedDefaultData(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to seed default data: %w", err)
	}

	return db, nil
}

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
		bitrate_kbps INTEGER NOT NULL DEFAULT 2048,
		codec TEXT NOT NULL DEFAULT 'H.264',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		folder_id TEXT,
		assigned_node_id TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS storage_pools (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		source_type TEXT NOT NULL DEFAULT 'local_drive',
		role TEXT NOT NULL DEFAULT 'HOT_VIDEO_BUFFER',
		node_or_server TEXT NOT NULL DEFAULT 'edge-node-01',
		path_or_endpoint TEXT,
		filesystem TEXT NOT NULL DEFAULT 'XFS',
		total_bytes INTEGER NOT NULL DEFAULT 0,
		used_bytes INTEGER NOT NULL DEFAULT 0,
		available_bytes INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'HEALTHY',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		is_spillover_active BOOLEAN NOT NULL DEFAULT 0,
		retention_days INTEGER NOT NULL DEFAULT 30,
		last_checked_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id TEXT,
		user_id TEXT,
		ip_address TEXT NOT NULL DEFAULT '127.0.0.1',
		action TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id TEXT,
		payload_json TEXT NOT NULL DEFAULT '{}',
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		rule_id TEXT,
		zone_id TEXT,
		event_type TEXT NOT NULL,
		severity TEXT NOT NULL DEFAULT 'warning',
		status TEXT NOT NULL DEFAULT 'new',
		triggered_at DATETIME NOT NULL,
		resolved_at DATETIME,
		resolved_by_user_id TEXT,
		object_class TEXT NOT NULL DEFAULT '',
		confidence REAL NOT NULL DEFAULT 0.0,
		bbox_normalized TEXT NOT NULL DEFAULT '{}',
		tracking_id INTEGER DEFAULT 0,
		snapshot_s3_key TEXT,
		crop_s3_key TEXT,
		clip_s3_key TEXT,
		is_pinned BOOLEAN NOT NULL DEFAULT 0,
		notes TEXT,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS camera_recording_profiles (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT 'Perfil de Gravacao',
		mode TEXT NOT NULL DEFAULT 'continuous',
		segment_duration_s INTEGER NOT NULL DEFAULT 60,
		pre_buffer_s INTEGER NOT NULL DEFAULT 5,
		post_buffer_s INTEGER NOT NULL DEFAULT 10,
		retention_days INTEGER NOT NULL DEFAULT 30,
		schedule_json TEXT NOT NULL DEFAULT '[]',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS recordings (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		recording_mode TEXT NOT NULL DEFAULT 'continuous',
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		duration_seconds INTEGER NOT NULL DEFAULT 0,
		file_size_bytes INTEGER NOT NULL DEFAULT 0,
		s3_bucket TEXT NOT NULL DEFAULT 'hydravms-recordings',
		s3_key TEXT NOT NULL,
		is_pinned BOOLEAN NOT NULL DEFAULT 0,
		is_purged BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS cluster_nodes (
		id TEXT PRIMARY KEY,
		tenant_id TEXT,
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

	CREATE INDEX IF NOT EXISTS idx_cameras_tenant ON cameras(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_folders_tenant ON folders(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_events_camera_time ON events(tenant_id, camera_id, triggered_at);
	CREATE INDEX IF NOT EXISTS idx_recordings_range ON recordings(tenant_id, camera_id, start_time, end_time);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON audit_logs(tenant_id, created_at);
	CREATE INDEX IF NOT EXISTS idx_tenant_plugins_tenant ON tenant_plugins(tenant_id, is_enabled);
	`
	_, err := db.Exec(schema)
	return err
}

func seedDefaultData(db *sql.DB) error {
	defaultTenantID := "00000000-0000-0000-0000-000000000001"
	adminID := "00000000-0000-0000-0000-000000000002"

	// 1. Insert default tenant
	_, err := db.Exec(`
		INSERT OR IGNORE INTO tenants (id, slug, name, plan, is_active, created_at, updated_at)
		VALUES (?, 'master', 'Empresa Alfa Matriz', 'enterprise', 1, datetime('now'), datetime('now'))
	`, defaultTenantID)
	if err != nil {
		return err
	}

	// 2. Insert default admin user (password: admin)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES (?, ?, 'Admin', 'admin@hydravms.io', ?, 'admin', 1, datetime('now'), datetime('now'))
	`, adminID, defaultTenantID, string(hashedPassword))
	if err != nil {
		return err
	}

	// 3. Insert default cameras if none exist
	var camCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM cameras WHERE tenant_id = ?", defaultTenantID).Scan(&camCount)
	if camCount == 0 {
		sampleCameras := []struct {
			id, name, protocol, rtsp, sub, loc, status, res string
			fps                                              float64
			bitrate                                          int
		}{
			{"cam_entrance_01", "PORTARIA PRINCIPAL // ENTRADA", "ONVIF", "rtsp://127.0.0.1:8554/cam_entrance_01", "rtsp://127.0.0.1:8554/cam_entrance_01_sub", "Portaria Principal - Acesso A", "ONLINE", "1920x1080", 30.0, 4096},
			{"cam_parking_02", "ESTACIONAMENTO // SETOR B", "RTSP", "rtsp://127.0.0.1:8554/cam_parking_02", "rtsp://127.0.0.1:8554/cam_parking_02_sub", "Estacionamento VIP - Vagas 01-40", "ONLINE", "1920x1080", 25.0, 2048},
			{"cam_server_room_03", "DATA CENTER // CORREDOR FRIO", "RTSP", "rtsp://127.0.0.1:8554/cam_server_room_03", "rtsp://127.0.0.1:8554/cam_server_room_03_sub", "Data Center Rack A-01 a A-12", "ONLINE", "1920x1080", 30.0, 3072},
			{"cam_perimeter_04", "PERÍMETRO NORTE // MURO", "RTSP", "rtsp://127.0.0.1:8554/cam_perimeter_04", "rtsp://127.0.0.1:8554/cam_perimeter_04_sub", "Cerca Elétrica Perímetro Norte", "ONLINE", "1920x1080", 20.0, 2048},
		}

		for _, sc := range sampleCameras {
			_, _ = db.Exec(`
				INSERT OR IGNORE INTO cameras (
					id, tenant_id, name, protocol, rtsp_url, sub_stream_url,
					location, status, resolution, fps, bitrate_kbps, codec, is_active, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'H.264', 1, datetime('now'), datetime('now'))
			`, sc.id, defaultTenantID, sc.name, sc.protocol, sc.rtsp, sc.sub, sc.loc, sc.status, sc.res, sc.fps, sc.bitrate)
		}
	}

	// 4. Seed Official SOTA Object Detection Plugin
	_, _ = db.Exec(`
		INSERT OR IGNORE INTO plugins (
			id, name, version, author, category, runtime, entrypoint, min_vms_version,
			permissions, config_schema, ui_schema, is_official, is_deprecated,
			package_url, package_checksum, hardware_req, description, created_at, updated_at
		) VALUES (
			'object_detection_sota',
			'DETECCAO DE OBJETOS',
			'1.0.0',
			'Hydra Vision AI',
			'analytics',
			'binary_elf',
			'main',
			'1.0.0',
			'["video:shm_read", "events:write", "cuda:ipc"]',
			'{}',
			'{}',
			1,
			0,
			'https://huggingface.co/hydra-vision/object-detection',
			'sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
			'CUDA 13.3 // RTX 5090',
			'Deteccao e rastreamento SOTA de pessoas, celular, animais e veiculos com Auto-SAHI, Motion-Gating e Poligono/Linha de contagem.',
			datetime('now'),
			datetime('now')
		)
	`)

	return nil
}
