package domain

import (
	"time"

	"github.com/google/uuid"
)

// FolderModule represents the category of resources grouped by a folder.
type FolderModule string

const (
	ModuleCameras   FolderModule = "cameras"
	ModuleLayouts   FolderModule = "layouts"
	ModuleMaps      FolderModule = "maps"
	ModuleTours     FolderModule = "tours"
	ModuleWorkflows FolderModule = "workflows"
	ModuleUsers     FolderModule = "users"
)

// Folder represents a unified hierarchical organizational node.
type Folder struct {
	ID        uuid.UUID    `json:"id"`
	TenantID  uuid.UUID    `json:"tenant_id"`
	Module    FolderModule `json:"module"`
	ParentID  *uuid.UUID   `json:"parent_id,omitempty"`
	Name      string       `json:"name"`
	ColorHex  string       `json:"color_hex"`
	Icon      string       `json:"icon"`
	SortOrder int          `json:"sort_order"`
	Depth     int          `json:"depth,omitempty"`
	Path      []string     `json:"path,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// IsRoot returns true if the folder has no parent.
func (f *Folder) IsRoot() bool {
	return f.ParentID == nil
}
