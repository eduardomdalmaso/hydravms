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

type CameraHandler struct {
	service *application.CameraService
}

func NewCameraHandler(service *application.CameraService) *CameraHandler {
	return &CameraHandler{service: service}
}

func (h *CameraHandler) HandleCameras(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	switch r.Method {
	case http.MethodGet:
		var folderID *uuid.UUID
		if fStr := r.URL.Query().Get("folder_id"); fStr != "" {
			if parsed, err := uuid.Parse(fStr); err == nil {
				folderID = &parsed
			}
		}

		cameras, err := h.service.ListCameras(r.Context(), tenantID, folderID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to list cameras: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"cameras": cameras,
			"total":   len(cameras),
		})

	case http.MethodPost:
		var cam domain.Camera
		if err := json.NewDecoder(r.Body).Decode(&cam); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid camera payload: "+err.Error())
			return
		}
		cam.TenantID = tenantID
		created, err := h.service.CreateCamera(r.Context(), &cam)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CameraHandler) HandleCameraByID(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	camID := strings.TrimPrefix(r.URL.Path, "/api/v1/cameras/")
	if camID == "" {
		writeError(w, http.StatusBadRequest, "Camera ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		cam, err := h.service.GetCamera(r.Context(), tenantID, camID)
		if err != nil {
			writeError(w, http.StatusNotFound, "Camera not found")
			return
		}
		writeJSON(w, http.StatusOK, cam)

	case http.MethodDelete:
		if err := h.service.DeleteCamera(r.Context(), tenantID, camID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
