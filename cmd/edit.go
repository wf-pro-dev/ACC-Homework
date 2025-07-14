package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/cmd/completion"
	"github.com/williamfotso/acc/internal/services/client"
	"github.com/williamfotso/acc/internal/types"
	//"github.com/williamfotso/acc/cmd/completion"
)

func ValidateColumn(col string) error {

	for _, column := range types.COLUMNS {
		if column[0:2] == col || column == col {
			return nil
		}
	}
	return fmt.Errorf("invalid column: %s", col)
}

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:               "edit",
	Short:             "Edit an existing assignment in the ACC Homework tracker.",
	Long:              `Edit an existing assignment in the ACC Homework tracker.`,
	ValidArgsFunction: completion.EditCompletion,
	Args: func(cmd *cobra.Command, args []string) error {

		/*db, DBerr := database.GetDB()
		if DBerr != nil {
			return DBerr
		}*/

		var err error

		if len(args) != 3 {
			return fmt.Errorf("edit-acc requires exactly 3 arguments")
		}

		/*if err = ValidateAssignmentId(args[0], db); err != nil {
			return err
		}*/

		if err = ValidateColumn(args[1]); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		assignment_id := args[0]
		col := args[1]
		newValue := args[2]

		err := client.UpdateAssignment(assignment_id, col, newValue)
		if err != nil {
			log.Fatalln("Error updating assignment: ", err)
		}

	}}
