package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

// InMemoryEventRepository provides a thread-safe in-memory store for events.
type InMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[uuid.UUID]*domain.Event
}

// NewInMemoryEventRepository creates a new instance.
func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[uuid.UUID]*domain.Event),
	}
}

func (r *InMemoryEventRepository) Create(ctx context.Context, e *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.TriggeredAt.IsZero() {
		e.TriggeredAt = now
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	if e.Status == "" {
		e.Status = domain.EventStatusNew
	}

	copied := *e
	r.events[e.ID] = &copied
	return nil
}

func (r *InMemoryEventRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.events[id]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	copied := *e
	return &copied, nil
}

func (r *InMemoryEventRepository) List(ctx context.Context, tenantID uuid.UUID, cameraID string, limit int) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	var results []*domain.Event
	for _, e := range r.events {
		if cameraID != "" && e.CameraID != cameraID {
			continue
		}
		copied := *e
		results = append(results, &copied)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (r *InMemoryEventRepository) UpdateStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EventStatus, resolvedBy *uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.events[id]
	if !ok {
		return fmt.Errorf("event not found")
	}
	now := time.Now().UTC()
	e.Status = status
	e.ResolvedAt = &now
	e.ResolvedByUserID = resolvedBy
	return nil
}
