package main

import (
	"fmt"
	"log"

	"github.com/williamfotso/acc/internal/storage/local"
)

func main() {
	userID, err := local.GetCurrentUserID()
	if err != nil {
		log.Fatalf("Failed to get current user ID: %v", err)
	}

	fmt.Printf("Current user ID: %d\n", userID)

	// 1. Initialize database
	db, err := local.GetLocalDB(userID)
	if err != nil {
		log.Fatalf("Failed to create local DB: %v", err)
	}
	defer local.CloseAll()
	fmt.Println("✅ Database connection established")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := local.InitializeSchema(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to seed initial data: %v", err)
	}
	fmt.Println("✅ Initial data seeded")

	tx.Commit()

}
