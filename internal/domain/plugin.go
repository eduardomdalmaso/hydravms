package domain

import (
	"time"

	"github.com/google/uuid"
)

type PluginCategory string
type PluginRuntime string
type PluginStatus string

const (
	PluginStatusAvailable PluginStatus = "available"
	PluginStatusInstalled PluginStatus = "installed"
	PluginStatusRunning   PluginStatus = "running"
	PluginStatusStopped   PluginStatus = "stopped"
	PluginStatusUpdating  PluginStatus = "updating"
	PluginStatusCrashLoop PluginStatus = "crash_loop"
)

type Plugin struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Version         string                 `json:"version"`
	Author          string                 `json:"author"`
	Category        PluginCategory         `json:"category"`
	Runtime         PluginRuntime          `json:"runtime"`
	Entrypoint      string                 `json:"entrypoint"`
	MinVMSVersion   string                 `json:"min_vms_version"`
	Permissions     []string               `json:"permissions"`
	ConfigSchema    map[string]interface{} `json:"config_schema"`
	UISchema        map[string]interface{} `json:"ui_schema"`
	IsOfficial      bool                   `json:"is_official"`
	IsDeprecated    bool                   `json:"is_deprecated"`
	PackageURL      string                 `json:"package_url,omitempty"`
	PackageChecksum string                 `json:"package_checksum,omitempty"`
	HardwareReq     string                 `json:"hardware_req,omitempty"`
	Description     string                 `json:"description,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`

	// Tenant runtime state
	IsInstalled    bool         `json:"is_installed"`
	Status         PluginStatus `json:"status"`
	InstancesCount int          `json:"instances_count"`
	EventsCount    int          `json:"events_count"`
}

type TenantPlugin struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	PluginID          string                 `json:"plugin_id"`
	InstalledVersion  string                 `json:"installed_version"`
	PreviousVersion   string                 `json:"previous_version,omitempty"`
	LastStableVersion string                 `json:"last_stable_version,omitempty"`
	IsEnabled         bool                   `json:"is_enabled"`
	Status            PluginStatus           `json:"status"`
	PID               int                    `json:"pid,omitempty"`
	AssignedGPUDevice string                 `json:"assigned_gpu_device"`
	ConfigValues      map[string]interface{} `json:"config_values"`
	LastHealthCheck   *time.Time             `json:"last_health_check,omitempty"`
	ErrorMessage      string                 `json:"error_message,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}
