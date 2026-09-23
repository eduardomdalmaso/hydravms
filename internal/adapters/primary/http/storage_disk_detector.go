package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hydravms/internal/application"
)

type lsblkDev struct {
	Name       string     `json:"name"`
	Model      *string    `json:"model"`
	Size       int64      `json:"size"`
	Type       string     `json:"type"`
	Mountpoint *string    `json:"mountpoint"`
	Tran       *string    `json:"tran"`
	Fstype     *string    `json:"fstype"`
	Children   []lsblkDev `json:"children"`
}

type lsblkRes struct {
	Blockdevices []lsblkDev `json:"blockdevices"`
}

func (h *StoragePoolHandler) DetectDisks(w http.ResponseWriter, r *http.Request) {
	out, err := exec.CommandContext(r.Context(), "lsblk", "-J", "-b", "-o", "NAME,MODEL,SIZE,TYPE,MOUNTPOINT,TRAN,FSTYPE").Output()
	var detected []map[string]any

	if err == nil {
		var parsed lsblkRes
		if err := json.Unmarshal(out, &parsed); err == nil {
			var processDevice func(d lsblkDev, parentModel string, parentTran string)
			processDevice = func(d lsblkDev, parentModel string, parentTran string) {
				model := parentModel
				if d.Model != nil && *d.Model != "" {
					model = *d.Model
				}
				tran := parentTran
				if d.Tran != nil && *d.Tran != "" {
					tran = strings.ToUpper(*d.Tran)
				}
				if tran == "" {
					tran = "LOCAL"
				}

				fstype := "EXT4"
				if d.Fstype != nil && *d.Fstype != "" {
					fstype = strings.ToUpper(*d.Fstype)
				}

				if len(d.Children) > 0 {
					for _, child := range d.Children {
						processDevice(child, model, tran)
					}
					return
				}

				isBoot := d.Mountpoint != nil && (strings.HasPrefix(*d.Mountpoint, "/boot") || strings.HasPrefix(*d.Mountpoint, "/efi"))
				if fstype != "SWAP" && !isBoot && !strings.HasPrefix(d.Name, "loop") && !strings.HasPrefix(d.Name, "zram") && d.Size >= 10*1024*1024*1024 {
					role := "WARM_ARCHIVE"
					if tran == "NVME" || strings.Contains(strings.ToLower(model), "nvme") || strings.Contains(strings.ToLower(model), "ssd") {
						role = "HOT_BUFFER"
					}

					devPath := "/dev/" + d.Name
					if d.Mountpoint != nil && *d.Mountpoint != "" && !strings.HasPrefix(*d.Mountpoint, "[") {
						devPath = *d.Mountpoint
					}

					label := model
					if label == "" {
						label = fmt.Sprintf("Disco %s (%s)", tran, d.Name)
					}

					detected = append(detected, map[string]any{
						"device_path":      devPath,
						"model":            label,
						"bus_type":         tran,
						"size_gb":          d.Size / (1024 * 1024 * 1024),
						"filesystem":       fstype,
						"recommended_role": role,
					})
				}
			}

			for _, dev := range parsed.Blockdevices {
				processDevice(dev, "", "")
			}
		}
	}

	localDirs := []struct {
		path string
		role string
		desc string
	}{
		{"storage/hot", "HOT_BUFFER", "Pasta Local SSD (Buffer Quente)"},
		{"storage/archive", "WARM_ARCHIVE", "Pasta Local SSD (Arquivo Longo Prazo)"},
	}

	for _, ld := range localDirs {
		abs, err := filepath.Abs(ld.path)
		if err == nil {
			if _, statErr := os.Stat(abs); statErr == nil {
				sizeGb := int64(1000)
				if total, _, _, err := application.GetDiskUsage(abs); err == nil && total > 0 {
					sizeGb = total / (1024 * 1024 * 1024)
				}
				detected = append(detected, map[string]any{
					"device_path":      abs,
					"model":            ld.desc,
					"bus_type":         "NVMe",
					"size_gb":          sizeGb,
					"filesystem":       "LOCAL_FS",
					"recommended_role": ld.role,
				})
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"disks": detected,
		"total": len(detected),
	})
}
