package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type MapRepository interface {
	Create(ctx context.Context, m *domain.InteractiveMap) error
	GetByID(ctx context.Context, tenantID, mapID uuid.UUID) (*domain.InteractiveMap, error)
	List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.InteractiveMap, error)
	Update(ctx context.Context, m *domain.InteractiveMap) error
	Delete(ctx context.Context, tenantID, mapID uuid.UUID) error
	SavePins(ctx context.Context, tenantID, mapID uuid.UUID, pins []domain.MapPin) error
}
