package domain

import (
	"time"

	"github.com/google/uuid"
)

// CameraRecordingProfile represents recording directives and weekly schedule for a camera.
type CameraRecordingProfile struct {
	ID               string    `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	CameraID         string    `json:"camera_id"`
	Name             string    `json:"name"`
	Mode             string    `json:"mode"` // 'continuous', 'motion', 'ai_event'
	SegmentDurationS int       `json:"segment_duration_s"`
	PreBufferS       int       `json:"pre_buffer_s"`
	PostBufferS      int       `json:"post_buffer_s"`
	RetentionDays    int       `json:"retention_days"`
	ScheduleJSON     string    `json:"schedule_json"` // JSON string representation of boolean[7][24]
	IsActive         bool      `json:"is_active"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RecordingSegment represents a stored video chunk on disk/storage.
type RecordingSegment struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	CameraID        string    `json:"camera_id"`
	RecordingMode   string    `json:"recording_mode"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationSeconds int       `json:"duration_seconds"`
	FileSizeBytes   int64     `json:"file_size_bytes"`
	S3Bucket        string    `json:"s3_bucket"`
	S3Key           string    `json:"s3_key"`
	IsPinned        bool      `json:"is_pinned"`
}
