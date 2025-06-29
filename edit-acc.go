package main

import (
	"ACC-HOMEWORK/assignment"
	"ACC-HOMEWORK/crud"
	"fmt"
	"log"
	"os"
)

func showUsage() {
	fmt.Println("Usage: edit-acc <ASSIGNMENT_ID> <COLUMN> <NEW_VALUE>")
	fmt.Println("")
	fmt.Println("Edit an existing assignment in the ACC Homework tracker.")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  ASSIGNMENT_ID    The ID of the assignment to edit")
	fmt.Println("  COLUMN           The column/field to update")
	fmt.Println("  NEW_VALUE        The new value for the specified column")
	fmt.Println("")
	fmt.Println("Available columns:")
	fmt.Println("  title            Assignment title")
	fmt.Println("  todo             Assignment todo/description")
	fmt.Println("  deadline         Due date (format: yyyy-mm-dd)")
	fmt.Println("  type             Assignment type (HW, Exam, etc.)")
	fmt.Println("  course_code      Course code")
	fmt.Println("  link             Assignment link/URL")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  edit-acc 5 title \"New Assignment Title\"")
	fmt.Println("  edit-acc 3 deadline 2024-12-15")
	fmt.Println("  edit-acc 1 todo \"Complete the final project\"")
	fmt.Println("  edit-acc 2 type HW")
	fmt.Println("  edit-acc 4 link \"https://example.com/assignment\"")
	fmt.Println("  edit-acc -h      # Show this help message")
	fmt.Println("")
	fmt.Println("Note: The assignment ID can be found using the ls-acc command.")
}

func main() {
	// Check if help is requested
	if len(os.Args) > 1 && (os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		showUsage()
		return
	}

	// Check if correct number of arguments is provided
	if len(os.Args) != 4 {
		fmt.Println("Error: Invalid number of arguments")
		fmt.Println("")
		showUsage()
		os.Exit(1)
	}

	db, err := crud.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	assignment_id := os.Args[1]
	col := os.Args[2]
	newValue := os.Args[3]

	// Validate column name
	validColumns := []string{"title", "todo", "deadline", "type", "course_code", "link"}
	isValidColumn := false
	for _, validCol := range validColumns {
		if col == validCol {
			isValidColumn = true
			break
		}
	}

	if !isValidColumn {
		fmt.Printf("Error: Invalid column '%s'\n", col)
		fmt.Println("")
		fmt.Println("Available columns:")
		for _, validCol := range validColumns {
			fmt.Printf("  %s\n", validCol)
		}
		fmt.Println("")
		showUsage()
		os.Exit(1)
	}

	assignment := assignment.GetAssignmentsbyId(assignment_id, db)

	err = assignment.Update(col, newValue, db)
	if err != nil {
		log.Fatalln("Error updating assignment: ", err)
	}

}
