package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type PluginRepository struct {
	db *sql.DB
}

func NewPluginRepository(db *sql.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

func (r *PluginRepository) ListPlugins(ctx context.Context, tenantID uuid.UUID) ([]*domain.Plugin, error) {
	query := `
		SELECT 
			p.id, p.name, p.version, p.author, p.category, p.runtime, p.entrypoint,
			p.min_vms_version, p.permissions, p.config_schema, p.ui_schema, p.is_official,
			p.is_deprecated, COALESCE(p.package_url, ''), COALESCE(p.package_checksum, ''),
			COALESCE(p.hardware_req, ''), COALESCE(p.description, ''), p.created_at, p.updated_at,
			CASE WHEN tp.id IS NOT NULL THEN 1 ELSE 0 END as is_installed,
			COALESCE(tp.status, 'available') as status
		FROM plugins p
		LEFT JOIN tenant_plugins tp ON tp.plugin_id = p.id AND tp.tenant_id = ?
		ORDER BY p.is_official DESC, p.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query plugins: %w", err)
	}
	defer rows.Close()

	var plugins []*domain.Plugin
	for rows.Next() {
		var (
			p                                    domain.Plugin
			permissionsJSON, configJSON, uiJSON string
			isInstalledInt                       int
			statusStr                            string
		)

		err := rows.Scan(
			&p.ID, &p.Name, &p.Version, &p.Author, &p.Category, &p.Runtime, &p.Entrypoint,
			&p.MinVMSVersion, &permissionsJSON, &configJSON, &uiJSON, &p.IsOfficial,
			&p.IsDeprecated, &p.PackageURL, &p.PackageChecksum,
			&p.HardwareReq, &p.Description, &p.CreatedAt, &p.UpdatedAt,
			&isInstalledInt, &statusStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plugin row: %w", err)
		}

		_ = json.Unmarshal([]byte(permissionsJSON), &p.Permissions)
		_ = json.Unmarshal([]byte(configJSON), &p.ConfigSchema)
		_ = json.Unmarshal([]byte(uiJSON), &p.UISchema)
		if p.Permissions == nil {
			p.Permissions = []string{}
		}
		if p.ConfigSchema == nil {
			p.ConfigSchema = make(map[string]interface{})
		}
		if p.UISchema == nil {
			p.UISchema = make(map[string]interface{})
		}

		p.IsInstalled = isInstalledInt == 1
		p.Status = domain.PluginStatus(statusStr)
		plugins = append(plugins, &p)
	}

	return plugins, nil
}

func (r *PluginRepository) GetPluginByID(ctx context.Context, tenantID uuid.UUID, id string) (*domain.Plugin, error) {
	query := `
		SELECT 
			p.id, p.name, p.version, p.author, p.category, p.runtime, p.entrypoint,
			p.min_vms_version, p.permissions, p.config_schema, p.ui_schema, p.is_official,
			p.is_deprecated, COALESCE(p.package_url, ''), COALESCE(p.package_checksum, ''),
			COALESCE(p.hardware_req, ''), COALESCE(p.description, ''), p.created_at, p.updated_at,
			CASE WHEN tp.id IS NOT NULL THEN 1 ELSE 0 END as is_installed,
			COALESCE(tp.status, 'available') as status
		FROM plugins p
		LEFT JOIN tenant_plugins tp ON tp.plugin_id = p.id AND tp.tenant_id = ?
		WHERE p.id = ?
	`

	var (
		p                                    domain.Plugin
		permissionsJSON, configJSON, uiJSON string
		isInstalledInt                       int
		statusStr                            string
	)

	err := r.db.QueryRowContext(ctx, query, tenantID.String(), id).Scan(
		&p.ID, &p.Name, &p.Version, &p.Author, &p.Category, &p.Runtime, &p.Entrypoint,
		&p.MinVMSVersion, &permissionsJSON, &configJSON, &uiJSON, &p.IsOfficial,
		&p.IsDeprecated, &p.PackageURL, &p.PackageChecksum,
		&p.HardwareReq, &p.Description, &p.CreatedAt, &p.UpdatedAt,
		&isInstalledInt, &statusStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plugin %s not found", id)
		}
		return nil, fmt.Errorf("failed to query plugin: %w", err)
	}

	_ = json.Unmarshal([]byte(permissionsJSON), &p.Permissions)
	_ = json.Unmarshal([]byte(configJSON), &p.ConfigSchema)
	_ = json.Unmarshal([]byte(uiJSON), &p.UISchema)
	if p.Permissions == nil {
		p.Permissions = []string{}
	}

	p.IsInstalled = isInstalledInt == 1
	p.Status = domain.PluginStatus(statusStr)
	return &p, nil
}

func (r *PluginRepository) InstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string, version string) error {
	newID := uuid.New().String()
	query := `
		INSERT INTO tenant_plugins (
			id, tenant_id, plugin_id, installed_version, is_enabled, status, config_values, created_at, updated_at
		) VALUES (?, ?, ?, ?, 1, 'running', '{}', datetime('now'), datetime('now'))
		ON CONFLICT(tenant_id, plugin_id) DO UPDATE SET
			installed_version = excluded.installed_version,
			is_enabled = 1,
			status = 'running',
			updated_at = datetime('now')
	`
	_, err := r.db.ExecContext(ctx, query, newID, tenantID.String(), pluginID, version)
	return err
}

func (r *PluginRepository) UninstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) error {
	query := `DELETE FROM tenant_plugins WHERE tenant_id = ? AND plugin_id = ?`
	_, err := r.db.ExecContext(ctx, query, tenantID.String(), pluginID)
	return err
}

func (r *PluginRepository) TogglePlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) (*domain.Plugin, error) {
	query := `
		UPDATE tenant_plugins
		SET status = CASE WHEN status = 'running' THEN 'stopped' ELSE 'running' END,
		    is_enabled = CASE WHEN status = 'running' THEN 0 ELSE 1 END,
		    updated_at = datetime('now')
		WHERE tenant_id = ? AND plugin_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, tenantID.String(), pluginID)
	if err != nil {
		return nil, err
	}
	return r.GetPluginByID(ctx, tenantID, pluginID)
}
