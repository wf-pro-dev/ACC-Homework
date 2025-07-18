package notes

import (
	"time"

	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/core/models/user"
	"gorm.io/gorm"
)

type Notes struct {
	gorm.Model
	Title      string
	Keywords   string
	CourseCode string
	UserID     uint
	NotionID   string    `gorm:"unique"`
	Date       time.Time `gorm:"not null"`
	Transcript string    `gorm:"not null"`
	Summary    string

	User   user.User     `gorm:"foreignKey:UserID;references:ID"`
	Course course.Course `gorm:"foreignKey:CourseCode;references:Code"`
}
