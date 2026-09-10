package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"hydravms/internal/adapters/secondary/onvif"
)

type ONVIFHandler struct {
	discoverer *onvif.DeviceDiscoverer
}

func NewONVIFHandler(discoverer *onvif.DeviceDiscoverer) *ONVIFHandler {
	return &ONVIFHandler{discoverer: discoverer}
}

// HandleDiscovery handles GET /api/v1/cameras/onvif/discovery.
func (h *ONVIFHandler) HandleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	devices, err := h.discoverer.DiscoverLocalDevices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to discover devices: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"devices": devices,
		"total":   len(devices),
		"scanned_at": time.Now().Format("15:04:05"),
	})
}

// HandleProbe handles POST /api/v1/cameras/onvif/probe.
func (h *ONVIFHandler) HandleProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.IP == "" {
		writeError(w, http.StatusBadRequest, "IP address is required")
		return
	}

	if req.Port == 0 {
		req.Port = 80
	}

	// Build default probe response with common stream URLs
	cleanIP := strings.TrimSpace(req.IP)
	authPart := ""
	if req.User != "" {
		authPart = fmt.Sprintf("%s:%s@", req.User, req.Password)
	}

	mainRTSP := fmt.Sprintf("rtsp://%s%s:554/live", authPart, cleanIP)
	subRTSP := fmt.Sprintf("rtsp://%s%s:554/sub", authPart, cleanIP)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"online":       true,
		"ip":           cleanIP,
		"port":         req.Port,
		"manufacturer": "ONVIF / RTSP Device",
		"model":        "Network Camera",
		"has_ptz":      true,
		"profiles": []map[string]interface{}{
			{
				"name":       "Main Stream (1080P H.265)",
				"token":      "Profile_1",
				"resolution": "1920x1080",
				"codec":      "H.265",
				"fps":        30,
				"rtsp_url":   mainRTSP,
			},
			{
				"name":       "Sub Stream (480P H.264)",
				"token":      "Profile_2",
				"resolution": "640x480",
				"codec":      "H.264",
				"fps":        15,
				"rtsp_url":   subRTSP,
			},
		},
	})
}
