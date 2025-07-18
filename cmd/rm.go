package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/cmd/completion"
	"github.com/williamfotso/acc/internal/services/client"
	"github.com/williamfotso/acc/internal/storage/local"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [assignment]",
	Short: "Remove an assignment from the ACC Homework tracker.",
	Long:  `Remove an assignment from the ACC Homework tracker.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completion.AssignmentIdCompletion()
	},
	Args: func(cmd *cobra.Command, args []string) error {
		userID, err := local.GetCurrentUserID()
		if err != nil {
			return err
		}
		db, err := local.GetLocalDB(userID)
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

		assigmnent_id := args[0]
		err := client.DeleteAssignment(assigmnent_id)
		if err != nil {
			log.Fatalln("Error deleting assignment: ", err)
		}
	}}
