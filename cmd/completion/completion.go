package completion

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/williamfotso/acc/internal/core/models/course"

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/storage/local"
	"github.com/williamfotso/acc/internal/types"
)

func AssignmentIdCompletion() ([]string, cobra.ShellCompDirective) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	baseName := filepath.Base(wd)

	userID, err := local.GetCurrentUserID()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	db, err := local.GetLocalDB(userID)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	query := fmt.Sprintf("SELECT remote_id, title FROM local_assignments WHERE course_code = '%s' ORDER BY remote_id ASC", baseName)
	var assignments []assignment.LocalAssignment
	err = db.Raw(query).Scan(&assignments).Error
	if err != nil {
		fmt.Println("[ERROR] Query error:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var assignmentIDs []string
	for _, assignment := range assignments {
		completion := fmt.Sprintf("%d\t%s", assignment.RemoteID, assignment.Title)
		assignmentIDs = append(assignmentIDs, completion)
	}
	return assignmentIDs, cobra.ShellCompDirectiveNoFileComp
}

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
