package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		json.NewEncoder(w).Encode(map[string]any{"nodes": nodes, "total": len(nodes)})
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

	targetURL := fmt.Sprintf("http://%s:%d/api/v1/info", req.IPAddress, req.HTTPPort)
	client := &http.Client{Timeout: 2 * time.Second}
	start := time.Now()

	resp, err := client.Get(targetURL)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		// Fallback test /healthz
		healthURL := fmt.Sprintf("http://%s:%d/healthz", req.IPAddress, req.HTTPPort)
		resp2, err2 := client.Get(healthURL)
		if err2 != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"online": false, "error": fmt.Sprintf("Conexão recusada em %s:%d", req.IPAddress, req.HTTPPort),
			})
			return
		}
		defer resp2.Body.Close()
		json.NewEncoder(w).Encode(map[string]any{
			"online": true, "latency_ms": latency, "app_name": "HydraStream Service",
			"suggested_role": "edge_ingest", "suggested_name": fmt.Sprintf("[NODE] HYDRASTREAM-%s", req.IPAddress),
		})
		return
	}
	defer resp.Body.Close()

	var info map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		info = map[string]any{"app_name": "HydraStream Engine"}
	}

	role := "edge_ingest"
	suggestedName := fmt.Sprintf("[NODE] HYDRASTREAM-%s", req.IPAddress)
	gpuModel, _ := info["gpu_model"].(string)

	if req.HTTPPort == 8081 {
		role = "gpu_worker"
		suggestedName = fmt.Sprintf("[NODE] HYDRAFORGE-GPU-%s", req.IPAddress)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"online": true, "latency_ms": latency, "app_name": info["app_name"], "version": info["version"],
		"engine_mode": info["engine_mode"], "gpu_model": gpuModel, "suggested_role": role,
		"suggested_name": suggestedName, "grpc_port": 50051, "webrtc_port": 8889,
	})
}
