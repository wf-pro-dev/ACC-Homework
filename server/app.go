package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartServer() {

	http.HandleFunc("/notion-webhooks/test", testHandler)
	http.HandleFunc("/notion-webhooks", notionWebhookHandler)
	log.Println("Server listening on :8080...") // Changed from fmt
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s test: %s", time.Now().Format(time.Stamp), r.URL.Path)

	var payload struct {
		Test string `json:"test"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("%s test: %s", time.Now().Format(time.Stamp), payload.Test)
}

func PrintLog(message string) {
	log.Printf("[INFO] %s", message)
}

func PrintERROR(w http.ResponseWriter, code int, message string) {
	log.Printf("[ERROR] [%d] %s", code, message)
	http.Error(w, message, code)
}

func notionWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verify it's a POST request
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode the payload
	var payload NotionWebhookPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		PrintERROR(w, http.StatusBadRequest, fmt.Sprintf("Error decoding payload: %s", err))
		return
	}

	// 3. Get the author of the change
	var author_type string
	if len(payload.Authors) > 0 {
		author_type = payload.Authors[0].Type
		if author_type == "bot" {
			PrintLog("Author is a bot, skipping")
			return
		}
	}

	// 4. Handle the payload
	switch payload.Type {
	case "page.properties_updated":
		WebhookUpdateHandler(w, r, payload)
	case "page.created":
		WebhookCreateHandler(w, r, payload)
	case "page.deleted":
		WebhookDeleteHandler(w, r, payload)
	default:
		PrintERROR(w, http.StatusBadRequest, fmt.Sprintf("Unknown payload type: %s", payload.Type))
	}
}
