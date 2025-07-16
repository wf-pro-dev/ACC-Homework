package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"github.com/williamfotso/acc/internal/storage/global"
)

type AssignmentData struct {
	ID         int
	CourseCode string
	Type       string
	Deadline   time.Time
	Title      string
	Todo       string
	NotionID   string
	Link       string
	Status     string
}

type CourseData struct {
	Name       string
	Code       string
	NotionID   string
	Duration   string
	RoomNumber string
}

func parseSQLFile(filePath string) ([]AssignmentData, []CourseData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	var assignments []AssignmentData
	var courses []CourseData
	var currentTable string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect table data sections
		if strings.HasPrefix(line, "COPY public.") {
			if strings.Contains(line, "assignments") {
				currentTable = "assignments"
			} else if strings.Contains(line, "courses") {
				currentTable = "courses"
			} else {
				currentTable = ""
			}
			continue
		}

		// Parse data lines
		if line == "\\." {
			currentTable = ""
			continue
		}

		switch currentTable {
		case "assignments":
			// Parse assignment data
			// Example line: 75	GOVT-2305	HW	2025-07-08 00:00:00	Self Intro	use this forum...	22940a21-a7e3-8169-8f22-f1d3e8e0f709	https://acconline.austincc.edu/ultra/stream	done
			parts := strings.Split(line, "\t")
			if len(parts) >= 9 {
				id := 0
				fmt.Sscanf(parts[0], "%d", &id)

				deadline, err := time.Parse("2006-01-02 15:04:05", parts[3])
				if err != nil {
					return nil, nil, fmt.Errorf("invalid deadline format: %w", err)
				}

				assignments = append(assignments, AssignmentData{
					ID:         id,
					CourseCode: parts[1],
					Type:       parts[2],
					Deadline:   deadline,
					Title:      parts[4],
					Todo:       parts[5],
					NotionID:   parts[6],
					Link:       parts[7],
					Status:     parts[8],
				})
			}

		case "courses":
			// Parse course data
			// Example line: Pre-Calculus	MATH-2412	17e40a21a7e380459b6fe9695d4edff9	\N	\N	1
			parts := strings.Split(line, "\t")
			if len(parts) >= 6 {
				courses = append(courses, CourseData{
					Name:       parts[0],
					Code:       parts[1],
					NotionID:   parts[2],
					Duration:   strings.Replace(parts[3], `\N`, "", -1),
					RoomNumber: strings.Replace(parts[4], `\N`, "", -1),
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error scanning SQL file: %w", err)
	}

	return assignments, courses, nil
}

func migrateData(db *gorm.DB, assignments []AssignmentData, courses []CourseData) error {
	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	currentTime := time.Now()
	userID := uint(1)

	// 1. Migrate status types
	statuses := []struct {
		ID       uint
		Name     string
		Color    string
		NotionID string
	}{
		{1, "Not started", "default", "3aa77cf8-c39e-4c7b-b7d2-ab15ae43ff23"},
		{2, "In progress", "blue", "97903420-1e83-4b3a-9eaf-a904354c968b"},
		{3, "Done", "green", "2fef8044-d8d7-4fcf-a3ee-393a1d558e94"},
	}

	for _, status := range statuses {
		if err := tx.Exec(`
			INSERT INTO public.assignment_statuses (id, name, color, notion_id)
			VALUES (?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			status.ID, status.Name, status.Color, status.NotionID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to migrate status types: %w", err)
		}
	}

	// 2. Migrate assignment types
	types := []struct {
		ID       uint
		Name     string
		Color    string
		NotionID string
	}{
		{1, "HW", "yellow", "Vn}Z"},
		{2, "Exam", "red", "oiNS"},
	}

	for _, t := range types {
		if err := tx.Exec(`
			INSERT INTO public.assignment_types (id, name, color, notion_id)
			VALUES (?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			t.ID, t.Name, t.Color, t.NotionID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to migrate assignment types: %w", err)
		}
	}

	// 3. Migrate courses
	for _, course := range courses {
		if err := tx.Exec(`
			INSERT INTO public.courses (created_at, updated_at, user_id, notion_id, code, name, duration, room_number)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (code) DO NOTHING`,
			currentTime, currentTime, userID, course.NotionID, course.Code, course.Name, course.Duration, course.RoomNumber).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert course: %w", err)
		}
	}

	// 4. Migrate assignments
	for _, assignment := range assignments {
		// Map status
		var statusName string
		switch assignment.Status {
		case "default":
			statusName = "Not started"
		case "start":
			statusName = "In progress"
		case "done":
			statusName = "Done"
		default:
			statusName = "Not started"
		}

		// Map type
		var typeName string
		switch assignment.Type {
		case "HW":
			typeName = "HW"
		case "Exam":
			typeName = "Exam"
		default:
			typeName = "HW"
		}

		if err := tx.Exec(`
			INSERT INTO public.assignments (
				created_at, updated_at, user_id, notion_id, title, todo, 
				deadline, link, course_code, type_name, status_name
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (notion_id) DO NOTHING`,
			currentTime, currentTime, userID, assignment.NotionID, assignment.Title, assignment.Todo,
			assignment.Deadline, assignment.Link, assignment.CourseCode, typeName, statusName).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert assignment: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}


func main() {
	// 1. Check command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path_to_sql_file>")
	}
	sqlFile := os.Args[1]
	
	db, err := global.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Parse the SQL file
	assignments, courses, err := parseSQLFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to parse SQL file: %v", err)
	}

	// 4. Perform the migration
	if err := migrateData(db, assignments, courses); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
