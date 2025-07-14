package main

import (
	"fmt"

	"github.com/williamfotso/acc/internal/storage/local"
)

func main() {
	err := local.StoreCredentials(1, "William Fotso")
	if err != nil {
		fmt.Printf("Failed to store credentials: %v\n", err)
	}
	fmt.Printf("Credentials stored successfully\n")

	userID, err := local.GetCurrentUserID()
	if err != nil {
		fmt.Printf("Failed to get current user ID: %v\n", err)
	}
	fmt.Printf("The user ID: %v\n", userID)
}
