package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/williamfotso/acc/internal/services/client"
	"github.com/williamfotso/acc/internal/storage/local"
)

var (
	// This will be managed by the 'listen' command's lifecycle.
	sseCancelFunc context.CancelFunc
)

// Login handles only authentication and saving the session cookie to a file.
func Login(username, password string) error {
	// Create a new http client for the login request.
	httpClient, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("could not create http client: %w", err)
	}

	loginData := map[string]string{"username": username, "password": password}
	jsonData, _ := json.Marshal(loginData)

	resp, err := httpClient.Post("https://newsroom.dedyn.io/acc-homework/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("http post failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		UserID string `json:"user_id"`
		Error  string `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}

	if err := client.SaveCookies(httpClient); err != nil {
		return fmt.Errorf("failed to save cookies: %w", err)
	}

	//==== User specific configuration ====

	userID, err := strconv.Atoi(response.UserID)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	if err := local.StoreCredentials(uint(userID), username); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	db, err := local.GetLocalDB(uint(userID))
	if err != nil {
		return fmt.Errorf("failed to create/get local database: %w", err)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := local.InitializeSchema(tx); err != nil {
		return fmt.Errorf("failed to initialize user database: %w", err)
	}

	if err := client.MigrateCourses(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to migrate courses: %w", err)
	}

	if err := client.MigrateAssignments(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to migrate assignments: %w", err)
	}

	tx.Commit()

	//==== User specific configuration ====

	// Save the session cookies from the client's jar to a file.
	// The 'listen' command will load these cookies to authenticate its requests.
	return nil
}

// StartListener initializes the SSE client using stored cookies and starts the connection.
func StartListener(ctx context.Context) (*client.SSEClient, error) {
	// Create a client and load the cookies saved by the login command.
	httpClient, err := client.NewClientWithCookies()
	if err != nil {
		return nil, fmt.Errorf("could not create http client for listener: %w", err)
	}

	// Create a new cancellable context for the SSE client.
	var sseCtx context.Context
	sseCtx, sseCancelFunc = context.WithCancel(ctx)

	sseClient := client.NewSSEClient(httpClient)

	// The connect method now accepts a context for graceful shutdown.
	go sseClient.Connect(sseCtx)

	return sseClient, nil
}

// StopListener signals the SSE connection to close.
func StopListener() {
	if sseCancelFunc != nil {
		log.Println("Signaling SSE client to disconnect...")
		sseCancelFunc() // This cancels the context passed to sseClient.Connect
	}
}

// GetDaemonPIDFilePath returns the canonical path for the daemon's PID file.
func GetDaemonPIDFilePath() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find user home directory: %w", err)
	}
	configDir := filepath.Join(home, "acc-homework")
	// Ensure the directory exists.
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config directory at %s: %w", configDir, err)
	}
	return filepath.Join(configDir, "daemon.pid"), nil
}
