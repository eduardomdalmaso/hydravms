package ports

import (
	"context"

	"hydravms/internal/domain"
)

// AuditLogRepository defines the secondary port for audit log persistence.
type AuditLogRepository interface {
	Save(ctx context.Context, log *domain.AuditLog) error
	List(ctx context.Context, tenantID string, category string, level string, limit int, offset int) ([]*domain.AuditLog, int, error)
	Count(ctx context.Context, tenantID string) (int, error)
}
