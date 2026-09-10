package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/secondary/memory"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type mockBroadcaster struct {
	events []*domain.CloudEvent
}

func (m *mockBroadcaster) PublishCloudEvent(ctx context.Context, event *domain.CloudEvent) error {
	m.events = append(m.events, event)
	return nil
}

type mockWSHub struct {
	broadcasts []string
}

func (m *mockWSHub) BroadcastToTenant(tenantID uuid.UUID, topic string, payload []byte) {
	m.broadcasts = append(m.broadcasts, topic)
}

func TestCameraWatchdog_OfflineDetection(t *testing.T) {
	cameraRepo := memory.NewInMemoryCameraRepository()
	eventRepo := memory.NewInMemoryEventRepository()
	broadcaster := &mockBroadcaster{}
	wsHub := &mockWSHub{}

	tenantID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	cam := &domain.Camera{
		ID:        "cam_test_01",
		TenantID:  tenantID,
		Name:      "Test Camera",
		Protocol:  domain.ProtocolRTSP,
		RTSPURL:   "rtsp://127.0.0.1:59999/dummy", // unreachable port
		Status:    domain.CameraStatusOnline,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	_ = cameraRepo.Create(context.Background(), cam)

	watchdog := application.NewCameraWatchdog(cameraRepo, eventRepo, broadcaster, wsHub, 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchdog.Start(ctx)

	// Wait for watchdog polling (requires 2 failed probes)
	time.Sleep(250 * time.Millisecond)

	updatedCam, err := cameraRepo.GetByID(context.Background(), tenantID, "cam_test_01")
	if err != nil {
		t.Fatalf("failed to get camera: %v", err)
	}

	if updatedCam.Status != domain.CameraStatusOffline {
		t.Errorf("expected camera status 'offline', got '%s'", updatedCam.Status)
	}

	if len(broadcaster.events) == 0 {
		t.Errorf("expected CloudEvent to be published, got 0")
	} else if broadcaster.events[0].Type != domain.TypeCameraOffline {
		t.Errorf("expected event type '%s', got '%s'", domain.TypeCameraOffline, broadcaster.events[0].Type)
	}

	if len(wsHub.broadcasts) == 0 {
		t.Errorf("expected WS broadcasts, got 0")
	}
}
