package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type RecordingRepository interface {
	SaveProfile(ctx context.Context, profile *domain.CameraRecordingProfile) error
	ListProfiles(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.CameraRecordingProfile, error)
	DeleteProfile(ctx context.Context, tenantID uuid.UUID, cameraID, profileID string) error
	ListRecordings(ctx context.Context, tenantID uuid.UUID, cameraID string, start, end time.Time) ([]*domain.RecordingSegment, error)
	InsertSegment(ctx context.Context, segment *domain.RecordingSegment) error
}
