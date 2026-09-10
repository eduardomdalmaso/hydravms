package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

func TestCloudEvent_Creation(t *testing.T) {
	tenantID := uuid.New()
	evt := domain.NewCloudEvent(
		tenantID,
		domain.TypeCameraOffline,
		domain.CategorySystemEvent,
		domain.SeverityCritical,
		"cam_01",
		map[string]string{"status": "offline"},
	)

	if evt.SpecVersion != "1.0" {
		t.Errorf("expected specversion '1.0', got '%s'", evt.SpecVersion)
	}
	if evt.Type != domain.TypeCameraOffline {
		t.Errorf("expected type '%s', got '%s'", domain.TypeCameraOffline, evt.Type)
	}
	if evt.TenantID != tenantID {
		t.Errorf("expected tenantID '%s', got '%s'", tenantID, evt.TenantID)
	}
	if evt.Category != domain.CategorySystemEvent {
		t.Errorf("expected category '%s', got '%s'", domain.CategorySystemEvent, evt.Category)
	}
	if evt.Severity != domain.SeverityCritical {
		t.Errorf("expected severity '%s', got '%s'", domain.SeverityCritical, evt.Severity)
	}
}
