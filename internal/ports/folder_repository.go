package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

// FolderRepository defines driven port operations for hierarchical folders.
type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) error
	GetByID(ctx context.Context, tenantID, folderID uuid.UUID) (*domain.Folder, error)
	ListTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error)
	Update(ctx context.Context, folder *domain.Folder) error
	Delete(ctx context.Context, tenantID, folderID uuid.UUID) error
}
