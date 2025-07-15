// internal/services/client/sse.go
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/williamfotso/acc/internal/storage/local"
)

type SSEClient struct {
	events     chan Event
	errors     chan error
	disconnect chan struct{}
	connected  bool
	mu         sync.Mutex
}

type Event struct {
	Type string
	Data json.RawMessage
}

func NewSSEClient() *SSEClient {
	return &SSEClient{
		events:     make(chan Event),
		errors:     make(chan error),
		disconnect: make(chan struct{}),
	}
}

func (c *SSEClient) Connect() error {
	userID, err := local.GetCurrentUserID()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	// Create new HTTP client with cookies
	httpClient, err := NewClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	req, err := http.NewRequest("GET", "https://newsroom.dedyn.io/acc-homework/events", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSE connection failed with status: %d", resp.StatusCode)
	}

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	go func() {
		defer resp.Body.Close()
		defer close(c.events)
		defer close(c.errors)

		reader := bufio.NewReader(resp.Body)
		for {
			select {
			case <-c.disconnect:
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						c.errors <- fmt.Errorf("SSE connection closed by server")
					} else {
						c.errors <- fmt.Errorf("error reading SSE: %w", err)
					}
					return
				}

				// Parse SSE event
				if len(line) > 0 {
					event := parseEvent(line)
					if event != nil {
						c.events <- *event
					}
				}
			}
		}
	}()

	return nil
}

func (c *SSEClient) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		close(c.disconnect)
		c.connected = false
	}
}

func (c *SSEClient) Events() <-chan Event {
	return c.events
}

func (c *SSEClient) Errors() <-chan error {
	return c.errors
}

func (c *SSEClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func parseEvent(data []byte) *Event {
	// Simple parser for SSE format
	// In a real implementation, you'd want something more robust
	var eventType string
	var eventData json.RawMessage

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("event:")) {
			eventType = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("event:"))))
		} else if bytes.HasPrefix(line, []byte("data:")) {
			eventData = bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		}
	}

	if len(eventData) > 0 {
		return &Event{
			Type: eventType,
			Data: eventData,
		}
	}

	return nil
}
