package domain

import (
	"time"

	"github.com/google/uuid"
)

type TourStep struct {
	ID               uuid.UUID  `json:"id"`
	TourID           uuid.UUID  `json:"tour_id"`
	StepOrder        int        `json:"step_order"`
	CameraID         string     `json:"camera_id,omitempty"`
	LayoutID         *uuid.UUID `json:"layout_id,omitempty"`
	DwellTimeSeconds int        `json:"dwell_time_seconds"`
	PTZPresetID      string     `json:"ptz_preset_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type VirtualTour struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	FolderID    *uuid.UUID `json:"folder_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	IsActive    bool       `json:"is_active"`
	LoopMode    bool       `json:"loop_mode"`
	Steps       []TourStep `json:"steps,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
