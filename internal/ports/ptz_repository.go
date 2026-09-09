package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type PTZRepository interface {
	// Presets
	CreatePreset(ctx context.Context, preset *domain.PTZPreset) error
	GetPresetByID(ctx context.Context, tenantID, presetID uuid.UUID) (*domain.PTZPreset, error)
	ListPresetsByCamera(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.PTZPreset, error)
	DeletePreset(ctx context.Context, tenantID, presetID uuid.UUID) error

	// Patrols
	CreatePatrol(ctx context.Context, patrol *domain.PTZPatrol) error
	GetPatrolByID(ctx context.Context, tenantID, patrolID uuid.UUID) (*domain.PTZPatrol, error)
	ListPatrolsByCamera(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.PTZPatrol, error)
	SavePatrolPoints(ctx context.Context, patrolID uuid.UUID, points []domain.PTZPatrolPoint) error
	DeletePatrol(ctx context.Context, tenantID, patrolID uuid.UUID) error
}
