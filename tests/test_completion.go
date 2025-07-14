package main

import (
	"fmt"

	"github.com/williamfotso/acc/cmd/completion"
)

func main() {
	course_codes, _ := completion.CourseCodeCompletion()
	fmt.Println(course_codes)
}
