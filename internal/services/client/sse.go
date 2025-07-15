package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type SSEClient struct {
	events        chan Event
	errors        chan error
	disconnect    chan struct{}
	connected     bool
	mu            sync.Mutex
	httpClient    *http.Client
	retryCount    int
	maxRetries    int
	baseDelay     time.Duration
	reconnect     chan struct{}
	lastEventTime time.Time
}

type Event struct {
	Type string
	Data json.RawMessage
}

func NewSSEClient(httpClient *http.Client) *SSEClient {
	return &SSEClient{
		events:     make(chan Event),
		errors:     make(chan error),
		disconnect: make(chan struct{}),
		httpClient: httpClient,
		maxRetries: 5,               // Maximum reconnection attempts
		baseDelay:  2 * time.Second, // Initial delay between retries
		reconnect:  make(chan struct{}, 1),
	}
}

func (c *SSEClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	go c.manageConnection()
	return nil
}

func (c *SSEClient) manageConnection() {
	for {
		select {
		case <-c.disconnect:
			return
		default:
			err := c.establishConnection()
			if err == nil {
				// Successful connection
				c.retryCount = 0
				continue
			}

			// Handle reconnection
			if c.retryCount >= c.maxRetries {
				c.errors <- fmt.Errorf("max reconnection attempts reached")
				return
			}

			delay := c.calculateBackoff()
			time.Sleep(delay)
			c.retryCount++
		}
	}
}

func (c *SSEClient) establishConnection() error {
	req, err := http.NewRequest("GET", "https://newsroom.dedyn.io/acc-homework/events", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE connection failed with status: %d", resp.StatusCode)
	}

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		resp.Body.Close()
	}()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-c.disconnect:
			return nil
		case <-c.reconnect:
			return fmt.Errorf("reconnection requested")
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return fmt.Errorf("server closed connection")
				}
				return fmt.Errorf("error reading SSE: %w", err)
			}

			if len(line) > 0 {
				event := c.parseEvent(line)
				if event != nil {
					c.lastEventTime = time.Now()
					c.events <- *event
				}
			}
		}
	}
}

func (c *SSEClient) calculateBackoff() time.Duration {
	if c.retryCount == 0 {
		return c.baseDelay
	}
	return c.baseDelay * time.Duration(c.retryCount*c.retryCount)
}

func (c *SSEClient) RequestReconnect() {
	select {
	case c.reconnect <- struct{}{}:
	default:
	}
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

func (c *SSEClient) parseEvent(data []byte) *Event {
	var eventType string
	var eventData json.RawMessage

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("event:")) {
			eventType = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("event:"))))
		} else if bytes.HasPrefix(line, []byte("data:")) {
			eventData = bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		} else if bytes.HasPrefix(line, []byte(":")) {
			// Heartbeat comment, update last activity
			c.lastEventTime = time.Now()
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
