package http

import (
	"net/http"

	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/primary/ws"
)

type Router struct {
	folderHandler      *FolderHandler
	cameraHandler      *CameraHandler
	storagePoolHandler *StoragePoolHandler
	clusterNodeHandler *ClusterNodeHandler
	wsHandler          *ws.WebSocketHandler
}

func NewRouter(
	folderHandler *FolderHandler,
	cameraHandler *CameraHandler,
	storagePoolHandler *StoragePoolHandler,
	clusterNodeHandler *ClusterNodeHandler,
	wsHandler *ws.WebSocketHandler,
) *Router {
	return &Router{
		folderHandler:      folderHandler,
		cameraHandler:      cameraHandler,
		storagePoolHandler: storagePoolHandler,
		clusterNodeHandler: clusterNodeHandler,
		wsHandler:          wsHandler,
	}
}

// BuildHandler constructs the HTTP handler pipeline with security middlewares.
func (rt *Router) BuildHandler() http.Handler {
	mux := http.NewServeMux()

	// REST API Endpoints
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

	// Wrap in middleware chain: CORS -> Logger -> Auth
	handler := middleware.CORSMiddleware(mux)
	handler = middleware.LoggerMiddleware(handler)

	return handler
}
