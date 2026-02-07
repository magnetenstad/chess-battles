package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type EventSocket struct {
	conn      *websocket.Conn
	peerID    string
	incoming  chan GameEvent
	closeOnce sync.Once
	writeMu   sync.Mutex
}

func NewEventSocket(rawURL, room, peerID string) (*EventSocket, error) {
	if strings.TrimSpace(room) == "" {
		return nil, fmt.Errorf("room codeword cannot be empty")
	}
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("websocket URL cannot be empty")
	}
	if peerID == "" {
		peerID = randomPeerID()
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid websocket URL: %w", err)
	}

	query := u.Query()
	query.Set("room", room)
	query.Set("peer", peerID)
	u.RawQuery = query.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial websocket: %w", err)
	}
	conn.SetReadLimit(1 << 20)

	socket := &EventSocket{
		conn:     conn,
		peerID:   peerID,
		incoming: make(chan GameEvent, 256),
	}

	go socket.readLoop()
	return socket, nil
}

func (s *EventSocket) Incoming() <-chan GameEvent {
	return s.incoming
}

func (s *EventSocket) Send(event GameEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	_ = s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := s.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return fmt.Errorf("write websocket message: %w", err)
	}

	return nil
}

func (s *EventSocket) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		closeErr = s.conn.Close()
	})
	return closeErr
}

func (s *EventSocket) PeerID() string {
	return s.peerID
}

func (s *EventSocket) readLoop() {
	defer close(s.incoming)

	for {
		_, payload, err := s.conn.ReadMessage()
		if err != nil {
			return
		}

		var event GameEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("ignoring malformed websocket payload: %v", err)
			continue
		}

		select {
		case s.incoming <- event:
		default:
			log.Printf("dropping remote event because incoming queue is full")
		}
	}
}

func randomPeerID() string {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("peer-%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(bytes)
}
