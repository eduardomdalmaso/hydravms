package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type SQLiteEventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *SQLiteEventRepository {
	return &SQLiteEventRepository{db: db}
}

func (r *SQLiteEventRepository) Create(ctx context.Context, e *domain.Event) error {
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

	var ruleIDStr, zoneIDStr, resolvedByStr sql.NullString
	if e.RuleID != nil {
		ruleIDStr = sql.NullString{String: e.RuleID.String(), Valid: true}
	}
	if e.ZoneID != nil {
		zoneIDStr = sql.NullString{String: e.ZoneID.String(), Valid: true}
	}
	if e.ResolvedByUserID != nil {
		resolvedByStr = sql.NullString{String: e.ResolvedByUserID.String(), Valid: true}
	}
	var trackingIDVal sql.NullInt64
	if e.TrackingID != nil {
		trackingIDVal = sql.NullInt64{Int64: *e.TrackingID, Valid: true}
	}

	query := `
		INSERT INTO events (
			id, tenant_id, camera_id, rule_id, zone_id, event_type,
			severity, status, triggered_at, resolved_at, resolved_by_user_id,
			object_class, confidence, bbox_normalized, tracking_id,
			snapshot_s3_key, crop_s3_key, clip_s3_key, is_pinned, notes, created_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?, ?, ?
		)
	`
	_, err = r.db.ExecContext(ctx, query,
		e.ID.String(), e.TenantID.String(), e.CameraID, ruleIDStr, zoneIDStr, e.EventType,
		string(e.Severity), string(e.Status), e.TriggeredAt, e.ResolvedAt, resolvedByStr,
		e.ObjectClass, e.Confidence, string(bboxJSON), trackingIDVal,
		e.SnapshotS3Key, e.CropS3Key, e.ClipS3Key, e.IsPinned, e.Notes, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	return nil
}

func (r *SQLiteEventRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Event, error) {
	query := `
		SELECT 
			e.id, e.tenant_id, e.camera_id, COALESCE(c.name, ''), e.rule_id, e.zone_id,
			e.event_type, e.severity, e.status, e.triggered_at, e.resolved_at, e.resolved_by_user_id,
			e.object_class, e.confidence, e.bbox_normalized, e.tracking_id,
			COALESCE(e.snapshot_s3_key, ''), COALESCE(e.crop_s3_key, ''), COALESCE(e.clip_s3_key, ''),
			e.is_pinned, COALESCE(e.notes, ''), e.created_at
		FROM events e
		LEFT JOIN cameras c ON c.tenant_id = e.tenant_id AND c.id = e.camera_id
		WHERE e.id = ? AND (e.tenant_id = ? OR e.tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	row := r.db.QueryRowContext(ctx, query, id.String(), tenantID.String())

	var e domain.Event
	var idStr, tenantIDStr, severityStr, statusStr, bboxStr string
	var ruleIDStr, zoneIDStr, resolvedByStr sql.NullString
	var trackingIDNull sql.NullInt64

	err := row.Scan(
		&idStr, &tenantIDStr, &e.CameraID, &e.CameraName, &ruleIDStr, &zoneIDStr,
		&e.EventType, &severityStr, &statusStr, &e.TriggeredAt, &e.ResolvedAt, &resolvedByStr,
		&e.ObjectClass, &e.Confidence, &bboxStr, &trackingIDNull,
		&e.SnapshotS3Key, &e.CropS3Key, &e.ClipS3Key,
		&e.IsPinned, &e.Notes, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	e.ID, _ = uuid.Parse(idStr)
	e.TenantID, _ = uuid.Parse(tenantIDStr)
	e.Severity = domain.EventSeverity(severityStr)
	e.Status = domain.EventStatus(statusStr)
	if ruleIDStr.Valid {
		rID, _ := uuid.Parse(ruleIDStr.String)
		e.RuleID = &rID
	}
	if zoneIDStr.Valid {
		zID, _ := uuid.Parse(zoneIDStr.String)
		e.ZoneID = &zID
	}
	if resolvedByStr.Valid {
		uID, _ := uuid.Parse(resolvedByStr.String)
		e.ResolvedByUserID = &uID
	}
	if trackingIDNull.Valid {
		val := trackingIDNull.Int64
		e.TrackingID = &val
	}
	_ = json.Unmarshal([]byte(bboxStr), &e.BBoxNormalized)
	return &e, nil
}

func (r *SQLiteEventRepository) List(ctx context.Context, tenantID uuid.UUID, cameraID string, limit int) ([]*domain.Event, error) {
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
			WHERE (e.tenant_id = ? OR e.tenant_id = '00000000-0000-0000-0000-000000000001')
			  AND e.camera_id = ?
			ORDER BY e.triggered_at DESC
			LIMIT ?
		`
		args = []any{tenantID.String(), cameraID, limit}
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
			WHERE (e.tenant_id = ? OR e.tenant_id = '00000000-0000-0000-0000-000000000001')
			ORDER BY e.triggered_at DESC
			LIMIT ?
		`
		args = []any{tenantID.String(), limit}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var e domain.Event
		var idStr, tenantIDStr, severityStr, statusStr, bboxStr string
		var ruleIDStr, zoneIDStr, resolvedByStr sql.NullString
		var trackingIDNull sql.NullInt64

		if err := rows.Scan(
			&idStr, &tenantIDStr, &e.CameraID, &e.CameraName, &ruleIDStr, &zoneIDStr,
			&e.EventType, &severityStr, &statusStr, &e.TriggeredAt, &e.ResolvedAt, &resolvedByStr,
			&e.ObjectClass, &e.Confidence, &bboxStr, &trackingIDNull,
			&e.SnapshotS3Key, &e.CropS3Key, &e.ClipS3Key,
			&e.IsPinned, &e.Notes, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		e.ID, _ = uuid.Parse(idStr)
		e.TenantID, _ = uuid.Parse(tenantIDStr)
		e.Severity = domain.EventSeverity(severityStr)
		e.Status = domain.EventStatus(statusStr)
		if ruleIDStr.Valid {
			rID, _ := uuid.Parse(ruleIDStr.String)
			e.RuleID = &rID
		}
		if zoneIDStr.Valid {
			zID, _ := uuid.Parse(zoneIDStr.String)
			e.ZoneID = &zID
		}
		if resolvedByStr.Valid {
			uID, _ := uuid.Parse(resolvedByStr.String)
			e.ResolvedByUserID = &uID
		}
		if trackingIDNull.Valid {
			val := trackingIDNull.Int64
			e.TrackingID = &val
		}
		_ = json.Unmarshal([]byte(bboxStr), &e.BBoxNormalized)
		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return events, nil
}

func (r *SQLiteEventRepository) UpdateStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EventStatus, resolvedBy *uuid.UUID) error {
	now := time.Now().UTC()
	var resolvedByStr sql.NullString
	if resolvedBy != nil {
		resolvedByStr = sql.NullString{String: resolvedBy.String(), Valid: true}
	}
	query := `
		UPDATE events
		SET status = ?, resolved_at = ?, resolved_by_user_id = ?
		WHERE id = ? AND (tenant_id = ? OR tenant_id = '00000000-0000-0000-0000-000000000001')
	`
	_, err := r.db.ExecContext(ctx, query, string(status), now, resolvedByStr, id.String(), tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}
	return nil
}
