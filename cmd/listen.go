package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/events"
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

var listenCmd = &cobra.Command{
	Use:    "listen",
	Short:  "Background listener (auto-started by login)",
	Hidden: true, // Hide from help since users shouldn't call this directly
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the listener service
		sseClient, err := auth.StartListener(ctx)
		if err != nil {
			log.Fatalf("Failed to start listener: %v", err)
		}

		// Initialize event handler
		eventHandler := events.NewEventHandler()
		eventHandler.OnAssignmentCreate(events.HandleAssignmentCreate)
		eventHandler.OnAssignmentUpdate(events.HandleAssignmentUpdate)
		eventHandler.OnAssignmentDelete(events.HandleAssignmentDelete)
		eventHandler.Start(sseClient)

		// Wait for shutdown signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down listener...")
		eventHandler.Stop()
		auth.StopListener()
	},
}
