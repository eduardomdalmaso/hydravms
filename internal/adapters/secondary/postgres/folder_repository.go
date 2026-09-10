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

// PostgresFolderRepository handles SQL operations for folders in PostgreSQL.
type PostgresFolderRepository struct {
	pool *pgxpool.Pool
}

// NewFolderRepository creates a new instance of PostgresFolderRepository.
func NewFolderRepository(pool *pgxpool.Pool) *PostgresFolderRepository {
	return &PostgresFolderRepository{pool: pool}
}

func (r *PostgresFolderRepository) Create(ctx context.Context, f *domain.Folder) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	query := `
		INSERT INTO folders (id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		f.ID, f.TenantID, string(f.Module), f.ParentID, f.Name, f.ColorHex, f.Icon, f.SortOrder, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert folder: %w", err)
	}
	return nil
}

func (r *PostgresFolderRepository) GetByID(ctx context.Context, tenantID, folderID uuid.UUID) (*domain.Folder, error) {
	query := `
		SELECT id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at
		FROM folders
		WHERE id = $1 AND (tenant_id = $2 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	row := r.pool.QueryRow(ctx, query, folderID, tenantID)

	var f domain.Folder
	var modStr string
	err := row.Scan(
		&f.ID, &f.TenantID, &modStr, &f.ParentID, &f.Name, &f.ColorHex, &f.Icon, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to query folder: %w", err)
	}
	f.Module = domain.FolderModule(modStr)
	return &f, nil
}

func (r *PostgresFolderRepository) ListTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT 
				id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at,
				0 AS depth,
				ARRAY[name::text] AS path
			FROM folders
			WHERE (tenant_id = $1 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			  AND module = $2
			  AND parent_id IS NULL

			UNION ALL

			SELECT 
				f.id, f.tenant_id, f.module, f.parent_id, f.name, f.color_hex, f.icon, f.sort_order, f.created_at, f.updated_at,
				ft.depth + 1 AS depth,
				ft.path || f.name::text AS path
			FROM folders f
			JOIN folder_tree ft ON f.parent_id = ft.id
			WHERE (f.tenant_id = $1 OR f.tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			  AND f.module = $2
		)
		CYCLE id SET is_cycle USING path_cycle
		SELECT id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, depth, path, created_at, updated_at
		FROM folder_tree
		WHERE NOT is_cycle
		ORDER BY sort_order, name
	`

	rows, err := r.pool.Query(ctx, query, tenantID, string(module))
	if err != nil {
		return nil, fmt.Errorf("failed to list folder tree: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.Folder, 0)
	for rows.Next() {
		var f domain.Folder
		var modStr string
		var path []string
		if err := rows.Scan(
			&f.ID, &f.TenantID, &modStr, &f.ParentID, &f.Name, &f.ColorHex, &f.Icon, &f.SortOrder,
			&f.Depth, &path, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan folder row: %w", err)
		}
		f.Module = domain.FolderModule(modStr)
		f.Path = path
		results = append(results, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *PostgresFolderRepository) Update(ctx context.Context, f *domain.Folder) error {
	f.UpdatedAt = time.Now()
	query := `
		UPDATE folders
		SET name = $1, color_hex = $2, icon = $3, parent_id = $4, sort_order = $5, updated_at = $6
		WHERE id = $7 AND (tenant_id = $8 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	cmdTag, err := r.pool.Exec(ctx, query, f.Name, f.ColorHex, f.Icon, f.ParentID, f.SortOrder, f.UpdatedAt, f.ID, f.TenantID)
	if err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("folder not found or permission denied")
	}
	return nil
}

func (r *PostgresFolderRepository) Delete(ctx context.Context, tenantID, folderID uuid.UUID) error {
	query := `
		DELETE FROM folders
		WHERE id = $1 AND (tenant_id = $2 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	cmdTag, err := r.pool.Exec(ctx, query, folderID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("folder not found or permission denied")
	}
	return nil
}
