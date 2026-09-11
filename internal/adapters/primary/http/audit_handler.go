package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type AuditHandler struct {
	auditService *application.AuditService
}

func NewAuditHandler(service *application.AuditService) *AuditHandler {
	return &AuditHandler{auditService: service}
}

// HandleLogs serves GET and POST on /api/v1/system/logs and /api/v1/audit/logs.
func (h *AuditHandler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		tenantID := q.Get("tenant_id")
		category := q.Get("category")
		level := q.Get("level")

		limit := 100
		offset := 0
		if lStr := q.Get("limit"); lStr != "" {
			if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
				limit = l
			}
		}
		if oStr := q.Get("offset"); oStr != "" {
			if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
				offset = o
			}
		}

		logs, total, err := h.auditService.ListLogs(r.Context(), tenantID, category, level, limit, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		// Convert to frontend LogEntry structure
		type jsonLog struct {
			ID         string `json:"id"`
			Timestamp  string `json:"timestamp"`
			Category   string `json:"category"`
			Level      string `json:"level"`
			Actor      string `json:"actor"`
			TenantID   string `json:"tenantId"`
			TenantName string `json:"tenantName"`
			Action     string `json:"action"`
			Target     string `json:"target"`
			Details    string `json:"details"`
			IPAddress  string `json:"ipAddress"`
		}

		formatted := make([]jsonLog, 0, len(logs))
		for _, l := range logs {
			formatted = append(formatted, jsonLog{
				ID:         l.ID,
				Timestamp:  l.CreatedAt.Format(time.RFC3339),
				Category:   string(l.Category),
				Level:      string(l.Level),
				Actor:      l.Actor,
				TenantID:   l.TenantID,
				TenantName: l.TenantName,
				Action:     l.Action,
				Target:     l.Target,
				Details:    l.Details,
				IPAddress:  l.IPAddress,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"logs":   formatted,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})

	case http.MethodPost:
		var req struct {
			TenantID   string                 `json:"tenantId"`
			Actor      string                 `json:"actor"`
			IPAddress  string                 `json:"ipAddress"`
			Action     string                 `json:"action"`
			Category   string                 `json:"category"`
			Level      string                 `json:"level"`
			EntityType string                 `json:"entityType"`
			EntityID   string                 `json:"entityId"`
			Details    string                 `json:"details"`
			Payload    map[string]interface{} `json:"payload,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		cat := domain.AuditCategoryAudit
		if req.Category == "SYSTEM" {
			cat = domain.AuditCategorySystem
		}

		lvl := domain.AuditLevelInfo
		if req.Level != "" {
			lvl = domain.AuditLevel(req.Level)
		}

		ip := req.IPAddress
		if ip == "" {
			ip = r.RemoteAddr
		}

		err := h.auditService.RecordAction(
			r.Context(),
			req.TenantID,
			req.Actor,
			ip,
			req.Action,
			req.EntityType,
			req.EntityID,
			req.Details,
			cat,
			lvl,
			req.Payload,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
