// internal/services/events/listener.go
package events

import (
	"encoding/json"
	"log"

	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/client"
)

type EventHandler struct {
	onAssignmentCreate func(data json.RawMessage)
	onAssignmentUpdate func(data json.RawMessage)
	onAssignmentDelete func(data json.RawMessage)
	// onCourseCreate     func(data json.RawMessage)
	// onCourseUpdate     func(data json.RawMessage)
	// onCourseDelete     func(data json.RawMessage)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (h *EventHandler) Start() {
	go func() {
		sseClient := auth.GetSSEClient()
		if sseClient == nil {
			return
		}

		for {
			select {
			case event, ok := <-sseClient.Events():
				if !ok {
					return
				}
				h.handleEvent(event)
			case err, ok := <-sseClient.Errors():
				if !ok {
					return
				}
				log.Printf("SSE error: %v", err)
			}
		}
	}()
}

func (h *EventHandler) handleEvent(event client.Event) {
	var notification struct {
		Type    string          `json:"type"`
		Entity  string          `json:"entity"`
		ID      string          `json:"id"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(event.Data, &notification); err != nil {
		log.Printf("Error parsing notification: %v", err)
		return
	}

	switch notification.Entity {
	case "assignment":
		switch notification.Type {
		case "create":
			if h.onAssignmentCreate != nil {
				h.onAssignmentCreate(notification.Data)
			}
		case "update":
			if h.onAssignmentUpdate != nil {
				h.onAssignmentUpdate(notification.Data)
			}
		case "delete":
			if h.onAssignmentDelete != nil {
				h.onAssignmentDelete(notification.Data)
			}
		}
	case "course":
		// Similar logic for courses
	}
}

// Set up handler functions
func (h *EventHandler) OnAssignmentCreate(f func(data json.RawMessage)) {
	h.onAssignmentCreate = f
}

func (h *EventHandler) OnAssignmentUpdate(f func(data json.RawMessage)) {
	h.onAssignmentUpdate = f
}

func (h *EventHandler) OnAssignmentDelete(f func(data json.RawMessage)) {
	h.onAssignmentDelete = f
}

// ... similar for course events ...
