package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"hydravms/internal/domain"
)

type PostgresRecordingRepository struct {
	pool *pgxpool.Pool
}

func NewRecordingRepository(pool *pgxpool.Pool) *PostgresRecordingRepository {
	return &PostgresRecordingRepository{pool: pool}
}

func (r *PostgresRecordingRepository) SaveProfile(ctx context.Context, p *domain.CameraRecordingProfile) error {
	p.UpdatedAt = time.Now()
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
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10::jsonb,
			$11, $12
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			mode = EXCLUDED.mode,
			segment_duration_s = EXCLUDED.segment_duration_s,
			pre_buffer_s = EXCLUDED.pre_buffer_s,
			post_buffer_s = EXCLUDED.post_buffer_s,
			retention_days = EXCLUDED.retention_days,
			schedule_json = EXCLUDED.schedule_json,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.pool.Exec(ctx, query,
		profileUUID, p.TenantID, p.CameraID, p.Name, p.Mode, p.SegmentDurationS,
		p.PreBufferS, p.PostBufferS, p.RetentionDays, p.ScheduleJSON,
		p.IsActive, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save recording profile: %w", err)
	}
	return nil
}

func (r *PostgresRecordingRepository) ListProfiles(ctx context.Context, tenantID uuid.UUID, cameraID string) ([]*domain.CameraRecordingProfile, error) {
	query := `
		SELECT 
			id, tenant_id, camera_id, name, mode, segment_duration_s,
			pre_buffer_s, post_buffer_s, retention_days, COALESCE(schedule_json::text, '[]'),
			is_active, updated_at
		FROM camera_recording_profiles
		WHERE (tenant_id = $1 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
		  AND camera_id = $2
		ORDER BY updated_at DESC
	`
	rows, err := r.pool.Query(ctx, query, tenantID, cameraID)
	if err != nil {
		return nil, fmt.Errorf("failed to list recording profiles: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.CameraRecordingProfile, 0)
	for rows.Next() {
		var p domain.CameraRecordingProfile
		var profUUID uuid.UUID
		if err := rows.Scan(
			&profUUID, &p.TenantID, &p.CameraID, &p.Name, &p.Mode, &p.SegmentDurationS,
			&p.PreBufferS, &p.PostBufferS, &p.RetentionDays, &p.ScheduleJSON,
			&p.IsActive, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recording profile: %w", err)
		}
		p.ID = profUUID.String()
		results = append(results, &p)
	}
	return results, nil
}

func (r *PostgresRecordingRepository) DeleteProfile(ctx context.Context, tenantID uuid.UUID, cameraID string, profileID string) error {
	profUUID, err := uuid.Parse(profileID)
	if err != nil {
		return fmt.Errorf("invalid profile id: %w", err)
	}

	query := `
		DELETE FROM camera_recording_profiles
		WHERE id = $1 AND camera_id = $2 AND (tenant_id = $3 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	cmdTag, err := r.pool.Exec(ctx, query, profUUID, cameraID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete recording profile: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("recording profile not found")
	}
	return nil
}

func (r *PostgresRecordingRepository) ListRecordings(ctx context.Context, tenantID uuid.UUID, cameraID string, start, end time.Time) ([]*domain.RecordingSegment, error) {
	query := `
		SELECT 
			id, tenant_id, camera_id, recording_mode, start_time, end_time,
			duration_seconds, file_size_bytes, s3_bucket, s3_key, is_pinned
		FROM recordings
		WHERE (tenant_id = $1 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
		  AND camera_id = $2
		  AND end_time >= $3
		  AND start_time <= $4
		  AND is_purged = FALSE
		ORDER BY start_time ASC
	`
	rows, err := r.pool.Query(ctx, query, tenantID, cameraID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query recordings: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.RecordingSegment, 0)
	for rows.Next() {
		var s domain.RecordingSegment
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.CameraID, &s.RecordingMode, &s.StartTime, &s.EndTime,
			&s.DurationSeconds, &s.FileSizeBytes, &s.S3Bucket, &s.S3Key, &s.IsPinned,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recording segment: %w", err)
		}
		results = append(results, &s)
	}
	return results, nil
}

func (r *PostgresRecordingRepository) InsertSegment(ctx context.Context, s *domain.RecordingSegment) error {
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
		INSERT INTO recordings (
			id, tenant_id, camera_id, recording_mode, start_time, end_time,
			duration_seconds, file_size_bytes, s3_bucket, s3_key, is_pinned, is_purged
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, FALSE
		)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query,
		s.ID, s.TenantID, s.CameraID, s.RecordingMode, s.StartTime, s.EndTime,
		s.DurationSeconds, s.FileSizeBytes, s.S3Bucket, s.S3Key, s.IsPinned,
	)
	return err
}
