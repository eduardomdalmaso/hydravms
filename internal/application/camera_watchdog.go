package application

import (
	"context"
	"log"
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
	cameraRepo     ports.CameraRepository
	eventRepo      ports.EventRepository
	eventPublisher EventBroadcaster
	wsHub          WSHubBroadcaster
	checkInterval  time.Duration
	mu             sync.RWMutex
	lastStates     map[string]domain.CameraStatus
	failureCount   map[string]int
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
		checkInterval = 2000 * time.Millisecond
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
		if w.failureCount[cam.ID] >= 2 && lastStatus != domain.CameraStatusOffline {
			w.lastStates[cam.ID] = domain.CameraStatusOffline
			w.mu.Unlock()
			w.onCameraOffline(ctx, cam)
			return
		}
	}
	w.mu.Unlock()
}
