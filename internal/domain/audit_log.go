package domain

import (
	"errors"
	"time"
)

// AuditCategory categorizes events between administrative audit and system runtime events.
type AuditCategory string

// AuditLevel defines the severity level of the audit event.
type AuditLevel string

const (
	AuditCategorySystem AuditCategory = "SYSTEM"
	AuditCategoryAudit  AuditCategory = "AUDIT"

	AuditLevelInfo     AuditLevel = "INFO"
	AuditLevelWarning  AuditLevel = "WARNING"
	AuditLevelCritical AuditLevel = "CRITICAL"
)

// AuditLog represents a forensically verifiable audit record in the database.
type AuditLog struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenantId"`
	TenantName  string                 `json:"tenantName"`
	UserID      string                 `json:"userId,omitempty"`
	Actor       string                 `json:"actor"`
	IPAddress   string                 `json:"ipAddress"`
	Action      string                 `json:"action"`
	Category    AuditCategory          `json:"category"`
	Level       AuditLevel             `json:"level"`
	EntityType  string                 `json:"entityType"`
	EntityID    string                 `json:"entityId"`
	Target      string                 `json:"target"`
	Details     string                 `json:"details"`
	PayloadJSON map[string]interface{} `json:"payloadJson,omitempty"`
	CreatedAt   time.Time              `json:"timestamp"`
}

// Validate ensures core audit log invariants are satisfied.
func (a *AuditLog) Validate() error {
	if a.Action == "" {
		return errors.New("action is required")
	}
	if a.EntityType == "" {
		return errors.New("entity_type is required")
	}
	return nil
}
