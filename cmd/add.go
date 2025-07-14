package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/client"
)

var assignmentFlag bool
var courseFlag bool

func init() {
	addCmd.Flags().BoolVarP(&assignmentFlag, "assignment", "a", false, "Add a new assignment")
	addCmd.Flags().BoolVarP(&courseFlag, "course", "c", false, "Add a new course")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new assignment or course to the ACC Homework tracker",
	Long:  `Add a new assignment or course through the ACC Homework API`,
	Run: func(cmd *cobra.Command, args []string) {
		if assignmentFlag {
			addAssignment()
		} else if courseFlag {
			addCourse()
		} else {
			fmt.Println("Please specify either --assignment (-a) or --course (-c)")
			cmd.Help()
			os.Exit(1)
		}
	},
}

func addAssignment() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== Create New Assignment ===")

	assignmentData := map[string]string{

		"course_code": promptInput(scanner, "Course Code"),
		"title":       promptInput(scanner, "Title"),
		"type_name":   promptInput(scanner, "Type (HW/Exam)"),
		"deadline":    promptDate(scanner, "Deadline (YYYY-MM-DD)"),
		"todo":        promptInput(scanner, "Todo/Description"),
	}

	assignmentData["status_name"] = "Not started"

	fmt.Println("\nCreating assignment...")
	_, err := client.CreateAssignment(assignmentData)
	if err != nil {
		log.Fatalf("Error creating assignment: %v", err)
	}

	fmt.Printf("\n✅ Assignment created successfully!\n")
}

func addCourse() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== Create New Course ===")

	courseData := map[string]string{

		"code":        promptInput(scanner, "Course Code"),
		"name":        promptInput(scanner, "Name"),
		"duration":    promptInput(scanner, "Duration (D,D HH:MM-HH:MM)"),
		"room_number": promptInput(scanner, "Room Number (e.g. 101 or Online)"),
	}

	fmt.Println("\nCreating course...")
	_, err := client.CreateCourse(courseData)
	if err != nil {
		log.Fatalf("Error creating course: %v", err)
	}

	fmt.Printf("\n✅ Course created successfully!\n")
}

func promptInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Printf("%s: ", prompt)
	scanner.Scan()
	return scanner.Text()
}

func promptDate(scanner *bufio.Scanner, prompt string) string {
	for {
		dateStr := promptInput(scanner, prompt)
		_, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			dateDB, err := convertToDBFormatWithLocation(dateStr, time.Local)
			if err == nil {
				return dateDB
			}
		}
		fmt.Println("Invalid date format. Please use YYYY-MM-DD")
	}
}

func convertToDBFormatWithLocation(dateStr string, loc *time.Location) (string, error) {
	t, err := time.ParseInLocation(time.DateOnly, dateStr, loc)
	if err != nil {
		return "", err
	}

	// Convert to RFC3339 format with the specified location
	return t.Format(time.RFC3339), nil
}
