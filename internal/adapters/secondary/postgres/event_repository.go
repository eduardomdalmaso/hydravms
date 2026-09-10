package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hydravms/internal/domain"
)

// PostgresEventRepository handles SQL persistence for AI and system events in PostgreSQL.
type PostgresEventRepository struct {
	pool *pgxpool.Pool
}

// NewEventRepository creates a new PostgresEventRepository.
func NewEventRepository(pool *pgxpool.Pool) *PostgresEventRepository {
	return &PostgresEventRepository{pool: pool}
}

func (r *PostgresEventRepository) Create(ctx context.Context, e *domain.Event) error {
	now := time.Now().UTC()
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.TriggeredAt.IsZero() {
		e.TriggeredAt = now
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	if e.Status == "" {
		e.Status = domain.EventStatusNew
	}

	bboxJSON, err := json.Marshal(e.BBoxNormalized)
	if err != nil {
		bboxJSON = []byte(`{"x_center":0,"y_center":0,"width":0,"height":0}`)
	}

	query := `
		INSERT INTO events (
			id, tenant_id, camera_id, rule_id, zone_id, event_type,
			severity, status, triggered_at, resolved_at, resolved_by_user_id,
			object_class, confidence, bbox_normalized, tracking_id,
			snapshot_s3_key, crop_s3_key, clip_s3_key, is_pinned, notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21
		)
	`
	_, err = r.pool.Exec(ctx, query,
		e.ID, e.TenantID, e.CameraID, e.RuleID, e.ZoneID, e.EventType,
		string(e.Severity), string(e.Status), e.TriggeredAt, e.ResolvedAt, e.ResolvedByUserID,
		e.ObjectClass, e.Confidence, bboxJSON, e.TrackingID,
		e.SnapshotS3Key, e.CropS3Key, e.ClipS3Key, e.IsPinned, e.Notes, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	return nil
}

func (r *PostgresEventRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Event, error) {
	query := `
		SELECT 
			e.id, e.tenant_id, e.camera_id, COALESCE(c.name, ''), e.rule_id, e.zone_id,
			e.event_type, e.severity, e.status, e.triggered_at, e.resolved_at, e.resolved_by_user_id,
			e.object_class, e.confidence, e.bbox_normalized, e.tracking_id,
			COALESCE(e.snapshot_s3_key, ''), COALESCE(e.crop_s3_key, ''), COALESCE(e.clip_s3_key, ''),
			e.is_pinned, COALESCE(e.notes, ''), e.created_at
		FROM events e
		LEFT JOIN cameras c ON c.tenant_id = e.tenant_id AND c.id = e.camera_id
		WHERE e.id = $1 AND (e.tenant_id = $2 OR e.tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	row := r.pool.QueryRow(ctx, query, id, tenantID)

	var e domain.Event
	var severityStr, statusStr string
	var bboxRaw []byte

	err := row.Scan(
		&e.ID, &e.TenantID, &e.CameraID, &e.CameraName, &e.RuleID, &e.ZoneID,
		&e.EventType, &severityStr, &statusStr, &e.TriggeredAt, &e.ResolvedAt, &e.ResolvedByUserID,
		&e.ObjectClass, &e.Confidence, &bboxRaw, &e.TrackingID,
		&e.SnapshotS3Key, &e.CropS3Key, &e.ClipS3Key,
		&e.IsPinned, &e.Notes, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	e.Severity = domain.EventSeverity(severityStr)
	e.Status = domain.EventStatus(statusStr)
	_ = json.Unmarshal(bboxRaw, &e.BBoxNormalized)
	return &e, nil
}

func (r *PostgresEventRepository) List(ctx context.Context, tenantID uuid.UUID, cameraID string, limit int) ([]*domain.Event, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	var query string
	var args []any

	if cameraID != "" {
		query = `
			SELECT 
				e.id, e.tenant_id, e.camera_id, COALESCE(c.name, ''), e.rule_id, e.zone_id,
				e.event_type, e.severity, e.status, e.triggered_at, e.resolved_at, e.resolved_by_user_id,
				e.object_class, e.confidence, e.bbox_normalized, e.tracking_id,
				COALESCE(e.snapshot_s3_key, ''), COALESCE(e.crop_s3_key, ''), COALESCE(e.clip_s3_key, ''),
				e.is_pinned, COALESCE(e.notes, ''), e.created_at
			FROM events e
			LEFT JOIN cameras c ON c.tenant_id = e.tenant_id AND c.id = e.camera_id
			WHERE (e.tenant_id = $1 OR e.tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			  AND e.camera_id = $2
			ORDER BY e.triggered_at DESC
			LIMIT $3
		`
		args = []any{tenantID, cameraID, limit}
	} else {
		query = `
			SELECT 
				e.id, e.tenant_id, e.camera_id, COALESCE(c.name, ''), e.rule_id, e.zone_id,
				e.event_type, e.severity, e.status, e.triggered_at, e.resolved_at, e.resolved_by_user_id,
				e.object_class, e.confidence, e.bbox_normalized, e.tracking_id,
				COALESCE(e.snapshot_s3_key, ''), COALESCE(e.crop_s3_key, ''), COALESCE(e.clip_s3_key, ''),
				e.is_pinned, COALESCE(e.notes, ''), e.created_at
			FROM events e
			LEFT JOIN cameras c ON c.tenant_id = e.tenant_id AND c.id = e.camera_id
			WHERE (e.tenant_id = $1 OR e.tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			ORDER BY e.triggered_at DESC
			LIMIT $2
		`
		args = []any{tenantID, limit}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var e domain.Event
		var severityStr, statusStr string
		var bboxRaw []byte

		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.CameraID, &e.CameraName, &e.RuleID, &e.ZoneID,
			&e.EventType, &severityStr, &statusStr, &e.TriggeredAt, &e.ResolvedAt, &e.ResolvedByUserID,
			&e.ObjectClass, &e.Confidence, &bboxRaw, &e.TrackingID,
			&e.SnapshotS3Key, &e.CropS3Key, &e.ClipS3Key,
			&e.IsPinned, &e.Notes, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		e.Severity = domain.EventSeverity(severityStr)
		e.Status = domain.EventStatus(statusStr)
		_ = json.Unmarshal(bboxRaw, &e.BBoxNormalized)
		events = append(events, &e)
	}
	return events, nil
}

func (r *PostgresEventRepository) UpdateStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EventStatus, resolvedBy *uuid.UUID) error {
	now := time.Now().UTC()
	query := `
		UPDATE events
		SET status = $1, resolved_at = $2, resolved_by_user_id = $3
		WHERE id = $4 AND (tenant_id = $5 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	_, err := r.pool.Exec(ctx, query, string(status), now, resolvedBy, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}
	return nil
}
