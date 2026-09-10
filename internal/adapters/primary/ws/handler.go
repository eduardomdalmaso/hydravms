package ws

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"hydravms/internal/adapters/primary/http/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configured via CORSMiddleware
	},
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// ServeWS upgrades the HTTP connection to WebSocket and registers the client.
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.GetTenantID(r.Context())
	if err != nil || tenantID == uuid.Nil {
		tenantID, _ = uuid.Parse("00000000-0000-0000-0000-000000000001")
	}

	userID, _ := middleware.GetUserID(r.Context())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub:      h.hub,
		Conn:     conn,
		TenantID: tenantID,
		UserID:   userID,
		Send:     make(chan []byte, 256),
		topics:   map[string]bool{"*": true},
	}

	client.Hub.Register <- client

	// Start concurrent pumps
	go client.WritePump()
	go client.ReadPump()
}
