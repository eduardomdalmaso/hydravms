package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/secondary/gpu"
)

func (h *ClusterNodeHandler) HandleProbe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenantID, err := middleware.GetTenantID(r.Context())
	if err != nil || tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IPAddress string `json:"ip_address"`
		HTTPPort  int    `json:"http_port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.IPAddress == "" { req.IPAddress = "127.0.0.1" }
	if req.HTTPPort == 0 { req.HTTPPort = 8080 }

	if req.HTTPPort <= 0 || req.HTTPPort > 65535 {
		http.Error(w, `{"error":"invalid port range"}`, http.StatusBadRequest)
		return
	}
	for _, blocked := range []string{"169.254.169.254", "0.0.0.0", "255.255.255.255"} {
		if strings.EqualFold(req.IPAddress, blocked) {
			http.Error(w, `{"error":"target address is restricted"}`, http.StatusForbidden)
			return
		}
	}

	existing, _ := h.repo.List(r.Context(), tenantID)
	streamCount, forgeCount, vmsCount := 0, 0, 0
	for _, en := range existing {
		if en.NodeRole == "edge_ingest" { streamCount++ }
		if en.NodeRole == "gpu_worker" { forgeCount++ }
		if en.NodeRole == "control_plane" { vmsCount++ }
	}

	targetURL := fmt.Sprintf("http://%s:%d/api/v1/info", req.IPAddress, req.HTTPPort)
	if req.HTTPPort == 8081 {
		targetURL = fmt.Sprintf("http://%s:%d/api/v1/training/telemetry", req.IPAddress, req.HTTPPort)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	start := time.Now()
	resp, err := client.Get(targetURL)
	latency := time.Since(start).Milliseconds()

	liveGPU := gpu.QueryGPU()
	defaultGPUModel := ""
	if liveGPU.Detected { defaultGPUModel = liveGPU.Model }

	if err != nil {
		healthURL := fmt.Sprintf("http://%s:%d/healthz", req.IPAddress, req.HTTPPort)
		resp2, err2 := client.Get(healthURL)
		if err2 != nil {
			json.NewEncoder(w).Encode(map[string]any{"online": false, "error": fmt.Sprintf("Conexão recusada em %s:%d", req.IPAddress, req.HTTPPort)})
			return
		}
		defer resp2.Body.Close()
		appName := "HydraStream Service"
		role := "edge_ingest"
		suggestedName := fmt.Sprintf("HYDRA-STREAM-NODE-%d", streamCount)
		if req.HTTPPort == 8081 {
			appName = "HydraForge Studio"
			role = "gpu_worker"
			suggestedName = fmt.Sprintf("HYDRA-FORGE-NODE-%d", forgeCount)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"online": true, "latency_ms": latency, "app_name": appName,
			"suggested_role": role, "suggested_name": suggestedName, "gpu_model": defaultGPUModel,
		})
		return
	}
	defer resp.Body.Close()

	var info map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		info = map[string]any{"app_name": "Hydra Engine"}
	}

	role := "edge_ingest"
	suggestedName := fmt.Sprintf("HYDRA-STREAM-NODE-%d", streamCount)
	gpuModel, _ := info["gpu_model"].(string)
	if gpuStats, ok := info["gpu_stats"].(map[string]any); ok {
		if m, ok := gpuStats["model"].(string); ok && m != "" {
			gpuModel = m
		}
	}
	if gpuModel == "" && liveGPU.Detected {
		gpuModel = liveGPU.Model
	}

	if req.HTTPPort == 8081 {
		role = "gpu_worker"
		suggestedName = fmt.Sprintf("HYDRA-FORGE-NODE-%d", forgeCount)
		if info["app_name"] == nil { info["app_name"] = "HydraForge AI Training & Inference" }
	} else if req.HTTPPort == 8083 {
		role = "control_plane"
		suggestedName = fmt.Sprintf("HYDRA-VMS-NODE-%d", vmsCount)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"online": true, "latency_ms": latency, "app_name": info["app_name"], "version": info["version"],
		"engine_mode": info["engine_mode"], "gpu_model": gpuModel, "suggested_role": role,
		"suggested_name": suggestedName, "grpc_port": 50051, "webrtc_port": 8889,
	})
}
