package course

import (
	
        "gorm.io/gorm"
)

type SyncStatus string

const (
    SyncStatusPending SyncStatus = "pending" // Needs to be synced
    SyncStatusSynced  SyncStatus = "synced"  // Already synced
)

type LocalCourse struct {
	gorm.Model
	Code        string    `gorm:"unique"`
	Name        string    `gorm:"not null"`
	NotionID    string    `gorm:"unique"`     // Empty until synced
	Duration    string
	RoomNumber  string
	UserID      uint      `gorm:"not null;index"`
	SyncStatus  SyncStatus `gorm:"not null;default:'pending'"`
}
