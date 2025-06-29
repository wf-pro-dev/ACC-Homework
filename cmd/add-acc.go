package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/course"
	"github.com/williamfotso/acc/crud"
)

var assignment_flag bool
var course_flag bool

func init() {
	addCmd.Flags().BoolVarP(&assignment_flag, "assignment", "a", false, "Add a new assignment")
	addCmd.Flags().BoolVarP(&course_flag, "course", "c", false, "Add a new course")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new assignment or course to the ACC Homework tracker.",
	Long:  `Add a new assignment or course to the ACC Homework tracker.`,
	Run: func(cmd *cobra.Command, args []string) {

		db, err := crud.GetDB()
		if err != nil {
			log.Fatal(err)
		}

		if assignment_flag {
			assignment := assignment.NewAssignment()
			assignment.Add(db)
		} else if course_flag {
			course := course.NewCourse()
			course.Add(db)
		}
	},
}
