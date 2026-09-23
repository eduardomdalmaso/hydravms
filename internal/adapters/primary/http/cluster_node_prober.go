package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"hydravms/internal/adapters/secondary/gpu"
	"hydravms/internal/adapters/secondary/postgres"
)

func probeNodeMetrics(ctx context.Context, n *postgres.ClusterNode, sys gpu.SystemMetrics, liveGPU gpu.GPUTelemetry) {
	client := &http.Client{Timeout: 1500 * time.Millisecond}

	if n.NodeRole == "control_plane" && (n.IPAddress == "127.0.0.1" || n.IPAddress == "localhost") {
		n.Status = "online"
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
		return
	}

	if n.NodeRole == "edge_ingest" || n.HTTPPort == 8080 {
		probeEdgeNode(ctx, client, n, sys, liveGPU)
		return
	}

	if n.NodeRole == "gpu_worker" || n.HTTPPort == 8081 {
		probeGPUWorkerNode(ctx, client, n, sys, liveGPU)
		return
	}

	probeGenericNode(ctx, client, n)
}

func probeEdgeNode(ctx context.Context, client *http.Client, n *postgres.ClusterNode, sys gpu.SystemMetrics, liveGPU gpu.GPUTelemetry) {
	url := fmt.Sprintf("http://%s:%d/api/v1/telemetry/hardware", n.IPAddress, n.HTTPPort)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var data struct {
			Host struct {
				HostRAMUsagePct float64 `json:"host_ram_usage_pct"`
				HostUsedRAMMB   float64 `json:"host_used_ram_mb"`
				HostTotalRAMMB  float64 `json:"host_total_ram_mb"`
			} `json:"host"`
			GPU struct {
				Detected    bool    `json:"detected"`
				Model       string  `json:"model"`
				TotalVRAMMB float64 `json:"total_vram_mb"`
				UsedVRAMMB  float64 `json:"used_vram_mb"`
				GPUUtilPct  float64 `json:"gpu_util_pct"`
				TempCelsius float64 `json:"temp_celsius"`
				PowerWatts  float64 `json:"power_watts"`
			} `json:"gpu"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			n.Status = "online"
			n.CPUModel = sys.CPUModel
			n.CPUUsagePct = sys.CPUUsagePct
			if data.Host.HostTotalRAMMB > 0 {
				n.RAMUsedGB = data.Host.HostUsedRAMMB / 1024.0
				n.RAMTotalGB = data.Host.HostTotalRAMMB / 1024.0
				n.RAMUsagePct = data.Host.HostRAMUsagePct
			} else {
				n.RAMUsedGB = sys.RAMUsedGB
				n.RAMTotalGB = sys.RAMTotalGB
				n.RAMUsagePct = sys.RAMUsagePct
			}
			if data.GPU.Detected {
				n.GPUInfo = data.GPU.Model
				n.VRAMUsedMB = data.GPU.UsedVRAMMB
				n.VRAMTotalMB = data.GPU.TotalVRAMMB
				n.GPUUsagePct = data.GPU.GPUUtilPct
				n.TempCelsius = data.GPU.TempCelsius
				n.PowerWatts = data.GPU.PowerWatts
			}
			return
		}
	}
	healthURL := fmt.Sprintf("http://%s:%d/healthz", n.IPAddress, n.HTTPPort)
	reqH, _ := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	respH, errH := client.Do(reqH)
	if errH == nil && respH.StatusCode == http.StatusOK {
		defer respH.Body.Close()
		n.Status = "online"
		n.CPUModel = sys.CPUModel
		n.CPUUsagePct = sys.CPUUsagePct
		n.RAMUsedGB = sys.RAMUsedGB
		n.RAMTotalGB = sys.RAMTotalGB
		n.RAMUsagePct = sys.RAMUsagePct
		if liveGPU.Detected {
			n.GPUInfo = liveGPU.Model
			n.VRAMUsedMB = liveGPU.UsedVRAMMB
			n.VRAMTotalMB = liveGPU.TotalVRAMMB
			n.GPUUsagePct = liveGPU.GPUUtilPct
			n.TempCelsius = liveGPU.TempCelsius
			n.PowerWatts = liveGPU.PowerWatts
		}
		return
	}
	n.Status = "offline"
	n.CPUUsagePct, n.RAMUsagePct, n.GPUUsagePct, n.VRAMUsedMB = 0, 0, 0, 0
}

func probeGPUWorkerNode(ctx context.Context, client *http.Client, n *postgres.ClusterNode, sys gpu.SystemMetrics, liveGPU gpu.GPUTelemetry) {
	url := fmt.Sprintf("http://%s:%d/api/v1/training/telemetry", n.IPAddress, n.HTTPPort)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var data struct {
			GPUStats struct {
				Model       string  `json:"model"`
				TotalVRAMMB float64 `json:"total_vram_mb"`
				UsedVRAMMB  float64 `json:"used_vram_mb"`
				GPUUtilPct  float64 `json:"gpu_util_pct"`
				TempCelsius float64 `json:"temp_celsius"`
				PowerWatts  float64 `json:"power_watts"`
			} `json:"gpu_stats"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			n.Status = "online"
			n.CPUModel = sys.CPUModel
			n.CPUUsagePct = sys.CPUUsagePct
			n.RAMUsedGB = sys.RAMUsedGB
			n.RAMTotalGB = sys.RAMTotalGB
			n.RAMUsagePct = sys.RAMUsagePct
			if data.GPUStats.Model != "" {
				n.GPUInfo = data.GPUStats.Model
				n.VRAMUsedMB = data.GPUStats.UsedVRAMMB
				n.VRAMTotalMB = data.GPUStats.TotalVRAMMB
				n.GPUUsagePct = data.GPUStats.GPUUtilPct
				n.TempCelsius = data.GPUStats.TempCelsius
				n.PowerWatts = data.GPUStats.PowerWatts
			} else if liveGPU.Detected {
				n.GPUInfo = liveGPU.Model
				n.VRAMUsedMB = liveGPU.UsedVRAMMB
				n.VRAMTotalMB = liveGPU.TotalVRAMMB
				n.GPUUsagePct = liveGPU.GPUUtilPct
				n.TempCelsius = liveGPU.TempCelsius
				n.PowerWatts = liveGPU.PowerWatts
			}
			return
		}
	}
	healthURL := fmt.Sprintf("http://%s:%d/healthz", n.IPAddress, n.HTTPPort)
	reqH, _ := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	respH, errH := client.Do(reqH)
	if errH == nil && respH.StatusCode == http.StatusOK {
		defer respH.Body.Close()
		n.Status = "online"
		n.CPUModel = sys.CPUModel
		n.CPUUsagePct = sys.CPUUsagePct
		n.RAMUsedGB = sys.RAMUsedGB
		n.RAMTotalGB = sys.RAMTotalGB
		n.RAMUsagePct = sys.RAMUsagePct
		if liveGPU.Detected {
			n.GPUInfo = liveGPU.Model
			n.VRAMUsedMB = liveGPU.UsedVRAMMB
			n.VRAMTotalMB = liveGPU.TotalVRAMMB
			n.GPUUsagePct = liveGPU.GPUUtilPct
			n.TempCelsius = liveGPU.TempCelsius
			n.PowerWatts = liveGPU.PowerWatts
		}
		return
	}
	n.Status = "offline"
	n.CPUUsagePct, n.RAMUsagePct, n.GPUUsagePct, n.VRAMUsedMB = 0, 0, 0, 0
}

func probeGenericNode(ctx context.Context, client *http.Client, n *postgres.ClusterNode) {
	probeURL := fmt.Sprintf("http://%s:%d/api/v1/info", n.IPAddress, n.HTTPPort)
	reqP, _ := http.NewRequestWithContext(ctx, http.MethodGet, probeURL, nil)
	respP, errP := client.Do(reqP)
	if errP == nil && respP.StatusCode == http.StatusOK {
		defer respP.Body.Close()
		n.Status = "online"
		return
	}
	healthURL := fmt.Sprintf("http://%s:%d/healthz", n.IPAddress, n.HTTPPort)
	reqH, _ := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	respH, errH := client.Do(reqH)
	if errH == nil && respH.StatusCode == http.StatusOK {
		defer respH.Body.Close()
		n.Status = "online"
		return
	}
	n.Status = "offline"
	n.CPUUsagePct, n.RAMUsagePct, n.GPUUsagePct, n.VRAMUsedMB = 0, 0, 0, 0
}
