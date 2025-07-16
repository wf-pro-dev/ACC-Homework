package events

import (
	"encoding/json"
	"log"

	"github.com/williamfotso/acc/internal/services/client"
)

type EventHandler struct {
	onAssignmentCreate func(data json.RawMessage, message string)
	onAssignmentUpdate func(data json.RawMessage, message string)
	onAssignmentDelete func(data json.RawMessage, message string)
	stopChan           chan struct{}
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		stopChan: make(chan struct{}),
	}
}

// Start now accepts the sseClient as a parameter.
func (h *EventHandler) Start(sseClient *client.SSEClient) {
	if sseClient == nil {
		log.Fatal("[EventHandler] Fatal: SSE client is nil.")
		return
	}

	go func() {
		log.Println("[EventHandler] Starting to listen for events...")
		for {
			select {
			case <-h.stopChan:
				log.Println("[EventHandler] Stop signal received, shutting down.")
				return
			// Listen for events from the client's public channel.
			case event, ok := <-sseClient.Events():
				if !ok {
					log.Println("[EventHandler] SSE events channel closed.")
					return
				}
				h.handleEvent(event)
			// Listen for errors.
			case err, ok := <-sseClient.Errors():
				if !ok {
					log.Println("[EventHandler] SSE errors channel closed.")
					return
				}
				log.Printf("[EventHandler] Received SSE error: %v", err)
			}
		}
	}()
}

// Stop signals the event handling goroutine to terminate.
func (h *EventHandler) Stop() {
	close(h.stopChan)
}

// handleEvent is now a private method.
func (h *EventHandler) handleEvent(event client.Event) {
	var notification struct {
		Type    string          `json:"type"`
		Entity  string          `json:"entity"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(event.Data, &notification); err != nil {
		log.Printf("[EventHandler] Error parsing notification: %v", err)
		return
	}

	// Route the event based on its entity and type.
	switch notification.Entity {
	case "assignment":
		switch notification.Type {
		case "create":
			if h.onAssignmentCreate != nil {
				h.onAssignmentCreate(notification.Data, notification.Message)
			}
		case "update":
			if h.onAssignmentUpdate != nil {
				h.onAssignmentUpdate(notification.Data, notification.Message)
			}
		case "delete":
			if h.onAssignmentDelete != nil {
				h.onAssignmentDelete(notification.Data, notification.Message)
			}
		}
	case "course":
		// Placeholder for future course event handling.
	}
}

// OnAssignmentCreate registers a handler function for assignment creation events.
func (h *EventHandler) OnAssignmentCreate(f func(data json.RawMessage, message string)) {
	h.onAssignmentCreate = f
}

// OnAssignmentUpdate registers a handler function for assignment update events.
func (h *EventHandler) OnAssignmentUpdate(f func(data json.RawMessage, message string)) {
	h.onAssignmentUpdate = f
}

// OnAssignmentDelete registers a handler function for assignment deletion events.
func (h *EventHandler) OnAssignmentDelete(f func(data json.RawMessage, message string)) {
	h.onAssignmentDelete = f
}
