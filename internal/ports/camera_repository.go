package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type CameraRepository interface {
	Create(ctx context.Context, camera *domain.Camera) error
	GetByID(ctx context.Context, tenantID uuid.UUID, cameraID string) (*domain.Camera, error)
	List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.Camera, error)
	Update(ctx context.Context, camera *domain.Camera) error
	Delete(ctx context.Context, tenantID uuid.UUID, cameraID string) error
}
