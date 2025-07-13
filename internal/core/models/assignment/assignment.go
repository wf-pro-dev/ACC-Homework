package assignment

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/williamfotso/acc/internal/core/models"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/core/models/user"
	"gorm.io/gorm"
)

// Assignment represents a homework or exam assignment
type Assignment struct {
	gorm.Model
	UserID     uint
	NotionID   string    `gorm:"unique"`
	Title      string    `gorm:"not null"`
	Todo       string
	Deadline   time.Time `gorm:"not null"`
	Link       string    `gorm:"default:https://acconline.austincc.edu/ultra/stream"`
	CourseCode string
	TypeName   string                  `gorm:"not null"`
	StatusName string                  `gorm:"not null"`

	User       user.User `gorm:"foreignKey:UserID;references:ID"`
	Course     course.Course           `gorm:"foreignKey:CourseCode;references:Code"`
	Type       models.AssignmentType   `gorm:"foreignKey:TypeName;references:Name"`
	Status     models.AssignmentStatus `gorm:"foreignKey:StatusName;references:Name"`
}

type Filter struct {
	Column string
	Value  string
}

// NewAssignment creates a new Assignment by prompting user for input
// This is equivalent to the createAssign function but returns a struct
func NewAssignment() *Assignment {

	fmt.Println("===== Creating new Assignement =====")

	assignment := &Assignment{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("The type (HW or Exam): ")
	scanner.Scan()
	assignment.TypeName = scanner.Text()

	fmt.Printf("The deadline (yyyy-mm-dd): ")
	scanner.Scan()
	deadline, err := time.Parse(time.DateOnly, scanner.Text())
	if err != nil {
		log.Fatal("Error parsing deadline: ", err)
	}
	assignment.Deadline = deadline

	fmt.Printf("The title: ")
	scanner.Scan()
	assignment.Title = scanner.Text()

	fmt.Printf("The todo: ")
	scanner.Scan()
	assignment.Todo = scanner.Text()

	// Get course code from current directory name
	pwd := os.Getenv("PWD")
	cmd := exec.Command("basename", pwd)
	output, _ := cmd.CombinedOutput()
	assignment.CourseCode = strings.TrimSpace(string(output))

	assignment.Link = "https://acconline.austincc.edu/ultra/stream"

	return assignment
}


func GetAssignmentsbyCourse(course_code string, columns []string, filters []Filter, up_to_date bool, db *gorm.DB) {

	col_length := 15
	query := fmt.Sprintf("SELECT %s FROM assignements WHERE course_code='%v'", strings.Join(columns, ","), course_code)

	for _, filter := range filters {
		query += fmt.Sprintf(" AND %s='%v'", filter.Column, filter.Value)
	}

	if up_to_date {
		query += " AND deadline > NOW()"
	}
	query += " ORDER BY deadline ASC"
	assignments := []Assignment{}
	err := db.Raw(query).Scan(&assignments).Error
	if err != nil {
		log.Fatal(err)
	}

	if len(assignments) == 0 {
		fmt.Println("No assignments found")
		os.Exit(0)
	}

	// Create column headers based on requested columns
	headers := make([]string, len(columns))
	for i, col := range columns {
		// Convert column names to display headers
		switch col {
		case "id":
			headers[i] = "ID"
		case "type":
			headers[i] = "Type"
		case "deadline":
			headers[i] = "Deadline"
		case "title":
			headers[i] = "Title"
		case "todo":
			headers[i] = "Todo"
		case "course_code":
			headers[i] = "Course Code"
		case "notion_id":
			headers[i] = "Notion ID"
		case "status":
			headers[i] = "Status"
		default:
			headers[i] = col
		}
	}

	// Print top border
	fmt.Print("┌")
	for range columns {
		fmt.Printf("%-*s┬", col_length, strings.Repeat("-", col_length+2))
	}
	fmt.Println("")

	// Print header row
	fmt.Print("│")
	for _, header := range headers {
		fmt.Printf(" %-*s │", col_length, header)
	}
	fmt.Println("")

	// Print separator
	fmt.Print("├")
	for range columns {
		fmt.Printf("%-*s┼", col_length, strings.Repeat("-", col_length+2))
	}
	fmt.Println("")

	// Print data rows
	for _, assignment := range assignments {
		obj_assign := assignment.ToMap()
		fmt.Print("│")
		for _, col := range columns {
			value := obj_assign[col]
			if col == "deadline" {
				value = value[:10]
			}

			// Truncate or pad to exactly 10 characters
			if len(value) > 15 && len(columns) > 2 {
				value = value[:12] + "..."
			}
			fmt.Printf(" %-*s │", col_length, value)
		}
		fmt.Println("")
	}

	// Print bottom border
	fmt.Print("└")
	for range columns {
		fmt.Printf("%-*s┴", col_length, strings.Repeat("-", col_length+2))
	}
	fmt.Println("")
}

func Get_Assignment_byId(id uint, db *gorm.DB) (*Assignment, error) {
    assignment := &Assignment{}
    err := db.Preload("User").
              Preload("Course").
              Preload("Type").
              Preload("Status").
              Where("id = ?", id).
              First(assignment).Error

    if err != nil {
        return nil, err
    }
    return assignment, nil
}

func Get_Assignment_byNotionID(notion_id string, db *gorm.DB) ( *Assignment, error ){

	assignment := &Assignment{}
	err := db.Where("notion_id = ?", notion_id).First(assignment).Error

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return assignment, nil
}

// ToMap converts the Assignment struct to a map[string]string
// This maintains compatibility with the existing database operations
func (a *Assignment) ToMap() map[string]string {

	return map[string]string{
		"id":          strconv.Itoa(int(a.ID)),
		"user_id":     strconv.Itoa(int(a.UserID)),
		"notion_id":   a.NotionID,
		"type":        a.TypeName,
		"deadline":    a.Deadline.Format(time.DateOnly),
		"title":       a.Title,
		"todo":        a.Todo,
		"course_code": a.CourseCode,
		"link":        a.Link,
		"status":      a.StatusName,
		"created_at":  a.CreatedAt.Format(time.DateOnly),
		"updated_at":  a.UpdatedAt.Format(time.DateOnly),
	}
}

func (a *Assignment) Add(db *gorm.DB) (err error) {

	assignment := a.ToMap()

	delete(assignment, "id")

	err = db.Create(a).Error

	if err != nil {
		log.Fatalln("Error adding assignment to database: ", err)
		return err
	}

	notion_id, err_notion := a.Add_Notion()

	if err_notion != nil {
		log.Fatalln("Error adding assignment to Notion: ", err_notion)
		return err_notion
	}

	var lastVal int
	err = db.Raw("SELECT MAX(id) FROM assignements").Scan(&lastVal).Error
	if err != nil {
		log.Fatal(err)
	}
	err = db.Model(&Assignment{}).Where("id = ?", lastVal).Update("notion_id", notion_id).Error

	if err != nil {
		log.Fatalln("Error updating assignment: ", err)
		return err
	}

	return nil
}


func (a *Assignment) Update(col, value string, db *gorm.DB) (err error) {

	err = db.Model(&Assignment{}).Where("id = ?", a.ID).Update(col, value).Error

	if err != nil {
		log.Fatalln("Error updating assignment in database: ", err)
		return err
	}

	assignment := a.ToMap()
	assignment[col] = value

	if col == "course_code" {
		value = course.Get_Course_byCode(value, db).NotionID
	}

	var obj map[string]string

	if col == "status" {
		obj = models.Get_AssignmentStatus_byName(value, db).ToMap()
	} else {
		obj = models.Get_AssignmentType_byName(value, db).ToMap()
	}

	err = a.Update_Notion(col, value, obj)

	if err != nil {
		log.Fatalln("Error updating assignment to Notion: ", err)
		return err
	}

	return nil
}

func (a *Assignment) Delete(db *gorm.DB) (err error) {

	err = db.Delete(a).Error 

	if err != nil {
		log.Fatalln(err)
	}

	err = a.Delete_Notion()
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
