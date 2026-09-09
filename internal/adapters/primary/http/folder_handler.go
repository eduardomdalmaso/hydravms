package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type FolderHandler struct {
	service *application.FolderService
}

func NewFolderHandler(service *application.FolderService) *FolderHandler {
	return &FolderHandler{service: service}
}

type CreateFolderRequest struct {
	Module   domain.FolderModule `json:"module"`
	Name     string              `json:"name"`
	ParentID *uuid.UUID          `json:"parent_id,omitempty"`
	ColorHex string              `json:"color_hex,omitempty"`
}

func (h *FolderHandler) HandleFolders(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	switch r.Method {
	case http.MethodGet:
		moduleStr := r.URL.Query().Get("module")
		if moduleStr == "" {
			moduleStr = "cameras"
		}
		folders, err := h.service.ListFolderTree(r.Context(), tenantID, domain.FolderModule(moduleStr))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to list folders: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"module":  moduleStr,
			"folders": folders,
			"total":   len(folders),
		})

	case http.MethodPost:
		var req CreateFolderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
			return
		}
		folder, err := h.service.CreateFolder(r.Context(), tenantID, req.Module, req.Name, req.ParentID, req.ColorHex)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, folder)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *FolderHandler) HandleFolderByID(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/folders/")
	folderID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid folder UUID")
		return
	}

	if r.Method == http.MethodDelete {
		if err := h.service.DeleteFolder(r.Context(), tenantID, folderID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"type":   "https://hydravms.domain.com/errors/api-error",
		"title":  "API Error",
		"status": status,
		"detail": detail,
	})
}
