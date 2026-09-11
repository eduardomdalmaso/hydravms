package ports

import (
	"context"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type UserRepository interface {
	FindByEmailOrUsername(ctx context.Context, identifier string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, user *domain.User) error
	List(ctx context.Context, tenantID uuid.UUID) ([]*domain.User, error)
}
