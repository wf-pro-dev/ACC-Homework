package user

import (
	"time"
	"gorm.io/gorm"
)

// User represents the application user
type User struct {
	gorm.Model
	Username     	string `gorm:"unique;not null"`
	Email        	string `gorm:"unique;not null"`
	PasswordHash 	string `gorm:"not null"`
	NotionAPIKey 	string // Encrypted in application layer
	AssignmentsDbId string
	CoursesDbId	string
	LastSync     	*time.Time
}

func (u *User) ToMap() map[string]interface{} {
	if u == nil {
		return nil
	}

	return map[string]interface{}{
		"id":              u.ID,
		"username":        u.Username,
		"email":           u.Email,
		"assignments_db":  u.AssignmentsDbId,
		"courses_db":      u.CoursesDbId,
		"last_sync":       u.LastSync,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
	}
}
