package application

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		cam.ID = fmt.Sprintf("cam_%d", time.Now().UnixMilli())
	}
	if strings.TrimSpace(cam.Name) == "" {
		return nil, fmt.Errorf("camera name is required")
	}
	if cam.Status == "" {
		cam.Status = domain.CameraStatusOnline
	}
	if cam.Protocol == "" {
		cam.Protocol = domain.ProtocolRTSP
	}
	if cam.Codec == "" {
		cam.Codec = "H.265"
	}
	if cam.Resolution == "" {
		cam.Resolution = "1920x1080"
	}
	if cam.FPS <= 0 {
		cam.FPS = 30.0
	}
	if cam.BitrateKbps <= 0 {
		cam.BitrateKbps = 4096
	}
	cam.IsActive = true

	if err := s.repo.Create(ctx, cam); err != nil {
		return nil, err
	}

	return cam, nil
}

func (s *CameraService) UpdateCamera(ctx context.Context, cam *domain.Camera) (*domain.Camera, error) {
	if strings.TrimSpace(cam.ID) == "" {
		return nil, fmt.Errorf("camera ID is required")
	}
	existing, err := s.repo.GetByID(ctx, cam.TenantID, cam.ID)
	if err != nil {
		return nil, fmt.Errorf("camera not found: %w", err)
	}
	if strings.TrimSpace(cam.Name) != "" {
		existing.Name = cam.Name
	}
	if cam.RTSPURL != "" {
		existing.RTSPURL = cam.RTSPURL
	}
	if cam.SubStreamURL != "" {
		existing.SubStreamURL = cam.SubStreamURL
	}
	if cam.Protocol != "" {
		existing.Protocol = cam.Protocol
	}
	if cam.Codec != "" {
		existing.Codec = cam.Codec
	}
	if cam.Resolution != "" {
		existing.Resolution = cam.Resolution
	}
	if cam.FPS > 0 {
		existing.FPS = cam.FPS
	}
	if cam.BitrateKbps > 0 {
		existing.BitrateKbps = cam.BitrateKbps
	}
	if cam.Location != "" {
		existing.Location = cam.Location
	}
	if cam.ONVIFIP != "" {
		existing.ONVIFIP = cam.ONVIFIP
	}
	if cam.ONVIFPort > 0 {
		existing.ONVIFPort = cam.ONVIFPort
	}
	if cam.ONVIFUser != "" {
		existing.ONVIFUser = cam.ONVIFUser
	}
	if cam.ONVIFPass != "" {
		existing.ONVIFPass = cam.ONVIFPass
	}
	if cam.Brand != "" {
		existing.Brand = cam.Brand
	}
	if cam.Model != "" {
		existing.Model = cam.Model
	}
	if cam.SerialNumber != "" {
		existing.SerialNumber = cam.SerialNumber
	}
	if cam.Firmware != "" {
		existing.Firmware = cam.Firmware
	}
	if cam.MACAddress != "" {
		existing.MACAddress = cam.MACAddress
	}
	if cam.Latitude != 0 {
		existing.Latitude = cam.Latitude
	}
	if cam.Longitude != 0 {
		existing.Longitude = cam.Longitude
	}
	existing.HasPTZ = cam.HasPTZ
	if cam.FolderID != nil {
		existing.FolderID = cam.FolderID
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
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

func (s *CameraService) UpdateCameraFolder(ctx context.Context, tenantID uuid.UUID, cameraID string, folderID *uuid.UUID) error {
	cam, err := s.repo.GetByID(ctx, tenantID, cameraID)
	if err != nil {
		return err
	}
	cam.FolderID = folderID
	return s.repo.Update(ctx, cam)
}
