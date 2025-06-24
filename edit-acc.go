package main

import (
	"ACC-HOMEWORK/assignment"
	"ACC-HOMEWORK/crud"
	"log"
	"os"
)

func main() {

	db, err := crud.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	assigmnent_id := os.Args[1]
	assignment := assignment.GetAssignmentsbyId(assigmnent_id, db)

	col := os.Args[2]
	newValue := os.Args[3]

	assignment.Update(col, newValue, db)
}
