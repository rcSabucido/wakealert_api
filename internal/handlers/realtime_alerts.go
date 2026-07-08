package handlers

import (
	"os"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	conn 	*websocket.Conn
	addr 	string
	send 	chan []byte
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[*client]struct{}
	upgrader websocket.Upgrader
}

var hub = &Hub{
	clients: make(map[*client]struct{}),
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
	        return origin == os.Getenv("CORS_ORIGIN")
	        //return true
		},
	},
}

func (h *Hub) register(c *client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	log.Println("client connected:", c.addr)
}

func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
		log.Println("client disconnected:", c.addr)
	}
	h.mu.Unlock()
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	real_ip := r.Header.Get("X-Real-IP")
	if real_ip == "" {
		real_ip = conn.RemoteAddr().String()
	}

	c := &client{conn: conn, addr: real_ip, send: make(chan []byte, 256)}
	h.register(c)

	// writePump: serializes all writes to this connection
	go func() {
		defer conn.Close()
		for msg := range c.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}()

	// readPump: detects disconnects and triggers cleanup
	go func() {
		defer h.unregister(c)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}


func HandleWS(w http.ResponseWriter, r *http.Request) {
	hub.HandleWS(w, r)
}

// Broadcast sends JSON-encoded data to all connected clients.
func (h *Hub) Broadcast(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			// client buffer full — drop the message
			log.Println("dropping message for slow client:", c.addr)
		}
	}
	return nil
}

// Broadcast is a package-level shortcut to hub.Broadcast.
func Broadcast(v any) error {
	return hub.Broadcast(v)
}