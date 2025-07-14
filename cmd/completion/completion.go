package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/storage/local"
	"github.com/williamfotso/acc/internal/types"
)

func CourseCodeCompletion() ([]string, cobra.ShellCompDirective) {

	userID, err := local.GetCurrentUserID()
	if err != nil {
		fmt.Println("Error getting current user ID:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		fmt.Println("Error getting local DB:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var courses []course.Course
	err = db.Raw("SELECT code, name FROM local_courses").Scan(&courses).Error
	if err != nil {
		fmt.Println("Error getting courses:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	course_codes := []string{}
	for _, course := range courses {
		completion := fmt.Sprintf("%s\t%s", course.Code, course.Name)
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
