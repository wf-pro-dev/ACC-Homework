package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"ACC-HOMEWORK/crud"
	"ACC-HOMEWORK/notion"
)

func createAssign() map[string]string {

	newAssign := map[string]string{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("The type (HW or Exam): ")
	scanner.Scan()
	newAssign["type"] = scanner.Text()

	fmt.Printf("The deadline (yyyy-mm-dd): ")
	scanner.Scan()
	newAssign["deadline"] = scanner.Text()

	fmt.Printf("The title: ")
	scanner.Scan()
	newAssign["title"] = scanner.Text()

	fmt.Printf("The todo: ")
	scanner.Scan()
	newAssign["todo"] = scanner.Text()

	return newAssign

}

func getCourse(course_code string, db *sql.DB) map[string]string {

	course, err := crud.GetHandler(fmt.Sprintf("SELECT id FROM courses WHERE code='%v'", course_code), db)
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

func main() {

	db, err_conn := crud.GetDB()

	if err_conn != nil {
		log.Fatalln(err_conn)
	}

	fmt.Println("===== Creating new Assignement =====")
	newAssign := createAssign()

	pwd := os.Getenv("PWD")
	cmd := exec.Command("basename", pwd)

	// Capture the output
	output, _ := cmd.CombinedOutput()

	newAssign["course_code"] = strings.TrimSpace(string(output))

	err_query := crud.PostHandler(newAssign, "assignements", db)

	if err_query != nil {
		log.Fatalln(err_query)
	}

	notion_id, err_notion := notion.AddAssignmentToNotion(newAssign, getType(newAssign["type"], db), getCourse(newAssign["course_code"], db))

	if err_notion != nil {
		log.Fatalln(err_notion)
	}

	var lastVal int
	err := db.QueryRow("SELECT last_value FROM assignements_id_seq").Scan(&lastVal)
	if err != nil {
		log.Fatal(err)
	}
	crud.PutHanlder(lastVal, "notion_id", notion_id, db)
	fmt.Printf("\nSucceful new Assignement ! %#v", notion_id)

}
