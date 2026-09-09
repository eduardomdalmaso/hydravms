package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type LayoutRepository interface {
	Create(ctx context.Context, layout *domain.MosaicLayout) error
	GetByID(ctx context.Context, tenantID, layoutID uuid.UUID) (*domain.MosaicLayout, error)
	List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.MosaicLayout, error)
	Update(ctx context.Context, layout *domain.MosaicLayout) error
	Delete(ctx context.Context, tenantID, layoutID uuid.UUID) error
}
