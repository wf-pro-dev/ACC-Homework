package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/database"
)

func ColumnValueCompletion(args []string) ([]string, cobra.ShellCompDirective) {

	column := args[1]
	switch column {
	case "course_code":
		return CourseCodeCompletion()
	case "deadline":
		assignment_id := args[0]
		return DeadlineCompletion(assignment_id)
	case "type":
		return []string{"HW", "Exam"}, cobra.ShellCompDirectiveNoFileComp
	case "status":
		return []string{"done", "start", "default"}, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveError
	}
}

func DeadlineCompletion(assignment_id string) ([]string, cobra.ShellCompDirective) {

	db, err := database.GetDB()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	assignment := assignment.GetAssignmentsbyId(assignment_id, db)

	deadline_date, err := time.Parse(time.DateOnly, assignment.Deadline[:10])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	a_day_before_date := deadline_date.AddDate(0, 0, -1)
	a_day_before := a_day_before_date.Format(time.DateOnly)

	a_day_after_date := deadline_date.AddDate(0, 0, 1)
	a_day_after := a_day_after_date.Format(time.DateOnly)

	next_week := deadline_date.AddDate(0, 0, 7).Format(time.DateOnly)

	deadline_n_day := int(deadline_date.Weekday())
	diff := 5 - deadline_n_day
	if diff <= 0 {
		diff = 7 + diff
	}

	next_friday := deadline_date.AddDate(0, 0, diff).Format(time.DateOnly)

	return []string{
		fmt.Sprintf("%s\t%s", a_day_before, fmt.Sprintf("%s, a day before", a_day_before_date.Weekday())),
		fmt.Sprintf("%s\t%s", a_day_after, fmt.Sprintf("%s, a day after", a_day_after_date.Weekday())),
		fmt.Sprintf("%s\t%s", next_friday, "Next Friday"),
		fmt.Sprintf("%s\t%s", next_week, "A week from deadline"),
	}, cobra.ShellCompDirectiveNoFileComp

}

func EditCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	db, err := database.GetDB()
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
		query := fmt.Sprintf("SELECT id, title FROM assignements WHERE course_code = '%s' ORDER BY id ASC", baseName)
		assignments, err := database.GetHandler(query, db)
		if err != nil {
			fmt.Println("[ERROR] Query error:", err)
			return nil, cobra.ShellCompDirectiveError
		}

		var assignmentIDs []string
		for _, assignment := range assignments {
			completion := fmt.Sprintf("%s\t%s", assignment["id"], assignment["title"])
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
