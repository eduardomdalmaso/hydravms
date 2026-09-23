package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"hydravms/internal/domain"
)

func (w *CameraWatchdog) probeCamera(cam *domain.Camera) bool {
	var targets []string

	if cam.RTSPURL != "" {
		raw := cam.RTSPURL
		if atIdx := strings.LastIndex(raw, "@"); atIdx != -1 {
			raw = raw[atIdx+1:]
		} else {
			raw = strings.TrimPrefix(raw, "rtsp://")
		}
		if slashIdx := strings.Index(raw, "/"); slashIdx != -1 {
			raw = raw[:slashIdx]
		}
		if !strings.Contains(raw, ":") {
			raw += ":554"
		}
		targets = append(targets, raw)
	}

	if cam.ONVIFIP != "" {
		port := cam.ONVIFPort
		if port <= 0 {
			port = 80
		}
		targets = append(targets, fmt.Sprintf("%s:%d", cam.ONVIFIP, port))
	}

	for _, target := range targets {
		conn, err := net.DialTimeout("tcp", target, 1200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
	}

	return false
}

func (w *CameraWatchdog) onCameraOffline(ctx context.Context, cam *domain.Camera) {
	log.Printf("⚠️ [CameraWatchdog] Camera '%s' (%s) is OFFLINE! Generating automatic system event...", cam.ID, cam.Name)

	cam.Status = domain.CameraStatusOffline
	_ = w.cameraRepo.Update(ctx, cam)

	if w.eventRepo != nil {
		evt := &domain.Event{
			TenantID:       cam.TenantID,
			CameraID:       cam.ID,
			CameraName:     cam.Name,
			EventType:      domain.TypeCameraOffline,
			Severity:       domain.SeverityCritical,
			Status:         domain.EventStatusNew,
			TriggeredAt:    time.Now().UTC(),
			ObjectClass:    "camera",
			Confidence:     1.0,
			BBoxNormalized: domain.BoundingBox{},
			Notes:          "Perda de sinal RTSP / Câmera desconectada da rede local",
		}
		_ = w.eventRepo.Create(ctx, evt)
	}

	cloudEvt := domain.NewCloudEvent(
		cam.TenantID,
		domain.TypeCameraOffline,
		domain.CategorySystemEvent,
		domain.SeverityCritical,
		cam.ID,
		map[string]interface{}{
			"camera_id":   cam.ID,
			"camera_name": cam.Name,
			"status":      "offline",
			"reason":      "RTSP connection timeout / link dropped",
			"time":        time.Now().UTC().Format(time.RFC3339),
		},
	)

	payload, _ := json.Marshal(cloudEvt)

	if w.eventPublisher != nil {
		_ = w.eventPublisher.PublishCloudEvent(ctx, cloudEvt)
	}

	if w.wsHub != nil {
		w.wsHub.BroadcastToTenant(cam.TenantID, "events.system.camera.offline", payload)
		w.wsHub.BroadcastToTenant(cam.TenantID, fmt.Sprintf("cameras.%s.events", cam.ID), payload)
	}
}

func (w *CameraWatchdog) onCameraOnline(ctx context.Context, cam *domain.Camera) {
	log.Printf("✅ [CameraWatchdog] Camera '%s' (%s) is back ONLINE! Restoring operational telemetry...", cam.ID, cam.Name)

	cam.Status = domain.CameraStatusOnline
	_ = w.cameraRepo.Update(ctx, cam)

	if w.eventRepo != nil {
		evt := &domain.Event{
			TenantID:       cam.TenantID,
			CameraID:       cam.ID,
			CameraName:     cam.Name,
			EventType:      domain.TypeCameraOnline,
			Severity:       domain.SeverityLow,
			Status:         domain.EventStatusResolved,
			TriggeredAt:    time.Now().UTC(),
			ObjectClass:    "camera",
			Confidence:     1.0,
			BBoxNormalized: domain.BoundingBox{},
			Notes:          "Sinal RTSP restabelecido / Câmera online",
		}
		_ = w.eventRepo.Create(ctx, evt)
	}

	cloudEvt := domain.NewCloudEvent(
		cam.TenantID,
		domain.TypeCameraOnline,
		domain.CategorySystemEvent,
		domain.SeverityLow,
		cam.ID,
		map[string]interface{}{
			"camera_id":   cam.ID,
			"camera_name": cam.Name,
			"status":      "online",
			"reason":      "RTSP stream restored",
			"time":        time.Now().UTC().Format(time.RFC3339),
		},
	)

	payload, _ := json.Marshal(cloudEvt)

	if w.eventPublisher != nil {
		_ = w.eventPublisher.PublishCloudEvent(ctx, cloudEvt)
	}

	if w.wsHub != nil {
		w.wsHub.BroadcastToTenant(cam.TenantID, "events.system.camera.online", payload)
		w.wsHub.BroadcastToTenant(cam.TenantID, fmt.Sprintf("cameras.%s.events", cam.ID), payload)
	}
}
