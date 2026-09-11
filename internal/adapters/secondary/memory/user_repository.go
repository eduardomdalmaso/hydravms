package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"hydravms/internal/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		users: make(map[uuid.UUID]*domain.User),
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	repo.users[adminID] = &domain.User{
		ID:           adminID,
		TenantID:     tenantID,
		Name:         "Super Admin",
		Email:        "admin@hydravms.io",
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return repo
}

func (r *InMemoryUserRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanID := strings.ToLower(strings.TrimSpace(identifier))
	for _, u := range r.users {
		if strings.ToLower(u.Email) == cleanID ||
			strings.ToLower(u.Email) == cleanID+"@hydravms.io" ||
			(strings.EqualFold(u.Name, cleanID) && u.IsActive) {
			copy := *u
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.users[id]; ok {
		copy := *u
		return &copy, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.users[id]; ok {
		now := time.Now()
		u.LastLoginAt = &now
		u.UpdatedAt = now
	}
	return nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.User
	for _, u := range r.users {
		if u.TenantID == tenantID || u.TenantID == uuid.Nil {
			copy := *u
			list = append(list, &copy)
		}
	}
	return list, nil
}
