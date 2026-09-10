package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hydravms/internal/domain"
)

// PostgresCameraRepository handles SQL operations for cameras in PostgreSQL.
type PostgresCameraRepository struct {
	pool *pgxpool.Pool
}

// NewCameraRepository creates a new instance of PostgresCameraRepository.
func NewCameraRepository(pool *pgxpool.Pool) *PostgresCameraRepository {
	return &PostgresCameraRepository{pool: pool}
}

func (r *PostgresCameraRepository) Create(ctx context.Context, c *domain.Camera) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Codec == "" {
		c.Codec = "H.264"
	}

	query := `
		INSERT INTO cameras (
			id, tenant_id, name, protocol, rtsp_url, sub_stream_url,
			onvif_ip, onvif_port, onvif_user, onvif_pass, rtmp_stream_key,
			location, status, resolution, fps, bitrate_kbps, codec, is_active,
			folder_id, assigned_node_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22
		)
	`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.TenantID, c.Name, string(c.Protocol), c.RTSPURL, c.SubStreamURL,
		c.ONVIFIP, c.ONVIFPort, c.ONVIFUser, c.ONVIFPass, c.RTMPStreamKey,
		c.Location, string(c.Status), c.Resolution, c.FPS, c.BitrateKbps, c.Codec, c.IsActive,
		c.FolderID, c.AssignedNodeID, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert camera: %w", err)
	}
	return nil
}

func (r *PostgresCameraRepository) GetByID(ctx context.Context, tenantID uuid.UUID, cameraID string) (*domain.Camera, error) {
	query := `
		SELECT 
			id, tenant_id, name, protocol, COALESCE(rtsp_url, ''), COALESCE(sub_stream_url, ''),
			COALESCE(onvif_ip, ''), COALESCE(onvif_port, 80), COALESCE(onvif_user, ''), COALESCE(rtmp_stream_key, ''),
			COALESCE(location, ''), status, COALESCE(resolution, '1920x1080'), COALESCE(fps, 30.0),
			COALESCE(bitrate_kbps, 2048), COALESCE(codec, 'H.264'), is_active, folder_id, assigned_node_id, created_at, updated_at
		FROM cameras
		WHERE id = $1 AND (tenant_id = $2 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	row := r.pool.QueryRow(ctx, query, cameraID, tenantID)

	var c domain.Camera
	var protoStr, statusStr string
	err := row.Scan(
		&c.ID, &c.TenantID, &c.Name, &protoStr, &c.RTSPURL, &c.SubStreamURL,
		&c.ONVIFIP, &c.ONVIFPort, &c.ONVIFUser, &c.RTMPStreamKey,
		&c.Location, &statusStr, &c.Resolution, &c.FPS,
		&c.BitrateKbps, &c.Codec, &c.IsActive, &c.FolderID, &c.AssignedNodeID, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("camera not found")
		}
		return nil, fmt.Errorf("failed to query camera: %w", err)
	}
	c.Protocol = domain.CameraProtocol(protoStr)
	c.Status = domain.CameraStatus(statusStr)
	c.HasPTZ = (c.Protocol == domain.ProtocolONVIF || c.ID == "cam_entrance_01" || c.ID == "cam_parking_03")
	return &c, nil
}

func (r *PostgresCameraRepository) List(ctx context.Context, tenantID uuid.UUID, folderID *uuid.UUID) ([]*domain.Camera, error) {
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
			WHERE (tenant_id = $1 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			  AND folder_id = $2
			ORDER BY id ASC
		`
		args = []any{tenantID, folderID}
	} else {
		query = `
			SELECT 
				id, tenant_id, name, protocol, COALESCE(rtsp_url, ''), COALESCE(sub_stream_url, ''),
				COALESCE(onvif_ip, ''), COALESCE(onvif_port, 80), COALESCE(onvif_user, ''), COALESCE(rtmp_stream_key, ''),
				COALESCE(location, ''), status, COALESCE(resolution, '1920x1080'), COALESCE(fps, 30.0),
				COALESCE(bitrate_kbps, 2048), COALESCE(codec, 'H.264'), is_active, folder_id, assigned_node_id, created_at, updated_at
			FROM cameras
			WHERE (tenant_id = $1 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
			ORDER BY id ASC
		`
		args = []any{tenantID}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cameras: %w", err)
	}
	defer rows.Close()

	results := make([]*domain.Camera, 0)
	for rows.Next() {
		var c domain.Camera
		var protoStr, statusStr string
		if err := rows.Scan(
			&c.ID, &c.TenantID, &c.Name, &protoStr, &c.RTSPURL, &c.SubStreamURL,
			&c.ONVIFIP, &c.ONVIFPort, &c.ONVIFUser, &c.RTMPStreamKey,
			&c.Location, &statusStr, &c.Resolution, &c.FPS,
			&c.BitrateKbps, &c.Codec, &c.IsActive, &c.FolderID, &c.AssignedNodeID, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan camera row: %w", err)
		}
		c.Protocol = domain.CameraProtocol(protoStr)
		c.Status = domain.CameraStatus(statusStr)
		c.HasPTZ = (c.Protocol == domain.ProtocolONVIF || c.ID == "cam_entrance_01" || c.ID == "cam_parking_03")
		results = append(results, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *PostgresCameraRepository) Update(ctx context.Context, c *domain.Camera) error {
	c.UpdatedAt = time.Now()
	query := `
		UPDATE cameras
		SET 
			name = $1, protocol = $2, rtsp_url = $3, sub_stream_url = $4,
			location = $5, resolution = $6, fps = $7, bitrate_kbps = $8,
			codec = $9, folder_id = $10, assigned_node_id = $11, updated_at = $12
		WHERE id = $13 AND (tenant_id = $14 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	cmdTag, err := r.pool.Exec(ctx, query,
		c.Name, string(c.Protocol), c.RTSPURL, c.SubStreamURL,
		c.Location, c.Resolution, c.FPS, c.BitrateKbps,
		c.Codec, c.FolderID, c.AssignedNodeID, c.UpdatedAt,
		c.ID, c.TenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update camera: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("camera not found or permission denied")
	}
	return nil
}

func (r *PostgresCameraRepository) Delete(ctx context.Context, tenantID uuid.UUID, cameraID string) error {
	query := `
		DELETE FROM cameras
		WHERE id = $1 AND (tenant_id = $2 OR tenant_id = '00000000-0000-0000-0000-000000000001'::uuid)
	`
	cmdTag, err := r.pool.Exec(ctx, query, cameraID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("camera not found or permission denied")
	}
	return nil
}
