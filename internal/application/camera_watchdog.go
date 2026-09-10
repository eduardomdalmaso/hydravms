package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type EventBroadcaster interface {
	PublishCloudEvent(ctx context.Context, event *domain.CloudEvent) error
}

type WSHubBroadcaster interface {
	BroadcastToTenant(tenantID uuid.UUID, topic string, payload []byte)
}

// CameraWatchdog continuously monitors RTSP camera connectivity and emits real-time events.
type CameraWatchdog struct {
	cameraRepo       ports.CameraRepository
	eventRepo        ports.EventRepository
	eventPublisher   EventBroadcaster
	wsHub            WSHubBroadcaster
	checkInterval    time.Duration
	mu               sync.RWMutex
	lastStates       map[string]domain.CameraStatus
	failureCount     map[string]int
}

// NewCameraWatchdog creates a new CameraWatchdog.
func NewCameraWatchdog(
	cameraRepo ports.CameraRepository,
	eventRepo ports.EventRepository,
	eventPublisher EventBroadcaster,
	wsHub WSHubBroadcaster,
	checkInterval time.Duration,
) *CameraWatchdog {
	if checkInterval <= 0 {
		checkInterval = 2500 * time.Millisecond
	}
	return &CameraWatchdog{
		cameraRepo:     cameraRepo,
		eventRepo:      eventRepo,
		eventPublisher: eventPublisher,
		wsHub:          wsHub,
		checkInterval:  checkInterval,
		lastStates:     make(map[string]domain.CameraStatus),
		failureCount:   make(map[string]int),
	}
}

// Start launches the watchdog loop in a background goroutine.
func (w *CameraWatchdog) Start(ctx context.Context) {
	go func() {
		log.Printf("🛡️ [CameraWatchdog] Active (polling every %v for camera offline/recovery detection)\n", w.checkInterval)
		ticker := time.NewTicker(w.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[CameraWatchdog] Stopped gracefully.")
				return
			case <-ticker.C:
				w.pollCameras(ctx)
			}
		}
	}()
}

func (w *CameraWatchdog) pollCameras(ctx context.Context) {
	// List cameras for default tenant
	defaultTenantID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	cameras, err := w.cameraRepo.List(ctx, defaultTenantID, nil)
	if err != nil || len(cameras) == 0 {
		return
	}

	for _, cam := range cameras {
		if !cam.IsActive {
			continue
		}
		w.checkCamera(ctx, cam)
	}
}

func (w *CameraWatchdog) checkCamera(ctx context.Context, cam *domain.Camera) {
	isReachable := w.probeCamera(cam)

	w.mu.Lock()
	lastStatus, hadPrevious := w.lastStates[cam.ID]
	if !hadPrevious {
		lastStatus = cam.Status
		if lastStatus == "" {
			lastStatus = domain.CameraStatusOnline
		}
		w.lastStates[cam.ID] = lastStatus
	}

	if isReachable {
		w.failureCount[cam.ID] = 0
		if lastStatus != domain.CameraStatusOnline {
			w.lastStates[cam.ID] = domain.CameraStatusOnline
			w.mu.Unlock()
			w.onCameraOnline(ctx, cam)
			return
		}
	} else {
		w.failureCount[cam.ID]++
		// Require 2 consecutive failed probes (avoid jitter) before declaring offline
		if w.failureCount[cam.ID] >= 2 && lastStatus != domain.CameraStatusOffline {
			w.lastStates[cam.ID] = domain.CameraStatusOffline
			w.mu.Unlock()
			w.onCameraOffline(ctx, cam)
			return
		}
	}
	w.mu.Unlock()
}

func (w *CameraWatchdog) probeCamera(cam *domain.Camera) bool {
	// 1. Probe RTSP Host:Port
	if cam.RTSPURL != "" {
		u, err := url.Parse(cam.RTSPURL)
		if err == nil && u.Host != "" {
			host := u.Host
			if !strings.Contains(host, ":") {
				host += ":554"
			}
			conn, err := net.DialTimeout("tcp", host, 1500*time.Millisecond)
			if err == nil {
				_ = conn.Close()
				return true
			}
		}
	}

	// 2. Probe ONVIF IP:Port if configured
	if cam.ONVIFIP != "" {
		port := cam.ONVIFPort
		if port <= 0 {
			port = 80
		}
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cam.ONVIFIP, port), 1500*time.Millisecond)
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

	// 1. Create and persist System Event in Database
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

	// 2. Build CNCF CloudEvents v1.0 Envelope
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

	// 3. Publish to NATS JetStream
	if w.eventPublisher != nil {
		_ = w.eventPublisher.PublishCloudEvent(ctx, cloudEvt)
	}

	// 4. Direct Broadcast to WebSocket Hub
	if w.wsHub != nil {
		w.wsHub.BroadcastToTenant(cam.TenantID, "events.system.camera.offline", payload)
		w.wsHub.BroadcastToTenant(cam.TenantID, fmt.Sprintf("cameras.%s.events", cam.ID), payload)
	}
}

func (w *CameraWatchdog) onCameraOnline(ctx context.Context, cam *domain.Camera) {
	log.Printf("✅ [CameraWatchdog] Camera '%s' (%s) is back ONLINE! Restoring operational telemetry...", cam.ID, cam.Name)

	cam.Status = domain.CameraStatusOnline
	_ = w.cameraRepo.Update(ctx, cam)

	// 1. Create and persist Recovery Event in Database
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

	// 2. Build CNCF CloudEvents v1.0 Envelope
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

	// 3. Publish to NATS JetStream
	if w.eventPublisher != nil {
		_ = w.eventPublisher.PublishCloudEvent(ctx, cloudEvt)
	}

	// 4. Direct Broadcast to WebSocket Hub
	if w.wsHub != nil {
		w.wsHub.BroadcastToTenant(cam.TenantID, "events.system.camera.online", payload)
		w.wsHub.BroadcastToTenant(cam.TenantID, fmt.Sprintf("cameras.%s.events", cam.ID), payload)
	}
}
