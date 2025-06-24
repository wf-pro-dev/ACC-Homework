package main

import (
	"ACC-HOMEWORK/assignment"
	"ACC-HOMEWORK/course"
	"ACC-HOMEWORK/crud"
	"log"
	"os"
)

func main() {

	db, err := crud.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	param_type := os.Args[1]

	switch param_type {
	case "-a":
		assignment := assignment.NewAssignment()
		assignment.Add(db)
	case "-c":
		course := course.NewCourse()
		course.Add(db)
	default:
		assignment := assignment.NewAssignment()
		assignment.Add(db)
	}

}
