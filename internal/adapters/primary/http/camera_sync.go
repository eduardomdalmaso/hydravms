package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"hydravms/internal/domain"
)

// syncCameraWithHydraStream informs HydraStream Data Plane about a newly registered camera.
func syncCameraWithHydraStream(cam *domain.Camera) {
	if cam == nil || cam.RTSPURL == "" {
		return
	}
	payload := map[string]interface{}{
		"tenant_id":       cam.TenantID.String(),
		"stream_id":       cam.ID,
		"source_url":      cam.RTSPURL,
		"decoding_engine": "nvidia_nvdec",
		"status":          "online",
		"resolution":      cam.Resolution,
		"codec":           cam.Codec,
		"ingest_fps":      cam.FPS,
	}
	b, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(fmt.Sprintf("%s/api/v1/streams", getStreamBaseURL()), "application/json", bytes.NewReader(b))
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}

// syncAllCamerasWithHydraStream synchronizes all cameras in background.
func syncAllCamerasWithHydraStream(cameras []*domain.Camera) {
	for _, cam := range cameras {
		syncCameraWithHydraStream(cam)
	}
}

func getStreamBaseURL() string {
	if u := os.Getenv("HYDRASTREAM_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "http://localhost:8080"
}
