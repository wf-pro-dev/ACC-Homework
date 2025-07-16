// cmd/acc/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/events"
)

func main() {
	log.Println("ACC Application started. Available commands: login, logout, exit")

	// eventHandler will be nil until the user logs in.
	var eventHandler *events.EventHandler
	scanner := bufio.NewScanner(os.Stdin)

	// This loop runs continuously, acting as a simple command prompt.
	for {
		// Only show the prompt if we are not in the middle of a command
		fmt.Print("> ")

		// Wait for user input
		if !scanner.Scan() {
			break // Exit on EOF (Ctrl+D)
		}
		command := strings.TrimSpace(scanner.Text())

		switch command {
		case "login":
			if eventHandler != nil {
				log.Println("Already logged in. Please 'logout' first.")
				continue
			}

			// In a real app, you would prompt for username/password.
			username := "testuser"
			password := "securepassword"
			log.Printf("[Main] Attempting to log in as user: %s\n", username)

			// auth.Login now handles starting the SSE connection.
			err := auth.Login(username, password)
			if err != nil {
				// Use log.Printf for non-fatal errors to keep the CLI running.
				log.Printf("[Main] Login failed: %v", err)
				continue
			}
			log.Println("[Main] Login successful. SSE client is connecting in the background.")

			// Once logged in, initialize and start the event handler.
			eventHandler = events.NewEventHandler()

			// Register the functions to handle specific events.
			eventHandler.OnAssignmentCreate(events.HandleAssignmentCreate)
			eventHandler.OnAssignmentUpdate(events.HandleAssignmentUpdate)
			eventHandler.OnAssignmentDelete(events.HandleAssignmentDelete)

			// Start the handler's goroutine to listen for events from the SSE client.
			eventHandler.Start()
			log.Println("[Main] Event listener started. Receiving events...")

		case "logout":
			if eventHandler == nil {
				log.Println("You are not logged in.")
				continue
			}

			log.Println("[Main] Initiating logout...")

			// 1. Stop the event handler first to prevent it from processing
			//    more events from a connection that is about to close.
			eventHandler.Stop()
			log.Println("[Main] Event listener stopped.")

			// 2. Call the Logout function to close the SSE connection and clean up.
			auth.Logout()

			// 3. Clear the local eventHandler instance.
			eventHandler = nil

		case "exit":
			log.Println("[Main] Exiting application...")
			// If the user is still logged in, perform a clean logout first.
			if eventHandler != nil {
				log.Println("[Main] Performing automatic logout before exiting...")
				eventHandler.Stop()
				auth.Logout()
			}
			return // Exit the main function and the program.

		default:
			if command != "" {
				log.Printf("Unknown command: '%s'. Available commands: login, logout, exit.", command)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}

	log.Println("Application shut down.")
}
