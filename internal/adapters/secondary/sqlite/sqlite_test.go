package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

func TestSQLiteRepositories(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_hydravms.db")

	db, err := OpenDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to open SQLite DB: %v", err)
	}
	defer db.Close()
	defer os.Remove(dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// 1. Test User Repository
	userRepo := NewUserRepository(db)
	user, err := userRepo.FindByEmailOrUsername(ctx, "admin")
	if err != nil {
		t.Fatalf("Admin user not seeded properly: %v", err)
	}
	if user.Email != "admin@hydravms.io" {
		t.Errorf("Expected email admin@hydravms.io, got %s", user.Email)
	}

	// 2. Test Camera Repository
	camRepo := NewCameraRepository(db)
	cams, err := camRepo.List(ctx, tenantID, nil)
	if err != nil {
		t.Fatalf("Failed to list cameras: %v", err)
	}
	if len(cams) == 0 {
		t.Errorf("Expected seeded cameras, got 0")
	}

	newCam := &domain.Camera{
		ID:        "cam_test_01",
		TenantID:  tenantID,
		Name:      "TEST CAMERA",
		Protocol:  domain.ProtocolRTSP,
		RTSPURL:   "rtsp://127.0.0.1:8554/test",
		Status:    domain.CameraStatusOnline,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := camRepo.Create(ctx, newCam); err != nil {
		t.Fatalf("Failed to create test camera: %v", err)
	}
	fetchedCam, err := camRepo.GetByID(ctx, tenantID, "cam_test_01")
	if err != nil || fetchedCam.Name != "TEST CAMERA" {
		t.Fatalf("Failed to fetch created camera: %v", err)
	}

	// 3. Test Folder Repository
	folderRepo := NewFolderRepository(db)
	newFolder := &domain.Folder{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Module:    domain.ModuleCameras,
		Name:      "Setor Norte",
		ColorHex:  "#3b82f6",
		Icon:      "folder",
		SortOrder: 1,
	}
	if err := folderRepo.Create(ctx, newFolder); err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}
	tree, err := folderRepo.ListTree(ctx, tenantID, domain.ModuleCameras)
	if err != nil || len(tree) == 0 {
		t.Fatalf("Failed to list folder tree: %v", err)
	}

	// 4. Test Event Repository
	eventRepo := NewEventRepository(db)
	eventID := uuid.New()
	newEvent := &domain.Event{
		ID:          eventID,
		TenantID:    tenantID,
		CameraID:    "cam_test_01",
		EventType:   "intrusion",
		Severity:    domain.SeverityWarning,
		Status:      domain.EventStatusNew,
		TriggeredAt: time.Now().UTC(),
		ObjectClass: "person",
		Confidence:  0.92,
	}
	if err := eventRepo.Create(ctx, newEvent); err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	events, err := eventRepo.List(ctx, tenantID, "cam_test_01", 10)
	if err != nil || len(events) == 0 {
		t.Fatalf("Failed to list events: %v", err)
	}

	// 5. Test Audit Log Repository
	auditRepo := NewAuditLogRepository(db)
	newLog := &domain.AuditLog{
		TenantID:   tenantID.String(),
		Action:     "USER_LOGIN",
		EntityType: "user",
		EntityID:   user.ID.String(),
		IPAddress:  "127.0.0.1",
		Details:    "User logged in successfully",
	}
	if err := auditRepo.Save(ctx, newLog); err != nil {
		t.Fatalf("Failed to save audit log: %v", err)
	}
	logs, total, err := auditRepo.List(ctx, tenantID.String(), "", "", 10, 0)
	if err != nil || total == 0 || len(logs) == 0 {
		t.Fatalf("Failed to list audit logs: %v", err)
	}
}
