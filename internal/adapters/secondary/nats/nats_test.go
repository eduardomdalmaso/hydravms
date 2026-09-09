package nats_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/secondary/nats"
	"hydravms/internal/domain"
)

func TestNATSIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := nats.DefaultConfig()
	client, err := nats.NewNATSClient(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer client.Close()

	publisher := nats.NewEventPublisher(client)
	watcher := nats.NewCameraStateWatcher(client)

	tenantID := uuid.New()
	eventID := uuid.New()

	// 1. Test Publish AI Event
	event := &domain.Event{
		ID:          eventID,
		TenantID:    tenantID,
		CameraID:    "cam_test_01",
		EventType:   "INTRUSION_DETECTED",
		Severity:    domain.SeverityCritical,
		Status:      domain.EventStatusNew,
		TriggeredAt: time.Now(),
		ObjectClass: "person",
		Confidence:  0.95,
		BBoxNormalized: domain.BoundingBox{
			XCenter: 0.5,
			YCenter: 0.5,
			Width:   0.2,
			Height:  0.4,
		},
		CreatedAt: time.Now(),
	}

	if err := publisher.PublishAIEvent(ctx, event); err != nil {
		t.Fatalf("Failed to publish AI event: %v", err)
	}

	// 2. Test Set Camera State in KV Bucket
	state := map[string]interface{}{
		"status":     "online",
		"fps":        30.0,
		"bitrate":    2048,
		"codec":      "H.265",
		"updated_at": time.Now().Format(time.RFC3339),
	}
	if err := watcher.SetCameraState(ctx, tenantID.String(), "cam_test_01", state); err != nil {
		t.Fatalf("Failed to set camera state in KV: %v", err)
	}
}
