package domain

import (
	"time"

	"github.com/google/uuid"
)

type NodeRole string

const (
	NodeRoleControlPlane NodeRole = "control_plane"
	NodeRoleEdgeIngest   NodeRole = "edge_ingest"
	NodeRoleGPUWorker    NodeRole = "gpu_worker"
)

type NodeStatus string

const (
	NodeStatusOnline   NodeStatus = "online"
	NodeStatusOffline  NodeStatus = "offline"
	NodeStatusDraining NodeStatus = "draining"
	NodeStatusError    NodeStatus = "error"
)

// ClusterNode represents a physical or virtual compute/ingest server in the distributed cluster.
type ClusterNode struct {
	ID                    uuid.UUID  `json:"id"`
	TenantID              *uuid.UUID `json:"tenant_id,omitempty"` // nil for shared cluster infrastructure
	NodeName              string     `json:"node_name"`
	NodeRole              NodeRole   `json:"node_role"`
	IPAddress             string     `json:"ip_address"`
	GRPCPort              int        `json:"grpc_port"`
	WebRTCPort            int        `json:"webrtc_port"`
	HTTPPort              int        `json:"http_port"`
	GPUDeviceInfo         string     `json:"gpu_device_info,omitempty"`
	CUDAComputeCapability string     `json:"cuda_compute_capability,omitempty"`
	CPUUsagePct           float64    `json:"cpu_usage_pct"`
	RAMUsagePct           float64    `json:"ram_usage_pct"`
	GPUUsagePct           float64    `json:"gpu_usage_pct"`
	VRAMUsedMB            int64      `json:"vram_used_mb"`
	ActiveStreamsCount    int        `json:"active_streams_count"`
	MaxStreamsCapacity    int        `json:"max_streams_capacity"`
	Status                NodeStatus `json:"status"`
	LastHeartbeatAt       time.Time  `json:"last_heartbeat_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// IsShared returns true if the cluster node is shared across all tenants.
func (n *ClusterNode) IsShared() bool {
	return n.TenantID == nil
}
