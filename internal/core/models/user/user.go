package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents the application user
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	NotionAPIKey string // Encrypted in application layer
	LastSync     *time.Time
}
