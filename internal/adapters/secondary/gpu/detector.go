package gpu

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GPUTelemetry represents physical NVIDIA GPU hardware metrics.
type GPUTelemetry struct {
	Detected     bool    `json:"detected"`
	Model        string  `json:"model"`
	TotalVRAMMB  float64 `json:"total_vram_mb"`
	UsedVRAMMB   float64 `json:"used_vram_mb"`
	VRAMUsagePct float64 `json:"vram_usage_pct"`
	GPUUtilPct   float64 `json:"gpu_util_pct"`
	TempCelsius  float64 `json:"temp_celsius"`
	PowerWatts   float64 `json:"power_watts"`
	EngineName   string  `json:"engine_name"`
}

// QueryGPU queries physical NVIDIA GPU sensors via nvidia-smi.
func QueryGPU() GPUTelemetry {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=name,memory.total,memory.used,utilization.gpu,temperature.gpu,power.draw", "--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 {
			fields := strings.Split(lines[0], ",")
			if len(fields) >= 5 {
				model := strings.TrimSpace(fields[0])
				totalMB, _ := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
				usedMB, _ := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
				utilPct, _ := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
				tempC, _ := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
				powerW := 0.0
				if len(fields) >= 6 {
					powerW, _ = strconv.ParseFloat(strings.TrimSpace(fields[5]), 64)
				}

				vramPct := 0.0
				if totalMB > 0 {
					vramPct = (usedMB / totalMB) * 100.0
				}

				return GPUTelemetry{
					Detected:     true,
					Model:        model,
					TotalVRAMMB:  totalMB,
					UsedVRAMMB:   usedMB,
					VRAMUsagePct: vramPct,
					GPUUtilPct:   utilPct,
					TempCelsius:  tempC,
					PowerWatts:   powerW,
					EngineName:   "NVDEC CUDA IPC (Zero-Copy VRAM)",
				}
			}
		}
	}

	return GPUTelemetry{
		Detected:     false,
		Model:        "CPU Host (Hardware Acceleration Disabled)",
		TotalVRAMMB:  0,
		UsedVRAMMB:   0,
		VRAMUsagePct: 0,
		GPUUtilPct:   0,
		TempCelsius:  0,
		PowerWatts:   0,
		EngineName:   "FFmpeg POSIX SHM (/dev/shm)",
	}
}
