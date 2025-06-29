package main

import (
	"ACC-HOMEWORK/assignment"
	"ACC-HOMEWORK/course"
	"ACC-HOMEWORK/crud"
	"fmt"
	"log"
	"os"
)

func showUsage() {
	fmt.Println("Usage: add-acc [OPTION]")
	fmt.Println("")
	fmt.Println("Add a new assignment or course to the ACC Homework tracker.")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -a, --assignment    Add a new assignment (default)")
	fmt.Println("  -c, --course        Add a new course")
	fmt.Println("  -h, --help          Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  add-acc             # Add a new assignment")
	fmt.Println("  add-acc -a          # Add a new assignment")
	fmt.Println("  add-acc -c          # Add a new course")
	fmt.Println("  add-acc -h          # Show this help message")
	fmt.Println("")
	fmt.Println("The program will prompt you for the necessary information")
	fmt.Println("to create the assignment or course entry.")
}

func main() {
	// Check if help is requested
	if len(os.Args) > 1 && (os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		showUsage()
		return
	}

	db, err := crud.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	var param_type string
	if len(os.Args) > 0 {
		param_type = os.Args[1]
	} else {
		param_type = "-a"
	}

	switch param_type {
	case "-a", "--assignment":
		assignment := assignment.NewAssignment()
		assignment.Add(db)
	case "-c", "--course":
		course := course.NewCourse()
		course.Add(db)
	default:
		showUsage()
		return
	}

}
