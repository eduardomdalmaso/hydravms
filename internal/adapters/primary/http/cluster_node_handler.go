package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/secondary/gpu"
	"hydravms/internal/adapters/secondary/postgres"
)

type ClusterNodeHandler struct {
	repo *postgres.PostgresClusterNodeRepository
}

func NewClusterNodeHandler(repo *postgres.PostgresClusterNodeRepository) *ClusterNodeHandler {
	return &ClusterNodeHandler{repo: repo}
}

func (h *ClusterNodeHandler) HandleNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	if r.Method == http.MethodGet {
		nodes, err := h.repo.List(r.Context(), tenantID)
		if err != nil {
			http.Error(w, `{"error":"failed to fetch cluster nodes"}`, http.StatusInternalServerError)
			return
		}
		liveGPU := gpu.QueryGPU()
		sys := gpu.QuerySystemMetrics()

		for _, n := range nodes {
			if n.IPAddress == "127.0.0.1" || n.IPAddress == "localhost" {
				n.CPUModel = sys.CPUModel
				n.CPUUsagePct = sys.CPUUsagePct
				n.RAMUsedGB = sys.RAMUsedGB
				n.RAMTotalGB = sys.RAMTotalGB
				n.RAMUsagePct = sys.RAMUsagePct

				if liveGPU.Detected {
					if n.GPUInfo == "" || strings.Contains(n.GPUInfo, "RTX") {
						n.GPUInfo = liveGPU.Model
					}
					n.VRAMUsedMB = liveGPU.UsedVRAMMB
					n.VRAMTotalMB = liveGPU.TotalVRAMMB
					n.GPUUsagePct = liveGPU.GPUUtilPct
					n.TempCelsius = liveGPU.TempCelsius
					n.PowerWatts = liveGPU.PowerWatts
				}
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"nodes": nodes, "total": len(nodes), "host_system": sys})
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			NodeName   string `json:"node_name"`
			NodeRole   string `json:"node_role"`
			IPAddress  string `json:"ip_address"`
			GRPCPort   int    `json:"grpc_port"`
			WebRTCPort int    `json:"webrtc_port"`
			HTTPPort   int    `json:"http_port"`
			GPUInfo    string `json:"gpu_device_info"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if req.GRPCPort == 0 { req.GRPCPort = 50051 }
		if req.WebRTCPort == 0 { req.WebRTCPort = 8889 }
		if req.HTTPPort == 0 { req.HTTPPort = 8080 }
		if req.NodeRole == "" { req.NodeRole = "edge_ingest" }

		if req.NodeName == "" {
			existing, _ := h.repo.List(r.Context(), tenantID)
			count := 0
			for _, en := range existing {
				if en.NodeRole == req.NodeRole { count++ }
			}
			switch req.NodeRole {
			case "gpu_worker":
				req.NodeName = fmt.Sprintf("HYDRA-FORGE-NODE-%d", count)
			case "control_plane":
				req.NodeName = fmt.Sprintf("HYDRA-VMS-NODE-%d", count)
			default:
				req.NodeName = fmt.Sprintf("HYDRA-STREAM-NODE-%d", count)
			}
		}

		if req.GPUInfo == "" && (req.IPAddress == "127.0.0.1" || req.IPAddress == "localhost") {
			liveGPU := gpu.QueryGPU()
			if liveGPU.Detected {
				req.GPUInfo = liveGPU.Model
			}
		}

		node := &postgres.ClusterNode{
			TenantID: tenantID, NodeName: req.NodeName, NodeRole: req.NodeRole,
			IPAddress: req.IPAddress, GRPCPort: req.GRPCPort, WebRTCPort: req.WebRTCPort,
			HTTPPort: req.HTTPPort, GPUInfo: req.GPUInfo, Status: "online",
		}
		if err := h.repo.Create(r.Context(), node); err != nil {
			http.Error(w, `{"error":"failed to create node"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(node)
		return
	}

	http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
}

func (h *ClusterNodeHandler) HandleNodeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodDelete {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/cluster/nodes/")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid node id"}`, http.StatusBadRequest)
			return
		}
		if err := h.repo.Delete(r.Context(), id); err != nil {
			http.Error(w, `{"error":"failed to delete node"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"deleted"}`))
		return
	}
	http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
}

func (h *ClusterNodeHandler) HandleProbe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
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
	if liveGPU.Detected {
		defaultGPUModel = liveGPU.Model
	}

	if err != nil {
		healthURL := fmt.Sprintf("http://%s:%d/healthz", req.IPAddress, req.HTTPPort)
		resp2, err2 := client.Get(healthURL)
		if err2 != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"online": false, "error": fmt.Sprintf("Conexão recusada em %s:%d", req.IPAddress, req.HTTPPort),
			})
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

	switch req.HTTPPort {
	case 8081:
		role = "gpu_worker"
		suggestedName = fmt.Sprintf("HYDRA-FORGE-NODE-%d", forgeCount)
		if info["app_name"] == nil {
			info["app_name"] = "HydraForge AI Training & Inference"
		}
	case 8083:
		role = "control_plane"
		suggestedName = fmt.Sprintf("HYDRA-VMS-NODE-%d", vmsCount)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"online": true, "latency_ms": latency, "app_name": info["app_name"], "version": info["version"],
		"engine_mode": info["engine_mode"], "gpu_model": gpuModel, "suggested_role": role,
		"suggested_name": suggestedName, "grpc_port": 50051, "webrtc_port": 8889,
	})
}
