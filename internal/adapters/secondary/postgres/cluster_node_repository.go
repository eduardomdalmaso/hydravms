package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClusterNode struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	NodeName      string    `json:"node_name"`
	NodeRole      string    `json:"node_role"`
	IPAddress     string    `json:"ip_address"`
	GRPCPort      int       `json:"grpc_port"`
	WebRTCPort    int       `json:"webrtc_port"`
	HTTPPort      int       `json:"http_port"`
	GPUInfo       string    `json:"gpu_device_info"`
	CPUModel      string    `json:"cpu_model"`
	CPUUsagePct   float64   `json:"cpu_usage_pct"`
	RAMUsedGB     float64   `json:"ram_used_gb"`
	RAMTotalGB    float64   `json:"ram_total_gb"`
	RAMUsagePct   float64   `json:"ram_usage_pct"`
	GPUUsagePct   float64   `json:"gpu_usage_pct"`
	VRAMUsedMB    float64   `json:"vram_used_mb"`
	VRAMTotalMB   float64   `json:"vram_total_mb"`
	TempCelsius   float64   `json:"temp_celsius"`
	PowerWatts    float64   `json:"power_watts"`
	ActiveStreams int       `json:"active_streams_count"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostgresClusterNodeRepository struct {
	pool *pgxpool.Pool
}

func NewClusterNodeRepository(pool *pgxpool.Pool) *PostgresClusterNodeRepository {
	return &PostgresClusterNodeRepository{pool: pool}
}

func (r *PostgresClusterNodeRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*ClusterNode, error) {
	q := `
		SELECT 
			id, 
			COALESCE(tenant_id, $1), 
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
		WHERE tenant_id IS NULL OR tenant_id = $1 
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query cluster nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*ClusterNode
	for rows.Next() {
		var n ClusterNode
		var vramUsedInt int64
		if err := rows.Scan(
			&n.ID, &n.TenantID, &n.NodeName, &n.NodeRole, &n.IPAddress,
			&n.GRPCPort, &n.WebRTCPort, &n.HTTPPort, &n.GPUInfo,
			&n.CPUUsagePct, &n.RAMUsagePct, &n.GPUUsagePct, &vramUsedInt,
			&n.ActiveStreams, &n.Status, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		n.VRAMUsedMB = float64(vramUsedInt)
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (r *PostgresClusterNodeRepository) Create(ctx context.Context, n *ClusterNode) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	n.CreatedAt = time.Now()

	// Check if already exists by IP and HTTP port
	var existingID uuid.UUID
	checkQuery := `SELECT id FROM cluster_nodes WHERE (tenant_id = $1 OR tenant_id IS NULL) AND ip_address = $2 AND http_port = $3 LIMIT 1`
	err := r.pool.QueryRow(ctx, checkQuery, n.TenantID, n.IPAddress, n.HTTPPort).Scan(&existingID)
	if err == nil && existingID != uuid.Nil {
		n.ID = existingID
		updateQuery := `
			UPDATE cluster_nodes 
			SET node_name = $1, node_role = $2, grpc_port = $3, webrtc_port = $4, gpu_device_info = $5, status = $6, updated_at = NOW() 
			WHERE id = $7
		`
		_, err = r.pool.Exec(ctx, updateQuery, n.NodeName, n.NodeRole, n.GRPCPort, n.WebRTCPort, n.GPUInfo, n.Status, existingID)
		return err
	}

	q := `
		INSERT INTO cluster_nodes (id, tenant_id, node_name, node_role, ip_address, grpc_port, webrtc_port, http_port, gpu_device_info, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
	`
	_, err = r.pool.Exec(ctx, q, n.ID, n.TenantID, n.NodeName, n.NodeRole, n.IPAddress, n.GRPCPort, n.WebRTCPort, n.HTTPPort, n.GPUInfo, n.Status, n.CreatedAt)
	return err
}

func (r *PostgresClusterNodeRepository) UpdateMetrics(ctx context.Context, id uuid.UUID, status string, cpuPct, ramPct, gpuPct, vramUsedMB float64, activeStreams int) error {
	q := `
		UPDATE cluster_nodes 
		SET status = $1, cpu_usage_pct = $2, ram_usage_pct = $3, gpu_usage_pct = $4, vram_used_mb = $5, active_streams_count = $6, last_heartbeat_at = NOW(), updated_at = NOW() 
		WHERE id = $7
	`
	_, err := r.pool.Exec(ctx, q, status, cpuPct, ramPct, gpuPct, int64(vramUsedMB), activeStreams, id)
	return err
}

func (r *PostgresClusterNodeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM cluster_nodes WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}
