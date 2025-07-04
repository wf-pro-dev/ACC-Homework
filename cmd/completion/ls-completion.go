package completion

import (
	"strings"

	"github.com/spf13/cobra"
)

func CompleteFilterFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	// Get parts of the filter flag
	parts := strings.Split(toComplete, "=")

	// Get all possible columns
	allColumns, _ := CompleteColumns()

	var completions []string
	for _, column := range allColumns {

		if len(parts) > 1 {
			switch parts[0] {
			case "course_code":
				course_codes, _ := CourseCodeCompletion()
				for _, course_code := range course_codes {
					completions = append(completions, column+"="+course_code)
				}
			case "type":
				for _, value := range []string{"HW", "Exam"} {
					completions = append(completions, column+"="+value)
				}
			case "status":
				for _, value := range []string{"done", "start", "default"} {
					completions = append(completions, column+"="+value)
				}
			default:
				completions = append(completions, column+"=")
			}
		} else {
			completions = append(completions, column+"=")
		}

	}

	return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp

}

// CompleteMultiColumn handles completion for flags that accept multiple columns
func CompleteMultiColumn(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	// Get already selected columns from the current flag value
	selected := strings.Split(toComplete, ",")

	// If there's a partial value at the end, we'll complete that
	lastPart := selected[len(selected)-1]

	// Get all possible columns
	allColumns, _ := CompleteColumns()

	// Filter out already selected columns (except the one we're completing)
	var available []string
	for _, col := range allColumns {
		alreadySelected := false
		for i, s := range selected {
			// Skip the last part (the one we're completing)
			if i < len(selected)-1 && s == col {
				alreadySelected = true
				break
			}
		}
		if !alreadySelected && strings.HasPrefix(col, lastPart) {
			// Build the completion suggestion
			if len(selected) > 1 {
				// If we have multiple values, suggest the full chain
				suggestion := strings.Join(selected[:len(selected)-1], ",") + "," + col
				available = append(available, suggestion)
			} else {
				available = append(available, col)
			}
		}
	}

	return available, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
