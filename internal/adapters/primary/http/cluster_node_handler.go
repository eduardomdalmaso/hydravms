package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/secondary/gpu"
	"hydravms/internal/adapters/secondary/postgres"
)

type ClusterNodeStore interface {
	List(ctx context.Context, tenantID uuid.UUID) ([]*postgres.ClusterNode, error)
	Create(ctx context.Context, n *postgres.ClusterNode) error
	UpdateMetrics(ctx context.Context, id uuid.UUID, status string, cpuPct, ramPct, gpuPct, vramUsedMB float64, activeStreams int) error
	Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error
}

type ClusterNodeHandler struct {
	repo ClusterNodeStore
}

func NewClusterNodeHandler(repo ClusterNodeStore) *ClusterNodeHandler {
	return &ClusterNodeHandler{repo: repo}
}

func (h *ClusterNodeHandler) HandleNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenantID, err := middleware.GetTenantID(r.Context())
	if err != nil || tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if r.Method == http.MethodGet {
		nodes, err := h.repo.List(r.Context(), tenantID)
		if err != nil {
			http.Error(w, `{"error":"failed to fetch cluster nodes"}`, http.StatusInternalServerError)
			return
		}
		liveGPU := gpu.QueryGPU()
		sys := gpu.QuerySystemMetrics()

		var wg sync.WaitGroup
		for _, n := range nodes {
			wg.Add(1)
			go func(node *postgres.ClusterNode) {
				defer wg.Done()
				probeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				probeNodeMetrics(probeCtx, node, sys, liveGPU)
				_ = h.repo.UpdateMetrics(context.Background(), node.ID, node.Status, node.CPUUsagePct, node.RAMUsagePct, node.GPUUsagePct, node.VRAMUsedMB, node.ActiveStreams)
			}(n)
		}
		wg.Wait()

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
	tenantID, err := middleware.GetTenantID(r.Context())
	if err != nil || tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if r.Method == http.MethodDelete {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/cluster/nodes/")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid node id"}`, http.StatusBadRequest)
			return
		}
		if err := h.repo.Delete(r.Context(), id, tenantID); err != nil {
			http.Error(w, `{"error":"failed to delete node"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"deleted"}`))
		return
	}
	http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
}
