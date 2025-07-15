package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/events"
	"github.com/williamfotso/acc/internal/storage/local"
	"github.com/williamfotso/acc/internal/types"
	"gorm.io/gorm"
)

var eventHandler *events.EventHandler

func init() {
	// Initialize event handler
	eventHandler = events.NewEventHandler()

	// Set up event handlers
	eventHandler.OnAssignmentCreate(events.HandleAssignmentCreate)
	eventHandler.OnAssignmentUpdate(events.HandleAssignmentUpdate)
	eventHandler.OnAssignmentDelete(events.HandleAssignmentDelete)

	// Start listening for events if logged in
	if _, err := local.GetCurrentUserID(); err == nil {
		startEventHandling()
	}
}

func startEventHandling() {
	// Ensure clean shutdown on exit
	cobra.OnFinalize(auth.CleanupSSE)

	// Start event handler
	eventHandler.Start()

	// Monitor SSE connection
	go func() {
		sseClient := auth.GetSSEClient()
		if sseClient == nil {
			return
		}

		for {
			select {
			case event := <-sseClient.Events():
				eventHandler.HandleEvent(event)
			case err := <-sseClient.Errors():
				log.Printf("SSE error: %v", err)
				// Consider automatic reconnection here
			}
		}
	}()
}

func ValidateAssignmentId(id string, db *gorm.DB) error {

	if id == "" {
		return fmt.Errorf("assignment ID is required")
	}

	int_id, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("failed to convert assignment ID to int: %s", err)
	}

	assignment, _ := assignment.Get_Local_Assignment_byId(uint(int_id), db)
	if assignment == nil {
		return fmt.Errorf("assignment not found")
	}

	return nil
}

func ColumnError(message string) {

	fmt.Println(message)
	fmt.Println("Available columns:")
	for _, column := range types.COLUMNS {
		fmt.Printf("  -%s (%s)\n", column[0:2], column)
	}
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "acc",
	Short: "ACC is a CLI tool for managing assignments and courses",
	Long: `ACC is a CLI tool for managing assignments and courses from austin community college.
				  Complete documentation is available at https://github.com/williamfotso/acc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func Execute() {
	if _, err := local.GetCurrentUserID(); err == nil {
		startEventHandling()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
