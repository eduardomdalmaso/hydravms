package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

// InMemoryFolderRepository provides thread-safe folder persistence with initial seed data.
type InMemoryFolderRepository struct {
	mu      sync.RWMutex
	folders map[uuid.UUID]*domain.Folder
}

func NewInMemoryFolderRepository() *InMemoryFolderRepository {
	return &InMemoryFolderRepository{
		folders: make(map[uuid.UUID]*domain.Folder),
	}
}

func (r *InMemoryFolderRepository) Create(ctx context.Context, folder *domain.Folder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if folder.ID == uuid.Nil {
		folder.ID = uuid.New()
	}
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	r.folders[folder.ID] = folder
	return nil
}

func (r *InMemoryFolderRepository) GetByID(ctx context.Context, tenantID, folderID uuid.UUID) (*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	folder, exists := r.folders[folderID]
	if !exists || (folder.TenantID != tenantID && folder.TenantID != uuid.MustParse("00000000-0000-0000-0000-000000000001")) {
		return nil, fmt.Errorf("folder not found")
	}

	return folder, nil
}

func (r *InMemoryFolderRepository) ListTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Folder, 0)
	for _, f := range r.folders {
		if (f.TenantID == tenantID || f.TenantID == uuid.MustParse("00000000-0000-0000-0000-000000000001")) && f.Module == module {
			result = append(result, f)
		}
	}

	return result, nil
}

func (r *InMemoryFolderRepository) Update(ctx context.Context, folder *domain.Folder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.folders[folder.ID]
	if !exists {
		return fmt.Errorf("folder not found")
	}

	existing.Name = folder.Name
	existing.ColorHex = folder.ColorHex
	existing.Icon = folder.Icon
	existing.ParentID = folder.ParentID
	existing.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryFolderRepository) Delete(ctx context.Context, tenantID, folderID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.folders, folderID)
	return nil
}
