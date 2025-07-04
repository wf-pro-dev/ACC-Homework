package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/cmd/completion"
	"github.com/williamfotso/acc/crud"
)

// splitCommaSeparated splits a string by commas and trims whitespace
func splitCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getColumns(arg string) (col string) {
	for _, column := range types.COLUMNS {
		if column[0:2] == arg || column == arg {
			col = column
		}
	}

	if col == "" {
		ColumnError(fmt.Sprintf("Invalid argument: -%s\n", arg))
	}

	return col
}

func CourseError(message string, COURSES_CODES []map[string]string) {
	fmt.Println(message)
	fmt.Println("Available courses:")
	for _, course := range COURSES_CODES {
		fmt.Printf("  %s\n", course["code"])
	}
	os.Exit(1)
}

func handleFlag(arg string) (columns []string, filters []assignment.Filter) {

	if len(arg) < 3 {
		ColumnError(fmt.Sprintf("Invalid argument: %s\n", arg))
	}

	if len(arg) > 3 && arg[3] == '=' {
		filters = append(filters, assignment.Filter{Column: getColumns(arg[1:3]), Value: arg[4:]})
	} else {
		columns = append(columns, getColumns(arg[1:3]))
	}

	return columns, filters

}

var courseName string
var filter string
var up_to_date bool
var include []string
var exclude []string

func init() {
	// Handle --course -c flag
	lsCmd.Flags().StringVarP(&courseName, "course", "c", "", "Course to list assignments for")
	_ = lsCmd.RegisterFlagCompletionFunc("course", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completion.CourseCodeCompletion()
	})

	lsCmd.Flags().BoolP("up-to-date", "d", false, "List only assignments that are up to date")

	// Handle --filter -f flag
	lsCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter assignments by a specific column and value")
	_ = lsCmd.RegisterFlagCompletionFunc("filter", completion.CompleteFilterFlag)

	// Handle --include -i flag
	lsCmd.Flags().StringArrayVarP(&include, "include", "i", []string{}, "Include columns to display")
	_ = lsCmd.RegisterFlagCompletionFunc("include", completion.CompleteMultiColumn)

	// Handle --exclude -e flag
	lsCmd.Flags().StringArrayVarP(&exclude, "exclude", "e", []string{}, "Exclude columns to display")
	_ = lsCmd.RegisterFlagCompletionFunc("exclude", completion.CompleteMultiColumn)

	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all assignments for a course",
	Long:  `List all assignments for a course`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get database connection
		db, err := crud.GetDB()
		if err != nil {
			log.Fatal(err)
		}

		// Get the courses codes from the database
		COURSES_CODES, err := crud.GetHandler("SELECT code FROM courses", db)
		if err != nil {
			log.Fatal(err)
		}

		// Get the current working directory to get the course name
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Get the base name of the current working directory
		baseName := filepath.Base(wd)

		if courseName == "" {
			courseName = baseName
		}
		// Check if the course is valid
		validCourse := false
		for _, course := range COURSES_CODES {
			if course["code"] == courseName {
				validCourse = true
			}
		}

		if !validCourse {
			CourseError(fmt.Sprintf("Invalid course code: %s\n", courseName), COURSES_CODES)
		}

		// Get the up-to-date flag
		up_to_date, err := cmd.Flags().GetBool("up-to-date")
		if err != nil {
			log.Fatal(err)
		}

		// Handle the filter flag
		var filters []assignment.Filter
		if filter != "" {
			filterParts := strings.Split(filter, "=")
			if len(filterParts) != 2 {
				log.Fatal("Invalid filter format. Use -f <column>=<value>")
			}
			filters = append(filters, assignment.Filter{Column: getColumns(filterParts[0]), Value: filterParts[1]})
		}

		// Initialize the columns and filters
		var columns []string

		if len(include) > 0 {
			// Process include flags - split comma-separated values
			for _, includeFlag := range include {
				columns = append(columns, splitCommaSeparated(includeFlag)...)
			}
		} else {
			// Add the default columns to the columns slice
			columns = types.DEFAULT_COLUMNS_FOR_LS
		}

		if len(exclude) > 0 {
			// Process exclude flags - split comma-separated values
			var excludeColumns []string
			for _, excludeFlag := range exclude {
				excludeColumns = append(excludeColumns, splitCommaSeparated(excludeFlag)...)
			}

			// Remove the excluded columns from the columns slice
			for _, column := range excludeColumns {
				index := slices.Index(columns, column)
				if index != -1 {
					columns = slices.Delete(columns, index, index+1)
				}
			}
		}

		assignment.GetAssignmentsbyCourse(courseName, columns, filters, up_to_date, db)

	},
}
