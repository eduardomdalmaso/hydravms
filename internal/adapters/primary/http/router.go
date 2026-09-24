package http

import (
	"net/http"

	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/primary/ws"
	"hydravms/internal/application"
)

type Router struct {
	authHandler        *AuthHandler
	folderHandler      *FolderHandler
	cameraHandler      *CameraHandler
	pluginHandler      *PluginHandler
	storagePoolHandler *StoragePoolHandler
	clusterNodeHandler *ClusterNodeHandler
	auditHandler       *AuditHandler
	adminDataHandler   *AdminDataHandler
	authService        *application.AuthService
	auditService       *application.AuditService
	wsHandler          *ws.WebSocketHandler
}

func NewRouter(
	authHandler *AuthHandler,
	folderHandler *FolderHandler,
	cameraHandler *CameraHandler,
	pluginHandler *PluginHandler,
	storagePoolHandler *StoragePoolHandler,
	clusterNodeHandler *ClusterNodeHandler,
	auditHandler *AuditHandler,
	adminDataHandler *AdminDataHandler,
	authService *application.AuthService,
	auditService *application.AuditService,
	wsHandler *ws.WebSocketHandler,
) *Router {
	return &Router{
		authHandler:        authHandler,
		folderHandler:      folderHandler,
		cameraHandler:      cameraHandler,
		pluginHandler:      pluginHandler,
		storagePoolHandler: storagePoolHandler,
		clusterNodeHandler: clusterNodeHandler,
		auditHandler:       auditHandler,
		adminDataHandler:   adminDataHandler,
		authService:        authService,
		auditService:       auditService,
		wsHandler:          wsHandler,
	}
}

// BuildHandler constructs the HTTP handler pipeline with security middlewares.
func (rt *Router) BuildHandler() http.Handler {
	mux := http.NewServeMux()

	// Authentication Endpoints (Public Login & Token Exchange)
	if rt.authHandler != nil {
		mux.HandleFunc("/api/v1/auth/login", rt.authHandler.HandleLogin)
	}

	adminOnly := middleware.RequireRole("admin", "superadmin")

	// REST API Endpoints - Folders & Cameras
	mux.HandleFunc("/api/v1/folders", rt.folderHandler.HandleFolders)
	mux.HandleFunc("/api/v1/folders/", rt.folderHandler.HandleFolderByID)
	mux.HandleFunc("/api/v1/cameras", rt.cameraHandler.HandleCameras)
	mux.HandleFunc("/api/v1/cameras/", rt.cameraHandler.HandleCameraByID)

	// Plugins & Marketplace Endpoints
	if rt.pluginHandler != nil {
		mux.HandleFunc("/api/v1/plugins", rt.pluginHandler.HandlePlugins)
		mux.HandleFunc("/api/v1/plugins/", rt.pluginHandler.HandlePluginAction)
	}

	// Cluster Nodes (HydraStream & HydraForge Instances) - Admin Only
	if rt.clusterNodeHandler != nil {
		mux.Handle("/api/v1/cluster/nodes", adminOnly(http.HandlerFunc(rt.clusterNodeHandler.HandleNodes)))
		mux.Handle("/api/v1/cluster/nodes/", adminOnly(http.HandlerFunc(rt.clusterNodeHandler.HandleNodeByID)))
		mux.Handle("/api/v1/cluster/probe", adminOnly(http.HandlerFunc(rt.clusterNodeHandler.HandleProbe)))
	}

	// Forensic Audit & System Logs Endpoints - Admin Only
	if rt.auditHandler != nil {
		mux.Handle("/api/v1/system/logs", adminOnly(http.HandlerFunc(rt.auditHandler.HandleLogs)))
		mux.Handle("/api/v1/audit/logs", adminOnly(http.HandlerFunc(rt.auditHandler.HandleLogs)))
	}

	// Admin Center: Users, Layouts, Maps, Tours
	if rt.adminDataHandler != nil {
		mux.Handle("/api/v1/users", adminOnly(http.HandlerFunc(rt.adminDataHandler.HandleUsers)))
		mux.HandleFunc("/api/v1/layouts", rt.adminDataHandler.HandleLayouts)
		mux.HandleFunc("/api/v1/maps", rt.adminDataHandler.HandleMaps)
		mux.HandleFunc("/api/v1/tours", rt.adminDataHandler.HandleTours)
	}

	// Storage & MinIO S3 Endpoints - Admin Only
	if rt.storagePoolHandler != nil {
		mux.Handle("/api/v1/storage/pools", adminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				rt.storagePoolHandler.CreatePool(w, r)
			} else {
				rt.storagePoolHandler.ListPools(w, r)
			}
		})))
		mux.Handle("/api/v1/storage/pools/", adminOnly(http.HandlerFunc(rt.storagePoolHandler.DeletePool)))
		mux.Handle("/api/v1/storage/disks", adminOnly(http.HandlerFunc(rt.storagePoolHandler.DetectDisks)))
		mux.Handle("/api/v1/storage/telemetry", adminOnly(http.HandlerFunc(rt.storagePoolHandler.GetTelemetry)))
		mux.Handle("/api/v1/storage/spillover/drain", adminOnly(http.HandlerFunc(rt.storagePoolHandler.TriggerDrain)))
		mux.HandleFunc("/api/v1/storage/presigned-url", rt.storagePoolHandler.GetPresignedURL)
	}

	// WebSocket Real-time Gateway (Multiplexed Pub/Sub)
	mux.HandleFunc("/ws/v1/live", rt.wsHandler.ServeWS)
	mux.HandleFunc("/ws/v1/events", rt.wsHandler.ServeWS)

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"hydravms-controlplane"}`))
	})

	// Wrap in middleware chain: (Mux -> Audit -> Auth -> Logger -> CORS)
	handler := http.Handler(mux)
	if rt.auditService != nil {
		handler = middleware.AuditMiddleware(rt.auditService)(handler)
	}
	if rt.authService != nil {
		validator := func(tokenStr string) (*middleware.TokenClaims, error) {
			claims, err := rt.authService.ValidateToken(tokenStr)
			if err != nil {
				return nil, err
			}
			return &middleware.TokenClaims{
				TenantID:         claims.TenantID,
				UserID:           claims.UserID,
				Role:             claims.Role,
				Scopes:           claims.Scopes,
				RegisteredClaims: claims.RegisteredClaims,
			}, nil
		}
		handler = middleware.AuthMiddleware(validator)(handler)
	}
	handler = middleware.LoggerMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)

	return handler
}
