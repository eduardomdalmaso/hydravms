package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/secondary/postgres"
)

type SQLiteClusterNodeRepository struct {
	db *sql.DB
}

func NewClusterNodeRepository(db *sql.DB) *SQLiteClusterNodeRepository {
	return &SQLiteClusterNodeRepository{db: db}
}

func (r *SQLiteClusterNodeRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*postgres.ClusterNode, error) {
	q := `
		SELECT 
			id, 
			COALESCE(tenant_id, ?), 
			node_name, 
			node_role, 
			ip_address, 
			grpc_port, 
			webrtc_port, 
			http_port, 
			COALESCE(gpu_device_info, ''), 
			cpu_usage_pct, 
			ram_usage_pct, 
			gpu_usage_pct, 
			vram_used_mb, 
			active_streams_count, 
			status, 
			created_at 
		FROM cluster_nodes 
		WHERE tenant_id IS NULL OR tenant_id = '' OR tenant_id = ? 
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, q, tenantID.String(), tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query cluster nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*postgres.ClusterNode
	for rows.Next() {
		var n postgres.ClusterNode
		var idStr, tenantIDStr string
		var vramUsedInt int64
		if err := rows.Scan(
			&idStr, &tenantIDStr, &n.NodeName, &n.NodeRole, &n.IPAddress,
			&n.GRPCPort, &n.WebRTCPort, &n.HTTPPort, &n.GPUInfo,
			&n.CPUUsagePct, &n.RAMUsagePct, &n.GPUUsagePct, &vramUsedInt,
			&n.ActiveStreams, &n.Status, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.TenantID, _ = uuid.Parse(tenantIDStr)
		n.VRAMUsedMB = float64(vramUsedInt)
		nodes = append(nodes, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return nodes, nil
}

func (r *SQLiteClusterNodeRepository) Create(ctx context.Context, n *postgres.ClusterNode) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	n.CreatedAt = time.Now().UTC()

	var existingIDStr string
	checkQuery := `SELECT id FROM cluster_nodes WHERE (tenant_id = ? OR tenant_id IS NULL OR tenant_id = '') AND ip_address = ? AND http_port = ? LIMIT 1`
	err := r.db.QueryRowContext(ctx, checkQuery, n.TenantID.String(), n.IPAddress, n.HTTPPort).Scan(&existingIDStr)
	if err == nil && existingIDStr != "" {
		existingID, _ := uuid.Parse(existingIDStr)
		n.ID = existingID
		updateQuery := `
			UPDATE cluster_nodes 
			SET node_name = ?, node_role = ?, grpc_port = ?, webrtc_port = ?, gpu_device_info = ?, status = ?, updated_at = datetime('now') 
			WHERE id = ?
		`
		_, err = r.db.ExecContext(ctx, updateQuery, n.NodeName, n.NodeRole, n.GRPCPort, n.WebRTCPort, n.GPUInfo, n.Status, existingIDStr)
		return err
	}

	q := `
		INSERT INTO cluster_nodes (id, tenant_id, node_name, node_role, ip_address, grpc_port, webrtc_port, http_port, gpu_device_info, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, q, n.ID.String(), n.TenantID.String(), n.NodeName, n.NodeRole, n.IPAddress, n.GRPCPort, n.WebRTCPort, n.HTTPPort, n.GPUInfo, n.Status, n.CreatedAt, n.CreatedAt)
	return err
}

func (r *SQLiteClusterNodeRepository) UpdateMetrics(ctx context.Context, id uuid.UUID, status string, cpuPct, ramPct, gpuPct, vramUsedMB float64, activeStreams int) error {
	q := `
		UPDATE cluster_nodes 
		SET status = ?, cpu_usage_pct = ?, ram_usage_pct = ?, gpu_usage_pct = ?, vram_used_mb = ?, active_streams_count = ?, last_heartbeat_at = datetime('now'), updated_at = datetime('now') 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, q, status, cpuPct, ramPct, gpuPct, int64(vramUsedMB), activeStreams, id.String())
	return err
}

func (r *SQLiteClusterNodeRepository) Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error {
	q := `DELETE FROM cluster_nodes WHERE id = ? AND (tenant_id = ? OR tenant_id IS NULL OR tenant_id = '')`
	_, err := r.db.ExecContext(ctx, q, id.String(), tenantID.String())
	return err
}
