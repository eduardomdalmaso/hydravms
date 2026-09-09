package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type CameraService struct {
	repo ports.CameraRepository
}

func NewCameraService(repo ports.CameraRepository) *CameraService {
	return &CameraService{repo: repo}
}

func (s *CameraService) CreateCamera(ctx context.Context, cam *domain.Camera) (*domain.Camera, error) {
	if strings.TrimSpace(cam.ID) == "" {
		return nil, fmt.Errorf("camera ID is required")
	}
	if strings.TrimSpace(cam.Name) == "" {
		return nil, fmt.Errorf("camera name is required")
	}
	if cam.Status == "" {
		cam.Status = domain.CameraStatusOnline
	}

	if err := s.repo.Create(ctx, cam); err != nil {
		return nil, err
	}

	return cam, nil
}

func (s *CameraService) ListCameras(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.Camera, error) {
	return s.repo.List(ctx, tenantID, folderID)
}

func (s *CameraService) GetCamera(ctx context.Context, tenantID uuid.UUID, cameraID string) (*domain.Camera, error) {
	return s.repo.GetByID(ctx, tenantID, cameraID)
}

func (s *CameraService) DeleteCamera(ctx context.Context, tenantID uuid.UUID, cameraID string) error {
	return s.repo.Delete(ctx, tenantID, cameraID)
}
