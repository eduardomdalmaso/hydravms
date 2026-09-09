package domain

import (
	"time"

	"github.com/google/uuid"
)

type LayoutSlot struct {
	SlotIndex  int    `json:"slot_index"`
	CameraID   string `json:"camera_id"`
	CameraName string `json:"camera_name,omitempty"`
	PTZLock    bool   `json:"ptz_lock"`
	Muted      bool   `json:"muted"`
}

// MosaicLayout represents a saved operator screen grid composition.
type MosaicLayout struct {
	ID          uuid.UUID    `json:"id"`
	TenantID    uuid.UUID    `json:"tenant_id"`
	FolderID    *uuid.UUID   `json:"folder_id,omitempty"`
	UserID      *uuid.UUID   `json:"user_id,omitempty"`
	Name        string       `json:"name"`
	GridType    string       `json:"grid_type"` // '1x1', '2x2', '3x3', '4x4', 'custom'
	IsShared    bool         `json:"is_shared"`
	SlotsConfig []LayoutSlot `json:"slots_config"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
