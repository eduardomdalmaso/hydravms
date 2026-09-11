package application_test

import (
	"context"
	"testing"

	"hydravms/internal/adapters/secondary/memory"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

func TestAuditServiceRecordAndList(t *testing.T) {
	repo := memory.NewInMemoryAuditLogRepository()
	service := application.NewAuditService(repo)
	ctx := context.Background()

	err := service.RecordAction(
		ctx,
		"tenant_alpha",
		"admin@hydravms.io",
		"192.168.1.100",
		"CAMERA_REGISTERED",
		"camera",
		"cam_lobby_01",
		"Câmera cadastrada com sucesso",
		domain.AuditCategoryAudit,
		domain.AuditLevelInfo,
		map[string]interface{}{"fps": 30, "codec": "h264"},
	)

	if err != nil {
		t.Fatalf("expected no error recording action, got: %v", err)
	}

	logs, total, err := service.ListLogs(ctx, "tenant_alpha", "ALL", "ALL", 10, 0)
	if err != nil {
		t.Fatalf("expected no error listing logs, got: %v", err)
	}

	if total == 0 || len(logs) == 0 {
		t.Fatalf("expected at least 1 log entry")
	}

	if logs[0].Action != "CAMERA_REGISTERED" {
		t.Errorf("expected CAMERA_REGISTERED, got %s", logs[0].Action)
	}
	if logs[0].Target != "camera // cam_lobby_01" {
		t.Errorf("expected camera // cam_lobby_01, got %s", logs[0].Target)
	}
}
