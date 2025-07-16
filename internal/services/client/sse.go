package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type SSEClient struct {
	httpClient *http.Client
	events     chan Event
	errors     chan error
	mu         sync.Mutex
}

type Event struct {
	Type string
	Data json.RawMessage
}

func NewSSEClient(httpClient *http.Client) *SSEClient {
	return &SSEClient{
		httpClient: httpClient,
		events:     make(chan Event, 1), // Buffered channel
		errors:     make(chan error, 1), // Buffered channel
	}
}

// Connect now accepts a context to handle cancellation.
func (c *SSEClient) Connect(ctx context.Context) {
	defer func() {
		close(c.events)
		close(c.errors)
		log.Println("[SSEClient] Connection loop terminated and channels closed.")
	}()

	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, exit the connection loop.
			return
		default:
			log.Println("[SSEClient] Attempting to establish SSE connection...")
			err := c.establishAndStream(ctx)
			if err != nil {
				// Don't push to error channel if it was a graceful shutdown.
				if ctx.Err() == nil {
					log.Printf("[SSEClient] Connection error: %v. Retrying in 5 seconds...", err)
					c.errors <- err
				}
			}

			// Wait before retrying, but exit immediately if cancelled.
			select {
			case <-time.After(5 * time.Second):
				// Continue to the next iteration of the loop.
			case <-ctx.Done():
				// Exit immediately.
				return
			}
		}
	}
}

func (c *SSEClient) establishAndStream(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://newsroom.dedyn.io/acc-homework/events", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	log.Println("[SSEClient] Connection established. Streaming events...")
	reader := bufio.NewReader(resp.Body)
	for {
		// Check for context cancellation before each read.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return errors.New("server closed connection (EOF)")
			}
			return fmt.Errorf("error reading from stream: %w", err)
		}

		if bytes.HasPrefix(line, []byte("data:")) {
			data := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
			if len(data) > 0 {
				c.events <- Event{Data: data}
			}
		}
	}
}

// Events returns the read-only channel for receiving events.
func (c *SSEClient) Events() <-chan Event {
	return c.events
}

// Errors returns the read-only channel for receiving errors.
func (c *SSEClient) Errors() <-chan error {
	return c.errors
}
