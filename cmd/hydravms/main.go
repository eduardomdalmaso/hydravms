package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	httpAdapter "hydravms/internal/adapters/primary/http"
	"hydravms/internal/adapters/primary/ws"
	"hydravms/internal/adapters/secondary/memory"
	natsAdapter "hydravms/internal/adapters/secondary/nats"
	onvifAdapter "hydravms/internal/adapters/secondary/onvif"
	postgresAdapter "hydravms/internal/adapters/secondary/postgres"
	s3Adapter "hydravms/internal/adapters/secondary/s3"
	"hydravms/internal/application"
	"hydravms/internal/ports"
)

func main() {
	log.Println("[HydraVMS] Starting Control Plane & Event Orchestrator Backend...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize Database Repositories (PostgreSQL with fallback to In-Memory)
	var folderRepo ports.FolderRepository
	var cameraRepo ports.CameraRepository
	var storagePoolRepo ports.StoragePoolRepository
	var pgPool *pgxpool.Pool

	dbURL := os.Getenv("DATABASE_URL")
	pgCfg := postgresAdapter.DefaultConfig()
	if dbURL != "" {
		pgCfg.URL = dbURL
	}

	pool, err := postgresAdapter.NewPool(ctx, pgCfg)
	if err != nil {
		log.Printf("⚠️ [HydraVMS] PostgreSQL not reachable: %v (falling back to In-Memory repository)\n", err)
		folderRepo = memory.NewInMemoryFolderRepository()
		cameraRepo = memory.NewInMemoryCameraRepository()
	} else {
		log.Println("✅ [HydraVMS] PostgreSQL relational database connected and connection pool initialized")
		pgPool = pool
		defer pool.Close()
		folderRepo = postgresAdapter.NewFolderRepository(pool)
		cameraRepo = postgresAdapter.NewCameraRepository(pool)
		storagePoolRepo = postgresAdapter.NewStoragePoolRepository(pool)
	}

	// 2. Initialize Application Core Services (Hexagonal Architecture)
	folderService := application.NewFolderService(folderRepo)
	cameraService := application.NewCameraService(cameraRepo)

	// 3. Connect to S3 / MinIO Object Storage for Video Recordings & Snapshots
	s3Cfg := s3Adapter.DefaultConfig()
	if endpoint := os.Getenv("MINIO_ENDPOINT"); endpoint != "" {
		s3Cfg.Endpoint = endpoint
	}
	var storagePoolService *application.StoragePoolService
	minioClient, err := s3Adapter.NewMinIOClient(s3Cfg)
	if err != nil {
		log.Printf("⚠️ [HydraVMS] MinIO S3 storage not reachable: %v (video archiving disabled)\n", err)
	} else {
		if err := minioClient.EnsureBuckets(ctx, "hydravms-recordings", "hydravms-snapshots", "hydravms-reports", "hydravms-maps"); err != nil {
			log.Printf("⚠️ [HydraVMS] MinIO bucket auto-creation warning: %v\n", err)
		} else {
			log.Println("✅ [HydraVMS] MinIO S3 Object Storage connected & buckets verified")
		}
		if storagePoolRepo != nil {
			storagePoolService = application.NewStoragePoolService(storagePoolRepo, minioClient)
		}
	}

	// 4. Initialize Real-Time WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	// 5. Connect to NATS Event Mesh & JetStream if available
	natsCfg := natsAdapter.DefaultConfig()
	natsClient, err := natsAdapter.NewNATSClient(ctx, natsCfg)
	if err != nil {
		log.Printf("⚠️ [HydraVMS] NATS Event Mesh not reachable: %v (continuing standalone mode)\n", err)
	} else {
		log.Println("✅ [HydraVMS] NATS JetStream Event Mesh connected and active on", natsCfg.URL)
		defer natsClient.Close()

		// Start NATS to WebSocket Bridge
		bridge := ws.NewNATSWebSocketBridge(natsClient.Conn(), wsHub)
		if err := bridge.Start(ctx); err != nil {
			log.Printf("⚠️ [HydraVMS] Failed to start NATS bridge: %v\n", err)
		} else {
			log.Println("✅ [HydraVMS] NATS to WebSocket Bridge active for all tenants (hydra.v1.*.>)")
		}
	}

	// 6. Initialize HTTP Handlers & Router
	folderHandler := httpAdapter.NewFolderHandler(folderService)
	cameraHandler := httpAdapter.NewCameraHandler(cameraService)
	onvifDiscoverer := onvifAdapter.NewDeviceDiscoverer()
	onvifHandler := httpAdapter.NewONVIFHandler(onvifDiscoverer)

	var storagePoolHandler *httpAdapter.StoragePoolHandler
	if storagePoolService != nil {
		storagePoolHandler = httpAdapter.NewStoragePoolHandler(storagePoolService)
	}
	var clusterNodeHandler *httpAdapter.ClusterNodeHandler
	if pgPool != nil {
		clusterNodeRepo := postgresAdapter.NewClusterNodeRepository(pgPool)
		clusterNodeHandler = httpAdapter.NewClusterNodeHandler(clusterNodeRepo)
	}
	wsHandler := ws.NewWebSocketHandler(wsHub)

	router := httpAdapter.NewRouter(folderHandler, cameraHandler, storagePoolHandler, clusterNodeHandler, onvifHandler, wsHandler)
	handler := router.BuildHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// 6. Production-Ready HTTP Server with Strict Timeouts (Uber Go Style Guide & Slowloris Guard)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	go func() {
		log.Printf("🚀 [HydraVMS] Control Plane REST API & WebSockets running on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HydraVMS] HTTP Server failed: %v", err)
		}
	}()

	// 7. Graceful Shutdown & Drain Connections
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 [HydraVMS] Shutting down Control Plane gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ [HydraVMS] Server forced to shutdown: %v\n", err)
	}

	log.Println("✅ [HydraVMS] Control Plane stopped cleanly.")
}
