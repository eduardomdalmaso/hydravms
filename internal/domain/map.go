package domain

import (
	"time"

	"github.com/google/uuid"
)

type MapType string

const (
	MapTypeImage MapType = "image"
	MapTypeGIS   MapType = "gis"
)

type MapPin struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	MapID          uuid.UUID `json:"map_id"`
	CameraID       string    `json:"camera_id,omitempty"`
	PinType        string    `json:"pin_type"` // 'camera', 'sensor', 'alarm_zone'
	PositionXPct   float64   `json:"position_x_pct"`
	PositionYPct   float64   `json:"position_y_pct"`
	GPSLat         *float64  `json:"gps_lat,omitempty"`
	GPSLng         *float64  `json:"gps_lng,omitempty"`
	FOVAngleDeg    int       `json:"fov_angle_deg"`
	FOVDirectionDeg int      `json:"fov_direction_deg"`
	FOVRadiusPx    int       `json:"fov_radius_px"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type InteractiveMap struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	FolderID      *uuid.UUID `json:"folder_id,omitempty"`
	Name          string    `json:"name"`
	MapType       MapType   `json:"map_type"`
	ImageS3Bucket string    `json:"image_s3_bucket,omitempty"`
	ImageS3Key    string    `json:"image_s3_key,omitempty"`
	CenterLat     *float64  `json:"center_lat,omitempty"`
	CenterLng     *float64  `json:"center_lng,omitempty"`
	ZoomLevel     int       `json:"zoom_level"`
	Pins          []MapPin  `json:"pins,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
