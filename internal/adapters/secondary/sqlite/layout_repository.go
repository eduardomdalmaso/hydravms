package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type SQLiteLayoutRepository struct {
	db *sql.DB
}

func NewLayoutRepository(db *sql.DB) *SQLiteLayoutRepository {
	return &SQLiteLayoutRepository{db: db}
}

func (r *SQLiteLayoutRepository) Create(ctx context.Context, l *domain.MosaicLayout) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	now := time.Now().UTC()
	l.CreatedAt = now
	l.UpdatedAt = now

	slotsJSON, _ := json.Marshal(l.SlotsConfig)
	usersJSON, _ := json.Marshal(l.AllowedUserIDs)

	var folderIDStr *string
	if l.FolderID != nil {
		str := l.FolderID.String()
		folderIDStr = &str
	}
	var userIDStr *string
	if l.UserID != nil {
		str := l.UserID.String()
		userIDStr = &str
	}

	query := `
		INSERT INTO layouts (
			id, tenant_id, folder_id, user_id, name, grid_type,
			is_locked, is_shared, target_monitor, allowed_user_ids,
			slots_config, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		l.ID.String(), l.TenantID.String(), folderIDStr, userIDStr, l.Name, l.GridType,
		l.IsLocked, l.IsShared, l.TargetMonitor, string(usersJSON),
		string(slotsJSON), l.CreatedAt, l.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create layout: %w", err)
	}
	return nil
}

func (r *SQLiteLayoutRepository) GetByID(ctx context.Context, tenantID, layoutID uuid.UUID) (*domain.MosaicLayout, error) {
	query := `
		SELECT 
			id, tenant_id, folder_id, user_id, name, grid_type,
			is_locked, is_shared, target_monitor, allowed_user_ids,
			slots_config, created_at, updated_at
		FROM layouts
		WHERE id = ? AND (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	row := r.db.QueryRowContext(ctx, query, layoutID.String(), tenantID.String())
	return r.scanLayout(row)
}

func (r *SQLiteLayoutRepository) List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.MosaicLayout, error) {
	query := `
		SELECT 
			id, tenant_id, folder_id, user_id, name, grid_type,
			is_locked, is_shared, target_monitor, allowed_user_ids,
			slots_config, created_at, updated_at
		FROM layouts
		WHERE (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	args := []interface{}{tenantID.String()}
	if folderID != nil {
		query += " AND folder_id = ?"
		args = append(args, folderID.String())
	}
	query += " ORDER BY updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list layouts: %w", err)
	}
	defer rows.Close()

	var results []*domain.MosaicLayout
	for rows.Next() {
		l, err := r.scanLayoutRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, l)
	}
	return results, nil
}

func (r *SQLiteLayoutRepository) Update(ctx context.Context, l *domain.MosaicLayout) error {
	l.UpdatedAt = time.Now().UTC()
	slotsJSON, _ := json.Marshal(l.SlotsConfig)
	usersJSON, _ := json.Marshal(l.AllowedUserIDs)

	var folderIDStr *string
	if l.FolderID != nil {
		str := l.FolderID.String()
		folderIDStr = &str
	}

	query := `
		UPDATE layouts SET
			name = ?, folder_id = ?, grid_type = ?, is_locked = ?,
			is_shared = ?, target_monitor = ?, allowed_user_ids = ?,
			slots_config = ?, updated_at = ?
		WHERE id = ? AND (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	_, err := r.db.ExecContext(ctx, query,
		l.Name, folderIDStr, l.GridType, l.IsLocked,
		l.IsShared, l.TargetMonitor, string(usersJSON),
		string(slotsJSON), l.UpdatedAt,
		l.ID.String(), l.TenantID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update layout: %w", err)
	}
	return nil
}

func (r *SQLiteLayoutRepository) Delete(ctx context.Context, tenantID, layoutID uuid.UUID) error {
	query := `DELETE FROM layouts WHERE id = ? AND (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')`
	_, err := r.db.ExecContext(ctx, query, layoutID.String(), tenantID.String())
	return err
}

func (r *SQLiteLayoutRepository) scanLayout(row *sql.Row) (*domain.MosaicLayout, error) {
	var l domain.MosaicLayout
	var idStr, tenantIDStr string
	var folderIDStr, userIDStr sql.NullString
	var usersJSON, slotsJSON string

	if err := row.Scan(
		&idStr, &tenantIDStr, &folderIDStr, &userIDStr, &l.Name, &l.GridType,
		&l.IsLocked, &l.IsShared, &l.TargetMonitor, &usersJSON,
		&slotsJSON, &l.CreatedAt, &l.UpdatedAt,
	); err != nil {
		return nil, err
	}
	l.ID, _ = uuid.Parse(idStr)
	l.TenantID, _ = uuid.Parse(tenantIDStr)
	if folderIDStr.Valid {
		fID, _ := uuid.Parse(folderIDStr.String)
		l.FolderID = &fID
	}
	if userIDStr.Valid {
		uID, _ := uuid.Parse(userIDStr.String)
		l.UserID = &uID
	}
	_ = json.Unmarshal([]byte(usersJSON), &l.AllowedUserIDs)
	_ = json.Unmarshal([]byte(slotsJSON), &l.SlotsConfig)
	return &l, nil
}

func (r *SQLiteLayoutRepository) scanLayoutRow(rows *sql.Rows) (*domain.MosaicLayout, error) {
	var l domain.MosaicLayout
	var idStr, tenantIDStr string
	var folderIDStr, userIDStr sql.NullString
	var usersJSON, slotsJSON string

	if err := rows.Scan(
		&idStr, &tenantIDStr, &folderIDStr, &userIDStr, &l.Name, &l.GridType,
		&l.IsLocked, &l.IsShared, &l.TargetMonitor, &usersJSON,
		&slotsJSON, &l.CreatedAt, &l.UpdatedAt,
	); err != nil {
		return nil, err
	}
	l.ID, _ = uuid.Parse(idStr)
	l.TenantID, _ = uuid.Parse(tenantIDStr)
	if folderIDStr.Valid {
		fID, _ := uuid.Parse(folderIDStr.String)
		l.FolderID = &fID
	}
	if userIDStr.Valid {
		uID, _ := uuid.Parse(userIDStr.String)
		l.UserID = &uID
	}
	_ = json.Unmarshal([]byte(usersJSON), &l.AllowedUserIDs)
	_ = json.Unmarshal([]byte(slotsJSON), &l.SlotsConfig)
	return &l, nil
}
