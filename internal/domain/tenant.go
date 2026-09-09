package domain

import (
	"time"

	"github.com/google/uuid"
)

type TenantPlan string

const (
	PlanStandard   TenantPlan = "standard"
	PlanEnterprise TenantPlan = "enterprise"
	PlanForensic   TenantPlan = "forensic"
)

type Tenant struct {
	ID               uuid.UUID  `json:"id"`
	Slug             string     `json:"slug"`
	Name             string     `json:"name"`
	Plan             TenantPlan `json:"plan"`
	MaxCameras       int        `json:"max_cameras"`
	MaxRetentionDays int        `json:"max_retention_days"`
	MaxStorageBytes  int64      `json:"max_storage_bytes"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
