package domain

import (
	"time"

	"github.com/google/uuid"
)

type StorageSourceType string

const (
	StorageSourceLocalDisk   StorageSourceType = "LOCAL_DISK"
	StorageSourceNetworkNAS  StorageSourceType = "NETWORK_NAS"
	StorageSourceObjectS3    StorageSourceType = "OBJECT_S3"
)

type StorageRole string

const (
	StorageRoleHotBuffer   StorageRole = "HOT_BUFFER"
	StorageRoleWarmArchive StorageRole = "WARM_ARCHIVE"
	StorageRoleSnapshots   StorageRole = "SNAPSHOTS"
	StorageRoleDatabase    StorageRole = "DATABASE"
)

type StorageStatus string

const (
	StorageStatusOnline      StorageStatus = "ONLINE"
	StorageStatusStandby     StorageStatus = "STANDBY"
	StorageStatusMaintenance StorageStatus = "MAINTENANCE"
)

// StoragePool represents a managed physical disk, NAS share, or MinIO S3 bucket tier.
type StoragePool struct {
	ID                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	SourceType        StorageSourceType `json:"source_type"`
	Role              StorageRole       `json:"role"`
	NodeOrServer      string            `json:"node_or_server"`
	PathOrEndpoint    string            `json:"path_or_endpoint"`
	Filesystem        string            `json:"filesystem"`
	TotalBytes        int64             `json:"total_bytes"`
	UsedBytes         int64             `json:"used_bytes"`
	AvailableBytes    int64             `json:"available_bytes"`
	Status            StorageStatus     `json:"status"`
	IsActive          bool              `json:"is_active"`
	IsSpilloverActive bool              `json:"is_spillover_active"`
	RetentionDays     int               `json:"retention_days"`
	LastCheckedAt     time.Time         `json:"last_checked_at"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
