package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

type SQLiteCameraRepository struct {
	db *sql.DB
}

func NewCameraRepository(db *sql.DB) *SQLiteCameraRepository {
	return &SQLiteCameraRepository{db: db}
}

func (r *SQLiteCameraRepository) Create(ctx context.Context, c *domain.Camera) error {
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Codec == "" {
		c.Codec = "H.264"
	}

	var folderIDStr, nodeIDStr sql.NullString
	if c.FolderID != nil {
		folderIDStr = sql.NullString{String: c.FolderID.String(), Valid: true}
	}
	if c.AssignedNodeID != nil {
		nodeIDStr = sql.NullString{String: c.AssignedNodeID.String(), Valid: true}
	}

	query := `
		INSERT INTO cameras (
			id, tenant_id, name, protocol, rtsp_url, sub_stream_url,
			onvif_ip, onvif_port, onvif_user, onvif_pass, rtmp_stream_key,
			location, status, resolution, fps, bitrate_kbps, codec, is_active,
			folder_id, assigned_node_id, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.TenantID.String(), c.Name, string(c.Protocol), c.RTSPURL, c.SubStreamURL,
		c.ONVIFIP, c.ONVIFPort, c.ONVIFUser, c.ONVIFPass, c.RTMPStreamKey,
		c.Location, string(c.Status), c.Resolution, c.FPS, c.BitrateKbps, c.Codec, c.IsActive,
		folderIDStr, nodeIDStr, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert camera: %w", err)
	}
	return nil
}

func (r *SQLiteCameraRepository) GetByID(ctx context.Context, tenantID uuid.UUID, cameraID string) (*domain.Camera, error) {
	query := `
		SELECT 
			id, tenant_id, name, protocol, COALESCE(rtsp_url, ''), COALESCE(sub_stream_url, ''),
			COALESCE(onvif_ip, ''), COALESCE(onvif_port, 80), COALESCE(onvif_user, ''), COALESCE(rtmp_stream_key, ''),
			COALESCE(location, ''), status, COALESCE(resolution, '1920x1080'), COALESCE(fps, 30.0),
			COALESCE(bitrate_kbps, 2048), COALESCE(codec, 'H.264'), is_active, folder_id, assigned_node_id, created_at, updated_at
		FROM cameras
		WHERE id = ? AND tenant_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, cameraID, tenantID.String())

	var c domain.Camera
	var tenantIDStr, protoStr, statusStr string
	var folderIDStr, nodeIDStr sql.NullString
	err := row.Scan(
		&c.ID, &tenantIDStr, &c.Name, &protoStr, &c.RTSPURL, &c.SubStreamURL,
		&c.ONVIFIP, &c.ONVIFPort, &c.ONVIFUser, &c.RTMPStreamKey,
		&c.Location, &statusStr, &c.Resolution, &c.FPS,
		&c.BitrateKbps, &c.Codec, &c.IsActive, &folderIDStr, &nodeIDStr, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("camera not found")
		}
		return nil, fmt.Errorf("failed to query camera: %w", err)
	}
	c.TenantID, _ = uuid.Parse(tenantIDStr)
	c.Protocol = domain.CameraProtocol(protoStr)
	c.Status = domain.CameraStatus(statusStr)
	c.HasPTZ = (c.Protocol == domain.ProtocolONVIF || c.ID == "cam_entrance_01" || c.ID == "cam_parking_03")
	if folderIDStr.Valid {
		fID, _ := uuid.Parse(folderIDStr.String)
		c.FolderID = &fID
	}
	if nodeIDStr.Valid {
		nID, _ := uuid.Parse(nodeIDStr.String)
		c.AssignedNodeID = &nID
	}
	return &c, nil
}

func (r *SQLiteCameraRepository) List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.Camera, error) {
	var query string
	var args []any

	if folderID != nil {
		query = `
			SELECT 
				id, tenant_id, name, protocol, COALESCE(rtsp_url, ''), COALESCE(sub_stream_url, ''),
				COALESCE(onvif_ip, ''), COALESCE(onvif_port, 80), COALESCE(onvif_user, ''), COALESCE(rtmp_stream_key, ''),
				COALESCE(location, ''), status, COALESCE(resolution, '1920x1080'), COALESCE(fps, 30.0),
				COALESCE(bitrate_kbps, 2048), COALESCE(codec, 'H.264'), is_active, folder_id, assigned_node_id, created_at, updated_at
			FROM cameras
			WHERE tenant_id = ? AND folder_id = ?
			ORDER BY id ASC
		`
		args = []any{tenantID.String(), folderID.String()}
	} else {
		query = `
			SELECT 
				id, tenant_id, name, protocol, COALESCE(rtsp_url, ''), COALESCE(sub_stream_url, ''),
				COALESCE(onvif_ip, ''), COALESCE(onvif_port, 80), COALESCE(onvif_user, ''), COALESCE(rtmp_stream_key, ''),
				COALESCE(location, ''), status, COALESCE(resolution, '1920x1080'), COALESCE(fps, 30.0),
				COALESCE(bitrate_kbps, 2048), COALESCE(codec, 'H.264'), is_active, folder_id, assigned_node_id, created_at, updated_at
			FROM cameras
			WHERE tenant_id = ?
			ORDER BY id ASC
		`
		args = []any{tenantID.String()}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cameras: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.Camera, 0)
	for rows.Next() {
		var c domain.Camera
		var tenantIDStr, protoStr, statusStr string
		var folderIDStr, nodeIDStr sql.NullString
		if err := rows.Scan(
			&c.ID, &tenantIDStr, &c.Name, &protoStr, &c.RTSPURL, &c.SubStreamURL,
			&c.ONVIFIP, &c.ONVIFPort, &c.ONVIFUser, &c.RTMPStreamKey,
			&c.Location, &statusStr, &c.Resolution, &c.FPS,
			&c.BitrateKbps, &c.Codec, &c.IsActive, &folderIDStr, &nodeIDStr, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan camera row: %w", err)
		}
		c.TenantID, _ = uuid.Parse(tenantIDStr)
		c.Protocol = domain.CameraProtocol(protoStr)
		c.Status = domain.CameraStatus(statusStr)
		c.HasPTZ = (c.Protocol == domain.ProtocolONVIF || c.ID == "cam_entrance_01" || c.ID == "cam_parking_03")
		if folderIDStr.Valid {
			fID, _ := uuid.Parse(folderIDStr.String)
			c.FolderID = &fID
		}
		if nodeIDStr.Valid {
			nID, _ := uuid.Parse(nodeIDStr.String)
			c.AssignedNodeID = &nID
		}
		results = append(results, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *SQLiteCameraRepository) Update(ctx context.Context, c *domain.Camera) error {
	c.UpdatedAt = time.Now().UTC()
	var folderIDStr, nodeIDStr sql.NullString
	if c.FolderID != nil {
		folderIDStr = sql.NullString{String: c.FolderID.String(), Valid: true}
	}
	if c.AssignedNodeID != nil {
		nodeIDStr = sql.NullString{String: c.AssignedNodeID.String(), Valid: true}
	}

	query := `
		UPDATE cameras
		SET 
			name = ?, protocol = ?, rtsp_url = ?, sub_stream_url = ?,
			location = ?, resolution = ?, fps = ?, bitrate_kbps = ?,
			codec = ?, folder_id = ?, assigned_node_id = ?, status = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`
	res, err := r.db.ExecContext(ctx, query,
		c.Name, string(c.Protocol), c.RTSPURL, c.SubStreamURL,
		c.Location, c.Resolution, c.FPS, c.BitrateKbps,
		c.Codec, folderIDStr, nodeIDStr, string(c.Status), c.UpdatedAt,
		c.ID, c.TenantID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update camera: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("camera not found or permission denied")
	}
	return nil
}

func (r *SQLiteCameraRepository) Delete(ctx context.Context, tenantID uuid.UUID, cameraID string) error {
	query := `DELETE FROM cameras WHERE id = ? AND tenant_id = ?`
	res, err := r.db.ExecContext(ctx, query, cameraID, tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("camera not found or permission denied")
	}
	return nil
}
