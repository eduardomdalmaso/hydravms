package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"hydravms/internal/adapters/primary/http/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	if err != nil {
		http.Error(w, "Unauthorized: missing tenant identity", http.StatusUnauthorized)
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

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
	}

	client.Hub.Register <- client

	// Start concurrent pumps
	go client.WritePump()
	go client.ReadPump()
}
