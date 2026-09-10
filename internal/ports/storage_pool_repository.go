package ports

import (
	"context"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type StoragePoolRepository interface {
	Create(ctx context.Context, pool *domain.StoragePool) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.StoragePool, error)
	List(ctx context.Context, role *domain.StorageRole) ([]*domain.StoragePool, error)
	Update(ctx context.Context, pool *domain.StoragePool) error
	Delete(ctx context.Context, id uuid.UUID) error
}
