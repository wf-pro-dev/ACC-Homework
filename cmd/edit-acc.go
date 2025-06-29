package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/crud"
)

// func showUsage() {
// 	fmt.Println("Usage: edit-acc <ASSIGNMENT_ID> <COLUMN> <NEW_VALUE>")
// 	fmt.Println("")
// 	fmt.Println("Edit an existing assignment in the ACC Homework tracker.")
// 	fmt.Println("")
// 	fmt.Println("Arguments:")
// 	fmt.Println("  ASSIGNMENT_ID    The ID of the assignment to edit")
// 	fmt.Println("  COLUMN           The column/field to update")
// 	fmt.Println("  NEW_VALUE        The new value for the specified column")
// 	fmt.Println("")
// 	fmt.Println("Available columns:")
// 	fmt.Println("  title            Assignment title")
// 	fmt.Println("  todo             Assignment todo/description")
// 	fmt.Println("  deadline         Due date (format: yyyy-mm-dd)")
// 	fmt.Println("  type             Assignment type (HW, Exam, etc.)")
// 	fmt.Println("  course_code      Course code")
// 	fmt.Println("  link             Assignment link/URL")
// 	fmt.Println("")
// 	fmt.Println("Options:")
// 	fmt.Println("  -h, --help       Show this help message")
// 	fmt.Println("")
// 	fmt.Println("Examples:")
// 	fmt.Println("  edit-acc 5 title \"New Assignment Title\"")
// 	fmt.Println("  edit-acc 3 deadline 2024-12-15")
// 	fmt.Println("  edit-acc 1 todo \"Complete the final project\"")
// 	fmt.Println("  edit-acc 2 type HW")
// 	fmt.Println("  edit-acc 4 link \"https://example.com/assignment\"")
// 	fmt.Println("  edit-acc -h      # Show this help message")
// 	fmt.Println("")
// 	fmt.Println("Note: The assignment ID can be found using the ls-acc command.")
// }

func ValidateColumn(col string) error {
	for _, columns := range types.COLUMNS {
		if columns[0:2] == col || columns == col {
			return nil
		}
	}
	return fmt.Errorf("invalid column: %s", col)
}

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing assignment in the ACC Homework tracker.",
	Long:  `Edit an existing assignment in the ACC Homework tracker.`,
	Args: func(cmd *cobra.Command, args []string) error {

		db, DBerr := crud.GetDB()
		if DBerr != nil {
			return DBerr
		}

		var err error

		if len(args) != 3 {
			return fmt.Errorf("edit-acc requires exactly 3 arguments")
		}

		if err = ValidateAssignmentId(args[0], db); err != nil {
			return err
		}

		if err = ValidateColumn(args[1]); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		db, err := crud.GetDB()
		if err != nil {
			log.Fatal(err)
		}

		assignment_id := args[0]
		col := args[1]
		newValue := args[2]

		assignment := assignment.GetAssignmentsbyId(assignment_id, db)

		err = assignment.Update(col, newValue, db)
		if err != nil {
			log.Fatalln("Error updating assignment: ", err)
		}

	}}
