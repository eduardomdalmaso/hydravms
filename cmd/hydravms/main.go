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
	"hydravms/internal/adapters/secondary/memory"
	natsAdapter "hydravms/internal/adapters/secondary/nats"
	"hydravms/internal/application"
)

func main() {
	log.Println("[HydraVMS] Starting Control Plane & Event Orchestrator Backend...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize Repositories (Memory with initial seed)
	folderRepo := memory.NewInMemoryFolderRepository()
	cameraRepo := memory.NewInMemoryCameraRepository()

	// 2. Initialize Application Services
	folderService := application.NewFolderService(folderRepo)
	cameraService := application.NewCameraService(cameraRepo)

	// 3. Initialize WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	// 4. Connect to NATS Event Mesh & JetStream if available
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

	// 5. Initialize HTTP Handlers & Router
	folderHandler := httpAdapter.NewFolderHandler(folderService)
	cameraHandler := httpAdapter.NewCameraHandler(cameraService)
	wsHandler := ws.NewWebSocketHandler(wsHub)

	router := httpAdapter.NewRouter(folderHandler, cameraHandler, wsHandler)
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
