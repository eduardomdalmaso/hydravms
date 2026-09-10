package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512 KB
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ClientCommand struct {
	Action string   `json:"action"` // "subscribe", "unsubscribe", "ping"
	Topics []string `json:"topics"`
}

// Client represents a single active WebSocket subscriber.
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	TenantID uuid.UUID
	UserID   uuid.UUID
	Send     chan []byte
	mu       sync.RWMutex
	topics   map[string]bool
}

func (c *Client) Matches(topic string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.topics) == 0 || c.topics["*"] {
		return true
	}
	if c.topics[topic] {
		return true
	}
	for pattern := range c.topics {
		if strings.HasSuffix(pattern, ".*") {
			prefix := strings.TrimSuffix(pattern, ".*")
			if strings.HasPrefix(topic, prefix) {
				return true
			}
		} else if strings.HasSuffix(pattern, ".>") {
			prefix := strings.TrimSuffix(pattern, ".>")
			if strings.HasPrefix(topic, prefix) {
				return true
			}
		}
	}
	return false
}

func (c *Client) handleMessage(msg []byte) {
	var cmd ClientCommand
	if err := json.Unmarshal(msg, &cmd); err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	switch strings.ToLower(cmd.Action) {
	case "subscribe":
		for _, t := range cmd.Topics {
			c.topics[strings.TrimSpace(t)] = true
		}
	case "unsubscribe":
		for _, t := range cmd.Topics {
			delete(c.topics, strings.TrimSpace(t))
		}
	case "ping":
		select {
		case c.Send <- []byte(`{"type":"pong"}`):
		default:
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Read error for client %s: %v", c.UserID, err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleMessage(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
