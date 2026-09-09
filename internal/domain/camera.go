package domain

import (
	"time"

	"github.com/google/uuid"
)

type CameraProtocol string

const (
	ProtocolRTSP  CameraProtocol = "rtsp"
	ProtocolONVIF CameraProtocol = "onvif"
	ProtocolRTMP  CameraProtocol = "rtmp"
)

type CameraStatus string

const (
	CameraStatusOnline    CameraStatus = "online"
	CameraStatusOffline   CameraStatus = "offline"
	CameraStatusRecording CameraStatus = "recording"
	CameraStatusError     CameraStatus = "error"
)

// Camera represents an edge video ingestion endpoint.
type Camera struct {
	ID             string         `json:"id"`
	TenantID       uuid.UUID      `json:"tenant_id"`
	FolderID       *uuid.UUID     `json:"folder_id,omitempty"`
	AssignedNodeID *uuid.UUID     `json:"assigned_node_id,omitempty"`
	Name           string         `json:"name"`
	Protocol       CameraProtocol `json:"protocol"`
	RTSPURL        string         `json:"rtsp_url,omitempty"`
	SubStreamURL   string         `json:"sub_stream_url,omitempty"`
	ONVIFIP        string         `json:"onvif_ip,omitempty"`
	ONVIFPort      int            `json:"onvif_port,omitempty"`
	ONVIFUser      string         `json:"onvif_user,omitempty"`
	ONVIFPass      string         `json:"-"`
	RTMPStreamKey  string         `json:"rtmp_stream_key,omitempty"`
	Location       string         `json:"location,omitempty"`
	Status         CameraStatus   `json:"status"`
	Resolution     string         `json:"resolution"`
	FPS            float64        `json:"fps"`
	BitrateKbps    int            `json:"bitrate_kbps"`
	Codec          string         `json:"codec"`
	HasPTZ         bool           `json:"has_ptz"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
