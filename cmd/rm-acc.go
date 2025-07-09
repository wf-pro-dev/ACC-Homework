package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/database"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an assignment from the ACC Homework tracker.",
	Long:  `Remove an assignment from the ACC Homework tracker.`,
	Args: func(cmd *cobra.Command, args []string) error {
		db, err := database.GetDB()
		if err != nil {
			return fmt.Errorf("error getting database: %s", err)
		}
		if len(args) != 1 {
			return fmt.Errorf("rm-acc requires exactly 1 argument")
		}

		if err = ValidateAssignmentId(args[0], db); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.GetDB()
		if err != nil {
			log.Fatalln("Error getting database: ", err)
		}

		assigmnent_id := args[0]
		assignment := assignment.Get_Assignment_byId(assigmnent_id, db)

		assignment.Delete(db)
	}}
