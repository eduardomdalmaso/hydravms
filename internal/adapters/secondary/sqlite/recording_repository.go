package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type SQLiteRecordingRepository struct {
	db *sql.DB
}

func NewRecordingRepository(db *sql.DB) *SQLiteRecordingRepository {
	return &SQLiteRecordingRepository{db: db}
}

func (r *SQLiteRecordingRepository) SaveProfile(ctx context.Context, p *domain.CameraRecordingProfile) error {
	p.UpdatedAt = time.Now().UTC()
	var profileUUID uuid.UUID
	if parsed, err := uuid.Parse(p.ID); err == nil {
		profileUUID = parsed
	} else {
		profileUUID = uuid.New()
		p.ID = profileUUID.String()
	}

	if p.PreBufferS <= 0 {
		p.PreBufferS = 5
	}
	if p.PostBufferS <= 0 {
		p.PostBufferS = 10
	}
	if p.RetentionDays <= 0 {
		p.RetentionDays = 30
	}
	if p.SegmentDurationS <= 0 {
		p.SegmentDurationS = 60
	}
	if p.Name == "" {
		p.Name = "Perfil de Gravacao"
	}
	if p.ScheduleJSON == "" {
		p.ScheduleJSON = "[]"
	}

	query := `
		INSERT INTO camera_recording_profiles (
			id, tenant_id, camera_id, name, mode, segment_duration_s,
			pre_buffer_s, post_buffer_s, retention_days, schedule_json,
			is_active, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?
		)
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			mode = excluded.mode,
			segment_duration_s = excluded.segment_duration_s,
			pre_buffer_s = excluded.pre_buffer_s,
			post_buffer_s = excluded.post_buffer_s,
			retention_days = excluded.retention_days,
			schedule_json = excluded.schedule_json,
			is_active = excluded.is_active,
			updated_at = excluded.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		profileUUID.String(), p.TenantID.String(), p.CameraID, p.Name, p.Mode, p.SegmentDurationS,
		p.PreBufferS, p.PostBufferS, p.RetentionDays, p.ScheduleJSON,
		p.IsActive, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save recording profile: %w", err)
	}
	return nil
}

func (r *SQLiteRecordingRepository) ListProfiles(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.CameraRecordingProfile, error) {
	query := `
		SELECT 
			id, tenant_id, camera_id, name, mode, segment_duration_s,
			pre_buffer_s, post_buffer_s, retention_days, COALESCE(schedule_json, '[]'),
			is_active, updated_at
		FROM camera_recording_profiles
		WHERE (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
		  AND camera_id = ?
		ORDER BY updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), cameraID)
	if err != nil {
		return nil, fmt.Errorf("failed to list recording profiles: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.CameraRecordingProfile, 0)
	for rows.Next() {
		var p domain.CameraRecordingProfile
		var idStr, tenantIDStr string
		if err := rows.Scan(
			&idStr, &tenantIDStr, &p.CameraID, &p.Name, &p.Mode, &p.SegmentDurationS,
			&p.PreBufferS, &p.PostBufferS, &p.RetentionDays, &p.ScheduleJSON,
			&p.IsActive, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recording profile: %w", err)
		}
		p.ID = idStr
		p.TenantID, _ = uuid.Parse(tenantIDStr)
		results = append(results, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return results, nil
}

func (r *SQLiteRecordingRepository) DeleteProfile(ctx context.Context, tenantID uuid.UUID, cameraID string, profileID string) error {
	query := `
		DELETE FROM camera_recording_profiles
		WHERE id = ? AND camera_id = ? AND (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	res, err := r.db.ExecContext(ctx, query, profileID, cameraID, tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to delete recording profile: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("recording profile not found")
	}
	return nil
}

func (r *SQLiteRecordingRepository) ListRecordings(ctx context.Context, tenantID uuid.UUID, cameraID string, start, end time.Time) ([]*domain.RecordingSegment, error) {
	query := `
		SELECT 
			id, tenant_id, camera_id, recording_mode, start_time, end_time,
			duration_seconds, file_size_bytes, s3_bucket, s3_key, is_pinned
		FROM recordings
		WHERE (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
		  AND camera_id = ?
		  AND end_time >= ?
		  AND start_time <= ?
		  AND is_purged = 0
		ORDER BY start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID.String(), cameraID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query recordings: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.RecordingSegment, 0)
	for rows.Next() {
		var s domain.RecordingSegment
		var idStr, tenantIDStr string
		if err := rows.Scan(
			&idStr, &tenantIDStr, &s.CameraID, &s.RecordingMode, &s.StartTime, &s.EndTime,
			&s.DurationSeconds, &s.FileSizeBytes, &s.S3Bucket, &s.S3Key, &s.IsPinned,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recording segment: %w", err)
		}
		s.ID, _ = uuid.Parse(idStr)
		s.TenantID, _ = uuid.Parse(tenantIDStr)
		results = append(results, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return results, nil
}

func (r *SQLiteRecordingRepository) InsertSegment(ctx context.Context, s *domain.RecordingSegment) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.TenantID == uuid.Nil {
		s.TenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	if s.StartTime.IsZero() {
		s.StartTime = time.Now().Add(-time.Duration(s.DurationSeconds) * time.Second)
	}
	if s.EndTime.IsZero() {
		s.EndTime = time.Now()
	}
	if s.DurationSeconds <= 0 {
		s.DurationSeconds = int(s.EndTime.Sub(s.StartTime).Seconds())
	}
	if s.RecordingMode == "" {
		s.RecordingMode = "motion"
	}

	query := `
		INSERT OR IGNORE INTO recordings (
			id, tenant_id, camera_id, recording_mode, start_time, end_time,
			duration_seconds, file_size_bytes, s3_bucket, s3_key, is_pinned, is_purged
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, 0
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID.String(), s.TenantID.String(), s.CameraID, s.RecordingMode, s.StartTime, s.EndTime,
		s.DurationSeconds, s.FileSizeBytes, s.S3Bucket, s.S3Key, s.IsPinned,
	)
	return err
}
