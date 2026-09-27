package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type LayoutHandler struct {
	repo ports.LayoutRepository
}

func NewLayoutHandler(repo ports.LayoutRepository) *LayoutHandler {
	return &LayoutHandler{repo: repo}
}

func (h *LayoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	tenantID, err := middleware.GetTenantID(r.Context())
	if err != nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/layouts")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path != "" {
			layoutID, err := uuid.Parse(path)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid layout ID")
				return
			}
			l, err := h.repo.GetByID(r.Context(), tenantID, layoutID)
			if err != nil {
				writeError(w, http.StatusNotFound, "Layout not found")
				return
			}
			writeJSON(w, http.StatusOK, l)
			return
		}

		var folderID *uuid.UUID
		if fStr := r.URL.Query().Get("folder_id"); fStr != "" {
			if parsed, err := uuid.Parse(fStr); err == nil {
				folderID = &parsed
			}
		}

		if h.repo == nil {
			writeJSON(w, http.StatusOK, map[string]interface{}{"layouts": []interface{}{}, "total": 0})
			return
		}

		list, err := h.repo.List(r.Context(), tenantID, folderID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"layouts": list,
			"total":   len(list),
		})

	case http.MethodPost:
		var l domain.MosaicLayout
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid payload: "+err.Error())
			return
		}
		l.TenantID = tenantID
		if l.ID == uuid.Nil {
			l.ID = uuid.New()
		}

		if h.repo != nil {
			_ = h.repo.Create(r.Context(), &l)
		}
		writeJSON(w, http.StatusCreated, l)

	case http.MethodPut:
		var l domain.MosaicLayout
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid payload: "+err.Error())
			return
		}
		l.TenantID = tenantID
		if path != "" {
			if parsed, err := uuid.Parse(path); err == nil {
				l.ID = parsed
			}
		}

		if h.repo != nil {
			_ = h.repo.Update(r.Context(), &l)
		}
		writeJSON(w, http.StatusOK, l)

	case http.MethodDelete:
		if path == "" {
			writeError(w, http.StatusBadRequest, "Layout ID is required")
			return
		}
		layoutID, err := uuid.Parse(path)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid layout ID")
			return
		}
		if h.repo != nil {
			_ = h.repo.Delete(r.Context(), tenantID, layoutID)
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
