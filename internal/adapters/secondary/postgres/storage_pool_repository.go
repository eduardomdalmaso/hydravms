package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hydravms/internal/domain"
)

// PostgresStoragePoolRepository handles SQL operations for storage pools in PostgreSQL.
type PostgresStoragePoolRepository struct {
	pool *pgxpool.Pool
}

// NewStoragePoolRepository creates a new instance of PostgresStoragePoolRepository.
func NewStoragePoolRepository(pool *pgxpool.Pool) *PostgresStoragePoolRepository {
	return &PostgresStoragePoolRepository{pool: pool}
}

func (r *PostgresStoragePoolRepository) Create(ctx context.Context, p *domain.StoragePool) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.LastCheckedAt = now

	query := `
		INSERT INTO storage_pools (
			id, name, source_type, role, node_or_server, path_or_endpoint, filesystem,
			total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
			retention_days, last_checked_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17
		)
	`
	_, err := r.pool.Exec(ctx, query,
		p.ID, p.Name, string(p.SourceType), string(p.Role), p.NodeOrServer, p.PathOrEndpoint, p.Filesystem,
		p.TotalBytes, p.UsedBytes, p.AvailableBytes, string(p.Status), p.IsActive, p.IsSpilloverActive,
		p.RetentionDays, p.LastCheckedAt, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert storage pool: %w", err)
	}
	return nil
}

func (r *PostgresStoragePoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StoragePool, error) {
	query := `
		SELECT 
			id, name, source_type, role, node_or_server, COALESCE(path_or_endpoint, ''), filesystem,
			total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
			retention_days, last_checked_at, created_at, updated_at
		FROM storage_pools
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var p domain.StoragePool
	var srcStr, roleStr, statusStr string
	err := row.Scan(
		&p.ID, &p.Name, &srcStr, &roleStr, &p.NodeOrServer, &p.PathOrEndpoint, &p.Filesystem,
		&p.TotalBytes, &p.UsedBytes, &p.AvailableBytes, &statusStr, &p.IsActive, &p.IsSpilloverActive,
		&p.RetentionDays, &p.LastCheckedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("storage pool not found")
		}
		return nil, fmt.Errorf("failed to query storage pool: %w", err)
	}
	p.SourceType = domain.StorageSourceType(srcStr)
	p.Role = domain.StorageRole(roleStr)
	p.Status = domain.StorageStatus(statusStr)
	return &p, nil
}

func (r *PostgresStoragePoolRepository) List(ctx context.Context, role *domain.StorageRole) ([]*domain.StoragePool, error) {
	var query string
	var args []any

	if role != nil {
		query = `
			SELECT 
				id, name, source_type, role, node_or_server, COALESCE(path_or_endpoint, ''), filesystem,
				total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
				retention_days, last_checked_at, created_at, updated_at
			FROM storage_pools
			WHERE role = $1 AND is_active = TRUE
			ORDER BY created_at ASC
		`
		args = []any{string(*role)}
	} else {
		query = `
			SELECT 
				id, name, source_type, role, node_or_server, COALESCE(path_or_endpoint, ''), filesystem,
				total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
				retention_days, last_checked_at, created_at, updated_at
			FROM storage_pools
			WHERE is_active = TRUE
			ORDER BY created_at ASC
		`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list storage pools: %w", err)
	}
	defer rows.Close()

	var results []*domain.StoragePool
	for rows.Next() {
		var p domain.StoragePool
		var srcStr, roleStr, statusStr string
		if err := rows.Scan(
			&p.ID, &p.Name, &srcStr, &roleStr, &p.NodeOrServer, &p.PathOrEndpoint, &p.Filesystem,
			&p.TotalBytes, &p.UsedBytes, &p.AvailableBytes, &statusStr, &p.IsActive, &p.IsSpilloverActive,
			&p.RetentionDays, &p.LastCheckedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan storage pool row: %w", err)
		}
		p.SourceType = domain.StorageSourceType(srcStr)
		p.Role = domain.StorageRole(roleStr)
		p.Status = domain.StorageStatus(statusStr)
		results = append(results, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *PostgresStoragePoolRepository) Update(ctx context.Context, p *domain.StoragePool) error {
	p.UpdatedAt = time.Now()
	query := `
		UPDATE storage_pools
		SET 
			name = $1, source_type = $2, role = $3, node_or_server = $4,
			path_or_endpoint = $5, filesystem = $6, total_bytes = $7,
			used_bytes = $8, available_bytes = $9, status = $10,
			is_active = $11, is_spillover_active = $12, retention_days = $13,
			last_checked_at = $14, updated_at = $15
		WHERE id = $16
	`
	cmdTag, err := r.pool.Exec(ctx, query,
		p.Name, string(p.SourceType), string(p.Role), p.NodeOrServer,
		p.PathOrEndpoint, p.Filesystem, p.TotalBytes,
		p.UsedBytes, p.AvailableBytes, string(p.Status),
		p.IsActive, p.IsSpilloverActive, p.RetentionDays,
		p.LastCheckedAt, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update storage pool: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("storage pool not found")
	}
	return nil
}

func (r *PostgresStoragePoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM storage_pools WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete storage pool: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("storage pool not found")
	}
	return nil
}
