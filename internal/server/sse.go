package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"
)

type SSEClient struct {
	UserID    uint
	Messages  chan []byte
	Connected bool
}

type SSEServer struct {
	clients map[uint]*SSEClient
	mu      sync.RWMutex
	db      *gorm.DB
}

func NewSSEServer(db *gorm.DB) *SSEServer {
	return &SSEServer{
		clients: make(map[uint]*SSEClient),
		db:      db,
	}
}

func (s *SSEServer) AddClient(userID uint) *SSEClient {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := &SSEClient{
		UserID:    userID,
		Messages:  make(chan []byte, 100),
		Connected: true,
	}

	s.clients[userID] = client
	return client
}

func (s *SSEServer) RemoveClient(userID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, ok := s.clients[userID]; ok {
		close(client.Messages)
		delete(s.clients, userID)
	}
}

func (s *SSEServer) SendToUser(userID uint, message []byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if client, ok := s.clients[userID]; ok {
		select {
		case client.Messages <- message:
			return true
		default:
			// Channel full, client might be slow
			return false
		}
	}
	return false
}

func (s *SSEServer) Broadcast(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		select {
		case client.Messages <- message:
		default:
			// Skip if channel is full
		}
	}
}

func (s *SSEServer) SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by AuthMiddleware)
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Add client to server
	client := s.AddClient(userID)
	defer s.RemoveClient(userID)

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\ndata: %s\n\n", "SSE connection established")
	flusher.Flush()

	// Keep connection alive and send messages
	for {
		select {
		case msg := <-client.Messages:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		case <-time.After(30 * time.Second):
			// Send heartbeat to keep connection alive
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

type NotificationMessage struct {
	Type    string      `json:"type"`    // "create", "update", "delete"
	Entity  string      `json:"entity"`  // "assignment", "course"
	ID      string      `json:"id"`      // Entity ID
	Message string      `json:"message"` // Human-readable message
	Data    interface{} `json:"data"`    // The actual data
}

func (s *SSEServer) SendNotification(userID uint, msgType, entity, id, message string, data interface{}) {
	notification := NotificationMessage{
		Type:    msgType,
		Entity:  entity,
		ID:      id,
		Message: message,
		Data:    data,
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return
	}

	s.SendToUser(userID, jsonData)
}
