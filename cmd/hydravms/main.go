package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "hydravms/internal/adapters/primary/http"
	"hydravms/internal/adapters/primary/ws"
	natsAdapter "hydravms/internal/adapters/secondary/nats"
	s3Adapter "hydravms/internal/adapters/secondary/s3"
	"hydravms/internal/application"
	"hydravms/internal/domain"
)

func main() {
	log.Println("[HydraVMS] Starting Control Plane (Relational DB, RBAC, StorageGuard, JetStream)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repos := initRepositories(ctx)
	if repos.pgPool != nil {
		defer repos.pgPool.Close()
	}

	folderService := application.NewFolderService(repos.folderRepo)
	cameraService := application.NewCameraService(repos.cameraRepo)
	auditService := application.NewAuditService(repos.auditRepo)
	var pluginService *application.PluginService
	if repos.pluginRepo != nil {
		pluginService = application.NewPluginService(repos.pluginRepo)
	}

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
		if repos.storagePoolRepo != nil {
			storagePoolService = application.NewStoragePoolService(repos.storagePoolRepo, minioClient)
		}
	}

	wsHub := ws.NewHub()
	go wsHub.Run()

	var recordingService *application.RecordingService
	if repos.recordingRepo != nil {
		recordingService = application.NewRecordingService(repos.recordingRepo)
	}

	var eventPublisher *natsAdapter.EventPublisher
	natsCfg := natsAdapter.DefaultConfig()
	natsClient, err := natsAdapter.NewNATSClient(ctx, natsCfg)
	if err != nil {
		log.Printf("⚠️ [HydraVMS] NATS Event Mesh not reachable: %v (continuing standalone mode)\n", err)
	} else {
		log.Println("✅ [HydraVMS] NATS JetStream Event Mesh connected and active on", natsCfg.URL)
		defer natsClient.Close()

		eventPublisher = natsAdapter.NewEventPublisher(natsClient)

		bridge := ws.NewNATSWebSocketBridge(natsClient.Conn(), wsHub)
		if err := bridge.Start(ctx); err != nil {
			log.Printf("⚠️ [HydraVMS] Failed to start NATS bridge: %v\n", err)
		} else {
			log.Println("✅ [HydraVMS] NATS to WebSocket Bridge active for all tenants (hydra.v1.*.>)")
		}

		if repos.recordingRepo != nil {
			recordingConsumer := natsAdapter.NewRecordingConsumer(natsClient, repos.recordingRepo)
			if err := recordingConsumer.Start(ctx); err != nil {
				log.Printf("⚠️ [HydraVMS] Failed to start NATS recording consumer: %v\n", err)
			}
		}
	}

	var broadcaster application.EventBroadcaster
	if eventPublisher != nil {
		broadcaster = eventPublisher
	}
	watchdog := application.NewCameraWatchdog(repos.cameraRepo, repos.eventRepo, broadcaster, wsHub, 2500*time.Millisecond)
	watchdog.Start(ctx)

	authService := application.NewAuthService(repos.userRepo, auditService)
	authHandler := httpAdapter.NewAuthHandler(authService)
	folderHandler := httpAdapter.NewFolderHandler(folderService)
	cameraHandler := httpAdapter.NewCameraHandler(cameraService, recordingService)
	auditHandler := httpAdapter.NewAuditHandler(auditService)
	adminDataHandler := httpAdapter.NewAdminDataHandler(repos.pgPool)

	var pluginHandler *httpAdapter.PluginHandler
	if pluginService != nil {
		pluginHandler = httpAdapter.NewPluginHandler(pluginService)
	}
	var storagePoolHandler *httpAdapter.StoragePoolHandler
	if storagePoolService != nil {
		storagePoolHandler = httpAdapter.NewStoragePoolHandler(storagePoolService)
	}
	var clusterNodeHandler *httpAdapter.ClusterNodeHandler
	if repos.clusterNodeStore != nil {
		clusterNodeHandler = httpAdapter.NewClusterNodeHandler(repos.clusterNodeStore)
	}
	layoutHandler := httpAdapter.NewLayoutHandler(repos.layoutRepo)
	eventHandler := httpAdapter.NewEventHandler(repos.eventRepo, wsHub)
	wsHandler := ws.NewWebSocketHandler(wsHub)

	router := httpAdapter.NewRouter(
		authHandler,
		folderHandler,
		cameraHandler,
		eventHandler,
		pluginHandler,
		storagePoolHandler,
		clusterNodeHandler,
		auditHandler,
		adminDataHandler,
		layoutHandler,
		authService,
		auditService,
		wsHandler,
	)
	handler := router.BuildHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		log.Printf("🚀 [HydraVMS] Control Plane REST API & WebSockets running on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HydraVMS] HTTP Server failed: %v", err)
		}
	}()

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
