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
	natsAdapter "hydravms/internal/adapters/secondary/nats"
	postgresAdapter "hydravms/internal/adapters/secondary/postgres"
	s3Adapter "hydravms/internal/adapters/secondary/s3"
	"hydravms/internal/adapters/secondary/memory"
	"hydravms/internal/application"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

func main() {
	log.Println("[HydraVMS] Starting Control Plane (PostgreSQL, RBAC, StorageGuard, JetStream)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize Database Repositories (PostgreSQL with fallback to In-Memory)
	var folderRepo ports.FolderRepository
	var cameraRepo ports.CameraRepository
	var storagePoolRepo ports.StoragePoolRepository
	var auditRepo ports.AuditLogRepository
	var userRepo ports.UserRepository
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
		auditRepo = memory.NewInMemoryAuditLogRepository()
		userRepo = memory.NewInMemoryUserRepository()
	} else {
		log.Println("✅ [HydraVMS] PostgreSQL relational database connected and connection pool initialized")
		pgPool = pool
		defer pool.Close()
		folderRepo = postgresAdapter.NewFolderRepository(pool)
		cameraRepo = postgresAdapter.NewCameraRepository(pool)
		storagePoolRepo = postgresAdapter.NewStoragePoolRepository(pool)
		auditRepo = postgresAdapter.NewAuditLogRepository(pool)
		userRepo = postgresAdapter.NewUserRepository(pool)
	}

	// Initialize Event Repository
	var eventRepo ports.EventRepository
	if pgPool != nil {
		eventRepo = postgresAdapter.NewEventRepository(pgPool)
	} else {
		eventRepo = memory.NewInMemoryEventRepository()
	}

	// 2. Initialize Application Core Services (Hexagonal Architecture)
	folderService := application.NewFolderService(folderRepo)
	cameraService := application.NewCameraService(cameraRepo)
	auditService := application.NewAuditService(auditRepo)

	// Record initial system boot audit record
	_ = auditService.RecordAction(
		ctx,
		"00000000-0000-0000-0000-000000000001",
		"system",
		"127.0.0.1",
		"SYSTEM_BOOT",
		"system",
		"core",
		"HydraVMS Control Plane inicializado com sucesso",
		domain.AuditCategorySystem,
		domain.AuditLevelInfo,
		map[string]interface{}{"version": "1.0.0", "status": "online"},
	)

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

	// 5. Initialize Recording Service & Repositories
	var recordingRepo ports.RecordingRepository
	var recordingService *application.RecordingService
	if pgPool != nil {
		recordingRepo = postgresAdapter.NewRecordingRepository(pgPool)
		recordingService = application.NewRecordingService(recordingRepo)
	}

	// 6. Connect to NATS Event Mesh & JetStream if available
	var eventPublisher *natsAdapter.EventPublisher
	natsCfg := natsAdapter.DefaultConfig()
	natsClient, err := natsAdapter.NewNATSClient(ctx, natsCfg)
	if err != nil {
		log.Printf("⚠️ [HydraVMS] NATS Event Mesh not reachable: %v (continuing standalone mode)\n", err)
	} else {
		log.Println("✅ [HydraVMS] NATS JetStream Event Mesh connected and active on", natsCfg.URL)
		defer natsClient.Close()

		eventPublisher = natsAdapter.NewEventPublisher(natsClient)

		// Start NATS to WebSocket Bridge
		bridge := ws.NewNATSWebSocketBridge(natsClient.Conn(), wsHub)
		if err := bridge.Start(ctx); err != nil {
			log.Printf("⚠️ [HydraVMS] Failed to start NATS bridge: %v\n", err)
		} else {
			log.Println("✅ [HydraVMS] NATS to WebSocket Bridge active for all tenants (hydra.v1.*.>)")
		}

		// Start Durable NATS Recording Consumer
		if recordingRepo != nil {
			recordingConsumer := natsAdapter.NewRecordingConsumer(natsClient, recordingRepo)
			if err := recordingConsumer.Start(ctx); err != nil {
				log.Printf("⚠️ [HydraVMS] Failed to start NATS recording consumer: %v\n", err)
			}
		}
	}

	// 7. Launch Automatic Camera Health Watchdog Service
	var broadcaster application.EventBroadcaster
	if eventPublisher != nil {
		broadcaster = eventPublisher
	}
	watchdog := application.NewCameraWatchdog(cameraRepo, eventRepo, broadcaster, wsHub, 2500*time.Millisecond)
	watchdog.Start(ctx)

	// 8. Initialize HTTP Handlers & Router
	authService := application.NewAuthService(userRepo, auditService)
	authHandler := httpAdapter.NewAuthHandler(authService)
	folderHandler := httpAdapter.NewFolderHandler(folderService)
	cameraHandler := httpAdapter.NewCameraHandler(cameraService, recordingService)
	auditHandler := httpAdapter.NewAuditHandler(auditService)
	adminDataHandler := httpAdapter.NewAdminDataHandler(pgPool)

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

	router := httpAdapter.NewRouter(
		authHandler,
		folderHandler,
		cameraHandler,
		storagePoolHandler,
		clusterNodeHandler,
		auditHandler,
		adminDataHandler,
		authService,
		auditService,
		wsHandler,
	)
	handler := router.BuildHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Production-Ready HTTP Server with Strict Timeouts
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

	// Graceful Shutdown & Drain Connections
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
