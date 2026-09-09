package domain

import (
	"time"

	"github.com/google/uuid"
)

// PTZPreset represents a saved hardware pan/tilt/zoom position with snapshot thumbnail.
type PTZPreset struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	CameraID       string    `json:"camera_id"`
	PresetToken    string    `json:"preset_token"`
	Name           string    `json:"name"`
	PanCoord       *float64  `json:"pan_coord,omitempty"`
	TiltCoord      *float64  `json:"tilt_coord,omitempty"`
	ZoomCoord      *float64  `json:"zoom_coord,omitempty"`
	ThumbnailS3Key string    `json:"thumbnail_s3_key,omitempty"`
	IsHomePreset   bool      `json:"is_home_preset"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PTZPatrolPoint represents a single step in an automated PTZ guard tour.
type PTZPatrolPoint struct {
	ID                  uuid.UUID `json:"id"`
	PatrolID            uuid.UUID `json:"patrol_id"`
	PresetID            uuid.UUID `json:"preset_id"`
	PresetName          string    `json:"preset_name,omitempty"`
	StepOrder           int       `json:"step_order"`
	DwellTimeSeconds    int       `json:"dwell_time_seconds"`
	TransitionSpeedPct  int       `json:"transition_speed_pct"`
	CreatedAt           time.Time `json:"created_at"`
}

// PTZPatrol represents an automated guard tour sequence executed by the PTZ camera.
type PTZPatrol struct {
	ID        uuid.UUID        `json:"id"`
	TenantID  uuid.UUID        `json:"tenant_id"`
	CameraID  string           `json:"camera_id"`
	Name      string           `json:"name"`
	IsActive  bool             `json:"is_active"`
	LoopMode  bool             `json:"loop_mode"`
	Points    []PTZPatrolPoint `json:"points,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
