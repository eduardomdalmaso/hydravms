package ws

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// NATSWebSocketBridge listens to NATS subjects and routes events to the WebSocket Hub.
type NATSWebSocketBridge struct {
	nc  *nats.Conn
	hub *Hub
}

// NewNATSWebSocketBridge creates a bridge between NATS and the WebSocket Hub.
func NewNATSWebSocketBridge(nc *nats.Conn, hub *Hub) *NATSWebSocketBridge {
	return &NATSWebSocketBridge{nc: nc, hub: hub}
}

// Start listens to all tenant subjects and broadcasts them in real time.
func (b *NATSWebSocketBridge) Start(ctx context.Context) error {
	// Wildcard subscription for all tenant events: hydra.v1.{tenant_id}.>
	sub, err := b.nc.Subscribe("hydra.v1.*.>", func(msg *nats.Msg) {
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 3 {
			return
		}

		tenantUUID, err := uuid.Parse(parts[2])
		if err != nil {
			return
		}

		b.hub.BroadcastToTenant(tenantUUID, msg.Data)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS wildcard subject: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = sub.Unsubscribe()
		log.Println("[WS-NATS] Bridge unsubscribed gracefully")
	}()

	return nil
}
