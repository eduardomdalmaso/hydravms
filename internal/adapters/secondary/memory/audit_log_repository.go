package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type InMemoryAuditLogRepository struct {
	mu   sync.RWMutex
	logs []*domain.AuditLog
}

func NewInMemoryAuditLogRepository() *InMemoryAuditLogRepository {
	return &InMemoryAuditLogRepository{
		logs: []*domain.AuditLog{
			{
				ID:         "90000000-0000-0000-0000-000000000001",
				TenantID:   "00000000-0000-0000-0000-000000000001",
				TenantName: "EMPRESA ALFA // MATRIZ",
				Actor:      "SISTEMA // AUTOMATICO",
				IPAddress:  "127.0.0.1",
				Action:     "SYSTEM_BOOT",
				Category:   domain.AuditCategorySystem,
				Level:      domain.AuditLevelInfo,
				EntityType: "system",
				EntityID:   "core",
				Target:     "system // core",
				Details:    "HydraVMS Control Plane inicializado com sucesso",
				CreatedAt:  time.Now().UTC(),
			},
		},
	}
}

func (m *InMemoryAuditLogRepository) Save(ctx context.Context, log *domain.AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if log.ID == "" {
		log.ID = fmt.Sprintf("audit_%d", time.Now().UnixNano())
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	m.logs = append([]*domain.AuditLog{log}, m.logs...)
	return nil
}

func (m *InMemoryAuditLogRepository) List(ctx context.Context, tenantID string, category string, level string, limit int, offset int) ([]*domain.AuditLog, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []*domain.AuditLog
	for _, l := range m.logs {
		if tenantID != "" && l.TenantID != tenantID && tenantID != "tenant_alpha" {
			continue
		}
		if category != "" && category != "ALL" && string(l.Category) != category {
			continue
		}
		if level != "" && level != "ALL" && string(l.Level) != level {
			continue
		}
		filtered = append(filtered, l)
	}

	total := len(filtered)
	if offset >= total {
		return []*domain.AuditLog{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

func (m *InMemoryAuditLogRepository) Count(ctx context.Context, tenantID string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.logs), nil
}

var _ ports.AuditLogRepository = (*InMemoryAuditLogRepository)(nil)
