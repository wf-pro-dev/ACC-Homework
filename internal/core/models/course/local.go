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
	Code       string `gorm:"unique"`
	Name       string `gorm:"not null"`
	NotionID   string `gorm:"unique"` // Empty until synced
	Duration   string
	RoomNumber string
	SyncStatus SyncStatus `gorm:"not null;default:'pending'"`
}

func (c *LocalCourse) ToMap() map[string]string {
	return map[string]string{
		"code":        c.Code,
		"name":        c.Name,
		"notion_id":   c.NotionID,
		"duration":    c.Duration,
		"room_number": c.RoomNumber,
		"sync_status": string(c.SyncStatus),
	}
}
