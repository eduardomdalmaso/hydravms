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
	repo := &InMemoryCameraRepository{
		cameras: make(map[string]*domain.Camera),
	}
	repo.seedInitialCameras()
	return repo
}

func (r *InMemoryCameraRepository) seedInitialCameras() {
	defaultTenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	initial := []*domain.Camera{
		{
			ID:          "cam_01",
			TenantID:    defaultTenantID,
			Name:        "CAM 01 // PORTARIA ENTRADA",
			Protocol:    domain.ProtocolRTSP,
			RTSPURL:     "rtsp://192.168.1.101:554/live",
			Status:      domain.CameraStatusOnline,
			Resolution:  "1920x1080",
			FPS:         30.0,
			BitrateKbps: 4096,
			Codec:       "H.265",
			HasPTZ:      true,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "cam_02",
			TenantID:    defaultTenantID,
			Name:        "CAM 02 // ESTACIONAMENTO VIP",
			Protocol:    domain.ProtocolONVIF,
			RTSPURL:     "rtsp://192.168.1.102:554/live",
			Status:      domain.CameraStatusOnline,
			Resolution:  "1920x1080",
			FPS:         30.0,
			BitrateKbps: 3072,
			Codec:       "H.265",
			HasPTZ:      false,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "cam_03",
			TenantID:    defaultTenantID,
			Name:        "CAM 03 // DOCAS DE CARGA",
			Protocol:    domain.ProtocolRTSP,
			RTSPURL:     "rtsp://192.168.1.103:554/live",
			Status:      domain.CameraStatusOnline,
			Resolution:  "2560x1440",
			FPS:         25.0,
			BitrateKbps: 6144,
			Codec:       "H.265",
			HasPTZ:      true,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "cam_04",
			TenantID:    defaultTenantID,
			Name:        "CAM 04 // PERIMETRO DOS FUNDOS",
			Protocol:    domain.ProtocolRTSP,
			RTSPURL:     "rtsp://192.168.1.104:554/live",
			Status:      domain.CameraStatusOnline,
			Resolution:  "1920x1080",
			FPS:         30.0,
			BitrateKbps: 3584,
			Codec:       "H.265",
			HasPTZ:      false,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, cam := range initial {
		r.cameras[cam.ID] = cam
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

	var result []*domain.Camera
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
	existing.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryCameraRepository) Delete(ctx context.Context, tenantID uuid.UUID, cameraID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cameras, cameraID)
	return nil
}
