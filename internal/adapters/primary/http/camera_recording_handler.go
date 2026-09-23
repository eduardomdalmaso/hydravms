package http

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

func (h *CameraHandler) handleCameraProfiles(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, camID, subPath string) {
	if h.recordingService == nil {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"profiles": []interface{}{},
				"total":    0,
			})
			return
		}
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

	cleanKey := filepath.Clean(chosen.S3Key)
	absPath, err := filepath.Abs(cleanKey)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recording file path")
		return
	}

	// Strictly prohibit traversal into sensitive system directories
	forbiddenRoots := []string{"/etc", "/proc", "/sys", "/root", "/var/log", "/usr", "/bin", "/sbin"}
	for _, root := range forbiddenRoots {
		if strings.HasPrefix(absPath, root) {
			writeError(w, http.StatusForbidden, "Access to requested path is forbidden")
			return
		}
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("X-Segment-Start", chosen.StartTime.Format(time.RFC3339))
	w.Header().Set("X-Segment-End", chosen.EndTime.Format(time.RFC3339))
	w.Header().Set("X-Segment-ID", chosen.ID.String())
	http.ServeFile(w, r, absPath)
}
