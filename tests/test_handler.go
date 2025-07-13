package main

import (
	"os"
	"fmt"

	"github.com/williamfotso/acc/internal/services/client"
)


func main() {
	
	assignments, err := client.GetAssignments()
	if err != nil {
		fmt.Printf("ERROR : %s",err)
		os.Exit(1)
	}

	for _,a := range assignments {
		fmt.Printf("title : %s, course code : %s, status : %s\n",a["title"],a["course_code"],a["status"])
	}

	
	courses, err := client.GetCourses()
	if err != nil {
		fmt.Printf("ERROR : %s",err)
		os.Exit(1)
	}

	fmt.Println("Courses :")
	for _,c := range courses {
		fmt.Printf("name : %s, course code : %s\n",c["name"],c["code"])
	}


}
