package http

import (
	"net/http"
	"strings"

	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/secondary/gpu"
	"hydravms/internal/application"

	"github.com/google/uuid"
)

type PluginHandler struct {
	service *application.PluginService
}

func NewPluginHandler(service *application.PluginService) *PluginHandler {
	return &PluginHandler{service: service}
}

func (h *PluginHandler) HandlePlugins(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	plugins, err := h.service.ListPlugins(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list plugins: "+err.Error())
		return
	}

	liveGPU := gpu.QueryGPU()
	for _, p := range plugins {
		if liveGPU.Detected {
			p.HardwareReq = "CUDA 13.3 // " + liveGPU.Model
		} else {
			p.HardwareReq = "CPU // ZERO-COPY SHM (AVX2)"
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"plugins":       plugins,
		"total":         len(plugins),
		"gpu_detected":  liveGPU.Detected,
		"gpu_telemetry": liveGPU,
	})
}

func (h *PluginHandler) HandlePluginAction(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/plugins/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeError(w, http.StatusBadRequest, "Invalid plugin ID")
		return
	}

	pluginID := parts[0]

	// GET /api/v1/plugins/{id}
	if len(parts) == 1 && r.Method == http.MethodGet {
		p, err := h.service.GetPluginByID(r.Context(), tenantID, pluginID)
		if err != nil {
			writeError(w, http.StatusNotFound, "Plugin not found")
			return
		}
		writeJSON(w, http.StatusOK, p)
		return
	}

	// Actions: POST /api/v1/plugins/{id}/install, /uninstall, /toggle
	if len(parts) == 2 && r.Method == http.MethodPost {
		action := parts[1]
		switch action {
		case "install":
			if err := h.service.InstallPlugin(r.Context(), tenantID, pluginID); err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to install plugin: "+err.Error())
				return
			}
			p, _ := h.service.GetPluginByID(r.Context(), tenantID, pluginID)
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"status":  "installed",
				"message": "Plugin instalado com sucesso",
				"plugin":  p,
			})
			return

		case "uninstall":
			if err := h.service.UninstallPlugin(r.Context(), tenantID, pluginID); err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to uninstall plugin: "+err.Error())
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{
				"status":  "uninstalled",
				"message": "Plugin desinstalado com sucesso",
			})
			return

		case "toggle":
			p, err := h.service.TogglePlugin(r.Context(), tenantID, pluginID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to toggle plugin: "+err.Error())
				return
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"status": "toggled",
				"plugin": p,
			})
			return

		default:
			writeError(w, http.StatusBadRequest, "Unknown action: "+action)
			return
		}
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}
