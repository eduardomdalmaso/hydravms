package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type RecordingService struct {
	repo ports.RecordingRepository
}

func NewRecordingService(repo ports.RecordingRepository) *RecordingService {
	return &RecordingService{repo: repo}
}

func (s *RecordingService) SaveProfile(ctx context.Context, profile *domain.CameraRecordingProfile) error {
	return s.repo.SaveProfile(ctx, profile)
}

func (s *RecordingService) ListProfiles(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.CameraRecordingProfile, error) {
	return s.repo.ListProfiles(ctx, tenantID, cameraID)
}

func (s *RecordingService) DeleteProfile(ctx context.Context, tenantID uuid.UUID, cameraID, profileID string) error {
	return s.repo.DeleteProfile(ctx, tenantID, cameraID, profileID)
}

func (s *RecordingService) ListRecordings(ctx context.Context, tenantID uuid.UUID, cameraID string, start, end time.Time) ([]*domain.RecordingSegment, error) {
	return s.repo.ListRecordings(ctx, tenantID, cameraID, start, end)
}
