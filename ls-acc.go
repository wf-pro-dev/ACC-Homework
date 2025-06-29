package main

import (
	"ACC-HOMEWORK/assignment"
	"ACC-HOMEWORK/assignment/notion/types"
	"ACC-HOMEWORK/crud"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getColumns(arg string) (col string) {
	for _, column := range types.COLUMNS {
		if column[0:2] == arg {
			col = column
		}
	}

	if col == "" {
		ColumnError(fmt.Sprintf("Invalid argument: -%s\n", arg))
	}

	return col
}

func ColumnError(message string) {

	fmt.Println(message)
	fmt.Println("Available columns:\n")
	for _, column := range types.COLUMNS {
		fmt.Printf("  -%s (%s)\n", column[0:2], column)
	}
	os.Exit(1)
}

func CourseError(message string, COURSES_CODES []map[string]string) {
	fmt.Println(message)
	fmt.Println("Available courses:\n")
	for _, course := range COURSES_CODES {
		fmt.Printf("  %s\n", course["code"])
	}
	os.Exit(1)
}

func handleFlag(arg string) (columns []string, filters []assignment.Filter) {

	if len(arg) < 3 {
		ColumnError(fmt.Sprintf("Invalid argument: %s\n", arg))
	}

	if len(arg) > 3 && arg[3] == '=' {
		filters = append(filters, assignment.Filter{Column: getColumns(arg[1:3]), Value: arg[4:]})
	} else {
		columns = append(columns, getColumns(arg[1:3]))
	}

	return columns, filters

}

func main() {

	db, err := crud.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	baseName := filepath.Base(wd)
	courseName := baseName
	columns := []string{}
	filters := []assignment.Filter{}
	up_to_date := false

	for _, arg := range os.Args[1:] {
		switch arg[0] {
		case '-':
			switch arg {
			case "-d":
				up_to_date = true
			default:
				col, fil := handleFlag(arg)
				columns = append(columns, col...)
				filters = append(filters, fil...)
			}

		default:
			courseName = arg
		}
	}

	if len(columns) == 0 {
		columns = types.DEFAULT_COLUMNS_FOR_LS
	}

	COURSES_CODES, err := crud.GetHandler("SELECT code FROM courses", db)

	validCourse := false
	for _, course := range COURSES_CODES {
		if course["code"] == courseName {
			validCourse = true
		}
	}

	if !validCourse {
		CourseError(fmt.Sprintf("Invalid course code: %s\n", courseName), COURSES_CODES)
	}

	if err != nil {
		log.Fatal(err)
	}

	assignment.GetAssignmentsbyCourse(courseName, columns, filters, up_to_date, db)

}
