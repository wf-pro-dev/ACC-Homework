package assignment

import (
        "time"

	"github.com/williamfotso/acc/internal/core/models"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/core/models/user"
        "gorm.io/gorm"
)

type SyncStatus string

const (
    SyncStatusPending SyncStatus = "pending" // Needs to be synced
    SyncStatusSynced  SyncStatus = "synced"  // Already synced
)

type LocalAssignment struct {
	gorm.Model
	NotionID    string    `gorm:"unique"`     // Empty until synced
	Title       string    `gorm:"not null"`
	Todo        string
	Deadline    time.Time `gorm:"not null;index"`
	Link        string    `gorm:"default:https://acconline.austincc.edu/ultra/stream"`
	CourseCode  string    `gorm:"not null;index"`
	TypeName    string    `gorm:"not null"`
	StatusName  string    `gorm:"not null"`
	UserID      uint      `gorm:"not null;index"`
	SyncStatus  SyncStatus `gorm:"not null;default:'pending'"`

	User       user.LocalUser `gorm:"foreignKey:UserID;references:ID"`
	Course     course.LocalCourse           `gorm:"foreignKey:CourseCode;references:Code"`
	Type       models.LocalAssignmentType   `gorm:"foreignKey:TypeName;references:Name"`
	Status     models.LocalAssignmentStatus `gorm:"foreignKey:StatusName;references:Name"`
}


func ( a *LocalAssignment ) ToMap() map[string]string {
	return map[string]string{

                "course_code":	a.CourseCode,
		"title":	a.Title,
                "type_name":	a.TypeName,
                "deadline":	a.Deadline.Format(time.DateOnly),
		"todo":		a.Todo,
		"status_name":	a.StatusName,
		"link":		a.Link,
	}
}
