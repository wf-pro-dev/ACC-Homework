package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/crud"
)

func CourseCodeCompletion() ([]string, cobra.ShellCompDirective) {

	db, err := crud.GetDB()
	if err != nil {
		fmt.Println("DEBUG: DB error:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	courses, err := crud.GetHandler("SELECT code,name FROM courses", db)
	if err != nil {
		fmt.Println("DEBUG: Query error:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	course_codes := []string{}
	for _, course := range courses {
		completion := fmt.Sprintf("%s\t%s", course["code"], course["name"])
		course_codes = append(course_codes, completion)
	}

	return course_codes, cobra.ShellCompDirectiveNoFileComp
}

// CompleteColumns provides completion for column names
func CompleteColumns() ([]string, cobra.ShellCompDirective) {
	var completions []string
	for _, column := range types.COLUMNS {
		completions = append(completions, column)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
