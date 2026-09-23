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

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE LOWER(email) = LOWER(?) 
		   OR LOWER(email) = LOWER(? || '@hydravms.io') 
		   OR (LOWER(name) = LOWER(?) AND is_active = 1)
		LIMIT 1
	`
	var u domain.User
	var idStr, tenantIDStr, roleStr string
	var lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx, q, identifier, identifier, identifier).Scan(
		&idStr, &tenantIDStr, &u.Name, &u.Email, &u.PasswordHash,
		&roleStr, &u.IsActive, &lastLogin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	u.ID, _ = uuid.Parse(idStr)
	u.TenantID, _ = uuid.Parse(tenantIDStr)
	u.Role = domain.UserRole(roleStr)
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	return &u, nil
}

func (r *SQLiteUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE id = ? 
		LIMIT 1
	`
	var u domain.User
	var idStr, tenantIDStr, roleStr string
	var lastLogin sql.NullTime
	err := r.db.QueryRowContext(ctx, q, id.String()).Scan(
		&idStr, &tenantIDStr, &u.Name, &u.Email, &u.PasswordHash,
		&roleStr, &u.IsActive, &lastLogin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	u.ID, _ = uuid.Parse(idStr)
	u.TenantID, _ = uuid.Parse(tenantIDStr)
	u.Role = domain.UserRole(roleStr)
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	return &u, nil
}

func (r *SQLiteUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE users SET last_login_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id.String())
	return err
}

func (r *SQLiteUserRepository) Create(ctx context.Context, u *domain.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	q := `
		INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, q, u.ID.String(), u.TenantID.String(), u.Name, u.Email, u.PasswordHash, string(u.Role), u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *SQLiteUserRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE tenant_id = ? 
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, tenantID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var idStr, tenantIDStr, roleStr string
		var lastLogin sql.NullTime
		if err := rows.Scan(
			&idStr, &tenantIDStr, &u.Name, &u.Email, &u.PasswordHash,
			&roleStr, &u.IsActive, &lastLogin, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.ID, _ = uuid.Parse(idStr)
		u.TenantID, _ = uuid.Parse(tenantIDStr)
		u.Role = domain.UserRole(roleStr)
		if lastLogin.Valid {
			u.LastLoginAt = &lastLogin.Time
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return users, nil
}
