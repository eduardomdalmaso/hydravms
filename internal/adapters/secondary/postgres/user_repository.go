package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"hydravms/internal/domain"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	repo := &PostgresUserRepository{pool: pool}
	// Bootstrap default superadmin user with admin/admin if needed
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = repo.EnsureAdminUser(ctx)
	}()
	return repo
}

func (r *PostgresUserRepository) EnsureAdminUser(ctx context.Context) error {
	defaultTenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Ensure default tenant exists
	_, _ = r.pool.Exec(ctx, `
		INSERT INTO tenants (id, slug, name, plan, is_active, created_at, updated_at) 
		VALUES ($1, 'master', 'Empresa Alfa Matriz', 'enterprise', true, NOW(), NOW()) 
		ON CONFLICT (id) DO NOTHING
	`, defaultTenantID)

	// Generate standard bcrypt hash for default password 'admin'
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert or update default superadmin
	upsertQuery := `
		INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, 'Super Admin', 'admin@hydravms.io', $3, 'admin', true, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE 
		SET password_hash = $3, is_active = true, role = 'admin', updated_at = NOW()
	`
	_, err = r.pool.Exec(ctx, upsertQuery, adminID, defaultTenantID, string(hashedPassword))
	return err
}

func (r *PostgresUserRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE LOWER(email) = LOWER($1) 
		   OR LOWER(email) = LOWER($1 || '@hydravms.io') 
		   OR (name ILIKE $1 AND is_active = true)
		LIMIT 1
	`
	var u domain.User
	var roleStr string
	err := r.pool.QueryRow(ctx, q, identifier).Scan(
		&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash,
		&roleStr, &u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	u.Role = domain.UserRole(roleStr)
	return &u, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE id = $1 
		LIMIT 1
	`
	var u domain.User
	var roleStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash,
		&roleStr, &u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	u.Role = domain.UserRole(roleStr)
	return &u, nil
}

func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *domain.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	q := `
		INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, q, u.ID, u.TenantID, u.Name, u.Email, u.PasswordHash, string(u.Role), u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.User, error) {
	q := `
		SELECT id, tenant_id, name, email, password_hash, role, is_active, last_login_at, created_at, updated_at 
		FROM users 
		WHERE tenant_id = $1 OR tenant_id IS NULL 
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var roleStr string
		if err := rows.Scan(
			&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash,
			&roleStr, &u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.Role = domain.UserRole(roleStr)
		users = append(users, &u)
	}
	return users, nil
}
