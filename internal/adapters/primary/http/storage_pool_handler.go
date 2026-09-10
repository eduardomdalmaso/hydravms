package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type StoragePoolHandler struct {
	service *application.StoragePoolService
}

func NewStoragePoolHandler(service *application.StoragePoolService) *StoragePoolHandler {
	return &StoragePoolHandler{service: service}
}

func (h *StoragePoolHandler) ListPools(w http.ResponseWriter, r *http.Request) {
	var roleFilter *domain.StorageRole
	roleParam := r.URL.Query().Get("role")
	if roleParam != "" {
		r := domain.StorageRole(roleParam)
		roleFilter = &r
	}

	pools, err := h.service.ListPools(r.Context(), roleFilter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list storage pools: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pools": pools,
		"total": len(pools),
	})
}

func (h *StoragePoolHandler) CreatePool(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name           string                   `json:"name"`
		SourceType     domain.StorageSourceType `json:"source_type"`
		Role           domain.StorageRole       `json:"role"`
		NodeOrServer   string                   `json:"node_or_server"`
		PathOrEndpoint string                   `json:"path_or_endpoint"`
		Filesystem     string                   `json:"filesystem"`
		TotalGb        int                      `json:"total_gb"`
		RetentionDays  int                      `json:"retention_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	pool := &domain.StoragePool{
		Name:           req.Name,
		SourceType:     req.SourceType,
		Role:           req.Role,
		NodeOrServer:   req.NodeOrServer,
		PathOrEndpoint: req.PathOrEndpoint,
		Filesystem:     req.Filesystem,
		TotalBytes:     int64(req.TotalGb) * 1024 * 1024 * 1024,
		RetentionDays:  req.RetentionDays,
		Status:         domain.StorageStatusOnline,
	}

	if err := h.service.CreatePool(r.Context(), pool); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create storage pool: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, pool)
}

func (h *StoragePoolHandler) DeletePool(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/storage/pools/")
	id, err := uuid.Parse(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid storage pool ID: "+err.Error())
		return
	}

	if err := h.service.DeletePool(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete storage pool: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *StoragePoolHandler) GetTelemetry(w http.ResponseWriter, r *http.Request) {
	telemetry, err := h.service.GetTelemetry(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch storage telemetry: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, telemetry)
}

func (h *StoragePoolHandler) TriggerDrain(w http.ResponseWriter, r *http.Request) {
	msg, err := h.service.TriggerDrain(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to trigger spillover drain: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func (h *StoragePoolHandler) GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	key := r.URL.Query().Get("key")
	if bucket == "" || key == "" {
		writeError(w, http.StatusBadRequest, "bucket and key query parameters are required")
		return
	}

	url, err := h.service.GetPresignedPlaybackURL(r.Context(), bucket, key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate presigned URL: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

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

				// Se o dispositivo possui partições filhas, processar apenas os filhos
				if len(d.Children) > 0 {
					for _, child := range d.Children {
						processDevice(child, model, tran)
					}
					return
				}

				// Filtrar swap, boot, efi, loop e zram (< 10GB ou sistemas de boot)
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

	// Adicionar diretórios de storage locais se existirem
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
				var stat syscall.Statfs_t
				sizeGb := int64(1000)
				if err := syscall.Statfs(abs, &stat); err == nil {
					sizeGb = (int64(stat.Blocks) * int64(stat.Bsize)) / (1024 * 1024 * 1024)
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
