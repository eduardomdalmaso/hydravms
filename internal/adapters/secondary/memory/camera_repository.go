package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type InMemoryCameraRepository struct {
	mu      sync.RWMutex
	cameras map[string]*domain.Camera
}

func NewInMemoryCameraRepository() *InMemoryCameraRepository {
	return &InMemoryCameraRepository{
		cameras: make(map[string]*domain.Camera),
	}
}

func (r *InMemoryCameraRepository) Create(ctx context.Context, camera *domain.Camera) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	camera.CreatedAt = time.Now()
	camera.UpdatedAt = time.Now()
	r.cameras[camera.ID] = camera
	return nil
}

func (r *InMemoryCameraRepository) GetByID(ctx context.Context, tenantID uuid.UUID, cameraID string) (*domain.Camera, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cam, exists := r.cameras[cameraID]
	if !exists {
		return nil, fmt.Errorf("camera not found")
	}
	return cam, nil
}

func (r *InMemoryCameraRepository) List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.Camera, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Camera, 0)
	for _, cam := range r.cameras {
		if folderID != nil {
			if cam.FolderID != nil && *cam.FolderID == *folderID {
				result = append(result, cam)
			}
		} else {
			result = append(result, cam)
		}
	}
	return result, nil
}

func (r *InMemoryCameraRepository) Update(ctx context.Context, camera *domain.Camera) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.cameras[camera.ID]
	if !exists {
		return fmt.Errorf("camera not found")
	}

	existing.Name = camera.Name
	existing.Protocol = camera.Protocol
	existing.RTSPURL = camera.RTSPURL
	existing.Resolution = camera.Resolution
	existing.FPS = camera.FPS
	existing.BitrateKbps = camera.BitrateKbps
	existing.Codec = camera.Codec
	existing.FolderID = camera.FolderID
	existing.Status = camera.Status
	existing.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryCameraRepository) Delete(ctx context.Context, tenantID uuid.UUID, cameraID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cameras, cameraID)
	return nil
}
