package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type SQLiteStoragePoolRepository struct {
	db *sql.DB
}

func NewStoragePoolRepository(db *sql.DB) *SQLiteStoragePoolRepository {
	return &SQLiteStoragePoolRepository{db: db}
}

func (r *SQLiteStoragePoolRepository) Create(ctx context.Context, p *domain.StoragePool) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.LastCheckedAt = now

	query := `
		INSERT INTO storage_pools (
			id, name, source_type, role, node_or_server, path_or_endpoint, filesystem,
			total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
			retention_days, last_checked_at, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID.String(), p.Name, string(p.SourceType), string(p.Role), p.NodeOrServer, p.PathOrEndpoint, p.Filesystem,
		p.TotalBytes, p.UsedBytes, p.AvailableBytes, string(p.Status), p.IsActive, p.IsSpilloverActive,
		p.RetentionDays, p.LastCheckedAt, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert storage pool: %w", err)
	}
	return nil
}

func (r *SQLiteStoragePoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StoragePool, error) {
	query := `
		SELECT 
			id, name, source_type, role, node_or_server, COALESCE(path_or_endpoint, ''), filesystem,
			total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
			retention_days, last_checked_at, created_at, updated_at
		FROM storage_pools
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var p domain.StoragePool
	var idStr, srcStr, roleStr, statusStr string
	err := row.Scan(
		&idStr, &p.Name, &srcStr, &roleStr, &p.NodeOrServer, &p.PathOrEndpoint, &p.Filesystem,
		&p.TotalBytes, &p.UsedBytes, &p.AvailableBytes, &statusStr, &p.IsActive, &p.IsSpilloverActive,
		&p.RetentionDays, &p.LastCheckedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("storage pool not found")
		}
		return nil, fmt.Errorf("failed to query storage pool: %w", err)
	}
	p.ID, _ = uuid.Parse(idStr)
	p.SourceType = domain.StorageSourceType(srcStr)
	p.Role = domain.StorageRole(roleStr)
	p.Status = domain.StorageStatus(statusStr)
	return &p, nil
}

func (r *SQLiteStoragePoolRepository) List(ctx context.Context, role *domain.StorageRole) ([]*domain.StoragePool, error) {
	var query string
	var args []any

	if role != nil {
		query = `
			SELECT 
				id, name, source_type, role, node_or_server, COALESCE(path_or_endpoint, ''), filesystem,
				total_bytes, used_bytes, available_bytes, status, is_active, is_spillover_active,
				retention_days, last_checked_at, created_at, updated_at
			FROM storage_pools
			WHERE role = ? AND is_active = 1
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
			WHERE is_active = 1
			ORDER BY created_at ASC
		`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list storage pools: %w", err)
	}
	defer rows.Close()

	var results []*domain.StoragePool
	for rows.Next() {
		var p domain.StoragePool
		var idStr, srcStr, roleStr, statusStr string
		if err := rows.Scan(
			&idStr, &p.Name, &srcStr, &roleStr, &p.NodeOrServer, &p.PathOrEndpoint, &p.Filesystem,
			&p.TotalBytes, &p.UsedBytes, &p.AvailableBytes, &statusStr, &p.IsActive, &p.IsSpilloverActive,
			&p.RetentionDays, &p.LastCheckedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan storage pool row: %w", err)
		}
		p.ID, _ = uuid.Parse(idStr)
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

func (r *SQLiteStoragePoolRepository) Update(ctx context.Context, p *domain.StoragePool) error {
	p.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE storage_pools
		SET 
			name = ?, source_type = ?, role = ?, node_or_server = ?,
			path_or_endpoint = ?, filesystem = ?, total_bytes = ?,
			used_bytes = ?, available_bytes = ?, status = ?,
			is_active = ?, is_spillover_active = ?, retention_days = ?,
			last_checked_at = ?, updated_at = ?
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query,
		p.Name, string(p.SourceType), string(p.Role), p.NodeOrServer,
		p.PathOrEndpoint, p.Filesystem, p.TotalBytes,
		p.UsedBytes, p.AvailableBytes, string(p.Status),
		p.IsActive, p.IsSpilloverActive, p.RetentionDays,
		p.LastCheckedAt, p.UpdatedAt, p.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update storage pool: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("storage pool not found")
	}
	return nil
}

func (r *SQLiteStoragePoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM storage_pools WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete storage pool: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("storage pool not found")
	}
	return nil
}
