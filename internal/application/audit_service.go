package application

import (
	"context"
	"fmt"
	"time"

	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

// AuditService coordinates forensic logging and auditing queries.
type AuditService struct {
	repo ports.AuditLogRepository
}

func NewAuditService(repo ports.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Record(ctx context.Context, log *domain.AuditLog) error {
	if s.repo == nil {
		return nil
	}
	if err := log.Validate(); err != nil {
		return err
	}
	return s.repo.Save(ctx, log)
}

func (s *AuditService) RecordAction(
	ctx context.Context,
	tenantID, actor, ipAddress, action, entityType, entityID, details string,
	category domain.AuditCategory,
	level domain.AuditLevel,
	payload map[string]interface{},
) error {
	if s.repo == nil {
		return nil
	}

	if category == "" {
		category = domain.AuditCategoryAudit
	}
	if level == "" {
		level = domain.AuditLevelInfo
	}
	if actor == "" {
		actor = "operador_master"
	}
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}

	target := fmt.Sprintf("%s // %s", entityType, entityID)
	if entityID == "" {
		target = entityType
	}

	if payload == nil {
		payload = make(map[string]interface{})
	}
	payload["message"] = details
	payload["level"] = string(level)

	log := &domain.AuditLog{
		TenantID:    tenantID,
		Actor:       actor,
		IPAddress:   ipAddress,
		Action:      action,
		Category:    category,
		Level:       level,
		EntityType:  entityType,
		EntityID:    entityID,
		Target:      target,
		Details:     details,
		PayloadJSON: payload,
		CreatedAt:   time.Now().UTC(),
	}

	return s.repo.Save(ctx, log)
}

func (s *AuditService) ListLogs(ctx context.Context, tenantID, category, level string, limit, offset int) ([]*domain.AuditLog, int, error) {
	if s.repo == nil {
		return []*domain.AuditLog{}, 0, nil
	}
	return s.repo.List(ctx, tenantID, category, level, limit, offset)
}
