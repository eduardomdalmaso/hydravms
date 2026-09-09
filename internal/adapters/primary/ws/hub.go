package ws

import (
	"sync"

	"github.com/google/uuid"
)

type TenantMessage struct {
	TenantID uuid.UUID
	Payload  []byte
}

// Hub maintains the set of active clients and broadcasts messages scoped by tenant.
type Hub struct {
	mu         sync.RWMutex
	tenants    map[uuid.UUID]map[*Client]bool
	Broadcast  chan *TenantMessage
	Register   chan *Client
	Unregister chan *Client
}

// NewHub creates a new multi-tenant WebSocket Hub.
func NewHub() *Hub {
	return &Hub{
		tenants:    make(map[uuid.UUID]map[*Client]bool),
		Broadcast:  make(chan *TenantMessage, 1024),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
	}
}

// Run executes the hub event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, exists := h.tenants[client.TenantID]; !exists {
				h.tenants[client.TenantID] = make(map[*Client]bool)
			}
			h.tenants[client.TenantID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, exists := h.tenants[client.TenantID]; exists {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.tenants, client.TenantID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			clients, exists := h.tenants[msg.TenantID]
			if exists {
				for client := range clients {
					select {
					case client.Send <- msg.Payload:
					default:
						// Slow consumer protection: drop connection if buffer is full
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToTenant sends a message to all connected clients of a specific tenant.
func (h *Hub) BroadcastToTenant(tenantID uuid.UUID, payload []byte) {
	h.Broadcast <- &TenantMessage{
		TenantID: tenantID,
		Payload:  payload,
	}
}
