package course

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/williamfotso/acc/internal/core/models/user"
	"gorm.io/gorm"
)

// Course represents a school course
type Course struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`
	User       user.User `gorm:"foreignKey:UserID;references:ID"`
	NotionID   string    `gorm:"unique;not null"`
	Code       string    `gorm:"unique;not null"`
	Name       string    `gorm:"not null"`
	Duration   string
	RoomNumber string
}

func NewCourse() *Course {
	fmt.Println("===== Creating new Course =====")

	course := &Course{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("The Name: ")
	scanner.Scan()
	course.Name = scanner.Text()

	fmt.Printf("The Code: ")
	scanner.Scan()
	course.Code = scanner.Text()

	fmt.Printf("The Room Number: ")
	scanner.Scan()
	course.RoomNumber = scanner.Text()

	fmt.Printf("The Duration: ")
	scanner.Scan()
	course.Duration = scanner.Text()

	return course
}

func (c *Course) SetNotionID(notionID string) {
	c.NotionID = notionID
}

func (c *Course) GetNotionID() string {
	return c.NotionID
}

func (c *Course) SetName(name string) {
	c.Name = name
}

func (c *Course) GetName() string {
	return c.Name
}

func (c *Course) SetCode(code string) {
	c.Code = code
}

func (c *Course) GetCode() string {
	return c.Code
}

func (c *Course) SetRoomNumber(roomNumber string) {
	c.RoomNumber = roomNumber
}

func (c *Course) GetRoomNumber() string {
	return c.RoomNumber
}

func (c *Course) SetDuration(duration string) {
	c.Duration = duration
}

func (c *Course) GetDuration() string {
	return c.Duration
}


func Get_Course_byId(id uint, db *gorm.DB) (*Course, error) {
	course := &Course{}
	err := db.Preload("User").
		Where("id = ?", id).
		First(course).Error

	if err != nil {
		return nil, err
	}
	return course, nil
}
func Get_Course_byCode(code string, db *gorm.DB) *Course {
	course := &Course{}
	err := db.Where("code = ?", code).First(course).Error
	if err != nil {
		log.Fatalln("Error getting course with code: ", err)
		return nil
	}

	return course
}

func Get_Course_byNotionID(notion_id string, db *gorm.DB) *Course {

	course := &Course{}
	err := db.Where("notion_id = ?", notion_id).First(course).Error
	if err != nil {
		log.Fatalln("Error getting course with notion id: ", err)
		return nil
	}

	return course
}

func (c *Course) ToMap() map[string]string {
	return map[string]string{
		"id":          strconv.Itoa(int(c.ID)),
		"user_id":     strconv.Itoa(int(c.UserID)),
		"notion_id":   c.NotionID,
		"name":        c.Name,
		"code":        c.Code,
		"room_number": c.RoomNumber,
		"duration":    c.Duration,
		"created_at":  c.CreatedAt.Format(time.DateOnly),
		"updated_at":  c.UpdatedAt.Format(time.DateOnly),
	}
}
