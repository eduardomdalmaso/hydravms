package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	if h.recordingService == nil {
		writeError(w, http.StatusServiceUnavailable, "Recording service not initialized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		profiles, err := h.recordingService.ListProfiles(r.Context(), tenantID, camID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"profiles": profiles,
			"total":    len(profiles),
		})

	case http.MethodPost:
		var profile domain.CameraRecordingProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid recording profile payload: "+err.Error())
			return
		}
		profile.TenantID = tenantID
		profile.CameraID = camID
		if err := h.recordingService.SaveProfile(r.Context(), &profile); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, profile)

	case http.MethodDelete:
		if subPath == "" {
			writeError(w, http.StatusBadRequest, "Profile ID is required")
			return
		}
		if err := h.recordingService.DeleteProfile(r.Context(), tenantID, camID, subPath); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CameraHandler) handleListRecordings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, camID string) {
	if h.recordingService == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"recordings": []interface{}{}, "total": 0})
		return
	}

	switch r.Method {
	case http.MethodGet:
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		now := time.Now()
		startTime := now.Add(-24 * time.Hour)
		endTime := now.Add(1 * time.Hour)

		if startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = t
			}
		}
		if endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = t
			}
		}

		segments, err := h.recordingService.ListRecordings(r.Context(), tenantID, camID, startTime, endTime)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"recordings": segments,
			"total":      len(segments),
		})

	case http.MethodPost:
		var seg domain.RecordingSegment
		if err := json.NewDecoder(r.Body).Decode(&seg); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid segment payload: "+err.Error())
			return
		}
		seg.TenantID = tenantID
		seg.CameraID = camID
		if err := h.recordingService.InsertSegment(r.Context(), &seg); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to insert segment: "+err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, seg)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CameraHandler) handleRecordingStream(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, camID string) {
	if h.recordingService == nil {
		writeError(w, http.StatusServiceUnavailable, "Recording service not initialized")
		return
	}

	tStr := r.URL.Query().Get("t")
	targetTime := time.Now()
	if tStr != "" {
		if ms, err := strconv.ParseInt(tStr, 10, 64); err == nil {
			targetTime = time.UnixMilli(ms)
		} else if t, err := time.Parse(time.RFC3339, tStr); err == nil {
			targetTime = t
		}
	}

	startTime := targetTime.Add(-1 * time.Minute)
	endTime := targetTime.Add(1 * time.Minute)
	segments, err := h.recordingService.ListRecordings(r.Context(), tenantID, camID, startTime, endTime)
	if err != nil || len(segments) == 0 {
		segments, _ = h.recordingService.ListRecordings(r.Context(), tenantID, camID, targetTime.Add(-24*time.Hour), targetTime.Add(24*time.Hour))
	}

	if len(segments) == 0 {
		writeError(w, http.StatusNotFound, "No recording segment available for this time")
		return
	}

	var chosen *domain.RecordingSegment = segments[0]
	minDiff := time.Duration(1<<63 - 1)
	for _, seg := range segments {
		if (targetTime.Equal(seg.StartTime) || targetTime.After(seg.StartTime)) && (targetTime.Equal(seg.EndTime) || targetTime.Before(seg.EndTime)) {
			chosen = seg
			break
		}
		diff := targetTime.Sub(seg.StartTime)
		if diff < 0 {
			diff = -diff
		}
		if diff < minDiff {
			minDiff = diff
			chosen = seg
		}
	}

	if chosen.S3Key == "" {
		writeError(w, http.StatusNotFound, "Recording file path not found")
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("X-Segment-Start", chosen.StartTime.Format(time.RFC3339))
	w.Header().Set("X-Segment-End", chosen.EndTime.Format(time.RFC3339))
	w.Header().Set("X-Segment-ID", chosen.ID.String())
	http.ServeFile(w, r, chosen.S3Key)
}

// syncCameraWithHydraStream informs HydraStream Data Plane about a newly registered camera.
func syncCameraWithHydraStream(cam *domain.Camera) {
	if cam == nil || cam.RTSPURL == "" {
		return
	}
	payload := map[string]interface{}{
		"tenant_id":        cam.TenantID.String(),
		"stream_id":        cam.ID,
		"source_url":       cam.RTSPURL,
		"decoding_engine":  "nvidia_nvdec",
		"status":           "online",
		"resolution":       cam.Resolution,
		"codec":            cam.Codec,
		"ingest_fps":       cam.FPS,
	}
	b, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(fmt.Sprintf("%s/api/v1/streams", getStreamBaseURL()), "application/json", bytes.NewReader(b))
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}

func getStreamBaseURL() string {
	if u := os.Getenv("HYDRASTREAM_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "http://localhost:8080"
}

