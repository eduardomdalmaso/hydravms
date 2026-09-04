-- ============================================================================
-- HYDRAVMS DATABASE SCHEMA MIGRATION v1
-- PostgreSQL Dialect with Multi-Tenant Scoping, Indices, and JSONB Support
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Tenants
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(128) NOT NULL,
    plan VARCHAR(32) NOT NULL DEFAULT 'standard',
    max_cameras INT NOT NULL DEFAULT 16,
    max_retention_days INT NOT NULL DEFAULT 30,
    max_storage_bytes BIGINT NOT NULL DEFAULT 1099511627776,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Users & RBAC
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'operator',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_tenant_user_email UNIQUE (tenant_id, email)
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- API Keys / Service Tokens (M2M Token-Only Access)
CREATE TABLE IF NOT EXISTS api_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    token_prefix VARCHAR(16) NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    scopes JSONB NOT NULL DEFAULT '["cameras:read", "events:read", "streams:read"]',
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Cameras (Support for ONVIF, RTSP, and RTMP)
CREATE TABLE IF NOT EXISTS cameras (
    id VARCHAR(64) NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    protocol VARCHAR(16) NOT NULL DEFAULT 'rtsp', -- 'rtsp', 'onvif', 'rtmp'
    rtsp_url TEXT,
    sub_stream_url TEXT,
    
    -- ONVIF Specific
    onvif_ip VARCHAR(64),
    onvif_port INT DEFAULT 80,
    onvif_user VARCHAR(64),
    onvif_pass VARCHAR(128),
    onvif_profile_token VARCHAR(64),
    
    -- RTMP Specific
    rtmp_stream_key VARCHAR(64) UNIQUE,
    
    location VARCHAR(128),
    status VARCHAR(32) NOT NULL DEFAULT 'offline',
    resolution VARCHAR(32) DEFAULT '1920x1080',
    fps NUMERIC(4,1) DEFAULT 30.0,
    bitrate_kbps INT DEFAULT 2048,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, id)
);

-- 4. Recording Profiles
CREATE TABLE IF NOT EXISTS camera_recording_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    mode VARCHAR(32) NOT NULL DEFAULT 'continuous',
    segment_duration_s INT NOT NULL DEFAULT 60,
    pre_buffer_s INT NOT NULL DEFAULT 5,
    post_buffer_s INT NOT NULL DEFAULT 10,
    retention_days INT NOT NULL DEFAULT 30,
    schedule_json JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_recording_profile_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE,
    CONSTRAINT uq_camera_recording_profile UNIQUE (tenant_id, camera_id)
);

-- 5. Camera Zones
CREATE TABLE IF NOT EXISTS camera_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    zone_type VARCHAR(32) NOT NULL DEFAULT 'intrusion',
    polygon_normalized JSONB NOT NULL,
    color_hex VARCHAR(9) DEFAULT '#00f0ff',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_zone_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE
);

-- 6. AI Rules
CREATE TABLE IF NOT EXISTS ai_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    zone_id UUID REFERENCES camera_zones(id) ON DELETE SET NULL,
    name VARCHAR(128) NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    target_classes JSONB NOT NULL,
    confidence_threshold NUMERIC(3,2) DEFAULT 0.65,
    min_duration_seconds INT DEFAULT 0,
    cooldown_seconds INT DEFAULT 15,
    sahi_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sahi_config JSONB DEFAULT '{"slice_size": 640, "overlap": 0.20, "nmm_threshold": 0.5}',
    actions JSONB NOT NULL DEFAULT '{"record_clip": true, "notify_ws": true, "webhook_url": null}',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_rule_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE
);

-- 7. Events (Gold Layer)
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    rule_id UUID REFERENCES ai_rules(id) ON DELETE SET NULL,
    zone_id UUID REFERENCES camera_zones(id) ON DELETE SET NULL,
    event_type VARCHAR(64) NOT NULL,
    severity VARCHAR(32) NOT NULL DEFAULT 'medium',
    status VARCHAR(32) NOT NULL DEFAULT 'new',
    triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    resolved_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    object_class VARCHAR(64) NOT NULL,
    confidence NUMERIC(3,2) NOT NULL,
    bbox_normalized JSONB NOT NULL,
    tracking_id BIGINT,
    snapshot_s3_key TEXT,
    crop_s3_key TEXT,
    clip_s3_key TEXT,
    metadata_s3_key TEXT,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_event_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE
);

-- Event Audit Actions
CREATE TABLE IF NOT EXISTS event_audit_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(64) NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. Recordings
CREATE TABLE IF NOT EXISTS recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    camera_id VARCHAR(64) NOT NULL,
    recording_mode VARCHAR(32) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    duration_seconds INT NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    s3_bucket VARCHAR(64) NOT NULL DEFAULT 'hydravms-recordings',
    s3_key TEXT NOT NULL,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_purged BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_recording_camera FOREIGN KEY (tenant_id, camera_id) REFERENCES cameras(tenant_id, id) ON DELETE CASCADE
);

-- 9. Storage Pools & Purge Logs
CREATE TABLE IF NOT EXISTS storage_pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    mount_path TEXT NOT NULL UNIQUE,
    total_bytes BIGINT NOT NULL,
    used_bytes BIGINT NOT NULL,
    available_bytes BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'healthy',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS storage_purge_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    storage_pool_id UUID REFERENCES storage_pools(id) ON DELETE SET NULL,
    bytes_reclaimed BIGINT NOT NULL,
    files_deleted_count INT NOT NULL,
    trigger_reason VARCHAR(64) NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 10. Video Exports
CREATE TABLE IF NOT EXISTS video_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    camera_id VARCHAR(64) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    export_type VARCHAR(32) NOT NULL DEFAULT 'zero_copy',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    file_size_bytes BIGINT,
    s3_key TEXT,
    error_message TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- 11. Relatórios Assíncronos (Report Jobs)
CREATE TABLE IF NOT EXISTS report_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    report_type VARCHAR(32) NOT NULL, -- 'events_executive_pdf', 'events_csv', 'events_xlsx', 'camera_uptime'
    format VARCHAR(16) NOT NULL,      -- 'pdf', 'csv', 'xlsx'
    filter_params JSONB NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    file_size_bytes BIGINT,
    s3_bucket VARCHAR(64) NOT NULL DEFAULT 'hydravms-reports',
    s3_key TEXT,
    error_message TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- 12. Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ip_address VARCHAR(45),
    action VARCHAR(128) NOT NULL,
    entity_type VARCHAR(64) NOT NULL,
    entity_id VARCHAR(64),
    payload_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 13. Índices de Alta Performance
CREATE INDEX IF NOT EXISTS idx_api_tokens_hash ON api_tokens (token_hash) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_api_tokens_tenant ON api_tokens (tenant_id);

CREATE INDEX IF NOT EXISTS idx_events_tenant_triggered ON events (tenant_id, triggered_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_tenant_camera ON events (tenant_id, camera_id, triggered_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_tenant_status ON events (tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_events_tenant_type ON events (tenant_id, event_type);
CREATE INDEX IF NOT EXISTS idx_events_pinned ON events (tenant_id, is_pinned) WHERE is_pinned = TRUE;

-- Índices Parciais para Performance de Dashboards
CREATE INDEX IF NOT EXISTS idx_events_unresolved ON events (tenant_id, triggered_at DESC) WHERE status = 'new';
CREATE INDEX IF NOT EXISTS idx_events_critical ON events (tenant_id, triggered_at DESC) WHERE severity = 'critical';

CREATE INDEX IF NOT EXISTS idx_recordings_timeline ON recordings (tenant_id, camera_id, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_recordings_purge_eval ON recordings (s3_bucket, is_pinned, is_purged, start_time ASC);

CREATE INDEX IF NOT EXISTS idx_report_jobs_tenant ON report_jobs (tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_zones_camera ON camera_zones (tenant_id, camera_id);
CREATE INDEX IF NOT EXISTS idx_rules_camera ON ai_rules (tenant_id, camera_id);

CREATE INDEX IF NOT EXISTS idx_events_bbox_gin ON events USING GIN (bbox_normalized);
CREATE INDEX IF NOT EXISTS idx_rules_classes_gin ON ai_rules USING GIN (target_classes);
CREATE INDEX IF NOT EXISTS idx_audit_payload_gin ON audit_logs USING GIN (payload_json);

-- 14. Catálogo Global de Plugins e Marketplace
CREATE TABLE IF NOT EXISTS plugins (
    id VARCHAR(64) PRIMARY KEY, -- ex: hydra-analytics-lpr
    name VARCHAR(128) NOT NULL,
    version VARCHAR(32) NOT NULL,
    author VARCHAR(128),
    category VARCHAR(64) NOT NULL DEFAULT 'analytics', -- analytics, device_driver, notification, storage
    runtime VARCHAR(32) NOT NULL DEFAULT 'python3', -- python3, binary_elf, wasm, docker
    entrypoint VARCHAR(255) NOT NULL,
    min_vms_version VARCHAR(32) NOT NULL DEFAULT '1.0.0',
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    config_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    ui_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_official BOOLEAN NOT NULL DEFAULT FALSE,
    is_deprecated BOOLEAN NOT NULL DEFAULT FALSE,
    package_url TEXT,
    package_checksum VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 15. Plugins Ativos por Tenant
CREATE TABLE IF NOT EXISTS tenant_plugins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plugin_id VARCHAR(64) NOT NULL REFERENCES plugins(id) ON DELETE CASCADE,
    installed_version VARCHAR(32) NOT NULL,
    previous_version VARCHAR(32),
    last_stable_version VARCHAR(32),
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(32) NOT NULL DEFAULT 'stopped', -- running, stopped, starting, crash_loop, error
    pid INT,
    assigned_gpu_device VARCHAR(32) DEFAULT '0',
    config_values JSONB NOT NULL DEFAULT '{}'::jsonb,
    last_health_check TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, plugin_id)
);

CREATE INDEX IF NOT EXISTS idx_tenant_plugins_tenant ON tenant_plugins (tenant_id, is_enabled);
CREATE INDEX IF NOT EXISTS idx_plugins_category ON plugins (category);


-- 16. Canais de Notificação (Telegram, Webhooks, Email, MQTT, WebSockets)
CREATE TABLE IF NOT EXISTS notification_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    channel_type VARCHAR(32) NOT NULL, -- telegram, webhook, email, mqtt, websocket
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    config_json JSONB NOT NULL DEFAULT '{}'::jsonb, -- bot_token, chat_id, endpoint_url, smtp_host, etc.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 17. Workflows de Eventos e Automação de Alertas
CREATE TABLE IF NOT EXISTS notification_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    trigger_type VARCHAR(64) NOT NULL DEFAULT 'ai_event', -- ai_event, camera_offline, storage_warning
    filter_conditions JSONB NOT NULL DEFAULT '{}'::jsonb, -- camera_ids, event_types, min_confidence, severity
    schedule_cron VARCHAR(64), -- ex: 0 20 * * * (somente à noite) ou NULL para 24/7
    cooldown_seconds INT NOT NULL DEFAULT 30, -- evita spam de alertas repetidos
    actions_pipeline JSONB NOT NULL DEFAULT '[]'::jsonb, -- lista ordenada de canais e transformações
    last_triggered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 18. Histórico e Auditoria de Disparos de Notificações
CREATE TABLE IF NOT EXISTS notification_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    workflow_id UUID REFERENCES notification_workflows(id) ON DELETE SET NULL,
    channel_id UUID REFERENCES notification_channels(id) ON DELETE SET NULL,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL DEFAULT 'sent', -- sent, failed, rate_limited
    response_code INT,
    error_message TEXT,
    payload_snapshot JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notification_channels_tenant ON notification_channels (tenant_id, channel_type);
CREATE INDEX IF NOT EXISTS idx_notification_workflows_tenant ON notification_workflows (tenant_id, is_enabled);
CREATE INDEX IF NOT EXISTS idx_notification_logs_tenant ON notification_logs (tenant_id, created_at DESC);
