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

type SQLiteFolderRepository struct {
	db *sql.DB
}

func NewFolderRepository(db *sql.DB) *SQLiteFolderRepository {
	return &SQLiteFolderRepository{db: db}
}

func (r *SQLiteFolderRepository) Create(ctx context.Context, f *domain.Folder) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now

	var parentIDStr sql.NullString
	if f.ParentID != nil {
		parentIDStr = sql.NullString{String: f.ParentID.String(), Valid: true}
	}

	query := `
		INSERT INTO folders (id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		f.ID.String(), f.TenantID.String(), string(f.Module), parentIDStr, f.Name, f.ColorHex, f.Icon, f.SortOrder, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert folder: %w", err)
	}
	return nil
}

func (r *SQLiteFolderRepository) GetByID(ctx context.Context, tenantID, folderID uuid.UUID) (*domain.Folder, error) {
	query := `
		SELECT id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at
		FROM folders
		WHERE id = ? AND tenant_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, folderID.String(), tenantID.String())

	var f domain.Folder
	var idStr, tenantIDStr, modStr string
	var parentIDStr sql.NullString
	err := row.Scan(
		&idStr, &tenantIDStr, &modStr, &parentIDStr, &f.Name, &f.ColorHex, &f.Icon, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to query folder: %w", err)
	}
	f.ID, _ = uuid.Parse(idStr)
	f.TenantID, _ = uuid.Parse(tenantIDStr)
	f.Module = domain.FolderModule(modStr)
	if parentIDStr.Valid {
		pID, _ := uuid.Parse(parentIDStr.String)
		f.ParentID = &pID
	}
	return &f, nil
}

func (r *SQLiteFolderRepository) ListTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT 
				id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, created_at, updated_at,
				0 AS depth,
				name AS path_str
			FROM folders
			WHERE tenant_id = ?
			  AND module = ?
			  AND (parent_id IS NULL OR parent_id = '')

			UNION ALL

			SELECT 
				f.id, f.tenant_id, f.module, f.parent_id, f.name, f.color_hex, f.icon, f.sort_order, f.created_at, f.updated_at,
				ft.depth + 1 AS depth,
				ft.path_str || '/' || f.name AS path_str
			FROM folders f
			JOIN folder_tree ft ON f.parent_id = ft.id
			WHERE f.tenant_id = ?
			  AND f.module = ?
		)
		SELECT id, tenant_id, module, parent_id, name, color_hex, icon, sort_order, depth, path_str, created_at, updated_at
		FROM folder_tree
		ORDER BY sort_order, name
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), string(module), tenantID.String(), string(module))
	if err != nil {
		return nil, fmt.Errorf("failed to list folder tree: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.Folder, 0)
	for rows.Next() {
		var f domain.Folder
		var idStr, tenantIDStr, modStr, pathStr string
		var parentIDStr sql.NullString
		if err := rows.Scan(
			&idStr, &tenantIDStr, &modStr, &parentIDStr, &f.Name, &f.ColorHex, &f.Icon, &f.SortOrder,
			&f.Depth, &pathStr, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan folder row: %w", err)
		}
		f.ID, _ = uuid.Parse(idStr)
		f.TenantID, _ = uuid.Parse(tenantIDStr)
		f.Module = domain.FolderModule(modStr)
		if parentIDStr.Valid {
			pID, _ := uuid.Parse(parentIDStr.String)
			f.ParentID = &pID
		}
		f.Path = []string{pathStr}
		results = append(results, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *SQLiteFolderRepository) Update(ctx context.Context, f *domain.Folder) error {
	f.UpdatedAt = time.Now().UTC()
	var parentIDStr sql.NullString
	if f.ParentID != nil {
		parentIDStr = sql.NullString{String: f.ParentID.String(), Valid: true}
	}

	query := `
		UPDATE folders
		SET name = ?, color_hex = ?, icon = ?, parent_id = ?, sort_order = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`
	res, err := r.db.ExecContext(ctx, query, f.Name, f.ColorHex, f.Icon, parentIDStr, f.SortOrder, f.UpdatedAt, f.ID.String(), f.TenantID.String())
	if err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("folder not found or permission denied")
	}
	return nil
}

func (r *SQLiteFolderRepository) Delete(ctx context.Context, tenantID, folderID uuid.UUID) error {
	query := `DELETE FROM folders WHERE id = ? AND tenant_id = ?`
	res, err := r.db.ExecContext(ctx, query, folderID.String(), tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("folder not found or permission denied")
	}
	return nil
}
