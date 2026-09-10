package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

// EventRepository defines persistence operations for system and AI security events.
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Event, error)
	List(ctx context.Context, tenantID uuid.UUID, cameraID string, limit int) ([]*domain.Event, error)
	UpdateStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EventStatus, resolvedBy *uuid.UUID) error
}
