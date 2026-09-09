package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventSeverity string

const (
	SeverityLow      EventSeverity = "low"
	SeverityMedium   EventSeverity = "medium"
	SeverityWarning  EventSeverity = "warning"
	SeverityCritical EventSeverity = "critical"
)

type EventStatus string

const (
	EventStatusNew          EventStatus = "new"
	EventStatusAcknowledged EventStatus = "acknowledged"
	EventStatusResolved     EventStatus = "resolved"
	EventStatusFalsePositive EventStatus = "false_positive"
)

type BoundingBox struct {
	XCenter float64 `json:"x_center"`
	YCenter float64 `json:"y_center"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
}

// Event represents a processed AI detection, smart alarm, or security incident (Gold Layer).
type Event struct {
	ID               uuid.UUID     `json:"id"`
	TenantID         uuid.UUID     `json:"tenant_id"`
	CameraID         string        `json:"camera_id"`
	CameraName       string        `json:"camera_name,omitempty"`
	RuleID           *uuid.UUID    `json:"rule_id,omitempty"`
	ZoneID           *uuid.UUID    `json:"zone_id,omitempty"`
	EventType        string        `json:"event_type"`
	Severity         EventSeverity `json:"severity"`
	Status           EventStatus   `json:"status"`
	TriggeredAt      time.Time     `json:"triggered_at"`
	ResolvedAt       *time.Time    `json:"resolved_at,omitempty"`
	ResolvedByUserID *uuid.UUID    `json:"resolved_by_user_id,omitempty"`
	ObjectClass      string        `json:"object_class"`
	Confidence       float64       `json:"confidence"`
	BBoxNormalized   BoundingBox   `json:"bbox_normalized"`
	TrackingID       *int64        `json:"tracking_id,omitempty"`
	SnapshotS3Key    string        `json:"snapshot_s3_key,omitempty"`
	CropS3Key        string        `json:"crop_s3_key,omitempty"`
	ClipS3Key        string        `json:"clip_s3_key,omitempty"`
	IsPinned         bool          `json:"is_pinned"`
	Notes            string        `json:"notes,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
}
