package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion/types"
)

func ValidateAssignmentId(id string, db *sql.DB) error {

	if id == "" {
		return fmt.Errorf("assignment ID is required")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("failed to convert assignment ID to int: %s", err)
	}

	assignment := assignment.GetAssignmentsbyId(id, db)
	if assignment == nil {
		return fmt.Errorf("assignment not found")
	}

	return nil
}

func ColumnError(message string) {

	fmt.Println(message)
	fmt.Println("Available columns:")
	for _, column := range types.COLUMNS {
		fmt.Printf("  -%s (%s)\n", column[0:2], column)
	}
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "acc",
	Short: "ACC is a CLI tool for managing assignments and courses",
	Long: `ACC is a CLI tool for managing assignments and courses from austin community college.
				  Complete documentation is available at https://github.com/williamfotso/acc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
