package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/storage/local"
	"github.com/williamfotso/acc/internal/types"
)

func ColumnValueCompletion(args []string) ([]string, cobra.ShellCompDirective) {

	column := args[1]
	switch column {
	case "course_code":
		return CourseCodeCompletion()
	case "deadline":
		assignment_id := args[0]
		return DeadlineCompletion(assignment_id)
	case "type_name":
		return []string{"HW", "Exam"}, cobra.ShellCompDirectiveNoFileComp
	case "status_name":
		return []string{"Not started", "In progress", "Done"}, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveError
	}
}

func DeadlineCompletion(assignment_id string) ([]string, cobra.ShellCompDirective) {

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

	assignment_id_int, err := strconv.Atoi(assignment_id)
	if err != nil {
		fmt.Println("Error converting assignment ID to int:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var assignment assignment.LocalAssignment
	err = db.First(&assignment, "remote_id = ?", assignment_id_int).Error
	if err != nil {
		fmt.Println("Error getting assignment:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	a_day_before_date := assignment.Deadline.AddDate(0, 0, -1)
	a_day_before := a_day_before_date.Format(time.DateOnly)

	a_day_after_date := assignment.Deadline.AddDate(0, 0, 1)
	a_day_after := a_day_after_date.Format(time.DateOnly)

	next_week := assignment.Deadline.AddDate(0, 0, 7).Format(time.DateOnly)

	deadline_n_day := int(assignment.Deadline.Weekday())
	diff := 5 - deadline_n_day
	if diff <= 0 {
		diff = 7 + diff
	}

	next_friday := assignment.Deadline.AddDate(0, 0, diff).Format(time.DateOnly)

	return []string{
		fmt.Sprintf("%s\t%s", a_day_before, fmt.Sprintf("%s, a day before", a_day_before_date.Weekday())),
		fmt.Sprintf("%s\t%s", a_day_after, fmt.Sprintf("%s, a day after", a_day_after_date.Weekday())),
		fmt.Sprintf("%s\t%s", next_friday, "Next Friday"),
		fmt.Sprintf("%s\t%s", next_week, "A week from deadline"),
	}, cobra.ShellCompDirectiveNoFileComp

}

func EditCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	userID, err := local.GetCurrentUserID()
	if err != nil {
		fmt.Println("Error getting current user ID:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		fmt.Println("[ERROR] DB connection error:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	switch len(args) {
	case 0:
		// First argument: assignment IDs
		wd, err := os.Getwd()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		baseName := filepath.Base(wd)
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

	case 1:
		// Second argument: columns
		var columns []string
		for _, column := range types.COLUMNS {
			columns = append(columns, column)
		}
		return columns, cobra.ShellCompDirectiveNoFileComp

	case 2:
		// Third argument: column-specific values
		return ColumnValueCompletion(args)

	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}
