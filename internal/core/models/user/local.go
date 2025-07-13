package user

import (
	"fmt"
	"gorm.io/gorm"
)

// User represents the application user
type LocalUser struct {
	gorm.Model
	Username     	string `gorm:"unique;not null"`
	Email        	string `gorm:"unique;not null"`
	PasswordHash 	string `gorm:"not null"`
	NotionAPIKey 	string // Encrypted in application layer
	AssignmentsDbId string
	NotionID	string
	CoursesDbId	string
}

func (u *LocalUser) ToMap() map[string]interface{} {
	if u == nil {
		return nil
	}

	return map[string]interface{}{
		"id":              u.ID,
		"username":        u.Username,
		"email":           u.Email,
		"assignments_db":  u.AssignmentsDbId,
		"courses_db":      u.CoursesDbId,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
	}
}


func Get_Local_User_by_NotionID(notion_id string, db *gorm.DB) (*LocalUser, error) {
	u := &LocalUser{}
	err := db.Where("notion_id = ?", notion_id).First(u).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting user with notion id: ", err)
	}
	return u, nil
}


