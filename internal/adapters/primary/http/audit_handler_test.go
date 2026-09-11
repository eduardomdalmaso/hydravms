package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpAdapter "hydravms/internal/adapters/primary/http"
	"hydravms/internal/adapters/secondary/memory"
	"hydravms/internal/application"
)

func TestAuditHandlerLogs(t *testing.T) {
	repo := memory.NewInMemoryAuditLogRepository()
	auditService := application.NewAuditService(repo)
	auditHandler := httpAdapter.NewAuditHandler(auditService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/system/logs", auditHandler.HandleLogs)

	// 1. Ingest log via POST
	postPayload := `{"tenantId":"tenant_alpha","actor":"admin","action":"USER_CREATED","entityType":"user","entityId":"usr_01","details":"Novo operador cadastrado"}`
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/system/logs", strings.NewReader(postPayload))
	postReq.Header.Set("Content-Type", "application/json")
	wPost := httptest.NewRecorder()
	mux.ServeHTTP(wPost, postReq)

	if wPost.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", wPost.Code)
	}

	// 2. Query logs via GET
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/system/logs", nil)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, getReq)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wGet.Code)
	}

	var res struct {
		Logs []struct {
			Action string `json:"action"`
			Target string `json:"target"`
		} `json:"logs"`
		Total int `json:"total"`
	}

	if err := json.NewDecoder(wGet.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Total == 0 || len(res.Logs) == 0 {
		t.Fatalf("expected non-empty logs list")
	}

	found := false
	for _, l := range res.Logs {
		if l.Action == "USER_CREATED" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find USER_CREATED action in logs")
	}
}
