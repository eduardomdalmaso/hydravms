package http

import (
	"encoding/json"
	"net/http"
	"strings"

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
