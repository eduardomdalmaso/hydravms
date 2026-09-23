package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type CameraHandler struct {
	service          *application.CameraService
	recordingService *application.RecordingService
}

func NewCameraHandler(service *application.CameraService, recordingService *application.RecordingService) *CameraHandler {
	return &CameraHandler{service: service, recordingService: recordingService}
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
		if len(cameras) > 0 {
			go syncAllCamerasWithHydraStream(cameras)
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

		// Asynchronously register/sync stream with HydraStream Data Plane
		go syncCameraWithHydraStream(created)

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

	rawID := strings.TrimPrefix(r.URL.Path, "/api/v1/cameras/")
	if rawID == "" {
		writeError(w, http.StatusBadRequest, "Camera ID is required")
		return
	}

	if strings.HasSuffix(rawID, "/live") || strings.HasSuffix(rawID, "/mjpeg") {
		camID := strings.TrimSuffix(strings.TrimSuffix(rawID, "/live"), "/mjpeg")
		target := fmt.Sprintf("%s/api/v1/streams/%s/mjpeg", getStreamBaseURL(), camID)
		if r.URL.RawQuery != "" { target += "?" + r.URL.RawQuery }
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		return
	}

	if strings.HasSuffix(rawID, "/snapshot") {
		camID := strings.TrimSuffix(rawID, "/snapshot")
		target := fmt.Sprintf("%s/api/v1/streams/%s/snapshot.jpg", getStreamBaseURL(), camID)
		if r.URL.RawQuery != "" { target += "?" + r.URL.RawQuery }
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		return
	}

	if strings.Contains(rawID, "/recording-profiles") {
		h.handleRecordingProfiles(w, r, tenantID, rawID)
		return
	}

	if strings.Contains(rawID, "/recordings") {
		camID := strings.Split(rawID, "/recordings")[0]
		if strings.Contains(rawID, "/recordings/stream") {
			h.handleRecordingStream(w, r, tenantID, camID)
			return
		}
		h.handleListRecordings(w, r, tenantID, camID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		cam, err := h.service.GetCamera(r.Context(), tenantID, rawID)
		if err != nil {
			writeError(w, http.StatusNotFound, "Camera not found")
			return
		}
		writeJSON(w, http.StatusOK, cam)

	case http.MethodDelete:
		if err := h.service.DeleteCamera(r.Context(), tenantID, rawID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CameraHandler) handleRecordingProfiles(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, path string) {
	parts := strings.Split(path, "/recording-profiles")
	camID := parts[0]
	subPath := ""
	if len(parts) > 1 {
		subPath = strings.TrimPrefix(parts[1], "/")
	}
	h.handleCameraProfiles(w, r, tenantID, camID, subPath)
}

