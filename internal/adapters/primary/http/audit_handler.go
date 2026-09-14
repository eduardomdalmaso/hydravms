package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
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
		tenantUUID, err := middleware.GetTenantID(r.Context())
		if err != nil || tenantUUID == uuid.Nil {
			tenantUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		}

		userRole := middleware.GetUserRole(r.Context())
		tenantID := tenantUUID.String()
		// Only superadmins on master tenant can filter by other tenants
		if userRole == "superadmin" && tenantUUID.String() == "00000000-0000-0000-0000-000000000001" {
			if qTenant := r.URL.Query().Get("tenant_id"); qTenant != "" {
				tenantID = qTenant
			}
		}

		q := r.URL.Query()
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
			writeError(w, http.StatusInternalServerError, err.Error())
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

		tenantUUID, err := middleware.GetTenantID(r.Context())
		if err != nil || tenantUUID == uuid.Nil {
			tenantUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		}

		if req.TenantID == "" || req.TenantID != tenantUUID.String() {
			req.TenantID = tenantUUID.String()
		}

		err = h.auditService.RecordAction(
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
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
