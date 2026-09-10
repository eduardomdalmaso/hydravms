package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventCategory classifies the architectural domain of an event.
type EventCategory string

const (
	CategoryAIEvent      EventCategory = "AI_EVENT"
	CategorySystemEvent  EventCategory = "SYSTEM_EVENT"
	CategoryStorageEvent EventCategory = "STORAGE_EVENT"
	CategoryAuditEvent   EventCategory = "AUDIT_EVENT"
)

// Standard CloudEvents v1.0 Types
const (
	TypeCameraOffline         = "system.camera.offline"
	TypeCameraOnline          = "system.camera.online"
	TypeRecordingStarted      = "system.recording.started"
	TypeRecordingInterrupted  = "system.recording.interrupted"
	TypeStorageDiskPressure   = "system.storage.disk_pressure"
	TypeStoragePurgeCompleted = "system.storage.purge_completed"
	TypeHardwareGPUError      = "system.hardware.gpu_error"
	TypeAIMotionDetected      = "ai.detection.motion"
	TypeAIZoneIntrusion       = "ai.detection.zone_intrusion"
	TypeAILineCross           = "ai.detection.line_cross"
	TypeAILPRDetected         = "ai.detection.lpr"
)

// CloudEvent represents a CNCF CloudEvents v1.0 standard event envelope.
type CloudEvent struct {
	SpecVersion string        `json:"specversion"` // Always "1.0"
	ID          string        `json:"id"`          // Unique event UUID
	Source      string        `json:"source"`      // e.g. "hydravms/edge-01" or "hydrastream"
	Type        string        `json:"type"`        // Hierarchical event type
	Subject     string        `json:"subject"`     // Primary subject (camera_id, pool_id, user_id)
	Time        time.Time     `json:"time"`        // Timestamp UTC
	TenantID    uuid.UUID     `json:"tenant_id"`   // Multi-tenant isolation ID
	Category    EventCategory `json:"category"`   // AI_EVENT, SYSTEM_EVENT, etc.
	Severity    EventSeverity `json:"severity"`   // low, medium, warning, critical
	Data        interface{}   `json:"data"`        // Typed domain payload
}

// NewCloudEvent creates a standard CloudEvents v1.0 envelope.
func NewCloudEvent(
	tenantID uuid.UUID,
	eventType string,
	category EventCategory,
	severity EventSeverity,
	subject string,
	data interface{},
) *CloudEvent {
	return &CloudEvent{
		SpecVersion: "1.0",
		ID:          "evt_" + uuid.New().String(),
		Source:      "hydravms/control-plane",
		Type:        eventType,
		Subject:     subject,
		Time:        time.Now().UTC(),
		TenantID:    tenantID,
		Category:    category,
		Severity:    severity,
		Data:        data,
	}
}
