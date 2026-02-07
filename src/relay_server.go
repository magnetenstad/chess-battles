package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var relayUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

type relayHub struct {
	mu    sync.Mutex
	rooms map[string]map[*relayClient]struct{}
}

type relayClient struct {
	hub       *relayHub
	room      string
	conn      *websocket.Conn
	sendQueue chan []byte
	closeOnce sync.Once
	closed    bool
	mu        sync.Mutex
}

func RunRelayServer(addr string) error {
	hub := &relayHub{
		rooms: map[string]map[*relayClient]struct{}{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.handleWebSocket)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}

func (h *relayHub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	room := strings.TrimSpace(r.URL.Query().Get("room"))
	peer := strings.TrimSpace(r.URL.Query().Get("peer"))
	if room == "" {
		http.Error(w, "room query parameter is required", http.StatusBadRequest)
		return
	}
	if peer == "" {
		peer = "anonymous"
	}

	conn, err := relayUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &relayClient{
		hub:       h,
		room:      room,
		conn:      conn,
		sendQueue: make(chan []byte, 256),
	}

	h.addClient(client)
	log.Printf("relay client connected: room=%s peer=%s", room, peer)

	go client.writeLoop()
	client.readLoop()
}

func (h *relayHub) addClient(client *relayClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.room] == nil {
		h.rooms[client.room] = map[*relayClient]struct{}{}
	}
	h.rooms[client.room][client] = struct{}{}
}

func (h *relayHub) removeClient(client *relayClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomClients, ok := h.rooms[client.room]
	if !ok {
		return
	}

	delete(roomClients, client)
	if len(roomClients) == 0 {
		delete(h.rooms, client.room)
	}
}

func (h *relayHub) broadcast(room string, sender *relayClient, payload []byte) {
	h.mu.Lock()
	roomClients := h.rooms[room]
	recipients := make([]*relayClient, 0, len(roomClients))
	for client := range roomClients {
		if client == sender {
			continue
		}
		recipients = append(recipients, client)
	}
	h.mu.Unlock()

	for _, client := range recipients {
		if ok := client.enqueue(payload); !ok {
			client.close()
		}
	}
}

func (c *relayClient) readLoop() {
	defer c.close()

	c.conn.SetReadLimit(1 << 20)

	for {
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}

		var event GameEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			continue
		}

		canonicalPayload, err := json.Marshal(event)
		if err != nil {
			continue
		}
		c.hub.broadcast(c.room, c, canonicalPayload)
	}
}

func (c *relayClient) writeLoop() {
	defer c.close()

	for {
		payload, ok := <-c.sendQueue
		if !ok {
			return
		}

		_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return
		}
	}
}

func (c *relayClient) enqueue(payload []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return false
	}

	select {
	case c.sendQueue <- payload:
		return true
	default:
		return false
	}
}

func (c *relayClient) close() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.closed = true
		c.mu.Unlock()

		c.hub.removeClient(c)
		close(c.sendQueue)
		_ = c.conn.Close()
	})
}
