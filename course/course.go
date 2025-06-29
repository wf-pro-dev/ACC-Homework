package course

import (
	"ACC-HOMEWORK/course/notion"
	"ACC-HOMEWORK/crud"
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Course struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	RoomNumber string `json:"room_number"`
	Duration   string `json:"duration"`
	NotionID   string `json:"notion_id"`
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

func NewCourseFromMap(_course map[string]string) *Course {
	course := &Course{}
	course.NotionID = _course["notion_id"]
	course.Name = _course["name"]
	course.Code = _course["code"]
	course.RoomNumber = _course["room_number"]
	course.Duration = _course["duration"]
	return course
}

func (c *Course) ToMap() map[string]string {
	return map[string]string{
		"notion_id":   c.NotionID,
		"name":        c.Name,
		"code":        c.Code,
		"room_number": c.RoomNumber,
		"duration":    c.Duration,
	}
}

func (c *Course) Add(db *sql.DB) (err error) {

	course := c.ToMap()

	err = crud.PostHandler(course, "courses", db)

	if err != nil {
		log.Fatalln("Error adding course to database: ", err)
		return err
	}

	notion_id, err_notion := notion.AddCourseToNotion(course)

	if err_notion != nil {
		log.Fatalln("Error adding course to Notion: ", err_notion)
		return err_notion
	}

	var lastVal int
	err = db.QueryRow("SELECT MAX(id) FROM courses").Scan(&lastVal)
	if err != nil {
		log.Fatalln("Error getting course id: ", err)
		return err
	}

	err = crud.PutHanlder(lastVal+1, "notion_id", "courses", notion_id, db)

	if err != nil {
		log.Fatalln("Error updating course: ", err)
		return err
	}

	return nil
}

func (c *Course) Update(col, value string, db *sql.DB) (err error) {
	query := fmt.Sprintf("UPDATE courses SET %v = '%v' WHERE id = '%v'", col, value, c.NotionID)
	_, err = db.Exec(query)
	return err
}
