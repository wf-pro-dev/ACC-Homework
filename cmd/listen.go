package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/events"
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start the background listener for notifications",
	Long:  `This command starts a background process to listen for real-time events. It will run until you explicitly call 'acc logout'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get path for the PID file.
		pidFile, err := auth.GetDaemonPIDFilePath()
		if err != nil {
			log.Fatalf("Could not get PID file path: %v", err)
		}

		// Check if the daemon is already running.
		if _, err := os.Stat(pidFile); err == nil {
			fmt.Println("Listener is already running.")
			return
		}

		// Create the daemon process. The `&` is not needed in Go.
		// We re-execute the command with a special flag to run as a daemon.
		// This is a common pattern to detach a process.
		// For simplicity in this example, we will run it in the foreground
		// and let the user background it with `&` in their shell.
		fmt.Println("Starting listener... Press Ctrl+C or run 'acc logout' to stop.")

		// Write the current process ID to the PID file.
		pid := os.Getpid()
		if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
			log.Fatalf("Could not create config directory: %v", err)
		}
		if err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Fatalf("Could not write PID file: %v", err)
		}
		// Ensure PID file is removed on exit.
		defer os.Remove(pidFile)

		// Set up a context that is cancelled on shutdown signals.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the listener service.
		sseClient, err := auth.StartListener(ctx)
		if err != nil {
			log.Fatalf("Failed to start listener: %v", err)
		}

		// Initialize and start the event handler.
		eventHandler := events.NewEventHandler()
		eventHandler.OnAssignmentCreate(events.HandleAssignmentCreate)
		eventHandler.OnAssignmentUpdate(events.HandleAssignmentUpdate)
		eventHandler.OnAssignmentDelete(events.HandleAssignmentDelete)
		eventHandler.Start(sseClient) // Start listening for events.

		// --- Graceful Shutdown ---
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Block until a signal is received.
		<-quit

		log.Println("Shutdown signal received, stopping listener...")

		// Stop the event handler and SSE client.
		eventHandler.Stop()
		auth.StopListener()

		log.Println("Listener stopped gracefully.")
	},
}
