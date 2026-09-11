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

	// REST API Endpoints - Folders & Cameras
	mux.HandleFunc("/api/v1/folders", rt.folderHandler.HandleFolders)
	mux.HandleFunc("/api/v1/folders/", rt.folderHandler.HandleFolderByID)
	mux.HandleFunc("/api/v1/cameras", rt.cameraHandler.HandleCameras)
	mux.HandleFunc("/api/v1/cameras/", rt.cameraHandler.HandleCameraByID)

	// Cluster Nodes (HydraStream & HydraForge Instances)
	if rt.clusterNodeHandler != nil {
		mux.HandleFunc("/api/v1/cluster/nodes", rt.clusterNodeHandler.HandleNodes)
		mux.HandleFunc("/api/v1/cluster/nodes/", rt.clusterNodeHandler.HandleNodeByID)
		mux.HandleFunc("/api/v1/cluster/probe", rt.clusterNodeHandler.HandleProbe)
	}

	// Forensic Audit & System Logs Endpoints
	if rt.auditHandler != nil {
		mux.HandleFunc("/api/v1/system/logs", rt.auditHandler.HandleLogs)
		mux.HandleFunc("/api/v1/audit/logs", rt.auditHandler.HandleLogs)
	}

	// Admin Center: Users, Layouts, Maps, Tours
	if rt.adminDataHandler != nil {
		mux.HandleFunc("/api/v1/users", rt.adminDataHandler.HandleUsers)
		mux.HandleFunc("/api/v1/layouts", rt.adminDataHandler.HandleLayouts)
		mux.HandleFunc("/api/v1/maps", rt.adminDataHandler.HandleMaps)
		mux.HandleFunc("/api/v1/tours", rt.adminDataHandler.HandleTours)
	}

	// Storage & MinIO S3 Endpoints
	if rt.storagePoolHandler != nil {
		mux.HandleFunc("/api/v1/storage/pools", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				rt.storagePoolHandler.CreatePool(w, r)
			} else {
				rt.storagePoolHandler.ListPools(w, r)
			}
		})
		mux.HandleFunc("/api/v1/storage/pools/", rt.storagePoolHandler.DeletePool)
		mux.HandleFunc("/api/v1/storage/disks", rt.storagePoolHandler.DetectDisks)
		mux.HandleFunc("/api/v1/storage/telemetry", rt.storagePoolHandler.GetTelemetry)
		mux.HandleFunc("/api/v1/storage/spillover/drain", rt.storagePoolHandler.TriggerDrain)
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

	// Wrap in middleware chain: CORS -> Logger -> Audit -> Auth
	handler := middleware.CORSMiddleware(mux)
	handler = middleware.LoggerMiddleware(handler)
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

	return handler
}
