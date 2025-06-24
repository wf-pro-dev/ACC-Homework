package assignment

import (
	"ACC-HOMEWORK/assignment/notion"
	"ACC-HOMEWORK/crud"
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Assignment represents a homework or exam assignment
type Assignment struct {
	ID         int    `db:"id,omitempty"`
	NotionID   string `db:"notion_id,omitempty"`
	Type       string
	Deadline   string
	Title      string
	Todo       string
	CourseCode string
	Link       string
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
	assignment.Type = scanner.Text()

	fmt.Printf("The deadline (yyyy-mm-dd): ")
	scanner.Scan()
	assignment.Deadline = scanner.Text()

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

	return assignment
}

func NewAssignmentFromMap(assign map[string]string) *Assignment {
	assignment := &Assignment{}
	assignment.NotionID = assign["notion_id"]
	assignment.Type = assign["type"]
	assignment.Deadline = assign["deadline"]
	assignment.Title = assign["title"]
	assignment.Todo = assign["todo"]
	assignment.CourseCode = assign["course_code"]

	return assignment
}

func (a *Assignment) SetNotionID(notionID string) {
	a.NotionID = notionID
}

func (a *Assignment) GetNotionID() string {
	return a.NotionID
}

// SetCourseCode allows setting the course code manually
func (a *Assignment) SetCourseCode(courseCode string) {
	a.CourseCode = courseCode
}

// GetType returns the assignment type
func (a *Assignment) GetType() string {
	return a.Type
}

// GetDeadline returns the assignment deadline
func (a *Assignment) GetDeadline() string {
	return a.Deadline
}

// GetTitle returns the assignment title
func (a *Assignment) GetTitle() string {
	return a.Title
}

// GetTodo returns the assignment todo
func (a *Assignment) GetTodo() string {
	return a.Todo
}

// GetCourseCode returns the assignment course code
func (a *Assignment) GetCourseCode() string {
	return a.CourseCode
}

func getCourse(course_code string, db *sql.DB) map[string]string {
	query := fmt.Sprintf("SELECT notion_id FROM courses WHERE code='%v'", course_code)
	fmt.Println(query)
	course, err := crud.GetHandler(query, db)
	if err != nil {
		panic(err)
	}
	return course[0]
}

func getType(type_name string, db *sql.DB) map[string]string {

	type_info, err := crud.GetHandler(fmt.Sprintf("SELECT * FROM type WHERE name='%v'", type_name), db)

	if err != nil {
		panic(err)
	}
	return type_info[0]
}

// ToMap converts the Assignment struct to a map[string]string
// This maintains compatibility with the existing database operations
func (a *Assignment) ToMap() map[string]string {
	fmt.Println(a.ID)
	return map[string]string{
		"id":          strconv.Itoa(a.ID),
		"notion_id":   a.NotionID,
		"type":        a.Type,
		"deadline":    a.Deadline,
		"title":       a.Title,
		"todo":        a.Todo,
		"course_code": a.CourseCode,
		"link":        a.Link,
	}
}

func (a *Assignment) Add(db *sql.DB) (err error) {

	assignment := a.ToMap()

	assignment["link"] = ""

	delete(assignment, "id")

	err = crud.PostHandler(assignment, "assignements", db)

	if err != nil {
		log.Fatalln(err)
	}

	notion_id, err_notion := notion.AddAssignmentToNotion(assignment, getType(assignment["type"], db), getCourse(assignment["course_code"], db))

	if err_notion != nil {
		log.Fatalln(err_notion)
	}

	var lastVal int
	err = db.QueryRow("SELECT MAX(id) FROM assignements").Scan(&lastVal)
	if err != nil {
		log.Fatal(err)
	}
	err = crud.PutHanlder(lastVal, "notion_id", "assignements", notion_id, db)

	if err == nil {
		fmt.Printf("\nSucceful new Assignement ! %#v", notion_id)
	} else {
		fmt.Println("Error updating assignment")
	}

	return err
}

func GetAssignmentsbyCourse(course_code string, columns []string, filters []Filter, up_to_date bool, db *sql.DB) {

	col_length := 15
	query := fmt.Sprintf("SELECT %s FROM assignements WHERE course_code='%v'", strings.Join(columns, ","), course_code)

	for _, filter := range filters {
		query += fmt.Sprintf(" AND %s='%v'", filter.Column, filter.Value)
	}

	if up_to_date {
		query += " AND deadline > NOW()"
	}
	query += " ORDER BY deadline ASC"
	assignments, err := crud.GetHandler(query, db)
	if err != nil {
		panic(err)
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
		fmt.Print("│")
		for _, col := range columns {
			value := assignment[col]
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

func GetAssignmentsbyId(id string, db *sql.DB) *Assignment {

	assignments, err := crud.GetHandler(fmt.Sprintf("SELECT * FROM assignements WHERE id='%v'", id), db)
	if err != nil {
		panic(err)
	}

	return NewAssignmentFromMap(assignments[0])
}

func (a *Assignment) getID(db *sql.DB) int {

	assignment, err := crud.GetHandler(fmt.Sprintf("SELECT id FROM assignements WHERE notion_id='%v'", a.NotionID), db)
	if err != nil {
		panic(err)
	}
	id := assignment[0]["id"]

	int_id, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	return int_id
}

func (a *Assignment) Update(col, value string, db *sql.DB) (err error) {

	fmt.Println("Updating assignment in database")

	err = crud.PutHanlder(a.getID(db), col, "assignements", value, db)

	if err != nil {
		log.Fatalln(err)
		return err
	}

	fmt.Println("Updating assignment to Notion")

	err = notion.UpdateAssignementToNotion(a.ToMap(), col, value)

	if err != nil {
		log.Fatalln(err)
		return err
	}

	fmt.Println("Successfully updated assignment")

	return nil
}

func (a *Assignment) Delete(db *sql.DB) (err error) {

	err = crud.DeleteHandler("assignements", "notion_id", a.NotionID, db)

	if err != nil {
		log.Fatalln(err)
	}

	err = notion.DeleteAssignementFromNotion(a.ToMap())
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
