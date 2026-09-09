package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type TourRepository interface {
	Create(ctx context.Context, tour *domain.VirtualTour) error
	GetByID(ctx context.Context, tenantID, tourID uuid.UUID) (*domain.VirtualTour, error)
	List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.VirtualTour, error)
	Update(ctx context.Context, tour *domain.VirtualTour) error
	Delete(ctx context.Context, tenantID, tourID uuid.UUID) error
}
